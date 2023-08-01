package okta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the HashiCups client is properly configured.
	// It is also possible to use the HASHICUPS_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"okta": providerserver.NewProtocol6WithError(NewFWProvider("test")),
	}
)

func TestAccPolicyDeviceAssuranceAndroid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource okta_policy_device_assurance_android test{
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
				Config: providerConfig + `
				resource okta_policy_device_assurance_android test{
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
