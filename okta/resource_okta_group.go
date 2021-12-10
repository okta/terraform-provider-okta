package okta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
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
				if !isValidSkipArg(importID[1]) {
					return nil, fmt.Errorf("'%s' is invalid value to be used as part of import ID, it can only be 'skip_users'", importID[1])
				}
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
			"users": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Users associated with the group. This can also be done per user.",
				Deprecated:  "The `users` field is now deprecated for the resource `okta_group`, please replace all uses of this with: `okta_group_memberships`",
			},
			"skip_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Ignore users sync. This is a temporary solution until 'users' field is supported in this resource",
				Default:     false,
			},
			"custom_profile_attributes": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        normalizeDataJSON,
				Description:      "JSON formatted custom attributes for a group. It must be JSON due to various types Okta allows.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
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
	bOff := backoff.NewExponentialBackOff()
	bOff.MaxElapsedTime = time.Second * 10
	bOff.InitialInterval = time.Second
	err = backoff.Retry(func() error {
		g, resp, err := getOktaClientFromMetadata(m).Group.GetGroup(ctx, responseGroup.Id)
		if err := suppressErrorOn404(resp, err); err != nil {
			return backoff.Permanent(err)
		}
		if g == nil {
			return fmt.Errorf("group '%s' hasn't been created after multiple checks", responseGroup.Id)
		}
		return nil
	}, bOff)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(responseGroup.Id)
	err = updateGroupUsers(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to update group users on group create: %v", err)
	}
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

	err = syncGroupUsers(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to get group users: %v", err)
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
	err = updateGroupUsers(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to update group users on group update: %v", err)
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

func syncGroupUsers(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	// Only sync when the user opts in by outlining users in the group config
	if _, exists := d.GetOk("users"); !exists {
		return nil
	}
	// temp solution until 'users' field is supported
	if d.Get("skip_users").(bool) {
		return nil
	}
	userIDList, err := listGroupUserIDs(ctx, m, d.Id())
	if err != nil {
		return err
	}
	return d.Set("users", convertStringSliceToSet(userIDList))
}

func updateGroupUsers(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	if !d.HasChange("users") {
		return nil
	}
	// temp solution until 'users' field is supported
	if d.Get("skip_users").(bool) {
		return nil
	}
	client := getOktaClientFromMetadata(m)
	oldGM, newGM := d.GetChange("users")
	oldSet := oldGM.(*schema.Set)
	newSet := newGM.(*schema.Set)
	usersToAdd := convertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	usersToRemove := convertInterfaceArrToStringArr(oldSet.Difference(newSet).List())
	err := addGroupMembers(ctx, client, d.Id(), usersToAdd)
	if err != nil {
		return err
	}
	return removeGroupMembers(ctx, client, d.Id(), usersToRemove)
}

func buildGroup(d *schema.ResourceData) *okta.Group {
	var customAttrs okta.GroupProfileMap
	if rawAttrs, ok := d.GetOk("custom_profile_attributes"); ok {
		str := rawAttrs.(string)

		// We validate the JSON, no need to check error
		_ = json.Unmarshal([]byte(str), &customAttrs)
	}

	return &okta.Group{
		Profile: &okta.GroupProfile{
			Name:            d.Get("name").(string),
			Description:     d.Get("description").(string),
			GroupProfileMap: customAttrs,
		},
	}
}
