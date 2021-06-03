package main

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

var generateTemplate = `

{{ $b := . }}var Responses []Response = []Response{
{{ range $l, $r := $b.ResponseMap }}    Response{
	Name: "{{$l}}",
        Sensors: []Sensor{
  {{ range $v := $r.Sensors }}    Sensor{
                SensorList: []string{{"{"}}{{range $s := $v.Sensors}}"{{$s}}",{{end}}{{"}"}},
                FilterList: []string{{"{"}}{{range $s := $v.Filters}}"{{$s}}",{{end}}{{"}"}},
		Stages: []ResponseStage{{"{"}}{{range $s := $v.Filters}}{{with $f := $b.Filter $s}}
		{{ range $v := $f }}ResponseStage{
		Type: "{{$v.Type}}",
		Lookup: "{{$v.Lookup}}",
		Filter: "{{$s}}",{{if eq $v.Type "paz" }}{{with $b.PAZ $v.Lookup}}
		StageSet: PAZ{
			Name: "{{$v.Lookup}}",
			Code: {{.PzTransferFunction}},
			Type: "{{.Type}}",{{if .Notes }}
			Notes: "{{.Notes|escape}}",{{end}}{{if .Poles }}
			Poles: []complex128{{"{"}}{{ range $p := .Poles}}{{ $p }},{{end}}{{"}"}},{{ end }}{{if .Zeros }}
			Zeros: []complex128{{"{"}}{{ range $z := .Zeros}}{{ $z }},{{end}}{{"}"}},{{ end }}
		},{{end}}{{end}}{{if eq $v.Type "a2d" }}{{with $b.PAZ $v.Lookup}}
		StageSet: A2D{
			Name: "{{$v.Lookup}}",
			Code: {{.PzTransferFunction}},
			Type: "{{.Type}}",{{if .Notes }}
			Notes: "{{.Notes|escape}}",{{end}}
		},{{end}}{{end}}{{if eq $v.Type "fir" }}{{with $b.FIR $v.Lookup}}
		StageSet: FIR{
			Name: "{{$v.Lookup}}",
			Causal: {{.Causal}},
			Symmetry: {{.SymmetryLookup}},
			Decimation: {{.Decimation}},
			Gain: {{.Gain}},{{if .Notes }}
			Notes: &[]string{"{{.Notes|escape}}"}[0],{{end}}
			Factors: []float64{{"{"}}{{ range $z := .Factors}}{{ $z }},{{end}}{{"}"}},
		},{{end}}{{end}}{{if eq $v.Type "poly" }}{{with $b.Polynomial $v.Lookup}}
		StageSet: Polynomial{
			Name: "{{$v.Lookup}}",
			Gain: {{.Gain}},
			ApproximationType: {{.ApproximationTypeLookup}},
			FrequencyLowerBound: {{.FrequencyLowerBound}},
			FrequencyUpperBound: {{.FrequencyUpperBound}},
			ApproximationLowerBound: {{.ApproximationLowerBound}},
			ApproximationUpperBound: {{.ApproximationUpperBound}},
			MaximumError: {{.MaximumError}},{{if .Notes }}
			Notes: &[]string{"{{.Notes|escape}}"}[0],{{end}}
			Coefficients: []Coefficient{{"{"}}{{ range $z := .Coefficients}}Coefficient{Value: {{ $z }}{{"}"}},{{end}}{{"}"}},
		},{{end}}{{end}}
		Frequency: {{$v.Frequency}},
		SampleRate: {{rate $v.SampleRate}},
		Decimate: {{$v.Decimate}},
		Gain: {{$v.Gain}},
		//Scale: {{$v.Scale}},
		Correction: {{$v.Correction}},
		Delay: {{$v.Delay}},
		InputUnits: "{{$v.InputUnits}}",
		OutputUnits: "{{$v.OutputUnits}}",
		},{{end}}{{end}}{{end}}
    		{{"}"}},
                Channels: "{{$v.Channels}}",
                Reversed: {{$v.Reversed}},
        },{{end}}
        },
        Dataloggers: []Datalogger{
  {{ range $v := $r.Dataloggers }}    Datalogger{
                DataloggerList: []string{{"{"}}{{range $s := $v.Dataloggers}}"{{$s}}",{{end}}{{"}"}},
                Type: "{{$v.Type}}",
                Label: "{{$v.Label}}",
                SampleRate: {{rate $v.SampleRate}},
                Frequency: {{$v.Frequency}},
                StorageFormat: "{{$v.StorageFormat}}",
                ClockDrift: {{$v.ClockDrift}},
                FilterList: []string{{"{"}}{{range $s := $v.Filters}}"{{$s}}",{{end}}{{"}"}},
		Stages: []ResponseStage{{"{"}}{{range $s := $v.Filters}}{{with $f := $b.Filter $s}}
		{{ range $v := $f }}ResponseStage{
		Type: "{{$v.Type}}",
		Lookup: "{{$v.Lookup}}",
		Filter: "{{$s}}",{{if eq $v.Type "paz" }}{{with $b.PAZ $v.Lookup}}
		StageSet: PAZ{
			Name: "{{$v.Lookup}}",
			Code: {{.PzTransferFunction}},
			Type: "{{.Type}}",{{if .Notes }}
			Notes: "{{.Notes|escape}}",{{end}}{{if .Poles }}
			Poles: []complex128{{"{"}}{{ range $p := .Poles}}{{ $p }},{{end}}{{"}"}},{{ end }}{{if .Zeros }}
			Zeros: []complex128{{"{"}}{{ range $z := .Zeros}}{{ $z }},{{end}}{{"}"}},{{ end }}
		},{{end}}{{end}}{{if eq $v.Type "a2d" }}{{with $b.PAZ $v.Lookup}}
		StageSet: A2D{
			Name: "{{$v.Lookup}}",
			Code: {{.PzTransferFunction}},
			Type: "{{.Type}}",{{if .Notes }}
			Notes: "{{.Notes|escape}}",{{end}}
		},{{end}}{{end}}{{if eq $v.Type "fir" }}{{with $b.FIR $v.Lookup}}
		StageSet: FIR{
			Name: "{{$v.Lookup}}",
			Causal: {{.Causal}},
			Symmetry: {{.SymmetryLookup}},
			Decimation: {{.Decimation}},
			Gain: {{.Gain}},{{if .Notes }}
			Notes: &[]string{"{{.Notes|escape}}"}[0],{{end}}
			Factors: []float64{{"{"}}{{ range $z := .Factors}}{{ $z }},{{end}}{{"}"}},
		},{{end}}{{end}}{{if eq $v.Type "poly" }}{{with $b.Polynomial $v.Lookup}}
		StageSet: Polynomial{
			Name: "{{$v.Lookup}}",
			Gain: {{.Gain}},
			ApproximationType: {{.ApproximationTypeLookup}},
			FrequencyLowerBound: {{.FrequencyLowerBound}},
			FrequencyUpperBound: {{.FrequencyUpperBound}},
			ApproximationLowerBound: {{.ApproximationLowerBound}},
			ApproximationUpperBound: {{.ApproximationUpperBound}},
			MaximumError: {{.MaximumError}},{{if .Notes }}
			Notes: &[]string{"{{.Notes|escape}}"}[0],{{end}}
			Coefficients: []Coefficient{{"{"}}{{ range $z := .Coefficients}}Coefficient{Value: {{ $z }}{{"}"}},{{end}}{{"}"}},
		},{{end}}{{end}}
		Frequency: {{$v.Frequency}},
		{{if $v.InputRate}}InputRate: {{$v.InputRate}},
{{end }}                SampleRate: {{rate $v.SampleRate}},
		Decimate: {{if eq $v.Type "fir"}}{{with $b.FIR $v.Lookup}}{{.Decimation}}{{end}}{{else}}{{$v.Decimate}}{{end}},
		Gain: {{$v.Gain}},
		//Scale: {{$v.Scale}},{{if eq $v.Type "fir"}}{{with $b.FIR $v.Lookup}}{{if and (eq $v.Correction 0.0) (gt .Decimation 1.0)}}
		Correction: {{.Correction (rate $v.SampleRate)}},
		Delay: {{.Correction (rate $v.SampleRate)}},{{else}}
		Correction: {{$v.Correction}},
		Delay: {{$v.Correction}},{{end}}{{end}}{{end}}
		InputUnits: "{{$v.InputUnits}}",
		OutputUnits: "{{$v.OutputUnits}}",
		},{{end}}{{end}}{{end}}
    		{{"}"}},
                Reversed: {{$v.Reversed}},
        },{{end}}
        },
  },{{end}}
{{"}"}}
`

//{{with $b.FIR $v.Lookup}}{{if and (gt .Decimation 1) (eq $v.Correction 0.0)}}{{.Correction $v.SampleRate}}{{else}}{{$v.Correction}}{{end}}{{end}}{{else}}{{$v.Correction}}{{end}},
//{{if eq $v.Type "fir"}}{{with $b.FIR $v.Lookup}}{{if gt .Decimation 1}}Correction: {{if eq $v.Correction 0.0}}{{.Correction $v.SampleRate}}{{else}}{{$v.Correction}}{{end}},
//{{if gt .Decimation 1}}Delay: {{if eq $v.Correction 0.0}}{{.Correction $v.SampleRate}}{{else}}{{$v.Correction}}{{end}},{{else}}Correction: {{$v.Correction}},
//Delay: {{$v.Delay}},{{end}}

type Generate struct {
	ResponseMap responseMap
	FilterMap   filterMap
	PazMap      pazMap
	FirMap      firMap
	PolyMap     polynomialMap
}

func (g Generate) Filter(filter string) *[]ResponseStage {
	if f, ok := g.FilterMap[filter]; ok {
		return &f
	}
	return nil
}

func (g Generate) PAZ(paz string) *PAZ {
	if p, ok := g.PazMap[paz]; ok {
		return &p
	}
	return nil
}

func (g Generate) FIR(fir string) *FIR {
	if f, ok := g.FirMap[fir]; ok {
		return &f
	}
	return nil
}

func (g Generate) Polynomial(poly string) *Polynomial {
	if p, ok := g.PolyMap[poly]; ok {
		return &p
	}
	return nil
}

func (g Generate) Validate() error {

	// check responses
	for k, r := range g.ResponseMap {
		// check dataloger responses
		for _, d := range r.Dataloggers {
			label := fmt.Sprintf("\"%s\" [%s]", k, d.Label)

			// check sampling rates
			var r float64
			for _, f := range d.Filters {
				// look for a configured map
				m, ok := g.FilterMap[f]
				if !ok {
					return fmt.Errorf("invalid filter %s: %s", label, f)
				}
				// choose the last sample rate
				for _, s := range m {
					r = s.SampleRate
				}
			}
			if r != d.SampleRate {
				return fmt.Errorf("invalid sample rate %s: found %v, expected %v", label, r, d.SampleRate)
			}
		}
	}
	return nil
}

func (g Generate) generate(w io.Writer) error {

	if err := g.Validate(); err != nil {
		return err
	}

	t, err := template.New("generate").Funcs(
		template.FuncMap{
			"escape": func(s string) string { return strings.Join(strings.Fields(s), " ") },
			"rate": func(f float64) float64 {
				if f < 0.0 {
					return -1.0 / f
				}
				return f
			},
		},
	).Parse(generateTemplate)
	if err != nil {
		return err
	}
	if err := t.Execute(w, g); err != nil {
		return err
	}

	return nil
}
