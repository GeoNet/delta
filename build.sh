#!/bin/bash -x

errcount=0

error_handler () {
    echo "Trapped error - ${1:-"Unknown Error"}" 1>&2
    (( errcount++ ))       # or (( errcount += $? ))
}

trap error_handler ERR

mkdir -p .tmp/geonet-meta/stationxml || exit 255

(cd ./tools/stationxml; go build; ./stationxml -output ../../.tmp/geonet-meta/stationxml/complete.xml)

exit $errcount

# vim: tabstop=4 expandtab shiftwidth=4 softtabstop=4
