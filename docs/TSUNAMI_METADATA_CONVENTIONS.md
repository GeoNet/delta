# Tsunami Metadata Conventions

In this document we describe the metadata conventions used for GeoNet sensors that monitor tsunami.

## Overview

GeoNet has two types of tsunami monitoring stations: tide gauges and DART buoys. The distinction is made by the primary purpose of the station; whether it exists principally to record tsunami arrivals on-shore (via a tide gauge) or in the deep ocean (via a DART buoy).

The management of DART data and metadata is an evolving space, and until this is more established we refrain from documenting it here.

Tide gauge stations operated by GeoNet are not measuring absolute sea level, as they are not fully calibrated to the mean sea level. Data obtained by GeoNet tide gauges are measuring water height above the sensor.

Tide gauge data and metadata mostly follow the conventions of seismic data and metadata, and unless stated otherwise in this document these conventions are what apply to tide gauge data and metadata. The sibling document `SEISMIC_METADATA_CONVENTIONS.md` details these conventions. Here we make reference only to the differences between the tide gauge metadata conventions and those of seismic metadata.

## Location Code

For tide gauges, the location code indicates primary and backup sensors at the station. For example, a station may have sites at location codes `41` and `42`, and these location code are used to distinguish the respective data streams from the sensors at each site. In this way, we can think of the location code for tide gauges as a "sensor code". 

## Naming Conventions When Moving or Installating Different Sensors

Tide gauges are sensitive to the overlying water column, which is principally influenced by sensor depth and local tides. To correctly represent data measured at tsunami monitoring sites, tidal constituents must be derived.

Where in seismic data small changes to sensor position can have substantial impacts on the data, the same is not true for tide gauges. 

If a station position changes such that tidal constituents are still the same, the station code will not change.
If a station position changes significantly and new tidal constituents need to be derived, then a new station code will be created and the old station will be decommissioned. 

If a sensor position changes such its depth changes significantly, this will influence the data. Such change is not presently captured in the metadata and is represented in the data with an offset in the measured water level (the sensors measure water height). Sensor position changes that change the depth will only cause changes in the station code if the tidal constituents change as described above.

In either case, the location code does not change, as the "sensor code" has not changed.

