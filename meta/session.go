package meta

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	sessionMark = iota
	sessionOperator
	sessionAgency
	sessionModel
	sessionSatelliteSystem
	sessionInterval
	sessionElevationMask
	sessionHeaderComment
	sessionFormat
	sessionStart
	sessionEnd
	sessionLast
)

var sessionHeaders Header = map[string]int{
	"Mark":             sessionMark,
	"Operator":         sessionOperator,
	"Agency":           sessionAgency,
	"Model":            sessionModel,
	"Satellite System": sessionSatelliteSystem,
	"Interval":         sessionInterval,
	"Elevation Mask":   sessionElevationMask,
	"Header Comment":   sessionHeaderComment,
	"Format":           sessionFormat,
	"Start Date":       sessionStart,
	"End Date":         sessionEnd,
}

type Session struct {
	Span

	Mark            string
	Operator        string
	Agency          string
	Model           string
	SatelliteSystem string
	Interval        time.Duration
	ElevationMask   float64
	HeaderComment   string
	Format          string

	elevationMask string // shadow variable to maintain formatting
}

func (s Session) Less(session Session) bool {
	switch {
	case s.Mark < session.Mark:
		return true
	case s.Mark > session.Mark:
		return false
	case s.Model < session.Model:
		return true
	case s.Model > session.Model:
		return false
	case s.Interval < session.Interval:
		return true
	case s.Interval > session.Interval:
		return false
	case s.Start.Before(session.Start):
		return true
	case s.Start.After(session.Start):
		return false
	default:
		return false
	}
}

type SessionList []Session

func (s SessionList) Len() int           { return len(s) }
func (s SessionList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SessionList) Less(i, j int) bool { return s[i].Less(s[j]) }

func (s SessionList) encode() [][]string {
	var data [][]string

	data = append(data, sessionHeaders.Columns())

	for _, row := range s {
		data = append(data, []string{
			strings.TrimSpace(row.Mark),
			strings.TrimSpace(row.Operator),
			strings.TrimSpace(row.Agency),
			strings.TrimSpace(row.Model),
			strings.TrimSpace(row.SatelliteSystem),
			strings.TrimSpace(row.Interval.String()),
			strings.TrimSpace(row.elevationMask),
			strings.TrimSpace(row.HeaderComment),
			strings.TrimSpace(row.Format),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (s *SessionList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var sessions []Session

	fields := sessionHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		interval, err := time.ParseDuration(d[sessionInterval])
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[sessionStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[sessionEnd])
		if err != nil {
			return err
		}

		mask, err := strconv.ParseFloat(d[sessionElevationMask], 64)
		if err != nil {
			return err
		}

		sessions = append(sessions, Session{
			Mark:            strings.TrimSpace(d[sessionMark]),
			Operator:        strings.TrimSpace(d[sessionOperator]),
			Agency:          strings.TrimSpace(d[sessionAgency]),
			Model:           strings.TrimSpace(d[sessionModel]),
			SatelliteSystem: strings.TrimSpace(d[sessionSatelliteSystem]),
			Interval:        interval,
			ElevationMask:   mask,
			HeaderComment:   strings.TrimSpace(d[sessionHeaderComment]),
			Format:          strings.TrimSpace(d[sessionFormat]),
			Span: Span{
				Start: start,
				End:   end,
			},

			elevationMask: strings.TrimSpace(d[sessionElevationMask]),
		})
	}

	*s = SessionList(sessions)

	return nil
}

func LoadSessions(path string) ([]Session, error) {
	var s []Session

	if err := LoadList(path, (*SessionList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(SessionList(s))

	return s, nil
}
