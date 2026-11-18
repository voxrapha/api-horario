package service

import (
	_ "embed"
	"strings"
)

// timezones.txt is generated from Go 1.27.4's bundled IANA timezone archive.
//
//go:embed timezones.txt
var timezonesData string

var timezoneIDs = strings.Fields(timezonesData)

var timezoneIDSet = func() map[string]struct{} {
	identifiers := make(map[string]struct{}, len(timezoneIDs))
	for _, identifier := range timezoneIDs {
		identifiers[identifier] = struct{}{}
	}
	return identifiers
}()

// TimezonesResponse contains the IANA timezone identifiers supported by the API.
type TimezonesResponse struct {
	Timezones []string `json:"timezones"`
}

// IsTimezoneSupported reports whether the identifier belongs to the bundled dataset.
func IsTimezoneSupported(timezone string) bool {
	_, supported := timezoneIDSet[timezone]
	return supported
}

// ListarTimezones returns the identifiers from the versioned timezone dataset.
func ListarTimezones() TimezonesResponse {
	return TimezonesResponse{
		Timezones: append([]string(nil), timezoneIDs...),
	}
}
