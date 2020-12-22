package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
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
		Schema:        buildSchema(userBaseSchemaSchema, userTypeSchema),
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
func resourceUserBaseSchemaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// create or modify a subschema
func updateBaseSubschema(ctx context.Context, client *sdk.ApiSupplement, schemaUrl string, d *schema.ResourceData) error {
	subSchema := &sdk.UserSubSchema{
		Master: getNullableMaster(d),
		Title:  d.Get("title").(string),
		Type:   d.Get("type").(string),

		Pattern: d.Get("pattern").(string),
		Permissions: []*sdk.UserSchemaPermission{
			{
				Action:    d.Get("permissions").(string),
				Principal: "SELF",
			},
		},
		Required: boolPtr(d.Get("required").(bool)),
	}
	_, _, err := client.UpdateBaseUserSchemaProperty(ctx, schemaUrl, d.Get("index").(string), subSchema)
	if err != nil {
		return fmt.Errorf("failed to update base user schema property: %v", err)
	}
	return nil
}
