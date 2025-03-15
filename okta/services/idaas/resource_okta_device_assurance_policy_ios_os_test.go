package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccResourceOktaPolicyDeviceAssuranceIOS_crud(t *testing.T) {
	acctest.OktaResourceTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: `resource okta_policy_device_assurance_ios test{
					name = "test"
					os_version = "12.4.5"
					jailbreak = false
					screenlock_type = toset(["BIOMETRIC"])
				  }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_policy_device_assurance_ios.test", "name", "test"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_ios.test", "os_version", "12.4.5"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_ios.test", "jailbreak", "false"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_ios.test", "screenlock_type.#", "1"),
				),
			},
			{
				Config: `resource okta_policy_device_assurance_ios test{
					name = "test"
					os_version = "12.4.6"
					jailbreak = false
					screenlock_type = toset(["BIOMETRIC", "PASSCODE"])
				  }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_policy_device_assurance_ios.test", "name", "test"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_ios.test", "os_version", "12.4.6"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_ios.test", "screenlock_type.#", "2"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_ios.test", "jailbreak", "false"),
				),
			},
		},
	})
}
