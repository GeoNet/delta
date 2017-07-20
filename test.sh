#!/bin/bash -x

errcount=0

error_handler () {
    echo "Trapped error - ${1:-"Unknown Error"}" 1>&2
    (( errcount++ ))       # or (( errcount += $? ))
}

trap error_handler ERR

go test ./meta
go test ./tests
go test ./tools/stationxml
go test ./tools/altus
go test ./tools/cusp
go test ./tools/amplitude
go test ./tools/spectra
go test ./tools/chart
go test ./tools/impact
go test ./tools/rinexml
go test .

exit $errcount

# vim: tabstop=4 expandtab shiftwidth=4 softtabstop=4
