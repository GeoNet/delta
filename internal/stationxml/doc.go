//go:generate bash -c "go run $(find generate -name \"*.go\" -and -not -name \"*_test.go\" -maxdepth 1) -version 1.0 -schema https://www.fdsn.org/xml/station/fdsn-station-1.0.xsd -insecure -output v1.0 -future -format 2006-01-02T15:04:05"
//go:generate bash -c "go run $(find generate -name \"*.go\" -and -not -name \"*_test.go\" -maxdepth 1) -version 1.1 -schema https://www.fdsn.org/xml/station/fdsn-station-1.1.xsd -insecure -output v1.1 -future -format 2006-01-02T15:04:05"
//go:generate bash -c "go run $(find generate -name \"*.go\" -and -not -name \"*_test.go\" -maxdepth 1) -version 1.2 -schema https://www.fdsn.org/xml/station/fdsn-station-1.2.xsd -insecure -output v1.2"

// Automatically generate stationxml bindings from the schema files.
package stationxml
