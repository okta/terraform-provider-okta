package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaPolicyDeviceAssuranceWindows(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource okta_policy_device_assurance_windows test{
					name = "test"
					os_version = "12.4.5"
					disk_encryption_type = toset(["ALL_INTERNAL_VOLUMES"])
					secure_hardware_present = true
					screenlock_type = toset(["BIOMETRIC"])
				  }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "name", "test"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "os_version", "12.4.5"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "disk_encryption_type.#", "1"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "secure_hardware_present", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "screenlock_type.#", "1"),
					resource.TestCheckNoResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_browser_version"),
				),
			},
			{
				Config: providerConfig + `
				resource okta_policy_device_assurance_windows test{
					name = "test"
					os_version = "12.4.6"
					disk_encryption_type = toset(["ALL_INTERNAL_VOLUMES"])
					secure_hardware_present = true
					screenlock_type = toset(["BIOMETRIC", "PASSCODE"])
					third_party_signal_providers = true
					tpsp_browser_version = "15393.27.0"
					tpsp_builtin_dns_client_enabled = true
					tpsp_chrome_remote_desktop_app_blocked = true
					tpsp_crowd_strike_agent_id = "testAgentId"
					tpsp_crowd_strike_customer_id = "testCustomerId"
					tpsp_device_enrollment_domain = "testDomain"
					tpsp_disk_encrypted = true
					tpsp_key_trust_level = "CHROME_BROWSER_HW_KEY"
					tpsp_os_firewall = true
					tpsp_os_version = "10.0.19041"
					tpsp_password_proctection_warning_trigger = "PASSWORD_PROTECTION_OFF"
					tpsp_realtime_url_check_mode = true
					tpsp_safe_browsing_protection_level = "ENHANCED_PROTECTION"
					tpsp_screen_lock_secured = true
					tpsp_secure_boot_enabled = true
					tpsp_site_isolation_enabled = true
					tpsp_third_party_blocking_enabled = true
					tpsp_windows_machine_domain = "testMachineDomain"
					tpsp_windows_user_domain = "testUserDomain"
				  }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "name", "test"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "os_version", "12.4.6"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "disk_encryption_type.#", "1"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "secure_hardware_present", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "screenlock_type.#", "2"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_browser_version", "15393.27.0"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_builtin_dns_client_enabled", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_chrome_remote_desktop_app_blocked", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_crowd_strike_agent_id", "testAgentId"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_crowd_strike_customer_id", "testCustomerId"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_device_enrollment_domain", "testDomain"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_disk_encrypted", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_key_trust_level", "CHROME_BROWSER_HW_KEY"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_os_firewall", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_os_version", "10.0.19041"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_password_proctection_warning_trigger", "PASSWORD_PROTECTION_OFF"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_realtime_url_check_mode", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_safe_browsing_protection_level", "ENHANCED_PROTECTION"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_screen_lock_secured", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_secure_boot_enabled", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_site_isolation_enabled", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_third_party_blocking_enabled", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_windows_machine_domain", "testMachineDomain"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_windows.test", "tpsp_windows_user_domain", "testUserDomain"),
				),
			},
		},
	})
}
