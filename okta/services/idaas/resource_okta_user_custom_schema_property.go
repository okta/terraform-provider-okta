package idaas

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceUserCustomSchemaProperty() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserSchemaCreateOrUpdate,
		ReadContext:   resourceUserSchemaRead,
		UpdateContext: resourceUserSchemaCreateOrUpdate,
		DeleteContext: resourceUserSchemaDelete,
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
		Description: "Creates a User Schema property. This resource allows you to create and configure a custom user schema property.",
		Schema: utils.BuildSchema(
			userBaseSchemaSchema,
			userSchemaSchema,
			userTypeSchema,
			userPatternSchema,
			map[string]*schema.Schema{
				"scope": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "NONE",
					Description: "determines whether an app user attribute can be set at the Individual or Group Level. Default: `NONE`",
				},
				"master": {
					Type:     schema.TypeString,
					Optional: true,
					// Accepting an empty value to allow for zero value (when provisioning is off)
					Description: " Master priority for the user schema property. It can be set to `PROFILE_MASTER`, `OVERRIDE` or `OKTA`.",
					Default:     "PROFILE_MASTER",
				},
				"master_override_priority": {
					Type:     schema.TypeList,
					Optional: true,
					Description: `Prioritized list of profile sources (required when 'master' is 'OVERRIDE').
	- 'type' - (Optional) - Type of profile source.
	- 'value' - (Required) - ID of profile source.`,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"type": {
								Type:     schema.TypeString,
								Optional: true,
								Default:  "APP",
							},
							"value": {
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
			},
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
	return &schema.Resource{Schema: utils.BuildSchema(userBaseSchemaSchema, userSchemaSchema, map[string]*schema.Schema{
		"scope": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "NONE",
		},
	})}
}

func resourceUserSchemaCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("creating user custom schema property", "name", d.Get("index").(string))
	err := validateUserSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}
	userCustomSchemaAttribute, err := buildUserCustomSchemaAttribute(d)
	if err != nil {
		return diag.FromErr(err)
	}
	custom := BuildCustomUserSchema(d.Get("index").(string), userCustomSchemaAttribute)
	subSchema, err := alterCustomUserSchema(ctx, meta, d.Get("user_type").(string), d.Get("index").(string), custom, false)
	if err != nil {
		return diag.Errorf("failed to create or update user custom schema property %s: %v", d.Get("index").(string), err)
	}
	d.SetId(d.Get("index").(string))
	err = syncCustomUserSchema(d, subSchema)
	if err != nil {
		return diag.Errorf("failed to set user custom schema property: %v", err)
	}
	return nil
}

func resourceUserSchemaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("reading user custom schema property", "name", d.Get("index").(string))
	typeSchemaID, err := GetUserTypeSchemaID(ctx, getOktaClientFromMetadata(meta), d.Get("user_type").(string))
	if err != nil {
		return diag.Errorf("failed to get user custom schema property: %v", err)
	}
	s, _, err := getOktaClientFromMetadata(meta).UserSchema.GetUserSchema(ctx, typeSchemaID)
	if err != nil {
		return diag.Errorf("failed to get user custom schema property: %v", err)
	}
	customAttribute := UserSchemaCustomAttribute(s, d.Id())
	if customAttribute == nil {
		d.SetId("")
		return nil
	}
	err = syncCustomUserSchema(d, customAttribute)
	if err != nil {
		return diag.Errorf("failed to set user custom schema property: %v", err)
	}
	return nil
}

func alterCustomUserSchema(ctx context.Context, meta interface{}, userType, index string, schema *sdk.UserSchema, isDeleteOperation bool) (*sdk.UserSchemaAttribute, error) {
	typeSchemaID, err := GetUserTypeSchemaID(ctx, getOktaClientFromMetadata(meta), userType)
	if err != nil {
		return nil, err
	}
	var schemaAttribute *sdk.UserSchemaAttribute

	boc := utils.NewExponentialBackOffWithContext(ctx, 120*time.Second)
	err = backoff.Retry(func() error {
		// NOTE: Enums on the schema can be typed other than string but the
		// Terraform SDK is staticly defined at runtime for string so we need to
		// juggle types on the fly.

		retypeUserSchemaPropertyEnums(schema)
		updated, resp, err := getOktaClientFromMetadata(meta).UserSchema.UpdateUserProfile(ctx, typeSchemaID, *schema)
		stringifyUserSchemaPropertyEnums(schema)
		if doNotRetry(meta, err) {
			return backoff.Permanent(err)
		}

		if err != nil {
			if resp != nil && resp.StatusCode == 500 {
				return fmt.Errorf("updating user custom schema property caused 500 error: %w", err)
			}
			if strings.Contains(err.Error(), "Wait until the data clean up process finishes and then try again") {
				return err
			}
			return backoff.Permanent(err)
		}
		s, _, err := getOktaClientFromMetadata(meta).UserSchema.GetUserSchema(ctx, typeSchemaID)
		if err != nil {
			return backoff.Permanent(fmt.Errorf("failed to get user custom schema property: %v", err))
		}
		schemaAttribute = UserSchemaCustomAttribute(s, index)
		if isDeleteOperation && schemaAttribute == nil {
			return nil
		} else if schemaAttribute != nil && reflect.DeepEqual(schemaAttribute, updated.Definitions.Custom.Properties[index]) {
			return nil
		}
		return errors.New("failed to apply changes after several retries")
	}, boc)
	if err != nil {
		logger(meta).Error("failed to apply changes after several retries", err)
	}
	return schemaAttribute, err
}

func resourceUserSchemaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	custom := BuildCustomUserSchema(d.Id(), nil)
	_, err := alterCustomUserSchema(ctx, meta, d.Get("user_type").(string), d.Get("index").(string), custom, true)
	if err != nil {
		return diag.Errorf("failed to delete user schema property %s: %v", d.Get("index").(string), err)
	}
	return nil
}

func validateUserSchema(d *schema.ResourceData) error {
	v, ok := d.GetOk("master")
	if !ok || v.(string) != "OVERRIDE" {
		return nil
	}
	mop, _ := d.Get("master_override_priority").([]interface{})
	if len(mop) == 0 {
		return errors.New("when setting profile master type to 'OVERRIDE' at least one 'master_override_priority' should be provided")
	}
	return nil
}
