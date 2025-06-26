package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccResourceOktaSMTPServer_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSEmailSMTPServer, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSEmailSMTPServer)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "host", "192.168.2.0"),
					resource.TestCheckResourceAttr(resourceName, "port", "8086"),
					resource.TestCheckResourceAttr(resourceName, "username", "test_user"),
					resource.TestCheckResourceAttr(resourceName, "alias", "CustomisedServer"),
				),
			},
		},
	})
}
