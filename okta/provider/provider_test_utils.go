package provider

import (
	"os"
	"testing"

	"github.com/okta/okta-sdk-golang/v4/okta"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/sdk"
)

var (
	testSdkV5Client         *v5okta.APIClient
	testSdkV3Client         *okta.APIClient
	testSdkV2Client         *sdk.Client
	testSdkSupplementClient *sdk.APISupplement
)

const (
	ResourcePrefixForTest = "testAcc"
)

func SdkV5ClientForTest() *v5okta.APIClient {
	if testSdkV5Client != nil {
		return testSdkV5Client
	}
	return config.SdkV5ClientForTest
}

func SdkV3ClientForTest() *okta.APIClient {
	if testSdkV3Client != nil {
		return testSdkV3Client
	}
	return config.SdkV3ClientForTest
}

func SdkV2ClientForTest() *sdk.Client {
	if testSdkV2Client != nil {
		return testSdkV2Client
	}
	return config.SdkV2ClientForTest
}

func SdkSupplementClientForTest() *sdk.APISupplement {
	if testSdkSupplementClient != nil {
		return testSdkSupplementClient
	}
	return config.SdkSupplementClientForTest
}

func SkipVCRTest(t *testing.T) bool {
	skip := os.Getenv("OKTA_VCR_TF_ACC") != ""
	if skip {
		t.Skipf("test %q is not VCR compatible", t.Name())
	}
	return skip
}
