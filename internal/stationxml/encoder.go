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
	file, err := os.Create(path) //nolint:gosec // disable G304
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close() // ignore close error as it will likely be related to an existing error
	}()

	if err := r.Write(file, version); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return nil
}
