package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppUserSchema() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppUserSchemaCreate,
		ReadContext:   resourceAppUserSchemaRead,
		UpdateContext: resourceAppUserSchemaUpdate,
		DeleteContext: resourceAppUserSchemaDelete,
		Importer:      createNestedResourceImporter([]string{"app_id", "id"}),
		Schema: buildSchema(
			userSchemaSchema,
			userBaseSchemaSchema,
			userTypeSchema,
			userPatternSchema,
			map[string]*schema.Schema{
				"app_id": {
					Type:     schema.TypeString,
					Required: true,
				},
			}),
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type: resourceAppUserSchemaResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
					rawState["user_type"] = "default"
					return rawState, nil
				},
				Version: 0,
			},
		},
	}
}

func resourceAppUserSchemaResourceV0() *schema.Resource {
	return &schema.Resource{Schema: buildSchema(map[string]*schema.Schema{
		"app_id": {
			Type:     schema.TypeString,
			Required: true,
		},
	}, userSchemaSchema, userBaseSchemaSchema)}
}

func resourceAppUserSchemaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if err := updateAppUserSubschema(ctx, d, m); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s/%s", d.Get("app_id").(string), d.Get("index").(string)))
	return resourceAppUserSchemaRead(ctx, d, m)
}

func resourceAppUserSchemaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	us, resp, err := getSupplementFromMetadata(m).GetAppUserSchema(ctx, d.Get("app_id").(string))
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get app user schema: %v", err)
	}
	subschema := getCustomProperty(us, d.Get("index").(string))
	if subschema == nil {
		d.SetId("")
		return nil
	}
	err = syncUserSchema(d, subschema)
	if err != nil {
		return diag.Errorf("failed to set user schema properties: %v", err)
	}
	return nil
}

func resourceAppUserSchemaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if err := updateAppUserSubschema(ctx, d, m); err != nil {
		return err
	}
	return resourceAppUserSchemaRead(ctx, d, m)
}

func resourceAppUserSchemaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, err := getSupplementFromMetadata(m).DeleteAppUserSchemaProperty(ctx, d.Get("index").(string), d.Get("app_id").(string))
	if err != nil {
		return diag.Errorf("failed to delete user schema property")
	}
	return nil
}

func updateAppUserSubschema(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, _, err := getSupplementFromMetadata(m).UpdateCustomAppUserSchemaProperty(
		ctx,
		d.Get("index").(string),
		d.Get("app_id").(string),
		userSubSchema(d),
	)
	if err != nil {
		return diag.Errorf("failed to update custom app user schema property: %v", err)
	}
	return nil
}
