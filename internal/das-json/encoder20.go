package dasjson

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	dasjson "github.com/GeoNet/delta/internal/das-json/dasv2.0"
)

type Encoder20 struct{}

func Encode20(r Root, schema string) ([]byte, error) {
	return Encoder20{}.MarshalRoot(r, schema)
}

// func (e Encoder20) toDateTime(at time.Time) dasjson.DateTime {
// 	if time.Since(at) > 0 {
// 		return dasjson.DateTime{Time: at}
// 	}
// 	return dasjson.DateTime{}
// }

// func (e Encoder20) toDateTimePtr(at time.Time) *dasjson.DateTime {
// 	t := dasjson.DateTime{Time: at}
// 	return &t
// }

func (e Encoder20) dateFromDateTime(t time.Time) string {
	return t.Format("2006-01-02")
}
func (e Encoder20) timeFromDateTime(t time.Time) string {
	return t.Format("2006-01-02T15:04:05")
}

func (e Encoder20) MarshalRoot(root Root, schema string) ([]byte, error) {
	var installs []dasjson.Root
	for _, d := range root.DAS {
		installs = append(installs, e.DASroot(root, d))
	}

	sort.Slice(installs, func(i, j int) bool {
		return installs[i].NetworkCode < installs[j].NetworkCode
	})
	var out []byte
	for _, install := range installs {
		d := install
		d.SchemaVersion = "2.0"
		b, err := json.MarshalIndent(d, "", "  ")

		// b, err := json.Marshal(r)
		// b, err := d.MarshalJSON()
		out = append(out, append(b, '\n')...)
		out = append([]byte("{\n  \"schema\": \""+schema+"\","), out[1:]...)
		// fmt.Print("r", r, "\n\n")
		// b, err := xml.MarshalIndent(r, "", "  ")
		if err != nil {
			return nil, err
		}

	}

	return append([]byte(""), append(out, '\n')...), nil
}

func (e Encoder20) DASroot(root Root, das DASoverview) dasjson.Root {

	var interrogators []*dasjson.Interrogator

	for _, i := range das.Interrogators {
		var acquisistions []*dasjson.Acquisition
		for _, a := range i.Acquisitions {
			var channel_groups []*dasjson.ChannelGroup
			for _, c := range a.ChannelGroups {

				channels := &dasjson.Channelarray{
					ChannelIds:              c.Channels.ChannelIDs,
					DepthsBelowSurface:      c.Channels.Dips,
					Dips:                    c.Channels.Dips,
					DistancesAlongFiber:     c.Channels.Dips,
					ElevationsAboveSeaLevel: c.Channels.Dips,
					Strikes:                 c.Channels.Dips,
					XCoordinates:            c.Channels.Dips,
					YCoordinates:            c.Channels.Dips,
				}
				channel_groups = append(channel_groups, &dasjson.ChannelGroup{
					CableId:                      c.CableID,
					ChannelGroupId:               c.ChannelGroupID,
					Channels:                     channels,
					CoordinateGenerationDate:     e.dateFromDateTime(c.CoordinateGenerationDate),
					CoordinateSystem:             c.CoordinateSystem,
					ReferenceFrame:               c.ReferenceFrame,
					LocationMethod:               c.LocationMethod,
					FiberId:                      c.FiberID,
					FirstUsableChannelId:         strconv.Itoa(c.FirstUsableChannelID),
					LastUsableChannelId:          strconv.Itoa(c.LastUsableChannelID),
					DistanceAlongFiberUnit:       "m",
					DepthBelowSurfaceUnit:        "m",
					UncertaintyInDepth:           c.UncertaintyInDepth,
					UncertaintyInDepthUnit:       "m",
					DipUnit:                      "degree",
					UncertaintyInDip:             c.UncertaintyInDip,
					UncertaintyInDipUnit:         "degree",
					ElevationAboveSeaLevelUnit:   "m",
					UncertaintyInElevation:       c.UncertaintyInElevation,
					UncertaintyInElevationUnit:   "m",
					StrikeUnit:                   "degree",
					UncertaintyInStrike:          c.UncertaintyInStrike,
					UncertaintyInStrikeUnit:      "degree",
					UncertaintyInXCoordinate:     c.UncertaintyInXCoordinate,
					UncertaintyInXCoordinateUnit: "m",
					UncertaintyInYCoordinate:     c.UncertaintyInYCoordinate,
					UncertaintyInYCoordinateUnit: "m",
					XCoordinateUnit:              "degree",
					YCoordinateUnit:              "degree",
				})
			}
			fmt.Print(a.ChannelGroups)
			nativeHeaders := &dasjson.NativeHeaders{
				AdditionalProperties: a.NativeHeaders,
			}

			acquisistions = append(acquisistions, &dasjson.Acquisition{
				AcquisitionId:               a.AcquisitionID,
				AcquisitionSampleRate:       a.AcquisitionSampleRate,
				AcquisitionSampleRateUnit:   "Hz",
				ChannelGroups:               channel_groups,
				Comment:                     a.Comment,
				GaugeLength:                 a.GaugeLength,
				GaugeLengthUnit:             "m",
				NativeHeaders:               nativeHeaders,
				NumberOfChannels:            a.NumberOfChannels,
				PulseRate:                   a.PulseRate,
				PulseRateUnit:               "Hz",
				PulseWidth:                  a.PulseWidth,
				PulseWidthUnit:              "m",
				ScaleFactor:                 1,
				SpatialSamplingInterval:     a.SpatialSamplingInterval,
				SpatialSamplingIntervalUnit: "m",
				UnitOfMeasure:               a.UnitOfMeasure,
				AcquisitionStartTime:        e.timeFromDateTime(a.AcquisitionStartTime),
				AcquisitionEndTime:          e.timeFromDateTime(a.AcquisitionEndTime),
			})
		}

		interrogators = append(interrogators, &dasjson.Interrogator{
			InterrogatorId:  i.InterrogatorID,
			Manufacturer:    i.Manufacturer,
			Model:           i.Model,
			SerialNumber:    i.SerialNumber,
			FirmwareVersion: i.FirmwareVersion,
			Comment:         i.Comment,
			Acquisitions:    acquisistions,
		})
	}
	var cables []*dasjson.Cable
	for _, c := range das.Cables {
		var fibers []*dasjson.Fiber
		for _, f := range c.Fibers {
			fibers = append(fibers, &dasjson.Fiber{
				FiberStartLocation:         f.FiberStartLocation,
				FiberStartLocationUnit:     "m",
				FiberEndLocation:           f.FiberEndLocation,
				FiberEndLocationUnit:       "m",
				FiberGeometry:              f.FiberGeometry,
				FiberId:                    f.FiberID,
				FiberMode:                  f.FiberMode,
				FiberOneWayAttenuation:     f.FiberOneWayAttenuation,
				FiberOneWayAttenuationUnit: "decibels/meter",
				FiberOpticLength:           f.FiberOpticLength,
				FiberOpticLengthUnit:       "m",
				FiberRefractionIndex:       f.FiberRefractionIndex,
				FiberWindingAngle:          f.FiberWindingAngle,
			})
		}

		cables = append(cables, &dasjson.Cable{
			CableId:                      c.CableID,
			CableOwner:                   c.CableOwner,
			CableModel:                   c.CableModel,
			CableCharacteristics:         c.CableCharacteristics,
			CableOutsideDiameter:         c.CableOutsideDiameter,
			CableOutsideDiameterUnit:     "mm",
			CableEnvironment:             c.CableEnvironment,
			CableInstallationEnvironment: c.CableInstallationEnvironment,
			CableBoundingBox:             c.CableBoundingBox,
			CableInstallationDate:        e.dateFromDateTime(c.CableInstallationDate),
			CableRemovalDate:             e.dateFromDateTime(c.CableRemovalDate),

			Comment: c.Comment,
			Fibers:  fibers,
		})
	}
	// fmt.Print("###", *das.Cables[0])
	return dasjson.Root{
		SchemaVersion: "2.0",
		NetworkCode:   das.Code,
		Country:       das.Country,
		Location:      das.Location,
		PrincipalInvestigator: []*dasjson.PrincipalInvestigatorFields{&dasjson.PrincipalInvestigatorFields{
			das.PrincipalInvestigator.Address,
			das.PrincipalInvestigator.Email,
			das.PrincipalInvestigator.Name}},
		PointOfContact:          das.PointOfContact,
		PointOfContactAddress:   das.PointOfContactAddress,
		PointOfContactEmail:     das.PointOfContactEmail,
		Cables:                  cables,
		Interrogators:           interrogators,
		FundingAgency:           das.FundingAgency,
		ProjectNumber:           das.ProjectNumber,
		DigitalObjectIdentifier: das.DigitalObjectIdentifier,
		PurposeOfDataCollection: das.PurposeOfDataCollection,
		// Comment: das.Comment,
		StartDate: e.dateFromDateTime(das.StartDate),
		EndDate:   e.dateFromDateTime(das.EndDate),
	}

}
