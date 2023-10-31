package main

import (
	"strconv"
	"strings"

	"github.com/GeoNet/delta/internal/ntrip"
	"github.com/GeoNet/delta/meta"
)

// NewConfig is a function that extracts information from delta meta data and local csv files.
func NewConfig(set *meta.Set, caster *ntrip.Caster, extra bool) (*Config, error) {

	var config Config

	mounts := make(map[string]ntrip.Mount)

	for _, m := range caster.Mounts() {
		mark, ok := set.Mark(m.Mark)
		if !ok {
			continue
		}
		mounts[m.Mount] = m

		config.Mounts = append(config.Mounts, MountConfig{
			Mount:     m.Mount,
			Mark:      mark.Code,
			Name:      mark.Name,
			Country:   m.Country,
			Latitude:  strconv.FormatFloat(mark.Latitude, 'f', 2, 64),
			Longitude: strconv.FormatFloat(mark.Longitude, 'f', 2, 64),
			Format:    strings.Replace(m.Format, " ", "_", -1),
		})
	}

	if extra {
		for _, a := range caster.Aliases() {
			mount, ok := mounts[a.Mount]
			if !ok {
				continue
			}
			mark, ok := set.Mark(mount.Mark)
			if !ok {
				continue
			}

			config.Mounts = append(config.Mounts, MountConfig{
				Mount:     a.Alias,
				Mark:      mark.Code,
				Name:      mark.Name,
				Country:   mount.Country,
				Latitude:  strconv.FormatFloat(mark.Latitude, 'f', 2, 64),
				Longitude: strconv.FormatFloat(mark.Longitude, 'f', 2, 64),
				Format:    strings.Replace(mount.Format, " ", "_", -1),
			})
		}
	}

	config.Sort()

	return &config, nil
}
