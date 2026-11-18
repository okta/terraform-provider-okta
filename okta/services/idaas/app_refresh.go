package idaas

import (
	"context"
	"fmt"

	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/sdk"
)

func readAppAssignment(ctx context.Context, c *config.Config, app, group string) (*sdk.ApplicationGroupAssignment, *sdk.Response, error) {
	client := c.OktaIDaaSClient.OktaSDKClientV2()
	if !c.AppReadBatching || c.AppReads.Written() {
		return client.Application.GetApplicationGroupAssignment(ctx, app, group, nil)
	}
	value, found, err := c.AppReads.Get(ctx, "assignments/"+app, group, func(ctx context.Context) (map[string]any, error) {
		groups, _, err := listApplicationGroupAssignments(ctx, client, app)
		if err != nil {
			return nil, err
		}
		values := make(map[string]any, len(groups))
		for _, g := range groups {
			if g == nil || g.Id == "" {
				return nil, fmt.Errorf("app assignment collection contains a missing ID")
			}
			values[g.Id] = g
		}
		return values, nil
	})
	if err != nil {
		return nil, nil, err
	}
	if found && !c.AppReads.Written() {
		return value.(*sdk.ApplicationGroupAssignment), nil, nil
	}
	// List omission is not proof of deletion (pagination and eventual consistency).
	return client.Application.GetApplicationGroupAssignment(ctx, app, group, nil)
}

func readPushMapping(ctx context.Context, c *config.Config, app, id string) (*v6okta.GroupPushMapping, error) {
	api := c.OktaIDaaSClient.OktaSDKClientV6().GroupPushMappingAPI
	if !c.AppReadBatching || c.AppReads.Written() {
		mapping, _, err := api.GetGroupPushMapping(ctx, app, id).Execute()
		return mapping, err
	}
	value, found, err := c.AppReads.Get(ctx, "push/"+app, id, func(ctx context.Context) (map[string]any, error) {
		values := make(map[string]any)
		cursor := ""
		seen := make(map[string]bool)
		for {
			mappings, response, err := api.ListGroupPushMappings(ctx, app).After(cursor).Execute()
			if err != nil {
				return nil, err
			}
			for i := range mappings {
				mapping := mappings[i]
				if mapping.Id == nil || *mapping.Id == "" {
					return nil, fmt.Errorf("push mapping collection contains a missing ID")
				}
				values[*mapping.Id] = &mapping
			}
			if response == nil {
				return nil, fmt.Errorf("push mapping collection has no pagination response")
			}
			if !response.HasNextPage() {
				break
			}
			if len(mappings) == 0 {
				return nil, fmt.Errorf("push mapping collection has an empty page with a next link")
			}
			cursor = mappings[len(mappings)-1].GetId()
			if seen[cursor] {
				return nil, fmt.Errorf("push mapping collection repeated a pagination cursor")
			}
			seen[cursor] = true
		}
		return values, nil
	})
	if err != nil {
		return nil, err
	}
	if found && !c.AppReads.Written() {
		return value.(*v6okta.GroupPushMapping), nil
	}
	mapping, _, err := api.GetGroupPushMapping(ctx, app, id).Execute()
	return mapping, err
}
