package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func dataSourceAppUserSchemaProperty() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppUserSchemaPropertyRead,
		Description: `Gets an application user schema property.

This data source allows you to retrieve information about an existing application user schema property without managing it in Terraform.

~> **Note:** App user schema properties may be automatically created by Okta when provisioning features are enabled on an application (e.g., PUSH_NEW_USERS, PUSH_PROFILE_UPDATES). Common auto-created properties include userName, email, givenName, familyName, displayName, title, department, and manager. This data source is useful for referencing these properties without bringing them under Terraform management.`,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Application's ID the user schema property is associated with.",
			},
			"index": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subschema unique string identifier",
			},
			"title": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The property's display title.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the schema property. It can be `string`, `boolean`, `number`, `integer`, `array`, or `object`.",
			},
			"array_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the array elements if `type` is set to `array`.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the user schema property.",
			},
			"min_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The minimum length of the user property value. Only applies to type `string`.",
			},
			"max_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum length of the user property value. Only applies to type `string`.",
			},
			"enum": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of values a primitive property can be set to. See `array_enum` for arrays.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"one_of": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of maps containing a mapping for display name to enum value.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"const": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enum value",
						},
						"title": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enum title",
						},
					},
				},
			},
			"array_enum": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of values that an array property's items can be set to.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"array_one_of": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Display name and value an enum array can be set to.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"const": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value mapping to member of `array_enum`.",
						},
						"title": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name for the enum value.",
						},
					},
				},
			},
			"external_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External name of the user schema property.",
			},
			"external_namespace": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External namespace of the user schema property.",
			},
			"unique": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether the property should be unique. It can be set to `UNIQUE_VALIDATED` or `NOT_UNIQUE`.",
			},
			"permissions": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access control permissions for the property. It can be set to `READ_WRITE`, `READ_ONLY`, `HIDE`.",
			},
			"required": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the subschema is required.",
			},
			"scope": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Determines whether an app user attribute can be set at the Personal `SELF` or Group `NONE` level.",
			},
			"master": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Master priority for the user schema property. It can be `PROFILE_MASTER` or `OKTA`.",
			},
			"union": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If `type` is `array`, determines whether attribute value is determined by group priority (false) or combine values across groups (true).",
			},
		},
	}
}

func dataSourceAppUserSchemaPropertyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appId := d.Get("app_id").(string)
	index := d.Get("index").(string)

	us, resp, err := getOktaClientFromMetadata(meta).UserSchema.GetApplicationUserSchema(ctx, appId)
	if err != nil {
		if resp != nil && utils.SuppressErrorOn404(resp, err) == nil {
			return diag.Errorf("application with ID '%s' not found or has no user schema", appId)
		}
		return diag.Errorf("failed to get application user schema: %v", err)
	}

	subschema := UserSchemaCustomAttribute(us, index)
	if subschema == nil {
		return diag.Errorf(
			"application user schema property with index '%s' not found for app '%s'.\n\n"+
				"This property may not exist yet. If you're trying to reference a property that should "+
				"be auto-created when provisioning is enabled, make sure provisioning features are enabled "+
				"on the application first.",
			index, appId,
		)
	}

	d.SetId(fmt.Sprintf("%s/%s", appId, index))

	// Sync all properties
	if err := syncCustomUserSchema(d, subschema); err != nil {
		return diag.Errorf("failed to set application user schema properties: %v", err)
	}

	// Handle union attribute
	if subschema.Union != "" {
		if subschema.Union == "DISABLE" {
			_ = d.Set("union", false)
		} else {
			_ = d.Set("union", true)
		}
	}

	return nil
}
