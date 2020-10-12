package main

import (
	"encoding/xml"
)

type FormInformation struct {
	XMLName xml.Name `xml:"formInformation"`

	PreparedBy   string `xml:"mi:preparedBy"`
	DatePrepared string `xml:"mi:datePrepared"`
	ReportType   string `xml:"mi:reportType"`
}

type FormInformationInput struct {
	XMLName xml.Name `xml:"formInformation"`

	PreparedBy   string `xml:"preparedBy"`
	DatePrepared string `xml:"datePrepared"`
	ReportType   string `xml:"reportType"`
}
