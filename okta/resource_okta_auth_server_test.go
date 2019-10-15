package okta

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/okta/okta-sdk-golang/okta/query"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
)

func findTestAuthServer(name string) bool {
	return strings.HasPrefix(name, testResourcePrefix)
}

func deleteAuthServers(client *testClient) error {
	servers, err := client.apiSupplement.FilterAuthServers(&query.Params{}, []*sdk.AuthorizationServer{}, findTestAuthServer)
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
	server, resp, err := client.GetAuthorizationServer(id)
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return server.Id != "" && err == nil, err
}

func TestAccOktaAuthServer_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := fmt.Sprintf("%s.sun_also_rises", authServer)
	name := buildResourceName(ri)
	mgr := newFixtureManager(authServer)
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

func TestAccOktaAuthServer_fullStack(t *testing.T) {
	ri := acctest.RandInt()
	name := buildResourceName(ri)
	resourceName := fmt.Sprintf("%s.test", authServer)
	claimName := fmt.Sprintf("%s.test", authServerClaim)
	ruleName := fmt.Sprintf("%s.test", authServerPolicyRule)
	policyName := fmt.Sprintf("%s.test", authServerPolicy)
	scopeName := fmt.Sprintf("%s.test", authServerScope)
	mgr := newFixtureManager(authServer)
	config := mgr.GetFixtures("full_stack.tf", ri, t)
	updatedConfig := mgr.GetFixtures("full_stack_with_client.tf", ri, t)

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
					resource.TestCheckResourceAttr(policyName, "client_whitelist.#", "1"),
					resource.TestCheckResourceAttr(ruleName, "name", "test"),
				),
			},
		},
	})
}

func TestAccOktaAuthServer_gh299(t *testing.T) {
	ri := acctest.RandInt()
	name := buildResourceName(ri)
	resourceName := fmt.Sprintf("%s.test", authServer)
	resource2Name := fmt.Sprintf("%s.test1", authServer)
	mgr := newFixtureManager(authServer)
	config := mgr.GetFixtures("dependency.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(authServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, authServerExists),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "audiences.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials_rotation_mode", "AUTO"),

					resource.TestCheckResourceAttr(resource2Name, "name", name+"1"),
					resource.TestCheckResourceAttr(resource2Name, "audiences.#", "1"),
					resource.TestCheckResourceAttr(resource2Name, "credentials_rotation_mode", "MANUAL"),
				),
			},
		},
	})
}
