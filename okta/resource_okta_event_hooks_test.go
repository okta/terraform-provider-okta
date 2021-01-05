package okta

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func TestAccOktaEventHook_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "okta_event_hook.test"
	mgr := newFixtureManager(eventHook)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	activatedConfig := mgr.GetFixtures("basic_activated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProvidersFactories,
		CheckDestroy:      createCheckResourceDestroy(eventHook, eventHookExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, eventHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),
					testCheckResourceSetAttr(
						resourceName,
						"events",
						eventSet(&okta.EventSubscriptions{
							Type:  "EVENT_TYPE",
							Items: []string{"user.lifecycle.create", "user.lifecycle.delete.initiated"},
						}),
					),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, eventHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusInactive),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/testUpdated"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),
					testCheckResourceSetAttr(
						resourceName,
						"headers",
						testMakeEventHookHeadersSet([]*okta.EventHookChannelConfigHeader{
							{
								Key:   "x-test-header",
								Value: "test stuff",
							},
							{
								Key:   "x-another-header",
								Value: "more test stuff",
							},
						}),
					),
					testCheckResourceSetAttr(
						resourceName,
						"events",
						eventSet(&okta.EventSubscriptions{
							Type: "EVENT_TYPE",
							Items: []string{
								"user.lifecycle.create",
								"user.lifecycle.delete.initiated",
								"user.account.update_profile",
							},
						},
						),
					),
				),
			},
			{
				Config: activatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, eventHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", statusActive),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),
					testCheckResourceSetAttr(
						resourceName,
						"events",
						eventSet(&okta.EventSubscriptions{
							Type:  "EVENT_TYPE",
							Items: []string{"user.lifecycle.create", "user.lifecycle.delete.initiated"},
						}),
					),
				),
			},
		},
	})
}

func eventHookExists(id string) (bool, error) {
	eh, resp, err := getOktaClientFromMetadata(testAccProvider.Meta()).EventHook.GetEventHook(context.Background(), id)
	if err := suppressErrorOn404(resp, err); err != nil {
		return false, err
	}
	return eh != nil, nil
}

func testMakeEventHookHeadersSet(headers []*okta.EventHookChannelConfigHeader) *schema.Set {
	h := make([]interface{}, len(headers))
	for i, header := range headers {
		h[i] = map[string]interface{}{
			"key":   header.Key,
			"value": header.Value,
		}
	}
	return schema.NewSet(schema.HashResource(headerSchema), h)
}

// Create a TestCheckFunc that compares a Set to the current state
// Works for TypeSet attributes with TypeString elements or Resource elements
// Resource elements cannot be nested
func testCheckResourceSetAttr(resourceName, attribute string, set *schema.Set) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s in %s", resourceName, ms.Path)
		}

		is := rs.Primary
		if is == nil {
			return fmt.Errorf("no primary instance: %s in %s", resourceName, ms.Path)
		}

		// Look through all attributes looking for attributes that match the one we're looking for
		// and building a map keyed by the item's hash.  In the end, for a set of strings, the found
		// map will look something like this:
		// map[
		//   "123456": "value",
		//   "234567": "value2",
		// ]
		// and for a set of maps it will look something like this:
		// map[
		//   "12345": map[
		//     "key1": "value1",
		//     "key2": "value2",
		//   ]
		//   "23456": map[
		//     "key1": "value3",
		//     "key2": "value4",
		//   ]
		// ]
		found := make(map[string]interface{})
		for k, v := range is.Attributes {
			if strings.HasPrefix(k, attribute) &&
				!strings.HasSuffix(k, ".#") &&
				!strings.HasSuffix(k, ".%") {
				// k will be "attribute.12345" or "attribute.12345.subAttribute"
				// This will split the attribute key into either two or three pieces.
				// If this attribute is a set of strings, then it will be two elements:
				//    { attributeName, hash }
				// If this attribute is a set of maps, then it will be three elements:
				//    { attributeName, hash, subAttribute }
				bits := strings.SplitN(k, ".", 3)
				itemHash := bits[1]

				if len(bits) == 2 {
					found[itemHash] = v
				} else {
					subAttribute := bits[2]
					entry := found[itemHash]
					if entry == nil {
						entry = make(map[string]interface{})
					}
					entry.(map[string]interface{})[subAttribute] = v
					found[itemHash] = entry
				}
			}
		}

		newSet := &schema.Set{F: set.F}
		for _, item := range found {
			newSet.Add(item)
		}

		if !set.Equal(newSet) {
			return fmt.Errorf(
				"%s: Attribute %s does not match expecation.  Expected: %v, Found: %v",
				resourceName,
				attribute,
				set.List(),
				newSet.List(),
			)
		}

		return nil
	}
}
