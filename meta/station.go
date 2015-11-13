package meta

import (
	/*
		"bytes"
		"encoding/csv"
		"fmt"
		"os"
		"path/filepath"
		"reflect"
		"strconv"
		"strings"
	*/
	"time"
)

type Station struct {
	Code      string `csv:"Station Code",`
	Name      string `csv:"Station Name",`
	Count     int16
	StartTime time.Time `csv:"Start Time"`
	/*
		Location  string  `csv:"Radio Location",`
		Target    string  `csv:"Radio Target Location",`
		Role      string  `csv:"Radio Role",`
		Model     string  `csv:"Radio Model",`
		Serial    string  `csv:"Radio Serial Number",`
		Polarity  string  `csv:"Antenna Polarity",`
		Frequency float64 `csv:"Frequency Key",`
	*/
}

type Stations []Station

func (s Stations) list() {}

/*

type List interface {
	List()
}

func Strings(list List) string {
	var b bytes.Buffer

	if lines, err := Encode(list); err == nil {
		csv.NewWriter(&b).WriteAll(lines)
	}

	return b.String()
}

func Decode(data [][]string, list List) error {
	// check for correct types
	rv := reflect.ValueOf(list)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("list decode requires a pointer")
	}
	rv = reflect.Indirect(rv)
	if rv.Kind() != reflect.Slice {
		return fmt.Errorf("list decode requires a pointer to a slice")
	}

	// no data ...
	if !(len(data) > 1) {
		return nil
	}

	// where to start ...
	offset := rv.Len()

	// make space ...
	if rv.IsNil() {
		rv.Set(reflect.MakeSlice(rv.Type(), len(data)-1, len(data)-1))
	} else {
		rv.Set(reflect.AppendSlice(rv, reflect.MakeSlice(rv.Type(), len(data)-1, len(data)-1)))
	}

	// skip the header line ...
	for i, n := 1, len(data); i < n; i++ {
		// gather the slot to store the data
		ri := reflect.Indirect(rv.Index(offset + i - 1))
		if ri.Kind() != reflect.Struct {
			return fmt.Errorf("list decode requires a pointer to a slice of structs")
		}
		if ri.NumField() != len(data[i]) {
			return fmt.Errorf("list decode found incorrect number of fields")
		}
		// decode each field, in order ...
		for j := 0; j < ri.NumField(); j++ {
			v := ri.Field(j)
			s := strings.TrimSpace(data[i][j])
			switch v.Kind() {
			case reflect.String:
				v.SetString(s)
			case reflect.Int32:
				i, err := strconv.ParseInt(s, 10, 32)
				if err != nil {
					return err
				}
				v.SetInt(i)
			case reflect.Float32:
				f, err := strconv.ParseFloat(s, 32)
				if err != nil {
					return err
				}
				v.SetFloat(f)
			case reflect.Float64:
				f, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return err
				}
				v.SetFloat(f)
			default:
				t, err := time.Parse(DateTimeFormat, s)
				if err != nil {
					return err
				}
				v.Set(reflect.ValueOf(t))
			}
		}
	}

	return nil
}

func Encode(list List) ([][]string, error) {
	var data [][]string

	rv := reflect.ValueOf(list)
	if rv.Kind() != reflect.Slice {
		return nil, fmt.Errorf("list decode requires a pointer to a slice")
	}
	if rv.IsNil() {
		return data, nil
	}
	for i, n := 0, rv.Len(); i < n; i++ {
		v := rv.Index(i).Interface()
		ri := reflect.ValueOf(v)
		if i == 0 {
			var header []string
			t := reflect.TypeOf(v)
			for j := 0; j < t.NumField(); j++ {
				f := t.Field(j)
				tags := strings.SplitAfter(strings.TrimSpace(f.Tag.Get("csv")), ",")
				if len(tags) > 0 && len(tags[0]) > 0 {
					header = append(header, tags[0])
				} else {
					header = append(header, strings.Title(f.Name))
				}
			}

			data = append(data, header)
		}
		var line []string
		for j := 0; j < ri.NumField(); j++ {
			f := ri.Field(j)
			switch f.Kind() {
			case reflect.String:
				line = append(line, f.String())
			case reflect.Int32:
				line = append(line, strconv.FormatInt(f.Int(), 10))
			case reflect.Float64:
				line = append(line, strconv.FormatFloat(f.Float(), 'g', -1, 64))
			default:
				line = append(line, f.Interface().(time.Time).UTC().Format(DateTimeFormat))
			}
		}
		data = append(data, line)
	}

	return data, nil
}
*/

/*
func LoadStations(path string) (Stations, error) {
	var list Stations

	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	if err := Decode(data, list); err != nil {
		return err
	}

	return nil
}
*/

/*
func LoadLists(dirname, filename string, list List) error {

	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err == nil && filepath.Base(path) == filename {
			if err := LoadList(path, list); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
*/
