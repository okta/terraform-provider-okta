package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func resourceAppUserBaseSchema() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppUserBaseSchemaCreate,
		ReadContext:   resourceAppUserBaseSchemaRead,
		UpdateContext: resourceAppUserBaseSchemaUpdate,
		DeleteContext: resourceAppUserBaseSchemaDelete,
		Importer:      createNestedResourceImporter([]string{"app_id", "id"}),
		Schema: buildBaseUserSchema(map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		}),
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type: resourceAppUserBaseSchemaResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
					rawState["user_type"] = "default"
					return rawState, nil
				},
				Version: 0,
			},
		},
	}
}

func resourceAppUserBaseSchemaResourceV0() *schema.Resource {
	return &schema.Resource{Schema: buildSchema(map[string]*schema.Schema{
		"app_id": {
			Type:     schema.TypeString,
			Required: true,
		},
	}, userBaseSchemaSchema)}
}

func resourceAppUserBaseSchemaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if err := updateAppUserBaseSubschema(ctx, d, m); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s/%s", d.Get("app_id").(string), d.Get("index").(string)))
	return resourceAppUserBaseSchemaRead(ctx, d, m)
}

func resourceAppUserBaseSchemaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	us, _, err := getSupplementFromMetadata(m).GetAppUserSchema(ctx, d.Get("app_id").(string))
	if err != nil {
		return diag.Errorf("failed to get app user base schema: %v", err)
	}
	subschema := getBaseProperty(us, d.Get("index").(string))
	if subschema == nil {
		d.SetId("")
		return nil
	}
	syncBaseUserSchema(d, subschema)
	return nil
}

func resourceAppUserBaseSchemaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if err := updateAppUserBaseSubschema(ctx, d, m); err != nil {
		return err
	}
	return resourceAppUserBaseSchemaRead(ctx, d, m)
}

// can't delete Base
func resourceAppUserBaseSchemaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// create or modify a subschema
func updateAppUserBaseSubschema(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	subSchema := &sdk.UserSubSchema{
		Master: getNullableMaster(d),
		Title:  d.Get("title").(string),
		Type:   d.Get("type").(string),
		Permissions: []*sdk.UserSchemaPermission{
			{
				Action:    d.Get("permissions").(string),
				Principal: "SELF",
			},
		},
		Required: boolPtr(d.Get("required").(bool)),
	}
	_, _, err := getSupplementFromMetadata(m).UpdateBaseAppUserSchemaProperty(
		ctx,
		d.Get("index").(string),
		d.Get("app_id").(string),
		subSchema,
	)
	if err != nil {
		return diag.Errorf("failed to update application user base schema: %v", err)
	}
	return nil
}
