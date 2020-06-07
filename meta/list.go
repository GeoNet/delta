package meta

import (
	"bytes"
	"encoding/csv"
	"os"
	"path/filepath"
	"sort"
)

// ListEncoder is an interface for encoding a type into a slice of string slices,
// suitable for storing in a CSV file.
type ListEncoder interface {
	encode() [][]string
}

// ListDecoder is an interface for decoding a slice of string slices into a type,
// suitable for reading from a CSV file.
type ListDecoder interface {
	decode([][]string) error
}

// List is an interface the encapsulates the ListEncoder, ListDecoder and sort.Interface interfaces,
// suitable for reading and writing CSV files.
type List interface {
	ListEncoder
	ListDecoder

	sort.Interface
}

// MarshalList converts a type that implements the ListEncoder interface into a byte slice,
// or returns an error otherwise.
func MarshalList(l ListEncoder) ([]byte, error) {
	var b bytes.Buffer

	if err := csv.NewWriter(&b).WriteAll(EncodeList(l)); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// UnmarshalList decodes a byte slice into a type that implements the ListDecoder interface,
// it returns a non empty error otherwise.
func UnmarshalList(b []byte, l ListDecoder) error {

	v, err := csv.NewReader(bytes.NewBuffer(b)).ReadAll()
	if err != nil {
		return err
	}
	if err := DecodeList(v, l); err != nil {
		return err
	}

	return nil
}

// DecodeList converts a slice of string slices into a type that implements the ListDecoder interface,
// it returns a non empty error otherwise.
func DecodeList(data [][]string, l ListDecoder) error {
	return l.decode(data)
}

// EncodeList converts a type that implements the ListEncoder interface into a slice of string slices.
func EncodeList(l ListEncoder) [][]string {
	return l.encode()
}

// LoadList reads and decodes a file into a type that implements the ListDecoder interface,
// it returns a non empty error otherwise.
func LoadList(path string, l ListDecoder) error {

	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return err
	}

	if err := DecodeList(data, l); err != nil {
		return err
	}

	return nil
}

// StoreList encodes and writes into a file the contents of a type that implements the ListEncoder interface,
// it returns a non empty error otherwise.
func StoreList(path string, l ListEncoder) error {

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = csv.NewWriter(file).WriteAll(EncodeList(l))
	if err != nil {
		return err
	}

	return nil
}
