# sensor-config

Build a sensor description file

## output

The output XML or JSON file provides a description of all sensor installations and manual collection points.
This is aimed at a mapping or site discovery service.

## usage

  ./sensor-config [options]

## options

  -base string
    	delta base files
  -building string
    	building camera network code (default "BC")
  -coastal string
    	coastal tsunami gauge network code (default "TG")
  -combined value
    	combined sensor location codes (default ^[123])
  -dart string
    	dart buoy network code (default "TD")
  -doas string
    	doas network code (default "EN")
  -enviro string
    	envirosensor network code (default "EN")
  -geomag value
    	geomag sensor codes (default ^5)
  -json
    	use JSON for output format
  -lentic string
    	lentic tsunami gauge network code (default "LG")
  -manual string
    	manualcollect network code (default "MC")
  -networks string
    	installed network codes (default "AK,CB,CH,EC,FI,HB,KI,NM,NZ,OT,RA,RT,SC,SI,SM,SP,TP,TR,WL")
  -output string
    	output sensor description file
  -volcano string
    	volcano camera network code (default "VC")
  -water value
    	water pressue sensor codes (default ^4)

