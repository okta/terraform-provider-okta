package okta

import (
	"context"
	"errors"
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
		Importer:      createNestedResourceImporter([]string{"app_id", "index"}),
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
				"union": {
					Type:          schema.TypeBool,
					Optional:      true,
					Description:   "Allows to assign attribute's group priority",
					Default:       false,
					ConflictsWith: []string{"enum"},
				},
				"scope": {
					Type:             schema.TypeString,
					Optional:         true,
					Default:          "NONE",
					ValidateDiagFunc: stringInSlice([]string{"SELF", "NONE", ""}),
					ForceNew:         true, // since the `scope` is read-only attribute, the resource should be recreated
				},
				"master": {
					Type:     schema.TypeString,
					Optional: true,
					// Accepting an empty value to allow for zero value (when provisioning is off)
					ValidateDiagFunc: stringInSlice([]string{"PROFILE_MASTER", "OKTA", ""}),
					Description:      "SubSchema profile manager, if not set it will inherit its setting.",
					Default:          "PROFILE_MASTER",
				},
			}),
		SchemaVersion: 2,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type: resourceAppUserSchemaResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
					rawState["user_type"] = "default"
					return rawState, nil
				},
				Version: 0,
			},
			{
				Type: resourceAppUserSchemaResourceV1().CoreConfigSchema().ImpliedType(),
				Upgrade: func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
					rawState["union"] = false
					return rawState, nil
				},
				Version: 1,
			},
		},
	}
}

func resourceAppUserSchemaResourceV1() *schema.Resource {
	return &schema.Resource{Schema: buildSchema(map[string]*schema.Schema{
		"app_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"scope": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "NONE",
			ValidateDiagFunc: stringInSlice([]string{"SELF", "NONE", ""}),
			ForceNew:         true, // since the `scope` is read-only attribute, the resource should be recreated
		},
	}, userSchemaSchema, userBaseSchemaSchema, userTypeSchema, userPatternSchema)}
}

func resourceAppUserSchemaResourceV0() *schema.Resource {
	return &schema.Resource{Schema: buildSchema(map[string]*schema.Schema{
		"app_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"scope": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "NONE",
			ValidateDiagFunc: stringInSlice([]string{"SELF", "NONE", ""}),
			ForceNew:         true, // since the `scope` is read-only attribute, the resource should be recreated
		},
	}, userSchemaSchema, userBaseSchemaSchema)}
}

func resourceAppUserSchemaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateAppUserSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}
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
	if subschema.Union != "" {
		if subschema.Union == "DISABLE" {
			_ = d.Set("union", false)
		} else {
			_ = d.Set("union", true)
		}
	}
	if err != nil {
		return diag.Errorf("failed to set user schema properties: %v", err)
	}
	return nil
}

func resourceAppUserSchemaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateAppUserSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}
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
	subSchema := userSubSchema(d)
	if d.Get("union").(bool) {
		subSchema.Union = "ENABLE"
	} else {
		subSchema.Union = "DISABLE"
	}
	_, _, err := getSupplementFromMetadata(m).UpdateCustomAppUserSchemaProperty(
		ctx,
		d.Get("index").(string),
		d.Get("app_id").(string),
		subSchema,
	)
	if err != nil {
		return diag.Errorf("failed to update custom app user schema property: %v", err)
	}
	return nil
}

func validateAppUserSchema(d *schema.ResourceData) error {
	if scope, ok := d.GetOk("scope"); ok {
		if union, ok := d.GetOk("union"); ok {
			if scope == "SELF" && union.(bool) {
				return errors.New("you can not use combine values across groups (union=true) for self scoped " +
					"attribute (scope=SELF). Either change scope to 'NONE', or use group priority option by setting union to 'false'")
			}
		}
	}
	return nil
}
