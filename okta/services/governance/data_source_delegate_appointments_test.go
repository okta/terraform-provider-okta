package governance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestAccDelegateAppointments_dataSource(t *testing.T) {
	mgr := newFixtureManager("data-sources", resources.OktaGovernanceDelegateAppointments, t.Name())
	config := mgr.ConfigReplace(mgr.GetFixtures("datasource.tf", t))
	allDataSourceName := fmt.Sprintf("data.%s.all", resources.OktaGovernanceDelegateAppointments)
	filteredDataSourceName := fmt.Sprintf("data.%s.by_principal", resources.OktaGovernanceDelegateAppointments)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(allDataSourceName, "id"),
					resource.TestCheckResourceAttrSet(filteredDataSourceName, "id"),
					resource.TestCheckResourceAttrSet(filteredDataSourceName, "principal_id"),
				),
			},
		},
	})
}
