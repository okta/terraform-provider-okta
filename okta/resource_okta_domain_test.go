package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDomain(t *testing.T) {
	mgr := newFixtureManager(domain, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", domain)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(domain, domainExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, domainExists),
					resource.TestCheckResourceAttr(resourceName, "name", "example.com"),
					resource.TestCheckResourceAttr(resourceName, "dns_records.#", "2"),
				),
			},
		},
	})
}

func domainExists(id string) (bool, error) {
	client := oktaClientForTest()
	domain, resp, err := client.Domain.GetDomain(context.Background(), id)
	if err := suppressErrorOn404(resp, err); err != nil {
		return false, err
	}
	return domain != nil, nil
}
