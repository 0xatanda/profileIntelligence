package main

import (
	"log"
	"net/http"
	"os"

	"github.com/0xatanda/profileIntelligence/internal/db"
	"github.com/0xatanda/profileIntelligence/internal/handler"
	repository "github.com/0xatanda/profileIntelligence/internal/respository"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	database := db.Connect()
	defer database.Close()

	repo := &repository.Repo{DB: database}
	h := &handler.Handler{Repo: repo}

	http.HandleFunc("/api/profiles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateProfile(w, r)
		case http.MethodGet:
			h.GetProfiles(w, r)
		default:
			w.WriteHeader(405)
		}
	})

	http.HandleFunc("/api/profiles/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetProfile(w, r)
		case http.MethodDelete:
			h.DeleteProfile(w, r)
		default:
			w.WriteHeader(405)
		}
	})

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
