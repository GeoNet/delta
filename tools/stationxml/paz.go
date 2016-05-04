package main

var PAZs map[string]PAZ = map[string]PAZ{
	"FBA-ES-T": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "Standard response of an Kinemetric's EpiSensor FBA-ES sensor, they are built\n\t\t\twith a wide range of gains. We use +/- 20V @ +/-2 g for the National Network,\n\t\t\tand +/- 2.5V @ +/- 2g for the ETNA strong motion recorders.",
		Poles: []complex128{(-981 + 1009i), (-981 - 1009i), (-3290 + 1263i), (-3290 - 1263i)},
	},

	"FBA-23-50Hz": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "This is the standard response of a Kinemetrics FBA-23 @ a gain of +/- .5g\n\t\t\tfull scale. It is assumed that the natural frequency is 50Hz.",
		Poles: []complex128{(-1000 + 0i), (-222.1 + 222.1i), (-222.1 - 222.1i)},
	},

	"FBA-23-100Hz": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "This is the standard response of a Kinemetrics FBA-23 @ a gain of +/- .5g\n\t\t\tfull scale. It is assumed that the natural frequency is 100Hz.",
		Poles: []complex128{(-1000 + 0i), (-444.2 + 444.2i), (-444.2 - 444.2i)},
	},

	"FBA-23-50Hz-2g": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "This is the standard response of a Kinemetrics FBA-23 @ a gain of +/- 2g\n\t\t\tfull scale. It is assumed that the natural frequency is 50Hz.",
		Poles: []complex128{(-222.1 + 222.1i), (-222.1 - 222.1i), (-1500 + 0i)},
	},

	"A2D": PAZ{
		Code:  "D",
		Type:  "Digital (Z-transform).",
		Notes: "This filter is used to represent an Analogue to Digital converter stage, it has no poles or zeros.",
	},

	"CMG-3ESP-UC": PAZ{
		Code:  "B",
		Type:  "Analogue response, in Hz.",
		Notes: "These are the Canterbury Ring Laser sensor response specifications.",
		Poles: []complex128{(-0.025356 + 0.025356i), (-0.025356 - 0.025356i), (-50 - 32.2i), (-50 + 32.2i)},
		Zeros: []complex128{(138 - 144i), (0 + 0i), (0 + 0i), (138 + 144i)},
	},

	"CMG-3ESP-GN": PAZ{
		Code:  "B",
		Type:  "Analogue response, in Hz.",
		Notes: "derived from Guralp Test &amp; Calibration Data [Cal 025]",
		Poles: []complex128{(-0.01178 + 0.01178i), (-0.01178 - 0.01178i), (-180 + 0i), (-160 + 0i), (-80 + 0i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i)},
	},

	"CMG-3TB-CTBTO": PAZ{
		Code:  "B",
		Type:  "Analogue response, in Hz.",
		Notes: "Standard CTBTO specification downhole broadband sensor.",
		Poles: []complex128{(-0.00589 + 0.00589i), (-0.00589 - 0.00589i), (-73.2 - 37.6i), (-73.2 + 37.6i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i), (146.5 + 0i)},
	},

	"CMG-40T-60S-GN": PAZ{
		Code:  "B",
		Type:  "Analogue response, in Hz.",
		Notes: "Standard GeoNet specification CMG-40T as initially deployed",
		Poles: []complex128{(-0.01178 + 0.01178i), (-0.01178 - 0.01178i), (-48.4 + 0i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i), (140 + 0i)},
	},

	"CMG-40T-30S-GN": PAZ{
		Code:  "B",
		Type:  "Analogue response, in Hz.",
		Notes: "Standard GeoNet specification CMG-40T as initially deployed",
		Poles: []complex128{(-0.02365 + 0.02365i), (-0.02365 - 0.02365i), (-159 + 0i), (-66 + 0i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i), (189 + 0i)},
	},

	"L4C": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "Standard response for a Mark Products L4C (Sercel)",
		Poles: []complex128{(-4.2097 + 4.6644i), (-4.2097 - 4.6644i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i)},
	},

	"CMG-40T-30S-GN-H": PAZ{
		Code:  "B",
		Type:  "Analogue response, in Hz.",
		Notes: "Standard GeoNet specification CMG-40T as initially deployed",
		Poles: []complex128{(-91.297 - 20.011i), (-0.02365 + 0.02365i), (-0.02365 - 0.02365i), (-91.297 + 20.011i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i), (174.208 + 0i)},
	},

	"CMG-40T-30S-GNS": PAZ{
		Code:  "B",
		Type:  "Analogue response, in Hz.",
		Notes: "Standard GNS specification CMG-40T as initially deployed",
		Poles: []complex128{(-0.02356 + 0.02356i), (-0.02356 - 0.02356i), (-50 + 0i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i), (159 + 0i)},
	},

	"CMG-40T-30S-GNS-Z": PAZ{
		Code:  "B",
		Type:  "Analogue response, in Hz.",
		Notes: "Standard GeoNet specification CMG-40T as initially deployed",
		Poles: []complex128{(-0.02356 + 0.02356i), (-0.02356 - 0.02356i), (-250 + 0i), (-61 + 0i), (-200 + 0i), (-200 + 0i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i)},
	},

	"CMG-40T-30S-GNS-H": PAZ{
		Code:  "B",
		Type:  "Analogue response, in Hz.",
		Notes: "Standard GeoNet specification CMG-40T as initially deployed",
		Poles: []complex128{(-0.02356 + 0.02356i), (-250 + 0i), (-0.02356 - 0.02356i), (-61 + 0i), (-200 + 0i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i)},
	},

	"STS-2": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Poles: []complex128{(-0.03701 + 0.03701i), (-0.03701 - 0.03701i), (-131 + 467.3i), (-131 - 467.3i), (-251.3 + 0i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i)},
	},

	"SDP": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "A simple place holder for the SDP down-hole sensor",
	},

	"SSA-320": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "Taken from the UW network station ALST",
		Poles: []complex128{(-201.06 - 241.39i), (-201.06 + 241.39i)},
	},

	"EARSS-GAIN": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "A simple place holder indicating we have an analogue EARSS gain stage",
	},

	"EARSS-25Hz": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "EARSS running at 25Hz",
		Poles: []complex128{(-33.3216 - 33.3216i), (-33.3216 + 33.3216i)},
	},

	"EARSS-50Hz": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "EARSS running at 50Hz",
		Poles: []complex128{(-66.6432 - 66.6432i), (-66.6432 + 66.6432i)},
	},

	"EARSS-100Hz": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "EARSS running at 100Hz",
		Poles: []complex128{(-133.2865 - 133.2865i), (-133.2865 + 133.2865i)},
	},

	"SNARE-50Hz": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "SNARE running at 50Hz",
		Poles: []complex128{(-66.6432 + 66.6432i), (-66.6432 - 66.6432i), (-66.6432 - 66.6432i), (-66.6432 + 66.6432i)},
	},

	"EARSS-HP": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "EARSS single pole high pass filter with corner @ 2 pi Hz",
		Poles: []complex128{(-6.283185 + 0i)},
		Zeros: []complex128{(0 + 0i)},
	},

	"W1": PAZ{
		Code:  "B",
		Type:  "Analogue response, in Hz.",
		Notes: "Standard response for a Willmore Mark I",
		Poles: []complex128{(-2.5132 + 5.7586i), (-2.5132 - 5.7586i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i)},
	},

	"W2": PAZ{
		Code:  "B",
		Type:  "Analogue response, in Hz.",
		Notes: "Standard response for a Willmore Mark II",
		Poles: []complex128{(-4.5239 + 4.3604i), (-4.5239 - 4.3604i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i)},
	},

	"L4C-OD": PAZ{
		Code:  "B",
		Type:  "Analogue response, in Hz.",
		Notes: "Overdamped response for a Mark Products L4C (Sercel) - using 560 ohm resistors",
		Poles: []complex128{(13.03 + 0i), (3.03 + 0i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i)},
	},

	"BENIOFF:SP": PAZ{
		Code:  "B",
		Type:  "Analogue response, in Hz.",
		Notes: "Guessed response based on 1s critical period and critical damping",
		Poles: []complex128{(-4.44288 - 4.44288i), (-4.44288 + 4.44288i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i)},
	},

	"EARSS-XHP": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "EARSS single pole high pass filter with corner @ 2 pi * 4.7 Hz as used for CRLZ",
		Poles: []complex128{(-29.530971 + 0i)},
		Zeros: []complex128{(0 + 0i)},
	},

	"L15": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "Standard response for a Mark Products L15 4.5 Hz sensor (Sercel)",
		Poles: []complex128{(-17.8128 + 21.9577i), (-17.8128 - 21.9577i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i)},
	},

	"SENSONICS:SP3": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "A dummy PAZ filter for the unknown SP3 response",
	},

	"MICROPHONE": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "A dummy PAZ filter for the unknown microphone response",
	},

	"LE-3D": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "Standard response for a Lennartz LE-3D (Lite)",
		Poles: []complex128{(-4.21 + 4.66i), (-4.21 - 4.66i), (-2.105 + 0i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i), (0 + 0i)},
	},

	"1500Hz Bessel 3P-LP": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Poles: []complex128{(-9904.799805 + 3786i), (-9904.799805 - 3786i), (-12507 + 0i)},
	},

	"CMG-6T": PAZ{
		Code:  "B",
		Type:  "Analogue response, in Hz.",
		Notes: "Standard CMG-6T as per T6262",
		Poles: []complex128{(-0.148597 + 0.148597i), (-0.148597 - 0.148597i), (-336.766 + 136.656i), (-336.766 - 136.656i), (-2469.36 + 0i), (-47.0636 + 0i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i), (-31.6 + 0i)},
	},

	"REFTEK": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Poles: []complex128{(-1517 + 406.6i), (-1517 - 406.6i), (-1111 + 1111i), (-1111 - 1111i), (-406.6 + 1517i), (-406.6 - 1517i)},
	},

	"LE-3D-MKII": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Poles: []complex128{(-4.44 + 4.44i), (-4.44 - 4.44i), (-1.083 + 0i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i), (0 + 0i)},
	},

	"MALIN": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "derived from information passed on from Eylon Shalev <shalev@duke.edu>",
		Poles: []complex128{(-7.6654860747591 + 9.95760983645695i), (-7.6654860747591 - 9.95760983645695i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i)},
	},

	"Q330-PREAMP": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "A simple place holder indicating the PREAMP has been enables",
	},

	"HYDROPHONE": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "A simple place holder for hydrophones",
	},

	"LE-3DliteMkII": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "derived from Lennartz documentation, i.e. 990-0003 page 15 via http://www.lennartz-electronic.de",
		Poles: []complex128{(-4.444 + 4.444i), (-4.444 - 4.444i), (-1.083 + 0i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i), (0 + 0i)},
	},

	"SM-6 Geophone": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "derived from information 10K shunt and instrument docs",
		Poles: []complex128{(-19.5092903787926 + 20.4652277144474i), (-19.5092903787926 - 20.4652277144474i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i)},
	},

	"LE-3Dlite": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "derived from Lennartz documentation, i.e. 990-0003 page 15 via http://www.lennartz-electronic.de",
		Poles: []complex128{(-4.444 + 4.444i), (-4.444 - 4.444i), (-1.083 + 0i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i), (0 + 0i)},
	},

	"CMG-3T": PAZ{
		Code:  "B",
		Type:  "Analogue response, in Hz.",
		Notes: "derived from IRIS web page",
		Poles: []complex128{(-0.03701 + 0.03701i), (-0.03701 - 0.03701i), (-1131 + 0i), (-1005 + 0i), (-502.7 + 0i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i)},
	},

	"CMG-40T-30S": PAZ{
		Code:  "B",
		Type:  "Analogue response, in Hz.",
		Notes: "Standard GNS specification CMG-40T as initially deployed",
		Poles: []complex128{(-0.02356 + 0.02356i), (-0.02356 - 0.02356i), (-50 + 0i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i), (159 + 0i)},
	},

	"CUSP": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "A simple place holder for the CSI CUSP strong motion sensors",
	},

	"POL": PAZ{
		Code:  "D",
		Type:  "Digital (Z-transform).",
		Notes: "This filter is used for a polynomial response stages, it has no poles or zeros.",
	},

	"CMG-3TB-GN": PAZ{
		Code:  "B",
		Type:  "Analogue response, in Hz.",
		Notes: "derived from Guralp documentation for T35920, i.e. caldoc@guralp.com",
		Poles: []complex128{(-0.00589 - 0.00589i), (-0.00589 + 0.00589i), (-180 + 0i), (-160 + 0i), (-80 + 0i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i)},
	},

	"IESE-4.5Hz-10K-SHUNT": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "from Carolin Boese IESE, damping 0.5 gain 50.4",
		Poles: []complex128{(-14.137167 + 24.486291i), (-14.137167 - 24.486291i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i)},
	},

	"TRILLIUM-120QA": PAZ{
		Code:  "A",
		Type:  "Laplace transform analog stage response, in rad/sec.",
		Notes: "derived from Nanometrics documentation for Trillium 120QA",
		Poles: []complex128{(-0.036614 + 0.037059i), (-0.036614 - 0.037059i), (-32.55 + 0i), (-142 + 0i), (-364 + 404i), (-364 - 404i), (-1260 + 0i), (-4900 + 5200i), (-4900 - 5200i), (-7100 + 1700i)},
		Zeros: []complex128{(0 + 0i), (0 + 0i), (-31.63 + 0i), (-160 + 0i), (-350 + 0i), (-3177 + 0i)},
	},
}
