package sdk

type LogRequest struct {
	IpChain []*LogIpAddress `json:"ipChain,omitempty"`
}
