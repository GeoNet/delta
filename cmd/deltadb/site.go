package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GeoNet/delta/meta/sqlite"
)

func (h handler) GetSites(w http.ResponseWriter, r *http.Request) {

	sites, err := sqlite.Sites(r.Context(), h.db)
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sites); err != nil {
		log.Println(err)
	}
}

func (h handler) GetSiteLocations(w http.ResponseWriter, r *http.Request) {

	location := r.PathValue("location")

	sites, err := sqlite.Sites(r.Context(), h.db, sqlite.Location(location))
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sites); err != nil {
		log.Println(err)
	}
}
