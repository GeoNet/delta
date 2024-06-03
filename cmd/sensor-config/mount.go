package main

import (
	"sort"

	"github.com/GeoNet/delta/meta"
)

func (n *Network) Camera(set *meta.Set, network, label string) error {

	net, ok := set.Network(network)
	if !ok {
		return nil
	}

	for _, mount := range set.Mounts() {
		if mount.Network != net.Code {
			continue
		}

		var views []View
		for _, view := range set.Views() {
			if mount.Code != view.Mount {
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
				cameras = append(cameras, Sensor{
					Make:  camera.Make,
					Model: camera.Model,
					Type:  label,

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

		sort.Slice(views, func(i, j int) bool {
			return views[i].Less(views[j])
		})

		n.Mounts = append(n.Mounts, Mount{
			Code:        mount.Code,
			Network:     net.External,
			Name:        mount.Name,
			Description: net.Description,
			Mount:       mount.Description,

			Latitude:  mount.Latitude,
			Longitude: mount.Longitude,
			Elevation: mount.Elevation,
			Datum:     mount.Datum,

			StartDate: mount.Start,
			EndDate:   mount.End,

			Views: views,
		})
	}

	return nil
}

func (n *Network) Doas(set *meta.Set, network, label string) error {

	net, ok := set.Network(network)
	if !ok {
		return nil
	}

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
					Make:  doas.Make,
					Model: doas.Model,
					Type:  label,

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

		sort.Slice(views, func(i, j int) bool {
			return views[i].Less(views[j])
		})

		n.Mounts = append(n.Mounts, Mount{
			Code:        mount.Code,
			Network:     net.External,
			Name:        mount.Name,
			Description: net.Description,
			Mount:       mount.Description,

			Latitude:  mount.Latitude,
			Longitude: mount.Longitude,
			Elevation: mount.Elevation,
			Datum:     mount.Datum,

			StartDate: mount.Start,
			EndDate:   mount.End,

			Views: views,
		})
	}

	return nil
}

func (s Settings) Cameras(set *meta.Set, name, network string) (Group, bool) {

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

			var cameras []Sensor

			for _, camera := range set.InstalledCameras() {
				if mount.Code != camera.Mount {
					continue
				}
				if view.Code != camera.View {
					continue
				}
				cameras = append(cameras, Sensor{
					Make:  camera.Make,
					Model: camera.Model,

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
					Make:  doas.Make,
					Model: doas.Model,

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
