package okta

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

const customSchema = "custom"

func resourceUserSchema() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserSchemaCreateOrUpdate,
		ReadContext:   resourceUserSchemaRead,
		UpdateContext: resourceUserSchemaCreateOrUpdate,
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
		Schema: buildSchema(
			userBaseSchemaSchema,
			userSchemaSchema,
			userTypeSchema,
			userPatternSchema,
		),
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

// Sometime Okta API does not update or create custom property on the first try, thus that require running
// `terraform apply` several times. This simple retry resolves that issue. (If) After  this issue will be resolved,
// this retry logic will be demolished.
func resourceUserSchemaCreateOrUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("creating user schema", "name", d.Get("index").(string))
	schemaUrl, err := getUserTypeSchemaUrl(ctx, getOktaClientFromMetadata(m), d.Get("user_type").(string))
	if err != nil {
		return diag.Errorf("failed to create user custom schema: %v", err)
	}
	var subschema *sdk.UserSubSchema
	timer := time.NewTimer(time.Second * 3)
	ticker := time.NewTicker(time.Millisecond * 500)
loop:
	for {
		select {
		case <-ctx.Done():
			return diag.Errorf("failed to create user custom schema: %v", ctx.Err())
		case <-timer.C:
			return diag.Errorf("failed to create user custom schema: no more attempts left")
		case <-ticker.C:
			updated, _, err := getSupplementFromMetadata(m).UpdateCustomUserSchemaProperty(ctx, schemaUrl, d.Get("index").(string), userSubSchema(d))
			if err != nil {
				return diag.Errorf("failed to create user custom schema: %v", err)
			}
			d.SetId(d.Get("index").(string))
			s, _, err := getSupplementFromMetadata(m).GetUserSchema(ctx, schemaUrl)
			if err != nil {
				return diag.Errorf("failed to get user custom schema: %v", err)
			}
			subschema = getCustomProperty(s, d.Id())
			if subschema != nil && reflect.DeepEqual(subschema, updated.Definitions.Custom.Properties[d.Id()]) {
				break loop
			}
		}
	}
	err = syncUserSchema(d, subschema)
	if err != nil {
		return diag.Errorf("failed to set user custom schema properties: %v", err)
	}
	return nil
}

func resourceUserSchemaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("reading user schema", "name", d.Get("index").(string))
	schemaUrl, err := getUserTypeSchemaUrl(ctx, getOktaClientFromMetadata(m), d.Get("user_type").(string))
	if err != nil {
		return diag.Errorf("failed to get user custom schema: %v", err)
	}
	s, _, err := getSupplementFromMetadata(m).GetUserSchema(ctx, schemaUrl)
	if err != nil {
		return diag.Errorf("failed to get user custom schema: %v", err)
	}
	subschema := getCustomProperty(s, d.Id())
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
