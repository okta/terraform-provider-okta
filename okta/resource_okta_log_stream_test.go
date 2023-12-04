package okta

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOktaLogStream_crud(t *testing.T) {
	mgr := newFixtureManager("resources", logStream, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	awsEventBridgeResourceName := fmt.Sprintf("%s.eventbridge_log_stream_example", logStream)
	splunkResourceName := fmt.Sprintf("%s.splunk_log_stream_example", logStream)

	oktaResourceTest(t, resource.TestCase{
		PreCheck:          testAccPreCheck(t),
		ErrorCheck:        testAccErrorChecks(t),
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      checkResourceDestroy(logStream, doesLogStreamExist),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(awsEventBridgeResourceName, "name", buildResourceName(mgr.Seed)+" EventBridge"),
					resource.TestCheckResourceAttr(awsEventBridgeResourceName, "type", "aws_eventbridge"),
					resource.TestCheckResourceAttr(awsEventBridgeResourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(awsEventBridgeResourceName, "settings.#", "1"),
					resource.TestCheckResourceAttr(splunkResourceName, "name", buildResourceName(mgr.Seed)+" Splunk"),
					resource.TestCheckResourceAttr(splunkResourceName, "type", "splunk_cloud_logstreaming"),
					resource.TestCheckResourceAttr(splunkResourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(splunkResourceName, "settings.#", "1"),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(5 * time.Second) // wait a bit for Okta to catchup between active->inactive->deleted status
				},
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(awsEventBridgeResourceName, "name", buildResourceName(mgr.Seed)+" EventBridge Updated"),
					resource.TestCheckResourceAttr(awsEventBridgeResourceName, "type", "aws_eventbridge"),
					resource.TestCheckResourceAttr(awsEventBridgeResourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(awsEventBridgeResourceName, "settings.#", "1"),
					resource.TestCheckResourceAttr(splunkResourceName, "name", buildResourceName(mgr.Seed)+" Splunk Updated"),
					resource.TestCheckResourceAttr(splunkResourceName, "type", "splunk_cloud_logstreaming"),
					resource.TestCheckResourceAttr(splunkResourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(splunkResourceName, "settings.#", "1"),
				),
			},
		},
	})
}

func doesLogStreamExist(id string) (bool, error) {
	client := sdkV3ClientForTest()
	_, response, err := client.LogStreamAPI.GetLogStream(context.Background(), id).Execute()
	return doesResourceExistV3(response, err)
}
