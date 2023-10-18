## REFERENCES ##

_Reference citations and external source information._

### FILES ###

Reference information for the GeoNet sensor networks and manual data collection systems.
 
* `citations.csv` - Reference citations
* `tilde_methods.csv` - Tilde  method descriptions

#### _CITATIONS_ ####

A list of _reference_ citations to published or otherwise information.
All fields, other than Key, are optional and can be left blank.

| Field | Description |
| --- | --- |
| _Key_ | The citation key used where appropriate
| _Author_ | Reference author or authors in a natural format
| _Year_ | Reference year, if appropriate
| _Title_ | The title of the reference in a natural format
| _Published_ | Where reference was published if appropriate and known
| _Volume_ | Series information for published reference if relavent
| _Pages_ | Optional page information for the published reference
| _DOI_ | The reference DOI (_Digital Object Identifier_) if known
| _Link_ | A URL link to the reference if available
| _Retrieved_ | The last time a valid URL was retrieved

##### EXAMPLE ######

    Key,Author,Year,Title,Published,Volume,Pages,DOI,Link,Retrieved
    Fry2020,"Fry, B., S.-J. McCurrach, K. Gledhill, W. Power, M. Williams, M. Angove, D. Arcas, and C. Moore",2020,Sensor network warns of stealth tsunamis,Eos,101,,https://doi.org/10.1029/2020EO144274,,

#### _METHODS_ ####
The range of methods used in the Tilde application is diverse. In some cases, the details of a method are well
know only to those familiar with a data set or with a similar data set from another institution.
This does not make it easy for non-specialists to use the data.

A list of methods used to collect the data is provided that contains a brief description of the method, and where available, a URL
link to publicly available resources that provide additional information.

| Field | Description |
| --- | --- |
| _Domain_ | The data domain. This is a concept used by the Tilde data appication, and refers to the broad collection method or data discipline |
| _Method_ | The data collection method. How the data are sampled, collected, processed, or analysed |
| _Description_ | A brief description of the data collection method |
| _Reference_ | A publicly available URL that provides additional information on the method |

### NOTES ###

- Dates should be given as in _ISO 8601_ (i.e. `2016-09-18T02:24:26Z`)
- Any field entries with commas should be enclosed within double quotes. (e.g. "Last, F., Name ...")
