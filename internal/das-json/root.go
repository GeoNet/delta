package dasjson

import (
	"time"

	"github.com/GeoNet/delta/meta"
	// dasjson "github.com/GeoNet/delta/internal/das-json/dasv2.0"
	//this needs to be more flexible somehow, hide this in the encoder section?
)

// // Channel forms the main part of a set of StationXML Channel elements
// type Channel struct {
// 	LocationCode string

// 	Latitude  float64
// 	Longitude float64
// 	Elevation float64
// }

type DASoverview struct {
	Code                    string
	Location                string
	Country                 string
	PointOfContact          string
	PointOfContactEmail     string
	PointOfContactAddress   string
	FundingAgency           string
	ProjectNumber           string
	DigitalObjectIdentifier string
	PurposeOfDataCollection string
	StartDate               time.Time
	EndDate                 time.Time

	PrincipalInvestigator PrincipalInvestigator
	Cables                []*Cable
	Interrogators         []*Interrogator
	// Comment: das.Comment,

}

type PrincipalInvestigator struct {
	// Physical address and institution of principal investigator.
	Address string
	// Email address of principal investigator.
	Email string
	// Name of principal investigator for the installation.
	Name string
}

type Cable struct {
	// List of approximate coordinates of the bounding box containing the array [minimum latitude, maximum latitude, minimum longitude, and maximum longitude].
	CableBoundingBox []float64
	// Describes how the fiber is encased in the cable.
	CableCharacteristics string
	// Describes the installation environment of cable [conduit| trench | outside borehole casing | wireline | other ].
	CableEnvironment string
	// Unique identifier for the specific physical cable used in the data acquisition. Identifier should have the maximum of eight alphanumeric characters with no special characters (e.g., underscores, period, and dash).
	CableID string
	// Cable installation date in Coordinated Universal Time (UTC). This may not be known, an optional field.
	CableInstallationDate time.Time
	// Plain language description of how the cable was installed.
	CableInstallationEnvironment string
	// Cable manufacturer model.
	CableModel string
	// Cable outside diameter.
	CableOutsideDiameter float64
	// Unit of outside diameter of cable. SI units should be used.
	CableOutsideDiameterUnit string
	// Proprietary owner of the cable.
	CableOwner string
	// Cable removal date in UTC. This may not be known, an optional field.
	CableRemovalDate time.Time
	// Additional comments.
	Comment string
	// One or more Fibers for this cable
	Fibers []*Fiber
}

type Fiber struct {
	Comment                    string
	FiberEndLocation           float64
	FiberEndLocationUnit       string
	FiberGeometry              string
	FiberID                    string
	FiberMode                  string
	FiberOneWayAttenuation     float64
	FiberOneWayAttenuationUnit string
	FiberOpticLength           float64
	FiberOpticLengthUnit       string
	FiberRefractionIndex       float64
	FiberStartLocation         float64
	FiberStartLocationUnit     string
	FiberWindingAngle          float64
	FiberWindingAngleUnit      string
}

type Interrogator struct {
	// One or more acquisitions with this interrogator.
	Acquisitions []*Acquisition
	// Additional comments.
	Comment string
	// Firmware version of the software used within the interrogator.
	FirmwareVersion string
	// Unique identifier of the interrogator unit used in the experiment, assigned by data provider. Identifier should have the maximum of eight alphanumeric characters with no special characters (e.g., underscores, period, and dash).
	InterrogatorID string
	// Manufacturer name of the interrogator.
	Manufacturer string
	// Model number of the interrogator.
	Model string
	// Serial number of the interrogator.
	SerialNumber string
	// Location to match against DAS installation list
	Location string
}

type Acquisition struct {
	// End time of this data acquisition in UTC. If data acquisition is still in operation, this field should be omitted and updated when it is known.
	AcquisitionEndTime time.Time
	// Unique identifier of the data acquisition, assigned by data provider. Identifier should have the maximum of eight alphanumeric characters with no special characters (e.g., underscores, period, and dash). One identifier per acquisition settings.
	AcquisitionID string
	// The rate at which the interrogator provides output data.
	AcquisitionSampleRate float64
	// Unit of acquisition sample rate. SI units should be used.
	AcquisitionSampleRateUnit string
	// Start time of this data acquisition in Coordinated Universal Time (UTC).
	AcquisitionStartTime time.Time
	// One or more channel groups for this acquisition.
	ChannelGroups []*ChannelGroup
	// Additional comments.
	Comment string
	// The averaging length along the fiber for a measurement, determined at experiment setup and used during acquisition.
	GaugeLength float64
	// Unit of gauge length. SI units should be used.
	GaugeLengthUnit string
	// Native headers from the interrogator for this acquisition as a key-value list.
	//   NativeHeaders *NativeHeaders
	NativeHeaders map[string]interface{}
	// The total number of sampling points along the fiber as output from the interrogator, referred to as NumberOfLoci in PRODML.
	NumberOfChannels int
	// Rate at which the interrogator unit interrogates the fiber sensor.
	PulseRate float64
	// Unit of pulse rate. SI units should be used.
	PulseRateUnit string
	// Width of the pulse sent down the fiber in unit of time.
	PulseWidth float64
	// Unit of pulse width. SI units should be used.
	PulseWidthUnit string
	// Scaling factor to multiply in order to express data in physical units.
	ScaleFactor float64
	// The channel spacing, or offset, between channels.
	SpatialSamplingInterval float64
	// Unit of spatial sampling interval. SI units should be used.
	SpatialSamplingIntervalUnit string
	// Unit of measure of archived data set. This may be the same unit as the Interrogator unit measure if the data are raw. SI units should be used.
	UnitOfMeasure string
}

type ChannelGroup struct {
	// Unique identifier of cable linked to this channel group.
	CableID string
	// Unique identifier of the channel group. Identifier should have the maximum of eight alphanumeric characters with no special characters (e.g., underscores, period, and dash).
	ChannelGroupID string
	// Channels for this channel group.
	Channels *Channels
	// Additional comments.
	Comment string
	// Date that the channel coordinates were generated.
	CoordinateGenerationDate time.Time
	// Coordinate system used in channel location. Select geographic for coordinates in latitude and longitude. Universal Transverse Mercator (UTM) is the preferred projected coordinate system.
	CoordinateSystem string
	// Unique identifier of fiber linked to this channel group. The fiber_id must nest within specified cable_id.
	FiberID string
	// The first channel considered physically meaningful. Typically, this will be a few channels after the last channel inside the interrogator. Must match one of the channel IDs.
	FirstUsableChannelID int
	// The last channel considered physically meaningful. Typically, an acquisition will add virtual channels at the end of the fiber or the optical path might not reach the end of the fiber, and this can be used to cut the data set into a working set. Must match one of the channel IDs.
	LastUsableChannelID int
	// Method used to determine the coordinates of each channel. Example: tap test, GPS tracker.
	LocationMethod string `json:"location_method,opt,omitempty"`
	// Geodetic datum and/or grid zone number in UTM.
	ReferenceFrame string `json:"reference_frame"`
	// Uncertainty in depth below surface, representing the average uncertainty in depth for all the coordinates.
	UncertaintyInDepth float64 `json:"uncertainty_in_depth,opt,omitempty"`
	// Uncertainty in the dip.
	UncertaintyInDip float64 `json:"uncertainty_in_dip,opt,omitempty"`
	// Uncertainty in elevation above sea level, representing the average uncertainty in elevation for all the coordinates.
	UncertaintyInElevation float64 `json:"uncertainty_in_elevation,opt,omitempty"`
	// Uncertainty in the strike.
	UncertaintyInStrike float64 `json:"uncertainty_in_strike,opt,omitempty"`
	// Uncertainty in the x-coordinate (UTM-x or longitude) representing the average uncertainty for all the channels.
	UncertaintyInXCoordinate float64 `json:"uncertainty_in_x_coordinate,opt,omitempty"`
	// Uncertainty in the y-coordinate (UTM-y or latitude) representing the average uncertainty for all the channels.
	UncertaintyInYCoordinate float64 `json:"uncertainty_in_y_coordinate,opt,omitempty"`
}

type Channels struct {
	// Array of Unique IDs of each individual channel, with the maximum of eight alphanumeric characters with no special characters (e.g., underscores, period, and dash).
	ChannelIDs []string
	// Array of vertical distances from the local ground surface to the fibre at each channel, with the ground surface defined as zero and positive values measured downward.
	DepthsBelowSurface []float64
	// Array of individual channel dips, measured from horizontal and positive downward: 0° is horizontal, −90° is vertically upward, and +90° is vertically downward.
	Dips []float64
	// Array of distances along fiber between the interrogator and channel. Zero meter at the connector to the interrogator; and negative distances for fiber section inside the interrogator.
	DistancesAlongFiber []float64
	// Array of local ground-surface elevations or heights above sea level at each channel.
	ElevationsAboveSeaLevel []float64
	// Array of local azimuths at each channel, measured clockwise from geographical north in the direction of increasing channel number.
	Strikes []float64
	// Array of X-coordinates (UTM-x or longitude) of the channels.
	XCoordinates []float64
	// Array of y-coordinates (UTM-y or latitude) of the channels.
	YCoordinates []float64
}

func (channels *Channels) AppendChannel(channel meta.DASChannel) {
	channels.ChannelIDs = append(channels.ChannelIDs, channel.ChannelID)
	channels.DepthsBelowSurface = append(channels.DepthsBelowSurface, channel.DepthBelowSurface)
	channels.Dips = append(channels.Dips, channel.Dip)
	channels.DistancesAlongFiber = append(channels.DistancesAlongFiber, channel.DistanceAlongFiber)
	channels.ElevationsAboveSeaLevel = append(channels.ElevationsAboveSeaLevel, channel.ElevationAboveSeaLevel)
	channels.Strikes = append(channels.Strikes, channel.Strike)
	channels.XCoordinates = append(channels.XCoordinates, channel.XCoordinate)
	channels.YCoordinates = append(channels.YCoordinates, channel.YCoordinate)
}

// Root describes the standard StationXML layout which can be used as the barebones for building version specific encoders.
type Root struct {
	DAS []DASoverview
}

// NetworkCode returns the network code of the first Network entry in the Root structure, this is aimed at building single entry file names.
func (r Root) NetworkCode() string {
	for _, d := range r.DAS {
		return d.Code
	}
	return ""
}

// // Single builds a Overview structure for the given station code.
// func (r Root) Single(code string) (Root, bool) {
// 	for _, e := range r.Externals {
// 		for _, n := range e.Networks {
// 			for _, s := range n.Stations {
// 				if s.Code != code {
// 					continue
// 				}

// 				root := Root{
// 					Source: r.Source,
// 					Sender: r.Sender,
// 					Module: r.Module,
// 					Create: r.Create,

// 					Externals: []External{{
// 						Code:        e.Code,
// 						Description: e.Description,
// 						StartDate:   e.StartDate,
// 						EndDate:     e.EndDate,
// 						Networks: []Network{{
// 							Code:        n.Code,
// 							Description: n.Description,
// 							Stations:    []Station{s},
// 						}},
// 					}},
// 				}

// 				return root, true
// 			}
// 		}
// 	}

// 	return Root{}, false
// }
