package main

import (
	"sort"
	"strings"

	"github.com/GeoNet/delta/meta"
)

func (s Settings) Cameras(set *meta.Set, name, networks string) (Group, bool) {

	nets := make(map[string]interface{})
	for _, n := range strings.Split(networks, ",") {
		if n = strings.TrimSpace(n); n != "" {
			nets[n] = true
		}
	}

	var mounts []Mount
	for _, mount := range set.Mounts() {
		net, ok := set.Network(mount.Network)
		if !ok {
			continue
		}
		if _, ok := nets[net.Code]; !ok {
			continue
		}

		var views []View
		for _, view := range set.Views() {
			if mount.Code != view.Mount {
				continue
			}

			if view.Start.After(mount.End) {
				continue
			}

			if view.End.Before(mount.Start) {
				continue
			}

			var cameras []Sensor

			for _, camera := range set.InstalledCameras() {
				if mount.Code != camera.Mount {
					continue
				}
				if view.Code != camera.View {
					continue
				}

				if camera.Start.After(view.End) {
					continue
				}

				if camera.End.Before(view.Start) {
					continue
				}
				cameras = append(cameras, Sensor{
					Make:   camera.Make,
					Model:  camera.Model,
					Serial: camera.Serial,
					Type:   name,

					Dip:     camera.Dip,
					Azimuth: camera.Azimuth,

					Vertical: camera.Vertical,
					North:    camera.North,
					East:     camera.East,

					StartDate: camera.Start,
					EndDate:   camera.End,
				})
			}

			sort.Slice(cameras, func(i, j int) bool {
				return cameras[i].Less(cameras[j])
			})

			views = append(views, View{
				Code:        view.Code,
				Label:       view.Label,
				Description: view.Description,

				Azimuth: view.Azimuth,
				Method:  view.Method,
				Dip:     view.Dip,

				StartDate: view.Start,
				EndDate:   view.End,

				Sensors: cameras,
			})
		}

		if !(len(views) > 0) {
			continue
		}

		sort.Slice(views, func(i, j int) bool {
			return views[i].Less(views[j])
		})

		mounts = append(mounts, Mount{
			Code:        mount.Code,
			Name:        mount.Name,
			Mount:       mount.Description,
			Network:     mount.Network,
			External:    net.External,
			Description: net.Description,

			Latitude:  mount.Latitude,
			Longitude: mount.Longitude,
			Elevation: mount.Elevation,
			Datum:     mount.Datum,

			StartDate: mount.Start,
			EndDate:   mount.End,

			Views: views,
		})
	}

	return Group{Name: name, Mounts: mounts}, true
}

func (s Settings) Doases(set *meta.Set, name, network string) (Group, bool) {

	net, ok := set.Network(network)
	if !ok {
		return Group{}, false
	}

	var mounts []Mount
	for _, mount := range set.Mounts() {
		if mount.Network != net.Code {
			continue
		}

		var views []View
		for _, view := range set.Views() {
			if mount.Code != view.Mount {
				continue
			}

			var doases []Sensor

			for _, doas := range set.Doases() {
				if mount.Code != doas.Mount {
					continue
				}
				if view.Code != doas.View {
					continue
				}
				doases = append(doases, Sensor{
					Make:   doas.Make,
					Model:  doas.Model,
					Serial: doas.Serial,
					Type:   name,

					Dip:     doas.Dip,
					Azimuth: doas.Azimuth,

					Vertical: doas.Vertical,
					North:    doas.North,
					East:     doas.East,

					StartDate: doas.Start,
					EndDate:   doas.End,
				})
			}

			sort.Slice(doases, func(i, j int) bool {
				return doases[i].Less(doases[j])
			})

			views = append(views, View{
				Code:        view.Code,
				Label:       view.Label,
				Description: view.Description,

				Azimuth: view.Azimuth,
				Method:  view.Method,
				Dip:     view.Dip,

				StartDate: view.Start,
				EndDate:   view.End,

				Sensors: doases,
			})
		}

		if !(len(views) > 0) {
			continue
		}

		sort.Slice(views, func(i, j int) bool {
			return views[i].Less(views[j])
		})

		mounts = append(mounts, Mount{
			Code:        mount.Code,
			Name:        mount.Name,
			Mount:       mount.Description,
			Network:     mount.Network,
			External:    net.External,
			Description: net.Description,

			Latitude:  mount.Latitude,
			Longitude: mount.Longitude,
			Elevation: mount.Elevation,
			Datum:     mount.Datum,

			StartDate: mount.Start,
			EndDate:   mount.End,

			Views: views,
		})
	}

	return Group{Name: name, Mounts: mounts}, true
}
