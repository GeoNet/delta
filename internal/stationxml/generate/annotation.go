package main

import (
	"strings"
)

type Example struct {
	LevelChoice string `xml:",attr,omitempty"`
	Text        string `xml:",innerxml"`
}

func (e Example) Comments(level string) []string {
	if level != e.LevelChoice {
		return nil
	}
	if strings.TrimSpace(e.Text) == "" {
		return nil
	}

	var lines []string
	lines = append(lines, "")
	lines = append(lines, "Example:")
	for _, l := range strings.Split(e.Text, "\n") {
		lines = append(lines, "   "+strings.TrimSpace(l))
	}
	return lines
}

type Warning struct {
	LevelChoice string `xml:",attr,omitempty"`
	Text        string `xml:",cdata"`
}

func (w Warning) Comments(level string) []string {
	if level != w.LevelChoice {
		return nil
	}
	if strings.TrimSpace(w.Text) == "" {
		return nil
	}

	var lines []string
	lines = append(lines, "")
	for _, l := range strings.Split(w.Text, "\n") {
		if v := strings.TrimSpace(l); v != "" {
			switch {
			case !(len(lines) > 1):
				lines = append(lines, "WARNING: "+v)
			default:
				lines = append(lines, v)
			}
		}
	}
	return lines
}

type LevelDesc struct {
	LevelChoice string `xml:",attr.omitempty"`
	Text        string `xml:",cdata"`
}

func (l LevelDesc) Comments(level string) []string {
	if l.LevelChoice != level {
		return nil
	}
	if strings.TrimSpace(l.Text) == "" {
		return nil
	}

	var lines []string
	lines = append(lines, "%L%"+"")
	for _, l := range strings.Split(l.Text, "\n") {
		lines = append(lines, "%L%"+strings.TrimSpace(l))
	}
	return lines
}

type Documentation struct {
	Text string `xml:",cdata"`

	LevelDesc *LevelDesc `xml:"levelDesc,omitempty"`
	Example   *Example   `xml:"example,omitempty"`
	Warning   *Warning   `xml:"warning,omitempty"`
}

func (d Documentation) Level() string {
	if l := d.LevelDesc; l != nil {
		return l.LevelChoice
	}
	return ""
}

func (d Documentation) Comments(level string) []string {
	var lines []string

	if l := d.LevelDesc; l != nil {
		lines = append(lines, l.Comments(level)...)
	}
	if d.Text != "" {
		for _, l := range strings.Split(d.Text, "\n") {
			if v := strings.TrimSpace(l); v != "" {
				lines = append(lines, v)
			}
		}
	}
	if e := d.Example; e != nil {
		lines = append(lines, e.Comments(level)...)
	}
	if w := d.Warning; w != nil {
		lines = append(lines, w.Comments(level)...)
	}

	return lines
}

type Annotation struct {
	Documentation []Documentation `xml:"documentation"`
}

func (a Annotation) Levels() []string {
	var levels []string

	buf := make(map[string]interface{})
	for _, d := range a.Documentation {
		if _, ok := buf[d.Level()]; !ok {
			levels = append(levels, d.Level())
			buf[d.Level()] = true
		}
	}

	return levels
}

func (a Annotation) Comments() []string {
	var comments []string
	for _, l := range a.Levels() {
		for _, d := range a.Documentation {
			comments = append(comments, d.Comments(l)...)
		}
	}
	return comments
}
