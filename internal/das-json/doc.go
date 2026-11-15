//go:generate bash -c "go run $(find generate -name \"*.go\" -and -not -name \"*_test.go\" -maxdepth 1) -version 2.0 -schema https://raw.githubusercontent.com/FDSN/DAS-metadata/refs/heads/main/schema/DAS-Metadata.v2.0.schema.json -insecure -output dasv2.0"

// Automatically generate das-json bindings from the schema files.
package dasjson
