package okta

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

var (
	userSchemaSchema = map[string]*schema.Schema{
		"array_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The type of the array elements if `type` is set to `array`",
			ForceNew:    true,
		},
		"array_enum": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Array of values that an array property's items can be set to.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"array_one_of": {
			Type:     schema.TypeList,
			Optional: true,
			Description: `Display name and value an enum array can be set to.
	- 'const' - (Required) value mapping to member of 'enum'.
	- 'title' - (Required) display name for the enum value.`,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"const": {
						Required:    true,
						Type:        schema.TypeString,
						Description: "Value mapping to member of `array_enum`",
					},
					"title": {
						Required:    true,
						Type:        schema.TypeString,
						Description: "Display name for the enum value.",
					},
				},
			},
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The description of the user schema property.",
		},
		"min_length": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The minimum length of the user property value. Only applies to type `string`",
		},
		"max_length": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The maximum length of the user property value. Only applies to type `string`",
		},
		"enum": {
			Type:          schema.TypeList,
			Optional:      true,
			Description:   "Array of values a primitive property can be set to. See `array_enum` for arrays.",
			ConflictsWith: []string{"array_type"},
			Elem:          &schema.Schema{Type: schema.TypeString},
		},
		"one_of": {
			Type:     schema.TypeList,
			Optional: true,
			Description: `Array of maps containing a mapping for display name to enum value.
	- 'const' - (Required) value mapping to member of 'enum'.
	- 'title' - (Required) display name for the enum value.`,
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
			Description: "External name of the user schema property.",
			ForceNew:    true,
		},
		"external_namespace": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "External namespace of the user schema property.",
			ForceNew:    true,
		},
		"unique": {
			Type:          schema.TypeString,
			Optional:      true,
			Description:   "Whether the property should be unique. It can be set to `UNIQUE_VALIDATED` or `NOT_UNIQUE`.",
			ConflictsWith: []string{"one_of", "enum", "array_type"},
			ForceNew:      true,
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
			Type:        schema.TypeString,
			Required:    true,
			Description: "The type of the schema property. It can be `string`, `boolean`, `number`, `integer`, `array`, or `object`",
			ForceNew:    true,
		},
		"permissions": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Access control permissions for the property. It can be set to `READ_WRITE`, `READ_ONLY`, `HIDE`. Default: `READ_ONLY`",
			Default:     "READ_ONLY",
		},
		"required": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether the subschema is required",
		},
	}

	userTypeSchema = map[string]*schema.Schema{
		"user_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "User type ID. By default, it is `default`",
			Default:     "default",
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

func syncCustomUserSchema(d *schema.ResourceData, subschema *sdk.UserSchemaAttribute) error {
	syncBaseUserSchema(d, subschema)
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
	// Terraform SDK is statically defined at runtime for string so we need to
	// juggle types on the fly.

	if subschema.Items != nil {
		stringifyOneOfSlice(subschema.Items.Type, &subschema.Items.OneOf)
		stringifyEnumSlice(subschema.Items.Type, &subschema.Items.Enum)
		_ = d.Set("array_type", subschema.Items.Type)
		_ = d.Set("array_one_of", flattenOneOf(subschema.Items.OneOf))
		_ = d.Set("array_enum", flattenArrayEnum(subschema.Items.Enum))
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

func syncBaseUserSchema(d *schema.ResourceData, subschema *sdk.UserSchemaAttribute) {
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

func getNullableMaster(d *schema.ResourceData) *sdk.UserSchemaAttributeMaster {
	v, ok := d.GetOk("master")
	if !ok {
		return nil
	}
	usm := &sdk.UserSchemaAttributeMaster{Type: v.(string)}
	if v.(string) == "OVERRIDE" {
		mop, ok := d.Get("master_override_priority").([]interface{})
		if ok && len(mop) > 0 {
			props := make([]*sdk.UserSchemaAttributeMasterPriority, len(mop))
			for i := range mop {
				props[i] = &sdk.UserSchemaAttributeMasterPriority{
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

func buildNullableItems(d *schema.ResourceData) (*sdk.UserSchemaAttributeItems, error) {
	at, ok := d.GetOk("array_type")
	if !ok {
		return nil, nil
	}
	arrayOneOf, okArrayOneOf := d.GetOk("array_one_of")
	arrayEnum, okArrayEnum := d.GetOk("array_enum")

	u := &sdk.UserSchemaAttributeItems{
		Type: at.(string),
	}
	if !okArrayOneOf && !okArrayEnum {
		return u, nil
	}
	if okArrayEnum {
		enumSlice, ok := arrayEnum.([]interface{})
		if !ok {
			return nil, fmt.Errorf("expected %q value %+v to be an array but was an %T", "array_enum", arrayEnum, arrayEnum)
		}
		if at == "object" {
			// If array type is object, assume each item is a JSON string.  Okta
			// API expects object to be a map.  Unmarshal that into a
			// map[string]interface{} so the array marshals to correct json.
			// Previous to this subtle distinction if the operator had JSON
			// strings they would be escaped extras as absolute strings.
			for i, item := range enumSlice {
				var object map[string]interface{}
				if err := json.Unmarshal([]byte(item.(string)), &object); err != nil {
					return nil, err
				}
				enumSlice[i] = object
			}
		}
		u.Enum = enumSlice
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

func buildOneOf(ae []interface{}, elemType string) ([]*sdk.UserSchemaAttributeEnum, error) {
	oneOf := make([]*sdk.UserSchemaAttributeEnum, len(ae))
	for i := range ae {
		valueMap := ae[i].(map[string]interface{})
		oneOf[i] = &sdk.UserSchemaAttributeEnum{
			Title: valueMap["title"].(string),
		}
		value := valueMap["const"]
		switch elemType {
		case "object":
			var object map[string]interface{}
			if err := json.Unmarshal([]byte(value.(string)), &object); err == nil {
				oneOf[i].Const = object
			}
		default:
			oneOf[i].Const = value.(string)
		}
	}
	return oneOf, nil
}

func flattenOneOf(oneOf []*sdk.UserSchemaAttributeEnum) []interface{} {
	result := make([]interface{}, len(oneOf))
	for i, v := range oneOf {
		var value string
		if obj, ok := v.Const.(map[string]interface{}); ok {
			objB, _ := json.Marshal(obj)
			value = string(objB)
		} else {
			value = v.Const.(string)
		}
		of := map[string]interface{}{
			"title": v.Title,
			"const": value,
		}
		result[i] = of
	}
	return result
}

func flattenArrayEnum(arrayEnum []interface{}) []interface{} {
	result := make([]interface{}, len(arrayEnum))
	for i, v := range arrayEnum {
		var value string
		if obj, ok := v.(map[string]interface{}); ok {
			objB, _ := json.Marshal(obj)
			value = string(objB)
		} else {
			value = v.(string)
		}
		result[i] = value
	}
	return result
}

func buildUserCustomSchemaAttribute(d *schema.ResourceData) (*sdk.UserSchemaAttribute, error) {
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
	attribute := &sdk.UserSchemaAttribute{
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

func buildUserBaseSchemaAttribute(d *schema.ResourceData) *sdk.UserSchemaAttribute {
	userSchemaAttribute := &sdk.UserSchemaAttribute{
		Master: getNullableMaster(d),
		Title:  d.Get("title").(string),
		Type:   d.Get("type").(string),
		Permissions: []*sdk.UserSchemaAttributePermission{
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
	us := &sdk.UserSchema{
		Definitions: &sdk.UserSchemaDefinitions{
			Base: &sdk.UserSchemaBase{
				Id: "#base",
				Properties: map[string]*sdk.UserSchemaAttribute{
					d.Get("index").(string): buildUserBaseSchemaAttribute(d),
				},
				Type: "object",
			},
		},
	}
	type localIDX sdk.UserSchema
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

func buildCustomUserSchema(index string, schema *sdk.UserSchemaAttribute) *sdk.UserSchema {
	return &sdk.UserSchema{
		Definitions: &sdk.UserSchemaDefinitions{
			Custom: &sdk.UserSchemaPublic{
				Id: "#custom",
				Properties: map[string]*sdk.UserSchemaAttribute{
					index: schema,
				},
				Type: "object",
			},
		},
	}
}

func userSchemaCustomAttribute(s *sdk.UserSchema, index string) *sdk.UserSchemaAttribute {
	if s == nil || s.Definitions == nil || s.Definitions.Custom == nil {
		return nil
	}
	return s.Definitions.Custom.Properties[index]
}

func userSchemaBaseAttribute(s *sdk.UserSchema, index string) *sdk.UserSchemaAttribute {
	if s == nil || s.Definitions == nil || s.Definitions.Base == nil {
		return nil
	}
	return s.Definitions.Base.Properties[index]
}

// retypeUserSchemaPropertyEnums takes a schema and ensures the enums in its
// UserSchemaAttribute(s) have the correct golang type values instead of the
// strings limitation due to the TF SDK.
func retypeUserSchemaPropertyEnums(schema *sdk.UserSchema) {
	if schema.Definitions != nil && schema.Definitions.Base != nil {
		retypeUserPropertiesEnum(schema.Definitions.Base.Properties)
	}
	if schema.Definitions != nil && schema.Definitions.Custom != nil {
		retypeUserPropertiesEnum(schema.Definitions.Custom.Properties)
	}
}

// stringifyUserSchemaPropertyEnums takes a schema and ensures the enums in its
// UserSchemaAttribute(s) have string values to satisfy the TF schema
func stringifyUserSchemaPropertyEnums(schema *sdk.UserSchema) {
	if schema.Definitions != nil && schema.Definitions.Base != nil {
		stringifyUserPropertiesEnum(schema.Definitions.Base.Properties)
	}
	if schema.Definitions != nil && schema.Definitions.Custom != nil {
		stringifyUserPropertiesEnum(schema.Definitions.Custom.Properties)
	}
}

func retypeUserPropertiesEnum(properties map[string]*sdk.UserSchemaAttribute) {
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

func stringifyUserPropertiesEnum(properties map[string]*sdk.UserSchemaAttribute) {
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

func retypeEnumSlice(elemType string, enum []interface{}) []interface{} {
	result := make([]interface{}, len(enum))
	for i, val := range enum {
		v, err := coerceCorrectTypedValue(elemType, val)
		if err == nil {
			result[i] = v
		}
	}
	return result
}

func stringifyEnumSlice(elemType string, enum *[]interface{}) {
	if enum == nil {
		return
	}
	for i, val := range *enum {
		v, err := coerceStringValue(elemType, val)
		if err == nil {
			(*enum)[i] = v
		}
	}
}

func retypeOneOfSlice(elemType string, enum []*sdk.UserSchemaAttributeEnum) []*sdk.UserSchemaAttributeEnum {
	result := make([]*sdk.UserSchemaAttributeEnum, len(enum))
	for i, val := range enum {
		ae := sdk.UserSchemaAttributeEnum{}
		if val != nil {
			ae.Title = val.Title
			if val.Const != nil {
				v, err := coerceCorrectTypedValue(elemType, val.Const)
				if err == nil {
					ae.Const = v
				}
			}
		}
		result[i] = &ae
	}
	return result
}

func stringifyOneOfSlice(elemType string, enum *[]*sdk.UserSchemaAttributeEnum) {
	if enum == nil {
		return
	}
	for _, val := range *enum {
		if val != nil {
			if val.Const != nil {
				v, err := coerceStringValue(elemType, val.Const)
				if err == nil {
					val.Const = v
				}
			}
		}
	}
}

func coerceCorrectTypedValue(elemType string, value interface{}) (interface{}, error) {
	switch elemType {
	case "number":
		return coerceFloat64(value)
	case "integer":
		return coerceInt(value)
	case "boolean":
		return coerceBool(value)
	case "object":
		return value, nil
	default:
		if str, ok := value.(string); ok {
			return str, nil
		}
		return nil, fmt.Errorf("could not coerce %+v of type %T to string", value, value)
	}
}

func coerceFloat64(value interface{}) (float64, error) {
	if v, ok := value.(float64); ok {
		return v, nil
	}
	if str, ok := value.(string); ok {
		return strconv.ParseFloat(str, 64)
	}
	return 0.0, fmt.Errorf("could not coerce %+v of type %T to float64", value, value)
}

func coerceInt(value interface{}) (int, error) {
	if v, ok := value.(float64); ok {
		return int(v), nil
	}
	if v, ok := value.(int); ok {
		return v, nil
	}
	if str, ok := value.(string); ok {
		return strconv.Atoi(str)
	}
	return 0, fmt.Errorf("could not coerce %+v of type %T to int", value, value)
}

func coerceBool(value interface{}) (bool, error) {
	if v, ok := value.(bool); ok {
		return v, nil
	}
	if str, ok := value.(string); ok {
		return strconv.ParseBool(str)
	}
	return false, fmt.Errorf("could not coerce %+v of type %T to bool", value, value)
}

func coerceStringValue(elemType string, value interface{}) (interface{}, error) {
	switch elemType {
	case "number":
		if v, ok := value.(float64); ok {
			return fmt.Sprintf("%g", v), nil
		}
		return nil, fmt.Errorf("could not coerce %+v of type %T to float64", value, value)
	case "integer":
		if v, ok := value.(float64); ok {
			return fmt.Sprintf("%d", int(v)), nil
		}
		if v, ok := value.(int); ok {
			return fmt.Sprintf("%d", v), nil
		}
		return nil, fmt.Errorf("could not coerce %+v of type %T to int", value, value)
	case "boolean":
		if v, ok := value.(bool); ok {
			return fmt.Sprintf("%t", v), nil
		}
		return nil, fmt.Errorf("could not coerce %+v of type %T to bool", value, value)
	case "object":
		return value, nil
	default:
		if str, ok := value.(string); ok {
			return str, nil
		}
		return nil, fmt.Errorf("could not coerce %+v of type %T to string", value, value)
	}
}
