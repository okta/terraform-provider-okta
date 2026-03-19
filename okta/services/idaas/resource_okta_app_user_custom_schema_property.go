package idaas

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceAppUserSchemaProperty() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppUserSchemaPropertyCreate,
		ReadContext:   resourceAppUserSchemaPropertyRead,
		UpdateContext: resourceAppUserSchemaPropertyUpdate,
		DeleteContext: resourceAppUserSchemaPropertyDelete,
		Importer:      utils.CreateNestedResourceImporter([]string{"app_id", "index"}),
		Description: `Creates an Application User Schema property.
This resource allows you to create and configure a custom user schema property and associate it with an application.
Make sure that the app instance is 'active' before creating the schema property, because in some cases API might return '404' error.
**IMPORTANT:** With 'enum', list its values as strings even though the 'type'
may be something other than string. This is a limitation of the schema definition
in the Terraform Plugin SDK runtime and we juggle the type correctly when making
Okta API calls. Same holds for the 'const' value of 'one_of' as well as the
'array_*' variation of 'enum' and 'one_of'`,
		Schema: utils.BuildSchema(
			userSchemaSchema,
			userBaseSchemaSchema,
			userTypeSchema,
			// userPatternSchema,
			map[string]*schema.Schema{
				"app_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The Application's ID the user custom schema property should be assigned to.",
				},
				"union": {
					Type:          schema.TypeBool,
					Optional:      true,
					Description:   "If `type` is set to `array`, used to set whether attribute value is determined by group priority `false`, or combine values across groups `true`. Can not be set to `true` if `scope` is set to `SELF`.",
					Default:       false,
					ConflictsWith: []string{"enum"},
				},
				"scope": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "NONE",
					ForceNew:    true, // since the `scope` is read-only attribute, the resource should be recreated
					Description: "determines whether an app user attribute can be set at the Personal `SELF` or Group `NONE` level. Default value is `NONE`.",
				},
				"master": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Master priority for the user schema property. It can be set to `PROFILE_MASTER` or `OKTA`",
					Default:     "PROFILE_MASTER",
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
	return &schema.Resource{Schema: utils.BuildSchema(map[string]*schema.Schema{
		"app_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"scope": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "NONE",
			ForceNew: true, // since the `scope` is read-only attribute, the resource should be recreated
		},
	}, userSchemaSchema, userBaseSchemaSchema, userTypeSchema, userPatternSchema)}
}

func resourceAppUserSchemaResourceV0() *schema.Resource {
	return &schema.Resource{Schema: utils.BuildSchema(map[string]*schema.Schema{
		"app_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"scope": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "NONE",
			ForceNew: true, // since the `scope` is read-only attribute, the resource should be recreated
		},
	}, userSchemaSchema, userBaseSchemaSchema)}
}

func resourceAppUserSchemaPropertyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Set the ID before calling set..., so if we taint in set... we won't overwrite it here.
	d.SetId(fmt.Sprintf("%s/%s", d.Get("app_id").(string), d.Get("index").(string)))
	if err := setAppUserSchemaProperty(ctx, d, meta); err != nil {
		return diag.FromErr(err)
	}
	if d.Id() == "" {
		// Tainted (parent or resource missing)
		return nil
	}
	return resourceAppUserSchemaPropertyRead(ctx, d, meta)
}

func setAppUserSchemaProperty(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	err := validateAppUserSchemaProperty(d)
	if err != nil {
		return err
	}
	boc := utils.NewExponentialBackOffWithContext(ctx, 30*time.Second)
	err = backoff.Retry(func() error {
		if err := updateAppUserSubSchemaProperty(ctx, d, meta); err != nil {
			if doNotRetry(meta, err) {
				return backoff.Permanent(err)
			}
			if errors.Is(err, utils.ErrInvalidElemFormat) {
				return backoff.Permanent(err)
			}
			return err
		}
		us, resp, err := getOktaClientFromMetadata(meta).UserSchema.GetApplicationUserSchema(ctx, d.Get("app_id").(string))
		if err != nil {
			if utils.SuppressErrorOn404(resp, err) == nil {
				logger(meta).Info(
					"okta_app_user_schema_property set: 404 from parent app schema; tainting resource",
					"app_id", d.Get("app_id").(string),
					"index", d.Get("index").(string),
				)
				d.SetId("")
				return nil
			}
			return err
		}
		subSchema := UserSchemaCustomAttribute(us, d.Get("index").(string))
		if subSchema == nil {
			return fmt.Errorf("application user schema property '%s' was not created/updated for '%s' app", d.Get("index").(string), d.Get("app_id").(string))
		}
		return nil
	}, boc)
	return err
}

func resourceAppUserSchemaPropertyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	us, resp, err := getOktaClientFromMetadata(meta).UserSchema.GetApplicationUserSchema(ctx, d.Get("app_id").(string))
	if err != nil {
		if utils.SuppressErrorOn404(resp, err) == nil {
			logger(meta).Info(
				"okta_app_user_schema_property read: 404 from parent app schema; tainting resource",
				"app_id", d.Get("app_id").(string),
				"index", d.Get("index").(string),
			)
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to get application user schema property: %v", err)
	}

	subschema := UserSchemaCustomAttribute(us, d.Get("index").(string))
	if subschema == nil {
		logger(meta).Info(
			"okta_app_user_schema_property read: property missing under parent schema; tainting resource",
			"app_id", d.Get("app_id").(string),
			"index", d.Get("index").(string),
		)
		d.SetId("")
		return nil
	}
	err = syncCustomUserSchema(d, subschema)
	if subschema.Union != "" {
		if subschema.Union == "DISABLE" {
			_ = d.Set("union", false)
		} else {
			_ = d.Set("union", true)
		}
	}
	if err != nil {
		return diag.Errorf("failed to set application user schema properties: %v", err)
	}
	return nil
}

func resourceAppUserSchemaPropertyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := setAppUserSchemaProperty(ctx, d, meta); err != nil {
		return diag.FromErr(err)
	}
	if d.Id() == "" {
		// Tainted during update (parent or resource missing)
		return nil
	}
	return resourceAppUserSchemaPropertyRead(ctx, d, meta)
}

func resourceAppUserSchemaPropertyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	custom := BuildCustomUserSchema(d.Get("index").(string), nil)
	retypeUserSchemaPropertyEnums(custom)
	_, _, err := getOktaClientFromMetadata(meta).UserSchema.UpdateApplicationUserProfile(ctx, d.Get("app_id").(string), *custom)
	if err != nil {
		return diag.Errorf("failed to delete application user schema property: %v", err)
	}
	return nil
}

func updateAppUserSubSchemaProperty(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	subSchema, err := buildUserCustomSchemaAttribute(d)
	if err != nil {
		return err
	}
	if d.Get("union").(bool) {
		subSchema.Union = "ENABLE"
	} else {
		subSchema.Union = "DISABLE"
	}
	custom := BuildCustomUserSchema(d.Get("index").(string), subSchema)
	retypeUserSchemaPropertyEnums(custom)
	boc := utils.NewExponentialBackOffWithContext(ctx, 10*time.Second)
	err = backoff.Retry(func() error {
		_, resp, err := getOktaClientFromMetadata(meta).UserSchema.UpdateApplicationUserProfile(ctx, d.Get("app_id").(string), *custom)
		// if error is nil, we skip this check entirely since utils.SuppressErrorOn404 will return nil either when err is nil or when resp.StatusCode is 404.
		if err != nil && utils.SuppressErrorOn404(resp, err) == nil {
			logger(meta).Info(
				"okta_app_user_schema_property update: 404 while updating; tainting resource",
				"app_id", d.Get("app_id").(string),
				"index", d.Get("index").(string),
			)
			d.SetId("")
			return nil
		}
		if doNotRetry(meta, err) {
			return backoff.Permanent(err)
		}
		if err == nil {
			return nil
		}
		var oktaErr *sdk.Error
		if errors.As(err, &oktaErr) {
			for i := range oktaErr.ErrorCauses {
				for _, sum := range oktaErr.ErrorCauses[i] {
					if strings.Contains(sum.(string), "deletion process for an attribute with the same variable name is incomplete") {
						return err
					}
				}
			}
			return backoff.Permanent(fmt.Errorf("failed to update custom app user schema property: %w", err))
		}
		return backoff.Permanent(fmt.Errorf("failed to update custom app user schema property: %w", err))
	}, boc)
	return err
}

func validateAppUserSchemaProperty(d *schema.ResourceData) error {
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
