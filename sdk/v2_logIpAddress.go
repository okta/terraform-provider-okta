package sdk

type LogIpAddress struct {
	GeographicalContext *LogGeographicalContext `json:"geographicalContext,omitempty"`
	Ip                  string                  `json:"ip,omitempty"`
	Source              string                  `json:"source,omitempty"`
	Version             string                  `json:"version,omitempty"`
}
