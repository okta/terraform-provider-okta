package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func ResourceAppUserBaseSchemaProperty() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppUserBaseSchemaCreate,
		ReadContext:   resourceAppUserBaseSchemaRead,
		UpdateContext: resourceAppUserBaseSchemaUpdate,
		DeleteContext: utils.ResourceFuncNoOp,
		Importer:      utils.CreateNestedResourceImporter([]string{"app_id", "index"}),
		Description:   "Manages an Application User Base Schema property. This resource allows you to configure a base app user schema property.",
		Schema: utils.BuildSchema(
			userBaseSchemaSchema,
			userTypeSchema,
			userPatternSchema,
			map[string]*schema.Schema{
				"master": {
					Type:     schema.TypeString,
					Optional: true,
					// Accepting an empty value to allow for zero value (when provisioning is off)
					Description: "Master priority for the user schema property. It can be set to `PROFILE_MASTER` or `OKTA`. Default: `PROFILE_MASTER`",
					Default:     "PROFILE_MASTER",
				},
				"app_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The Application's ID the user schema property should be assigned to.",
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
	return &schema.Resource{Schema: utils.BuildSchema(map[string]*schema.Schema{
		"app_id": {
			Type:     schema.TypeString,
			Required: true,
		},
	}, userBaseSchemaSchema)}
}

func resourceAppUserBaseSchemaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := updateAppUserBaseSubschema(ctx, d, meta); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s/%s", d.Get("app_id").(string), d.Get("index").(string)))
	return resourceAppUserBaseSchemaRead(ctx, d, meta)
}

func resourceAppUserBaseSchemaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	us, _, err := GetOktaClientFromMetadata(meta).UserSchema.GetApplicationUserSchema(ctx, d.Get("app_id").(string))
	if err != nil {
		return diag.Errorf("failed to get app user base schema: %v", err)
	}
	subschema := UserSchemaBaseAttribute(us, d.Get("index").(string))
	if subschema == nil {
		d.SetId("")
		return nil
	}
	syncBaseUserSchema(d, subschema)
	return nil
}

func resourceAppUserBaseSchemaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := updateAppUserBaseSubschema(ctx, d, meta); err != nil {
		return err
	}
	return resourceAppUserBaseSchemaRead(ctx, d, meta)
}

// create or modify a subschema
func updateAppUserBaseSubschema(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := validateAppUserBaseSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}
	base := buildBaseUserSchema(d)
	url := fmt.Sprintf("/api/v1/meta/schemas/apps/%v/default", d.Get("app_id").(string))
	re := GetOktaClientFromMetadata(meta).GetRequestExecutor()
	req, err := re.WithAccept("application/json").WithContentType("application/json").
		NewRequest("POST", url, base)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = re.Do(ctx, req, nil)
	if err != nil {
		return diag.Errorf("failed to update application user base schema: %v", err)
	}
	return nil
}

func validateAppUserBaseSchema(d *schema.ResourceData) error {
	_, ok := d.GetOk("pattern")
	if d.Get("index").(string) != "login" {
		if ok {
			return fmt.Errorf("'pattern' property is only allowed to be set for 'login'")
		}
		return nil
	}
	if !d.Get("required").(bool) {
		return fmt.Errorf("'login' base schema is always required attribute")
	}
	return nil
}
