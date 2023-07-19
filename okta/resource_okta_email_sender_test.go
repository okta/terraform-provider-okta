package okta

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaEmailSender(t *testing.T) {
	t.Skip("okta_email_sender is effectively deprecated as its API has been removed")

	mgr := newFixtureManager(emailSender, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", emailSender)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(emailSender, emailSenderExists),
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
	client := apiSupplementForTest()
	sender, resp, err := client.GetEmailSender(context.Background(), id)
	if err := suppressErrorOn404(resp, err); err != nil {
		return false, err
	}
	return sender != nil && sender.Status != "DELETED", nil
}
