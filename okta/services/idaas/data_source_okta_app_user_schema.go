package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func dataSourceAppUserSchema() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppUserSchemaRead,
		Description: `Gets the entire application user schema for an application.

This data source allows you to retrieve all custom properties in an application's user schema without managing them in Terraform.

~> **Note:** This is useful for referencing auto-created properties (from provisioning features like PUSH_NEW_USERS) or reading the complete schema configuration for use in other resources.`,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Application's ID",
			},
			"custom_property": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Custom properties in the schema",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The property name/index",
						},
						"title": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name for the property",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the schema property",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the property",
						},
						"required": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the property is required",
						},
						"scope": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scope of the property",
						},
						"min_length": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum length of the property value",
						},
						"max_length": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum length of the property value",
						},
						"enum": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Array of values a primitive property can be set to",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"external_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "External name of the property",
						},
						"external_namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "External namespace of the property",
						},
						"master": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Master priority for the property",
						},
						"permissions": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access control permissions for the property",
						},
						"union": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether attribute value is combined across groups",
						},
					},
				},
			},
		},
	}
}

func dataSourceAppUserSchemaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appId := d.Get("app_id").(string)
	us, resp, err := getOktaClientFromMetadata(meta).UserSchema.GetApplicationUserSchema(ctx, appId)
	if err != nil {
		if utils.SuppressErrorOn404(resp, err) == nil {
			return diag.Errorf("application with ID '%s' not found or has no user schema", appId)
		}
		return diag.Errorf("failed to get application user schema: %v", err)
	}

	d.SetId(appId)
	_ = d.Set("app_id", appId)

	// Read custom properties
	if us.Definitions != nil && us.Definitions.Custom != nil && us.Definitions.Custom.Properties != nil {
		customProps := make([]interface{}, 0)
		for index, attr := range us.Definitions.Custom.Properties {
			propMap := flattenSchemaProperty(index, attr)
			customProps = append(customProps, propMap)
		}
		_ = d.Set("custom_property", customProps)
	}

	return nil
}
