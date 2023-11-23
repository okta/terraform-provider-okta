package sdk

type LogStreamSettings struct {
	// AWS EventBridge fields
	AccountId       string `json:"accountId,omitempty"`
	EventSourceName string `json:"eventSourceName,omitempty"`
	Region          string `json:"region,omitempty"`
	// Splunk fields
	Edition string `json:"edition,omitempty"`
	Host    string `json:"host,omitempty"`
	Token   string `json:"token,omitempty"`
}
