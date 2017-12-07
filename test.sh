#!/bin/bash -x

errcount=0

error_handler () {
    echo "Trapped error - ${1:-"Unknown Error"}" 1>&2
    (( errcount++ ))       # or (( errcount += $? ))
}

trap error_handler ERR

echo "go formatting" 1>&2
test -z "$(find ./meta ./resp ./tests ./tools -name "*.go" -exec gofmt -l {} \; | tee /dev/stderr)"

echo "go vetting" 1>&2
go vet ./meta ./resp ./tests ./tools/...

echo "go testing" 1>&2
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

echo "go building" 1>&2
go build ./tools/stationxml
go build ./tools/altus
go build ./tools/cusp
go build ./tools/amplitude
go build ./tools/spectra
go build ./tools/chart
go build ./tools/impact
go build ./tools/rinexml

exit $errcount

# vim: tabstop=4 expandtab shiftwidth=4 softtabstop=4
