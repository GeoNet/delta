package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// https://igscb.jpl.nasa.gov/igscb/station/general/antenna.gra

// load prepared information
func loadInfo(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var info Config
	if err := yaml.Unmarshal(b, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

// load antenna diagrams
func loadDiagrams(path string) (map[string]string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var diagrams map[string]string
	if err := yaml.Unmarshal(b, &diagrams); err != nil {
		return nil, err
	}

	return diagrams, nil
}

func loadCountries(path string) (map[string]Country, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var countries map[string]Country
	if err := yaml.Unmarshal(b, &countries); err != nil {
		return nil, err
	}

	return countries, nil
}

func loadAgencies(path string) (map[string]Agency, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var agencies map[string]Agency
	if err := yaml.Unmarshal(b, &agencies); err != nil {
		return nil, err
	}

	return agencies, nil
}

func main() {

	var prepared string
	flag.StringVar(&prepared, "prepared", "config/prepared.yaml", "prepared file")

	var countries string
	flag.StringVar(&countries, "countries", "config/countries.yaml", "countries file")

	var antennas string
	flag.StringVar(&antennas, "antennas", "config/antennas.yaml", "antennas file")

	var agencies string
	flag.StringVar(&agencies, "agencies", "config/agencies.yaml", "agencies file")

	var contact string
	flag.StringVar(&contact, "contact", "GNS Science", "contact agency")

	var responsible string
	flag.StringVar(&responsible, "responsible", "Land Information New Zealand", "responsible agency")

	flag.Parse()

	info, err := loadInfo(prepared)
	if err != nil {
		log.Fatal(err)
	}

	diagrams, err := loadDiagrams(antennas)
	if err != nil {
		log.Fatal(err)
	}

	lookups, err := loadCountries(countries)
	if err != nil {
		log.Fatal(err)
	}

	who, err := loadAgencies(agencies)
	if err != nil {
		log.Fatal(err)
	}

	config := Config{
		PreparedBy:            info.PreparedBy,
		PrimaryDatacentre:     info.PrimaryDatacentre,
		URLForMoreInformation: info.URLForMoreInformation,
		ExtraNotes:            info.ExtraNotes,
		ContactAgency: func() Agency {
			if a, ok := who[contact]; ok {
				return a
			}
			return Agency{}
		}(),
		ResponsibleAgency: func() Agency {
			if a, ok := who[responsible]; ok {
				return a
			}
			return Agency{}
		}(),
		Countries: lookups,
		Diagrams:  diagrams,
	}

	if err := config.Generate(os.Stdout); err != nil {
		log.Fatal(err)
	}
}
