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
	"github.com/iokuru/sentry/internal/pairing"
	"github.com/iokuru/sentry/internal/risk"
	"github.com/iokuru/sentry/internal/store"
	"github.com/iokuru/sentry/internal/telegram"
)

type MemoryStore struct {
	mu        sync.RWMutex
	incidents map[string]*model.Incident
	uploads   map[string][]byte
	invites   map[string]*model.Invite
	links     map[string][]*model.CaregiverLink
}

func newMemoryStore() *MemoryStore {
	return &MemoryStore{
		incidents: make(map[string]*model.Incident),
		uploads:   make(map[string][]byte),
		invites:   make(map[string]*model.Invite),
		links:     make(map[string][]*model.CaregiverLink),
	}
}

func (s *MemoryStore) SaveInvite(inv *model.Invite) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.invites[inv.Code] = inv
}

func (s *MemoryStore) GetInvite(code string) (*model.Invite, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	inv, ok := s.invites[code]
	return inv, ok
}

func (s *MemoryStore) RedeemInvite(code, seniorID string) (*model.CaregiverLink, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	inv, ok := s.invites[code]
	if !ok {
		return nil, store.ErrInviteNotFound
	}
	if inv.Used {
		return nil, store.ErrInviteUsed
	}
	if time.Now().UTC().Unix() > inv.ExpiresAt {
		return nil, store.ErrInviteExpired
	}
	inv.Used = true
	name := inv.SeniorName
	if name == "" {
		name = "Senior"
	}
	link := &model.CaregiverLink{
		CaregiverID: inv.CaregiverID,
		SeniorID:    seniorID,
		SeniorName:  name,
		LinkedAt:    time.Now().UTC(),
	}
	s.links[inv.CaregiverID] = append(s.links[inv.CaregiverID], link)
	return link, nil
}

func (s *MemoryStore) ListLinks(caregiverID string) []*model.CaregiverLink {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := s.links[caregiverID]
	if list == nil {
		return []*model.CaregiverLink{}
	}
	return list
}

func (s *MemoryStore) DeleteLink(caregiverID, seniorID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	current := s.links[caregiverID]
	filtered := make([]*model.CaregiverLink, 0, len(current))
	for _, l := range current {
		if l.SeniorID != seniorID {
			filtered = append(filtered, l)
		}
	}
	s.links[caregiverID] = filtered
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
	mux.HandleFunc("/api/caregiver/invites", s.handleCaregiverInvites)
	mux.HandleFunc("/api/caregiver/links", s.handleCaregiverLinks)
	mux.HandleFunc("/api/caregiver/links/", s.handleCaregiverLinks)
	mux.HandleFunc("/api/senior/pair", s.handleSeniorPair)
	mux.HandleFunc("/telegram/simulate", s.handleTelegramSimulate)

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

		if r.URL.Path == "/" || r.URL.Path == "/api/health" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"healthy","service":"sentry"}`))
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

	var id string
	if strings.TrimSpace(req.ImageKey) != "" {
		hash := sha256.Sum256([]byte(strings.TrimSpace(req.ImageKey)))
		id = fmt.Sprintf("inc-%s", hex.EncodeToString(hash[:10]))
		if existing, ok := s.store.Get(id); ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"incident_id": existing.ID,
				"status":      string(existing.Status),
			})
			return
		}
	} else {
		randBytes := make([]byte, 8)
		_, _ = rand.Read(randBytes)
		id = fmt.Sprintf("inc-%d-%s", time.Now().Unix(), hex.EncodeToString(randBytes))
	}

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
	detectedOneLine := "suspicious pattern detected"
	if len(inc.Signals) > 0 {
		detectedOneLine = inc.Signals[0].Description
	}

	formattedTelegram, _ := telegram.FormatAlertMessage(
		string(inc.Category),
		string(inc.RiskLevel),
		detectedOneLine,
		inc.Output.Hotline.Name,
		inc.Output.Hotline.Number,
		inc.ID,
		inc.UpdatedAt,
	)

	fmt.Printf("\n==================== [DEV ALERT] ====================\n")
	fmt.Printf("[Email Alert]: Caregiver email notification queued for %s\n\n", inc.ID)
	fmt.Printf("[Telegram Alert Message]:\n%s\n\n", formattedTelegram)
	fmt.Printf("Simulated Telegram Action Buttons:\n")
	fmt.Printf("-> Acknowledge: http://localhost:8080/telegram/simulate?action=ack&id=%s\n", inc.ID)
	fmt.Printf("-> Resolve:     http://localhost:8080/telegram/simulate?action=res&id=%s\n", inc.ID)
	fmt.Printf("=====================================================\n\n")
}

func (s *Server) handleTelegramSimulate(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	id := r.URL.Query().Get("id")

	if action == "start" {
		code := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("code")))
		inv, ok := s.store.GetInvite(code)
		if !ok {
			http.Error(w, "invite code not found", http.StatusNotFound)
			return
		}
		links := s.store.ListLinks(inv.CaregiverID)
		for _, l := range links {
			l.TelegramChatID = "simulated-chat-id"
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<h3>Telegram /start processed</h3><p>Linked Telegram to caregiver account for <b>%s</b>.</p>", inv.SeniorName)
		return
	}

	var newStatus model.IncidentStatus
	if action == "ack" {
		newStatus = model.StatusAcknowledged
	} else if action == "res" {
		newStatus = model.StatusResolved
	} else {
		http.Error(w, "invalid action", http.StatusBadRequest)
		return
	}

	if !s.store.UpdateStatus(id, newStatus) {
		http.Error(w, "incident not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h3>Telegram callback processed successfully</h3><p>Incident <b>%s</b> is now marked as <b>%s</b>.</p><p><a href=\"/caregiver\">View Caregiver Portal</a></p>", id, newStatus)
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
	if inc.ImageKey != "" {
		inc.ThumbnailURL = "/upload/" + inc.ImageKey
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(inc)
}

func (s *Server) handleCaregiverList(w http.ResponseWriter, r *http.Request) {
	list := s.store.List()
	for _, inc := range list {
		if inc.ImageKey != "" {
			inc.ThumbnailURL = "/upload/" + inc.ImageKey
		}
	}
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

	if req.Status != model.StatusAcknowledged && req.Status != model.StatusResolved && req.Status != model.StatusAnalyzed {
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

func (s *Server) handleCaregiverInvites(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SeniorName string `json:"senior_name"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	code, err := pairing.GenerateCode()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	now := time.Now().UTC().Unix()
	inv := &model.Invite{
		Code:        code,
		CaregiverID: "caregiver",
		SeniorName:  strings.TrimSpace(req.SeniorName),
		CreatedAt:   now,
		ExpiresAt:   now + 86400,
	}
	s.store.SaveInvite(inv)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code":       code,
		"expires_at": inv.ExpiresAt,
	})
}

func (s *Server) handleCaregiverLinks(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	switch r.Method {
	case http.MethodGet:
		links := s.store.ListLinks("caregiver")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(links)

	case http.MethodDelete:
		path := r.URL.Path
		seniorID := strings.TrimPrefix(path, "/api/caregiver/links/")
		if seniorID == "" || seniorID == "/api/caregiver/links" {
			seniorID = r.URL.Query().Get("senior_id")
		}
		if seniorID == "" {
			var body struct {
				SeniorID string `json:"senior_id"`
			}
			_ = json.NewDecoder(r.Body).Decode(&body)
			seniorID = body.SeniorID
		}
		if seniorID == "" {
			http.Error(w, "senior_id is required", http.StatusBadRequest)
			return
		}
		s.store.DeleteLink("caregiver", seniorID)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleSeniorPair(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Code     string `json:"code"`
		DeviceID string `json:"device_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	code := strings.ToUpper(strings.TrimSpace(req.Code))
	deviceID := strings.TrimSpace(req.DeviceID)
	if code == "" || deviceID == "" {
		http.Error(w, "code and device_id are required", http.StatusBadRequest)
		return
	}

	link, err := s.store.RedeemInvite(code, deviceID)
	if err != nil {
		switch err {
		case store.ErrInviteNotFound:
			http.Error(w, "invite code not found", http.StatusNotFound)
		case store.ErrInviteExpired:
			http.Error(w, "invite code expired", http.StatusGone)
		case store.ErrInviteUsed:
			http.Error(w, "invite code already used", http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"senior_id":   link.SeniorID,
		"senior_name": link.SeniorName,
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
