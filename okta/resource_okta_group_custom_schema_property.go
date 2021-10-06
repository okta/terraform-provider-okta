package okta

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceGroupCustomSchemaProperty() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupSchemaCreateOrUpdate,
		ReadContext:   resourceGroupSchemaRead,
		UpdateContext: resourceGroupSchemaCreateOrUpdate,
		DeleteContext: resourceGroupSchemaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: buildSchema(
			userBaseSchemaSchema,
			userSchemaSchema,
			map[string]*schema.Schema{
				"scope": {
					Type:             schema.TypeString,
					Optional:         true,
					Default:          "NONE",
					ValidateDiagFunc: elemInSlice([]string{"SELF", "NONE", ""}),
				},
				"master": {
					Type:     schema.TypeString,
					Optional: true,
					// Accepting an empty value to allow for zero value (when provisioning is off)
					ValidateDiagFunc: elemInSlice([]string{"PROFILE_MASTER", "OKTA", "OVERRIDE", ""}),
					Description:      "SubSchema profile manager, if not set it will inherit its setting.",
					Default:          "PROFILE_MASTER",
				},
				"master_override_priority": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "",
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
	}
}

// Sometime Okta API does not update or create custom property on the first try, thus that require running
// `terraform apply` several times. This simple retry resolves that issue. (If) After this issue will be resolved,
// this retry logic will be demolished.
func resourceGroupSchemaCreateOrUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("creating group custom schema property", "name", d.Get("index").(string))
	err := validateUserSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}
	custom := buildCustomGroupSchema(d.Get("index").(string), buildGroupCustomSchemaAttribute(d))
	subSchema, err := alterCustomGroupSchema(ctx, m, d.Get("index").(string), custom, false)
	if err != nil {
		return diag.Errorf("failed to create or update group custom schema property %s: %v", d.Get("index").(string), err)
	}
	d.SetId(d.Get("index").(string))
	err = syncCustomGroupSchema(d, subSchema)
	if err != nil {
		return diag.Errorf("failed to set group custom schema property: %v", err)
	}
	return nil
}

func resourceGroupSchemaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("reading group custom schema property", "name", d.Get("index").(string))
	s, _, err := getOktaClientFromMetadata(m).GroupSchema.GetGroupSchema(ctx)
	if err != nil {
		return diag.Errorf("failed to get group custom schema property: %v", err)
	}
	customAttribute := groupSchemaCustomAttribute(s, d.Id())
	if customAttribute == nil {
		d.SetId("")
		return nil
	}
	err = syncCustomGroupSchema(d, customAttribute)
	if err != nil {
		return diag.Errorf("failed to set group custom schema property: %v", err)
	}
	return nil
}

func alterCustomGroupSchema(ctx context.Context, m interface{}, index string, schema *okta.GroupSchema, isDeleteOperation bool) (*okta.GroupSchemaAttribute, error) {
	var schemaAttribute *okta.GroupSchemaAttribute
	timer := time.NewTimer(time.Second * 120) // sometimes it takes some time (several attempts) to recreate/delete group schema property
	ticker := time.NewTicker(time.Second * 5)
loop:
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timer.C:
			return nil, errors.New("no more attempts left")
		case <-ticker.C:
			updated, resp, err := getOktaClientFromMetadata(m).GroupSchema.UpdateGroupSchema(ctx, *schema)
			if err != nil {
				if resp != nil && resp.StatusCode == 500 {
					logger(m).Debug("updating group custom schema property caused 500 error", err)
					continue
				}
				if strings.Contains(err.Error(), "Wait until the data clean up process finishes and then try again") {
					continue
				}
				return nil, err
			}
			s, _, err := getOktaClientFromMetadata(m).GroupSchema.GetGroupSchema(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get group custom schema property: %v", err)
			}
			schemaAttribute = groupSchemaCustomAttribute(s, index)
			if isDeleteOperation && schemaAttribute == nil {
				break loop
			} else if schemaAttribute != nil && reflect.DeepEqual(schemaAttribute, updated.Definitions.Custom.Properties[index]) {
				break loop
			}
		}
	}
	return schemaAttribute, nil
}

func resourceGroupSchemaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	custom := buildCustomGroupSchema(d.Id(), nil)
	_, err := alterCustomGroupSchema(ctx, m, d.Get("index").(string), custom, true)
	if err != nil {
		return diag.Errorf("failed to delete group schema property %s: %v", d.Get("index").(string), err)
	}
	return nil
}

func buildCustomGroupSchema(index string, schema *okta.GroupSchemaAttribute) *okta.GroupSchema {
	return &okta.GroupSchema{
		Definitions: &okta.GroupSchemaDefinitions{
			Custom: &okta.GroupSchemaCustom{
				Id: "#custom",
				Properties: map[string]*okta.GroupSchemaAttribute{
					index: schema,
				},
				Type: "object",
			},
		},
	}
}

func syncCustomGroupSchema(d *schema.ResourceData, subschema *okta.GroupSchemaAttribute) error {
	syncBaseGroupSchema(d, subschema)
	_ = d.Set("description", subschema.Description)
	_ = d.Set("min_length", subschema.MinLength)
	_ = d.Set("max_length", subschema.MaxLength)
	_ = d.Set("scope", subschema.Scope)
	_ = d.Set("external_name", subschema.ExternalName)
	_ = d.Set("external_namespace", subschema.ExternalNamespace)
	_ = d.Set("unique", subschema.Unique)
	if subschema.Items != nil {
		_ = d.Set("array_type", subschema.Items.Type)
		_ = d.Set("array_one_of", flattenOneOf(subschema.Items.OneOf))
		_ = d.Set("array_enum", convertStringSliceToInterfaceSlice(subschema.Items.Enum))
	}
	if len(subschema.Enum) > 0 {
		_ = d.Set("enum", convertStringSliceToInterfaceSlice(subschema.Enum))
	}
	return setNonPrimitives(d, map[string]interface{}{
		"one_of": flattenOneOf(subschema.OneOf),
	})
}

func syncBaseGroupSchema(d *schema.ResourceData, subschema *okta.GroupSchemaAttribute) {
	_ = d.Set("title", subschema.Title)
	_ = d.Set("type", subschema.Type)
	_ = d.Set("required", subschema.Required)
	if subschema.Master != nil {
		_ = d.Set("master", subschema.Master.Type)
		if subschema.Master.Type == "OVERRIDE" {
			arr := make([]map[string]interface{}, len(subschema.Master.Priority))
			for i, st := range subschema.Master.Priority {
				arr[i] = map[string]interface{}{
					"type":  st.Type,
					"value": st.Value,
				}
			}
			_ = setNonPrimitives(d, map[string]interface{}{"master_override_priority": arr})
		}
	}
	if len(subschema.Permissions) > 0 {
		_ = d.Set("permissions", subschema.Permissions[0].Action)
	}
}

func buildGroupCustomSchemaAttribute(d *schema.ResourceData) *okta.GroupSchemaAttribute {
	return &okta.GroupSchemaAttribute{
		Title:       d.Get("title").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Required:    boolPtr(d.Get("required").(bool)),
		Permissions: []*okta.UserSchemaAttributePermission{
			{
				Action:    d.Get("permissions").(string),
				Principal: "SELF",
			},
		},
		Scope:             d.Get("scope").(string),
		Enum:              convertInterfaceToStringArrNullable(d.Get("enum")),
		Master:            getNullableMaster(d),
		Items:             getNullableItem(d),
		MinLength:         int64(d.Get("min_length").(int)),
		MaxLength:         int64(d.Get("max_length").(int)),
		OneOf:             getNullableOneOf(d, "one_of"),
		ExternalName:      d.Get("external_name").(string),
		ExternalNamespace: d.Get("external_namespace").(string),
		Unique:            d.Get("unique").(string),
	}
}

func groupSchemaCustomAttribute(s *okta.GroupSchema, index string) *okta.GroupSchemaAttribute {
	if s == nil || s.Definitions == nil || s.Definitions.Custom == nil {
		return nil
	}
	return s.Definitions.Custom.Properties[index]
}
