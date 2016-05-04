// Package fdsn/stationxml provides a wrapper around the FDSN StationXML schema.
//
// FDSN StationXML schema. Designed as an XML representation of SEED metadata, the schema maps to
// the most important and commonly used structures of SEED 2.4. When definitions and usage are
// underdefined the SEED manual should be referred to for clarification.
//
// FDSN StationXML (www.fdsn.org/xml/station)
//
// The purpose of this schema is to define an XML representation of the most important
// and commonly used structures of SEED 2.4 metadata.
//
// The goal is to allow mapping between SEED 2.4 dataless SEED volumes and this schema with as
// little transformation or loss of information as possible while at the same time simplifying
// station metadata representation when possible.  Also, content and clarification has been added
// where lacking in the SEED standard.
//
// When definitions and usage are underdefined the SEED manual should be referred to for
// clarification.  SEED specifiation: http://www.fdsn.org/publications.htm
//
// Another goal is to create a base schema that can be extended to represent similar data types.
//
//
// Versioning for FDSN StationXML:
//
// The 'version' attribute of the schema definition identifies the version of the schema.  This
// version is not enforced when validating documents.
//
// The required 'schemaVersion' attribute of the root element identifies the version of the schema
// that the document is compatible with.  Validation only requires that a value is present but
// not that it matches the schema used for validation.
//
// The targetNamespace of the document identifies the major version of the schema and document,
// version 1.x of the schema uses a target namespace of "http://www.fdsn.org/xml/station/1".
// All minor versions of a will be backwards compatible with previous minor releases.  For
// example, all 1.x schemas are backwards compatible with and will validate documents for 1.0.
// Major changes to the schema that would break backwards compabibility will increment the major
// version number, e.g. 2.0, and the namespace, e.g. "http://www.fdsn.org/xml/station/2".
//
// This combination of attributes and targetNamespaces allows the schema and documents to be
// versioned and allows the schema to be updated with backward compatible changes (e.g. 1.2)
// and still validate documents created for previous major versions of the schema (e.g. 1.0).
package stationxml
