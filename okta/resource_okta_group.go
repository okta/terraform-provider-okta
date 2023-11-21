package okta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				importID := strings.Split(d.Id(), "/")
				if len(importID) == 1 {
					return []*schema.ResourceData{d}, nil
				}
				if len(importID) > 2 {
					return nil, errors.New("invalid format used for import ID, format must be 'group_id' or 'group_id/skip_users'")
				}
				d.SetId(importID[0])
				_ = d.Set(importID[1], true)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group description",
			},
			"skip_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Ignore users sync. This is a temporary solution until 'users' field is supported in all the app-like resources",
				Default:     false,
				Deprecated:  "Because users has been removed, this attribute is a no op and will be removed",
			},
			"custom_profile_attributes": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        normalizeDataJSON,
				Description:      "JSON formatted custom attributes for a group. It must be JSON due to various types Okta allows.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "" {
						return true
					}

					var oldCustomAttrs sdk.GroupProfileMap
					_ = json.Unmarshal([]byte(old), &oldCustomAttrs)
					oldCustomAttrs = normalizeGroupProfile(oldCustomAttrs)

					var newCustomAttrs sdk.GroupProfileMap
					_ = json.Unmarshal([]byte(new), &newCustomAttrs)
					newCustomAttrs = normalizeGroupProfile(newCustomAttrs)

					return reflect.DeepEqual(oldCustomAttrs, newCustomAttrs)
				},
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("creating group", "name", d.Get("name").(string))
	group := buildGroup(d)
	responseGroup, _, err := getOktaClientFromMetadata(m).Group.CreateGroup(ctx, *group)
	if err != nil {
		return diag.Errorf("failed to create group: %v", err)
	}
	boc := newExponentialBackOffWithContext(ctx, 10*time.Second)
	err = backoff.Retry(func() error {
		g, resp, err := getOktaClientFromMetadata(m).Group.GetGroup(ctx, responseGroup.Id)
		if doNotRetry(m, err) {
			return backoff.Permanent(err)
		}
		if err := suppressErrorOn404(resp, err); err != nil {
			return backoff.Permanent(err)
		}
		if g == nil {
			return fmt.Errorf("group '%s' hasn't been created after multiple checks", responseGroup.Id)
		}
		return nil
	}, boc)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(responseGroup.Id)

	return resourceGroupRead(ctx, d, m)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("reading group", "id", d.Id(), "name", d.Get("name").(string))
	g, resp, err := getOktaClientFromMetadata(m).Group.GetGroup(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get group: %v", err)
	}

	if g == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", g.Profile.Name)
	_ = d.Set("description", g.Profile.Description)

	if g.Profile.GroupProfileMap != nil {
		customProfile, err := json.Marshal(g.Profile.GroupProfileMap)
		if err != nil {
			return diag.Errorf("failed to read custom profile attributes from group: %s", g.Profile.Name)
		}
		customProfileStr := string(customProfile)
		_ = d.Set("custom_profile_attributes", customProfileStr)
	}

	return nil
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("updating group", "id", d.Id(), "name", d.Get("name").(string))
	group := buildGroup(d)
	_, _, err := getOktaClientFromMetadata(m).Group.UpdateGroup(ctx, d.Id(), *group)
	if err != nil {
		return diag.Errorf("failed to update group: %v", err)
	}
	return resourceGroupRead(ctx, d, m)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("deleting group", "id", d.Id(), "name", d.Get("name").(string))
	_, err := getOktaClientFromMetadata(m).Group.DeleteGroup(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete group: %v", err)
	}
	return nil
}

func buildGroup(d *schema.ResourceData) *sdk.Group {
	var customAttrs sdk.GroupProfileMap
	if rawAttrs, ok := d.GetOk("custom_profile_attributes"); ok {
		str := rawAttrs.(string)

		// We validate the JSON, no need to check error
		_ = json.Unmarshal([]byte(str), &customAttrs)
	}

	return &sdk.Group{
		Profile: &sdk.GroupProfile{
			Name:            d.Get("name").(string),
			Description:     d.Get("description").(string),
			GroupProfileMap: normalizeGroupProfile(customAttrs),
		},
	}
}
