# IMPACT

Provide configuration for real-time shaking information

## Output

The output json file provides a lookup for each expected stream, the information is encoded as a
hash with the key being the stream name, i.e. *<NN>_<SSS>_<LL>_<CCC>*.

The related fields are:

 * Longitude
 * Latitude
 * Q
 * Rate
 * Gain
 * Name

## Configuration

Each desired sample rate requires a filter _Q_ value given in the _build.go_ file. 

Each stream requires its Datalogger and Sensor models to be listed as _blessed_ models. This list is
maintained and can be updated in the _blessed.go_ file.
