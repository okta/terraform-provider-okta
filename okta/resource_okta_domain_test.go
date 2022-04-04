package okta

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktaDomain(t *testing.T) {
	ri := acctest.RandInt()
	mgr := newFixtureManager(domain)
	config := mgr.GetFixtures("basic.tf", ri, t)
	resourceName := fmt.Sprintf("%s.test", domain)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(domain, domainExists),
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
	domain, resp, err := getOktaClientFromMetadata(testAccProvider.Meta()).Domain.GetDomain(context.Background(), id)
	if err := suppressErrorOn404(resp, err); err != nil {
		return false, err
	}
	return domain != nil, nil
}
