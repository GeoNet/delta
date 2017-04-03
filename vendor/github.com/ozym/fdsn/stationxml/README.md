# FDSN - StationXML

A golang wrapper package for the __FDSN__ _StationXML_ format for describing seismic data collection parameters.

See http://www.fdsn.org/xml/station/ for schema details.

Currently the classes are capable of marshalling / unmarshalling a set of test _StationXML_ files derived from the NZ
__FDSN__ webservice.

Expected enhancements:

* Validation
* Enumeration of selections
* More tests

Progress: IsValid()

- [x] RootType
- [x] NetworkType
- [x] StationType
- [x] ChannelType
- [x] GainType
- [x] FrequencyRangeGroup
- [x] SensitivityType
- [x] EquipmentType
- [x] ResponseStageType
- [x] LogType
- [x] CommentType
- [x] PolesZerosType
- [x] FIRType
- [x] CoefficientsType
- [x] ResponseListElementType
- [x] ResponseListType
- [x] PolynomialType
- [x] DecimationType
- [x] uncertaintyDouble
- [x] FloatNoUnitType
- [x] FloatType
- [x] SecondType
- [x] VoltageType
- [x] AngleType
- [x] LatitudeBaseType
- [x] LatitudeType
- [x] LongitudeBaseType
- [x] LongitudeType
- [x] AzimuthType
- [x] DipType
- [x] DistanceType
- [x] FrequencyType
- [x] SampleRateGroup
- [x] SampleRateType
- [x] SampleRateRatioType
- [x] PoleZeroType
- [x] CounterType
- [x] PersonType
- [x] SiteType
- [x] ExternalReferenceType
- [x] NominalType
- [x] EmailType
- [x] PhoneNumberType
- [x] RestrictedStatusType
- [x] UnitsType
- [x] BaseFilterType
- [x] ResponseType
- [x] BaseNodeType
