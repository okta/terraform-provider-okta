package idaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
)

func TestAccResourceOktaAuthServerClaim_create(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAuthServerClaim)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAuthServerClaim, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:      checkResourceDestroy(resources.OktaIDaaSAuthServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					resource.TestCheckResourceAttr(resourceName, "value_type", "EXPRESSION"),
					resource.TestCheckResourceAttr(resourceName, "value", "cool"),
					resource.TestCheckResourceAttr(resourceName, "claim_type", "RESOURCE"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusInactive),
					resource.TestCheckResourceAttr(resourceName, "name", "test_updated"),
					resource.TestCheckResourceAttr(resourceName, "value_type", "EXPRESSION"),
					resource.TestCheckResourceAttr(resourceName, "value", "cool_updated"),
					resource.TestCheckResourceAttr(resourceName, "claim_type", "RESOURCE"),
				),
			},
		},
	})
}

func TestAccResourceOktaAuthServerClaim_groupType(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSAuthServerClaim)
	swResourceName := fmt.Sprintf("%s.test_sw", resources.OktaIDaaSAuthServerClaim)
	mgr := newFixtureManager("resources", resources.OktaIDaaSAuthServerClaim, t.Name())
	config := mgr.GetFixtures("basic_group.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:          acctest.AccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:      checkResourceDestroy(resources.OktaIDaaSAuthServer, authServerExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					resource.TestCheckResourceAttr(resourceName, "group_filter_type", "EQUALS"),
					resource.TestCheckResourceAttr(resourceName, "value_type", "GROUPS"),
					resource.TestCheckResourceAttr(resourceName, "value", idaas.GroupProfileEveryone),
					resource.TestCheckResourceAttr(resourceName, "claim_type", "RESOURCE"),

					resource.TestCheckResourceAttr(swResourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(swResourceName, "name", "test_sw"),
					resource.TestCheckResourceAttr(swResourceName, "group_filter_type", "STARTS_WITH"),
					resource.TestCheckResourceAttr(swResourceName, "value_type", "GROUPS"),
					resource.TestCheckResourceAttr(swResourceName, "value", "Every"),
					resource.TestCheckResourceAttr(swResourceName, "claim_type", "RESOURCE"),
				),
			},
		},
	})
}
