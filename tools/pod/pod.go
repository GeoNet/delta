package main

import (
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ozym/fdsn/stationxml"
)

type Channels []*stationxml.Channel

func (c Channels) Len() int           { return len(c) }
func (c Channels) Less(i, j int) bool { return c[i].StartDate.Time.Before(c[j].StartDate.Time) }
func (c Channels) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

type Pod struct {
	base string
}

func NewPod(base string) *Pod {
	return &Pod{
		base: base,
	}
}

func (p *Pod) Header() error {

	sort.Sort(GenericAbbreviations(genericAbbreviations))

	var lines []string
	for _, b := range dataFormatDictionary {
		lines = append(lines, Blockette{
			Type:    30,
			Content: b.String(),
		}.String())
	}
	for _, b := range commentDescriptions {
		lines = append(lines, Blockette{
			Type:    31,
			Content: b.String(),
		}.String())
	}
	lines = append(lines, Blockette{
		Type: 32,
		Content: CitedSourceDictionary{
			Code:      1,
			Name:      "Unknown",
			Publisher: "Unknown Publisher",
			Date:      time.Now(),
		}.String(),
	}.String())
	for _, b := range genericAbbreviations {
		lines = append(lines, Blockette{
			Type:    33,
			Content: b.String(),
		}.String())
	}
	for _, b := range unitsAbbreviation {
		lines = append(lines, Blockette{
			Type:    34,
			Content: b.String(),
		}.String())
	}

	header := filepath.Join(p.base, "HDR000", "H.A")
	if err := os.MkdirAll(filepath.Dir(header), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(header, []byte(strings.Join(lines, "\n")+"\n"), 0644); err != nil {
		return err
	}

	return nil
}

func (p *Pod) Network(net *stationxml.Network) error {
	for i, _ := range net.Stations {
		if err := p.Station(net, &net.Stations[i]); err != nil {
			return err
		}
	}
	return nil
}

func (p *Pod) Station(net *stationxml.Network, sta *stationxml.Station) error {

	b50 := filepath.Join(p.base, "HDR000", strings.Join([]string{sta.Code, net.Code}, "."), "B050")
	if err := os.MkdirAll(filepath.Dir(b50), 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(b50, []byte(Blockette{
		Type: 50,
		Content: StationIdentifier{
			Station:   sta.Code,
			Latitude:  sta.Latitude.Value,
			Longitude: sta.Longitude.Value,
			Elevation: sta.Elevation.Value,
			Name:      sta.Site.Name,
			Description: func() int {
				switch lookupGenericAbbreviation(sta.Description) {
				case 0:
					return lookupGenericAbbreviation("New Zealand National Seismograph Network")
				default:
					return lookupGenericAbbreviation(sta.Description)
				}
			}(),
			Opened: sta.CreationDate.Time,
			Closed: func() time.Time {
				if sta.TerminationDate != nil {
					return sta.TerminationDate.Time
				}
				return time.Now().AddDate(9999, 0, 0)
			}(),
			Network: net.Code,
		}.String(),
	}.String()+"\n"), 0644); err != nil {
		return err
	}

	comments := []Blockette{}

	for _, b := range sta.Comments {
		comments = append(comments,
			Blockette{
				Type: 51,
				Content: StationComment{
					Start: sta.CreationDate.Time,
					End: func() time.Time {
						if sta.TerminationDate != nil {
							return sta.TerminationDate.Time
						}
						return time.Now().AddDate(9999, 0, 0)
					}(),
					Lookup: lookupCommentDescription(b.Value),
				}.String(),
			})
	}

	b51 := filepath.Join(p.base, "HDR000", strings.Join([]string{sta.Code, net.Code}, "."), "B051")
	if err := os.MkdirAll(filepath.Dir(b51), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(b51, []byte(Blockettes(comments).Encode()+"\n"), 0644); err != nil {
		return err
	}

	var channels []*stationxml.Channel
	for i, _ := range sta.Channels {
		channels = append(channels, &sta.Channels[i])
	}

	sort.Sort(Channels(channels))

	blockettes := make(map[string][]Blockette)
	for _, c := range channels {
		code := strings.Join([]string{c.Code, c.LocationCode}, ".")
		blockettes[code] = append(blockettes[code], p.Channel(c)...)
	}

	for k, v := range blockettes {
		b52 := filepath.Join(p.base, "HDR000", strings.Join([]string{sta.Code, net.Code}, "."), k, "B052")
		if err := os.MkdirAll(filepath.Dir(b52), 0755); err != nil {
			return err
		}
		if err := ioutil.WriteFile(b52, []byte(Blockettes(v).Encode()+"\n"), 0644); err != nil {
			return err
		}
	}

	blockettes = make(map[string][]Blockette)
	for _, c := range channels {
		code := strings.Join([]string{c.Code, c.LocationCode}, ".")
		blockettes[code] = append(blockettes[code], p.ChannelComments(c)...)
	}

	for k, v := range blockettes {
		b59 := filepath.Join(p.base, "HDR000", strings.Join([]string{sta.Code, net.Code}, "."), k, "B059")
		if err := os.MkdirAll(filepath.Dir(b59), 0755); err != nil {
			return err
		}
		if err := ioutil.WriteFile(b59, []byte(Blockettes(v).Encode()+"\n"), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (p *Pod) Channel(cha *stationxml.Channel) []Blockette {
	var blockettes []Blockette

	blockettes = append(blockettes, Blockette{
		Type: 52,
		Content: ChannelIdentifier{
			LocationIdentifier:   cha.LocationCode,
			ChannelIdentifier:    cha.Code,
			SubchannelIdentifier: "",
			InstrumentIdentifier: func() int {
				if cha.Sensor != nil {
					var model string
					switch cha.Sensor.Model {
					case "InfraBSU microphone":
						model = "Boise State University InfraBSU microphone"
					case "2 Hz Duke Malin Seismometer":
						model = "Duke 2 Hz Duke Malin Seismometer"
					case "FBA-23-DECK":
						model = "Kinemetrics FBA-23-DECK"
					case "CUSP3A":
						model = "CSI CUSP3A SENSOR"
					case "CUSP3B":
						model = "CSI CUSP3B SENSOR"
					case "CUSP3C":
						model = "CSI CUSP3C SENSOR"
					case "CUSP3D":
						model = "CSI CUSP3D SENSOR"
					case "270-600/24V":
						model = "Setra 270-600/24V"
					case "270-600/12V":
						model = "Setra 270-600/12V"
					case "CMG-3ESP":
						model = "Guralp CMG-3ESP"
					case "CMG-3ESPC":
						model = "Guralp CMG-3ESPC"
					case "CMG-3ESP-Z":
						model = "Guralp CMG-3ESP-Z"
					case "CMG-3TB":
						model = "Guralp CMG-3TB"
					case "CMG-3TB-GN":
						model = "Guralp CMG-3TB-GN"
					case "CMG-40T-30S":
						model = "Guralp CMG-40T-30S"
					case "CMG-40T-60S":
						model = "Guralp CMG-40T-60S"
					case "FBA-ES-T":
						model = "Kinemetrics FBA-ES-T"
					case "FBA-ES-T-ISO":
						model = "Kinemetrics FBA-ES-T-ISO"
					case "FBA-ES-T-DECK":
						model = "Kinemetrics FBA-ES-T-DECK"
					case "FBA-ES-T-BASALT":
						model = "Kinemetrics FBA-ES-T-BASALT"
					case "FBA-ES-T-OBSIDIAN":
						model = "Kinemetrics FBA-ES-T-OBSIDIAN"
					case "L4C":
						model = "Sercel L4C"
					case "L4C-3D":
						model = "Sercel L4C-3D"
					case "LE-3Dlite":
						model = "Lennartz LE-3Dlite"
					case "LE-3DliteMkII":
						model = "Lennartz LE-3DliteMkII"
					case "STS-2":
						model = "Streckeisen STS-2"
					case "Nanometrics Trillium 120QA":
						model = "Nanometrics Nanometrics Trillium 120QA"
					default:
						log.Println("missing sensor model lookup ->>>", cha.Sensor.Model, "<<<--")
						model = cha.Sensor.Model
					}
					if abbr := lookupGenericAbbreviation(model); abbr > 0 {
						return abbr
					}
				}
				return 0
			}(),
			OptionalComment: func() string {
				if cha.Sensor != nil {
					return "S/N " + cha.Sensor.SerialNumber
				}
				return ""
			}(),
			UnitsOfSignalResponse: func() int {
				if len(cha.Response.Stages) > 0 {
					stage := cha.Response.Stages[0]
					switch {
					case stage.PolesZeros != nil:
						return lookupUnitsAbbreviation(stage.PolesZeros.InputUnits.Name)
					case stage.Coefficients != nil:
						return lookupUnitsAbbreviation(stage.Coefficients.InputUnits.Name)
					case stage.ResponseList != nil:
						return lookupUnitsAbbreviation(stage.ResponseList.InputUnits.Name)
					case stage.FIR != nil:
						return lookupUnitsAbbreviation(stage.FIR.InputUnits.Name)
					case stage.Polynomial != nil:
						if stage.Polynomial.InputUnits.Name == "hPa" {
							return lookupUnitsAbbreviation("PA")
							//return lookupUnitsAbbreviation("M")
						}
						return lookupUnitsAbbreviation(stage.Polynomial.InputUnits.Name)
					}
				}
				return 0
			}(),
			UnitsOfCalibrationInput: func() int {
				if cha.CalibrationUnits != nil {
					return lookupUnitsAbbreviation(cha.CalibrationUnits.Name)
				}
				return lookupUnitsAbbreviation("A")
			}(),
			Latitude:   cha.Latitude.Value,
			Longitude:  cha.Longitude.Value,
			Elevation:  cha.Elevation.Value,
			LocalDepth: cha.Depth.Value,
			Azimuth: func() float64 {
				if cha.Azimuth != nil {
					return cha.Azimuth.Value
				}
				return 0.0
			}(),
			Dip: func() float64 {
				if cha.Dip != nil {
					return cha.Dip.Value
				}
				return 0.0
			}(),
			DataFormatIdentifierCode: func() int {
				switch cha.StorageFormat {
				case "Steim2":
					return lookupDataFormatDictionary("Steim-2 Integer Compression Format")
				default:
					return lookupDataFormatDictionary(cha.StorageFormat)
				}
			}(),
			DataRecordLength: int(math.Log(4096.0) / math.Log(2.0)),
			SampleRate:       cha.SampleRateGroup.SampleRate.Value,
			MaxClockDrift: func() float64 {
				if cha.ClockDrift != nil {
					return cha.ClockDrift.Value
				}
				return 0.0
			}(),
			NumberOfComments: "",
			ChannelFlags: func() string {
				var flags string
				for _, t := range cha.Types {
					switch t {
					case stationxml.TypeContinuous:
						flags = flags + "C"
					case stationxml.TypeTriggered:
						flags = flags + "T"
					case stationxml.TypeWeather:
						flags = flags + "W"
					case stationxml.TypeGeophysical:
						flags = flags + "G"
					}
				}
				return flags
			}(),
			StartDate: func() time.Time {
				if cha.StartDate != nil {
					return cha.StartDate.Time
				}
				return time.Time{}
			}(),
			EndDate: func() time.Time {
				if cha.EndDate != nil {
					return cha.EndDate.Time
				}
				return time.Time{}
			}(),
			UpdateFlag: "N",
		}.String(),
	})

	if cha.Response != nil {
		for _, s := range cha.Response.Stages {
			if s.PolesZeros != nil {
				blockettes = append(blockettes, Blockette{
					Type: 53,
					Content: ResponsePolesZeros{
						TransferFunctionType: func() string {
							switch s.PolesZeros.PzTransferFunctionType {
							case stationxml.PZFunctionLaplaceRadiansPerSecond:
								return "A"
							case stationxml.PZFunctionLaplaceHertz:
								return "B"
							case stationxml.PZFunctionLaplaceZTransform:
								return "D"
							default:
								return " "
							}
						}(),
						StageSequenceNumber:    int(s.Number),
						StageSignalInputUnits:  lookupUnitsAbbreviation(s.PolesZeros.InputUnits.Name),
						StageSignalOutputUnits: lookupUnitsAbbreviation(s.PolesZeros.OutputUnits.Name),
						AONormalizationFactor:  s.PolesZeros.NormalizationFactor,
						NormalizationFrequency: s.PolesZeros.NormalizationFrequency.Value,
						Zeros: func() []ResponsePoleZero {
							var pz []ResponsePoleZero
							for _, z := range s.PolesZeros.Zeros {
								pz = append(pz, ResponsePoleZero{
									Real:      z.Real.Value,
									Imaginary: z.Imaginary.Value,
								})
							}
							return pz
						}(),
						Poles: func() []ResponsePoleZero {
							var pz []ResponsePoleZero
							for _, p := range s.PolesZeros.Poles {
								pz = append(pz, ResponsePoleZero{
									Real:      p.Real.Value,
									Imaginary: p.Imaginary.Value,
								})
							}
							return pz
						}(),
					}.String(),
				})
			}
			if s.Coefficients != nil {
				blockettes = append(blockettes, Blockette{
					Type: 54,
					Content: ResponseCoefficients{
						ResponseType: func() string {
							switch s.Coefficients.CfTransferFunctionType {
							case stationxml.CfFunctionAnalogRadiansPerSecond:
								return "A"
							case stationxml.CfFunctionAnalogHertz:
								return "B"
							case stationxml.CfFunctionDigital:
								return "D"
							default:
								return " "
							}
						}(),
						StageSequenceNumber:    int(s.Number),
						StageSignalInputUnits:  lookupUnitsAbbreviation(s.Coefficients.InputUnits.Name),
						StageSignalOutputUnits: lookupUnitsAbbreviation(s.Coefficients.OutputUnits.Name),
						Numerators: func() []ResponseCoefficient {
							var c []ResponseCoefficient
							for _, i := range s.Coefficients.Numerators {
								c = append(c, ResponseCoefficient{
									Coefficient: i.Value,
								})
							}
							return c
						}(),
						Denominators: func() []ResponseCoefficient {
							var c []ResponseCoefficient
							for _, i := range s.Coefficients.Denominators {
								c = append(c, ResponseCoefficient{
									Coefficient: i.Value,
								})
							}
							return c
						}(),
					}.String(),
				})
			}
			if s.FIR != nil {
				blockettes = append(blockettes, Blockette{
					Type: 61,
					Content: FIRResponse{
						StageSequenceNumber:    int(s.Number),
						StageSignalInputUnits:  lookupUnitsAbbreviation(s.FIR.InputUnits.Name),
						StageSignalOutputUnits: lookupUnitsAbbreviation(s.FIR.OutputUnits.Name),
						ResponseName: func() string {
							switch s.FIR.Name {
							case "TAURUS_100HZ_STAGE_1":
								return "TAURUSX100HZXSTAGEX1"
							case "TAURUS_100HZ_STAGE_2":
								return "TAURUSX100HZXSTAGEX2"
							case "TAURUS_100HZ_STAGE_3":
								return "TAURUSX100HZXSTAGEX3"
							case "ALTUS_BNC":
								return "ALTUSXBNC"
							case "ALTUS_A200":
								return "ALTUSXA200"
							case "BASALT_A4-50":
								return "BASALTXA4X50"
							case "BASALT_A5-50-S5C":
								return "BASALTXA5X50XS5C"
							case "BASALT_A3-50":
								return "BASALTXA3X50"
							case "BASALT_A5-50":
								return "BASALTXA5X50"
							case "BASALT_B2-80":
								return "BASALTXB2X80"
							case "QUANTERRA_VLP":
								return "QUANTERRAXVLP"
							case "Q330S+_FLbelow100-200":
								return "Q330S+XFLBELOW100X200"
							case "Q330S+_FLbelow100-100":
								return "Q330S+XFLBELOW100X100"
							case "Q330S+_FLbelow100-50":
								return "Q330S+XFLBELOW100X50"
							case "Q330S+_FLbelow100-1":
								return "Q330S+XFLBELOW100X1"
							case "Q330_FLbelow100-1":
								return "Q330XFLBELOW100X1"
							case "Q330_FLbelow100-50":
								return "Q330XFLBELOW100X50"
							case "Q330_FLbelow100-100":
								return "Q330XFLBELOW100X100"
							case "Q330_FLbelow100-200":
								return "Q330XFLBELOW100X200"
							case "QUANTERRA_F260":
								return "QUANTERRAXF260"
							case "QUANTERRA_F96C":
								return "QUANTERRAXF96C"
							case "QUANTERRA_FS2D5":
								return "QUANTERRAXFS2D5"
							case "QUANTERRA_F96CM":
								return "QUANTERRAXF96CM"
							case "QUANTERRA_FS2D5M":
								return "QUANTERRAXFS2D5M"
							case "QUANTERRA_A2D":
								return "QUANTERRAXA2D"
							case "ORION_FIR_1_FILTER":
								return "ORIONXFIRX1XFILTER"
							case "ORION_FIR_2_FILTER":
								return "ORIONXFIRX2XFILTER"
							case "ORION_FIR_3_FILTER":
								return "ORIONXFIRX3XFILTER"
							case "ORION_FIR_5_FILTER":
								return "ORIONXFIRX5XFILTER"
							default:
								return s.FIR.Name
							}
						}(),
						SymmetryCode: func() string {
							switch s.FIR.Symmetry {
							case stationxml.SymmetryNone:
								return "A"
							case stationxml.SymmetryOdd:
								return "B"
							case stationxml.SymmetryEven:
								return "C"
							default:
								return ""
							}
						}(),
						Coefficients: func() []float64 {
							var c []float64
							for _, i := range s.FIR.NumeratorCoefficients {
								c = append(c, i.Value)
							}
							return c
						}(),
					}.String(),
				})
			}

			if s.Polynomial != nil {
				blockettes = append(blockettes, Blockette{
					Type: 62,
					Content: ResponsePolynomial{
						TransferFunctionType: "P",
						StageSequenceNumber:  int(s.Number),
						StageSignalInputUnits: func() int {
							switch s.Polynomial.InputUnits.Name {
							case "hPa":
								// mistake in the original POD it seems
								return lookupUnitsAbbreviation("PA")
								//return lookupUnitsAbbreviation("M")
							default:
								return lookupUnitsAbbreviation(s.Polynomial.InputUnits.Name)
							}
						}(),
						StageSignalOutputUnits: lookupUnitsAbbreviation(s.Polynomial.OutputUnits.Name),
						PolynomialApproximationType: func() string {
							switch s.Polynomial.ApproximationType {
							case stationxml.ApproximationTypeMaclaurin:
								return "M"
							default:
								return ""
							}
						}(),
						ValidFrequencyUnits:      "B",
						LowerValidFrequencyBound: "", //s.Polynomial.FrequencyLowerBound.Value,
						UpperValidFrequencyBound: "", //s.Polynomial.FrequencyUpperBound.Value,
						LowerBoundOfApproximation: func() float64 {
							if f, err := strconv.ParseFloat(s.Polynomial.ApproximationLowerBound, 64); err == nil {
								return f
							}
							return 0.0
						}(),
						UpperBoundOfApproximation: func() float64 {
							if f, err := strconv.ParseFloat(s.Polynomial.ApproximationUpperBound, 64); err == nil {
								return f
							}
							return 0.0
						}(),
						MaximumAbsoluteError: s.Polynomial.MaximumError,
						Coefficients: func() []float64 {
							var c []float64
							for _, v := range s.Polynomial.Coefficients {
								c = append(c, v.Value)
							}
							return c
						}(),
					}.String(),
				})
			}

			if s.Decimation != nil {
				blockettes = append(blockettes, Blockette{
					Type: 57,
					Content: Decimation{
						StageSequenceNumber: int(s.Number),
						InputSampleRate:     s.Decimation.InputSampleRate.Value,
						DecimationFactor:    int(s.Decimation.Factor),
						DecimationOffset:    int(s.Decimation.Offset),
						EstimatedDelay:      s.Decimation.Delay.Value,
						CorrectionApplied:   s.Decimation.Correction.Value,
					}.String(),
				})
			}
			blockettes = append(blockettes, Blockette{
				Type: 58,
				Content: StageGain{
					StageSequenceNumber: int(s.Number),
					Gain: func() float64 {
						gain := float64(1.0)
						if s.StageGain.Value != 0.0 {
							return gain * s.StageGain.Value
						}
						return gain
					}(),
					Frequency: s.StageGain.Frequency,
				}.String(),
			})
		}

		blockettes = append(blockettes, Blockette{
			Type: 58,
			Content: StageGain{
				Gain: func() float64 {
					gain := float64(1.0)
					if cha.Response != nil {
						for _, s := range cha.Response.Stages {
							if s.Polynomial != nil {
								switch s.Polynomial.InputUnits.Name {
								case "hPa":
									gain = gain * 100.0
								}
							}
						}
						// microphones!!!
						if cha.Response.InstrumentSensitivity.Gain.Value == 107374.1824 {
							return 125 * gain * cha.Response.InstrumentSensitivity.Gain.Value / 100.0
						}
						if cha.Response.InstrumentSensitivity.Gain.Value != 0.0 {
							return gain * cha.Response.InstrumentSensitivity.Gain.Value
						}
					}
					return gain
				}(),
				Frequency: func() float64 {
					return cha.Response.InstrumentSensitivity.Gain.Frequency
				}(),
			}.String(),
		})
	}
	return blockettes
}

func (p *Pod) ChannelComments(cha *stationxml.Channel) []Blockette {
	var blockettes []Blockette

	for _, b := range cha.Comments {
		blockettes = append(blockettes,
			Blockette{
				Type: 59,
				Content: ChannelComment{
					BeginningEffectiveTime: func() time.Time {
						if cha.StartDate != nil {
							return cha.StartDate.Time
						}
						return time.Time{}
					}(),
					EndEffectiveTime: func() time.Time {
						if cha.EndDate != nil {
							return cha.EndDate.Time
						}
						return time.Time{}
					}(),
					CommentCodeKey: lookupCommentDescription(b.Value),
				}.String(),
			})
	}

	return blockettes
}
