package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func seedFormat(t time.Time) string {
	if time.Now().After(t) {
		return fmt.Sprintf("%04d,%03d,%s", t.Year(), t.YearDay(), t.Format("15:04:05.0000"))
	}
	return ""
}

type Blockette struct {
	Type    int
	Content string
}

func (b Blockette) String() string {
	return fmt.Sprintf("%03d%04d%s", b.Type, len(b.Content)+7, b.Content)
}

type Blockettes []Blockette

func (b Blockettes) Encode() string {
	var lines []string
	for _, v := range b {
		lines = append(lines, v.String())
	}
	return strings.Join(lines, "\n")
}

type UnitsAbbreviation struct {
	Key         string
	Code        int
	Description string
}

func (u UnitsAbbreviation) String() string {
	return fmt.Sprintf("%03d%s~%s~", u.Code, u.Key, u.Description)
}

var unitsAbbreviation = []UnitsAbbreviation{
	{Key: "M/S", Code: 1, Description: "Velocity in Meters Per Second"},
	{Key: "A", Code: 2, Description: "Amperes"},
	{Key: "V", Code: 3, Description: "Volts"},
	{Key: "COUNT", Code: 4, Description: "Digital Counts"},
	{Key: "M/S**2", Code: 5, Description: "Acceleration in Meters Per Second Per Second"},
	{Key: "C", Code: 6, Description: "Degrees Centigrade"},
	{Key: "M", Code: 7, Description: "Displacement in Meters"},
	{Key: "PA", Code: 8, Description: "Pressure in Pascals"},
}

func lookupUnitsAbbreviation(key string) int {
	for _, k := range []string{key, strings.ToUpper(key)} {
		for _, a := range unitsAbbreviation {
			if a.Key == k {
				return a.Code
			}
		}
	}
	log.Println("UNABLE TO FIND: ", key)
	return 0
}

type GenericAbbreviation struct {
	Code        int
	Description string
}

type GenericAbbreviations []GenericAbbreviation

func (g GenericAbbreviations) Len() int           { return len(g) }
func (g GenericAbbreviations) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g GenericAbbreviations) Less(i, j int) bool { return g[i].Code < g[j].Code }

func (g GenericAbbreviation) String() string {
	return fmt.Sprintf("%03d%s~", g.Code, g.Description)
}

var genericAbbreviations = []GenericAbbreviation{
	{Code: 1, Description: "New Zealand National Seismograph Network"},
	{Code: 2, Description: "Sercel L4C"},
	{Code: 3, Description: "Guralp CMG-40T-30S"},
	{Code: 4, Description: "Kinemetrics FBA-23-DECK"},
	{Code: 5, Description: "Kinemetrics FBA-ES-T"},
	{Code: 6, Description: "Guralp CMG-40T-60S"},
	{Code: 7, Description: "Kinemetrics FBA-ES-T-DECK"},
	{Code: 8, Description: "Guralp CMG-3TB"},
	{Code: 9, Description: "Sercel L4C-3D"},
	{Code: 10, Description: "Guralp CMG-3ESP"},
	{Code: 11, Description: "Streckeisen STS-2"},
	{Code: 12, Description: "Guralp CMG-6T"},
	{Code: 13, Description: "Setra 270-600/12V"},
	{Code: 14, Description: "Lennartz LE-3Dlite"},
	{Code: 15, Description: "CSI CUSP3A SENSOR"},
	{Code: 16, Description: "Setra 270-600/24V"},
	{Code: 17, Description: "CSI CUSP3B SENSOR"},
	{Code: 18, Description: "Lennartz LE-3DliteMkII"},
	{Code: 19, Description: "GNS Science SDP"},
	{Code: 20, Description: "Geospace GS-11D seismometer"},
	{Code: 21, Description: "CSI CUSP3C SENSOR"},
	{Code: 22, Description: "Duke 2 Hz Duke Malin Seismometer"},
	{Code: 23, Description: "Guralp CMG-3TB-GN"},
	{Code: 24, Description: "Kinemetrics FBA-ES-T-BASALT"},
	{Code: 25, Description: "Guralp CMG-3ESPC"},
	{Code: 26, Description: "Kinemetrics FBA-ES-T-ISO"},
	{Code: 27, Description: "CSI CUSP3D SENSOR"},
	{Code: 28, Description: "Kinemetrics Kinemetrics SBEPI"},
	{Code: 29, Description: "IESE IESE S10g-4.5 (with preamp)"},
	{Code: 30, Description: "Guralp CMG-3ESP-Z"},
	{Code: 31, Description: "Nanometrics Nanometrics Trillium 120QA"},
	{Code: 32, Description: "TSP Scout Hydrophone"},
	{Code: 33, Description: "IESE IESE S31f-15 (with preamp)"},
	{Code: 34, Description: "Kinemetrics FBA-ES-T-OBSIDIAN"},
	{Code: 35, Description: "IESE IESE HS-1-LT Mini 4.5Hz"},
	{Code: 36, Description: "Boise State University InfraBSU microphone"},
	{Code: 37, Description: "GE Sensing Druck PTX-1830-LAND"},
	{Code: 38, Description: "GNS Science LM35"},
	{Code: 39, Description: "IESE IESE S10g-4.5 (without preamp)"},
	{Code: 40, Description: "Duke 4.5 Hz Duke Malin Seismometer"},
	{Code: 41, Description: "National tsunami gauge network"},
	{Code: 42, Description: "GE Sensing Druck PTX-1830"},
}

func lookupGenericAbbreviation(desc string) int {
	for _, a := range genericAbbreviations {
		if a.Description == desc {
			return a.Code
		}
	}
	abbr := len(genericAbbreviations)
	genericAbbreviations = append(genericAbbreviations, GenericAbbreviation{
		Code:        abbr,
		Description: desc,
	})
	return abbr
}

type CitedSourceDictionary struct {
	Code      int
	Name      string
	Publisher string
	Date      time.Time
}

func (c CitedSourceDictionary) String() string {
	return fmt.Sprintf("%02d%s~%s~%s~", c.Code, c.Name, c.Date.Format("2006"), c.Publisher)
}

type DataFormatDictionary struct {
	Key    string
	Code   int
	Family int
	Keys   []string
}

func (d DataFormatDictionary) String() string {
	if len(d.Keys) > 0 {
		return fmt.Sprintf("%s~%04d%03d%02d%s", d.Key, d.Code, d.Family, len(d.Keys), strings.Join(d.Keys, "~")+"~")
	} else {
		return fmt.Sprintf("%s~%04d%03d%02d", d.Key, d.Code, d.Family, len(d.Keys))
	}
}

var dataFormatDictionary = []DataFormatDictionary{
	{Key: "Steim-2 Integer Compression Format",
		Code: 1, Family: 0, Keys: []string{
			"F1 P4 W4 D C2 R1 P8 W4 D C2",
			"P0 W4 N15 S2,0,1",
			"T0 X W4",
			"T1 Y4 W1 D C2",
			"T2 W4 I1 D2",
			"K0 X D30",
			"K1 N0 D30 C2",
			"K2 Y2 D15 C2",
			"K3 Y3 D10 C2",
			"T3 W4 I1 D2",
			"K0 Y5 D6 C2",
			"K1 Y6 D5 C2",
			"K2 X D2 Y7 D4 C2",
			"K3 X D30",
		},
	},
	{Key: "Steim-1 Integer Compression Format",
		Code: 2, Family: 0, Keys: []string{
			"F1 P4 W4 D C2 R1 P8 W4 D C2",
			"P0 W4 N15 S2,0,1",
			"T0 X W4",
			"T1 Y4 W1 D C2",
			"T2 Y2 W2 D C2",
			"T3 N0 W4 D C2",
		},
	},
	{Key: "16-Bit Integer Format",
		Code: 3, Family: 0, Keys: []string{
			"M0",
			"W2 D0-15 C2",
		},
	},
	{Key: "32-Bit Integer Format",
		Code: 4, Family: 0, Keys: []string{
			"M0",
			"W4 D0-31 C2",
		},
	},
	{Key: "Console Log",
		Code: 5, Family: 0, Keys: []string{},
	},
}

func lookupDataFormatDictionary(key string) int {
	for _, c := range dataFormatDictionary {
		if c.Key == key {
			return c.Code
		}
	}
	return 0
}

type CommentDescription struct {
	Key         int
	Class       string
	Description string
	Units       int
}

func (c CommentDescription) String() string {
	return fmt.Sprintf("%04d%1s%s~%03d", c.Key, c.Class, c.Description, c.Units)
}

var commentDescriptions = []CommentDescription{
	{Key: 7, Class: "L", Description: "Location estimated from topographic map"},
	{Key: 2, Class: "L", Description: "Location estimation method is unknown"},
	{Key: 1, Class: "L", Description: "Location is given in WGS84"},
	{Key: 8, Class: "L", Description: "Location is given in NZGD49"},
	{Key: 4, Class: "L", Description: "Location is given in NZGD2000"},
	{Key: 6, Class: "L", Description: "Location estimated from internal GPS clock"},
	{Key: 9, Class: "L", Description: "Location estimated from plans and survey to mark"},
	{Key: 5, Class: "L", Description: "Location estimated from external GPS measurement"},
	{Key: 3, Class: "O", Description: "Sensor orientation not known"},
}

func lookupCommentDescription(desc string) int {
	for _, c := range commentDescriptions {
		if c.Description == desc {
			return c.Key
		}
	}
	return 0
}

//nolint:deadcode // for completeness
func findSurveyComment(survey string) int {
	switch survey {
	case "Unknown":
		return lookupCommentDescription("Location estimation method is unknown")
	case "External GPS Device":
		return lookupCommentDescription("Location estimated from external GPS measurement")
	case "Internal GPS Clock":
		return lookupCommentDescription("Location estimated from internal GPS clock")
	case "Topographic Map":
		return lookupCommentDescription("Location estimated from topographic map")
	case "Site Survey":
		return lookupCommentDescription("Location estimated from plans and survey to mark")
	}
	return 0
}

type StationIdentifier struct {
	Station     string
	Latitude    float64
	Longitude   float64
	Elevation   float64
	Name        string
	Description int
	Opened      time.Time
	Closed      time.Time
	Network     string
}

func (s StationIdentifier) String() string {
	return fmt.Sprintf("%-5s%+010.6f%+011.6f%+07.1f%04d%3s%s~%3d%4d%2d%s~%s~%1s%2s",
		s.Station,
		s.Latitude,
		s.Longitude,
		s.Elevation,
		int(0),
		"",
		s.Name,
		s.Description,
		int(3210),
		int(10),
		seedFormat(s.Opened),
		func() string {
			if time.Now().After(s.Closed) {
				return seedFormat(s.Closed)
			}
			return ""
		}(),
		"N",
		s.Network,
	)
}

type StationComment struct {
	Start  time.Time
	End    time.Time
	Lookup int
}

func (s StationComment) String() string {
	return fmt.Sprintf("%s~%s~%04d%06d",
		seedFormat(s.Start),
		func() string {
			if time.Now().After(s.End) {
				return seedFormat(s.End)
			}
			return ""
		}(),
		s.Lookup,
		0,
	)
}

type ChannelIdentifier struct {
	LocationIdentifier       string
	ChannelIdentifier        string
	SubchannelIdentifier     string
	InstrumentIdentifier     int
	OptionalComment          string
	UnitsOfSignalResponse    int
	UnitsOfCalibrationInput  int
	Latitude                 float64
	Longitude                float64
	Elevation                float64
	LocalDepth               float64
	Azimuth                  float64
	Dip                      float64
	DataFormatIdentifierCode int
	DataRecordLength         int
	SampleRate               float64
	MaxClockDrift            float64
	NumberOfComments         string
	ChannelFlags             string
	StartDate                time.Time
	EndDate                  time.Time
	UpdateFlag               string
}

func (c ChannelIdentifier) String() string {
	return fmt.Sprintf(
		"%-2s%-3s%4s%03d%s~%03d%03d%+010.6f%+011.6f%+07.1f%05.1f%05.1f%+05.1f%04d%02d%10.4E%10.4E%4s%s~%s~%s~%1s",
		c.LocationIdentifier,
		c.ChannelIdentifier,
		c.SubchannelIdentifier,
		c.InstrumentIdentifier,
		c.OptionalComment,
		c.UnitsOfSignalResponse,
		c.UnitsOfCalibrationInput,
		c.Latitude,
		c.Longitude,
		c.Elevation,
		c.LocalDepth,
		c.Azimuth,
		c.Dip,
		c.DataFormatIdentifierCode,
		c.DataRecordLength,
		c.SampleRate,
		c.MaxClockDrift,
		c.NumberOfComments,
		c.ChannelFlags,
		seedFormat(c.StartDate),
		seedFormat(c.EndDate),
		c.UpdateFlag,
	)
}

type ChannelComment struct {
	BeginningEffectiveTime time.Time
	EndEffectiveTime       time.Time
	CommentCodeKey         int
	CommentLevel           int
}

func (c ChannelComment) String() string {
	return fmt.Sprintf("%s~%s~%04d%06d",
		seedFormat(c.BeginningEffectiveTime),
		seedFormat(c.EndEffectiveTime),
		c.CommentCodeKey,
		c.CommentLevel,
	)
}

type ResponsePoleZero struct {
	Real           float64
	Imaginary      float64
	RealError      float64
	ImaginaryError float64
}

func (r ResponsePoleZero) String() string {
	return fmt.Sprintf("%+12.5E%+12.5E%+12.5E%+12.5E",
		r.Real,
		r.Imaginary,
		r.RealError,
		r.ImaginaryError,
	)
}

type ResponsePolesZeros struct {
	TransferFunctionType   string
	StageSequenceNumber    int
	StageSignalInputUnits  int
	StageSignalOutputUnits int
	AONormalizationFactor  float64
	NormalizationFrequency float64
	Zeros                  []ResponsePoleZero
	Poles                  []ResponsePoleZero
}

func (r ResponsePolesZeros) String() string {
	return fmt.Sprintf("%1s%02d%03d%03d%+12.5E%+12.5E%03d%s%03d%s",
		r.TransferFunctionType,
		r.StageSequenceNumber,
		r.StageSignalInputUnits,
		r.StageSignalOutputUnits,
		r.AONormalizationFactor,
		r.NormalizationFrequency,
		len(r.Zeros),
		func() string {
			var pz []string
			for _, z := range r.Zeros {
				pz = append(pz, z.String())
			}
			return strings.Join(pz, "")
		}(),
		len(r.Poles),
		func() string {
			var pz []string
			for _, p := range r.Poles {
				pz = append(pz, p.String())
			}
			return strings.Join(pz, "")
		}(),
	)
}

type ResponseCoefficient struct {
	Coefficient      float64
	CoefficientError float64
}

func (r ResponseCoefficient) String() string {
	return fmt.Sprintf("%+12.5E%+12.5E",
		r.Coefficient,
		r.CoefficientError,
	)
}

type ResponseCoefficients struct {
	ResponseType           string
	StageSequenceNumber    int
	StageSignalInputUnits  int
	StageSignalOutputUnits int
	Numerators             []ResponseCoefficient
	Denominators           []ResponseCoefficient
}

func (r ResponseCoefficients) String() string {
	return fmt.Sprintf("%1s%02d%03d%03d%04d%s%04d%s",
		r.ResponseType,
		r.StageSequenceNumber,
		r.StageSignalInputUnits,
		r.StageSignalOutputUnits,
		len(r.Numerators),
		func() string {
			var c []string
			for _, n := range r.Numerators {
				c = append(c, n.String())
			}
			return strings.Join(c, "")
		}(),
		len(r.Denominators),
		func() string {
			var c []string
			for _, d := range r.Denominators {
				c = append(c, d.String())
			}
			return strings.Join(c, "")
		}(),
	)

}

type FIRResponse struct {
	StageSequenceNumber    int
	ResponseName           string
	SymmetryCode           string
	StageSignalInputUnits  int
	StageSignalOutputUnits int
	Coefficients           []float64
}

func (r FIRResponse) String() string {
	return fmt.Sprintf("%02d%s~%1s%03d%03d%04d%s",
		r.StageSequenceNumber,
		r.ResponseName,
		r.SymmetryCode,
		r.StageSignalInputUnits,
		r.StageSignalOutputUnits,
		len(r.Coefficients),
		func() string {
			var c []string
			for i := range r.Coefficients {
				if r.SymmetryCode == "A" {
					//c = append(c, fmt.Sprintf("%+14.7E", r.Coefficients[len(r.Coefficients)-i-1]))
					c = append(c, fmt.Sprintf("%+14.7E", r.Coefficients[i]))
				} else {
					c = append(c, fmt.Sprintf("%+14.7E", r.Coefficients[i]))
				}
			}
			return strings.Join(c, "")
		}(),
	)
}

type ResponsePolynomial struct {
	TransferFunctionType        string
	StageSequenceNumber         int
	StageSignalInputUnits       int
	StageSignalOutputUnits      int
	PolynomialApproximationType string
	ValidFrequencyUnits         string
	LowerValidFrequencyBound    string
	UpperValidFrequencyBound    string
	LowerBoundOfApproximation   float64
	UpperBoundOfApproximation   float64
	MaximumAbsoluteError        float64
	Coefficients                []float64
}

func (r ResponsePolynomial) String() string {
	return fmt.Sprintf("%1s%02d%03d%03d%1s%1s%12s%12s%12s%12s%+12.5E%03d%s",
		r.TransferFunctionType,
		r.StageSequenceNumber,
		r.StageSignalInputUnits,
		r.StageSignalOutputUnits,
		r.PolynomialApproximationType,
		r.ValidFrequencyUnits,
		r.LowerValidFrequencyBound,
		r.UpperValidFrequencyBound,
		func() string {
			if r.LowerBoundOfApproximation != 0.0 {
				return fmt.Sprintf("%+12.5E", r.LowerBoundOfApproximation)
			}
			return ""
		}(),
		func() string {
			if r.UpperBoundOfApproximation != 0.0 {
				return fmt.Sprintf("%+12.5E", r.UpperBoundOfApproximation)
			}
			return ""
		}(),
		r.MaximumAbsoluteError,
		len(r.Coefficients),
		func() string {
			var c []string
			for _, n := range r.Coefficients {
				c = append(c, fmt.Sprintf("%+12.5E%+12.5E", n, 0.0))
			}
			return strings.Join(c, "")
		}(),
	)
}

type Decimation struct {
	StageSequenceNumber int
	InputSampleRate     float64
	DecimationFactor    int
	DecimationOffset    int
	EstimatedDelay      float64
	CorrectionApplied   float64
}

func (d Decimation) String() string {
	return fmt.Sprintf("%02d%10.4E%05d%05d%+11.4E%+11.4E",
		d.StageSequenceNumber,
		d.InputSampleRate,
		d.DecimationFactor,
		d.DecimationOffset,
		d.EstimatedDelay,
		d.CorrectionApplied,
	)
}

type StageGain struct {
	StageSequenceNumber int
	Gain                float64
	Frequency           float64
	Something           int
	SomethingElse       []interface{}
}

func (s StageGain) String() string {
	return fmt.Sprintf("%02d%+12.5E%+12.5E%02d",
		s.StageSequenceNumber,
		s.Gain,
		s.Frequency,
		s.Something,
	)
}
