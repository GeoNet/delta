package metadb

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/GeoNet/delta/meta"
	"github.com/hashicorp/go-memdb"
)

type MetaDB struct {
	*memdb.MemDB
}

type IntegerFieldIndex struct {
	Field string
}

func integerIndex(v int64) string {
	return fmt.Sprintf("%021d", v)
}

func (i *IntegerFieldIndex) FromObject(obj interface{}) (bool, []byte, error) {
	v := reflect.ValueOf(obj)
	v = reflect.Indirect(v) // Dereference the pointer if any

	fv := v.FieldByName(i.Field)
	if !fv.IsValid() {
		return false, nil,
			fmt.Errorf("field '%s' for %#v is invalid", i.Field, obj)
	}

	// Add the null character as a terminator
	out := integerIndex(fv.Int()) + "\x00"

	return true, []byte(out), nil
}

func (i *IntegerFieldIndex) FromArgs(args ...interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("must provide only a single argument")
	}

	var out string
	switch args[0].(type) {
	case int:
		out = integerIndex(int64(args[0].(int))) + "\x00"
	case int32:
		out = integerIndex(int64(args[0].(int32))) + "\x00"
	case int64:
		out = integerIndex(args[0].(int64)) + "\x00"
	default:
		return nil, fmt.Errorf("argument must be an int: %#v", args[0])
	}

	return []byte(out), nil
}

func (i *IntegerFieldIndex) PrefixFromArgs(args ...interface{}) ([]byte, error) {
	val, err := i.FromArgs(args...)
	if err != nil {
		return nil, err
	}

	// Strip the null terminator, the rest is a prefix
	n := len(val)
	if n > 0 {
		return val[:n-1], nil
	}
	return val, nil
}

type TimeFieldIndex struct {
	Field string
}

func (t *TimeFieldIndex) FromObject(obj interface{}) (bool, []byte, error) {
	v := reflect.ValueOf(obj)
	v = reflect.Indirect(v) // Dereference the pointer if any

	fv := v.FieldByName(t.Field)
	if !fv.IsValid() {
		return false, nil,
			fmt.Errorf("field '%s' for %#v is invalid", t.Field, obj)
	}

	arg, ok := fv.Interface().(time.Time)
	if !ok {
		return false, nil,
			fmt.Errorf("field '%s' for %#v is invalid", t.Field, arg)
	}

	// Add the null character as a terminator
	out := arg.Format(meta.DateTimeFormat) + "\x00"

	return true, []byte(out), nil
}

func (t *TimeFieldIndex) FromArgs(args ...interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("must provide only a single argument")
	}

	var out string
	switch arg := args[0].(type) {
	case time.Time:
		out = arg.Format(meta.DateTimeFormat) + "\x00"
	default:
		return nil, fmt.Errorf("argument must be a time: %#v", args[0])
	}

	return []byte(out), nil
}

func (t *TimeFieldIndex) PrefixFromArgs(args ...interface{}) ([]byte, error) {
	val, err := t.FromArgs(args...)
	if err != nil {
		return nil, err
	}

	// Strip the null terminator, the rest is a prefix
	n := len(val)
	if n > 0 {
		return val[:n-1], nil
	}
	return val, nil
}

type DurationFieldIndex struct {
	Field string
}

func durationIndex(v time.Duration) string {
	return integerIndex(int64(v))
}

func (d *DurationFieldIndex) FromObject(obj interface{}) (bool, []byte, error) {
	v := reflect.ValueOf(obj)
	v = reflect.Indirect(v) // Dereference the pointer if any

	fv := v.FieldByName(d.Field)
	if !fv.IsValid() {
		return false, nil,
			fmt.Errorf("field '%s' for %#v is invalid", d.Field, obj)
	}

	arg, ok := fv.Interface().(time.Duration)
	if !ok {
		return false, nil,
			fmt.Errorf("field '%s' for %#v is invalid", d.Field, arg)
	}

	// Add the null character as a terminator
	out := durationIndex(arg) + "\x00"

	return true, []byte(out), nil
}

func (d *DurationFieldIndex) FromArgs(args ...interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("must provide only a single argument")
	}

	var out string
	switch arg := args[0].(type) {
	case time.Duration:
		out = durationIndex(arg) + "\x00"
	default:
		return nil, fmt.Errorf("argument must be a duration: %#v", args[0])
	}

	return []byte(out), nil
}

func (d *DurationFieldIndex) PrefixFromArgs(args ...interface{}) ([]byte, error) {
	val, err := d.FromArgs(args...)
	if err != nil {
		return nil, err
	}

	// Strip the null terminator, the rest is a prefix
	n := len(val)
	if n > 0 {
		return val[:n-1], nil
	}
	return val, nil
}

type SampleRateFieldIndex struct {
	Field string
}

func sampleRateIndex(r float64) string {
	if r != 0.0 {
		r = float64(time.Second) / r
	}
	return integerIndex(int64(r))
}

func (s *SampleRateFieldIndex) FromObject(obj interface{}) (bool, []byte, error) {
	v := reflect.ValueOf(obj)
	v = reflect.Indirect(v) // Dereference the pointer if any

	fv := v.FieldByName(s.Field)
	if !fv.IsValid() {
		return false, nil,
			fmt.Errorf("field '%s' for %#v is invalid", s.Field, obj)
	}

	arg, ok := fv.Interface().(float64)
	if !ok {
		return false, nil,
			fmt.Errorf("field '%s' for %#v is invalid", s.Field, arg)
	}

	// Add the null character as a terminator
	out := sampleRateIndex(arg) + "\x00"

	return true, []byte(out), nil
}

func (s *SampleRateFieldIndex) FromArgs(args ...interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("must provide only a single argument")
	}

	var f float64

	switch args[0].(type) {
	case float32:
		f = float64(args[0].(float32))
	case float64:
		f = args[0].(float64)
	default:
		return nil, fmt.Errorf("argument must be a time: %#v", args[0])
	}

	out := sampleRateIndex(f) + "\x00"

	return []byte(out), nil
}

func (s *SampleRateFieldIndex) PrefixFromArgs(args ...interface{}) ([]byte, error) {
	val, err := s.FromArgs(args...)
	if err != nil {
		return nil, err
	}

	// Strip the null terminator, the rest is a prefix
	n := len(val)
	if n > 0 {
		return val[:n-1], nil
	}
	return val, nil
}

type EmptyStringFieldIndex struct {
	Field     string
	Lowercase bool
}

func (s *EmptyStringFieldIndex) FromObject(obj interface{}) (bool, []byte, error) {
	v := reflect.ValueOf(obj)
	v = reflect.Indirect(v) // Dereference the pointer if any

	fv := v.FieldByName(s.Field)
	if !fv.IsValid() {
		return false, nil,
			fmt.Errorf("field '%s' for %#v is invalid", s.Field, obj)
	}

	val := fv.String()
	if val == "" {
		val = "-"
	}

	if s.Lowercase {
		val = strings.ToLower(val)
	}

	// Add the null character as a terminator
	val += "\x00"
	return true, []byte(val), nil
}

func (s *EmptyStringFieldIndex) FromArgs(args ...interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("must provide only a single argument")
	}
	arg, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("argument must be a string: %#v", args[0])
	}
	if arg == "" {
		arg = "-"
	}
	if s.Lowercase {
		arg = strings.ToLower(arg)
	}
	// Add the null character as a terminator
	arg += "\x00"
	return []byte(arg), nil
}

func (s *EmptyStringFieldIndex) PrefixFromArgs(args ...interface{}) ([]byte, error) {
	val, err := s.FromArgs(args...)
	if err != nil {
		return nil, err
	}

	// Strip the null terminator, the rest is a prefix
	n := len(val)
	if n > 0 {
		return val[:n-1], nil
	}
	return val, nil
}

func NewSchema() *memdb.DBSchema {
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"asset": &memdb.TableSchema{
				Name: "asset",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:   "id",
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
							},
						},
					},
					"model": &memdb.IndexSchema{
						Name: "model",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
							},
						},
					},
					"number": &memdb.IndexSchema{
						Name:         "number",
						AllowMissing: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Number"},
							},
						},
					},
				},
			},
			"constituent": &memdb.TableSchema{
				Name: "constituent",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:   "id",
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Gauge"},
								&IntegerFieldIndex{Field: "Number"},
							},
						},
					},
					"gauge": &memdb.IndexSchema{
						Name:    "gauge",
						Indexer: &memdb.StringFieldIndex{Field: "Gauge"},
					},
				},
			},
			"gauge": &memdb.TableSchema{
				Name: "gauge",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Code"},
					},
					"network": &memdb.IndexSchema{
						Name:    "network",
						Indexer: &memdb.StringFieldIndex{Field: "Network"},
					},
				},
			},
			"mark": &memdb.TableSchema{
				Name: "mark",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Code"},
					},
					"network": &memdb.IndexSchema{
						Name:    "network",
						Indexer: &memdb.StringFieldIndex{Field: "Network"},
					},
				},
			},
			"monument": &memdb.TableSchema{
				Name: "monument",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Mark"},
					},
				},
			},
			"mount": &memdb.TableSchema{
				Name: "mount",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Code"},
					},
					"network": &memdb.IndexSchema{
						Name:         "network",
						AllowMissing: true,
						Indexer:      &memdb.StringFieldIndex{Field: "Network"},
					},
				},
			},
			"network": &memdb.TableSchema{
				Name: "network",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Code"},
					},
				},
			},
			"site": &memdb.TableSchema{
				Name: "site",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:   "id",
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Station"},
								&memdb.StringFieldIndex{Field: "Location"},
							},
						},
					},
					"station": &memdb.IndexSchema{
						Name: "station",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Station"},
							},
						},
					},
				},
			},
			"station": &memdb.TableSchema{
				Name: "station",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Code"},
					},
					"network": &memdb.IndexSchema{
						Name:    "network",
						Indexer: &memdb.StringFieldIndex{Field: "Network"},
					},
				},
			},

			"antenna": &memdb.TableSchema{
				Name: "antenna",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:   "id",
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
								&memdb.StringFieldIndex{Field: "Mark"},
								&TimeFieldIndex{Field: "Start"},
							},
						},
					},
					"asset": &memdb.IndexSchema{
						Name: "asset",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
							},
						},
					},
					"mark": &memdb.IndexSchema{
						Name:    "mark",
						Indexer: &memdb.StringFieldIndex{Field: "Mark"},
					},
				},
			},
			"camera": &memdb.TableSchema{
				Name: "camera",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:   "id",
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
								&memdb.StringFieldIndex{Field: "Mount"},
								&TimeFieldIndex{Field: "Start"},
							},
						},
					},
					"asset": &memdb.IndexSchema{
						Name: "asset",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
							},
						},
					},
					"mount": &memdb.IndexSchema{
						Name:    "mount",
						Indexer: &memdb.StringFieldIndex{Field: "Mount"},
					},
				},
			},
			"datalogger": &memdb.TableSchema{
				Name: "datalogger",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:   "id",
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
								&memdb.StringFieldIndex{Field: "Place"},
								&EmptyStringFieldIndex{Field: "Role"},
								&TimeFieldIndex{Field: "Start"},
							},
						},
					},
					"asset": &memdb.IndexSchema{
						Name: "asset",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
							},
						},
					},
					"place": &memdb.IndexSchema{
						Name:    "place",
						Indexer: &memdb.StringFieldIndex{Field: "Place"},
					},
					"role": &memdb.IndexSchema{
						Name:         "role",
						AllowMissing: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Place"},
								&EmptyStringFieldIndex{Field: "Role"},
							},
						},
					},
				},
			},
			"metsensor": &memdb.TableSchema{
				Name: "metsensor",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:   "id",
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
								&memdb.StringFieldIndex{Field: "Mark"},
								&TimeFieldIndex{Field: "Start"},
							},
						},
					},
					"asset": &memdb.IndexSchema{
						Name: "asset",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
							},
						},
					},
					"mark": &memdb.IndexSchema{
						Name:    "mark",
						Indexer: &memdb.StringFieldIndex{Field: "Mark"},
					},
				},
			},
			"radome": &memdb.TableSchema{
				Name: "radome",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:   "id",
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
								&memdb.StringFieldIndex{Field: "Mark"},
								&TimeFieldIndex{Field: "Start"},
							},
						},
					},
					"asset": &memdb.IndexSchema{
						Name: "asset",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
							},
						},
					},
					"mark": &memdb.IndexSchema{
						Name:    "mark",
						Indexer: &memdb.StringFieldIndex{Field: "Mark"},
					},
				},
			},
			"receiver": &memdb.TableSchema{
				Name: "receiver",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:   "id",
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
								&memdb.StringFieldIndex{Field: "Mark"},
								&TimeFieldIndex{Field: "Start"},
							},
						},
					},
					"asset": &memdb.IndexSchema{
						Name: "asset",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
							},
						},
					},
					"mark": &memdb.IndexSchema{
						Name:    "mark",
						Indexer: &memdb.StringFieldIndex{Field: "Mark"},
					},
				},
			},
			"recorder": &memdb.TableSchema{
				Name: "recorder",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:   "id",
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
								&memdb.StringFieldIndex{Field: "Station"},
								&memdb.StringFieldIndex{Field: "Location"},
								&TimeFieldIndex{Field: "Start"},
							},
						},
					},
					"asset": &memdb.IndexSchema{
						Name: "asset",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
							},
						},
					},
					"station": &memdb.IndexSchema{
						Name:    "station",
						Indexer: &memdb.StringFieldIndex{Field: "Station"},
					},
					"site": &memdb.IndexSchema{
						Name: "site",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Station"},
								&memdb.StringFieldIndex{Field: "Location"},
							},
						},
					},
				},
			},
			"sensor": &memdb.TableSchema{
				Name: "sensor",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:   "id",
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
								&memdb.StringFieldIndex{Field: "Station"},
								&memdb.StringFieldIndex{Field: "Location"},
								&TimeFieldIndex{Field: "Start"},
							},
						},
					},
					"asset": &memdb.IndexSchema{
						Name: "asset",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
							},
						},
					},
					"station": &memdb.IndexSchema{
						Name:    "station",
						Indexer: &memdb.StringFieldIndex{Field: "Station"},
					},
					"site": &memdb.IndexSchema{
						Name: "site",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Station"},
								&memdb.StringFieldIndex{Field: "Location"},
							},
						},
					},
				},
			},
			"connection": &memdb.TableSchema{
				Name: "connection",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:   "id",
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Station"},
								&memdb.StringFieldIndex{Field: "Location"},
								&memdb.StringFieldIndex{Field: "Place"},
								&EmptyStringFieldIndex{Field: "Role"},
								&TimeFieldIndex{Field: "Start"},
							},
						},
					},
					"place": &memdb.IndexSchema{
						Name:    "place",
						Indexer: &memdb.StringFieldIndex{Field: "Place"},
					},
					"role": &memdb.IndexSchema{
						Name: "role",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Place"},
								&EmptyStringFieldIndex{Field: "Role"},
							},
						},
					},
					"station": &memdb.IndexSchema{
						Name:    "station",
						Indexer: &memdb.StringFieldIndex{Field: "Station"},
					},
					"site": &memdb.IndexSchema{
						Name: "site",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Station"},
								&memdb.StringFieldIndex{Field: "Location"},
							},
						},
					},
				},
			},
			"firmware": &memdb.TableSchema{
				Name: "firmware",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:   "id",
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
								&memdb.StringFieldIndex{Field: "Version"},
								&TimeFieldIndex{Field: "Start"},
							},
						},
					},
					"asset": &memdb.IndexSchema{
						Name: "asset",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Make"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "Serial"},
							},
						},
					},
				},
			},

			"session": &memdb.TableSchema{
				Name: "session",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:   "id",
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Mark"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "SatelliteSystem"},
								&DurationFieldIndex{Field: "Interval"},
								&TimeFieldIndex{Field: "Start"},
							},
						},
					},
					"mark": &memdb.IndexSchema{
						Name:    "mark",
						Indexer: &memdb.StringFieldIndex{Field: "Mark"},
					},
					"model": &memdb.IndexSchema{
						Name:    "model",
						Indexer: &memdb.StringFieldIndex{Field: "Model"},
					},
					"interval": &memdb.IndexSchema{
						Name: "interval",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Mark"},
								&memdb.StringFieldIndex{Field: "Model"},
								&memdb.StringFieldIndex{Field: "SatelliteSystem"},
								&DurationFieldIndex{Field: "Interval"},
							},
						},
					},
				},
			},
			//Station,Location,Sampling Rate,Axial,Reversed,Start Date,End Date
			"stream": &memdb.TableSchema{
				Name: "stream",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:   "id",
						Unique: true,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Station"},
								&memdb.StringFieldIndex{Field: "Location"},
								&SampleRateFieldIndex{Field: "SamplingRate"},
								&TimeFieldIndex{Field: "Start"},
							},
						},
					},
					"station": &memdb.IndexSchema{
						Name:    "station",
						Indexer: &memdb.StringFieldIndex{Field: "Station"},
					},
					"site": &memdb.IndexSchema{
						Name: "site",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Station"},
								&memdb.StringFieldIndex{Field: "Location"},
							},
						},
					},
					"interval": &memdb.IndexSchema{
						Name: "interval",
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Station"},
								&memdb.StringFieldIndex{Field: "Location"},
								&SampleRateFieldIndex{Field: "SamplingRate"},
							},
						},
					},
				},
			},
		},
	}
}

func NewMetaDB(base string) (*MetaDB, error) {
	db, err := memdb.NewMemDB(NewSchema())
	if err != nil {
		return nil, err
	}

	var networks = []struct {
		table string
		path  string
		list  meta.List
	}{
		{"gauge", "network/gauges.csv", &meta.GaugeList{}},
		{"constituent", "network/constituents.csv", &meta.ConstituentList{}},
		{"mark", "network/marks.csv", &meta.MarkList{}},
		{"monument", "network/monuments.csv", &meta.MonumentList{}},
		{"mount", "network/mounts.csv", &meta.MountList{}},
		{"network", "network/networks.csv", &meta.NetworkList{}},
		{"site", "network/sites.csv", &meta.SiteList{}},
		{"station", "network/stations.csv", &meta.StationList{}},
	}

	txn := db.Txn(true)
	for _, list := range networks {
		if _, err := os.Stat(filepath.Join(base, list.path)); os.IsNotExist(err) {
			continue
		}

		switch list.list.(type) {
		case *meta.ConstituentList:
			var input meta.ConstituentList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.GaugeList:
			var input meta.GaugeList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.MarkList:
			var input meta.MarkList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.MonumentList:
			var input meta.MonumentList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.MountList:
			var input meta.MountList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.NetworkList:
			var input meta.NetworkList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.SiteList:
			var input meta.SiteList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.StationList:
			var input meta.StationList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		default:
			return nil, fmt.Errorf("invalid type found for: %s", list.path)
		}
	}
	txn.Commit()

	var assets = []struct {
		table string
		path  string
	}{
		{"antenna", "assets/antennas.csv"},
		{"camera", "assets/cameras.csv"},
		{"datalogger", "assets/dataloggers.csv"},
		{"metsensor", "assets/metsensors.csv"},
		{"radome", "assets/radomes.csv"},
		{"receiver", "assets/receivers.csv"},
		{"recorder", "assets/recorders.csv"},
		{"sensor", "assets/sensors.csv"},
	}

	txn = db.Txn(true)
	for _, list := range assets {
		if _, err := os.Stat(filepath.Join(base, list.path)); os.IsNotExist(err) {
			continue
		}

		var input meta.AssetList
		if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
			return nil, err
		}
		for _, i := range input {
			if err := txn.Insert("asset", i); err != nil {
				return nil, err
			}
		}

	}
	txn.Commit()

	var installs = []struct {
		table string
		path  string
		list  meta.List
	}{
		{"antenna", "install/antennas.csv", &meta.InstalledAntennaList{}},
		{"camera", "install/cameras.csv", &meta.InstalledCameraList{}},
		{"datalogger", "install/dataloggers.csv", &meta.DeployedDataloggerList{}},
		{"metsensor", "install/metsensors.csv", &meta.InstalledMetSensorList{}},
		{"radome", "install/radomes.csv", &meta.InstalledRadomeList{}},
		{"receiver", "install/receivers.csv", &meta.DeployedReceiverList{}},
		{"recorder", "install/recorders.csv", &meta.InstalledRecorderList{}},
		{"sensor", "install/sensors.csv", &meta.InstalledSensorList{}},
		{"connection", "install/connections.csv", &meta.ConnectionList{}},
		{"firmware", "install/firmware.csv", &meta.FirmwareHistoryList{}},
		{"session", "install/sessions.csv", &meta.SessionList{}},
		{"stream", "install/streams.csv", &meta.StreamList{}},
	}

	txn = db.Txn(true)
	for _, list := range installs {
		if _, err := os.Stat(filepath.Join(base, list.path)); os.IsNotExist(err) {
			continue
		}

		switch list.list.(type) {
		case *meta.InstalledAntennaList:
			var input meta.InstalledAntennaList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.InstalledCameraList:
			var input meta.InstalledCameraList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.DeployedDataloggerList:
			var input meta.DeployedDataloggerList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.InstalledMetSensorList:
			var input meta.InstalledMetSensorList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.InstalledRadomeList:
			var input meta.InstalledRadomeList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.DeployedReceiverList:
			var input meta.DeployedReceiverList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.InstalledRecorderList:
			var input meta.InstalledRecorderList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.InstalledSensorList:
			var input meta.InstalledSensorList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.ConnectionList:
			var input meta.ConnectionList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.FirmwareHistoryList:
			var input meta.FirmwareHistoryList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.SessionList:
			var input meta.SessionList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		case *meta.StreamList:
			var input meta.StreamList
			if err := meta.LoadList(filepath.Join(base, list.path), &input); err != nil {
				return nil, err
			}
			for _, i := range input {
				if err := txn.Insert(list.table, i); err != nil {
					return nil, err
				}
			}
		default:
			return nil, fmt.Errorf("invalid type found for: %s", list.path)
		}
	}
	txn.Commit()

	return &MetaDB{MemDB: db}, nil
}
