package sdk

import (
	"github.com/okta/okta-sdk-golang/v2/okta"
)

// APISupplement not all APIs are supported by okta-sdk-golang, this will act as a supplement to the Okta SDK
type APISupplement struct {
	RequestExecutor *okta.RequestExecutor
}

// CloneRequestExecutor create a clone of the underlying request executor
func (m *APISupplement) cloneRequestExecutor() *okta.RequestExecutor {
	a := *m.RequestExecutor
	return &a
}
