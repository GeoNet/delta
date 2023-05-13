package main

import (
	"os"
	"path/filepath"
)

type Settings struct {
	baseDir string

	daysGrace  int
	purgeFiles bool
	outputDir  string

	includeNetworks []string
	excludeNetworks []string

	includeStations []Station
	excludeStations []Station
	extraStations   []Station
}

func (s Settings) Walk() ([]string, error) {
	var files []string
	if err := filepath.Walk(s.outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	}); err != nil {
		return nil, err
	}
	return files, nil
}

func (s Settings) ExcludeNetwork(network string) bool {

	// may want to exclude single networks
	for _, n := range s.excludeNetworks {
		if n != network {
			continue
		}
		return true
	}

	// no include networks listed, means all non-excluded networks wanted
	if len(s.includeNetworks) == 0 {
		return false
	}

	// otherwise check included networks
	for _, n := range s.includeNetworks {
		if n != network {
			continue
		}
		return false
	}

	return true
}

func (s Settings) ExcludeStation(stn Station) bool {

	// may want to exclude single networks
	for _, x := range s.excludeStations {
		if ok, err := filepath.Match(x.Key(), stn.Key()); ok || err != nil {
			return true
		}
	}

	// no include stations listed, means all non-excluded stations wanted
	if len(s.includeStations) == 0 {
		return false
	}

	// otherwise check included stations
	for _, x := range s.includeStations {
		if ok, err := filepath.Match(x.Key(), stn.Key()); ok && err == nil {
			return false
		}
	}

	return true
}
