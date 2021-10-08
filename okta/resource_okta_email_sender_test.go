package okta

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaEmailSender(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(emailSender)
	config := mgr.GetFixtures("basic.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", emailSender)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(emailSender, emailSenderExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, emailSenderExists),
					resource.TestCheckResourceAttr(resourceName, "from_name", "testAcc_"+strconv.Itoa(ri)),
					resource.TestCheckResourceAttr(resourceName, "from_address", "no-reply@example.com"),
					resource.TestCheckResourceAttr(resourceName, "subdomain", "mail"),
					resource.TestCheckResourceAttr(resourceName, "dns_records.#", "4"),
				),
			},
		},
	})
}

func emailSenderExists(id string) (bool, error) {
	sender, resp, err := getSupplementFromMetadata(testAccProvider.Meta()).GetEmailSender(context.Background(), id)
	if err := suppressErrorOn404(resp, err); err != nil {
		return false, err
	}
	return sender != nil && sender.Status != "DELETED", nil
}
