package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

const customSchema = "custom"

func resourceUserSchema() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserSchemaCreate,
		ReadContext:   resourceUserSchemaRead,
		UpdateContext: resourceUserSchemaUpdate,
		DeleteContext: resourceUserSchemaDelete,
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
		Schema:        buildBaseUserSchema(userSchemaSchema),
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type: resourceUserSchemaResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
					rawState["user_type"] = "default"
					return rawState, nil
				},
				Version: 0,
			},
		},
	}
}

func resourceUserSchemaResourceV0() *schema.Resource {
	return &schema.Resource{Schema: buildSchema(userBaseSchemaSchema, userSchemaSchema)}
}

func resourceUserSchemaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	schemaUrl, err := getUserTypeSchemaUrl(ctx, getOktaClientFromMetadata(m), d.Get("user_type").(string))
	if err != nil {
		return diag.Errorf("failed to create user custom schema: %v", err)
	}
	_, _, err = getSupplementFromMetadata(m).UpdateCustomUserSchemaProperty(ctx, schemaUrl, d.Get("index").(string), getUserSubSchema(d))
	if err != nil {
		return diag.Errorf("failed to create user custom schema: %v", err)
	}
	d.SetId(d.Get("index").(string))
	return resourceUserSchemaRead(ctx, d, m)
}

func resourceUserSchemaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	schemaUrl, err := getUserTypeSchemaUrl(ctx, getOktaClientFromMetadata(m), d.Get("user_type").(string))
	if err != nil {
		return diag.Errorf("failed to get user custom schema: %v", err)
	}
	subschema, err := getSubSchema(ctx, getSupplementFromMetadata(m), schemaUrl, d)
	if err != nil {
		return diag.Errorf("failed to get user custom schema: %v", err)
	}
	if subschema == nil {
		d.SetId("")
		return nil
	}
	err = syncUserSchema(d, subschema)
	if err != nil {
		return diag.Errorf("failed to set user custom schema properties: %v", err)
	}
	return nil
}

func getSubSchema(ctx context.Context, client *sdk.ApiSupplement, schemaUrl string, d *schema.ResourceData) (*sdk.UserSubSchema, error) {
	s, _, err := client.GetUserSchema(ctx, schemaUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get user custom schema: %v", err)
	}
	return getCustomProperty(s, d.Id()), err
}

func resourceUserSchemaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	schemaUrl, err := getUserTypeSchemaUrl(ctx, getOktaClientFromMetadata(m), d.Get("user_type").(string))
	if err != nil {
		return diag.Errorf("failed to update user custom schema: %v", err)
	}
	_, _, err = getSupplementFromMetadata(m).UpdateCustomUserSchemaProperty(ctx, schemaUrl, d.Get("index").(string), getUserSubSchema(d))
	if err != nil {
		return diag.Errorf("failed to update user custom schema: %v", err)
	}
	return resourceUserSchemaRead(ctx, d, m)
}

func resourceUserSchemaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	schemaUrl, err := getUserTypeSchemaUrl(ctx, getOktaClientFromMetadata(m), d.Get("user_type").(string))
	if err != nil {
		return diag.Errorf("failed to delete user custom schema: %v", err)
	}
	_, err = getSupplementFromMetadata(m).DeleteUserSchemaProperty(ctx, schemaUrl, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete user custom schema: %v", err)
	}
	return nil
}
