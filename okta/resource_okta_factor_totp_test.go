package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaFactorTOTP(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", factorTotp)
	mgr := newFixtureManager(factorTotp, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	oktaResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(factorTotp, doesFactorTOTPExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "otp_length", "10"),
					resource.TestCheckResourceAttr(resourceName, "hmac_algorithm", "HMacSHA256"),
					resource.TestCheckResourceAttr(resourceName, "time_step", "30"),
					resource.TestCheckResourceAttr(resourceName, "clock_drift_interval", "10"),
					resource.TestCheckResourceAttr(resourceName, "shared_secret_encoding", "hexadecimal"),
				),
			},
		},
	})
}

func doesFactorTOTPExist(id string) (bool, error) {
	client := apiSupplementForTest()
	_, response, err := client.GetHotpFactorProfile(context.Background(), id)
	return doesResourceExist(response, err)
}
