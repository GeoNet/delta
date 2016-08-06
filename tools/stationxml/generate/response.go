package main

type Sensor struct {
	Sensors  []string `yaml:"sensors"`
	Filters  []string `yaml:"filters"`
	Channels string   `yaml:"channels"`
	Reversed bool     `yaml:"reversed"`
	Match    string   `yaml:"match"`
	Skip     string   `yaml:"skip"`
}

type Datalogger struct {
	Dataloggers   []string `yaml:"dataloggers"`
	Type          string   `yaml:"type"`
	Label         string   `yaml:"label"`
	SampleRate    float64  `yaml:"samplerate"`
	Frequency     float64  `yaml:"frequency"`
	StorageFormat string   `yaml:"storageformat"`
	ClockDrift    float64  `yaml:"clockdrift"`
	Filters       []string `yaml:"filters"`
	Reversed      bool     `yaml:"reversed"`
	Match         string   `yaml:"match"`
	Skip          string   `yaml:"skip"`
}

type Response struct {
	Sensors     []Sensor     `yaml:"sensors"`
	Dataloggers []Datalogger `yaml:"dataloggers"`
}
