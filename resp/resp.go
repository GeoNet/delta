package resp

import (
	"embed"
	"encoding/xml"
	"fmt"
	"io/fs"

	"github.com/GeoNet/delta/internal/stationxml/v1.1"
)

//TODO: add embed when populated
var locations = []string{"files", "auto"}

//go:embed files/*.xml
//go:embed auto/*.xml
var files embed.FS

// Resp returns a pointer to an embeded stationxml Response if present.
func Lookup(response string) (*stationxml.ResponseType, error) {

	for _, l := range locations {
		names, err := fs.Glob(files, fmt.Sprintf("%s/%s.xml", l, response))
		if err != nil {
			return nil, err
		}

		for _, name := range names {
			data, err := fs.ReadFile(files, name)
			if err != nil {
				return nil, err
			}

			var resp stationxml.ResponseType
			if err := xml.Unmarshal(data, &resp); err != nil {
				return nil, err
			}
			return &resp, nil
		}
	}

	return nil, nil
}
