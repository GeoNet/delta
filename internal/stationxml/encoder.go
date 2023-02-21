package stationxml

import (
	"io"
	"os"
)

// Encode encodes the Root struct using the given Encoder.
func (r Root) MarshalVersion(version string) ([]byte, error) {
	switch version {
	case "1.0":
		return Encode10(r)
	case "1.1":
		return Encode11(r)
	case "1.2":
		return Encode12(r)
	default:
		return Encode10(r)
	}
}

// Write encodes and writes the output from the given Encoder for a Root struct.
func (r Root) Write(wr io.Writer, version string) error {
	res, err := r.MarshalVersion(version)
	if err != nil {
		return err
	}
	if _, err := wr.Write(res); err != nil {
		return err
	}

	return nil
}

// WriteFile encodes and stores the output from the given Encoder for a Root struct.
func (r Root) WriteFile(path string, version string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return r.Write(file, version)
}
