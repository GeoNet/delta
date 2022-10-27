package meta

import (
	"testing"
	"time"
)

func TestCitationList(t *testing.T) {
	t.Run("check citations", testListFunc("testdata/citations.csv", &CitationList{
		Citation{
			Key:       "fry2020",
			Author:    "Fry, B., S.-J. McCurrach, K. Gledhill, W. Power, M. Williams, M. Angove, D. Arcas, and C. Moore",
			Title:     "Sensor network warns of stealth tsunamis",
			Published: "Eos",
			Volume:    "101",
			Doi:       MustDoi("https://doi.org/10.1029/2020EO144274"),

			doi:  "https://doi.org/10.1029/2020EO144274",
			year: "2020",
		},
		Citation{
			Key:       "key2002",
			Author:    "A Writer",
			Title:     "A test",
			Year:      2002,
			Published: "Journal of test",
			Volume:    "1",
			Pages:     "1-2",
			Doi:       MustDoi("https://doi.org/10.21420/8TCZ-TV02"),
			Link:      "http://example.com",
			Retrieved: time.Date(2021, time.January, 12, 0, 0, 0, 0, time.UTC),

			doi:       "https://doi.org/10.21420/8TCZ-TV02",
			year:      "2002",
			retrieved: "2021-01-12T00:00:00Z",
		},
		Citation{
			Key: "unknown",
		},
	}))
}
