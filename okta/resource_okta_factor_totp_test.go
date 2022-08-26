package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaFactorTOTP(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.test", factorTotp)
	mgr := newFixtureManager(factorTotp)
	config := mgr.GetFixtures("basic.tf", ri, t)
	resource.Test(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(factorTotp, doesFactorTOTPExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
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
	_, response, err := getSupplementFromMetadata(testAccProvider.Meta()).GetHotpFactorProfile(context.Background(), id)
	return doesResourceExist(response, err)
}
