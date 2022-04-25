package okta

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

var (
	userSchemaSchema = map[string]*schema.Schema{
		"array_type": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: elemInSlice([]string{"string", "number", "integer", "reference"}),
			Description:      "Subschema array type: string, number, integer, reference. Type field must be an array.",
			ForceNew:         true,
		},
		"array_enum": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Custom Subschema enumerated value of a property of type array.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"array_one_of": {
			Type:        schema.TypeList,
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
			Description:   "Custom Subschema enumerated value of the property. see: developer.okta.com/docs/api/resources/schemas#user-profile-schema-property-object",
			ConflictsWith: []string{"array_type"},
			Elem:          &schema.Schema{Type: schema.TypeString},
		},
		"one_of": {
			Type:          schema.TypeList,
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
			ValidateDiagFunc: elemInSlice([]string{"UNIQUE_VALIDATED", "NOT_UNIQUE"}),
			ConflictsWith:    []string{"one_of", "enum", "array_type"},
			ForceNew:         true,
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
			ValidateDiagFunc: elemInSlice([]string{"string", "boolean", "number", "integer", "array", "object"}),
			Description:      "Subschema type: string, boolean, number, integer, array, or object",
			ForceNew:         true,
		},
		"permissions": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: elemInSlice([]string{"HIDE", "READ_ONLY", "READ_WRITE"}),
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

func syncCustomUserSchema(d *schema.ResourceData, subschema *okta.UserSchemaAttribute) error {
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
		_ = d.Set("array_enum", subschema.Items.Enum)
	}
	if len(subschema.Enum) > 0 {
		_ = d.Set("enum", subschema.Enum)
	}
	return setNonPrimitives(d, map[string]interface{}{
		"one_of": flattenOneOf(subschema.OneOf),
	})
}

func syncBaseUserSchema(d *schema.ResourceData, subschema *okta.UserSchemaAttribute) {
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
	}
}

func getNullableMaster(d *schema.ResourceData) *okta.UserSchemaAttributeMaster {
	v, ok := d.GetOk("master")
	if !ok {
		return nil
	}
	usm := &okta.UserSchemaAttributeMaster{Type: v.(string)}
	if v.(string) == "OVERRIDE" {
		mop, ok := d.Get("master_override_priority").([]interface{})
		if ok && len(mop) > 0 {
			props := make([]*okta.UserSchemaAttributeMasterPriority, len(mop))
			for i := range mop {
				props[i] = &okta.UserSchemaAttributeMasterPriority{
					Type:  d.Get(fmt.Sprintf("master_override_priority.%d.type", i)).(string),
					Value: d.Get(fmt.Sprintf("master_override_priority.%d.value", i)).(string),
				}
			}
			usm.Priority = props
		}
	}
	return usm
}

var errInvalidElemFormat = errors.New("element type does not match the value provided in 'array_type' or 'type'")

func buildNullableItems(d *schema.ResourceData) (*okta.UserSchemaAttributeItems, error) {
	at, ok := d.GetOk("array_type")
	if !ok {
		return nil, nil
	}
	arrayOneOf, okArrayOneOf := d.GetOk("array_one_of")
	arrayEnum, okArrayEnum := d.GetOk("array_enum")

	u := &okta.UserSchemaAttributeItems{
		Type: at.(string),
	}
	if !okArrayOneOf && !okArrayEnum {
		return u, nil
	}
	if okArrayEnum {
		enum := buildStringSlice(arrayEnum.([]interface{}))
		u.Enum = enum
	}
	if okArrayOneOf {
		oneOf, err := buildOneOf(arrayOneOf.([]interface{}), u.Type)
		if err != nil {
			return nil, err
		}
		u.OneOf = oneOf
	}
	return u, nil
}

func buildOneOf(ae []interface{}, elemType string) ([]*okta.UserSchemaAttributeEnum, error) {
	oneOf := make([]*okta.UserSchemaAttributeEnum, len(ae))
	for i := range ae {
		valueMap := ae[i].(map[string]interface{})
		oneOf[i] = &okta.UserSchemaAttributeEnum{
			Title: valueMap["title"].(string),
		}
		c := valueMap["const"].(string)
		oneOf[i].Const = c
	}
	return oneOf, nil
}

func buildStringSlice(enum []interface{}) []string {
	result := make([]string, len(enum))
	for i := range enum {
		result[i] = enum[i].(string)
	}
	return result
}

func strToInterfaceSlice(enum []string) []interface{} {
	result := make([]interface{}, len(enum))
	for i := range enum {
		result[i] = enum[i]
	}

	return result
}

func flattenOneOf(oneOf []*okta.UserSchemaAttributeEnum) []interface{} {
	result := make([]interface{}, len(oneOf))
	for i, v := range oneOf {
		of := map[string]interface{}{
			"title": v.Title,
			"const": v.Const,
		}
		result[i] = of
	}
	return result
}

func buildUserCustomSchemaAttribute(d *schema.ResourceData) (*okta.UserSchemaAttribute, error) {
	items, err := buildNullableItems(d)
	if err != nil {
		return nil, err
	}
	var oneOf []*okta.UserSchemaAttributeEnum
	if rawOneOf, ok := d.GetOk("one_of"); ok {
		oneOf, err = buildOneOf(rawOneOf.([]interface{}), d.Get("type").(string))
		if err != nil {
			return nil, err
		}
	}
	var enum []string
	if rawEnum, ok := d.GetOk("enum"); ok {
		enum = buildStringSlice(rawEnum.([]interface{}))
	}
	return &okta.UserSchemaAttribute{
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
		Enum:              enum,
		Master:            getNullableMaster(d),
		Items:             items,
		MinLength:         int64(d.Get("min_length").(int)),
		MaxLength:         int64(d.Get("max_length").(int)),
		OneOf:             oneOf,
		ExternalName:      d.Get("external_name").(string),
		ExternalNamespace: d.Get("external_namespace").(string),
		Unique:            d.Get("unique").(string),
	}, nil
}

func buildUserBaseSchemaAttribute(d *schema.ResourceData) *okta.UserSchemaAttribute {
	userSchemaAttribute := &okta.UserSchemaAttribute{
		Master: getNullableMaster(d),
		Title:  d.Get("title").(string),
		Type:   d.Get("type").(string),
		Permissions: []*okta.UserSchemaAttributePermission{
			{
				Action:    d.Get("permissions").(string),
				Principal: "SELF",
			},
		},
		Required: boolPtr(d.Get("required").(bool)),
	}
	if d.Get("index").(string) == "login" {
		p, ok := d.GetOk("pattern")
		if ok {
			userSchemaAttribute.Pattern = stringPtr(p.(string))
		}
	}
	return userSchemaAttribute
}

func buildBaseUserSchema(d *schema.ResourceData) []byte {
	us := &okta.UserSchema{
		Definitions: &okta.UserSchemaDefinitions{
			Base: &okta.UserSchemaBase{
				Id: "#base",
				Properties: map[string]*okta.UserSchemaAttribute{
					d.Get("index").(string): buildUserBaseSchemaAttribute(d),
				},
				Type: "object",
			},
		},
	}
	type localIDX okta.UserSchema
	m, _ := json.Marshal((*localIDX)(us))
	if d.Get("index").(string) != "login" {
		return m
	}
	var a interface{}
	_ = json.Unmarshal(m, &a)
	b := a.(map[string]interface{})
	p := us.Definitions.Base.Properties["login"].Pattern
	if p == nil {
		b["definitions"].(map[string]interface{})["base"].(map[string]interface{})["properties"].(map[string]interface{})["login"].(map[string]interface{})["pattern"] = nil
	}
	m, _ = json.Marshal(b)
	return m
}

func buildCustomUserSchema(index string, schema *okta.UserSchemaAttribute) *okta.UserSchema {
	return &okta.UserSchema{
		Definitions: &okta.UserSchemaDefinitions{
			Custom: &okta.UserSchemaPublic{
				Id: "#custom",
				Properties: map[string]*okta.UserSchemaAttribute{
					index: schema,
				},
				Type: "object",
			},
		},
	}
}

func userSchemaCustomAttribute(s *okta.UserSchema, index string) *okta.UserSchemaAttribute {
	if s == nil || s.Definitions == nil || s.Definitions.Custom == nil {
		return nil
	}
	return s.Definitions.Custom.Properties[index]
}

func userSchemaBaseAttribute(s *okta.UserSchema, index string) *okta.UserSchemaAttribute {
	if s == nil || s.Definitions == nil || s.Definitions.Base == nil {
		return nil
	}
	return s.Definitions.Base.Properties[index]
}
