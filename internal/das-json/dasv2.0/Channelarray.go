/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

package dasjson

import (
    "bytes"
    "encoding/json"
    "errors"
)

// Channelarray Nested under the Channel Group metadata, the metadata of the Channels is stored in arrays, including a unique identifier, coordinates, and for borehole cable, optional strike and dip. Note that the total amount of individual Channel metadata may be equal or less than the number_of_channels attribute specified under Acquisition metadata as only a subset of the channels may be included in the data. This is to accommodate cases where channels cannot be located or omitted due to various reasons.
type Channelarray struct {

  // Array of Unique IDs of each individual channel, with the maximum of eight alphanumeric characters with no special characters (e.g., underscores, period, and dash).
  ChannelIds []string `json:"channel_ids"`

  // Array of vertical distances from the local ground surface to the fibre at each channel, with the ground surface defined as zero and positive values measured downward.
  DepthsBelowSurface []float64 `json:"depths_below_surface,omitempty"`

  // Array of individual channel dips, measured from horizontal and positive downward: 0° is horizontal, −90° is vertically upward, and +90° is vertically downward.
  Dips []float64 `json:"dips,omitempty"`

  // Array of distances along fiber between the interrogator and channel. Zero meter at the connector to the interrogator; and negative distances for fiber section inside the interrogator.
  DistancesAlongFiber []float64 `json:"distances_along_fiber"`

  // Array of local ground-surface elevations or heights above sea level at each channel.
  ElevationsAboveSeaLevel []float64 `json:"elevations_above_sea_level,omitempty"`

  // Array of local azimuths at each channel, measured clockwise from geographical north in the direction of increasing channel number.
  Strikes []float64 `json:"strikes,omitempty"`

  // Array of X-coordinates (UTM-x or longitude) of the channels.
  XCoordinates []float64 `json:"x_coordinates"`

  // Array of y-coordinates (UTM-y or latitude) of the channels.
  YCoordinates []float64 `json:"y_coordinates"`
}

func (strct *Channelarray) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
    comma := false
    // "ChannelIds" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "channel_ids" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"channel_ids\": ")
	if tmp, err := json.Marshal(strct.ChannelIds); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "depths_below_surface" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"depths_below_surface\": ")
	if tmp, err := json.Marshal(strct.DepthsBelowSurface); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "dips" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"dips\": ")
	if tmp, err := json.Marshal(strct.Dips); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "DistancesAlongFiber" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "distances_along_fiber" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"distances_along_fiber\": ")
	if tmp, err := json.Marshal(strct.DistancesAlongFiber); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "elevations_above_sea_level" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"elevations_above_sea_level\": ")
	if tmp, err := json.Marshal(strct.ElevationsAboveSeaLevel); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "strikes" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"strikes\": ")
	if tmp, err := json.Marshal(strct.Strikes); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "XCoordinates" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "x_coordinates" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"x_coordinates\": ")
	if tmp, err := json.Marshal(strct.XCoordinates); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "YCoordinates" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "y_coordinates" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"y_coordinates\": ")
	if tmp, err := json.Marshal(strct.YCoordinates); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *Channelarray) UnmarshalJSON(b []byte) error {
    channel_idsReceived := false
    distances_along_fiberReceived := false
    x_coordinatesReceived := false
    y_coordinatesReceived := false
    var jsonMap map[string]json.RawMessage
    if err := json.Unmarshal(b, &jsonMap); err != nil {
        return err
    }
    // parse all the defined properties
    for k, v := range jsonMap {
        switch k {
        case "channel_ids":
            if err := json.Unmarshal([]byte(v), &strct.ChannelIds); err != nil {
                return err
             }
            channel_idsReceived = true
        case "depths_below_surface":
            if err := json.Unmarshal([]byte(v), &strct.DepthsBelowSurface); err != nil {
                return err
             }
        case "dips":
            if err := json.Unmarshal([]byte(v), &strct.Dips); err != nil {
                return err
             }
        case "distances_along_fiber":
            if err := json.Unmarshal([]byte(v), &strct.DistancesAlongFiber); err != nil {
                return err
             }
            distances_along_fiberReceived = true
        case "elevations_above_sea_level":
            if err := json.Unmarshal([]byte(v), &strct.ElevationsAboveSeaLevel); err != nil {
                return err
             }
        case "strikes":
            if err := json.Unmarshal([]byte(v), &strct.Strikes); err != nil {
                return err
             }
        case "x_coordinates":
            if err := json.Unmarshal([]byte(v), &strct.XCoordinates); err != nil {
                return err
             }
            x_coordinatesReceived = true
        case "y_coordinates":
            if err := json.Unmarshal([]byte(v), &strct.YCoordinates); err != nil {
                return err
             }
            y_coordinatesReceived = true
        }
    }
    // check if channel_ids (a required property) was received
    if !channel_idsReceived {
        return errors.New("\"channel_ids\" is required but was not present")
    }
    // check if distances_along_fiber (a required property) was received
    if !distances_along_fiberReceived {
        return errors.New("\"distances_along_fiber\" is required but was not present")
    }
    // check if x_coordinates (a required property) was received
    if !x_coordinatesReceived {
        return errors.New("\"x_coordinates\" is required but was not present")
    }
    // check if y_coordinates (a required property) was received
    if !y_coordinatesReceived {
        return errors.New("\"y_coordinates\" is required but was not present")
    }
    return nil
}
