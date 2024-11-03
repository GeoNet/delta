package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

type handler struct {
	db *sql.DB
}

func newHandler(db *sql.DB) *http.ServeMux {
	h := handler{
		db: db,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", homePage)
	mux.HandleFunc("GET /network", h.GetNetworks)
	mux.HandleFunc("GET /network/", h.GetNetworks)
	mux.HandleFunc("GET /network/{code}", h.GetNetworkCode)
	mux.HandleFunc("GET /network/{code}/", h.GetNetworkCode)
	mux.HandleFunc("GET /network/{code}/station", h.GetNetworkStations)
	mux.HandleFunc("GET /network/{code}/station/", h.GetNetworkStations)
	mux.HandleFunc("GET /network/{code}/mark", h.GetNetworkMarks)
	mux.HandleFunc("GET /network/{code}/mark/", h.GetNetworkMarks)
	mux.HandleFunc("GET /network/{code}/mount", h.GetNetworkMounts)
	mux.HandleFunc("GET /network/{code}/mount/", h.GetNetworkMounts)
	mux.HandleFunc("GET /network/{code}/sample", h.GetNetworkSamples)
	mux.HandleFunc("GET /network/{code}/sample/", h.GetNetworkSamples)

	mux.HandleFunc("GET /station", h.GetStations)
	mux.HandleFunc("GET /station/", h.GetStations)
	mux.HandleFunc("GET /station/{code}", h.GetStationCode)
	mux.HandleFunc("GET /station/{code}/", h.GetStationCode)
	mux.HandleFunc("GET /station/{code}/site", h.GetStationSites)
	mux.HandleFunc("GET /station/{code}/site/", h.GetStationSites)
	mux.HandleFunc("GET /station/{code}/site/{location}", h.GetStationSiteLocation)
	mux.HandleFunc("GET /station/{code}/site/{location}/", h.GetStationSiteLocation)
	mux.HandleFunc("GET /station/{code}/site/{location}/sensor", h.GetStationSiteSensors)
	mux.HandleFunc("GET /station/{code}/site/{location}/sensor/", h.GetStationSiteSensors)

	mux.HandleFunc("GET /site", h.GetSites)
	mux.HandleFunc("GET /site/", h.GetSites)
	mux.HandleFunc("GET /site/{location}", h.GetSiteLocations)
	mux.HandleFunc("GET /site/{location}/", h.GetSiteLocations)

	mux.HandleFunc("GET /mark", h.GetMarks)
	mux.HandleFunc("GET /mark/", h.GetMarks)
	mux.HandleFunc("GET /mark/{code}", h.GetMarkCode)
	mux.HandleFunc("GET /mark/{code}/", h.GetMarkCode)

	mux.HandleFunc("GET /monument", h.GetMonuments)
	mux.HandleFunc("GET /monument/", h.GetMonuments)
	mux.HandleFunc("GET /monument/{mark}", h.GetMonumentMark)
	mux.HandleFunc("GET /monument/{mark}/", h.GetMonumentMark)

	mux.HandleFunc("GET /sample", h.GetSamples)
	mux.HandleFunc("GET /sample/", h.GetSamples)
	mux.HandleFunc("GET /sample/{code}", h.GetSampleCode)
	mux.HandleFunc("GET /sample/{code}/", h.GetSampleCode)
	mux.HandleFunc("GET /sample/{code}/point", h.GetSamplePoints)
	mux.HandleFunc("GET /sample/{code}/point/", h.GetSamplePoints)
	mux.HandleFunc("GET /sample/{code}/point/{location}", h.GetSamplePointLocation)
	mux.HandleFunc("GET /sample/{code}/point/{location}/", h.GetSamplePointLocation)
	mux.HandleFunc("GET /sample/{code}/point/{location}/sensor", h.GetSamplePointLocationSensors)
	mux.HandleFunc("GET /sample/{code}/point/{location}/sensor/", h.GetSamplePointLocationSensors)

	mux.HandleFunc("GET /sensor", h.GetSensors)
	mux.HandleFunc("GET /sensor/", h.GetSensors)

	return mux
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the DELTA REST API")
	fmt.Println("DELTA REST API")
}
