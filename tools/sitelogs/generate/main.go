package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// https://igscb.jpl.nasa.gov/igscb/station/general/antenna.gra

func main() {

	var config string
	flag.StringVar(&config, "config", "config/config.yaml", "config file")

	flag.Parse()

	b, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatal(err)
	}

	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		log.Fatal(err)
	}

	if err := c.Generate(os.Stdout); err != nil {
		log.Fatal(err)
	}
}
