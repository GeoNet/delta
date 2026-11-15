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

// Fiber The intention for a separate Fiber metadata is to uniquely identify the fiber used to make the measurements. More than one cable may be used over the course of an experiment and there may be several fibers within a cable.
type Fiber struct {

  // Additional comments.
  Comment string `json:"comment,omitempty"`

  // The end location or distance along the fiber with respect to the interrogator.
  FiberEndLocation float64 `json:"fiber_end_location,omitempty"`

  // Unit of fiber end location. SI units should be used.
  FiberEndLocationUnit string `json:"fiber_end_location_unit,omitempty"`

  // Geometry of fiber [linear | helical | other].
  FiberGeometry string `json:"fiber_geometry"`

  // Unique identifier for the physical fiber, with the maximum of eight alphanumeric characters with no special characters (e.g., underscores, period, and dash).
  FiberId interface{} `json:"fiber_id"`

  // Mode of fiber [single-mode | multi-mode | other].
  FiberMode string `json:"fiber_mode"`

  // Power loss for one-way travel of a beam of light. Typical value for single mode fiber is 0.4.
  FiberOneWayAttenuation float64 `json:"fiber_one_way_attenuation,omitempty"`

  // Unit of fiber one way attenuation [dB/m | dB/km]. SI units should be used.
  FiberOneWayAttenuationUnit string `json:"fiber_one_way_attenuation_unit,omitempty"`

  // Total length of this fiber. Typically, it is the sum of fiber_start_location and fiber_end_location.
  FiberOpticLength float64 `json:"fiber_optic_length,omitempty"`

  // Unit of fiber-optical length. SI units should be used.
  FiberOpticLengthUnit string `json:"fiber_optic_length_unit,omitempty"`

  // Refractive index of the fiber (unitless). Typical value for single-mode fiber is 1.4681.
  FiberRefractionIndex float64 `json:"fiber_refraction_index"`

  // The start distance along the fiber with respect to the interrogator, referred to as StartLocus in PRODML.
  FiberStartLocation float64 `json:"fiber_start_location,omitempty"`

  // Unit of fiber start location. SI units should be used.
  FiberStartLocationUnit string `json:"fiber_start_location_unit,omitempty"`

  // If helical, this is the winding angle relative to the center of the cable.
  FiberWindingAngle float64 `json:"fiber_winding_angle,omitempty"`

  // Unit of fiber winding angle.
  FiberWindingAngleUnit string `json:"fiber_winding_angle_unit,omitempty"`
}

func (strct *Fiber) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
    comma := false
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
    // Marshal the "fiber_end_location" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"fiber_end_location\": ")
	if tmp, err := json.Marshal(strct.FiberEndLocation); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "fiber_end_location_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"fiber_end_location_unit\": ")
	if tmp, err := json.Marshal(strct.FiberEndLocationUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "FiberGeometry" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "fiber_geometry" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"fiber_geometry\": ")
	if tmp, err := json.Marshal(strct.FiberGeometry); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "FiberId" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "fiber_id" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"fiber_id\": ")
	if tmp, err := json.Marshal(strct.FiberId); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "FiberMode" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "fiber_mode" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"fiber_mode\": ")
	if tmp, err := json.Marshal(strct.FiberMode); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "fiber_one_way_attenuation" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"fiber_one_way_attenuation\": ")
	if tmp, err := json.Marshal(strct.FiberOneWayAttenuation); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "fiber_one_way_attenuation_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"fiber_one_way_attenuation_unit\": ")
	if tmp, err := json.Marshal(strct.FiberOneWayAttenuationUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "fiber_optic_length" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"fiber_optic_length\": ")
	if tmp, err := json.Marshal(strct.FiberOpticLength); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "fiber_optic_length_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"fiber_optic_length_unit\": ")
	if tmp, err := json.Marshal(strct.FiberOpticLengthUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "FiberRefractionIndex" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "fiber_refraction_index" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"fiber_refraction_index\": ")
	if tmp, err := json.Marshal(strct.FiberRefractionIndex); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "fiber_start_location" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"fiber_start_location\": ")
	if tmp, err := json.Marshal(strct.FiberStartLocation); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "fiber_start_location_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"fiber_start_location_unit\": ")
	if tmp, err := json.Marshal(strct.FiberStartLocationUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "fiber_winding_angle" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"fiber_winding_angle\": ")
	if tmp, err := json.Marshal(strct.FiberWindingAngle); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "fiber_winding_angle_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"fiber_winding_angle_unit\": ")
	if tmp, err := json.Marshal(strct.FiberWindingAngleUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *Fiber) UnmarshalJSON(b []byte) error {
    fiber_geometryReceived := false
    fiber_idReceived := false
    fiber_modeReceived := false
    fiber_refraction_indexReceived := false
    var jsonMap map[string]json.RawMessage
    if err := json.Unmarshal(b, &jsonMap); err != nil {
        return err
    }
    // parse all the defined properties
    for k, v := range jsonMap {
        switch k {
        case "comment":
            if err := json.Unmarshal([]byte(v), &strct.Comment); err != nil {
                return err
             }
        case "fiber_end_location":
            if err := json.Unmarshal([]byte(v), &strct.FiberEndLocation); err != nil {
                return err
             }
        case "fiber_end_location_unit":
            if err := json.Unmarshal([]byte(v), &strct.FiberEndLocationUnit); err != nil {
                return err
             }
        case "fiber_geometry":
            if err := json.Unmarshal([]byte(v), &strct.FiberGeometry); err != nil {
                return err
             }
            fiber_geometryReceived = true
        case "fiber_id":
            if err := json.Unmarshal([]byte(v), &strct.FiberId); err != nil {
                return err
             }
            fiber_idReceived = true
        case "fiber_mode":
            if err := json.Unmarshal([]byte(v), &strct.FiberMode); err != nil {
                return err
             }
            fiber_modeReceived = true
        case "fiber_one_way_attenuation":
            if err := json.Unmarshal([]byte(v), &strct.FiberOneWayAttenuation); err != nil {
                return err
             }
        case "fiber_one_way_attenuation_unit":
            if err := json.Unmarshal([]byte(v), &strct.FiberOneWayAttenuationUnit); err != nil {
                return err
             }
        case "fiber_optic_length":
            if err := json.Unmarshal([]byte(v), &strct.FiberOpticLength); err != nil {
                return err
             }
        case "fiber_optic_length_unit":
            if err := json.Unmarshal([]byte(v), &strct.FiberOpticLengthUnit); err != nil {
                return err
             }
        case "fiber_refraction_index":
            if err := json.Unmarshal([]byte(v), &strct.FiberRefractionIndex); err != nil {
                return err
             }
            fiber_refraction_indexReceived = true
        case "fiber_start_location":
            if err := json.Unmarshal([]byte(v), &strct.FiberStartLocation); err != nil {
                return err
             }
        case "fiber_start_location_unit":
            if err := json.Unmarshal([]byte(v), &strct.FiberStartLocationUnit); err != nil {
                return err
             }
        case "fiber_winding_angle":
            if err := json.Unmarshal([]byte(v), &strct.FiberWindingAngle); err != nil {
                return err
             }
        case "fiber_winding_angle_unit":
            if err := json.Unmarshal([]byte(v), &strct.FiberWindingAngleUnit); err != nil {
                return err
             }
        }
    }
    // check if fiber_geometry (a required property) was received
    if !fiber_geometryReceived {
        return errors.New("\"fiber_geometry\" is required but was not present")
    }
    // check if fiber_id (a required property) was received
    if !fiber_idReceived {
        return errors.New("\"fiber_id\" is required but was not present")
    }
    // check if fiber_mode (a required property) was received
    if !fiber_modeReceived {
        return errors.New("\"fiber_mode\" is required but was not present")
    }
    // check if fiber_refraction_index (a required property) was received
    if !fiber_refraction_indexReceived {
        return errors.New("\"fiber_refraction_index\" is required but was not present")
    }
    return nil
}
