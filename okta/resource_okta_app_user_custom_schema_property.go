package okta

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceAppUserSchemaProperty() *schema.Resource {
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
					ValidateDiagFunc: elemInSlice([]string{"SELF", "NONE", ""}),
					ForceNew:         true, // since the `scope` is read-only attribute, the resource should be recreated
				},
				"master": {
					Type:     schema.TypeString,
					Optional: true,
					// Accepting an empty value to allow for zero value (when provisioning is off)
					ValidateDiagFunc: elemInSlice([]string{"PROFILE_MASTER", "OKTA", ""}),
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
			ValidateDiagFunc: elemInSlice([]string{"SELF", "NONE", ""}),
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
			ValidateDiagFunc: elemInSlice([]string{"SELF", "NONE", ""}),
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
	bOff := backoff.NewExponentialBackOff()
	bOff.MaxElapsedTime = time.Second * 10
	bOff.InitialInterval = time.Second
	err = backoff.Retry(func() error {
		err := resourceAppUserSchemaRead(ctx, d, m)
		if err != nil {
			return backoff.Permanent(fmt.Errorf("%s", err[0].Summary))
		}
		if d.Id() != "" {
			return nil
		}
		return fmt.Errorf("application user schema property %s was not created in app %s", d.Get("index").(string), d.Get("app_id").(string))
	}, bOff)
	return diag.FromErr(err)
}

func resourceAppUserSchemaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	us, resp, err := getOktaClientFromMetadata(m).UserSchema.GetApplicationUserSchema(ctx, d.Get("app_id").(string))
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get application user schema property: %v", err)
	}
	subschema := userSchemaCustomAttribute(us, d.Get("index").(string))
	if subschema == nil {
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
	} else {
		out, err := json.Marshal(d)
		fmt.Println(string(out))
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
	bOff := backoff.NewExponentialBackOff()
	bOff.MaxElapsedTime = time.Second * 10
	bOff.InitialInterval = time.Second
	err = backoff.Retry(func() error {
		err := resourceAppUserSchemaRead(ctx, d, m)
		if err != nil {
			return backoff.Permanent(fmt.Errorf("%s", err[0].Summary))
		}
		if d.Id() != "" {
			return nil
		}
		return fmt.Errorf("application user schema property %s was not created in the app %s", d.Get("index").(string), d.Get("app_id").(string))
	}, bOff)
	return diag.FromErr(err)
}

func resourceAppUserSchemaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	custom := buildCustomUserSchema(d.Get("index").(string), nil)
	_, _, err := getOktaClientFromMetadata(m).UserSchema.
		UpdateApplicationUserProfile(ctx, d.Get("app_id").(string), *custom)
	if err != nil {
		return diag.Errorf("failed to delete application user schema property")
	}
	return nil
}

func updateAppUserSubschema(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	subSchema := buildUserCustomSchemaAttribute(d)
	if d.Get("union").(bool) {
		subSchema.Union = "ENABLE"
	} else {
		subSchema.Union = "DISABLE"
	}
	custom := buildCustomUserSchema(d.Get("index").(string), subSchema)
	bOff := backoff.NewExponentialBackOff()
	bOff.MaxElapsedTime = time.Second * 10
	bOff.InitialInterval = time.Second
	err := backoff.Retry(func() error {
		_, _, err := getOktaClientFromMetadata(m).UserSchema.
			UpdateApplicationUserProfile(ctx, d.Get("app_id").(string), *custom)
		if err != nil {
			var oktaErr *okta.Error
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
		}
		return nil
	}, bOff)
	return diag.FromErr(err)
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
