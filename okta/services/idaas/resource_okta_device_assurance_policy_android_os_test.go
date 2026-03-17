package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccResourceOktaPolicyDeviceAssuranceAndroid_crud(t *testing.T) {
	mgr := newFixtureManager("resources", "okta_policy_device_assurance_android", t.Name())
	acctest.OktaResourceTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`resource okta_policy_device_assurance_android test{
					name = "testAcc-replace_with_uuid"
					os_version = "12"
					disk_encryption_type = toset(["FULL", "USER"])
					jailbreak = false
					secure_hardware_present = true
					screenlock_type = toset(["BIOMETRIC"])
				  }`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "name", fmt.Sprintf("testAcc-%d", mgr.Seed)),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "os_version", "12"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "jailbreak", "false"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "secure_hardware_present", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "disk_encryption_type.#", "2"),
				),
			},
			{
				Config: mgr.ConfigReplace(`resource okta_policy_device_assurance_android test{
					name = "testAcc-replace_with_uuid"
					os_version = "13"
					disk_encryption_type = toset(["FULL", "USER"])
					jailbreak = false
					secure_hardware_present = true
					screenlock_type = toset(["BIOMETRIC"])
				  }`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "name", fmt.Sprintf("testAcc-%d", mgr.Seed)),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "os_version", "13"),
				),
			},
			{
				Config: mgr.ConfigReplace(`resource okta_policy_device_assurance_android test{
					name = "testAcc-replace_with_uuid"
					os_version = "13"
					disk_encryption_type = toset(["FULL", "USER"])
					jailbreak = false
					secure_hardware_present = true
					screenlock_type = toset(["BIOMETRIC"])
					grace_period {
						type = "BY_DURATION"
						expiry = "P30D"
					}
					display_remediation_mode = "SHOW"
				  }`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "name", fmt.Sprintf("testAcc-%d", mgr.Seed)),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "grace_period.type", "BY_DURATION"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "grace_period.expiry", "P30D"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "display_remediation_mode", "SHOW"),
				),
			},
			{
				Config: mgr.ConfigReplace(`resource okta_policy_device_assurance_android test{
					name = "testAcc-replace_with_uuid"
					os_version = "13"
					disk_encryption_type = toset(["FULL", "USER"])
					jailbreak = false
					secure_hardware_present = true
					screenlock_type = toset(["BIOMETRIC"])
				  }`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "name", fmt.Sprintf("testAcc-%d", mgr.Seed)),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "os_version", "13"),
					resource.TestCheckNoResourceAttr("okta_policy_device_assurance_android.test", "grace_period.type"),
					resource.TestCheckNoResourceAttr("okta_policy_device_assurance_android.test", "grace_period.expiry"),
				),
			},
		},
	})
}
