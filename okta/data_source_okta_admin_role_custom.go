package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAdminRoleCustom() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAdminRoleCustomRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"label"},
			},
			"label": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permissions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceAdminRoleCustomRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var identifier string
	if d.Get("id") != "" {
		identifier = d.Get("id").(string)
	} else if d.Get("label") != "" {
		identifier = d.Get("label").(string)
	} else {
		return diag.Errorf("you must provide either an 'id' or a 'label' to retrieve a custom admin role")
	}

	role, _, err := getOktaV3ClientFromMetadata(m).RoleApi.GetRole(ctx, identifier).Execute()
	if err != nil {
		return diag.Errorf("failed to find custom admin role: %v", err)
	}

	permissionsResp, _, err := getOktaV3ClientFromMetadata(m).RoleApi.ListRolePermissions(ctx, identifier).Execute()
	if err != nil {
		return diag.Errorf("failed to list permissions for custom admin role: %v", err)
	}

	permissions := make([]interface{}, len(permissionsResp.Permissions))
	for i := range permissionsResp.Permissions {
		permissions[i] = *permissionsResp.Permissions[i].Label
	}

	if *role.Id != d.Id() {
		d.SetId(*role.Id)
	}
	_ = d.Set("label", role.Label)
	_ = d.Set("description", role.Description)
	_ = d.Set("permissions", schema.NewSet(schema.HashString, permissions))

	return nil
}
