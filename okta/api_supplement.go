package okta

import (
	"github.com/okta/okta-sdk-golang/okta"
)

// Not all APIs are supported by okta-sdk-golang, this will act as a supplement to the Okta SDK
type ApiSupplement struct {
	requestExecutor *okta.RequestExecutor
}
