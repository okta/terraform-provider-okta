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
		UpdateContext: resourceUserBaseSchemaUpdate,
		DeleteContext: resourceUserBaseSchemaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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
		SchemaVersion: 1,
		Schema: buildSchema(
			userBaseSchemaSchema,
			userTypeSchema,
			userPatternSchema,
			map[string]*schema.Schema{
				"master": {
					Type:     schema.TypeString,
					Optional: true,
					// Accepting an empty value to allow for zero value (when provisioning is off)
					ValidateDiagFunc: elemInSlice([]string{"PROFILE_MASTER", "OKTA", "OVERRIDE", ""}),
					Description:      "SubSchema profile manager, if not set it will inherit its setting.",
					Default:          "PROFILE_MASTER",
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

func resourceUserBaseSchemaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateUserBaseSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}
	typeSchemaID, err := getUserTypeSchemaID(ctx, getOktaClientFromMetadata(m), d.Get("user_type").(string))
	if err != nil {
		return diag.Errorf("failed to create user base schema: %v", err)
	}
	base := buildBaseUserSchema(d.Get("index").(string), buildUserBaseSchemaAttribute(d))
	_, _, err = getOktaClientFromMetadata(m).UserSchema.UpdateUserProfile(ctx, typeSchemaID, *base)
	if err != nil {
		return diag.Errorf("failed to create user base schema: %v", err)
	}
	d.SetId(d.Get("index").(string))
	return resourceUserBaseSchemaRead(ctx, d, m)
}

func resourceUserBaseSchemaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	typeSchemaID, err := getUserTypeSchemaID(ctx, getOktaClientFromMetadata(m), d.Get("user_type").(string))
	if err != nil {
		return diag.Errorf("failed to get user base schema: %v", err)
	}
	us, _, err := getOktaClientFromMetadata(m).UserSchema.GetUserSchema(ctx, typeSchemaID)
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

func resourceUserBaseSchemaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateUserBaseSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}
	typeSchemaID, err := getUserTypeSchemaID(ctx, getOktaClientFromMetadata(m), d.Get("user_type").(string))
	if err != nil {
		return diag.Errorf("failed to update user base schema: %v", err)
	}
	base := buildBaseUserSchema(d.Get("index").(string), buildUserBaseSchemaAttribute(d))
	_, _, err = getOktaClientFromMetadata(m).UserSchema.UpdateUserProfile(ctx, typeSchemaID, *base)
	if err != nil {
		return diag.Errorf("failed to create user base schema: %v", err)
	}
	return resourceUserBaseSchemaRead(ctx, d, m)
}

// can't delete Base schema
func resourceUserBaseSchemaDelete(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
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
