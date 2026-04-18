package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/0xatanda/profileIntelligence/internal/model"
	repository "github.com/0xatanda/profileIntelligence/internal/respository"
	"github.com/0xatanda/profileIntelligence/internal/service"
	"github.com/gofrs/uuid"
)

type Handler struct {
	Repo *repository.Repo
}

// --------------------
// helper
// --------------------
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// --------------------
// CREATE PROFILE
// --------------------
func (h *Handler) CreateProfile(w http.ResponseWriter, r *http.Request) {

	var body struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 422, map[string]string{"status": "error", "message": "invalid request body"})
		return
	}

	name := strings.TrimSpace(body.Name)

	if name == "" {
		writeJSON(w, 400, map[string]string{"status": "error", "message": "Name is required"})
		return
	}

	// idempotency
	existing, err := h.Repo.FindByName(name)
	if err == nil && existing != nil && existing.ID != "" {
		writeJSON(w, 200, map[string]any{
			"status":  "success",
			"message": "Profile already exists",
			"data":    existing,
		})
		return
	}

	g, a, n, err := service.FetchAll(name)
	if err != nil {
		writeJSON(w, 502, map[string]string{"status": "error", "message": "external API failure"})
		return
	}

	// edge cases
	if g.Gender == nil || g.Count == nil || *g.Count == 0 {
		writeJSON(w, 502, map[string]string{"status": "error", "message": "Genderize invalid response"})
		return
	}

	if a.Age == nil {
		writeJSON(w, 502, map[string]string{"status": "error", "message": "Agify invalid response"})
		return
	}

	if len(n.Country) == 0 {
		writeJSON(w, 502, map[string]string{"status": "error", "message": "Nationalize invalid response"})
		return
	}

	// age group
	age := *a.Age
	ageGroup := "adult"

	switch {
	case age <= 12:
		ageGroup = "child"
	case age <= 19:
		ageGroup = "teenager"
	case age <= 59:
		ageGroup = "adult"
	default:
		ageGroup = "senior"
	}

	// best country
	best := n.Country[0]
	for _, c := range n.Country {
		if c.Probability > best.Probability {
			best = c
		}
	}

	id, _ := uuid.NewV7()

	profile := model.Profile{
		ID:                 id.String(),
		Name:               name,
		Gender:             *g.Gender,
		GenderProbability:  g.Probability,
		SampleSize:         *g.Count,
		Age:                age,
		AgeGroup:           ageGroup,
		CountryID:          best.CountryID,
		CountryProbability: best.Probability,
		CreatedAt:          time.Now().UTC(),
	}

	if err := h.Repo.Create(&profile); err != nil {
		writeJSON(w, 500, map[string]string{"status": "error", "message": "failed to save profile"})
		return
	}

	writeJSON(w, 201, map[string]any{"status": "success", "data": profile})
}

// --------------------
// GET ALL
// --------------------
func (h *Handler) GetProfiles(w http.ResponseWriter, r *http.Request) {

	gender := r.URL.Query().Get("gender")
	country := r.URL.Query().Get("country_id")
	ageGroup := r.URL.Query().Get("age_group")

	data, err := h.Repo.FindAll(gender, country, ageGroup)
	if err != nil {
		writeJSON(w, 500, map[string]string{"status": "error", "message": "failed to fetch profiles"})
		return
	}

	writeJSON(w, 200, map[string]any{
		"status": "success",
		"count":  len(data),
		"data":   data,
	})
}

// --------------------
// GET BY ID
// --------------------
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/api/profiles/")

	profile, err := h.Repo.FindByID(id)
	if err != nil {
		writeJSON(w, 404, map[string]string{"status": "error", "message": "Profile not found"})
		return
	}

	writeJSON(w, 200, map[string]any{"status": "success", "data": profile})
}

// --------------------
// DELETE
// --------------------
func (h *Handler) DeleteProfile(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/api/profiles/")

	if err := h.Repo.Delete(id); err != nil {
		writeJSON(w, 404, map[string]string{"status": "error", "message": "Profile not found"})
		return
	}

	w.WriteHeader(204)
}
