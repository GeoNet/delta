package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/GeoNet/delta/internal/ntrip"
	"github.com/GeoNet/delta/meta"
)

// Config is a function that extracts information from delta meta data and local csv files.
func Build(set *meta.Set, caster *ntrip.Caster) (*Config, error) {

	var config Config

	receivers := make(map[string]string)
	for _, r := range set.DeployedReceivers() {
		if t := time.Now().UTC(); t.After(r.End) {
			continue
		}
		receivers[r.Mark] = r.Model
	}

	for _, u := range caster.Users() {
		config.Users = append(config.Users, UserConfig{
			Username: u.Username,
			Password: u.Password,
		})
	}

	groups := make(map[string][]string)
	for _, u := range caster.Users() {
		for _, g := range u.Groups {
			groups[g] = append(groups[g], u.Username)
		}
	}
	for k, v := range groups {
		config.Groups = append(config.Groups, GroupConfig{
			Group: k,
			Users: v,
		})
	}

	for _, m := range caster.Mounts() {
		config.ClientMounts = append(config.ClientMounts, ClientMountConfig{
			Mount:  m.Mount,
			Groups: m.Groups,
		})
	}

	for _, m := range caster.Mounts() {
		mark, ok := set.Mark(m.Mark)
		if !ok {
			continue
		}

		receiver, ok := receivers[m.Mark]
		if !ok {
			continue
		}

		model, ok := caster.Model(receiver)
		if !ok {
			return nil, fmt.Errorf("unknown model: %s", receiver)
		}

		details, ok := caster.Format(m.Details)
		if !ok {
			return nil, fmt.Errorf("unknown format: %s", m.Details)
		}

		config.Mounts = append(config.Mounts, MountConfig{
			Mount:      m.Mount,
			Mark:       mark.Code,
			Name:       mark.Name,
			Country:    m.Country,
			Latitude:   strconv.FormatFloat(mark.Latitude, 'f', 2, 64),
			Longitude:  strconv.FormatFloat(mark.Longitude, 'f', 2, 64),
			Format:     m.Format,
			Details:    strings.Join(details, ","),
			Navigation: m.Navigation,
			Model:      model,
			User:       m.User,
			Address:    m.Address,
		})
	}

	for _, a := range caster.Aliases() {
		config.Aliases = append(config.Aliases, AliasConfig{
			Alias: a.Alias,
			Mount: a.Mount,
		})
	}

	config.Sort()

	return &config, nil
}
