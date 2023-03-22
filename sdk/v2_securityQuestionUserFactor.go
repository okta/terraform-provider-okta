package sdk

import (
	"time"
)

type SecurityQuestionUserFactor struct {
	Embedded    interface{}                        `json:"_embedded,omitempty"`
	Links       interface{}                        `json:"_links,omitempty"`
	Created     *time.Time                         `json:"created,omitempty"`
	FactorType  string                             `json:"factorType,omitempty"`
	Id          string                             `json:"id,omitempty"`
	LastUpdated *time.Time                         `json:"lastUpdated,omitempty"`
	Provider    string                             `json:"provider,omitempty"`
	Status      string                             `json:"status,omitempty"`
	Verify      *VerifyFactorRequest               `json:"verify,omitempty"`
	Profile     *SecurityQuestionUserFactorProfile `json:"profile,omitempty"`
}

func NewSecurityQuestionUserFactor() *SecurityQuestionUserFactor {
	return &SecurityQuestionUserFactor{
		FactorType: "question",
	}
}

func (a *SecurityQuestionUserFactor) IsUserFactorInstance() bool {
	return true
}
