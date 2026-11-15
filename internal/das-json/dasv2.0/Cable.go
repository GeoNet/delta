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
    "encoding/json"
    "errors"
    "bytes"
)

// Cable The Cable metadata is designed to allow investigators to describe the cable used in their experiment, which could be one continuous long cable, or cables with multiple connectors or across several installation environments. There are also parameters relating to the general geographical location of the cable to facilitate searchability in future databases. The Cable metadata should include available information on the cable and its installation as this knowledge may help in developing a system-level description of the exact cable response.
type Cable struct {

  // List of approximate coordinates of the bounding box containing the array [minimum latitude, maximum latitude, minimum longitude, and maximum longitude].
  CableBoundingBox []float64 `json:"cable_bounding_box"`

  // Describes how the fiber is encased in the cable.
  CableCharacteristics string `json:"cable_characteristics,omitempty"`

  // Describes the installation environment of cable [conduit| trench | outside borehole casing | wireline | other ].
  CableEnvironment string `json:"cable_environment,omitempty"`

  // Unique identifier for the specific physical cable used in the data acquisition. Identifier should have the maximum of eight alphanumeric characters with no special characters (e.g., underscores, period, and dash).
  CableId interface{} `json:"cable_id"`

  // Cable installation date in Coordinated Universal Time (UTC). This may not be known, an optional field.
  CableInstallationDate string `json:"cable_installation_date,omitempty"`

  // Plain language description of how the cable was installed.
  CableInstallationEnvironment string `json:"cable_installation_environment,omitempty"`

  // Cable manufacturer model.
  CableModel string `json:"cable_model,omitempty"`

  // Cable outside diameter.
  CableOutsideDiameter float64 `json:"cable_outside_diameter,omitempty"`

  // Unit of outside diameter of cable. SI units should be used.
  CableOutsideDiameterUnit string `json:"cable_outside_diameter_unit,omitempty"`

  // Proprietary owner of the cable.
  CableOwner string `json:"cable_owner"`

  // Cable removal date in UTC. This may not be known, an optional field.
  CableRemovalDate string `json:"cable_removal_date,omitempty"`

  // Additional comments.
  Comment string `json:"comment,omitempty"`

  // One or more Fibers for this cable
  Fibers []*Fiber `json:"fibers,omitempty"`
}

func (strct *Cable) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
    comma := false
    // "CableBoundingBox" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "cable_bounding_box" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"cable_bounding_box\": ")
	if tmp, err := json.Marshal(strct.CableBoundingBox); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "cable_characteristics" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"cable_characteristics\": ")
	if tmp, err := json.Marshal(strct.CableCharacteristics); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "cable_environment" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"cable_environment\": ")
	if tmp, err := json.Marshal(strct.CableEnvironment); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "CableId" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "cable_id" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"cable_id\": ")
	if tmp, err := json.Marshal(strct.CableId); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "cable_installation_date" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"cable_installation_date\": ")
	if tmp, err := json.Marshal(strct.CableInstallationDate); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "cable_installation_environment" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"cable_installation_environment\": ")
	if tmp, err := json.Marshal(strct.CableInstallationEnvironment); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "cable_model" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"cable_model\": ")
	if tmp, err := json.Marshal(strct.CableModel); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "cable_outside_diameter" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"cable_outside_diameter\": ")
	if tmp, err := json.Marshal(strct.CableOutsideDiameter); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "cable_outside_diameter_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"cable_outside_diameter_unit\": ")
	if tmp, err := json.Marshal(strct.CableOutsideDiameterUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "CableOwner" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "cable_owner" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"cable_owner\": ")
	if tmp, err := json.Marshal(strct.CableOwner); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "cable_removal_date" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"cable_removal_date\": ")
	if tmp, err := json.Marshal(strct.CableRemovalDate); err != nil {
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
    // Marshal the "fibers" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"fibers\": ")
	if tmp, err := json.Marshal(strct.Fibers); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *Cable) UnmarshalJSON(b []byte) error {
    cable_bounding_boxReceived := false
    cable_idReceived := false
    cable_ownerReceived := false
    var jsonMap map[string]json.RawMessage
    if err := json.Unmarshal(b, &jsonMap); err != nil {
        return err
    }
    // parse all the defined properties
    for k, v := range jsonMap {
        switch k {
        case "cable_bounding_box":
            if err := json.Unmarshal([]byte(v), &strct.CableBoundingBox); err != nil {
                return err
             }
            cable_bounding_boxReceived = true
        case "cable_characteristics":
            if err := json.Unmarshal([]byte(v), &strct.CableCharacteristics); err != nil {
                return err
             }
        case "cable_environment":
            if err := json.Unmarshal([]byte(v), &strct.CableEnvironment); err != nil {
                return err
             }
        case "cable_id":
            if err := json.Unmarshal([]byte(v), &strct.CableId); err != nil {
                return err
             }
            cable_idReceived = true
        case "cable_installation_date":
            if err := json.Unmarshal([]byte(v), &strct.CableInstallationDate); err != nil {
                return err
             }
        case "cable_installation_environment":
            if err := json.Unmarshal([]byte(v), &strct.CableInstallationEnvironment); err != nil {
                return err
             }
        case "cable_model":
            if err := json.Unmarshal([]byte(v), &strct.CableModel); err != nil {
                return err
             }
        case "cable_outside_diameter":
            if err := json.Unmarshal([]byte(v), &strct.CableOutsideDiameter); err != nil {
                return err
             }
        case "cable_outside_diameter_unit":
            if err := json.Unmarshal([]byte(v), &strct.CableOutsideDiameterUnit); err != nil {
                return err
             }
        case "cable_owner":
            if err := json.Unmarshal([]byte(v), &strct.CableOwner); err != nil {
                return err
             }
            cable_ownerReceived = true
        case "cable_removal_date":
            if err := json.Unmarshal([]byte(v), &strct.CableRemovalDate); err != nil {
                return err
             }
        case "comment":
            if err := json.Unmarshal([]byte(v), &strct.Comment); err != nil {
                return err
             }
        case "fibers":
            if err := json.Unmarshal([]byte(v), &strct.Fibers); err != nil {
                return err
             }
        }
    }
    // check if cable_bounding_box (a required property) was received
    if !cable_bounding_boxReceived {
        return errors.New("\"cable_bounding_box\" is required but was not present")
    }
    // check if cable_id (a required property) was received
    if !cable_idReceived {
        return errors.New("\"cable_id\" is required but was not present")
    }
    // check if cable_owner (a required property) was received
    if !cable_ownerReceived {
        return errors.New("\"cable_owner\" is required but was not present")
    }
    return nil
}
