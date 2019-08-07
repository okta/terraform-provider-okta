package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUserSchemaObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserSchemaObjectCreate,
		Read:   resourceUserSchemaObjectRead,
		Update: resourceUserSchemaObjectUpdate,
		Delete: resourceUserSchemaObjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"custom": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: userSchemaSchema,
				},
			},
		},
	}
}

func resourceUserSchemaObjectCreate(d *schema.ResourceData, m interface{}) error {
	if err := updateSubschema(d, m); err != nil {
		return err
	}
	d.SetId(d.Get("index").(string))

	return resourceUserSchemaRead(d, m)
}

func resourceUserSchemaObjectRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceUserSchemaObjectUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceUserSchemaRead(d, m)
}

func resourceUserSchemaObjectDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
