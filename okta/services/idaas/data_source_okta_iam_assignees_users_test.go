package idaas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDataSourceIamAssigneesUsers_basic(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSIamAssigneesUsers, t.Name())
	config := mgr.GetFixtures("test_iam_assignees_users.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_iam_assignees_users.test", "users.#"),
					resource.TestCheckResourceAttr("data.okta_iam_assignees_users.test", "limit", "200"),
				),
			},
		},
	})
}

func TestAccDataSourceIamAssigneesUsers_withLimit(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaIDaaSIamAssigneesUsers, t.Name())
	config := mgr.GetFixtures("test_iam_assignees_users_with_limit.tf", t)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.okta_iam_assignees_users.test", "users.#"),
					resource.TestCheckResourceAttr("data.okta_iam_assignees_users.test", "limit", "50"),
				),
			},
		},
	})
}
