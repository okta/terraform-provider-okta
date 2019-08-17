package okta

import (
	"fmt"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const baseSchema = "base"

func resourceUserBaseSchema() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserBaseSchemaCreate,
		Read:   resourceUserBaseSchemaRead,
		Update: resourceUserBaseSchemaUpdate,
		Delete: resourceUserBaseSchemaDelete,
		Exists: resourceUserBaseSchemaExists,
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
				ValidateFunc: validation.StringInSlice([]string{"string", "number", "integer", "reference"}, false),
				Description:  "Subschema array type: string, number, integer, reference. Type field must be an array.",
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
		},
	}
}

func resourceUserBaseSchemaCreate(d *schema.ResourceData, m interface{}) error {
	if err := updateBaseSubschema(d, m); err != nil {
		return err
	}
	d.SetId(d.Get("index").(string))

	return resourceUserBaseSchemaRead(d, m)
}

func resourceUserBaseSchemaExists(d *schema.ResourceData, m interface{}) (bool, error) {
	subschema, err := getBaseSubSchema(d, m)

	return subschema != nil, err
}

func resourceUserBaseSchemaRead(d *schema.ResourceData, m interface{}) error {
	subschema, err := getBaseSubSchema(d, m)
	if err != nil {
		return err
	} else if subschema == nil {
		subschema, err = getBaseSubSchema(d, m)
		if err != nil {
			return err
		} else if subschema == nil {
			return fmt.Errorf("Okta did not return a subschema for \"%s\"", d.Id())
		}
	}

	d.Set("title", subschema.Title)
	d.Set("type", subschema.Type)
	d.Set("index", subschema.Index)

	if subschema.Master != nil {
		d.Set("master", subschema.Master.Type)
	}

	if len(subschema.Permissions) > 0 {
		d.Set("permissions", subschema.Permissions[0].Action)
	}

	return nil
}

func getBaseSubSchema(d *schema.ResourceData, m interface{}) (subschema *articulateOkta.BaseSubSchema, err error) {
	var schema *articulateOkta.Schema
	id := d.Id()

	client := getClientFromMetadata(m)
	schema, _, err = client.Schemas.GetUserSchema()
	if err != nil {
		return
	}

	for _, part := range schema.Definitions.Base.Properties {
		if part.Index == id {
			subschema = &part
			return
		}
	}

	return
}

func resourceUserBaseSchemaUpdate(d *schema.ResourceData, m interface{}) error {
	if err := updateBaseSubschema(d, m); err != nil {
		return err
	}

	return resourceUserBaseSchemaRead(d, m)
}

// can't delete Base
func resourceUserBaseSchemaDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

// create or modify a  subschema
func updateBaseSubschema(d *schema.ResourceData, m interface{}) error {
	client := getClientFromMetadata(m)

	template := &articulateOkta.BaseSubSchema{
		Index: d.Get("index").(string),
		Title: d.Get("title").(string),
		Type:  d.Get("type").(string),
		Permissions: []articulateOkta.Permissions{
			{
				Action:    d.Get("permissions").(string),
				Principal: "SELF",
			},
		},
	}

	if v, ok := d.GetOk("master"); ok {
		template.Master = &articulateOkta.Master{Type: v.(string)}
	}

	_, _, err := client.Schemas.UpdateUserBaseSubSchema(*template)
	if err != nil {
		return fmt.Errorf("Error Updating User Base Subschema in Okta: %v", err)
	}

	return nil
}
