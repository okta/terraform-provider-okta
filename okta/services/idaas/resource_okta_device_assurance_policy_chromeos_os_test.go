package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
)

func TestAccResourceOktaPolicyDeviceAssuranceChromeOS_crud(t *testing.T) {
	mgr := newFixtureManager("resources", "okta_policy_device_assurance_chromeos", t.Name())
	acctest.OktaResourceTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: mgr.ConfigReplace(`resource okta_policy_device_assurance_chromeos test{
					name = "z-testAcc-replace_with_uuid"
					tpsp_allow_screen_lock = true
				  }`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_policy_device_assurance_chromeos.test", "name", fmt.Sprintf("z-testAcc-%d", mgr.Seed)),
				),
			},
			{
				Config: mgr.ConfigReplace(`resource okta_policy_device_assurance_chromeos test{
					name = "testAcc-replace_with_uuid"
					tpsp_allow_screen_lock = true
					tpsp_browser_version = "15393.27.0"
					tpsp_builtin_dns_client_enabled = true
					tpsp_chrome_remote_desktop_app_blocked = true
					tpsp_device_enrollment_domain = "testDomain"
					tpsp_disk_encrypted = true
					tpsp_key_trust_level = "CHROME_OS_VERIFIED_MODE"
					tpsp_os_firewall = true
					tpsp_os_version = "10.0.19041.1110"
					tpsp_password_proctection_warning_trigger = "PASSWORD_PROTECTION_OFF"
					tpsp_realtime_url_check_mode = true
					tpsp_safe_browsing_protection_level = "ENHANCED_PROTECTION"
					tpsp_screen_lock_secured = true
					tpsp_site_isolation_enabled = true
				  }`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("okta_policy_device_assurance_chromeos.test", "name", fmt.Sprintf("testAcc-%d", mgr.Seed)),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_chromeos.test", "tpsp_allow_screen_lock", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_chromeos.test", "tpsp_browser_version", "15393.27.0"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_chromeos.test", "tpsp_builtin_dns_client_enabled", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_chromeos.test", "tpsp_chrome_remote_desktop_app_blocked", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_chromeos.test", "tpsp_device_enrollment_domain", "testDomain"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_chromeos.test", "tpsp_disk_encrypted", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_chromeos.test", "tpsp_key_trust_level", "CHROME_OS_VERIFIED_MODE"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_chromeos.test", "tpsp_os_firewall", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_chromeos.test", "tpsp_os_version", "10.0.19041.1110"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_chromeos.test", "tpsp_password_proctection_warning_trigger", "PASSWORD_PROTECTION_OFF"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_chromeos.test", "tpsp_realtime_url_check_mode", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_chromeos.test", "tpsp_safe_browsing_protection_level", "ENHANCED_PROTECTION"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_chromeos.test", "tpsp_screen_lock_secured", "true"),
					resource.TestCheckResourceAttr("okta_policy_device_assurance_chromeos.test", "tpsp_site_isolation_enabled", "true"),
				),
			},
		},
	})
}
