package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GeoNet/delta/meta/sqlite"
)

func (h handler) GetSamples(w http.ResponseWriter, r *http.Request) {

	samples, err := sqlite.Samples(r.Context(), h.db)
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

func (h handler) GetSampleCode(w http.ResponseWriter, r *http.Request) {

	code := r.PathValue("code")

	samples, err := sqlite.Samples(r.Context(), h.db, sqlite.Code(code))
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

func (h handler) GetSamplePoints(w http.ResponseWriter, r *http.Request) {

	code := r.PathValue("code")

	points, err := sqlite.Points(r.Context(), h.db, sqlite.Sample(code))
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(points); err != nil {
		log.Println(err)
	}
}

func (h handler) GetSamplePointLocation(w http.ResponseWriter, r *http.Request) {

	code := r.PathValue("code")
	location := r.PathValue("location")

	points, err := sqlite.Points(r.Context(), h.db, sqlite.Sample(code), sqlite.Location(location))
	if err != nil {
		log.Println("failed to execute query", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(points); err != nil {
		log.Println(err)
	}
}

func (h handler) GetSamplePointLocationSensors(w http.ResponseWriter, r *http.Request) {

	code := r.PathValue("code")
	location := r.PathValue("location")

	sensors, err := sqlite.Sensors(r.Context(), h.db, sqlite.Station(code), sqlite.Location(location))
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
