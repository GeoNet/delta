package main

import (
	"sort"
)

type PAZ struct {
	Code  string       `yaml:"code"`
	Type  string       `yaml:"type"`
	Notes string       `yaml:"notes"`
	Poles []Complex128 `yaml:"poles"`
	Zeros []Complex128 `yaml:"zeros"`
}

func (p PAZ) PzTransferFunction() string {
	switch p.Code {
	case "A":
		return "PZFunctionLaplaceRadiansPerSecond"
	case "B":
		return "PZFunctionLaplaceHertz"
	case "D":
		return "PZFunctionLaplaceZTransform"
	default:
		return "PZFunctionUnknown"
	}
}

type pazMap map[string]PAZ

func (p pazMap) Keys() []string {
	var keys []string
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (p pazMap) Values() []PAZ {
	var values []PAZ
	for _, k := range p.Keys() {
		values = append(values, p[k])
	}
	return values
}
