package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaPolicyDeviceAssuranceMacOS(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource okta_policy_device_assurance_macos test{
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
					resource.TestCheckNoResourceAttr("okta_policy_device_assurance_macos.test", "tpsp_browser_version"),
				),
			},
			{
				Config: `resource okta_policy_device_assurance_macos test{
					name = "test"
					os_version = "12.4.6"
					disk_encryption_type = toset(["ALL_INTERNAL_VOLUMES"])
					secure_hardware_present = true
					screenlock_type = toset(["BIOMETRIC", "PASSCODE"])
					third_party_signal_providers = true
					tpsp_browser_version = "15393.27.0"
					tpsp_builtin_dns_client_enabled = true
					tpsp_chrome_remote_desktop_app_blocked = true
					tpsp_device_enrollment_domain = "testDomain"
					tpsp_disk_encrypted = true
					tpsp_key_trust_level = "CHROME_BROWSER_HW_KEY"
					tpsp_os_firewall = true
					tpsp_os_version = "10.0.19041"
					tpsp_password_proctection_warning_trigger = "PASSWORD_PROTECTION_OFF"
					tpsp_realtime_url_check_mode = true
					tpsp_safe_browsing_protection_level = "ENHANCED_PROTECTION"
					tpsp_screen_lock_secured = true
					tpsp_site_isolation_enabled = true
				  }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "name", "test"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "os_version", "12.4.6"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "disk_encryption_type.#", "1"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "secure_hardware_present", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "screenlock_type.#", "2"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "tpsp_browser_version", "15393.27.0"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "tpsp_builtin_dns_client_enabled", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "tpsp_chrome_remote_desktop_app_blocked", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "tpsp_device_enrollment_domain", "testDomain"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "tpsp_disk_encrypted", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "tpsp_key_trust_level", "CHROME_BROWSER_HW_KEY"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "tpsp_os_firewall", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "tpsp_os_version", "10.0.19041"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "tpsp_password_proctection_warning_trigger", "PASSWORD_PROTECTION_OFF"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "tpsp_realtime_url_check_mode", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "tpsp_safe_browsing_protection_level", "ENHANCED_PROTECTION"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "tpsp_screen_lock_secured", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_macos.test", "tpsp_site_isolation_enabled", "true"),
				),
			},
		},
	})
}
