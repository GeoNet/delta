package stationxml

/*
 *  WARNING: CODE GENERATED AUTOMATICALLY.
 *
 *  Use "go generate" to update these files.
 *
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 *
 */

type EquipmentType struct {
	ResourceId string `xml:"resourceId,attr,omitempty"`

	Type string `xml:"Type,omitempty"`

	Description string `xml:"Description,omitempty"`

	Manufacturer string `xml:"Manufacturer,omitempty"`

	Vendor string `xml:"Vendor,omitempty"`

	Model string `xml:"Model,omitempty"`

	SerialNumber string `xml:"SerialNumber,omitempty"`

	InstallationDate DateTime `xml:"InstallationDate,omitempty"`

	RemovalDate DateTime `xml:"RemovalDate,omitempty"`

	CalibrationDate []DateTime `xml:"CalibrationDate,omitempty"`
}
