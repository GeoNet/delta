package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GeoNet/delta/meta/sqlite"
)

func (h handler) GetMarks(w http.ResponseWriter, r *http.Request) {

	marks, err := sqlite.Marks(r.Context(), h.db)
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(marks); err != nil {
		log.Println(err)
	}
}

func (h handler) GetMarkCode(w http.ResponseWriter, r *http.Request) {

	code := r.PathValue("code")

	marks, err := sqlite.Marks(r.Context(), h.db, sqlite.Code(code))
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(marks); err != nil {
		log.Println(err)
	}
}
