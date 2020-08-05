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
