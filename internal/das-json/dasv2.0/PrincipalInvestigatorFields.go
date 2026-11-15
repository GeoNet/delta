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
    "errors"
    "bytes"
    "encoding/json"
)

// PrincipalInvestigatorFields The Principal Investigator metadata contains information about the principal investigator for the installation.
type PrincipalInvestigatorFields struct {

  // Physical address and institution of principal investigator.
  Address string `json:"address"`

  // Email address of principal investigator.
  Email string `json:"email"`

  // Name of principal investigator for the installation.
  Name string `json:"name"`
}

func (strct *PrincipalInvestigatorFields) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
    comma := false
    // "Address" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "address" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"address\": ")
	if tmp, err := json.Marshal(strct.Address); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "Email" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "email" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"email\": ")
	if tmp, err := json.Marshal(strct.Email); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "Name" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "name" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"name\": ")
	if tmp, err := json.Marshal(strct.Name); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *PrincipalInvestigatorFields) UnmarshalJSON(b []byte) error {
    addressReceived := false
    emailReceived := false
    nameReceived := false
    var jsonMap map[string]json.RawMessage
    if err := json.Unmarshal(b, &jsonMap); err != nil {
        return err
    }
    // parse all the defined properties
    for k, v := range jsonMap {
        switch k {
        case "address":
            if err := json.Unmarshal([]byte(v), &strct.Address); err != nil {
                return err
             }
            addressReceived = true
        case "email":
            if err := json.Unmarshal([]byte(v), &strct.Email); err != nil {
                return err
             }
            emailReceived = true
        case "name":
            if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
                return err
             }
            nameReceived = true
        }
    }
    // check if address (a required property) was received
    if !addressReceived {
        return errors.New("\"address\" is required but was not present")
    }
    // check if email (a required property) was received
    if !emailReceived {
        return errors.New("\"email\" is required but was not present")
    }
    // check if name (a required property) was received
    if !nameReceived {
        return errors.New("\"name\" is required but was not present")
    }
    return nil
}
