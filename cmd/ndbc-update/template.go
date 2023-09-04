package main

import (
	"text/template"
	"time"
)

var tmpl = template.Must(template.New("ndbc").Funcs(
	template.FuncMap{
		"active": func(at time.Time) bool {
			return time.Since(at) < 0
		},
		"now": func() string {
			return time.Now().UTC().Format("2006-01-02T15:04Z")
		},
	}).Parse(`
{{now}}

New Zealand DART Buoy Network
{{range $v := .}}{{if (active $v.End)}}


{{printf "%s/%s -- %s -- %s" $v.Buoy $v.Deployment $v.Region $v.Name}}
-------------------------------
{{printf "WMO Id            : %s" $v.Pid}}
{{printf "Payload Ids       : %s" ($v.Payload $v.Buoy)}}
{{printf "Platform Type     : %s" $v.Platform}}
{{printf "Paroscientific SN : %s" ""}}
{{printf "Deployment Start  : %s" ($v.Start.Format "2006-01-02T15:04Z")}}
{{printf "Water Depth       : %4.0f" $v.Depth}}
{{printf "BPR Drop Position : %8.4f %9.4f" $v.Latitude $v.Longitude}}
{{end}}{{end}}
`))
