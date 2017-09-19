package main

import (
	"sort"
	"strings"
	"time"

	"github.com/GeoNet/delta/internal/metadb"
)

type Channels []metadb.Channel

func (c Channels) Len() int      { return len(c) }
func (c Channels) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c Channels) Less(i, j int) bool {
	switch {
	case c[i].Location < c[j].Location:
		return true
	case c[i].Location > c[j].Location:
		return false
	default:
		return c[i].SampleRate > c[j].SampleRate
	}
}

func (cp ConfigPage) Wave(base string) ([]Plot, error) {
	var plots []Plot

	db := metadb.NewMetaDB(base)

	chans := make(map[string]interface{})
	for _, c := range cp.Options.Channels() {
		chans[c] = true
	}

	for _, net := range cp.Options.Networks {

		stns, err := db.NetworkStation(net)
		if err != nil {
			return nil, err
		}

		for _, s := range stns {
			var list []Stream

			chas, err := db.Channels(s.Code)
			if err != nil {
				return nil, err
			}

			sort.Sort(Channels(chas))

			for _, c := range chas {
				if cp.Options.GetLocation(s.Code) != c.Location {
					continue
				}
				if !strings.HasSuffix(c.Code, "H") {
					continue
				}
				if time.Now().After(c.End) {
					continue
				}
				if _, ok := chans[c.Code]; !ok {
					continue
				}

				list = append(list, Stream{
					Id:      cp.Id(StationChannel(s, c), "%n_%s_%l_%c"),
					Srcname: strings.ToUpper(cp.Id(StationChannel(s, c), "%s/%l-%c")),
					Rrd:     cp.Rrd(StationChannel(s, c), "amplitude/wave/%s.%n/%s.%l-%c.%n.rrd"),
				})
			}

			plots = append(plots, Plot{
				Caption: s.Name,
				Id: cp.Id(OptionStream{
					Station: s.Code,
					Network: s.Network,
				}, "height-%s-%n"),
				Png: cp.Png(OptionStream{
					Station: s.Code,
					Network: net,
				}, "tsunami/%n/%s/height.png"),
				Pdf: cp.Png(OptionStream{
					Station: s.Code,
					Network: net,
				}, "tsunami/%n/%s/height.pdf"),
				Streams: list,
			})
		}
	}

	for _, net := range cp.Options.Networks {

		stns, err := db.NetworkStation(net)
		if err != nil {
			return nil, err
		}

		for _, s := range stns {
			var list []Stream

			chas, err := db.Channels(s.Code)
			if err != nil {
				return nil, err
			}

			sort.Sort(Channels(chas))

			for _, c := range chas {
				if cp.Options.GetLocation(s.Code) != c.Location {
					continue
				}
				if !strings.HasSuffix(c.Code, "T") {
					continue
				}
				if time.Now().After(c.End) {
					continue
				}
				if _, ok := chans[c.Code]; !ok {
					continue
				}

				list = append(list, Stream{
					Id:      cp.Id(StationChannel(s, c), "%n_%s_%l_%c"),
					Srcname: strings.ToUpper(cp.Id(StationChannel(s, c), "%s/%l-%c")),
					Rrd:     cp.Rrd(StationChannel(s, c), "amplitude/wave/%s.%n/%s.%l-%c.%n.rrd"),
				})
			}

			plots = append(plots, Plot{
				Caption: s.Name,
				Id: cp.Id(OptionStream{
					Station: s.Code,
					Network: s.Network,
				}, "wave-%s-%n"),
				Png: cp.Png(OptionStream{
					Station: s.Code,
					Network: net,
				}, "tsunami/%n/%s/wave.png"),
				Pdf: cp.Png(OptionStream{
					Station: s.Code,
					Network: net,
				}, "tsunami/%n/%s/wave.pdf"),
				Streams: list,
			})
		}
	}

	return plots, nil
}
