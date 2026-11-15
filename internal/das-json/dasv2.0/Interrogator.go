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

// Interrogator The Interrogator metadata contains information about the interrogator unit used for the data collection. More than one interrogator may be operating during an experiment, therefore each unit is given a unique identifier and described in individual Interrogator metadata blocks.
type Interrogator struct {

  // One or more acquisitions with this interrogator.
  Acquisitions []*Acquisition `json:"acquisitions,omitempty"`

  // Additional comments.
  Comment string `json:"comment,omitempty"`

  // Firmware version of the software used within the interrogator.
  FirmwareVersion string `json:"firmware_version,omitempty"`

  // Unique identifier of the interrogator unit used in the experiment, assigned by data provider. Identifier should have the maximum of eight alphanumeric characters with no special characters (e.g., underscores, period, and dash).
  InterrogatorId interface{} `json:"interrogator_id"`

  // Manufacturer name of the interrogator.
  Manufacturer string `json:"manufacturer"`

  // Model number of the interrogator.
  Model string `json:"model"`

  // Serial number of the interrogator.
  SerialNumber string `json:"serial_number,omitempty"`
}

func (strct *Interrogator) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
    comma := false
    // Marshal the "acquisitions" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"acquisitions\": ")
	if tmp, err := json.Marshal(strct.Acquisitions); err != nil {
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
    // Marshal the "firmware_version" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"firmware_version\": ")
	if tmp, err := json.Marshal(strct.FirmwareVersion); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "InterrogatorId" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "interrogator_id" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"interrogator_id\": ")
	if tmp, err := json.Marshal(strct.InterrogatorId); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "Manufacturer" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "manufacturer" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"manufacturer\": ")
	if tmp, err := json.Marshal(strct.Manufacturer); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "Model" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "model" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"model\": ")
	if tmp, err := json.Marshal(strct.Model); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "serial_number" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"serial_number\": ")
	if tmp, err := json.Marshal(strct.SerialNumber); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *Interrogator) UnmarshalJSON(b []byte) error {
    interrogator_idReceived := false
    manufacturerReceived := false
    modelReceived := false
    var jsonMap map[string]json.RawMessage
    if err := json.Unmarshal(b, &jsonMap); err != nil {
        return err
    }
    // parse all the defined properties
    for k, v := range jsonMap {
        switch k {
        case "acquisitions":
            if err := json.Unmarshal([]byte(v), &strct.Acquisitions); err != nil {
                return err
             }
        case "comment":
            if err := json.Unmarshal([]byte(v), &strct.Comment); err != nil {
                return err
             }
        case "firmware_version":
            if err := json.Unmarshal([]byte(v), &strct.FirmwareVersion); err != nil {
                return err
             }
        case "interrogator_id":
            if err := json.Unmarshal([]byte(v), &strct.InterrogatorId); err != nil {
                return err
             }
            interrogator_idReceived = true
        case "manufacturer":
            if err := json.Unmarshal([]byte(v), &strct.Manufacturer); err != nil {
                return err
             }
            manufacturerReceived = true
        case "model":
            if err := json.Unmarshal([]byte(v), &strct.Model); err != nil {
                return err
             }
            modelReceived = true
        case "serial_number":
            if err := json.Unmarshal([]byte(v), &strct.SerialNumber); err != nil {
                return err
             }
        }
    }
    // check if interrogator_id (a required property) was received
    if !interrogator_idReceived {
        return errors.New("\"interrogator_id\" is required but was not present")
    }
    // check if manufacturer (a required property) was received
    if !manufacturerReceived {
        return errors.New("\"manufacturer\" is required but was not present")
    }
    // check if model (a required property) was received
    if !modelReceived {
        return errors.New("\"model\" is required but was not present")
    }
    return nil
}
