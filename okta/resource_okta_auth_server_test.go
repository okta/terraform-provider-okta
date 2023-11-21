package okta

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaAuthServer_crud(t *testing.T) {
	mgr := newFixtureManager("resources", authServer, t.Name())
	resourceName := fmt.Sprintf("%s.sun_also_rises", authServer)
	name := buildResourceName(mgr.Seed)
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(authServer, authServerExists),
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

func TestAccResourceOktaAuthServer_fullStack(t *testing.T) {
	mgr := newFixtureManager("resources", authServer, t.Name())
	name := buildResourceName(mgr.Seed)
	resourceName := fmt.Sprintf("%s.test", authServer)
	claimName := fmt.Sprintf("%s.test", authServerClaim)
	ruleName := fmt.Sprintf("%s.test", authServerPolicyRule)
	policyName := fmt.Sprintf("%s.test", authServerPolicy)
	scopeName := fmt.Sprintf("%s.test", authServerScope)
	config := mgr.GetFixtures("full_stack.tf", t)
	updatedConfig := mgr.GetFixtures("full_stack_with_client.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(authServer, authServerExists),
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

func TestAccResourceOktaAuthServer_gh299(t *testing.T) {
	mgr := newFixtureManager("resources", authServer, t.Name())
	name := buildResourceName(mgr.Seed)
	resourceName := fmt.Sprintf("%s.test", authServer)
	resource2Name := fmt.Sprintf("%s.test1", authServer)
	config := mgr.GetFixtures("dependency.tf", t)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(authServer, authServerExists),
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

func authServerExists(id string) (bool, error) {
	client := sdkV2ClientForTest()
	server, resp, err := client.AuthorizationServer.GetAuthorizationServer(context.Background(), id)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return server != nil && server.Id != "" && err == nil, err
}
