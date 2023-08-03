package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPolicyDeviceAssuranceMacOS(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource okta_policy_device_assurance_macos test{
					name = "test"
					os_version = "12.4.5"
					disk_encryption_type = toset(["ALL_INTERNAL_VOLUMES"])
					secure_hardware_present = true
					screenlock_type = toset(["BIOMETRIC"])
				  }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "name", "test"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "os_version", "12.4.5"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "disk_encryption_type.#", "1"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "secure_hardware_present", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "screenlock_type.#", "1"),
				),
			},
			{
				Config: providerConfig + `
				resource okta_policy_device_assurance_macos test{
					name = "test"
					os_version = "12.4.6"
					disk_encryption_type = toset(["ALL_INTERNAL_VOLUMES"])
					secure_hardware_present = true
					screenlock_type = toset(["BIOMETRIC", "PASSCODE"])
				  }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "name", "test"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "os_version", "12.4.6"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "disk_encryption_type.#", "1"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "secure_hardware_present", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "screenlock_type.#", "2"),
				),
			},
		},
	})
}
