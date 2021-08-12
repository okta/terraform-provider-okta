package sdk

import (
	"github.com/okta/okta-sdk-golang/v2/okta"
)

// APISupplement not all APIs are supported by okta-sdk-golang, this will act as a supplement to the Okta SDK
type APISupplement struct {
	RequestExecutor *okta.RequestExecutor
}
