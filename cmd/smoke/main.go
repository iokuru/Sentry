package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type StepTimer struct {
	start time.Time
}

func startStep() *StepTimer {
	return &StepTimer{start: time.Now()}
}

func (s *StepTimer) duration() string {
	d := time.Since(s.start)
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func pass(step, timing string) {
	fmt.Printf("[PASS] %-32s (%s)\n", step, timing)
}

func fail(step string, err error) {
	fmt.Printf("[FAIL] %-32s: %v\n", step, err)
	os.Exit(1)
}

func authenticateCognito(region, clientID, email, password string) (string, error) {
	endpoint := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/", region)
	body := map[string]interface{}{
		"AuthFlow": "USER_PASSWORD_AUTH",
		"ClientId": clientID,
		"AuthParameters": map[string]string{
			"USERNAME": email,
			"PASSWORD": password,
		},
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/x-amz-json-1.1")
	req.Header.Set("X-Amz-Target", "AWSCognitoIdentityProviderService.InitiateAuth")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		rb, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("cognito returned status %d: %s", res.StatusCode, string(rb))
	}

	var data struct {
		AuthenticationResult struct {
			IdToken string `json:"IdToken"`
		} `json:"AuthenticationResult"`
	}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return "", err
	}
	return data.AuthenticationResult.IdToken, nil
}

func main() {
	liveFlag := flag.Bool("live", false, "run smoke tests against live aws deployment")
	flag.Parse()

	var apiBase string
	var authToken string

	if *liveFlag {
		fmt.Println("Running smoke tests against LIVE AWS environment...")
		apiBase = os.Getenv("API_URL")
		if apiBase == "" {
			cmd := exec.Command("terraform", "-chdir=infra", "output", "-raw", "api_url")
			out, err := cmd.Output()
			if err == nil {
				apiBase = strings.TrimSpace(string(out))
			}
		}
		if apiBase == "" {
			fail("resolve live api_url", fmt.Errorf("API_URL environment variable or terraform output is required"))
		}

		userPoolClientID := os.Getenv("USER_POOL_CLIENT_ID")
		if userPoolClientID == "" {
			cmd := exec.Command("terraform", "-chdir=infra", "output", "-raw", "user_pool_client_id")
			out, err := cmd.Output()
			if err == nil {
				userPoolClientID = strings.TrimSpace(string(out))
			}
		}

		region := os.Getenv("AWS_REGION")
		if region == "" {
			cmd := exec.Command("terraform", "-chdir=infra", "output", "-raw", "region")
			out, err := cmd.Output()
			if err == nil {
				region = strings.TrimSpace(string(out))
			}
		}
		if region == "" {
			region = "us-east-1"
		}

		caregiverEmail := os.Getenv("CAREGIVER_EMAIL")
		caregiverPassword := os.Getenv("CAREGIVER_PASSWORD")

		if caregiverEmail == "" {
			if inData, err := os.ReadFile("inputs.local.json"); err == nil {
				var in struct {
					CaregiverEmail string `json:"caregiver_email"`
				}
				_ = json.Unmarshal(inData, &in)
				caregiverEmail = in.CaregiverEmail
			}
		}
		if caregiverPassword == "" {
			caregiverPassword = "SentrySafety2026!"
		}

		if userPoolClientID != "" && caregiverEmail != "" && caregiverPassword != "" {
			token, err := authenticateCognito(region, userPoolClientID, caregiverEmail, caregiverPassword)
			if err == nil && token != "" {
				authToken = token
			}
		}
	} else {
		fmt.Println("Running smoke tests against LOCAL devserver (http://localhost:8080) ...")
		apiBase = "http://localhost:8080"
		authToken = "dev-id-token"
	}

	apiBase = strings.TrimSuffix(apiBase, "/")
	client := &http.Client{Timeout: 15 * time.Second}

	// 1. presign upload URL
	t1 := startStep()
	presignPayload := map[string]string{
		"content_type": "image/png",
		"filename":     "bank_scam.png",
	}
	body, _ := json.Marshal(presignPayload)
	req1, _ := http.NewRequest(http.MethodPost, apiBase+"/api/presign", bytes.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	resp1, err := client.Do(req1)
	if err != nil {
		fail("presign upload url", err)
	}
	defer resp1.Body.Close()
	if resp1.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp1.Body)
		fail("presign upload url", fmt.Errorf("status %d: %s", resp1.StatusCode, string(b)))
	}

	var presignData struct {
		UploadURL string `json:"upload_url"`
		ImageKey  string `json:"image_key"`
	}
	if err := json.NewDecoder(resp1.Body).Decode(&presignData); err != nil {
		fail("decode presign response", err)
	}
	if presignData.UploadURL == "" || presignData.ImageKey == "" {
		fail("presign response validation", fmt.Errorf("missing upload_url or image_key"))
	}
	pass("presign upload url", t1.duration())

	// 2. PUT test scam screenshot
	t2 := startStep()
	uploadURL := presignData.UploadURL
	if !strings.HasPrefix(uploadURL, "http://") && !strings.HasPrefix(uploadURL, "https://") {
		uploadURL = apiBase + uploadURL
	}

	var imageBytes []byte
	if data, readErr := os.ReadFile("fixtures/bank_scam.png"); readErr == nil {
		imageBytes = data
	} else {
		// minimal 1x1 png fallback
		imageBytes = []byte{
			0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
			0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
			0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
			0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
			0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
			0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
			0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00,
			0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
			0x42, 0x60, 0x82,
		}
	}

	putReq, _ := http.NewRequest(http.MethodPut, uploadURL, bytes.NewReader(imageBytes))
	putReq.Header.Set("Content-Type", "image/png")
	putResp, err := client.Do(putReq)
	if err != nil {
		fail("upload screenshot", err)
	}
	defer putResp.Body.Close()
	if putResp.StatusCode < 200 || putResp.StatusCode >= 300 {
		fail("upload screenshot", fmt.Errorf("status %d", putResp.StatusCode))
	}
	pass("upload screenshot", t2.duration())

	// 3. POST intake
	t3 := startStep()
	intakePayload := map[string]string{
		"senior_id":    "smoke-senior-test-device",
		"image_key":    presignData.ImageKey,
		"user_context": "Chase bank alert debit card locked verify pin",
		"country":      "US",
	}
	body3, _ := json.Marshal(intakePayload)
	req3, _ := http.NewRequest(http.MethodPost, apiBase+"/api/intake", bytes.NewReader(body3))
	req3.Header.Set("Content-Type", "application/json")
	resp3, err := client.Do(req3)
	if err != nil {
		fail("submit intake", err)
	}
	defer resp3.Body.Close()
	if resp3.StatusCode != http.StatusAccepted && resp3.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp3.Body)
		fail("submit intake", fmt.Errorf("status %d: %s", resp3.StatusCode, string(b)))
	}

	var intakeData struct {
		IncidentID string `json:"incident_id"`
		Status     string `json:"status"`
	}
	if err := json.NewDecoder(resp3.Body).Decode(&intakeData); err != nil {
		fail("decode intake response", err)
	}
	if intakeData.IncidentID == "" {
		fail("intake response validation", fmt.Errorf("missing incident_id"))
	}
	pass("submit intake", t3.duration())

	// 4. Poll incident result
	t4 := startStep()
	var incidentDetails struct {
		ID        string `json:"id"`
		Status    string `json:"status"`
		RiskLevel string `json:"risk_level"`
		Category  string `json:"category"`
		Output    struct {
			HoldOff         string `json:"hold_off"`
			SafeAlternative string `json:"safe_alternative"`
			Hotline         struct {
				Name   string `json:"name"`
				Number string `json:"number"`
			} `json:"hotline"`
		} `json:"output"`
	}

	analyzed := false
	for attempt := 0; attempt < 25; attempt++ {
		time.Sleep(1 * time.Second)
		pollResp, pollErr := client.Get(fmt.Sprintf("%s/api/incidents?id=%s", apiBase, intakeData.IncidentID))
		if pollErr != nil {
			continue
		}
		if pollResp.StatusCode == http.StatusOK {
			_ = json.NewDecoder(pollResp.Body).Decode(&incidentDetails)
			pollResp.Body.Close()
			if incidentDetails.Output.HoldOff != "" &&
				incidentDetails.Output.SafeAlternative != "" &&
				incidentDetails.Output.Hotline.Name != "" &&
				incidentDetails.Output.Hotline.Number != "" {
				analyzed = true
				break
			}
		} else {
			pollResp.Body.Close()
		}
	}

	if !analyzed {
		fail("poll analysis result", fmt.Errorf("incident analysis timed out or missing three outputs"))
	}
	pass("poll analysis result", t4.duration())

	// 5. List caregiver incidents
	t5 := startStep()
	listReq, _ := http.NewRequest(http.MethodGet, apiBase+"/api/caregiver/incidents", nil)
	if authToken != "" {
		listReq.Header.Set("Authorization", "Bearer "+authToken)
	}
	listResp, err := client.Do(listReq)
	if err != nil {
		fail("list caregiver incidents", err)
	}
	defer listResp.Body.Close()
	if listResp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(listResp.Body)
		fail("list caregiver incidents", fmt.Errorf("status %d: %s", listResp.StatusCode, string(b)))
	}

	var incidentsList []struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(listResp.Body).Decode(&incidentsList); err != nil {
		fail("decode caregiver list", err)
	}
	if len(incidentsList) == 0 {
		fail("assert incident in list", fmt.Errorf("caregiver incidents list is empty"))
	}

	found := false
	for _, item := range incidentsList {
		if item.ID == intakeData.IncidentID {
			found = true
			break
		}
	}
	if !found {
		fail("assert incident in list", fmt.Errorf("incident %s not found in caregiver feed", intakeData.IncidentID))
	}
	pass("list caregiver incidents", t5.duration())

	// 6. Acknowledge incident
	t6 := startStep()
	ackPayload := map[string]string{
		"incident_id": intakeData.IncidentID,
		"status":      "acknowledged",
	}
	body6, _ := json.Marshal(ackPayload)
	ackReq, _ := http.NewRequest(http.MethodPost, apiBase+"/api/caregiver/status", bytes.NewReader(body6))
	ackReq.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		ackReq.Header.Set("Authorization", "Bearer "+authToken)
	}
	ackResp, err := client.Do(ackReq)
	if err != nil {
		fail("acknowledge incident", err)
	}
	defer ackResp.Body.Close()
	if ackResp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(ackResp.Body)
		fail("acknowledge incident", fmt.Errorf("status %d: %s", ackResp.StatusCode, string(b)))
	}
	pass("acknowledge incident", t6.duration())

	// 7. Resolve incident
	t7 := startStep()
	resPayload := map[string]string{
		"incident_id": intakeData.IncidentID,
		"status":      "resolved",
	}
	body7, _ := json.Marshal(resPayload)
	resReq, _ := http.NewRequest(http.MethodPost, apiBase+"/api/caregiver/status", bytes.NewReader(body7))
	resReq.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		resReq.Header.Set("Authorization", "Bearer "+authToken)
	}
	resResp, err := client.Do(resReq)
	if err != nil {
		fail("resolve incident", err)
	}
	defer resResp.Body.Close()
	if resResp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resResp.Body)
		fail("resolve incident", fmt.Errorf("status %d: %s", resResp.StatusCode, string(b)))
	}
	pass("resolve incident", t7.duration())

	fmt.Println("\nAll smoke tests completed successfully.")
}
