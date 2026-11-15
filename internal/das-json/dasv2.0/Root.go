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

// Root The Overview metadata gives a high-level overview of the DAS deployment and describes the geographic location of the installation, dates of the project, and information on the data provider.
type Root struct {

  // Metadata schema version. This identifies which version of the metadata schema that this document conforms to.
  SchemaVersion interface{} `json:"schema_version"`

  // Cables. Multiple cables can be documented.
  Cables []*Cable `json:"cables,omitempty"`

  // Additional comments.
  Comment string `json:"comment,omitempty"`

  // Country where the installation is located. Use ISO 3166-1 alpha-3 three-letter country code.
  Country string `json:"country,omitempty"`

  // Digital object identifier (DOI) that uniquely identifies the metadata; this identifier may only become available following archiving.
  DigitalObjectIdentifier string `json:"digital_object_identifier,omitempty"`

  // End date of data collection at the installation in UTC. If the installation end date is unknown this field should be omitted and updated when it is known.
  EndDate string `json:"end_date,omitempty"`

  // Name(s) of agency that funded the experiment.
  FundingAgency string `json:"funding_agency,omitempty"`

  // Interrogators. Multiple interrogators can be documented.
  Interrogators []*Interrogator `json:"interrogators,omitempty"`

  // Name of the geographic location of the installation.
  Location string `json:"location"`

  // Unique FDSN network code for the installation, following FDSN Source Identifier specification
  NetworkCode string `json:"network_code"`

  // Point of contact for the metadata.
  PointOfContact string `json:"point_of_contact"`

  // Physical address and institution of point of contact.
  PointOfContactAddress string `json:"point_of_contact_address"`

  // Email address of point of contact.
  PointOfContactEmail string `json:"point_of_contact_email"`

  // Principal investigator(s) for the installation.
  PrincipalInvestigator []*PrincipalInvestigatorFields `json:"principal_investigator"`

  // Funding project number. Should be supplied if a number has been assigned by funding agency(s).
  ProjectNumber string `json:"project_number,omitempty"`

  // Brief explanation of the purpose of experiment.
  PurposeOfDataCollection string `json:"purpose_of_data_collection,omitempty"`

  // Start date of data collection at the installation in Coordinated Universal Time (UTC).
  StartDate string `json:"start_date"`
}

func (strct *Root) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
    comma := false
    // "SchemaVersion" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "schema_version" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"schema_version\": ")
	if tmp, err := json.Marshal(strct.SchemaVersion); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "cables" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"cables\": ")
	if tmp, err := json.Marshal(strct.Cables); err != nil {
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
    // Marshal the "country" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"country\": ")
	if tmp, err := json.Marshal(strct.Country); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "digital_object_identifier" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"digital_object_identifier\": ")
	if tmp, err := json.Marshal(strct.DigitalObjectIdentifier); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "end_date" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"end_date\": ")
	if tmp, err := json.Marshal(strct.EndDate); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "funding_agency" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"funding_agency\": ")
	if tmp, err := json.Marshal(strct.FundingAgency); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "interrogators" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"interrogators\": ")
	if tmp, err := json.Marshal(strct.Interrogators); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "Location" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "location" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"location\": ")
	if tmp, err := json.Marshal(strct.Location); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "NetworkCode" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "network_code" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"network_code\": ")
	if tmp, err := json.Marshal(strct.NetworkCode); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "PointOfContact" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "point_of_contact" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"point_of_contact\": ")
	if tmp, err := json.Marshal(strct.PointOfContact); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "PointOfContactAddress" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "point_of_contact_address" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"point_of_contact_address\": ")
	if tmp, err := json.Marshal(strct.PointOfContactAddress); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "PointOfContactEmail" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "point_of_contact_email" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"point_of_contact_email\": ")
	if tmp, err := json.Marshal(strct.PointOfContactEmail); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "PrincipalInvestigator" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "principal_investigator" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"principal_investigator\": ")
	if tmp, err := json.Marshal(strct.PrincipalInvestigator); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "project_number" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"project_number\": ")
	if tmp, err := json.Marshal(strct.ProjectNumber); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "purpose_of_data_collection" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"purpose_of_data_collection\": ")
	if tmp, err := json.Marshal(strct.PurposeOfDataCollection); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "StartDate" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "start_date" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"start_date\": ")
	if tmp, err := json.Marshal(strct.StartDate); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *Root) UnmarshalJSON(b []byte) error {
    schema_versionReceived := false
    locationReceived := false
    network_codeReceived := false
    point_of_contactReceived := false
    point_of_contact_addressReceived := false
    point_of_contact_emailReceived := false
    principal_investigatorReceived := false
    start_dateReceived := false
    var jsonMap map[string]json.RawMessage
    if err := json.Unmarshal(b, &jsonMap); err != nil {
        return err
    }
    // parse all the defined properties
    for k, v := range jsonMap {
        switch k {
        case "schema_version":
            if err := json.Unmarshal([]byte(v), &strct.SchemaVersion); err != nil {
                return err
             }
            schema_versionReceived = true
        case "cables":
            if err := json.Unmarshal([]byte(v), &strct.Cables); err != nil {
                return err
             }
        case "comment":
            if err := json.Unmarshal([]byte(v), &strct.Comment); err != nil {
                return err
             }
        case "country":
            if err := json.Unmarshal([]byte(v), &strct.Country); err != nil {
                return err
             }
        case "digital_object_identifier":
            if err := json.Unmarshal([]byte(v), &strct.DigitalObjectIdentifier); err != nil {
                return err
             }
        case "end_date":
            if err := json.Unmarshal([]byte(v), &strct.EndDate); err != nil {
                return err
             }
        case "funding_agency":
            if err := json.Unmarshal([]byte(v), &strct.FundingAgency); err != nil {
                return err
             }
        case "interrogators":
            if err := json.Unmarshal([]byte(v), &strct.Interrogators); err != nil {
                return err
             }
        case "location":
            if err := json.Unmarshal([]byte(v), &strct.Location); err != nil {
                return err
             }
            locationReceived = true
        case "network_code":
            if err := json.Unmarshal([]byte(v), &strct.NetworkCode); err != nil {
                return err
             }
            network_codeReceived = true
        case "point_of_contact":
            if err := json.Unmarshal([]byte(v), &strct.PointOfContact); err != nil {
                return err
             }
            point_of_contactReceived = true
        case "point_of_contact_address":
            if err := json.Unmarshal([]byte(v), &strct.PointOfContactAddress); err != nil {
                return err
             }
            point_of_contact_addressReceived = true
        case "point_of_contact_email":
            if err := json.Unmarshal([]byte(v), &strct.PointOfContactEmail); err != nil {
                return err
             }
            point_of_contact_emailReceived = true
        case "principal_investigator":
            if err := json.Unmarshal([]byte(v), &strct.PrincipalInvestigator); err != nil {
                return err
             }
            principal_investigatorReceived = true
        case "project_number":
            if err := json.Unmarshal([]byte(v), &strct.ProjectNumber); err != nil {
                return err
             }
        case "purpose_of_data_collection":
            if err := json.Unmarshal([]byte(v), &strct.PurposeOfDataCollection); err != nil {
                return err
             }
        case "start_date":
            if err := json.Unmarshal([]byte(v), &strct.StartDate); err != nil {
                return err
             }
            start_dateReceived = true
        }
    }
    // check if schema_version (a required property) was received
    if !schema_versionReceived {
        return errors.New("\"schema_version\" is required but was not present")
    }
    // check if location (a required property) was received
    if !locationReceived {
        return errors.New("\"location\" is required but was not present")
    }
    // check if network_code (a required property) was received
    if !network_codeReceived {
        return errors.New("\"network_code\" is required but was not present")
    }
    // check if point_of_contact (a required property) was received
    if !point_of_contactReceived {
        return errors.New("\"point_of_contact\" is required but was not present")
    }
    // check if point_of_contact_address (a required property) was received
    if !point_of_contact_addressReceived {
        return errors.New("\"point_of_contact_address\" is required but was not present")
    }
    // check if point_of_contact_email (a required property) was received
    if !point_of_contact_emailReceived {
        return errors.New("\"point_of_contact_email\" is required but was not present")
    }
    // check if principal_investigator (a required property) was received
    if !principal_investigatorReceived {
        return errors.New("\"principal_investigator\" is required but was not present")
    }
    // check if start_date (a required property) was received
    if !start_dateReceived {
        return errors.New("\"start_date\" is required but was not present")
    }
    return nil
}
