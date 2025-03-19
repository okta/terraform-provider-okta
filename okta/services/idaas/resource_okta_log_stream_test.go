package idaas_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaLogStream_crud(t *testing.T) {
	mgr := newFixtureManager("resources", resources.OktaIDaaSLogStream, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	awsEventBridgeResourceName := fmt.Sprintf("%s.eventbridge", resources.OktaIDaaSLogStream)
	splunkResourceName := fmt.Sprintf("%s.splunk", resources.OktaIDaaSLogStream)
	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSLogStream, doesLogStreamExist),
		// Of note:
		//   Step 1:
		//     AWS log stream is created in an active status and Splunk log stream is created in an inactive status
		//   Step 2:
		//     Names and status are toggled
		//   Step 3:
		//     Import check
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(awsEventBridgeResourceName, "name", acctest.BuildResourceName(mgr.Seed)+" EventBridge"),
					resource.TestCheckResourceAttr(awsEventBridgeResourceName, "type", "aws_eventbridge"),
					resource.TestCheckResourceAttr(awsEventBridgeResourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(splunkResourceName, "name", acctest.BuildResourceName(mgr.Seed)+" Splunk"),
					resource.TestCheckResourceAttr(splunkResourceName, "type", "splunk_cloud_logstreaming"),
					resource.TestCheckResourceAttr(splunkResourceName, "status", "INACTIVE"),
				),
			},
			{
				PreConfig: func() {
					// lintignore:R018
					time.Sleep(2 * time.Second) // wait a bit for Okta to catchup between active->inactive->deleted status
				},
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(awsEventBridgeResourceName, "name", acctest.BuildResourceName(mgr.Seed)+" EventBridge Updated"),
					resource.TestCheckResourceAttr(awsEventBridgeResourceName, "type", "aws_eventbridge"),
					resource.TestCheckResourceAttr(awsEventBridgeResourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(splunkResourceName, "name", acctest.BuildResourceName(mgr.Seed)+" Splunk Updated"),
					resource.TestCheckResourceAttr(splunkResourceName, "type", "splunk_cloud_logstreaming"),
					resource.TestCheckResourceAttr(splunkResourceName, "status", "ACTIVE"),
				),
			},
			{
				ResourceName: splunkResourceName,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[splunkResourceName]
					if !ok {
						return "", fmt.Errorf("failed to find %s", splunkResourceName)
					}

					return rs.Primary.Attributes["id"], nil
				},
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					// import should only net one log stream
					if len(s) != 1 {
						return errors.New("failed to import into resource into state")
					}
					// simple check
					_type := s[0].Attributes["type"]
					scls := "splunk_cloud_logstreaming"
					if scls != "splunk_cloud_logstreaming" {
						return fmt.Errorf("expected imported log stream type to be %q got %q", scls, _type)
					}
					return nil
				},
			},
		},
	})
}

func doesLogStreamExist(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV3()
	_, response, err := client.LogStreamAPI.GetLogStream(context.Background(), id).Execute()
	return utils.DoesResourceExistV3(response, err)
}
