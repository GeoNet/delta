package main

import (
	"encoding/json"
)

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Feature struct {
	Id         string                 `json:"id"`
	Type       string                 `json:"type"`
	Geometry   Geometry               `json:"geometry"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	BBox       []float64              `json:"bbox,omitempty"`
}

func NewFeature() *Feature {
	return &Feature{
		Type:       "Feature",
		Properties: make(map[string]interface{}),
	}
}

func (f *Feature) SetId(id string) {
	f.Id = id
}

func (f *Feature) AddProperty(key string, value interface{}) {
	if f.Properties == nil {
		f.Properties = make(map[string]interface{})
	}
	f.Properties[key] = value
}

func (f *Feature) AddPointGeometry(lon, lat float64) {
	f.Geometry.Type = "Point"
	f.Geometry.Coordinates = []float64{lon, lat}
}

func (f *Feature) GetGeometry() (float64, float64) {
	if len(f.Geometry.Coordinates) > 1 {
		return f.Geometry.Coordinates[0], f.Geometry.Coordinates[1]
	}
	return 0.0, 0.0
}

func (f *Feature) GetProperty(key string) (interface{}, bool) {
	if f.Properties == nil {
		return nil, false
	}

	v, ok := f.Properties[key]

	return v, ok
}

func (f *Feature) GetStringProperty(key string) (string, bool) {
	if v, ok := f.GetProperty(key); ok {
		switch v := v.(type) {
		case string:
			return v, true
		}
	}
	return "", false
}

func (f *Feature) GetFloat64Property(key string) (float64, bool) {
	if v, ok := f.GetProperty(key); ok {
		switch v := v.(type) {
		case float32:
			return float64(v), true
		case float64:
			return v, true
		}
	}
	return 0.0, false
}

func (f *Feature) GetIntProperty(key string) (int, bool) {
	if v, ok := f.GetProperty(key); ok {
		switch v := v.(type) {
		case int16:
			return int(v), true
		case int32:
			return int(v), true
		case int:
			return v, true
		}
	}
	return 0, false
}

type FeatureCollection struct {
	Type     string                 `json:"type"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Features []Feature              `json:"features"`
	BBox     []float64              `json:"bbox,omitempty"`
}

func NewFeatureCollection() *FeatureCollection {
	return &FeatureCollection{
		Type:     "FeatureCollection",
		Metadata: make(map[string]interface{}),
	}
}

func (f *FeatureCollection) AddFeature(feature Feature) {
	f.Features = append(f.Features, feature)
}

func (f *FeatureCollection) AddMetadata(key string, value interface{}) {
	if f.Metadata == nil {
		f.Metadata = make(map[string]interface{})
	}
	f.Metadata[key] = value
}

func (f *FeatureCollection) GetMetadata(key string) (interface{}, bool) {
	if f.Metadata == nil {
		return nil, false
	}

	v, ok := f.Metadata[key]

	return v, ok
}

func (f *FeatureCollection) GetStringMetadata(key string) (string, bool) {
	if v, ok := f.GetMetadata(key); ok {
		switch v := v.(type) {
		case string:
			return v, true
		}
	}
	return "", false
}

func (f *FeatureCollection) MustGetStringMetadata(key string) string {
	if v, ok := f.GetStringMetadata(key); ok {
		return v
	}
	return ""
}

func (f *FeatureCollection) GetIntMetadata(key string) (int, bool) {
	if v, ok := f.GetMetadata(key); ok {
		switch v := v.(type) {
		case int16:
			return int(v), true
		case int32:
			return int(v), true
		case int:
			return v, true
		}
	}
	return 0, false
}

func (f *FeatureCollection) GetFloat64Metadata(key string) (float64, bool) {
	if v, ok := f.GetMetadata(key); ok {
		switch v := v.(type) {
		case float32:
			return float64(v), true
		case float64:
			return v, true
		}
	}
	return 0.0, false
}

func (f *FeatureCollection) Marshal() ([]byte, error) {
	return json.Marshal(f)
}

func (f *FeatureCollection) MarshalIndent(prefix, offset string) ([]byte, error) {
	return json.MarshalIndent(f, prefix, offset)
}
