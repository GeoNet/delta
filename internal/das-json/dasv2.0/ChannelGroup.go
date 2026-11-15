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

// ChannelGroup The Channel Group metadata describes the properties that are common to all the channels nested under the Interrogator and Acquisition metadata blocks. Having a group-level metadata removes the repetition of common attributes in the individual Channel metadata. This design allows multiple channel groups to exist under the same interrogator and data acquisition settings. The unique identifiers for the interrogator, acquisition, cable and fiber are recorded as well in this metadata block to allows mapping between the interrogator, acquisition, and fiber that were used for this set of coordinates during a particular data collection.
type ChannelGroup struct {

  // Unique identifier of cable linked to this channel group.
  CableId interface{} `json:"cable_id"`

  // Unique identifier of the channel group. Identifier should have the maximum of eight alphanumeric characters with no special characters (e.g., underscores, period, and dash).
  ChannelGroupId interface{} `json:"channel_group_id"`

  // Channels for this channel group.
  Channels *Channelarray `json:"channels,omitempty"`

  // Additional comments.
  Comment string `json:"comment,omitempty"`

  // Date that the channel coordinates were generated.
  CoordinateGenerationDate string `json:"coordinate_generation_date"`

  // Coordinate system used in channel location. Select geographic for coordinates in latitude and longitude. Universal Transverse Mercator (UTM) is the preferred projected coordinate system.
  CoordinateSystem interface{} `json:"coordinate_system"`

  // Unit of depth below surface. Applies to every individual channel nested under the channel group. SI units should be used.
  DepthBelowSurfaceUnit string `json:"depth_below_surface_unit,omitempty"`

  // Unit of dip. Applies to every individual channel nested under the channel group.
  DipUnit string `json:"dip_unit,omitempty"`

  // Unit of distance along fiber. SI units should be used.
  DistanceAlongFiberUnit string `json:"distance_along_fiber_unit"`

  // Unit of elevation above sea level. Applies to every individual channel nested under the channel group. SI units should be used.
  ElevationAboveSeaLevelUnit string `json:"elevation_above_sea_level_unit,omitempty"`

  // Unique identifier of fiber linked to this channel group. The fiber_id must nest within specified cable_id.
  FiberId interface{} `json:"fiber_id"`

  // The first channel considered physically meaningful. Typically, this will be a few channels after the last channel inside the interrogator. Must match one of the channel IDs.
  FirstUsableChannelId string `json:"first_usable_channel_id,omitempty"`

  // The last channel considered physically meaningful. Typically, an acquisition will add virtual channels at the end of the fiber or the optical path might not reach the end of the fiber, and this can be used to cut the data set into a working set. Must match one of the channel IDs.
  LastUsableChannelId string `json:"last_usable_channel_id,omitempty"`

  // Method used to determine the coordinates of each channel. Example: tap test, GPS tracker.
  LocationMethod string `json:"location_method,omitempty"`

  // Geodetic datum and/or grid zone number in UTM.
  ReferenceFrame string `json:"reference_frame"`

  // Unit of strike. Applies to every individual channel nested under the channel group.
  StrikeUnit string `json:"strike_unit,omitempty"`

  // Uncertainty in depth below surface, representing the average uncertainty in depth for all the coordinates.
  UncertaintyInDepth float64 `json:"uncertainty_in_depth,omitempty"`

  // Unit of uncertainty in depth. SI units should be used.
  UncertaintyInDepthUnit string `json:"uncertainty_in_depth_unit,omitempty"`

  // Uncertainty in the dip.
  UncertaintyInDip float64 `json:"uncertainty_in_dip,omitempty"`

  // Unit of uncertainty in dip.
  UncertaintyInDipUnit string `json:"uncertainty_in_dip_unit,omitempty"`

  // Uncertainty in elevation above sea level, representing the average uncertainty in elevation for all the coordinates.
  UncertaintyInElevation float64 `json:"uncertainty_in_elevation,omitempty"`

  // Unit of uncertainty in elevation. SI units should be used.
  UncertaintyInElevationUnit string `json:"uncertainty_in_elevation_unit,omitempty"`

  // Uncertainty in the strike.
  UncertaintyInStrike float64 `json:"uncertainty_in_strike,omitempty"`

  // Unit of uncertainty in strike.
  UncertaintyInStrikeUnit string `json:"uncertainty_in_strike_unit,omitempty"`

  // Uncertainty in the x-coordinate (UTM-x or longitude) representing the average uncertainty for all the channels.
  UncertaintyInXCoordinate float64 `json:"uncertainty_in_x_coordinate,omitempty"`

  // Unit of uncertainty in x-coordinate. SI units should be used.
  UncertaintyInXCoordinateUnit string `json:"uncertainty_in_x_coordinate_unit,omitempty"`

  // Uncertainty in the y-coordinate (UTM-y or latitude) representing the average uncertainty for all the channels.
  UncertaintyInYCoordinate float64 `json:"uncertainty_in_y_coordinate,omitempty"`

  // Unit of uncertainty in y-coordinate. SI units should be used.
  UncertaintyInYCoordinateUnit string `json:"uncertainty_in_y_coordinate_unit,omitempty"`

  // Unit of x-coordinate (UTM-x or longitude). Applies to every individual channel nested under the channel group.
  XCoordinateUnit string `json:"x_coordinate_unit"`

  // Unit of y-coordinate (UTM-y or latitude). Applies to every individual channel nested under the channel group.
  YCoordinateUnit string `json:"y_coordinate_unit"`
}

func (strct *ChannelGroup) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
    comma := false
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
    // "ChannelGroupId" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "channel_group_id" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"channel_group_id\": ")
	if tmp, err := json.Marshal(strct.ChannelGroupId); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "channels" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"channels\": ")
	if tmp, err := json.Marshal(strct.Channels); err != nil {
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
    // "CoordinateGenerationDate" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "coordinate_generation_date" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"coordinate_generation_date\": ")
	if tmp, err := json.Marshal(strct.CoordinateGenerationDate); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "CoordinateSystem" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "coordinate_system" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"coordinate_system\": ")
	if tmp, err := json.Marshal(strct.CoordinateSystem); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "depth_below_surface_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"depth_below_surface_unit\": ")
	if tmp, err := json.Marshal(strct.DepthBelowSurfaceUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "dip_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"dip_unit\": ")
	if tmp, err := json.Marshal(strct.DipUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "DistanceAlongFiberUnit" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "distance_along_fiber_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"distance_along_fiber_unit\": ")
	if tmp, err := json.Marshal(strct.DistanceAlongFiberUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "elevation_above_sea_level_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"elevation_above_sea_level_unit\": ")
	if tmp, err := json.Marshal(strct.ElevationAboveSeaLevelUnit); err != nil {
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
    // Marshal the "first_usable_channel_id" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"first_usable_channel_id\": ")
	if tmp, err := json.Marshal(strct.FirstUsableChannelId); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "last_usable_channel_id" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"last_usable_channel_id\": ")
	if tmp, err := json.Marshal(strct.LastUsableChannelId); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "location_method" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"location_method\": ")
	if tmp, err := json.Marshal(strct.LocationMethod); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "ReferenceFrame" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "reference_frame" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"reference_frame\": ")
	if tmp, err := json.Marshal(strct.ReferenceFrame); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "strike_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"strike_unit\": ")
	if tmp, err := json.Marshal(strct.StrikeUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "uncertainty_in_depth" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"uncertainty_in_depth\": ")
	if tmp, err := json.Marshal(strct.UncertaintyInDepth); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "uncertainty_in_depth_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"uncertainty_in_depth_unit\": ")
	if tmp, err := json.Marshal(strct.UncertaintyInDepthUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "uncertainty_in_dip" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"uncertainty_in_dip\": ")
	if tmp, err := json.Marshal(strct.UncertaintyInDip); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "uncertainty_in_dip_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"uncertainty_in_dip_unit\": ")
	if tmp, err := json.Marshal(strct.UncertaintyInDipUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "uncertainty_in_elevation" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"uncertainty_in_elevation\": ")
	if tmp, err := json.Marshal(strct.UncertaintyInElevation); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "uncertainty_in_elevation_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"uncertainty_in_elevation_unit\": ")
	if tmp, err := json.Marshal(strct.UncertaintyInElevationUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "uncertainty_in_strike" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"uncertainty_in_strike\": ")
	if tmp, err := json.Marshal(strct.UncertaintyInStrike); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "uncertainty_in_strike_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"uncertainty_in_strike_unit\": ")
	if tmp, err := json.Marshal(strct.UncertaintyInStrikeUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "uncertainty_in_x_coordinate" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"uncertainty_in_x_coordinate\": ")
	if tmp, err := json.Marshal(strct.UncertaintyInXCoordinate); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "uncertainty_in_x_coordinate_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"uncertainty_in_x_coordinate_unit\": ")
	if tmp, err := json.Marshal(strct.UncertaintyInXCoordinateUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "uncertainty_in_y_coordinate" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"uncertainty_in_y_coordinate\": ")
	if tmp, err := json.Marshal(strct.UncertaintyInYCoordinate); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // Marshal the "uncertainty_in_y_coordinate_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"uncertainty_in_y_coordinate_unit\": ")
	if tmp, err := json.Marshal(strct.UncertaintyInYCoordinateUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "XCoordinateUnit" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "x_coordinate_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"x_coordinate_unit\": ")
	if tmp, err := json.Marshal(strct.XCoordinateUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true
    // "YCoordinateUnit" field is required
    // only required object types supported for marshal checking (for now)
    // Marshal the "y_coordinate_unit" field
    if comma { 
        buf.WriteString(",") 
    }
    buf.WriteString("\"y_coordinate_unit\": ")
	if tmp, err := json.Marshal(strct.YCoordinateUnit); err != nil {
		return nil, err
 	} else {
 		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *ChannelGroup) UnmarshalJSON(b []byte) error {
    cable_idReceived := false
    channel_group_idReceived := false
    coordinate_generation_dateReceived := false
    coordinate_systemReceived := false
    distance_along_fiber_unitReceived := false
    fiber_idReceived := false
    reference_frameReceived := false
    x_coordinate_unitReceived := false
    y_coordinate_unitReceived := false
    var jsonMap map[string]json.RawMessage
    if err := json.Unmarshal(b, &jsonMap); err != nil {
        return err
    }
    // parse all the defined properties
    for k, v := range jsonMap {
        switch k {
        case "cable_id":
            if err := json.Unmarshal([]byte(v), &strct.CableId); err != nil {
                return err
             }
            cable_idReceived = true
        case "channel_group_id":
            if err := json.Unmarshal([]byte(v), &strct.ChannelGroupId); err != nil {
                return err
             }
            channel_group_idReceived = true
        case "channels":
            if err := json.Unmarshal([]byte(v), &strct.Channels); err != nil {
                return err
             }
        case "comment":
            if err := json.Unmarshal([]byte(v), &strct.Comment); err != nil {
                return err
             }
        case "coordinate_generation_date":
            if err := json.Unmarshal([]byte(v), &strct.CoordinateGenerationDate); err != nil {
                return err
             }
            coordinate_generation_dateReceived = true
        case "coordinate_system":
            if err := json.Unmarshal([]byte(v), &strct.CoordinateSystem); err != nil {
                return err
             }
            coordinate_systemReceived = true
        case "depth_below_surface_unit":
            if err := json.Unmarshal([]byte(v), &strct.DepthBelowSurfaceUnit); err != nil {
                return err
             }
        case "dip_unit":
            if err := json.Unmarshal([]byte(v), &strct.DipUnit); err != nil {
                return err
             }
        case "distance_along_fiber_unit":
            if err := json.Unmarshal([]byte(v), &strct.DistanceAlongFiberUnit); err != nil {
                return err
             }
            distance_along_fiber_unitReceived = true
        case "elevation_above_sea_level_unit":
            if err := json.Unmarshal([]byte(v), &strct.ElevationAboveSeaLevelUnit); err != nil {
                return err
             }
        case "fiber_id":
            if err := json.Unmarshal([]byte(v), &strct.FiberId); err != nil {
                return err
             }
            fiber_idReceived = true
        case "first_usable_channel_id":
            if err := json.Unmarshal([]byte(v), &strct.FirstUsableChannelId); err != nil {
                return err
             }
        case "last_usable_channel_id":
            if err := json.Unmarshal([]byte(v), &strct.LastUsableChannelId); err != nil {
                return err
             }
        case "location_method":
            if err := json.Unmarshal([]byte(v), &strct.LocationMethod); err != nil {
                return err
             }
        case "reference_frame":
            if err := json.Unmarshal([]byte(v), &strct.ReferenceFrame); err != nil {
                return err
             }
            reference_frameReceived = true
        case "strike_unit":
            if err := json.Unmarshal([]byte(v), &strct.StrikeUnit); err != nil {
                return err
             }
        case "uncertainty_in_depth":
            if err := json.Unmarshal([]byte(v), &strct.UncertaintyInDepth); err != nil {
                return err
             }
        case "uncertainty_in_depth_unit":
            if err := json.Unmarshal([]byte(v), &strct.UncertaintyInDepthUnit); err != nil {
                return err
             }
        case "uncertainty_in_dip":
            if err := json.Unmarshal([]byte(v), &strct.UncertaintyInDip); err != nil {
                return err
             }
        case "uncertainty_in_dip_unit":
            if err := json.Unmarshal([]byte(v), &strct.UncertaintyInDipUnit); err != nil {
                return err
             }
        case "uncertainty_in_elevation":
            if err := json.Unmarshal([]byte(v), &strct.UncertaintyInElevation); err != nil {
                return err
             }
        case "uncertainty_in_elevation_unit":
            if err := json.Unmarshal([]byte(v), &strct.UncertaintyInElevationUnit); err != nil {
                return err
             }
        case "uncertainty_in_strike":
            if err := json.Unmarshal([]byte(v), &strct.UncertaintyInStrike); err != nil {
                return err
             }
        case "uncertainty_in_strike_unit":
            if err := json.Unmarshal([]byte(v), &strct.UncertaintyInStrikeUnit); err != nil {
                return err
             }
        case "uncertainty_in_x_coordinate":
            if err := json.Unmarshal([]byte(v), &strct.UncertaintyInXCoordinate); err != nil {
                return err
             }
        case "uncertainty_in_x_coordinate_unit":
            if err := json.Unmarshal([]byte(v), &strct.UncertaintyInXCoordinateUnit); err != nil {
                return err
             }
        case "uncertainty_in_y_coordinate":
            if err := json.Unmarshal([]byte(v), &strct.UncertaintyInYCoordinate); err != nil {
                return err
             }
        case "uncertainty_in_y_coordinate_unit":
            if err := json.Unmarshal([]byte(v), &strct.UncertaintyInYCoordinateUnit); err != nil {
                return err
             }
        case "x_coordinate_unit":
            if err := json.Unmarshal([]byte(v), &strct.XCoordinateUnit); err != nil {
                return err
             }
            x_coordinate_unitReceived = true
        case "y_coordinate_unit":
            if err := json.Unmarshal([]byte(v), &strct.YCoordinateUnit); err != nil {
                return err
             }
            y_coordinate_unitReceived = true
        }
    }
    // check if cable_id (a required property) was received
    if !cable_idReceived {
        return errors.New("\"cable_id\" is required but was not present")
    }
    // check if channel_group_id (a required property) was received
    if !channel_group_idReceived {
        return errors.New("\"channel_group_id\" is required but was not present")
    }
    // check if coordinate_generation_date (a required property) was received
    if !coordinate_generation_dateReceived {
        return errors.New("\"coordinate_generation_date\" is required but was not present")
    }
    // check if coordinate_system (a required property) was received
    if !coordinate_systemReceived {
        return errors.New("\"coordinate_system\" is required but was not present")
    }
    // check if distance_along_fiber_unit (a required property) was received
    if !distance_along_fiber_unitReceived {
        return errors.New("\"distance_along_fiber_unit\" is required but was not present")
    }
    // check if fiber_id (a required property) was received
    if !fiber_idReceived {
        return errors.New("\"fiber_id\" is required but was not present")
    }
    // check if reference_frame (a required property) was received
    if !reference_frameReceived {
        return errors.New("\"reference_frame\" is required but was not present")
    }
    // check if x_coordinate_unit (a required property) was received
    if !x_coordinate_unitReceived {
        return errors.New("\"x_coordinate_unit\" is required but was not present")
    }
    // check if y_coordinate_unit (a required property) was received
    if !y_coordinate_unitReceived {
        return errors.New("\"y_coordinate_unit\" is required but was not present")
    }
    return nil
}
