# Tsunami Monitoring Sensors Metadata Conventions

In this document we describe the metadata conventions used for GeoNet sensors that monitor tsunami.

## Overview

GeoNet has two types of tsunami monitoring stations: tide gauges and DART buoys. The distinction is made by the primary purpose of the station; whether it exists principally to record tsunami arrivals on-shore (via a tide gauge) or in the deep ocean (via a DART buoy).

The management of DART data and metadata is an evolving space, and until this is more established we refrain from documenting it here.

Tide gauge stations operated by GeoNet are not measuring absolute sea level, as they are not fully calibrated to the mean sea level. Data obtained by GeoNet tide gauges are measuring water height above the sensor.

Tide gauge data and metadata mostly follow the conventions of seismic data and metadata, and unless stated otherwise in this document these conventions are what apply to tide gauge data and metadata. The sibling document `MINISEED_METADATA_CONVENTIONS.md` details these conventions. Here we make reference only to the differences between the tide gauge metadata conventions and those of seismic metadata.

## Location Code

For tide gauges, the location code indicates primary and backup sensors at the station. A station may have sites at location codes:
- `40`: primary sensor
- `41`: secondary sensor

For tsunami monitoring gauges, these location codes are used to distinguish the respective data streams from the two sensors installed at each site.

## Naming Conventions When Moving or Installating Different Sensors

Tide gauges are sensitive to the overlying water column, which is principally influenced by sensor depth and local tides. To correctly represent data measured at tsunami monitoring sites, tidal constituents must be derived.

Where in seismic data small changes to sensor position can have substantial impacts on the data, the same is not true for tide gauges.

If a sensor position changes such that tidal constituents are still the same, the sensor's station code will not change.

If a sensor position changes and new tidal constituents need to be derived, then a new station code will be created. The old station will be decommissioned if no sensors remain there.

If a sensor is moved and its depth changes, this will influence the data, and will present an offset in the measured water level (the sensors measure water height). However, because the tsunami monitoring changes measures relative water level, such a change is not captured in the metadata when tidal constituents remain unchanged as a result of the change; sensor depth changes will only cause changes in the station code if the tidal constituents change.
