package governance_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

func TestEndUserMyRequests_with_requesterFields(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceEndUserMyRequests, t.Name())
	config := mgr.GetFixtures("basic_with_requester_fields.tf", t)
	resourceName := fmt.Sprintf("%s.example", resources.OktaGovernanceEndUserMyRequests)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "SUBMITTED"),
					resource.TestCheckResourceAttr(resourceName, "entry_id", "cen123456789abcdefgh"),
					resource.TestCheckResourceAttr(resourceName, "requester_field_values.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "requester_field_values.0.id", "abcdefgh-0123-4567-8910-hgfedcba123"),
					resource.TestCheckResourceAttr(resourceName, "requester_field_values.0.value", "I need access to complete my certification."),
					resource.TestCheckResourceAttr(resourceName, "requester_field_values.1.id", "ijklmnop-a12b2-c3d4-e5f6-abcdefghi"),
					resource.TestCheckResourceAttr(resourceName, "requester_field_values.1.value", "For 5 days"),
					resource.TestCheckResourceAttr(resourceName, "requester_field_values.2.id", "tuvwxyz-0123-456-8910-zyxwvut0123"),
					resource.TestCheckResourceAttr(resourceName, "requester_field_values.2.value", "Yes"),
				),
			},
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "status", regexp.MustCompile(`^(APPROVED|CANCELED|DENIED|EXPIRED|PENDING|REJECTED)$`)),
					resource.TestCheckResourceAttr(resourceName, "entry_id", "cen123456789abcdefgh"),
					resource.TestCheckResourceAttr(resourceName, "requester_field_values.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "requester_field_values.0.id", "abcdefgh-0123-4567-8910-hgfedcba123"),
					resource.TestCheckResourceAttr(resourceName, "requester_field_values.0.value", "I need access to complete my certification."),
					resource.TestCheckResourceAttr(resourceName, "requester_field_values.1.id", "ijklmnop-a12b2-c3d4-e5f6-abcdefghi"),
					resource.TestCheckResourceAttr(resourceName, "requester_field_values.1.value", "For 5 days"),
					resource.TestCheckResourceAttr(resourceName, "requester_field_values.2.id", "tuvwxyz-0123-456-8910-zyxwvut0123"),
					resource.TestCheckResourceAttr(resourceName, "requester_field_values.2.value", "Yes"),
				),
			},
		},
	})
}

func TestEndUserMyRequests_without_requesterFields(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaGovernanceEndUserMyRequests, t.Name())
	config := mgr.GetFixtures("basic_without_requester_fields.tf", t)
	resourceName := fmt.Sprintf("%s.example", resources.OktaGovernanceEndUserMyRequests)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "SUBMITTED"),
					resource.TestCheckResourceAttr(resourceName, "entry_id", "ce5678abcdefghi12345"),
				),
			},
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "status", regexp.MustCompile(`^(APPROVED|CANCELED|DENIED|EXPIRED|PENDING|REJECTED)$`)),
				),
			},
		},
	})
}
