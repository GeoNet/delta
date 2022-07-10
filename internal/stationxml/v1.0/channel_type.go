package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type ChannelType struct {
	BaseNodeType

	LocationCode string `xml:"locationCode,attr"`

	ExternalReference []ExternalReferenceType `xml:"ExternalReference,omitempty"`

	Latitude LatitudeType `xml:"Latitude"`

	Longitude LongitudeType `xml:"Longitude"`

	Elevation DistanceType `xml:"Elevation"`

	Depth DistanceType `xml:"Depth"`

	Azimuth *AzimuthType `xml:"Azimuth,omitempty"`

	Dip *DipType `xml:"Dip,omitempty"`

	Type []Type `xml:"Type,omitempty"`

	SampleRate SampleRateType `xml:"SampleRate"`

	SampleRateRatio *SampleRateRatioType `xml:"SampleRateRatio,omitempty"`

	StorageFormat string `xml:"StorageFormat,omitempty"`

	ClockDrift *ClockDrift `xml:"ClockDrift,omitempty"`

	CalibrationUnits *UnitsType `xml:"CalibrationUnits,omitempty"`

	Sensor *EquipmentType `xml:"Sensor,omitempty"`

	PreAmplifier *EquipmentType `xml:"PreAmplifier,omitempty"`

	DataLogger *EquipmentType `xml:"DataLogger,omitempty"`

	Equipment *EquipmentType `xml:"Equipment,omitempty"`

	Response *ResponseType `xml:"Response,omitempty"`
}
