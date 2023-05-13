# sc3-build

Build a set of sc3 key files

### overview

The `sc3-build` application parses the delta metadata information and builds the `key` files used for SeisComP configuration.

### usage

The `delta` meta-data files can be selected using the `base` option by pointing to the appropriate `delta` directory.
The internal compiled files will be used if no base is set, this is mainly used for automatic file generation.

The `output` option selects a directory to store the generated files in. If the `purge` flag is set, then this directory
is first scanned for files and after generating config files, any not updated will be removed.

When selecting stations a `grace` period, in days, can be selected. This is to handle the situation where a station may have
been decommissioned but its data may still be valid and requires further processing. After the grace period this station
will no longer be selected.

Network selection is done by the `network` option. This can either be single network codes, e.g. `-network NZ -network CH` or
a comma separated list, e.g. `-network NZ,CH`.  If a network is to be excluded than a `!` can be used as a prefix, e.g. `-network !CH`.
Networks will be ignored if they have been explicitly excluded, otherwise if there are no actual networks given to include then all
other networks will be selected. If, however, at least one network is selected to be included then only those networks will be used.
External networks can be ignored using the station selection mechanisms using wildcards.

Station selection is similar to `network` selection. It uses a similar `station` option but expects values matching `NN_SSSS` or `SSSS`. 
Wildcards can be used, e.g. (`*_*` will select all stations, `SSSS` is just an alias for `*_SSSS`, and `NN_*` will select all stations
in the `NN` network, the network given is the external network as inserted into the key files. Again a prefix of `!` will exclude a selection.


```
Usage:

  sc3-build [options]

Options:

  -base string
    	delta base files
  -grace int
    	allow for a grace period in days after site changes (default 30)
  -network value
    	add specific network(s), will skip all others (use ! prefix to exclude specific network)
  -output string
    	output sc3 configuration directory (default "key")
  -purge
    	remove unknown single xml files
  -station value
    	add specific station(s) (requires SSSS, NN_SSSS, *_SSSS, or NN_*) (use ! prefix to exclude specific station) 

```
