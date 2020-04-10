package okta

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-okta/sdk"
)

func TestAccOktaEventHook_crud(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "okta_event_hook.test"
	mgr := newFixtureManager(eventHook)
	config := mgr.GetFixtures("basic.tf", ri, t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", ri, t)
	activatedConfig := mgr.GetFixtures("basic_activated.tf", ri, t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: createCheckResourceDestroy(eventHook, eventHookExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, eventHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", buildResourceName(ri)),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),
					testCheckResourceSetAttr(
						resourceName,
						"events",
						eventSet(&sdk.Events{
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
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/testUpdated"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),
					testCheckResourceSetAttr(
						resourceName,
						"headers",
						testMakeHeadersSet([]sdk.Header{
							sdk.Header{
								Key:   "x-test-header",
								Value: "test stuff",
							},
							sdk.Header{
								Key:   "x-another-header",
								Value: "more test stuff",
							},
						}),
					),
					testCheckResourceSetAttr(
						resourceName,
						"events",
						eventSet(
							&sdk.Events{
								Type:  "EVENT_TYPE",
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
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),
					testCheckResourceSetAttr(
						resourceName,
						"events",
						eventSet(&sdk.Events{
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
	_, res, err := getSupplementFromMetadata(testAccProvider.Meta()).GetEventHook(id)
	return err == nil && res.StatusCode != http.StatusNotFound, err
}

func testMakeHeadersSet(headers []sdk.Header) *schema.Set {
	h := make([]interface{}, len(headers))
	for i, header := range headers {
		h[i] = map[string]interface{}{
			"key":   header.Key,
			"value": header.Value,
		}
	}

	return schema.NewSet(schema.HashResource(headerSchema), h)
}

func testCheckResourceSetAttr(resourceName string, attribute string, set *schema.Set) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s in %s", resourceName, ms.Path)
		}

		is := rs.Primary
		if is == nil {
			return fmt.Errorf("No primary instance: %s in %s", resourceName, ms.Path)
		}

		found := make(map[string]map[string]interface{})
		for k, v := range is.Attributes {
			if strings.HasPrefix(k, attribute) && !strings.HasSuffix(k, ".#") {
				bits := strings.SplitN(k, ".", 3)[1:]
				entry := found[bits[0]]
				if entry == nil {
					entry = make(map[string]interface{})
				}
				if len(bits) > 1 {
					entry[bits[1]] = v
				} else {
					entry[""] = v
				}
				found[bits[0]] = entry
			}
		}

		newSet := &schema.Set{F: set.F}
		for _, item := range found {
			if value := item[""]; value != nil {
				newSet.Add(value)
			} else {
				newSet.Add(item)
			}
		}

		expected := set.Difference(newSet)
		unexpected := newSet.Difference(set)
		if expected.Len() != 0 || unexpected.Len() != 0 {
			return fmt.Errorf(
				"%s: Attribute %s does not match expecation.  Missing values: %v, Unexpected values: %v",
				resourceName,
				attribute,
				expected.List(),
				unexpected.List(),
			)
		}

		return nil
	}
}
