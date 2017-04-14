package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Archive    int32                 `yaml:"archive"`
	Components []string              `yaml:"components"`
	Days       int32                 `yaml:"days"`
	Daemon     string                `yaml:"daemon"`
	Delta      int32                 `yaml:"delta"`
	Dir        string                `yaml:"dir"`
	File       string                `yaml:"file"`
	Forever    int32                 `yaml:"forever"`
	Loggers    []string              `yaml:"loggers"`
	Networks   []string              `yaml:"networks"`
	Sensors    []string              `yaml:"sensors"`
	Stations   []string              `yaml:"stations"`
	Rates      map[string]ConfigRate `yaml:"rates"`
	Step       int32                 `yaml:"step"`
	Settling   string                `yaml:"settling"`
	Extras     []ConfigExtra         `yaml:"extra"`
}

type ConfigRate struct {
	Q      string                 `yaml:"q"`
	Scale  float64                `yaml:"scale"`
	Scales map[string]interface{} `yaml:"scales"`
	High   string                 `yaml:"high"`
	Low    string                 `yaml:"low"`
}

type ConfigExtra struct {
	Locations map[string]map[string]ConfigSite `yaml:"locations"`
	Name      string                           `yaml:"name"`
	NetworkId string                           `yaml:"network_id"`
	StationId string                           `yaml:"station_id"`
}

type ConfigSite struct {
	Gain float64 `yaml:"gain"`
	Q    string  `yaml:"q"`
	Rate float64 `yaml:"rate"`
	High string  `yaml:"high"`
	Low  string  `yaml:"low"`
}

func loadConfig(config string) (map[string]Config, error) {

	cfgs := make(map[string]Config)
	b, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(b, &cfgs); err != nil {
		return nil, err
	}

	return cfgs, nil
}
