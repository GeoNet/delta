#!/bin/bash -x

errcount=0

error_handler () {
    echo "Trapped error - ${1:-"Unknown Error"}" 1>&2
    (( errcount++ ))       # or (( errcount += $? ))
}

trap error_handler ERR

mkdir -p .tmp/geonet-meta/stationxml || exit 255
go build ./tools/stationxml || exit 255

./stationxml -base . -output .tmp/geonet-meta/stationxml/complete.xml
./stationxml -base . -output .tmp/geonet-meta/stationxml/scp.xml \
    -active \
    -installed \
    -operational \
    -networks '!(SB|.X)' \
    -channels '([EHB][HN][ZNE12])'

./stationxml -base . -output .tmp/geonet-meta/stationxml/iris.xml \
    -stations '(KHZ|QRZ|OUZ|HIZ|BKZ|ODZ|BFZ|CTZ|URZ|RPZ|WPVZ)' \
    -sensors '(STS-2|CMG-3TB|CMG-40T-60S|FBA-ES-T)' \
    -dataloggers '(Q330HR/6|Q4120/6|Q330/3)' \
    -channels '([HLV]H[ZNE12]|[HBL]N[ZNE])'

mkdir -p .tmp/geonet-meta/seed/pod || exit 255
go build ./tools/pod || exit 255

for input in .tmp/geonet-meta/stationxml/*.xml; do
    output=$(basename $input .xml)
    ./pod -output .tmp/pod/$output $input
    (cd .tmp/pod/$output; tar cfz ../../geonet-meta/seed/pod/$output.tar.gz HDR000)
done

mkdir -p .tmp/geonet-meta/config || exit 255
go build ./tools/impact || exit 255

./impact -base . -channels "[BH]N[ZNE12XY]" -output .tmp/geonet-meta/config/impact.json

mkdir -p .tmp/geonet-meta/config/rinex || exit 255
go build ./tools/rinexml || exit 255

./rinexml -install ./install -network ./network -output .tmp/geonet-meta/config/gnsssitexml

mkdir -p .tmp/geonet-meta/config/gloria || exit 255
go build ./tools/gloria || exit 255

./gloria --install install --network network -output .tmp/geonet-meta/config/gloria

mkdir -p .tmp/geonet-meta/config/sit || exit 255
go build ./tools/sit || exit 255

./sit -install install -network network -asset assets -output .tmp/geonet-meta/config/sit

exit $errcount

# vim: tabstop=4 expandtab shiftwidth=4 softtabstop=4
