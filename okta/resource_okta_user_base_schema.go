package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

const baseSchema = "base"

func resourceUserBaseSchema() *schema.Resource {
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
					ValidateDiagFunc: stringInSlice([]string{"PROFILE_MASTER", "OKTA", "OVERRIDE", ""}),
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
	schemaUrl, err := getUserTypeSchemaUrl(ctx, getOktaClientFromMetadata(m), d.Get("user_type").(string))
	if err != nil {
		return diag.Errorf("failed to create user base schema: %v", err)
	}
	if err := updateBaseSubschema(ctx, getSupplementFromMetadata(m), schemaUrl, d); err != nil {
		return diag.Errorf("failed to create user base schema: %v", err)
	}
	d.SetId(d.Get("index").(string))
	return resourceUserBaseSchemaRead(ctx, d, m)
}

func resourceUserBaseSchemaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	schemaUrl, err := getUserTypeSchemaUrl(ctx, getOktaClientFromMetadata(m), d.Get("user_type").(string))
	if err != nil {
		return diag.Errorf("failed to get user base schema: %v", err)
	}
	subschema, err := getBaseSubSchema(ctx, getSupplementFromMetadata(m), schemaUrl, d)
	if err != nil {
		return diag.Errorf("failed to get user base schema: %v", err)
	}
	if subschema == nil {
		d.SetId("")
		return nil
	}
	syncBaseUserSchema(d, subschema)
	return nil
}

func getBaseSubSchema(ctx context.Context, client *sdk.ApiSupplement, schemaUrl string, d *schema.ResourceData) (*sdk.UserSubSchema, error) {
	s, _, err := client.GetUserSchema(ctx, schemaUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get user schema: %v", err)
	}
	return getBaseProperty(s, d.Id()), err
}

func resourceUserBaseSchemaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateUserBaseSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}
	schemaUrl, err := getUserTypeSchemaUrl(ctx, getOktaClientFromMetadata(m), d.Get("user_type").(string))
	if err != nil {
		return diag.Errorf("failed to update user base schema: %v", err)
	}
	if err := updateBaseSubschema(ctx, getSupplementFromMetadata(m), schemaUrl, d); err != nil {
		return diag.Errorf("failed to update user base schema: %v", err)
	}
	return resourceUserBaseSchemaRead(ctx, d, m)
}

// can't delete Base schema
func resourceUserBaseSchemaDelete(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}

// create or modify a subschema
func updateBaseSubschema(ctx context.Context, client *sdk.ApiSupplement, schemaUrl string, d *schema.ResourceData) error {
	_, _, err := client.UpdateBaseUserSchemaProperty(
		ctx,
		schemaUrl,
		d.Get("index").(string),
		userBasedSubSchema(d),
	)
	if err != nil {
		return fmt.Errorf("failed to update base user schema property: %v", err)
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
