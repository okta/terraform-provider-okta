package okta

import (
	"fmt"

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
				ValidateFunc: validation.StringInSlice([]string{"string", "boolean", "number", "integer", "array", "object"}, false),
				Description:  "Subschema type: string, boolean, number, integer, array, or object",
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
		},
	}
}

func resourceUserSchemaCreate(d *schema.ResourceData, m interface{}) error {
	if err := updateSubschema(d, m); err != nil {
		return err
	}
	d.SetId(d.Get("index").(string))
	fmt.Println("CREATE --", d.Id())

	return resourceUserSchemaRead(d, m)
}

func resourceUserSchemaExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := getClientFromMetadata(m)
	subschemas, _, err := client.Schemas.GetUserSubSchemaIndex(customSchema)
	if err != nil {
		return false, fmt.Errorf("Error Listing User Subschemas in Okta: %v", err)
	}

	fmt.Println("EXISTS", d.Id(), " == ", subschemas)
	return contains(subschemas, d.Id()), nil
}

func resourceUserSchemaRead(d *schema.ResourceData, m interface{}) error {
	client := getClientFromMetadata(m)
	schema, _, err := client.Schemas.GetUserSchema()
	if err != nil {
		return err
	}
	fmt.Println("READ --", d.Id())
	subschema := getSubSchema(schema.Definitions.Custom.Properties, d.Id())
	d.Set("array_type", subschema.Items.Type)
	d.Set("title", subschema.Title)
	d.Set("type", subschema.Type)
	d.Set("description", subschema.Description)
	d.Set("required", subschema.Required)
	d.Set("index", subschema.Index)

	if subschema.Master != nil {
		d.Set("master", subschema.Master.Type)
	}

	if len(subschema.Permissions) > 0 {
		d.Set("permissions", subschema.Permissions[0].Action)
	}

	if subschema.MinLength > 0 {
		d.Set("min_length", subschema.MinLength)
	}

	if subschema.MaxLength > 0 {
		d.Set("max_length", subschema.MaxLength)
	}

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
	fmt.Println("UPDATE --", d.Id())
	if err := updateSubschema(d, m); err != nil {
		return err
	}

	return resourceUserSchemaRead(d, m)
}

func resourceUserSchemaDelete(d *schema.ResourceData, m interface{}) error {
	client := getClientFromMetadata(m)

	_, _, err := client.Schemas.DeleteUserCustomSubSchema(d.Id())
	fmt.Println("DELETE --", d.Id())
	return err
}

// create or modify a custom subschema
func updateSubschema(d *schema.ResourceData, m interface{}) error {
	client := getClientFromMetadata(m)

	template := &articulateOkta.CustomSubSchema{
		Index:       d.Get("index").(string),
		Title:       d.Get("title").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Required:    d.Get("required").(bool),
		Permissions: []articulateOkta.Permissions{
			{
				Action:    d.Get("permissions").(string),
				Principal: "SELF",
			},
		},
		Enum: convertInterfaceToStringArrNullable(d.Get("enum")),
	}

	if v, ok := d.GetOk("master"); ok {
		template.Master = &articulateOkta.Master{Type: v.(string)}
	}

	if v, ok := d.GetOk("array_type"); ok {
		template.Items.Type = v.(string)
	}

	if v, ok := d.GetOk("min_length"); ok {
		template.MinLength = v.(int)
	}

	if v, ok := d.GetOk("max_length"); ok {
		template.MaxLength = v.(int)
	}

	if oneOfList, ok := d.GetOk("one_of"); ok {
		for _, v := range oneOfList.([]interface{}) {
			valueMap := v.(map[string]interface{})
			template.OneOf = append(template.OneOf, articulateOkta.OneOf{
				Const: valueMap["const"].(string),
				Title: valueMap["title"].(string),
			})
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
