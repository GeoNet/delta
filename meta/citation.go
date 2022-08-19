package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	citationKey int = iota
	citationAuthor
	citationYear
	citationTitle
	citationPublished
	citationDoi
	citationLink
	citationRetrieved
	citationLast
)

type Citation struct {
	Key       string
	Author    string
	Year      int
	Title     string
	Doi       string
	Published string
	Link      string
	Retrieved time.Time

	retrieved string
}

type CitationList []Citation

func (c CitationList) Len() int           { return len(c) }
func (c CitationList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c CitationList) Less(i, j int) bool { return c[i].Key < c[j].Key }

func (c CitationList) encode() [][]string {
	data := [][]string{{"Key", "Author", "Year", "Title", "Published", "DOI", "Link", "Retrieved"}}
	for _, v := range c {
		data = append(data, []string{
			v.Key,
			v.Author,
			strconv.Itoa(v.Year),
			v.Title,
			v.Published,
			v.Doi,
			v.Link,
			v.retrieved,
		})
	}
	return data
}

func (c *CitationList) decode(data [][]string) error {
	var citations []Citation
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != citationLast {
				return fmt.Errorf("incorrect number of citation fields")
			}

			var year int
			if s := strings.TrimSpace(d[citationYear]); s != "" {
				i, err := strconv.Atoi(s)
				if err != nil {
					return err
				}
				year = i
			}

			var retrieved time.Time
			if s := strings.TrimSpace(d[citationRetrieved]); s != "" {
				t, err := time.Parse(DateTimeFormat, s)
				if err != nil {
					return err
				}
				retrieved = t
			}

			citations = append(citations, Citation{
				Key:       strings.TrimSpace(d[citationKey]),
				Author:    strings.TrimSpace(d[citationAuthor]),
				Year:      year,
				Title:     strings.TrimSpace(d[citationTitle]),
				Published: strings.TrimSpace(d[citationPublished]),
				Doi:       strings.TrimSpace(d[citationDoi]),
				Link:      strings.TrimSpace(d[citationLink]),
				Retrieved: retrieved,

				retrieved: strings.TrimSpace(d[citationRetrieved]),
			})
		}

		*c = CitationList(citations)
	}
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
