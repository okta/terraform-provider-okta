package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserBaseSchemaProperty() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserBaseSchemaCreate,
		ReadContext:   resourceUserBaseSchemaRead,
		UpdateContext: resourceUserBaseSchemaCreate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				resourceIndex := d.Id()
				resourceUserType := "default"
				if strings.Contains(d.Id(), ".") {
					resourceUserType = strings.Split(d.Id(), ".")[0]
					resourceIndex = strings.Split(d.Id(), ".")[1]
				}
				d.SetId(resourceIndex)
				_ = d.Set("index", resourceIndex)
				_ = d.Set("user_type", resourceUserType)
				return []*schema.ResourceData{d}, nil
			},
		},
		Description:   "Manages a User Base Schema property. This resource allows you to configure a base user schema property.",
		SchemaVersion: 1,
		Schema: buildSchema(
			userBaseSchemaSchema,
			userTypeSchema,
			userPatternSchema,
			map[string]*schema.Schema{
				"master": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Master priority for the user schema property. It can be set to `PROFILE_MASTER` or `OKTA`. Default: `PROFILE_MASTER`",
					Default:     "PROFILE_MASTER",
				},
			},
		),
		StateUpgraders: []schema.StateUpgrader{
			{
				Type: resourceUserBaseSchemaResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
					rawState["user_type"] = "default"
					return rawState, nil
				},
				Version: 0,
			},
		},
	}
}

func resourceUserBaseSchemaResourceV0() *schema.Resource {
	return &schema.Resource{Schema: userBaseSchemaSchema}
}

func resourceUserBaseSchemaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// NOTE: Okta API will ignore parallel calls to `POST
	// /api/v1/meta/schemas/user/{userId}` so a mutex to affect TF
	// `-parallelism=1` behavior is needed here.
	oktaMutexKV.Lock(userBaseSchemaProperty)
	defer oktaMutexKV.Unlock(userBaseSchemaProperty)

	if err := updateUserBaseSubschema(ctx, d, meta); err != nil {
		return err
	}
	d.SetId(d.Get("index").(string))
	return resourceUserBaseSchemaRead(ctx, d, meta)
}

func resourceUserBaseSchemaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	typeSchemaID, err := getUserTypeSchemaID(ctx, getOktaClientFromMetadata(meta), d.Get("user_type").(string))
	if err != nil {
		return diag.Errorf("failed to get user base schema: %v", err)
	}
	us, _, err := getOktaClientFromMetadata(meta).UserSchema.GetUserSchema(ctx, typeSchemaID)
	if err != nil {
		return diag.Errorf("failed to get user base schema: %v", err)
	}
	subschema := userSchemaBaseAttribute(us, d.Id())
	if subschema == nil {
		d.SetId("")
		return nil
	}
	syncBaseUserSchema(d, subschema)
	return nil
}

// create or modify a subschema
func updateUserBaseSubschema(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := validateUserBaseSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}
	schemaID, err := getUserTypeSchemaID(ctx, getOktaClientFromMetadata(meta), d.Get("user_type").(string))
	if err != nil {
		return diag.Errorf("failed to create user base schema: %v", err)
	}
	base := buildBaseUserSchema(d)
	url := fmt.Sprintf("/api/v1/meta/schemas/user/%v", schemaID)
	re := getOktaClientFromMetadata(meta).GetRequestExecutor()
	req, err := re.WithAccept("application/json").WithContentType("application/json").
		NewRequest("POST", url, base)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = re.Do(ctx, req, nil)
	if err != nil {
		return diag.Errorf("failed to create user base schema: %v", err)
	}
	return nil
}

func validateUserBaseSchema(d *schema.ResourceData) error {
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
