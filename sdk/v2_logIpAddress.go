// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type LogIpAddress struct {
	GeographicalContext *LogGeographicalContext `json:"geographicalContext,omitempty"`
	Ip                  string                  `json:"ip,omitempty"`
	Source              string                  `json:"source,omitempty"`
	Version             string                  `json:"version,omitempty"`
}
