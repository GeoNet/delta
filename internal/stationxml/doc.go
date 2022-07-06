//go:generate bash -c "go run generate/*.go -version v1.0 -schema https://www.fdsn.org/xml/station/fdsn-station-1.0.xsd -insecure -output v1.0"
//go:generate bash -c "go run generate/*.go -version v1.1 -schema https://www.fdsn.org/xml/station/fdsn-station-1.1.xsd -insecure -output v1.1"
//go:generate bash -c "go run generate/*.go -version v1.2 -schema https://www.fdsn.org/xml/station/fdsn-station-1.2.xsd -insecure -output v1.2"

// Automatically generate stationxml bindings from the schema files.
package stationxml
