package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GeoNet/delta/meta"
)

func (s Settings) Gaps(set *meta.Set) error {

	now := time.Now().UTC()

	files := make(map[string]string)

	if err := filepath.Walk(s.output, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		switch d := filepath.Base(filepath.Dir(path)); d {
		case "global":
			files[filepath.Join(d, filepath.Base(path))] = path
		default:
			files[filepath.Base(path)] = path
		}
		return nil
	}); err != nil {
		return err
	}

	check := make(map[string]interface{})
	for _, n := range strings.Split(s.configured, ",") {
		check[strings.TrimSpace(strings.ToUpper(n))] = true
	}

	extra := make(map[string]interface{})
	for _, v := range strings.Split(s.include, ",") {
		extra[strings.TrimSpace(strings.ToUpper(v))] = true
	}

	skip := make(map[string]interface{})
	for _, v := range strings.Split(s.skip, ",") {
		skip[strings.TrimSpace(strings.ToUpper(v))] = true
	}

	var unmatch []string
	for _, p := range strings.Split(s.exclude, ",") {
		if s := func() string {
			switch r := strings.Split(strings.ToUpper(p), "_"); len(r) {
			case 2:
				return strings.TrimSpace(r[0]) + "_" + strings.TrimSpace(r[1])
			case 1:
				if strings.ContainsAny(r[0], "*?\\[") {
					return strings.TrimSpace(r[0])
				}
				return "*_" + strings.TrimSpace(r[0])
			default:
				return ""
			}
		}(); s != "" {
			unmatch = append(unmatch, s)
		}
	}

	globals := make(map[string]Global)
	stationMap := make(map[string]Station)

	lookup := make(map[string]meta.Network)
	for _, n := range set.Networks() {
		lookup[n.Code] = n
	}

	for _, station := range set.Stations() {
		if _, ok := check[station.Network]; !ok {
			if _, ok := extra[station.Code]; !ok {
				continue
			}
		}
		network, ok := lookup[station.Network]
		if !ok {
			continue
		}

		for _, site := range set.Sites() {
			if site.Station != station.Code {
				continue
			}

			stream := ToStream(network.External, station.Code, site.Location)
			if _, ok := skip[stream]; ok {
				continue
			}

			for _, collection := range set.Collections(site) {
				if collection.Start.After(now) || collection.End.Before(now) {
					continue
				}

				if v, ok := stationMap[ToKey(network.External, station.Code)]; ok {
					if collection.InstalledSensor.Location >= v.Global.Location {
						continue
					}
				}

				var code [2]byte
				copy(code[:], collection.Code())

				global := Global{
					Stream:   string(code[:]),
					Location: collection.Stream.Location,
					Headless: true,
				}
				if global.Style() == "" {
					continue
				}

				if _, ok := globals[global.Key()]; !ok {
					globals[global.Key()] = global
				}

				station := Station{
					Global:   global,
					Code:     station.Code,
					Network:  network.External,
					Gaps:     true,
					Headless: true,
				}

				stationMap[station.Key()] = station
			}
		}
	}

	var count int
	for _, g := range globals {
		slog.Debug("gaps", slog.String("store", g.Path()))
		if err := Store(g, s.output); err != nil {
			return err
		}
		delete(files, g.Path())
		count++
	}

	for k, v := range stationMap {
		for _, m := range unmatch {
			ok, err := filepath.Match(m, k)
			if err != nil {
				return err
			}
			if ok {
				continue
			}

			slog.Debug("gaps", slog.String("store", v.Path()))
			if err := Store(v, s.output); err != nil {
				return err
			}

			delete(files, v.Path())

			count++
		}
	}

	for k, v := range files {
		slog.Info("gaps", slog.String("remove", k))
		if err := os.Remove(v); err != nil {
			return err
		}
	}

	slog.Info("gaps", slog.Int("found", count), slog.Int("removed", len(files)))

	return nil
}
