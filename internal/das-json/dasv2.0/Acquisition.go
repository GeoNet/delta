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

// Acquisition The Acquisition metadata describes the acquisition parameters set by the user on a particular interrogator. While different interrogator vendors may implement various photonics design and algorithms which are not easily translatable across interrogators, the Acquisition metadata intends to capture the critical acquisition attributes that are commonly required to analyze and interpret the output time series or matrices. These attributes are carefully defined to be standardized across all interrogator systems.
type Acquisition struct {

  // End time of this data acquisition in UTC. If data acquisition is still in operation, this field should be omitted and updated when it is known.
  AcquisitionEndTime string `json:"acquisition_end_time,omitempty"`

  // Unique identifier of the data acquisition, assigned by data provider. Identifier should have the maximum of eight alphanumeric characters with no special characters (e.g., underscores, period, and dash). One identifier per acquisition settings.
  AcquisitionId interface{} `json:"acquisition_id"`

  // The rate at which the interrogator provides output data.
  AcquisitionSampleRate float64 `json:"acquisition_sample_rate"`

  // Unit of acquisition sample rate. SI units should be used.
  AcquisitionSampleRateUnit string `json:"acquisition_sample_rate_unit"`

  // Start time of this data acquisition in Coordinated Universal Time (UTC).
  AcquisitionStartTime string `json:"acquisition_start_time"`

  // One or more channel groups for this acquisition.
  ChannelGroups []*ChannelGroup `json:"channel_groups,omitempty"`

  // Additional comments.
  Comment string `json:"comment,omitempty"`

  // The averaging length along the fiber for a measurement, determined at experiment setup and used during acquisition.
  GaugeLength float64 `json:"gauge_length"`

  // Unit of gauge length. SI units should be used.
  GaugeLengthUnit string `json:"gauge_length_unit"`

  // Native headers from the interrogator for this acquisition as a key-value list.
  NativeHeaders *NativeHeaders `json:"native_headers,omitempty"`

  // The total number of sampling points along the fiber as output from the interrogator, referred to as NumberOfLoci in PRODML.
  NumberOfChannels int `json:"number_of_channels"`

  // Rate at which the interrogator unit interrogates the fiber sensor.
  PulseRate float64 `json:"pulse_rate,omitempty"`

  // Unit of pulse rate. SI units should be used.
  PulseRateUnit string `json:"pulse_rate_unit,omitempty"`

  // Width of the pulse sent down the fiber in unit of time.
  PulseWidth float64 `json:"pulse_width,omitempty"`

  // Unit of pulse width. SI units should be used.
  PulseWidthUnit string `json:"pulse_width_unit,omitempty"`

  // Scaling factor to multiply in order to express data in physical units.
  ScaleFactor float64 `json:"scale_factor,omitempty"`

  // The channel spacing, or offset, between channels.
  SpatialSamplingInterval float64 `json:"spatial_sampling_interval"`

  // Unit of spatial sampling interval. SI units should be used.
  SpatialSamplingIntervalUnit string `json:"spatial_sampling_interval_unit"`

  // Unit of measure of archived data set. This may be the same unit as the Interrogator unit measure if the data are raw. SI units should be used.
  UnitOfMeasure interface{} `json:"unit_of_measure"`
}

func (strct *Acquisition) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
    comma := false
    // Marshal the "acquisition_end_time" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"acquisition_end_time\": ")
	if tmp, err := json.Marshal(strct.AcquisitionEndTime); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "AcquisitionId" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "acquisition_id" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"acquisition_id\": ")
	if tmp, err := json.Marshal(strct.AcquisitionId); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "AcquisitionSampleRate" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "acquisition_sample_rate" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"acquisition_sample_rate\": ")
	if tmp, err := json.Marshal(strct.AcquisitionSampleRate); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "AcquisitionSampleRateUnit" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "acquisition_sample_rate_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"acquisition_sample_rate_unit\": ")
	if tmp, err := json.Marshal(strct.AcquisitionSampleRateUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "AcquisitionStartTime" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "acquisition_start_time" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"acquisition_start_time\": ")
	if tmp, err := json.Marshal(strct.AcquisitionStartTime); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "channel_groups" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"channel_groups\": ")
	if tmp, err := json.Marshal(strct.ChannelGroups); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "comment" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"comment\": ")
	if tmp, err := json.Marshal(strct.Comment); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "GaugeLength" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "gauge_length" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"gauge_length\": ")
	if tmp, err := json.Marshal(strct.GaugeLength); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "GaugeLengthUnit" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "gauge_length_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"gauge_length_unit\": ")
	if tmp, err := json.Marshal(strct.GaugeLengthUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "native_headers" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"native_headers\": ")
	if tmp, err := json.Marshal(strct.NativeHeaders); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "NumberOfChannels" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "number_of_channels" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"number_of_channels\": ")
	if tmp, err := json.Marshal(strct.NumberOfChannels); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "pulse_rate" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"pulse_rate\": ")
	if tmp, err := json.Marshal(strct.PulseRate); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "pulse_rate_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"pulse_rate_unit\": ")
	if tmp, err := json.Marshal(strct.PulseRateUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "pulse_width" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"pulse_width\": ")
	if tmp, err := json.Marshal(strct.PulseWidth); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "pulse_width_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"pulse_width_unit\": ")
	if tmp, err := json.Marshal(strct.PulseWidthUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "scale_factor" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"scale_factor\": ")
	if tmp, err := json.Marshal(strct.ScaleFactor); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "SpatialSamplingInterval" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "spatial_sampling_interval" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"spatial_sampling_interval\": ")
	if tmp, err := json.Marshal(strct.SpatialSamplingInterval); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "SpatialSamplingIntervalUnit" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "spatial_sampling_interval_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"spatial_sampling_interval_unit\": ")
	if tmp, err := json.Marshal(strct.SpatialSamplingIntervalUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "UnitOfMeasure" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "unit_of_measure" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"unit_of_measure\": ")
	if tmp, err := json.Marshal(strct.UnitOfMeasure); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *Acquisition) UnmarshalJSON(b []byte) error {
    acquisition_idReceived := false
    acquisition_sample_rateReceived := false
    acquisition_sample_rate_unitReceived := false
    acquisition_start_timeReceived := false
    gauge_lengthReceived := false
    gauge_length_unitReceived := false
    number_of_channelsReceived := false
    spatial_sampling_intervalReceived := false
    spatial_sampling_interval_unitReceived := false
    unit_of_measureReceived := false
    var jsonMap map[string]json.RawMessage
    if err := json.Unmarshal(b, &jsonMap); err != nil {
        return err
    }
    // parse all the defined properties
    for k, v := range jsonMap {
        switch k {
        case "acquisition_end_time":
            if err := json.Unmarshal([]byte(v), &strct.AcquisitionEndTime); err != nil {
                return err
             }
        case "acquisition_id":
            if err := json.Unmarshal([]byte(v), &strct.AcquisitionId); err != nil {
                return err
             }
            acquisition_idReceived = true
        case "acquisition_sample_rate":
            if err := json.Unmarshal([]byte(v), &strct.AcquisitionSampleRate); err != nil {
                return err
             }
            acquisition_sample_rateReceived = true
        case "acquisition_sample_rate_unit":
            if err := json.Unmarshal([]byte(v), &strct.AcquisitionSampleRateUnit); err != nil {
                return err
             }
            acquisition_sample_rate_unitReceived = true
        case "acquisition_start_time":
            if err := json.Unmarshal([]byte(v), &strct.AcquisitionStartTime); err != nil {
                return err
             }
            acquisition_start_timeReceived = true
        case "channel_groups":
            if err := json.Unmarshal([]byte(v), &strct.ChannelGroups); err != nil {
                return err
             }
        case "comment":
            if err := json.Unmarshal([]byte(v), &strct.Comment); err != nil {
                return err
             }
        case "gauge_length":
            if err := json.Unmarshal([]byte(v), &strct.GaugeLength); err != nil {
                return err
             }
            gauge_lengthReceived = true
        case "gauge_length_unit":
            if err := json.Unmarshal([]byte(v), &strct.GaugeLengthUnit); err != nil {
                return err
             }
            gauge_length_unitReceived = true
        case "native_headers":
            if err := json.Unmarshal([]byte(v), &strct.NativeHeaders); err != nil {
                return err
             }
        case "number_of_channels":
            if err := json.Unmarshal([]byte(v), &strct.NumberOfChannels); err != nil {
                return err
             }
            number_of_channelsReceived = true
        case "pulse_rate":
            if err := json.Unmarshal([]byte(v), &strct.PulseRate); err != nil {
                return err
             }
        case "pulse_rate_unit":
            if err := json.Unmarshal([]byte(v), &strct.PulseRateUnit); err != nil {
                return err
             }
        case "pulse_width":
            if err := json.Unmarshal([]byte(v), &strct.PulseWidth); err != nil {
                return err
             }
        case "pulse_width_unit":
            if err := json.Unmarshal([]byte(v), &strct.PulseWidthUnit); err != nil {
                return err
             }
        case "scale_factor":
            if err := json.Unmarshal([]byte(v), &strct.ScaleFactor); err != nil {
                return err
             }
        case "spatial_sampling_interval":
            if err := json.Unmarshal([]byte(v), &strct.SpatialSamplingInterval); err != nil {
                return err
             }
            spatial_sampling_intervalReceived = true
        case "spatial_sampling_interval_unit":
            if err := json.Unmarshal([]byte(v), &strct.SpatialSamplingIntervalUnit); err != nil {
                return err
             }
            spatial_sampling_interval_unitReceived = true
        case "unit_of_measure":
            if err := json.Unmarshal([]byte(v), &strct.UnitOfMeasure); err != nil {
                return err
             }
            unit_of_measureReceived = true
        }
    }
    // check if acquisition_id (a required property) was received
    if !acquisition_idReceived {
        return errors.New("\"acquisition_id\" is required but was not present")
    }
    // check if acquisition_sample_rate (a required property) was received
    if !acquisition_sample_rateReceived {
        return errors.New("\"acquisition_sample_rate\" is required but was not present")
    }
    // check if acquisition_sample_rate_unit (a required property) was received
    if !acquisition_sample_rate_unitReceived {
        return errors.New("\"acquisition_sample_rate_unit\" is required but was not present")
    }
    // check if acquisition_start_time (a required property) was received
    if !acquisition_start_timeReceived {
        return errors.New("\"acquisition_start_time\" is required but was not present")
    }
    // check if gauge_length (a required property) was received
    if !gauge_lengthReceived {
        return errors.New("\"gauge_length\" is required but was not present")
    }
    // check if gauge_length_unit (a required property) was received
    if !gauge_length_unitReceived {
        return errors.New("\"gauge_length_unit\" is required but was not present")
    }
    // check if number_of_channels (a required property) was received
    if !number_of_channelsReceived {
        return errors.New("\"number_of_channels\" is required but was not present")
    }
    // check if spatial_sampling_interval (a required property) was received
    if !spatial_sampling_intervalReceived {
        return errors.New("\"spatial_sampling_interval\" is required but was not present")
    }
    // check if spatial_sampling_interval_unit (a required property) was received
    if !spatial_sampling_interval_unitReceived {
        return errors.New("\"spatial_sampling_interval_unit\" is required but was not present")
    }
    // check if unit_of_measure (a required property) was received
    if !unit_of_measureReceived {
        return errors.New("\"unit_of_measure\" is required but was not present")
    }
    return nil
}
