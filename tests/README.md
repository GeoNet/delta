## TESTS ##

The golang test mechanism is used to validate delta information.

Related groups of _checks_ are split into different files with the `golang` naming
scheme of `xxxx_test.go`. Each test should be in a function called `TestXXX`.

There is a helper function `loadListFile` which can reduce the amount of code needed to
load a `csv` formatted `List` file. This is only used inside the testing framework.

Generally the layout of the tests will be along the lines of:

```golang
func TestStations(t *testing.T) {

  var stations meta.StationList
  loadListFile(t, "../network/stations.csv", &stations)

  t.Run("check for duplicated stations", func(t *testing.T) {
    // actual test
  }

  ...

}
```

Tests can be run using `go test` in the directory containing the test files.

### Enforcing Checks ###

These checks will produce error messages which will generally mean the data
will not be allowed to be merged into the repository.

__antennas_test.go__

  * check for antenna installation equipment overlaps

__antennas_test.go__

  * check for antenna installation mark overlaps
  * check for missing antenna marks
  * check for missing antenna assets
  * check for missing antenna sessions

__assets_test.go__

  * check for duplicate asset numbers
  * check for duplicate equipment

__cameras_test.go__

  * check for cameras installation equipment overlaps
  * check for missing camera mounts
  * Load camera assets file

__combined_test.go__

  * check for sensor/recorder installation location overlaps

__connections_test.go__

  * check for connection overlaps
  * check for missing connection stations
  * check for missing connection sites
  * check for missing connection locations
  * check for connection span mismatch
  * check for missing connection places
  * check for missing sensor connections
  * check for missing datalogger connections

__consistency_test.go__

  * check file consistency

__constituents_test.go__

  * check for constituent duplications
  * check for missing constituent gauges

__dataloggers_test.go__

  * check for datalogger installation place overlaps
  * check for datalogger installation equipment overlaps
  * check for missing datalogger assets

__firmware_test.go__

  * check for firmware history overlaps
  * check for firmware receiver assets
  * check for latest installed receiver firmware

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

  * check for monument duplication
  * check for monument ground relationships
  * check for monument types

__networks_test.go__

  * check for network duplication

__radomes_test.go__

  * check for radomes installation equipment overlaps
  * check for overlapping radomes installations
  * check for missing radome marks
  * check for missing radome assets

__receivers_test.go__

  * check for particular receiver installation overlaps
  * check for receiver installation equipment overlaps
  * check for missing receiver marks
  * check for missing receiver assets

__recorders_test.go__

  * check for recorder installation equipment overlaps
  * check for missing recorder stations
  * check for recorder assets

__sensors_test.go__

  * check for missing sensors
  * check for sensor installation overlaps
  * check for missing sensor stations
  * check for missing sensor sites
  * check for invalid sensor azimuth
  * check for invalid sensor dip
  * check for missing sensor assets

__sessions_test.go__

  * check for session overlaps
  * check for missing session marks
  * check for session span mismatches
  * check for unknown session satellite systems

__sites_test.go__

  * check for duplicated sites
  * check for invalid dates: site within station

__stations_test.go__

  * check for duplicated stations

__streams_test.go__

  * check for invalid stream sample rates
  * check for invalid stream spans
  * check for invalid stream stations
  * check for invalid stream sites
  * check for invalid stream locations
  * check for invalid stream spans
  * check for missing recorder streams
  * check for missing sensor streams


### Non-Enforcing Checks ###

There are a number of checks which are non-enforcing and will only produce
warning messages when the verbose flag is used (`-v`) or when an unrelated
error is encountered in which case these warnings will be bundled into the
output.

Once these checks pass without error it is intended that they become enforcing checks.

__sites_test.go__

  * check for duplicated sites
  * check for invalid dates: site within station

__streams_test.go__

  * check for invalid dates: stream within station
  * check for invalid dates: stream within site
