package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GeoNet/delta/meta/sqlite"
)

func (h handler) GetMonuments(w http.ResponseWriter, r *http.Request) {

	monuments, err := sqlite.Monuments(r.Context(), h.db)
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(monuments); err != nil {
		log.Println(err)
	}
}

func (h handler) GetMonumentMark(w http.ResponseWriter, r *http.Request) {

	mark := r.PathValue("mark")

	monuments, err := sqlite.Monuments(r.Context(), h.db, sqlite.Mark(mark))
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(monuments); err != nil {
		log.Println(err)
	}
}
