## ENVIRONMENT ##

_Lists the geological and physical environment details of collection points._

### FILES ###

* `geology.csv` - Geologic Characteristic of GNSS marks.
* `marksnotes.csv` - additional notes for GNSS marks (monument description).
* `visibility.csv` - Sky View Descriptions.


#### _GEOLOGY_ ####

| Field | Description |
| --- | --- |
| _Code_ | Code to uniquely identify GNSS _Mark_ (or recording _Station_)
| _Geologic Characteristic_ | bedrock/clay/gravel/sand/sediments/conglomerate/etc
| _Bedrock Type_ | igneous/metamorphic/sedimentary
| _Bedrock Condition_ | fresh/jointed/weathered
| _Fracture Spacing_ | 1-10 cm/11-50 cm/51-200 cm/over 200 cm
| _Fault zones nearby_  | yes/no/name of the zone
| _Additional info | Free form description for additional information

#### _VISIBILITY_ ####

| Field | Description |
| --- | --- |
| _Code_ | Code to uniquely identify GNSS _Mark_ (or recording _Station_)
| _Sky Visibility_ | Free form description of the site sky visibility and obstructions
| _Start Date_ | General date and time at which the visibility was accurate
| _End Date_ | General date and time at which the visibility was no longer accurate



### CHECKS ###

Pre-commit checks will be made on these files to ensure:

### NOTES ###

Dates should be given as in _ISO 8601_ (i.e. `2016-09-18T02:24:26Z`), future dates should be given in the form: `9999-01-01T00:00:00Z`.

