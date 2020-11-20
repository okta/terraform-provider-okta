package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

var (
	userSchemaSchema = map[string]*schema.Schema{
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
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"string", "boolean", "number", "integer", "array", "object"}, false),
			Description:  "Subschema type: string, boolean, number, integer, array, or object",
			ForceNew:     true,
		},
		"array_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"string", "number", "integer", "reference"}, false),
			Description:  "Subschema array type: string, number, integer, reference. Type field must be an array.",
			ForceNew:     true,
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
		"required": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether the Subschema is required",
		},
		"min_length": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Subschema of type string minimum length",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"max_length": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Subschema of type string maximum length",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"enum": {
			Type:          schema.TypeList,
			Optional:      true,
			ForceNew:      true,
			Description:   "Custom Subschema enumerated value of the property. see: developer.okta.com/docs/api/resources/schemas#user-profile-schema-property-object",
			ConflictsWith: []string{"array_type"},
			Elem:          &schema.Schema{Type: schema.TypeString},
		},
		"scope": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "NONE",
			ValidateFunc: validation.StringInSlice([]string{"SELF", "NONE", ""}, false),
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
		"permissions": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"HIDE", "READ_ONLY", "READ_WRITE"}, false),
			Description:  "SubSchema permissions: HIDE, READ_ONLY, or READ_WRITE.",
			Default:      "READ_ONLY",
		},
		"master": {
			Type:     schema.TypeString,
			Optional: true,
			// Accepting an empty value to allow for zero value (when provisioning is off)
			ValidateFunc: validation.StringInSlice([]string{"PROFILE_MASTER", "OKTA", ""}, false),
			Description:  "SubSchema profile manager, if not set it will inherit its setting.",
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
			Type:          schema.TypeString,
			Optional:      true,
			Description:   "Subschema unique restriction",
			ValidateFunc:  validation.StringInSlice([]string{"UNIQUE_VALIDATED", "NOT_UNIQUE"}, false),
			ConflictsWith: []string{"one_of", "enum", "array_type"},
		},
		"user_type": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Custom Subschema usertype",
			Default:     "default",
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
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"string", "boolean", "number", "integer", "array", "object"}, false),
			Description:  "Subschema type: string, boolean, number, integer, array, or object",
			ForceNew:     true,
		},
		"permissions": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"HIDE", "READ_ONLY", "READ_WRITE"}, false),
			Description:  "SubSchema permissions: HIDE, READ_ONLY, or READ_WRITE.",
			Default:      "READ_ONLY",
		},
		"master": {
			Type:     schema.TypeString,
			Optional: true,
			// Accepting an empty value to allow for zero value (when provisioning is off)
			ValidateFunc: validation.StringInSlice([]string{"PROFILE_MASTER", "OKTA", ""}, false),
			Description:  "SubSchema profile manager, if not set it will inherit its setting.",
		},
		"required": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether the Subschema is required",
		},
		"user_type": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Custom Subschema usertype",
			Default:     "default",
		},
	}
)

func buildBaseUserSchema(target map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(userBaseSchemaSchema, target)
}

func buildCustomUserSchema(target map[string]*schema.Schema) map[string]*schema.Schema {
	return buildSchema(userSchemaSchema, target)
}

func syncUserSchema(d *schema.ResourceData, subschema *sdk.UserSubSchema) error {
	_ = d.Set("title", subschema.Title)
	_ = d.Set("type", subschema.Type)
	_ = d.Set("description", subschema.Description)
	_ = d.Set("required", subschema.Required)
	_ = d.Set("min_length", subschema.MinLength)
	_ = d.Set("max_length", subschema.MaxLength)
	_ = d.Set("scope", subschema.Scope)
	_ = d.Set("external_name", subschema.ExternalName)
	_ = d.Set("external_namespace", subschema.ExternalNamespace)
	_ = d.Set("unique", subschema.Unique)

	if subschema.Items != nil {
		_ = d.Set("array_type", subschema.Items.Type)
	}

	if subschema.Master != nil {
		_ = d.Set("master", subschema.Master.Type)
	}

	if len(subschema.Permissions) > 0 {
		_ = d.Set("permissions", subschema.Permissions[0].Action)
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
	}

	if len(subschema.Permissions) > 0 {
		_ = d.Set("permissions", subschema.Permissions[0].Action)
	}
}

func getBaseProperty(us *sdk.UserSchema, id string) *sdk.UserSubSchema {
	for key, part := range us.Definitions.Base.Properties {
		if key == id {
			return part
		}
	}
	return nil
}

func getCustomProperty(s *sdk.UserSchema, id string) *sdk.UserSubSchema {
	for key, part := range s.Definitions.Custom.Properties {
		if key == id {
			return part
		}
	}

	return nil
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
	if v, ok := d.GetOk("master"); ok {
		return &sdk.UserSchemaMaster{Type: v.(string)}
	}

	return nil
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

func getUserSubSchema(d *schema.ResourceData) *sdk.UserSubSchema {
	return &sdk.UserSubSchema{
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
}
