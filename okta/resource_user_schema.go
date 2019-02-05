package okta

import (
	"encoding/json"
	"fmt"
	"log"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const customSchema = "custom"

func resourceUserSchema() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserSchemaCreate,
		Read:   resourceUserSchemaRead,
		Update: resourceUserSchemaUpdate,
		Delete: resourceUserSchemaDelete,
		Exists: resourceUserSchemaExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
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
				ValidateFunc: validation.StringInSlice([]string{"string", "boolean", "number", "integer", "array"}, false),
				Description:  "Subschema type: string, boolean, number, integer, or array",
				ForceNew:     true,
			},
			"array_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"string", "number", "interger", "reference"}, false),
				Description:  "Subschema array type: string, number, interger, reference. Type field must be an array.",
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
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subschema of type string minlength",
			},
			"max_length": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subschema of type string maxlength",
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
				Elem: &schema.Schema{
					Type: schema.TypeMap,
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
			},
			"permissions": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"HIDE", "READ_ONLY", "READ_WRITE"}, false),
				Description:  "SubSchema permissions: HIDE, READ_ONLY, or READ_WRITE.",
				Default:      "READ_ONLY",
			},
			"master": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PROFILE_MASTER", "OKTA"}, false),
				Description:  "SubSchema profile manager: PROFILE_MASTER or OKTA.",
				Default:      "PROFILE_MASTER",
			},
		},
	}
}

func resourceUserSchemaCreate(d *schema.ResourceData, m interface{}) error {
	id := d.Get("index").(string)
	if err := userCustomSchemaTemplate(d, m); err != nil {
		return err
	}
	d.SetId(id)

	return nil
}

func resourceUserSchemaExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := getClientFromMetadata(m)
	id := d.Get("index").(string)

	subschemas, _, err := client.Schemas.GetUserSubSchemaIndex(customSchema)
	if err != nil {
		return false, fmt.Errorf("Error Listing User Subschemas in Okta: %v", err)
	}

	return contains(subschemas, id), nil
}

func resourceUserSchemaRead(d *schema.ResourceData, m interface{}) error {
	client := getClientFromMetadata(m)
	schema, _, err := client.Schemas.GetUserSchema()
	if err != nil {
		return err
	}
	subschema := getSubSchema(schema.Definitions.Custom.Properties, d.Id())
	d.Set("array_type", subschema.Items.Type)
	d.Set("title", subschema.Title)
	d.Set("type", subschema.Type)
	d.Set("description", subschema.Description)
	d.Set("required", subschema.Required)
	d.Set("min_length", subschema.MinLength)
	d.Set("max_length", subschema.MaxLength)
	d.Set("enum", subschema.Enum)
	d.Set("one_of", subschema.OneOf)

	return setNonPrimitives(d, map[string]interface{}{
		"enum":   subschema.Enum,
		"one_of": flattenOneOf(subschema.OneOf),
	})
}

func getSubSchema(props []articulateOkta.CustomSubSchema, id string) *articulateOkta.CustomSubSchema {
	for _, schema := range props {
		if schema.Index == id {
			return &schema
		}
	}
	return nil
}

func resourceUserSchemaUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Update User Schema %v", d.Id())

	d.Partial(true)
	if err := userCustomSchemaTemplate(d, m); err != nil {
		return err
	}
	d.Partial(false)

	return resourceUserRead(d, m)
}

func resourceUserSchemaDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	log.Printf("[INFO] Delete User Schema %v", id)
	client := getClientFromMetadata(m)

	_, _, err := client.Schemas.DeleteUserCustomSubSchema(id)
	if err != nil {
		return err
	}

	return resourceUserRead(d, m)
}

// create or modify a custom subschema
func userCustomSchemaTemplate(d *schema.ResourceData, m interface{}) error {
	client := getClientFromMetadata(m)

	template := &articulateOkta.CustomSubSchema{
		Index:       d.Id(),
		Title:       d.Get("title").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Required:    d.Get("required").(bool),
		MinLength:   d.Get("min_length").(int),
		MaxLength:   d.Get("max_length").(int),
		Permissions: []articulateOkta.Permissions{
			{
				Action:    d.Get("permissions").(string),
				Principal: "SELF",
			},
		},
		Enum: convertInterfaceToStringArrNullable(d.Get("enum")),
	}
	template.Master.Type = d.Get("master").(string)
	template.Items.Type = d.Get("array_type").(string)

	if _, ok := d.GetOk("one_of"); ok {
		var obj interface{}

		// Will never error, we validate
		json.Unmarshal([]byte(d.Get("one_of").(string)), &obj)
		for _, v := range obj.([]interface{}) {
			oneOf := client.Schemas.OneOf()
			for k2, v2 := range v.(map[string]interface{}) {
				switch k2 {
				case "const":
					oneOf.Const = v2.(string)
				case "title":
					oneOf.Title = v2.(string)
				}
			}
			template.OneOf = append(template.OneOf, oneOf)
		}
	}

	_, _, err := client.Schemas.UpdateUserCustomSubSchema(*template)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Custom User Subschema in Okta: %v", err)
	}

	return nil
}

func flattenOneOf(oneOf []articulateOkta.OneOf) []map[string]interface{} {
	result := make([]map[string]interface{}, len(oneOf))
	for i, v := range oneOf {
		result[i] = map[string]interface{}{
			"const": v.Const,
			"title": v.Title,
		}
	}
	return result
}
