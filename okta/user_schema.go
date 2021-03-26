package okta

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

var (
	userSchemaSchema = map[string]*schema.Schema{
		"array_type": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringInSlice([]string{"string", "number", "integer", "reference"}),
			Description:      "Subschema array type: string, number, integer, reference. Type field must be an array.",
			ForceNew:         true,
		},
		"array_enum": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Description: "Custom Subschema enumerated value of a property of type array.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"array_one_of": {
			Type:        schema.TypeList,
			ForceNew:    true,
			Optional:    true,
			Description: "array of valid JSON schemas for property type array.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"const": {
						Required:    true,
						Type:        schema.TypeString,
						Description: "Enum value",
					},
					"title": {
						Required:    true,
						Type:        schema.TypeString,
						Description: "Enum title",
					},
				},
			},
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Custom Subschema description",
		},
		"min_length": {
			Type:             schema.TypeInt,
			Optional:         true,
			Description:      "Subschema of type string minimum length",
			ValidateDiagFunc: intAtLeast(1),
		},
		"max_length": {
			Type:             schema.TypeInt,
			Optional:         true,
			Description:      "Subschema of type string maximum length",
			ValidateDiagFunc: intAtLeast(1),
		},
		"enum": {
			Type:          schema.TypeList,
			Optional:      true,
			ForceNew:      true,
			Description:   "Custom Subschema enumerated value of the property. see: developer.okta.com/docs/api/resources/schemas#user-profile-schema-property-object",
			ConflictsWith: []string{"array_type"},
			Elem:          &schema.Schema{Type: schema.TypeString},
		},
		"one_of": {
			Type:          schema.TypeList,
			ForceNew:      true,
			Optional:      true,
			Description:   "Custom Subschema json schemas. see: developer.okta.com/docs/api/resources/schemas#user-profile-schema-property-object",
			ConflictsWith: []string{"array_type"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"const": {
						Required:    true,
						Type:        schema.TypeString,
						Description: "Enum value",
					},
					"title": {
						Required:    true,
						Type:        schema.TypeString,
						Description: "Enum title",
					},
				},
			},
		},
		"external_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Subschema external name",
			ForceNew:    true,
		},
		"external_namespace": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Subschema external namespace",
			ForceNew:    true,
		},
		"unique": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Subschema unique restriction",
			ValidateDiagFunc: stringInSlice([]string{"UNIQUE_VALIDATED", "NOT_UNIQUE"}),
			ConflictsWith:    []string{"one_of", "enum", "array_type"},
		},
	}

	userBaseSchemaSchema = map[string]*schema.Schema{
		"index": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Subschema unique string identifier",
			ForceNew:    true,
		},
		"title": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Subschema title (display name)",
		},
		"type": {
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: stringInSlice([]string{"string", "boolean", "number", "integer", "array", "object"}),
			Description:      "Subschema type: string, boolean, number, integer, array, or object",
			ForceNew:         true,
		},
		"permissions": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringInSlice([]string{"HIDE", "READ_ONLY", "READ_WRITE"}),
			Description:      "SubSchema permissions: HIDE, READ_ONLY, or READ_WRITE.",
			Default:          "READ_ONLY",
		},
		"required": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether the subschema is required",
		},
	}

	userTypeSchema = map[string]*schema.Schema{
		"user_type": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Custom subschema user type",
			Default:          "default",
			ValidateDiagFunc: stringAtLeast(7),
		},
	}

	userPatternSchema = map[string]*schema.Schema{
		"pattern": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The validation pattern to use for the subschema. Must be in form of '.+', or '[<pattern>]+' if present.'",
			ForceNew:    false,
		},
	}
)

func syncUserSchema(d *schema.ResourceData, subschema *sdk.UserSubSchema) error {
	syncBaseUserSchema(d, subschema)
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
		_ = d.Set("array_enum", convertStringArrToInterface(subschema.Items.Enum))
	}
	return setNonPrimitives(d, map[string]interface{}{
		"enum":   subschema.Enum,
		"one_of": flattenOneOf(subschema.OneOf),
	})
}

func syncBaseUserSchema(d *schema.ResourceData, subschema *sdk.UserSubSchema) {
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
	if subschema.Pattern != nil {
		_ = d.Set("pattern", &subschema.Pattern)
	} else {
		_ = d.Set("pattern", "")
	}
}

func getBaseProperty(s *sdk.UserSchema, id string) *sdk.UserSubSchema {
	if s == nil || s.Definitions == nil || s.Definitions.Base == nil {
		return nil
	}
	return s.Definitions.Base.Properties[id]
}

func getCustomProperty(s *sdk.UserSchema, id string) *sdk.UserSubSchema {
	if s == nil || s.Definitions == nil || s.Definitions.Custom == nil {
		return nil
	}
	return s.Definitions.Custom.Properties[id]
}

func getNullableOneOf(d *schema.ResourceData, key string) (oneOf []*sdk.UserSchemaEnum) {
	oneOf = []*sdk.UserSchemaEnum{}

	if oneOfList, ok := d.GetOk(key); ok {
		for _, v := range oneOfList.([]interface{}) {
			valueMap := v.(map[string]interface{})
			oneOf = append(oneOf, &sdk.UserSchemaEnum{
				Const: valueMap["const"].(string),
				Title: valueMap["title"].(string),
			})
		}
	}

	return oneOf
}

func getNullableMaster(d *schema.ResourceData) *sdk.UserSchemaMaster {
	v, ok := d.GetOk("master")
	if !ok {
		return nil
	}
	usm := &sdk.UserSchemaMaster{Type: v.(string)}
	if v.(string) == "OVERRIDE" {
		mop, ok := d.Get("master_override_priority").([]interface{})
		if ok && len(mop) > 0 {
			props := make([]sdk.UserSchemaMasterPriority, len(mop))
			for i := range mop {
				props[i] = sdk.UserSchemaMasterPriority{
					Type:  d.Get(fmt.Sprintf("master_override_priority.%d.type", i)).(string),
					Value: d.Get(fmt.Sprintf("master_override_priority.%d.value", i)).(string),
				}
			}
			usm.Priority = props
		}
	}
	return usm
}

func getNullableItem(d *schema.ResourceData) *sdk.UserSchemaItem {
	if v, ok := d.GetOk("array_type"); ok {
		return &sdk.UserSchemaItem{
			Type:  v.(string),
			OneOf: getNullableOneOf(d, "array_one_of"),
			Enum:  convertInterfaceToStringArrNullable(d.Get("array_enum")),
		}
	}

	return nil
}

func flattenOneOf(oneOf []*sdk.UserSchemaEnum) []interface{} {
	result := make([]interface{}, len(oneOf))
	for i, v := range oneOf {
		result[i] = map[string]interface{}{
			"const": v.Const,
			"title": v.Title,
		}
	}
	return result
}

func userSubSchema(d *schema.ResourceData) *sdk.UserSubSchema {
	subschema := &sdk.UserSubSchema{
		Title:       d.Get("title").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Required:    boolPtr(d.Get("required").(bool)),
		Permissions: []*sdk.UserSchemaPermission{
			{
				Action:    d.Get("permissions").(string),
				Principal: "SELF",
			},
		},
		Scope:             d.Get("scope").(string),
		Enum:              convertInterfaceToStringArrNullable(d.Get("enum")),
		Master:            getNullableMaster(d),
		Items:             getNullableItem(d),
		MinLength:         getNullableInt(d, "min_length"),
		MaxLength:         getNullableInt(d, "max_length"),
		OneOf:             getNullableOneOf(d, "one_of"),
		ExternalName:      d.Get("external_name").(string),
		ExternalNamespace: d.Get("external_namespace").(string),
		Unique:            d.Get("unique").(string),
	}
	p, ok := d.GetOk("pattern")
	if ok {
		subschema.Pattern = stringPtr(p.(string))
	}
	return subschema
}

func userBasedSubSchema(d *schema.ResourceData) *sdk.UserSubSchema {
	subSchema := &sdk.UserSubSchema{
		Master: getNullableMaster(d),
		Title:  d.Get("title").(string),
		Type:   d.Get("type").(string),
		Permissions: []*sdk.UserSchemaPermission{
			{
				Action:    d.Get("permissions").(string),
				Principal: "SELF",
			},
		},
		Required: boolPtr(d.Get("required").(bool)),
	}
	if d.Get("index").(string) == "login" {
		subSchema.IsLogin = true
		p := d.Get("pattern").(string)
		if p != "" {
			subSchema.Pattern = stringPtr(p)
		}
	}
	return subSchema
}
