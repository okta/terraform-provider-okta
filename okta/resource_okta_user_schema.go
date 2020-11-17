package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
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
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				_ = d.Set("index", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: userSchemaSchema,
	}
}

func resourceUserSchemaCreate(d *schema.ResourceData, m interface{}) error {
	if err := updateSubschema(d, m); err != nil {
		return err
	}
	d.SetId(d.Get("index").(string))

	return resourceUserSchemaRead(d, m)
}

func resourceUserSchemaExists(d *schema.ResourceData, m interface{}) (bool, error) {
	subschema, err := getSubSchema(d, m)

	return subschema != nil, err
}

func resourceUserSchemaRead(d *schema.ResourceData, m interface{}) error {
	subschema, err := getSubSchema(d, m)
	if err != nil {
		return err
	} else if subschema == nil {
		d.SetId("")
		return nil
	}

	return syncUserSchema(d, subschema)
}

func getSubSchema(d *schema.ResourceData, m interface{}) (*sdk.UserSubSchema, error) {
	s, _, err := getSupplementFromMetadata(m).GetUserSchema()
	if err != nil {
		return nil, err
	}
	return getCustomProperty(s, d.Id()), err
}

func resourceUserSchemaUpdate(d *schema.ResourceData, m interface{}) error {
	if err := updateSubschema(d, m); err != nil {
		return err
	}

	return resourceUserSchemaRead(d, m)
}

func resourceUserSchemaDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getSupplementFromMetadata(m).DeleteUserSchemaProperty(d.Id())

	return err
}

// create or modify a custom subschema
func updateSubschema(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	_, _, err := client.UpdateCustomUserSchemaProperty(d.Get("index").(string), getUserSubSchema(d))

	return err
}
