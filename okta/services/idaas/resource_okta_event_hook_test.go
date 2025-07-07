package idaas_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaEventHook_basic(t *testing.T) {
	resourceName := "okta_event_hook.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy("okta_event_hook", eventHookExists),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOktaEventHook_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Example Event Hook"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "events.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "events.0", "user.lifecycle.create"),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/webhook"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "x-api-token"),
					// Note: auth.value is not stored in state due to ephemeral implementation
				),
			},
		},
	})
}

func TestAccResourceOktaEventHook_update(t *testing.T) {
	resourceName := "okta_event_hook.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy("okta_event_hook", eventHookExists),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOktaEventHook_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Example Event Hook"),
				),
			},
			{
				Config: testAccResourceOktaEventHook_updated(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Updated Event Hook"),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "events.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceOktaEventHook_withHeaders(t *testing.T) {
	resourceName := "okta_event_hook.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy("okta_event_hook", eventHookExists),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOktaEventHook_withHeaders(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Event Hook with Headers"),
					resource.TestCheckResourceAttr(resourceName, "headers.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceOktaEventHook_noAuth(t *testing.T) {
	resourceName := "okta_event_hook.test"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy("okta_event_hook", eventHookExists),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOktaEventHook_noAuth(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Event Hook No Auth"),
					resource.TestCheckNoResourceAttr(resourceName, "auth"),
				),
			},
		},
	})
}

func testAccResourceOktaEventHook_basic() string {
	return `
resource "okta_event_hook" "test" {
  name   = "Example Event Hook"
  status = "ACTIVE"
  events = ["user.lifecycle.create"]
  
  channel {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/webhook"
  }
  
  auth {
    type  = "HEADER"
    key   = "x-api-token"
    value = "secret-token"
  }
}
`
}

func testAccResourceOktaEventHook_updated() string {
	return `
resource "okta_event_hook" "test" {
  name   = "Updated Event Hook"
  status = "INACTIVE"
  events = ["user.lifecycle.create", "user.lifecycle.delete.initiated"]
  
  channel {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/webhook"
  }
  
  auth {
    type  = "HEADER"
    key   = "x-api-token"
    value = "updated-secret-token"
  }
}
`
}

func testAccResourceOktaEventHook_withHeaders() string {
	return `
resource "okta_event_hook" "test" {
  name   = "Event Hook with Headers"
  status = "ACTIVE"
  events = ["user.lifecycle.create"]
  
  channel {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/webhook"
  }
  
  auth {
    type  = "HEADER"
    key   = "x-api-token"
    value = "secret-token"
  }
  
  headers {
    key   = "Content-Type"
    value = "application/json"
  }
  
  headers {
    key   = "User-Agent"
    value = "Terraform-Provider-Okta"
  }
}
`
}

func testAccResourceOktaEventHook_noAuth() string {
	return `
resource "okta_event_hook" "test" {
  name   = "Event Hook No Auth"
  status = "ACTIVE"
  events = ["user.lifecycle.create"]
  
  channel {
    type    = "HTTP"
    version = "1.0.0"
    uri     = "https://example.com/webhook"
  }
}
`
}

func eventHookExists(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV2()
	eh, resp, err := client.EventHook.GetEventHook(context.Background(), id)
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return false, err
	}
	return eh != nil, nil
}
