package okta

import (
	"fmt"

	"github.com/articulate/terraform-provider-okta/sdk"
	"github.com/hashicorp/terraform/helper/schema"
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
		return fmt.Errorf("Okta did not return a subschema for \"%s\". This is a known limitation of Okta's API, these must be created one at a time to avoid overwritting. One way to do this is via depends_on, see link for example https://github.com/articulate/terraform-provider-okta/blob/master/examples/okta_user/custom_attributes.tf.", d.Id())
	}

	return syncUserSchema(d, subschema)
}

func getSubSchema(d *schema.ResourceData, m interface{}) (subschema *sdk.UserSubSchema, err error) {
	var schema *sdk.UserSchema

	schema, _, err = getSupplementFromMetadata(m).GetUserSchema()
	if err != nil {
		return
	}

	subschema = getCustomProperty(schema, d.Id())
	return
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
