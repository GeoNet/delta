package ntrip

import (
	"fmt"
	"strings"
)

const (
	modelModel int = iota
	modelName
	modelLast
)

// Model represents a set of GNSS receiver model name aliases.
type Model struct {
	Model string
	Name  string
}

func (m *Model) decode(row []string) error {
	if l := len(row); l != modelLast {
		return fmt.Errorf("incorrect \"model\" \"%s\": found %d items, expected %d", strings.Join(row, ","), l, modelLast)
	}

	*m = Model{
		Model: strings.TrimSpace(row[modelModel]),
		Name:  strings.TrimSpace(row[modelName]),
	}

	return nil
}

func (m Model) encode() []string {
	var row []string

	row = append(row, m.Model)
	row = append(row, m.Name)

	return row
}

// Models represents a set of GNSS receiver model aliases.
type Models []Model

// ReadModels extracts a map of GNNS receiver models from a csv file.
func ReadModels(path string) (map[string]string, error) {
	var models Models
	if err := ReadFile(path, &models); err != nil {
		return nil, err
	}

	res := make(map[string]string)
	for _, m := range models {
		if _, ok := res[m.Model]; ok {
			return nil, fmt.Errorf("duplicate model name: %s", m.Model)
		}
		res[m.Model] = m.Name
	}

	return res, nil
}

// Header returns a slice of expected csv file header values.
func (m Models) Header() []string {
	return []string{"#Model", "Name"}
}

// Fields returns the number of expected fields in a csv file representation.
func (m Models) Fields() int {
	return modelLast
}

// Decode extracts model alias information from string slices as expected from csv files.
func (m *Models) Decode(data [][]string) error {
	for _, v := range data {
		var model Model
		if err := model.decode(v); err != nil {
			return err
		}
		*m = append(*m, model)
	}
	return nil
}

// Encode builds model alias information as expected in csv files.
func (m Models) Encode() [][]string {
	var items [][]string

	for _, model := range m {
		items = append(items, model.encode())
	}

	return items
}
