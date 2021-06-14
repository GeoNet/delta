# Tsunami Metadata Conventions

In this document we describe the metadata conventions present in the GeoNet tsunami data and metadata.

## Overview

GeoNet has two types of tsunami data stations: tsunami gauges and DART buoys. The distinction is made by the primary purpose of the station; whether it exists principally to record tsunami arrivals on-shore (via a tsunami gauge) or in the deep ocean (via a DART buoy).

The management of DART data and metadata is an evolving space, and until this is more established we refrain from documenting it here.

Tsunami gauge data and metadata mostly follow the conventions of seismic data and metadata, and unless stated otherwise in this document these conventions are what apply to tsunami gauge data and metadata. The sibling document `SEISMIC_METADATA_CONVENTIONS.md` details these conventions. Here we make reference only to the differences between the tsunami gauge metadata conventions and those of seismic metadata.

## Location Code

For tsunami gauges, the location code does not denote position or type as with seismic sensors, but rather is used only to indicate one sensor from another at the station. For example, a station may have sites at location codes `41` and `42`, but these location code differences are only significant for distinguishing the respective data streams from the sensors at each site. In this way, we can think of the location code for tsunami gauges as a "sensor code". 

## Naming Conventions When Moving or Installating Different Sensors

Where in seismic data small changes to sensor position can have substantial impacts on the data, the same is not true for tsunami gauges. 

Tsunami gauges are sensitive to the overlying water column, which is principally influenced by sensor depth and tides. 

If a sensor position changes such that new constituents describing the tides at that sensor location do not need to be derived, the station code does not change.

If a sensor position changes such that the sensor depth changes, this will influence the data, but is not presently captured in the metadata as it is self-evident in the data itself (the sensors measure water height). Sensor position changes that change the depth will only cause changes in the station code if the tidal constituents change as described above.

In either case, the location code does not change, as the "sensor code" has not changed.

