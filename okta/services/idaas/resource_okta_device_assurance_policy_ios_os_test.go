package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccResourceOktaPolicyDeviceAssuranceIOS_crud(t *testing.T) {
	mgr := newFixtureManager("resources", "okta_policy_device_assurance_ios", t.Name())
	acctest.OktaResourceTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`resource okta_policy_device_assurance_ios test{
					name = "testAcc-replace_with_uuid"
					os_version = "12.4.5"
					jailbreak = false
					screenlock_type = toset(["BIOMETRIC"])
				  }`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_policy_device_assurance_ios.test", "name", fmt.Sprintf("testAcc-%d", mgr.Seed)),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_ios.test", "os_version", "12.4.5"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_ios.test", "jailbreak", "false"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_ios.test", "screenlock_type.#", "1"),
				),
			},
			{
				Config: mgr.ConfigReplace(`resource okta_policy_device_assurance_ios test{
					name = "testAcc-replace_with_uuid"
					os_version = "12.4.6"
					jailbreak = false
					screenlock_type = toset(["BIOMETRIC", "PASSCODE"])
				  }`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_policy_device_assurance_ios.test", "name", fmt.Sprintf("testAcc-%d", mgr.Seed)),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_ios.test", "os_version", "12.4.6"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_ios.test", "screenlock_type.#", "2"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_ios.test", "jailbreak", "false"),
				),
			},
		},
	})
}
