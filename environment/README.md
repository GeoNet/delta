## ENVIRONMENT ##

_Lists the geological and physical environment details of collection points._

### FILES ###

* `visibility.csv` - Sky View Descriptions.

#### _VISIBILITY_ ####

| Field | Description |
| --- | --- |
| _Code_ | Code to uniquely identify GNSS _Mark_ (or recording _Station_)
| _Sky Visibility_ | Free form description of the site sky visibility and obstructions
| _Start Date_ | General date and time at which the visibility was accurate
| _Start Date_ | General date and time at which the visibility was no longer accurate

### CHECKS ###

Pre-commit checks will be made on these files to ensure:

### NOTES ###

Dates should be given as in _ISO 8601_ (i.e. `2016-09-18T02:24:26Z`), future dates should be given in the form: `9999-01-01T00:00:00Z`.

