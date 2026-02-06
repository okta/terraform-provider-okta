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
		Schema: utils.BuildSchema(
			userSchemaSchema,
			userBaseSchemaSchema,
			userTypeSchema,
			map[string]*schema.Schema{
				"app_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The Application's ID the user schema property is associated with.",
				},
				"union": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "If `type` is `array`, determines whether attribute value is determined by group priority (false) or combine values across groups (true).",
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
			}),
	}
}

func dataSourceAppUserSchemaPropertyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appId := d.Get("app_id").(string)
	index := d.Get("index").(string)

	us, resp, err := getOktaClientFromMetadata(meta).UserSchema.GetApplicationUserSchema(ctx, appId)
	if err != nil {
		if utils.SuppressErrorOn404(resp, err) == nil {
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
