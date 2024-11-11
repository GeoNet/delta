package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GeoNet/delta/meta/sqlite"
)

func (h handler) GetNetworks(w http.ResponseWriter, r *http.Request) {

	networks, err := sqlite.Networks(r.Context(), h.db)
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(networks); err != nil {
		log.Println(err)
	}
}

func (h handler) GetNetworkCode(w http.ResponseWriter, r *http.Request) {

	code := r.PathValue("code")

	networks, err := sqlite.Networks(r.Context(), h.db, sqlite.Code(code))
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(networks); err != nil {
		log.Println(err)
	}
}

func (h handler) GetNetworkStations(w http.ResponseWriter, r *http.Request) {

	code := r.PathValue("code")

	stations, err := sqlite.Stations(r.Context(), h.db, sqlite.Network(code))
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stations); err != nil {
		log.Println(err)
	}
}

func (h handler) GetNetworkMarks(w http.ResponseWriter, r *http.Request) {

	code := r.PathValue("code")

	marks, err := sqlite.Marks(r.Context(), h.db, sqlite.Network(code))
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

func (h handler) GetNetworkSamples(w http.ResponseWriter, r *http.Request) {

	code := r.PathValue("code")

	samples, err := sqlite.Samples(r.Context(), h.db, sqlite.Network(code))
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(samples); err != nil {
		log.Println(err)
	}
}

func (h handler) GetNetworkMounts(w http.ResponseWriter, r *http.Request) {

	code := r.PathValue("code")

	mounts, err := sqlite.Mounts(r.Context(), h.db, sqlite.Network(code))
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(mounts); err != nil {
		log.Println(err)
	}
}
