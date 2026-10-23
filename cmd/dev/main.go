package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/iokuru/sentry/internal/hotline"
	"github.com/iokuru/sentry/internal/model"
	"github.com/iokuru/sentry/internal/risk"
)

type MemoryStore struct {
	mu        sync.RWMutex
	incidents map[string]*model.Incident
	uploads   map[string][]byte
}

var devStore = &MemoryStore{
	incidents: make(map[string]*model.Incident),
	uploads:   make(map[string][]byte),
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/presign", handlePresign)
	mux.HandleFunc("/upload/", handleUpload)
	mux.HandleFunc("/api/intake", handleIntake)
	mux.HandleFunc("/api/incidents", handleGetIncident)
	mux.HandleFunc("/api/caregiver/incidents", handleCaregiverList)
	mux.HandleFunc("/api/caregiver/status", handleCaregiverStatus)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		mux.ServeHTTP(w, r)
	})

	log.Printf("sentry local dev server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func handlePresign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	randBytes := make([]byte, 8)
	_, _ = rand.Read(randBytes)
	key := fmt.Sprintf("dev-%d-%s.png", time.Now().Unix(), hex.EncodeToString(randBytes))

	res := map[string]string{
		"upload_url": fmt.Sprintf("http://localhost:8080/upload/%s", key),
		"image_key":  key,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := strings.TrimPrefix(r.URL.Path, "/upload/")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	devStore.mu.Lock()
	devStore.uploads[key] = body
	devStore.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func handleIntake(w http.ResponseWriter, r *http.Request) {
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

	// stub ocr text when testing locally with sample images
	extracted := req.UserContext
	if req.ImageKey != "" && strings.Contains(req.ImageKey, "bank") {
		extracted = "Your Chase bank debit card has been suspended. Call immediately or pay $500 gift card."
	}

	input := risk.Input{
		Text:    extracted,
		Context: req.UserContext,
	}
	score, level, matches := risk.CalculateScore(input)
	cat := risk.ClassifyCategory(input)
	out := risk.BuildSeniorOutput(cat, level, country)

	now := time.Now().UTC().Format(time.RFC3339)
	inc := &model.Incident{
		ID:            id,
		SeniorID:      req.SeniorID,
		ImageKey:      req.ImageKey,
		ExtractedText: extracted,
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

	devStore.mu.Lock()
	devStore.incidents[id] = inc
	devStore.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"incident_id": id,
		"status":      string(model.StatusAnalyzed),
	})
}

func handleGetIncident(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		handleCaregiverList(w, r)
		return
	}

	devStore.mu.RLock()
	inc, ok := devStore.incidents[id]
	devStore.mu.RUnlock()

	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(inc)
}

func handleCaregiverList(w http.ResponseWriter, r *http.Request) {
	devStore.mu.RLock()
	var list []model.Incident
	for _, inc := range devStore.incidents {
		list = append(list, *inc)
	}
	devStore.mu.RUnlock()

	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt > list[j].CreatedAt
	})

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

func handleCaregiverStatus(w http.ResponseWriter, r *http.Request) {
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

	devStore.mu.Lock()
	inc, ok := devStore.incidents[req.IncidentID]
	if ok {
		inc.Status = req.Status
		inc.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	devStore.mu.Unlock()

	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"incident_id": req.IncidentID,
		"status":      string(req.Status),
	})
}
