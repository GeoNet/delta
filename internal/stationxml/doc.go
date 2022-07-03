//go:generate bash -c "go run generate/*.go -version v1.0 -schema https://www.fdsn.org/xml/station/fdsn-station-1.0.xsd -insecure"
//go:generate bash -c "go run generate/*.go -version v1.1 -schema https://www.fdsn.org/xml/station/fdsn-station-1.1.xsd -insecure"
//go:generate bash -c "go run generate/*.go -version v1.2 -schema https://www.fdsn.org/xml/station/fdsn-station-1.2.xsd -insecure"
package stationxml
