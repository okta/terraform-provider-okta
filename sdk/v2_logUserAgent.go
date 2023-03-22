package sdk

type LogUserAgent struct {
	Browser      string `json:"browser,omitempty"`
	Os           string `json:"os,omitempty"`
	RawUserAgent string `json:"rawUserAgent,omitempty"`
}
