package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

const globalTemplate = `###
### Delivered by puppet
###
# Defines the channel code of the preferred stream used eg by scautopick and
# scrttv. If no component code is given, 'Z' will be used by default.
detecStream = {{.Stream}}

# Defines the location code of the preferred stream used eg by scautopick and
# scrttv.
detecLocid = {{.Location}}
`

type Global struct {
	Stream   string
	Location string
}

func (p Global) Id() string {
	return "global"
}

func (p Global) Template() string {
	return globalTemplate
}

func (p Global) Style() string {
	switch {
	case strings.HasPrefix(p.Stream, "HH"):
		return "broadband"
	case strings.HasPrefix(p.Stream, "EH"):
		return "weak"
	case strings.HasPrefix(p.Stream, "SH"):
		return "weak"
	case strings.HasPrefix(p.Stream, "HN"):
		return "strong"
	case strings.HasPrefix(p.Stream, "BN"):
		return "strong"
	default:
		return ""
	}
}

func (p Global) Rate() string {
	switch {
	case strings.HasPrefix(p.Stream, "B"):
		return "lowrate"
	case strings.HasPrefix(p.Stream, "S"):
		return "lowrate"
	default:
		return ""
	}
}

func (p Global) Key() string {
	switch r := p.Rate(); r {
	case "":
		return fmt.Sprintf("%s_%s", p.Style(), p.Location)
	default:
		return fmt.Sprintf("%s_%s_%s", p.Style(), r, p.Location)
	}
}

func (p Global) Path() string {
	return filepath.Join(p.Id(), fmt.Sprintf("profile_%s", p.Key()))
}
