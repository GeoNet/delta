package main

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/GeoNet/delta/internal/metadb"
	"github.com/GeoNet/delta/meta"
)

type Plots struct {
	Configs map[string]ConfigPages `yaml:"plots"`
}

type ConfigPages struct {
	Filename string       `yaml:"filename"`
	Pages    []ConfigPage `yaml:"pages"`
}

type ConfigPage struct {
	Type    string
	Options Options
}

type Options struct {
	Reversed  string            `yaml:"reversed"`
	Label     string            `yaml:"label"`
	Base      string            `yaml:"base"`
	Png       string            `yaml:"png"`
	Id        string            `yaml:"id"`
	Tag       string            `yaml:"tag"`
	Rrd       string            `yaml:"rrd"`
	Networks  []string          `yaml:"networks"`
	Locations []string          `yaml:"locations"`
	Stations  []string          `yaml:"stations"`
	Location  string            `yaml:"location"`
	Override  map[string]string `yaml:"override"`
	Bands     []string          `yaml:"bands"`
	Orients   []string          `yaml:"orients"`
	Sensors   []string          `yaml:"sensors"`
	Detide    int               `yaml:"detide"`
	Reference int               `yaml:"reference"`
	Excludes  []string          `yaml:"exclude"`
	Thumb     string            `yaml:"thumb"`
	Gains     map[string]Gain   `yaml:"gains"`
	Streams   []OptionStream    `yaml:"streams"`
}

type Gain map[string]string

type OptionStream struct {
	Title    string `yaml:"title"`
	Station  string `yaml:"sta"`
	Location string `yaml:"loc"`
	Channel  string `yaml:"cha"`
	Network  string `yaml:"net"`
}

func StationChannel(s meta.Station, c metadb.Channel) OptionStream {
	return OptionStream{
		Network:  s.Network,
		Station:  s.Code,
		Location: c.Location,
		Channel:  c.Code,
	}
}

func (cp ConfigPage) Png(stream OptionStream, def string) string {
	png := func() string {
		if cp.Options.Png != "" {
			return cp.Options.Png
		}
		return def
	}()

	label := func() string {
		if cp.Options.Label != "" {
			return cp.Options.Label
		}
		return "trace"
	}()

	png = strings.Replace(png, "%t", strings.ToLower(stream.Title), -1)
	png = strings.Replace(png, "%n", strings.ToLower(stream.Network), -1)
	png = strings.Replace(png, "%s", strings.ToLower(stream.Station), -1)
	png = strings.Replace(png, "%l", strings.ToLower(stream.Location), -1)
	png = strings.Replace(png, "%c", strings.ToLower(stream.Channel), -1)
	png = strings.Replace(png, "%x", label, -1)

	png = strings.Replace(png, "%a", func() string {
		switch stream.Channel {
		case "HDF":
			return "acoustic"
		case "HDH":
			return "hydrophone"
		default:
			return "seismic"
		}
	}(), -1)

	return filepath.Join(cp.Options.Base, png)
}

func (cp ConfigPage) Id(stream OptionStream, def string) string {
	id := func() string {
		if cp.Options.Id != "" {
			return cp.Options.Id
		}
		return def
	}()

	tag := func() string {
		if cp.Options.Tag != "" {
			return cp.Options.Tag
		}
		return "earthquake"
	}()

	id = strings.Replace(id, "%t", tag, -1)
	id = strings.Replace(id, "%n", strings.ToLower(stream.Network), -1)
	id = strings.Replace(id, "%s", strings.ToLower(stream.Station), -1)
	id = strings.Replace(id, "%l", strings.ToLower(stream.Location), -1)
	id = strings.Replace(id, "%c", strings.ToLower(stream.Channel), -1)

	return id
}

func (cp ConfigPage) Rrd(stream OptionStream, def string) string {

	rrd := func() string {
		if def != "" {
			return def
		}
		return "/%s.%n/%s.%l-%c.%n.rrd"
	}()

	rrd = strings.Replace(rrd, "%t", strings.ToLower(stream.Title), -1)
	rrd = strings.Replace(rrd, "%n", strings.ToLower(stream.Network), -1)
	rrd = strings.Replace(rrd, "%s", strings.ToLower(stream.Station), -1)
	rrd = strings.Replace(rrd, "%l", strings.ToLower(stream.Location), -1)
	rrd = strings.Replace(rrd, "%c", strings.ToLower(stream.Channel), -1)

	return func() string {
		if def != "" {
			return filepath.Join(cp.Options.Rrd, rrd)
		}
		return rrd
	}()
}

func (cp ConfigPage) Gain(stream OptionStream) string {
	if lookup, ok := cp.Options.Gains[stream.Station]; ok {
		if gain, ok := lookup[stream.Location]; ok {
			return gain
		}
	}
	return ""
}
func (cp ConfigPage) Auto(stream OptionStream) string {
	if cp.Gain(stream) != "" {
		return "false"
	}
	return "yes"
}

func (c Options) GetLocation(sta string) string {
	if loc, ok := c.Override[sta]; ok {
		return loc
	}
	return c.Location
}

func (c Options) Channels() []string {
	var chans []string
	for _, b := range c.Bands {
		for _, s := range c.Sensors {
			for _, o := range c.Orients {
				chans = append(chans, b+s+o)
			}
		}
	}

	sort.Strings(chans)

	return chans
}

func ReadPlots(path string) (*Plots, error) {

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var plots Plots
	if err := yaml.Unmarshal(b, &plots); err != nil {
		return nil, err
	}

	return &plots, nil
}
