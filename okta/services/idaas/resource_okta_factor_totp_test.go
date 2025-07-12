package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaFactorTOTP_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSFactorTotp)
	mgr := newFixtureManager("resources", resources.OktaIDaaSFactorTotp, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		// NOTE: The publicly documented DELETE /api/v1/org/factors/hotp/profiles/{id} appears to only 501 at the present time.
		// CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSFactorTotp, doesFactorTOTPExist),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
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

// func doesFactorTOTPExist(id string) (bool, error) {
// 	client := iDaaSAPIClientForTestUtil.OktaSDKSupplementClient()
// 	_, response, err := client.GetHotpFactorProfile(context.Background(), id)
// 	return utils.DoesResourceExist(response, err)
// }
