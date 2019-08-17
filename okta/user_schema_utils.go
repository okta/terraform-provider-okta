package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

var (
	userSchemaSchema = map[string]*schema.Schema{
		"index": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Subschema unique string identifier",
			ForceNew:    true,
		},
		"title": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Subschema title (display name)",
		},
		"type": &schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"string", "boolean", "number", "integer", "array", "object"}, false),
			Description:  "Subschema type: string, boolean, number, integer, array, or object",
			ForceNew:     true,
		},
		"array_type": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"string", "number", "integer", "reference"}, false),
			Description:  "Subschema array type: string, number, integer, reference. Type field must be an array.",
			ForceNew:     true,
		},
		"description": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Custom Subschema description",
		},
		"required": &schema.Schema{
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether the Subschema is required",
		},
		"min_length": &schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Subschema of type string minimum length",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"max_length": &schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Subschema of type string maximum length",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"enum": &schema.Schema{
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Custom Subschema enumerated value of the property. see: developer.okta.com/docs/api/resources/schemas#user-profile-schema-property-object",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"one_of": &schema.Schema{
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Custom Subschema json schemas. see: developer.okta.com/docs/api/resources/schemas#user-profile-schema-property-object",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"const": &schema.Schema{
						Required:    true,
						Type:        schema.TypeString,
						Description: "Enum value",
					},
					"title": &schema.Schema{
						Required:    true,
						Type:        schema.TypeString,
						Description: "Enum title",
					},
				},
			},
		},
		"permissions": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"HIDE", "READ_ONLY", "READ_WRITE"}, false),
			Description:  "SubSchema permissions: HIDE, READ_ONLY, or READ_WRITE.",
			Default:      "READ_ONLY",
		},
		"master": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			// Accepting an empty value to allow for zero value (when provisioning is off)
			ValidateFunc: validation.StringInSlice([]string{"PROFILE_MASTER", "OKTA", ""}, false),
			Description:  "SubSchema profile manager, if not set it will inherit its setting.",
		},
	}

	userBaseSchemaSchema = map[string]*schema.Schema{
		"index": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Subschema unique string identifier",
			ForceNew:    true,
		},
		"title": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Subschema title (display name)",
		},
		"type": &schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"string", "boolean", "number", "integer", "array", "object"}, false),
			Description:  "Subschema type: string, boolean, number, integer, array, or object",
			ForceNew:     true,
		},
		"permissions": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"HIDE", "READ_ONLY", "READ_WRITE"}, false),
			Description:  "SubSchema permissions: HIDE, READ_ONLY, or READ_WRITE.",
			Default:      "READ_ONLY",
		},
		"master": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			// Accepting an empty value to allow for zero value (when provisioning is off)
			ValidateFunc: validation.StringInSlice([]string{"PROFILE_MASTER", "OKTA", ""}, false),
			Description:  "SubSchema profile manager, if not set it will inherit its setting.",
		},
	}
)
