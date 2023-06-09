## TESTS ##

The golang test mechanism is used to validate delta information.

Related groups of _checks_ are split into different files with the `golang` naming
scheme of `xxxx_test.go`. Each test should be in a function called `TestXXX`.

The tests are written as a series of function calls based around the template:

```golang
func TestXXXX(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range xxxxChecks {
		t.Run(k, v(set))
	}
}
```

where the specific `xxxx` checks are given as a map of functions, e.g.:

```golang
var xxxxChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check for xxxx in situtation yyyy": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

                /* test goes here and can use "set" to access delta fields */

                }
        },
        ...
}
```

Tests can be run using `go test` in the directory containing the test files.

### Enforcing Checks ###

These checks will produce error messages which will generally mean the data
will not be allowed to be merged into the repository.

__antennas_test.go__

* check antenna installation overlap
* check for antenna installation mark overlaps
* check for missing antenna marks
* check for missing antenna assets
* check for missing antenna sessions

__assets_test.go__

* check for duplicate asset numbers
* check for duplicate equipment

__calibrations_test.go__

* check for calibration overlaps
* check for missing assets

__cameras_test.go__

* check for cameras installation equipment overlaps
* check for cameras installation mount overlaps
* check for cameras installation views
* check cameras assets

__citations_test.go__

* check for duplicate citations

__combined_test.go__

* check for sensor/recorder installation location overlaps

__connections_test.go__

* check for connection overlaps
* check for connection span mismatch
* check for missing connection stations
* check for missing connection sites
* check for missing connection site locations
* check for missing datalogger places
* check for missing sensor connections
* check for missing datalogger connections
* check for gauge duplication
* check for missing constituent gauges

__dataloggers_test.go__

* check for datalogger installation place overlaps
* check for datalogger installation equipment overlaps
* check for missing datalogger assets

__doases_test.go__

* check for doases installation equipment overlaps
* check for doases installation mount overlaps
* check for doases installation views
* check doases assets

__features_test.go__

* check for duplicated site features
* check for duplicated features
* check for duplicated feature sites

__firmware_test.go__

* check for firmware history overlaps
* check for firmware non-changes
* check for firmware assets
* check for latest installed receiver firmware

__gains_test.go__

* check for gain installation overlaps
* check for missing sites

__gauges_test.go__

* check for gauge duplication
* check for missing gauge stations

__marks_test.go__

* check for duplicated marks

__metsensors_test.go__

* check for metsensors installation equipment overlaps
* check for missing metsensor marks
* check for missing metsensor assets

__monuments_test.go__

* check for duplicated monuments
* check for monument ground relationships
* check for monument types

__networks_test.go__

* check for duplicated networks

__placenames_test.go__

* check for placename duplication
* check for placename latitude longitudes
* check for placename levels

__polarities_test.go__

* check for duplicate polarities

__preamps_test.go__

* check for duplicate preamps

__radomes_test.go__

* check for radomes installation equipment overlaps
* check for overlapping radomes installations
* check for missing radome marks
* check for missing radome assets

__receivers_test.go__

* check for receiver installation overlaps
* check for receiver installation equipment overlaps
* check for missing receiver marks
* check for missing receiver assets
* check for recorder installation overlaps
* check for invalid sensor azimuth
* check for invalid sensor dip
* check for missing recorder stations
* check for missing recorder sites
* check for missing assets

__resps_test.go__

* check for component response files
* check for channel response files

__samples_test.go__

* check for duplicate sampling site
* check for missing sample networks

__sensors_test.go__

* check for sensor installation overlaps
* check for invalid sensor azimuth
* check for invalid sensor dip
* check for missing sensor stations
* check for missing sensor sites
* check for missing assets

__sessions_test.go__

* check session overlap
* check session spans
* check session satellite system
* check session marks

__sites_test.go__

* check for duplicated sites
* check for duplicated station sites

__stations_test.go__

* check for duplicated stations
* check for missing networks

__streams_test.go__

* check for invalid axial labels
* check for invalid stream span overlaps
* check for invalid stream spans
* check for invalid stream sample rates
* check for invalid stream stations
* check for invalid dates: stream within station
* check for invalid stream sites
* check for invalid stream locations
* check for invalid dates: stream within site
* check for invalid stream sensor sites
* check for invalid stream recorder sites

__telemetry_test.go__

* check for overlapping telemeties
* check for zero gain telemetries

__views_test.go__

* check for duplicated views

### Non-Enforcing Checks ###

There are a number of checks which are non-enforcing and will only produce
warning messages when the verbose flag is used (`-v`) or when an unrelated
error is encountered in which case these warnings will be bundled into the
output.

Once these checks pass without error it is intended that they become enforcing checks.
