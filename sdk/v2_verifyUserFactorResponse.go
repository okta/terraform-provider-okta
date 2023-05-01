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
