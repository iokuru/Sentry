package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/iokuru/sentry/internal/comprehend"
	"github.com/iokuru/sentry/internal/hotline"
	"github.com/iokuru/sentry/internal/model"
	"github.com/iokuru/sentry/internal/risk"
)

type MemoryStore struct {
	mu        sync.RWMutex
	incidents map[string]*model.Incident
	uploads   map[string][]byte
}

func newMemoryStore() *MemoryStore {
	return &MemoryStore{
		incidents: make(map[string]*model.Incident),
		uploads:   make(map[string][]byte),
	}
}

func (s *MemoryStore) Save(inc *model.Incident) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.incidents[inc.ID] = inc
}

func (s *MemoryStore) Get(id string) (*model.Incident, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	inc, ok := s.incidents[id]
	return inc, ok
}

func (s *MemoryStore) List() []*model.Incident {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Incident, 0, len(s.incidents))
	for _, inc := range s.incidents {
		list = append(list, inc)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt > list[j].CreatedAt
	})
	return list
}

func (s *MemoryStore) UpdateStatus(id string, status model.IncidentStatus) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	inc, ok := s.incidents[id]
	if !ok {
		return false
	}
	inc.Status = status
	inc.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return true
}

func (s *MemoryStore) SaveUpload(key string, data []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.uploads[key] = data
}

func (s *MemoryStore) GetUpload(key string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.uploads[key]
	return d, ok
}

type FakeOCR struct{}

func (o *FakeOCR) Extract(key string, data []byte) string {
	lowerKey := strings.ToLower(key)
	switch {
	case strings.Contains(lowerKey, "bank") || strings.Contains(lowerKey, "chase") || strings.Contains(lowerKey, "wells"):
		return "URGENT: Your Chase bank account has been locked due to suspicious activity. Verify your PIN and pay $500 in Target gift cards immediately to prevent permanent suspension."
	case strings.Contains(lowerKey, "delivery") || strings.Contains(lowerKey, "usps") || strings.Contains(lowerKey, "ups"):
		return "USPS Notice: Your package could not be delivered due to invalid zip code. Update your billing address and pay $1.99 redelivery fee at bit.ly/usps-track-99."
	case strings.Contains(lowerKey, "utility") || strings.Contains(lowerKey, "power") || strings.Contains(lowerKey, "electric"):
		return "DISCONNECTION NOTICE: Your electricity will be shut off in 45 minutes due to unpaid bill of $389. Pay immediately by MoneyGram or Bitcoin to avoid disconnection."
	case strings.Contains(lowerKey, "prize") || strings.Contains(lowerKey, "lottery") || strings.Contains(lowerKey, "winner"):
		return "Congratulations! You have been selected to receive $1,500,000 from the Mega Millions sweepstakes. To claim your prize, send $250 processing fee via cryptocurrency."
	case strings.Contains(lowerKey, "relative") || strings.Contains(lowerKey, "grandson") || strings.Contains(lowerKey, "family"):
		return "Hi Grandma, this is your grandson. I had an accident and was arrested. My phone is broken so I am using this number. Please wire $1,200 for bail right now, please do not tell Mom."
	case strings.Contains(lowerKey, "tax") || strings.Contains(lowerKey, "irs"):
		return "IRS FINAL WARNING: A federal warrant has been issued in your name for tax fraud. An officer is en route to arrest you. Call 1-800-555-0182 immediately to settle your fine."
	case strings.Contains(lowerKey, "romance"):
		return "My dearest love, I am stuck at customs on my way to visit you. They are demanding $800 in fees. Please send wire transfer so we can finally be together."
	default:
		if len(data) > 0 {
			h := sha256.Sum256(data)
			hexStr := hex.EncodeToString(h[:4])
			return fmt.Sprintf("Notice from security department regarding transaction %s. Please verify your one time passcode and contact support.", hexStr)
		}
		return ""
	}
}

type Server struct {
	store *MemoryStore
	ocr   *FakeOCR
}

func newServer() *Server {
	return &Server{
		store: newMemoryStore(),
		ocr:   &FakeOCR{},
	}
}

func (s *Server) enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Amz-Target")
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/config.json", s.handleConfig)
	mux.HandleFunc("/api/presign", s.handlePresign)
	mux.HandleFunc("/upload/", s.handleUpload)
	mux.HandleFunc("/api/intake", s.handleIntake)
	mux.HandleFunc("/api/incidents", s.handleGetIncident)
	mux.HandleFunc("/api/caregiver/incidents", s.handleCaregiverList)
	mux.HandleFunc("/api/caregiver/status", s.handleCaregiverStatus)

	// mock cognito idp endpoint for local auth
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.enableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		target := r.Header.Get("X-Amz-Target")
		if strings.Contains(target, "InitiateAuth") {
			s.handleMockCognitoAuth(w, r)
			return
		}

		if r.URL.Path == "/" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("sentry devserver running"))
			return
		}

		http.NotFound(w, r)
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.enableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		mux.ServeHTTP(w, r)
	})
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"api_url":               "",
		"user_pool_id":          "local_user_pool",
		"user_pool_client_id":   "local_client_id",
		"region":                "local",
	})
}

func (s *Server) handleMockCognitoAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-amz-json-1.1")
	resp := map[string]interface{}{
		"AuthenticationResult": map[string]interface{}{
			"AccessToken":  "dev-access-token-sentry-caregiver",
			"IdToken":      "dev-id-token-sentry-caregiver",
			"RefreshToken": "dev-refresh-token-sentry-caregiver",
			"ExpiresIn":    3600,
			"TokenType":    "Bearer",
		},
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handlePresign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ContentType string `json:"content_type"`
		Filename    string `json:"filename"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	ext := ".png"
	if req.Filename != "" {
		parsedExt := strings.ToLower(filepath.Ext(req.Filename))
		if parsedExt == ".jpg" || parsedExt == ".jpeg" || parsedExt == ".webp" || parsedExt == ".png" {
			ext = parsedExt
		}
	}

	randBytes := make([]byte, 8)
	_, _ = rand.Read(randBytes)
	key := fmt.Sprintf("screenshots/%d-%s%s", time.Now().Unix(), hex.EncodeToString(randBytes), ext)

	res := map[string]string{
		"upload_url": fmt.Sprintf("/upload/%s", key),
		"image_key":  key,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := strings.TrimPrefix(r.URL.Path, "/upload/")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	s.store.SaveUpload(key, data)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleIntake(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ImageKey    string `json:"image_key"`
		UserContext string `json:"user_context"`
		Locale      string `json:"locale"`
		Country     string `json:"country"`
		SeniorID    string `json:"senior_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	country := req.Country
	if country == "" {
		country = hotline.ResolveCountry(req.Locale)
	}

	randBytes := make([]byte, 8)
	_, _ = rand.Read(randBytes)
	id := fmt.Sprintf("inc-%d-%s", time.Now().Unix(), hex.EncodeToString(randBytes))

	var extractedText string
	if req.ImageKey != "" {
		data, _ := s.store.GetUpload(req.ImageKey)
		extractedText = s.ocr.Extract(req.ImageKey, data)
	}

	combinedText := extractedText
	if req.UserContext != "" {
		if combinedText != "" {
			combinedText += "\n" + req.UserContext
		} else {
			combinedText = req.UserContext
		}
	}

	compRes := comprehend.FallbackAnalysis(combinedText, "en")

	input := risk.Input{
		Text:       extractedText,
		Context:    req.UserContext,
		Sentiment:  compRes.Sentiment,
		KeyPhrases: compRes.KeyPhrases,
		Entities:   compRes.Entities,
	}

	score, level, matches := risk.CalculateScore(input)
	cat := risk.ClassifyCategory(input)
	out := risk.BuildSeniorOutput(cat, level, country)

	now := time.Now().UTC().Format(time.RFC3339)
	inc := &model.Incident{
		ID:            id,
		SeniorID:      req.SeniorID,
		ImageKey:      req.ImageKey,
		ExtractedText: extractedText,
		UserContext:   req.UserContext,
		Score:         score,
		RiskLevel:     level,
		Category:      cat,
		Signals:       matches,
		Output:        out,
		Status:        model.StatusAnalyzed,
		Country:       country,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	s.store.Save(inc)

	if level == model.RiskHigh {
		s.printAlert(inc)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"incident_id": id,
		"status":      string(model.StatusAnalyzed),
	})
}

func (s *Server) printAlert(inc *model.Incident) {
	fmt.Printf("\n==================== [DEV ALERT] ====================\n")
	fmt.Printf("High risk incident detected!\n")
	fmt.Printf("Incident ID: %s\n", inc.ID)
	fmt.Printf("Risk Score:  %d / 100 (%s)\n", inc.Score, inc.RiskLevel)
	fmt.Printf("Category:    %s\n", inc.Category)
	fmt.Printf("Hold off:    %s\n", inc.Output.HoldOff)
	fmt.Printf("Alternative: %s\n", inc.Output.SafeAlternative)
	fmt.Printf("Hotline:     %s (%s)\n", inc.Output.Hotline.Name, inc.Output.Hotline.Number)
	fmt.Printf("[Email Alert]:   caregiver alert queued\n")
	fmt.Printf("[Telegram Alert]: simulated telegram message dispatched\n")
	fmt.Printf("=====================================================\n\n")
}

func (s *Server) handleGetIncident(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		s.handleCaregiverList(w, r)
		return
	}

	inc, ok := s.store.Get(id)
	if !ok {
		http.Error(w, "incident not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(inc)
}

func (s *Server) handleCaregiverList(w http.ResponseWriter, r *http.Request) {
	list := s.store.List()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

func (s *Server) handleCaregiverStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		IncidentID string               `json:"incident_id"`
		Status     model.IncidentStatus `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Status != model.StatusAcknowledged && req.Status != model.StatusResolved {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	if !s.store.UpdateStatus(req.IncidentID, req.Status) {
		http.Error(w, "incident not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"incident_id": req.IncidentID,
		"status":      string(req.Status),
	})
}

func main() {
	backendOnly := flag.Bool("backend-only", false, "run only backend dev server without frontend")
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := newServer()
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: server.routes(),
	}

	go func() {
		log.Printf("sentry local backend listening on http://localhost:%s", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	var frontendCmd *exec.Cmd
	if !*backendOnly {
		npmCmd := "npm"
		if _, err := exec.LookPath("npm.cmd"); err == nil {
			npmCmd = "npm.cmd"
		}

		frontendCmd = exec.Command(npmCmd, "run", "dev")
		frontendCmd.Dir = "frontend"
		frontendCmd.Stdout = os.Stdout
		frontendCmd.Stderr = os.Stderr

		if err := frontendCmd.Start(); err != nil {
			log.Printf("could not start frontend: %v (run npm run dev in frontend manually)", err)
		} else {
			log.Printf("sentry frontend starting via vite")
		}
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down sentry dev server...")
	if frontendCmd != nil && frontendCmd.Process != nil {
		_ = frontendCmd.Process.Kill()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(ctx)
}
