package meta

import (
	"sort"
	"strings"
	"time"
)

const (
	citationKey int = iota
	citationAuthor
	citationYear
	citationTitle
	citationPublished
	citationVolume
	citationPages
	citationDoi
	citationLink
	citationRetrieved
	citationLast
)

var citationHeaders Header = map[string]int{
	"Key":       citationKey,
	"Author":    citationAuthor,
	"Year":      citationYear,
	"Title":     citationTitle,
	"Published": citationPublished,
	"Volume":    citationVolume,
	"Pages":     citationPages,
	"DOI":       citationDoi,
	"Link":      citationLink,
	"Retrieved": citationRetrieved,
}

type Citation struct {
	Key       string
	Author    string
	Year      int
	Title     string
	Published string
	Volume    string
	Pages     string
	Doi       Doi
	Link      string
	Retrieved time.Time

	year      string
	doi       string
	retrieved string
}

type CitationList []Citation

func (c CitationList) Len() int           { return len(c) }
func (c CitationList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c CitationList) Less(i, j int) bool { return c[i].Key < c[j].Key }

func (c CitationList) encode() [][]string {
	var data [][]string

	data = append(data, citationHeaders.Columns())

	for _, row := range c {
		data = append(data, []string{
			row.Key,
			row.Author,
			row.year,
			row.Title,
			row.Published,
			row.Volume,
			row.Pages,
			row.doi,
			row.Link,
			row.retrieved,
		})
	}

	return data
}

func (c *CitationList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var citations []Citation

	fields := citationHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		year, err := ParseInt(d[citationYear])
		if err != nil {
			return err
		}

		var retrieved time.Time
		if s := strings.TrimSpace(d[citationRetrieved]); s != "" {
			t, err := time.Parse(DateTimeFormat, s)
			if err != nil {
				return err
			}
			retrieved = t
		}

		var doi Doi
		if s := strings.TrimSpace(d[citationDoi]); s != "" {
			if err := doi.UnmarshalText([]byte(s)); err != nil {
				return err
			}
		}

		citations = append(citations, Citation{
			Key:       strings.TrimSpace(d[citationKey]),
			Author:    strings.TrimSpace(d[citationAuthor]),
			Year:      year,
			Title:     strings.TrimSpace(d[citationTitle]),
			Published: strings.TrimSpace(d[citationPublished]),
			Volume:    strings.TrimSpace(d[citationVolume]),
			Pages:     strings.TrimSpace(d[citationPages]),
			Doi:       doi,
			Link:      strings.TrimSpace(d[citationLink]),
			Retrieved: retrieved,

			year:      strings.TrimSpace(d[citationYear]),
			doi:       strings.TrimSpace(d[citationDoi]),
			retrieved: strings.TrimSpace(d[citationRetrieved]),
		})
	}

	*c = CitationList(citations)

	return nil
}

func LoadCitations(path string) ([]Citation, error) {
	var c []Citation

	if err := LoadList(path, (*CitationList)(&c)); err != nil {
		return nil, err
	}

	sort.Sort(CitationList(c))

	return c, nil
}
