## RESPONSE FILES

The goal of the `resp` directory is to build a `golang` source module that allows easy instrument
response information discovery for downstream programs but primarily to build `stationxml` representations
of the meta information.

### ResponseType

The _ResponseType_ is a struct that mimics the _StationXML ResponseType_ element, but is not constrained to
any particular version. For this reason it can be used to load base XML files which can then be used to
construct the full response needed for _StationXML_.

The location of the base XML files can be found in:

- `auto` -- where the original response files have been converted into an XML equivalent.
- `files` -- where response files have been handcrafted as needed.
- `nrl` -- where responses derived from the NRL (Nominal Response Library) can be stored.

Files found in any of the locations can be used in the channel and component configurations by removing
the `.xml` suffix.

### References

* https://ds.iris.edu/ds/nrl/
 
- Mary E. Templeton (2017): IRIS Library of Nominal Response for Seismic Instruments. Incorporated Research Institutions for Seismology. Dataset. https://doi.org/10.17611/S7159Q
