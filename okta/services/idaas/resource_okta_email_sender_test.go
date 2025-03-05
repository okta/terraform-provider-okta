package idaas_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaEmailSender_crud(t *testing.T) {
	t.Skip("okta_email_sender is effectively deprecated as its API has been removed")

	mgr := newFixtureManager("resources", resources.OktaIDaaSEmailSender, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSEmailSender)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSEmailSender, emailSenderExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, emailSenderExists),
					resource.TestCheckResourceAttr(resourceName, "from_name", "testAcc_"+strconv.Itoa(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "from_address", "no-reply@example.com"),
					resource.TestCheckResourceAttr(resourceName, "subdomain", "mail"),
					resource.TestCheckResourceAttr(resourceName, "dns_records.#", "4"),
				),
			},
		},
	})
}

func emailSenderExists(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKSupplementClient()
	sender, resp, err := client.GetEmailSender(context.Background(), id)
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return false, err
	}
	return sender != nil && sender.Status != "DELETED", nil
}
