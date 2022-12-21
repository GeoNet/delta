package stationxml

// ToSiteSurey returns the survey notes fro a given survey method.
func ToSiteSurvey(survey string) string {
	switch survey {
	case "Unknown":
		return "Location estimation method is unknown"
	case "External GPS Device":
		return "Location estimated from external GPS measurement"
	case "Internal GPS Clock":
		return "Location estimated from internal GPS clock"
	case "Topographic Map":
		return "Location estimated from topographic map"
	case "Site Survey":
		return "Location estimated from plans and survey to mark"
	default:
		return "Location estimation method is unknown"
	}
}
