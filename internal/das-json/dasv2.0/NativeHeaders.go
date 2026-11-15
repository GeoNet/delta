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
    "fmt"
    "encoding/json"
)

// NativeHeaders Native headers from the interrogator for this acquisition as a key-value list.
type NativeHeaders struct {
  AdditionalProperties map[string]interface{} `json:"-,omitempty"`
}

func (strct *NativeHeaders) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
    comma := false
    // Marshal any additional Properties
    for k, v := range strct.AdditionalProperties {
		if comma {
			buf.WriteString(",")
		}
        buf.WriteString(fmt.Sprintf("\"%s\":", k))
		if tmp, err := json.Marshal(v); err != nil {
			return nil, err
		} else {
			buf.Write(tmp)
		}
        comma = true
	}

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *NativeHeaders) UnmarshalJSON(b []byte) error {
    var jsonMap map[string]json.RawMessage
    if err := json.Unmarshal(b, &jsonMap); err != nil {
        return err
    }
    // parse all the defined properties
    for k, v := range jsonMap {
        switch k {
        default:
            // an additional "interface{}" value
            var additionalValue interface{}
            if err := json.Unmarshal([]byte(v), &additionalValue); err != nil {
                return err // invalid additionalProperty
            }
            if strct.AdditionalProperties == nil {
                strct.AdditionalProperties = make(map[string]interface{}, 0)
            }
            strct.AdditionalProperties[k]= additionalValue
        }
    }
    return nil
}
