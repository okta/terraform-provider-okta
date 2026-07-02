package idaas_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

// TestAccResourceOktaIdentitySourceImport_uploadError verifies that when an
// upload step fails mid-flow the defer in runImport deletes the session so
// that a subsequent re-apply is not blocked by Okta's 5-minute rate limit.
// The VCR cassette records: POST /sessions (201), POST /bulk-upsert (400),
// DELETE /sessions/{id} (204 — the cleanup).
func TestAccResourceOktaIdentitySourceImport_uploadError(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSIdentitySourceImport, t.Name())
	config := mgr.GetFixtures("upload_error.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`Error uploading upsert users`),
			},
		},
	})
}

func TestAccResourceOktaIdentitySourceImport_basic(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSIdentitySourceImport, t.Name())
	config := mgr.GetFixtures("resource.tf", t)
	resourceName := fmt.Sprintf("%s.test", resources.OktaIDaaSIdentitySourceImport)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactoriesForTestAcc(t),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "identity_source_id", "0oaxc95befZNgrJl71d7"),
					resource.TestCheckResourceAttrSet(resourceName, "session_id"),
					resource.TestCheckResourceAttrSet(resourceName, "session_status"),
				),
			},
		},
	})
}
