# tilde-build

Build the tilde domain config from delta.

### overview

The `tilde-build` application parses the delta metadata information and builds the config XML file for updating `tilde`.

### usage

The `delta` meta-data files can be selected using the `base` option by pointing to the appropriate `delta` directory.
The internal compiled files will be used if no base is set, this is mainly used for automatic file generation.

The `output` option selects a file to store the generated XML config, if no option is given the config is sent to `stdout`.

Each `tilde` domain is encoded using a set of network codes which can be comma separated if there are more than one network
associated with any domain.

```
Usage:

  tilde-build [options]

Options:

  -base string
    	delta base files
  -coastal string
    	coast tsunami gauge network code (default "TG,LG")
  -dart string
    	dart buoy network code (default "TD")
  -enviro string
    	envirosensor network code (default "EN")
  -manual string
    	manualcollect network code (default "MC")
  -output string
    	output tilde configuration file

```
