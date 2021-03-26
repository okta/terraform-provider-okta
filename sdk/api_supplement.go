package sdk

import (
	"github.com/okta/okta-sdk-golang/v2/okta"
)

// ApiSupplement not all APIs are supported by okta-sdk-golang, this will act as a supplement to the Okta SDK
type ApiSupplement struct {
	RequestExecutor *okta.RequestExecutor
}
