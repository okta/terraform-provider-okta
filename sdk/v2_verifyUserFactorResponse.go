// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"time"
)

type VerifyUserFactorResponse struct {
	Embedded            interface{} `json:"_embedded,omitempty"`
	Links               interface{} `json:"_links,omitempty"`
	ExpiresAt           *time.Time  `json:"expiresAt,omitempty"`
	FactorResult        string      `json:"factorResult,omitempty"`
	FactorResultMessage string      `json:"factorResultMessage,omitempty"`
}

func NewVerifyUserFactorResponse() *VerifyUserFactorResponse {
	return &VerifyUserFactorResponse{}
}

func (a *VerifyUserFactorResponse) IsUserFactorInstance() bool {
	return true
}
