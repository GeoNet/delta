#!/bin/bash -x

errcount=0

error_handler () {
    echo "Trapped error - ${1:-"Unknown Error"}" 1>&2
    (( errcount++ ))       # or (( errcount += $? ))
}

trap error_handler ERR

mkdir -p .tmp/geonet-meta/stationxml || exit 255

(cd ./tools/stationxml; go build; ./stationxml -output ../../.tmp/geonet-meta/stationxml/complete.xml)
(cd ./tools/stationxml; go build; ./stationxml -output ../../.tmp/geonet-meta/stationxml/scp.xml \
    -operational \
    -channels '([EHB][HN][ZNE12])' \
)
(cd ./tools/stationxml; go build; ./stationxml -output ../../.tmp/geonet-meta/stationxml/iris.xml \
    -stations '(KHZ|QRZ|OUZ|HIZ|BKZ|ODZ|BFZ|CTZ|URZ|RPZ|WPVZ)' \
    -sensors '(STS-2|CMG-3TB|CMG-40T-60S|FBA-ES-T)' \
    -dataloggers '(Q330HR/6|Q4120/6|Q330/3)' \
    -channels '([HLV]H[ZNE12]|[HBL]N[ZNE])' \
    -output ../../.tmp/geonet-meta/stationxml/iris.xml
)

mkdir -p .tmp/geonet-meta/seed/pod || exit 255

for output in complete scp iris; do
    (cd ./tools/pod; go build; ./pod --output  ../../.tmp/pod/$output ../../.tmp/geonet-meta/stationxml/$output.xml)
    (cd .tmp/pod/$output; tar cfz ../../geonet-meta/seed/pod/$output.tar.gz HDR000)
done

exit $errcount

# vim: tabstop=4 expandtab shiftwidth=4 softtabstop=4
