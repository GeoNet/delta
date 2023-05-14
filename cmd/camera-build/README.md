# camera-build

Provide configuration for labelling remote camera images.

## output

The output JSON file provides a lookup for each camera _view_ to allow a simple caption.

Build a camera caption configuration from delta meta information

## usage

    camera-build [options]

## options

    -active
        only output active camera information
    -base string
        base for custom delta files
    -networks string
        comma separated list of networks, an empty value matches all networks
    -output string
        where to store json formatted output
    -verbose
        make noise
