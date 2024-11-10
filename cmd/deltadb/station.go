package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GeoNet/delta/meta/sqlite"
)

func (h handler) GetStations(w http.ResponseWriter, r *http.Request) {

	stations, err := sqlite.Stations(r.Context(), h.db)
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

func (h handler) GetStationCode(w http.ResponseWriter, r *http.Request) {

	code := r.PathValue("code")

	stations, err := sqlite.Stations(r.Context(), h.db, sqlite.Code(code))
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

func (h handler) GetStationSites(w http.ResponseWriter, r *http.Request) {

	station := r.PathValue("code")

	sites, err := sqlite.Sites(r.Context(), h.db, sqlite.Station(station))
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

func (h handler) GetStationSiteLocation(w http.ResponseWriter, r *http.Request) {

	station := r.PathValue("code")
	location := r.PathValue("location")

	sites, err := sqlite.Sites(r.Context(), h.db, sqlite.Station(station), sqlite.Location(location))
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

func (h handler) GetStationSiteSensors(w http.ResponseWriter, r *http.Request) {

	station := r.PathValue("code")
	location := r.PathValue("location")

	sensors, err := sqlite.Sensors(r.Context(), h.db, sqlite.Station(station), sqlite.Location(location))
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sensors); err != nil {
		log.Println(err)
	}

}
