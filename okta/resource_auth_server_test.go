package okta

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/okta/okta-sdk-golang/okta/query"
)

func findTestAuthServer(name string) bool {
	return strings.HasPrefix(name, testResourcePrefix)
}

func deleteAuthServers(client *testClient) error {
	servers, err := client.apiSupplement.FilterAuthServers(&query.Params{}, []*AuthorizationServer{}, findTestAuthServer)
	if err != nil {
		return err
	}

	for _, s := range servers {
		if _, err := client.apiSupplement.DeactivateAuthorizationServer(s.Id); err != nil {
			return err
		}
		if _, err := client.apiSupplement.DeleteAuthorizationServer(s.Id); err != nil {
			return err
		}

	}
	return nil
}

func authServerExists(id string) (bool, error) {
	client := getSupplementFromMetadata(testAccProvider.Meta())
	_, resp, err := client.GetAuthorizationServer(id)
	return resp.StatusCode != 404 && err == nil, err
}

func TestAccOktaAuthServerCrud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.sun_also_rises", authServer)
	name := buildResourceName(ri)
	mgr := newFixtureManager("okta_auth_server")
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(authServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, authServerExists),
					resource.TestCheckResourceAttr(resourceName, "description", "The best way to find out if you can trust somebody is to trust them."),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "audiences.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials_rotation_mode", "AUTO"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, authServerExists),
					resource.TestCheckResourceAttr(resourceName, "description", "The past is not dead. In fact, it's not even past."),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "audiences.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials_rotation_mode", "AUTO"),
				),
			},
		},
	})
}

func TestAccOktaAuthServerFullStack(t *testing.T) {
	ri := acctest.RandInt()
	name := buildResourceName(ri)
	resourceName := fmt.Sprintf("%s.test", authServer)
	claimName := fmt.Sprintf("%s.test", authServerClaim)
	ruleName := fmt.Sprintf("%s.test", authServerPolicyRule)
	policyName := fmt.Sprintf("%s.test", authServerPolicy)
	scopeName := fmt.Sprintf("%s.test", authServerScope)
	mgr := newFixtureManager("okta_auth_server")
	config := mgr.GetFixtures("full_stack.tf", ri, t)
	updatedConfig := mgr.GetFixtures("full_stack_updated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(authServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, authServerExists),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "audiences.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials_rotation_mode", "AUTO"),

					resource.TestCheckResourceAttr(scopeName, "name", "test:something"),
					resource.TestCheckResourceAttr(claimName, "name", "test"),
					resource.TestCheckResourceAttr(policyName, "name", "test"),
					resource.TestCheckResourceAttr(ruleName, "name", "test"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, authServerExists),
					resource.TestCheckResourceAttr(resourceName, "description", "test_updated"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "audiences.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials_rotation_mode", "AUTO"),

					resource.TestCheckResourceAttr(scopeName, "name", "test:something"),
					resource.TestCheckResourceAttr(claimName, "name", "test"),
					resource.TestCheckResourceAttr(policyName, "name", "test"),
					resource.TestCheckResourceAttr(ruleName, "name", "test"),
				),
			},
		},
	})
}
