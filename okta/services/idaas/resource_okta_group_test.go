package idaas_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaGroup_crud(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroup)
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroup, t.Name())
	config := mgr.GetFixtures("okta_group.tf", t)
	updatedConfig := mgr.GetFixtures("okta_group_updated.tf", t)
	acctest.BuildResourceName(mgr.Seed)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroup, doesGroupExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceNameWithPrefix("testAcc", mgr.Seed)),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceNameWithPrefix("testAcc_Different", mgr.Seed)),
				),
			},
		},
	})
}

func TestAccResourceOktaGroup_customschema(t *testing.T) {
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroup)
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroup, t.Name())
	base := mgr.GetFixtures("okta_group_custom_base.tf", t)
	updated := mgr.GetFixtures("okta_group_custom_updated.tf", t)
	removal := mgr.GetFixtures("okta_group_custom_removal.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroup, doesGroupExist),
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%s", strconv.Itoa(mgr.Seed))),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", fmt.Sprintf("{\"testSchema1_%s\":\"testing1234\",\"testSchema2_%s\":true,\"testSchema3_%s\":54321}", strconv.Itoa(mgr.Seed), strconv.Itoa(mgr.Seed), strconv.Itoa(mgr.Seed))),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%s", strconv.Itoa(mgr.Seed))),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", fmt.Sprintf("{\"testSchema1_%s\":\"moretesting1234\",\"testSchema2_%s\":false,\"testSchema3_%s\":12345}", strconv.Itoa(mgr.Seed), strconv.Itoa(mgr.Seed), strconv.Itoa(mgr.Seed))),
				),
			},
			{
				Config: removal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%s", strconv.Itoa(mgr.Seed))),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", fmt.Sprintf("{\"testSchema1_%s\":\"moretesting1234\"}", strconv.Itoa(mgr.Seed))),
				),
			},
		},
	})
}

func TestAccResourceOktaGroup_customschema_null(t *testing.T) {
	if acctest.SkipVCRTest(t) {
		return
	}
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSGroup)
	mgr := newFixtureManager("resources", resources.OktaIDaaSGroup, t.Name())
	base := mgr.GetFixtures("okta_group_custom_base.tf", t)
	nulls := mgr.GetFixtures("okta_group_custom_nulls.tf", t)
	removal := mgr.GetFixtures("okta_group_custom_removal.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSGroup, doesGroupExist),
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%s", strconv.Itoa(mgr.Seed))),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", fmt.Sprintf("{\"testSchema1_%s\":\"testing1234\",\"testSchema2_%s\":true,\"testSchema3_%s\":54321}", strconv.Itoa(mgr.Seed), strconv.Itoa(mgr.Seed), strconv.Itoa(mgr.Seed))),
				),
			},
			{
				Config: nulls,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%s", strconv.Itoa(mgr.Seed))),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", fmt.Sprintf("{\"testSchema2_%s\":true}", strconv.Itoa(mgr.Seed))),
				),
			},
			{
				Config:   nulls,
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", fmt.Sprintf("{\"testSchema2_%s\":true}", strconv.Itoa(mgr.Seed))),
				),
			},
			{
				Config: removal,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("testAcc_%s", strconv.Itoa(mgr.Seed))),
					resource.TestCheckResourceAttr(resourceName, "custom_profile_attributes", fmt.Sprintf("{\"testSchema1_%s\":\"moretesting1234\"}", strconv.Itoa(mgr.Seed))),
				),
			},
		},
	})
}

func doesGroupExist(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
	_, response, err := client.Group.GetGroup(context.Background(), id)
	return utils.DoesResourceExist(response, err)
}
