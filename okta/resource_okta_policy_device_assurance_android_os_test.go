package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaPolicyDeviceAssuranceAndroid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource okta_policy_device_assurance_android test{
					name = "test"
					os_version = "12"
					disk_encryption_type = toset(["FULL", "USER"])
					jailbreak = false
					secure_hardware_present = true
					screenlock_type = toset(["BIOMETRIC"])
				  }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "name", "test"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "os_version", "12"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "jailbreak", "false"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "secure_hardware_present", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "disk_encryption_type.#", "2"),
				),
			},
			{
				Config: `resource okta_policy_device_assurance_android test{
					name = "test"
					os_version = "13"
					disk_encryption_type = toset(["FULL", "USER"])
					jailbreak = false
					secure_hardware_present = true
					screenlock_type = toset(["BIOMETRIC"])
				  }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "name", "test"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_android.test", "os_version", "13"),
				),
			},
		},
	})
}
