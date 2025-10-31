package idaas

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func resourceAdminRoleCustomAssignments() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdminRoleCustomAssignmentsCreate,
		ReadContext:   resourceAdminRoleCustomAssignmentsRead,
		UpdateContext: resourceAdminRoleCustomAssignmentsUpdate,
		DeleteContext: resourceAdminRoleCustomAssignmentsDelete,
		Importer:      utils.CreateNestedResourceImporter([]string{"resource_set_id", "custom_role_id"}),
		Description: `Resource to manage the assignment and unassignment of Custom Roles
These operations allow the creation and manipulation of custom roles as custom collections of permissions.
		
~> **NOTE:** This an Early Access feature.`,
		Schema: map[string]*schema.Schema{
			"resource_set_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the target Resource Set",
			},
			"custom_role_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the Custom Role",
			},
			"members": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The hrefs that point to User(s) and/or Group(s) that receive the Role",
			},
		},
	}
}

func resourceAdminRoleCustomAssignmentsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cr, err := buildAdminRoleCustomAssignment(d)
	if err != nil {
		return diag.Errorf("failed to create custom admin role assignment: %v", err)
	}
	client := getOktaV5ClientFromMetadata(meta)
	_, _, err = client.ResourceSetAPI.CreateResourceSetBinding(ctx, d.Get("resource_set_id").(string)).Instance(*cr).Execute()
	if err != nil {
		return diag.Errorf("failed to create custom admin role assignment: %v", err)
	}
	d.SetId(fmt.Sprintf("%s/%s", d.Get("resource_set_id").(string), d.Get("custom_role_id").(string)))
	return resourceAdminRoleCustomAssignmentsRead(ctx, d, meta)
}

func resourceAdminRoleCustomAssignmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(meta)
	members, _, err := client.ResourceSetAPI.ListMembersOfBinding(ctx, d.Get("resource_set_id").(string), d.Get("custom_role_id").(string)).Execute()
	if err != nil {
		return diag.Errorf("failed to list custom admin role assignment: %v", err)
	}
	if members == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("members", flattenAdminRoleCustomAssignments(members.Members))
	return nil
}

func resourceAdminRoleCustomAssignmentsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChange("members") {
		return nil
	}
	client := getAPISupplementFromMetadata(meta)
	oldMembers, newMembers := d.GetChange("members")
	oldSet := oldMembers.(*schema.Set)
	newSet := newMembers.(*schema.Set)
	membersToAdd := utils.ConvertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	membersToRemove := utils.ConvertInterfaceArrToStringArr(oldSet.Difference(newSet).List())
	err := assignMembersToCustomAdminRole(ctx, client, d.Get("resource_set_id").(string), d.Get("custom_role_id").(string), membersToAdd)
	if err != nil {
		return diag.FromErr(err)
	}
	err = removeMembersFromCustomAdminRole(ctx, client, d.Get("resource_set_id").(string), d.Get("custom_role_id").(string), membersToRemove)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAdminRoleCustomAssignmentsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(meta)
	members, _, _ := client.ResourceSetAPI.ListMembersOfBinding(ctx, d.Get("resource_set_id").(string), d.Get("custom_role_id").(string)).Execute()
	existingMembers := d.Get("members").(*schema.Set).List()
	for _, member := range members.Members {
		mem := member.Links.Self.GetHref()
		for _, v := range existingMembers {
			if mem == v {
				_, err := client.ResourceSetAPI.UnassignMemberFromBinding(ctx, d.Get("resource_set_id").(string), d.Get("custom_role_id").(string), member.GetId()).Execute()
				if err != nil {
					return diag.Errorf("failed to unassign member with id %s from binding with error: %v", member.GetId(), err)
				}
			}
		}
	}
	return nil
}

func buildAdminRoleCustomAssignment(d *schema.ResourceData) (*okta.ResourceSetBindingCreateRequest, error) {
	customRoleId := d.Get("custom_role_id").(string)
	rb := &okta.ResourceSetBindingCreateRequest{
		Role:    &customRoleId,
		Members: utils.ConvertInterfaceToStringSetNullable(d.Get("members")),
	}
	if len(rb.Members) == 0 {
		return nil, errors.New("at least one member must be specified when creating assignment")
	}
	return rb, nil
}

func flattenAdminRoleCustomAssignments(members []okta.ResourceSetBindingMember) *schema.Set {
	var arr []interface{}
	for _, member := range members {
		// Extract the URL from the links structure
		link := member.Links.Self.GetHref()
		arr = append(arr, link)
	}

	return schema.NewSet(schema.HashString, arr)
}

func listResourceSetBindingMembers(ctx context.Context, client *sdk.APISupplement, resourceSetID, customRoleID string) ([]*sdk.CustomRoleBindingMember, *sdk.Response, error) {
	var resMembers []*sdk.CustomRoleBindingMember
	resources, resp, err := client.ListResourceSetBindingMembers(ctx, resourceSetID, customRoleID, &query.Params{Limit: utils.DefaultPaginationLimit})
	if err != nil {
		return nil, resp, err
	}
	for {
		resMembers = append(resMembers, resources.Members...)
		if resp.HasNextPage() {
			resp, err = resp.Next(ctx, &resources)
			if err != nil {
				return nil, resp, err
			}
			continue
		} else {
			break
		}
	}
	return resMembers, nil, nil
}

func assignMembersToCustomAdminRole(ctx context.Context, client *sdk.APISupplement, resourceSetID, roleID string, links []string) error {
	if len(links) == 0 {
		return nil
	}
	_, err := client.AddResourceSetBindingMembers(ctx, resourceSetID, roleID, sdk.AddCustomRoleBindingMemberRequest{Additions: links})
	if err != nil {
		return fmt.Errorf("failed assign new members to the custom role: %v", err)
	}
	return nil
}

func removeMembersFromCustomAdminRole(ctx context.Context, client *sdk.APISupplement, resourceSetID, roleID string, urls []string) error {
	members, _, err := listResourceSetBindingMembers(ctx, client, resourceSetID, roleID)
	if err != nil {
		return fmt.Errorf("failed to list members assigned to the custom role: %v", err)
	}
	for _, member := range members {
		links := member.Links.(map[string]interface{})
		var url string
		for _, v := range links {
			for _, link := range v.(map[string]interface{}) {
				url = link.(string)
				break
			}
		}
		if utils.Contains(urls, url) {
			_, err := client.DeleteResourceSetBindingMember(ctx, resourceSetID, roleID, member.Id)
			if err != nil {
				return fmt.Errorf("failed to remove %s member from the custom role: %v", member.Id, err)
			}
		}
	}
	return nil
}
