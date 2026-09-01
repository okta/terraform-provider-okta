package idaas_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/acctest"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/services/idaas"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func TestAccResourceOktaEventHook_crud(t *testing.T) {
	resourceName := "okta_event_hook.test"
	mgr := newFixtureManager("resources", resources.OktaIDaaSEventHook, t.Name())
	config := mgr.GetFixtures("basic.tf", t)
	updatedConfig := mgr.GetFixtures("basic_updated.tf", t)
	activatedConfig := mgr.GetFixtures("basic_activated.tf", t)

	header_key_1 := "x-test-header"
	header_value_1 := "test stuff"
	header_key_2 := "x-another-header"
	header_value_2 := "more test stuff"

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSEventHook, eventHookExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, eventHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),
					testCheckResourceSetAttr(
						resourceName,
						"events",
						utils.ConvertStringSliceToSet([]string{"user.lifecycle.create", "user.lifecycle.delete.initiated"}),
					),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, eventHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusInactive),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/testUpdated"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),
					testCheckResourceSetAttr(
						resourceName,
						"headers",
						testMakeEventHookHeadersSet([]*v5okta.EventHookChannelConfigHeader{
							{Key: &header_key_1, Value: &header_value_1},
							{Key: &header_key_2, Value: &header_value_2},
						}),
					),
					testCheckResourceSetAttr(
						resourceName,
						"events",
						idaas.EventSet(v5okta.NewEventSubscriptions([]string{
							"user.lifecycle.create",
							"user.lifecycle.delete.initiated",
							"user.account.update_profile",
						}, "EVENT_TYPE")),
					),
				),
			},
			{
				Config: activatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, eventHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),
					testCheckResourceSetAttr(
						resourceName,
						"events",
						utils.ConvertStringSliceToSet([]string{"user.lifecycle.create", "user.lifecycle.delete.initiated"}),
					),
				),
			},
		},
	})
}

func TestAccResourceOktaEventHook_withFilter(t *testing.T) {
	resourceName := "okta_event_hook.test"
	mgr := newFixtureManager("resources", resources.OktaIDaaSEventHook, t.Name())
	config := mgr.GetFixtures("basic_with_filter.tf", t)
	updatedConfig := mgr.GetFixtures("basic_with_filter_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSEventHook, eventHookExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, eventHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),
					resource.TestCheckResourceAttr(resourceName, "filter.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "filter.*", map[string]string{
						"event":     "group.user_membership.add",
						"condition": "event.target.?[type eq 'UserGroup'].size()>0 && event.target.?[displayName eq 'Sales'].size()>0",
					}),
					testCheckResourceSetAttr(
						resourceName,
						"events",
						idaas.EventSet(v5okta.NewEventSubscriptions([]string{"group.user_membership.add"}, "EVENT_TYPE")),
					),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, eventHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test-updated"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),
					resource.TestCheckResourceAttr(resourceName, "filter.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "filter.*", map[string]string{
						"event":     "group.user_membership.add",
						"condition": "event.target.?[type eq 'UserGroup'].size()>0 && event.target.?[displayName eq 'Marketing'].size()>0",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "filter.*", map[string]string{
						"event":     "user.lifecycle.create",
						"condition": "event.actor.id != null",
					}),
					testCheckResourceSetAttr(
						resourceName,
						"events",
						idaas.EventSet(v5okta.NewEventSubscriptions([]string{"group.user_membership.add", "user.lifecycle.create"}, "EVENT_TYPE")),
					),
				),
			},
		},
	})
}

func TestAccResourceOktaEventHook_withNewFilter(t *testing.T) {
	resourceName := "okta_event_hook.test"
	mgr := newFixtureManager("resources", resources.OktaIDaaSEventHook, t.Name())
	config := mgr.GetFixtures("basic_with_new_filter.tf", t)
	updatedConfig := mgr.GetFixtures("basic_with_new_filter_updated.tf", t)

	acctest.OktaResourceTest(t, resource.TestCase{
		PreCheck:                 acctest.AccPreCheck(t),
		ErrorCheck:               testAccErrorChecks(t),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactoriesForTestAcc(t),
		CheckDestroy:             checkResourceDestroy(resources.OktaIDaaSEventHook, eventHookExists),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, eventHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),
					resource.TestCheckResourceAttr(resourceName, "filter.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "filter.*", map[string]string{
						"event":     "user.lifecycle.create",
						"condition": "event.actor.id != null",
					}),
					testCheckResourceSetAttr(
						resourceName,
						"events",
						idaas.EventSet(v5okta.NewEventSubscriptions([]string{"user.lifecycle.create"}, "EVENT_TYPE")),
					),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					ensureResourceExists(resourceName, eventHookExists),
					resource.TestCheckResourceAttr(resourceName, "name", acctest.BuildResourceName(mgr.Seed)),
					resource.TestCheckResourceAttr(resourceName, "status", idaas.StatusActive),
					resource.TestCheckResourceAttr(resourceName, "channel.type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "channel.version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "channel.uri", "https://example.com/test-updated"),
					resource.TestCheckResourceAttr(resourceName, "auth.type", "HEADER"),
					resource.TestCheckResourceAttr(resourceName, "auth.key", "Authorization"),
					resource.TestCheckResourceAttr(resourceName, "filter.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "filter.*", map[string]string{
						"event":     "user.lifecycle.create",
						"condition": "event.actor.id != null",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "filter.*", map[string]string{
						"event":     "user.lifecycle.delete.initiated",
						"condition": "event.actor.id != null",
					}),
					testCheckResourceSetAttr(
						resourceName,
						"events",
						idaas.EventSet(v5okta.NewEventSubscriptions([]string{"user.lifecycle.create", "user.lifecycle.delete.initiated"}, "EVENT_TYPE")),
					),
				),
			},
		},
	})
}

func eventHookExists(id string) (bool, error) {
	client := iDaaSAPIClientForTestUtil.OktaSDKClientV5()
	eh, resp, err := client.EventHookAPI.GetEventHook(context.Background(), id).Execute()
	if err := utils.SuppressErrorOn404_V5(resp, err); err != nil {
		return false, err
	}
	return eh != nil, nil
}

func testMakeEventHookHeadersSet(headers []*v5okta.EventHookChannelConfigHeader) *schema.Set {
	h := make([]interface{}, len(headers))
	for i, header := range headers {
		h[i] = map[string]interface{}{
			"key":   header.GetKey(),
			"value": header.GetValue(),
		}
	}
	return schema.NewSet(schema.HashResource(idaas.HeaderSchema), h)
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
