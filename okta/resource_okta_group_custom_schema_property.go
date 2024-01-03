package okta

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
	"github.com/okta/terraform-provider-okta/sdk"
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
					Type:     schema.TypeString,
					Optional: true,
					Default:  "NONE",
				},
				"master": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "SubSchema profile manager, if not set it will inherit its setting.",
					Default:     "PROFILE_MASTER",
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
	groupCustomSchemaAttribute, err := buildGroupCustomSchemaAttribute(d)
	if err != nil {
		return diag.FromErr(err)
	}
	custom := buildCustomGroupSchema(d.Get("index").(string), groupCustomSchemaAttribute)
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
	_ = d.Set("index", d.Id())
	err = syncCustomGroupSchema(d, customAttribute)
	if err != nil {
		return diag.Errorf("failed to set group custom schema property: %v", err)
	}
	return nil
}

func alterCustomGroupSchema(ctx context.Context, m interface{}, index string, schema *sdk.GroupSchema, isDeleteOperation bool) (*sdk.GroupSchemaAttribute, error) {
	var schemaAttribute *sdk.GroupSchemaAttribute

	boc := newExponentialBackOffWithContext(ctx, 120*time.Second)
	err := backoff.Retry(func() error {
		// NOTE: Enums on the schema can be typed other than string but the
		// Terraform SDK is staticly defined at runtime for string so we need to
		// juggle types on the fly.

		retypeGroupSchemaPropertyEnums(schema)
		updated, resp, err := getOktaClientFromMetadata(m).GroupSchema.UpdateGroupSchema(ctx, *schema)
		stringifyGroupSchemaPropertyEnums(schema)
		if doNotRetry(m, err) {
			return backoff.Permanent(err)
		}

		if err != nil {
			if resp != nil && resp.StatusCode == 500 {
				return fmt.Errorf("updating group custom schema property caused 500 error: %w", err)
			}
			if strings.Contains(err.Error(), "Wait until the data clean up process finishes and then try again") {
				return err
			}
			return backoff.Permanent(err)
		}
		s, _, err := getOktaClientFromMetadata(m).GroupSchema.GetGroupSchema(ctx)
		if err != nil {
			return backoff.Permanent(fmt.Errorf("failed to get group custom schema property: %v", err))
		}
		schemaAttribute = groupSchemaCustomAttribute(s, index)
		if isDeleteOperation && schemaAttribute == nil {
			return nil
		} else if schemaAttribute != nil && reflect.DeepEqual(schemaAttribute, updated.Definitions.Custom.Properties[index]) {
			return nil
		}
		return errors.New("failed to apply changes after several retries")
	}, boc)
	if err != nil {
		logger(m).Error("failed to apply changes after several retries", err)
	}
	return schemaAttribute, err
}

func resourceGroupSchemaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	custom := buildCustomGroupSchema(d.Id(), nil)
	_, err := alterCustomGroupSchema(ctx, m, d.Get("index").(string), custom, true)
	if err != nil {
		return diag.Errorf("failed to delete group schema property %s: %v", d.Get("index").(string), err)
	}
	return nil
}

func buildCustomGroupSchema(index string, schema *sdk.GroupSchemaAttribute) *sdk.GroupSchema {
	return &sdk.GroupSchema{
		Definitions: &sdk.GroupSchemaDefinitions{
			Custom: &sdk.GroupSchemaCustom{
				Id: "#custom",
				Properties: map[string]*sdk.GroupSchemaAttribute{
					index: schema,
				},
				Type: "object",
			},
		},
	}
}

func syncCustomGroupSchema(d *schema.ResourceData, subschema *sdk.GroupSchemaAttribute) error {
	syncBaseGroupSchema(d, subschema)
	_ = d.Set("description", subschema.Description)
	if subschema.MinLengthPtr != nil {
		_ = d.Set("min_length", *subschema.MinLengthPtr)
	}
	if subschema.MaxLengthPtr != nil {
		_ = d.Set("max_length", *subschema.MaxLengthPtr)
	}
	_ = d.Set("scope", subschema.Scope)
	_ = d.Set("external_name", subschema.ExternalName)
	_ = d.Set("external_namespace", subschema.ExternalNamespace)
	_ = d.Set("unique", subschema.Unique)

	// NOTE: Enums on the schema can be typed other than string but the
	// Terraform SDK is staticly defined at runtime for string so we need to
	// juggle types on the fly.

	if subschema.Items != nil {
		stringifyOneOfSlice(subschema.Items.Type, &subschema.Items.OneOf)
		stringifyEnumSlice(subschema.Items.Type, &subschema.Items.Enum)
		_ = d.Set("array_type", subschema.Items.Type)
		_ = d.Set("array_one_of", flattenOneOf(subschema.Items.OneOf))
		_ = d.Set("array_enum", subschema.Items.Enum)
	}

	stringifyOneOfSlice(subschema.Type, &subschema.OneOf)
	stringifyEnumSlice(subschema.Type, &subschema.Enum)

	if len(subschema.Enum) > 0 {
		_ = d.Set("enum", subschema.Enum)
	}

	return setNonPrimitives(d, map[string]interface{}{
		"one_of": flattenOneOf(subschema.OneOf),
	})
}

func syncBaseGroupSchema(d *schema.ResourceData, subschema *sdk.GroupSchemaAttribute) {
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

func buildGroupCustomSchemaAttribute(d *schema.ResourceData) (*sdk.GroupSchemaAttribute, error) {
	items, err := buildNullableItems(d)
	if err != nil {
		return nil, err
	}
	var oneOf []*sdk.UserSchemaAttributeEnum
	if rawOneOf, ok := d.GetOk("one_of"); ok {
		oneOf, err = buildOneOf(rawOneOf.([]interface{}), d.Get("type").(string))
		if err != nil {
			return nil, err
		}
	}
	attribute := &sdk.GroupSchemaAttribute{
		Title:       d.Get("title").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Required:    boolPtr(d.Get("required").(bool)),
		Permissions: []*sdk.UserSchemaAttributePermission{
			{
				Action:    d.Get("permissions").(string),
				Principal: "SELF",
			},
		},
		Scope:             d.Get("scope").(string),
		Master:            getNullableMaster(d),
		Items:             items,
		OneOf:             oneOf,
		ExternalName:      d.Get("external_name").(string),
		ExternalNamespace: d.Get("external_namespace").(string),
		Unique:            d.Get("unique").(string),
	}
	if min, ok := d.GetOk("min_length"); ok {
		attribute.MinLengthPtr = int64Ptr(min.(int))
	}
	if max, ok := d.GetOk("max_length"); ok {
		attribute.MaxLengthPtr = int64Ptr(max.(int))
	}
	if rawEnum, ok := d.GetOk("enum"); ok {
		attribute.Enum = rawEnum.([]interface{})
	}
	return attribute, nil
}

func groupSchemaCustomAttribute(s *sdk.GroupSchema, index string) *sdk.GroupSchemaAttribute {
	if s == nil || s.Definitions == nil || s.Definitions.Custom == nil {
		return nil
	}
	return s.Definitions.Custom.Properties[index]
}

// retypeGroupSchemaPropertyEnums takes a schema and ensures the enums in its
// GroupSchemaAttribute(s) have the correct golang type values instead of the
// strings limitation due to the TF SDK.
func retypeGroupSchemaPropertyEnums(schema *sdk.GroupSchema) {
	if schema.Definitions != nil && schema.Definitions.Base != nil {
		retypeGroupPropertiesEnum(schema.Definitions.Base.Properties)
	}
	if schema.Definitions != nil && schema.Definitions.Custom != nil {
		retypeGroupPropertiesEnum(schema.Definitions.Custom.Properties)
	}
}

// stringifyGroupSchemaPropertyEnums takes a schema and ensures the enums in its
// GroupSchemaAttribute(s) have string values to satisfy the TF schema
func stringifyGroupSchemaPropertyEnums(schema *sdk.GroupSchema) {
	if schema.Definitions != nil && schema.Definitions.Base != nil {
		stringifyGroupPropertiesEnum(schema.Definitions.Base.Properties)
	}
	if schema.Definitions != nil && schema.Definitions.Custom != nil {
		stringifyGroupPropertiesEnum(schema.Definitions.Custom.Properties)
	}
}

func retypeGroupPropertiesEnum(properties map[string]*sdk.GroupSchemaAttribute) {
	for _, val := range properties {
		if val == nil {
			continue
		}
		enum := retypeEnumSlice(val.Type, val.Enum)
		val.Enum = enum
		attributeEnum := retypeOneOfSlice(val.Type, val.OneOf)
		val.OneOf = attributeEnum
		if val.Items != nil {
			enum := retypeEnumSlice(val.Items.Type, val.Items.Enum)
			val.Items.Enum = enum
			retypeOneOfSlice(val.Type, val.OneOf)
			attributeEnum := retypeOneOfSlice(val.Items.Type, val.Items.OneOf)
			val.Items.OneOf = attributeEnum
		}

	}
}

func stringifyGroupPropertiesEnum(properties map[string]*sdk.GroupSchemaAttribute) {
	for _, val := range properties {
		if val != nil && val.Enum != nil {
			stringifyEnumSlice(val.Type, &val.Enum)
		}
		if val != nil && val.OneOf != nil {
			stringifyOneOfSlice(val.Type, &val.OneOf)
		}
		if val != nil && val.Items != nil {
			stringifyEnumSlice(val.Items.Type, &val.Items.Enum)
			stringifyOneOfSlice(val.Items.Type, &val.Items.OneOf)
		}

	}
}
