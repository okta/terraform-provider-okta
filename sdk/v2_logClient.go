package sdk

type LogClient struct {
	Device              string                  `json:"device,omitempty"`
	GeographicalContext *LogGeographicalContext `json:"geographicalContext,omitempty"`
	Id                  string                  `json:"id,omitempty"`
	IpAddress           string                  `json:"ipAddress,omitempty"`
	UserAgent           *LogUserAgent           `json:"userAgent,omitempty"`
	Zone                string                  `json:"zone,omitempty"`
}
