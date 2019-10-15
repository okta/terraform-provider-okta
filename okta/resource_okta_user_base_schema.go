package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
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

		Schema: userBaseSchemaSchema,
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
		d.SetId("")
		return nil
	}

	syncBaseUserSchema(d, subschema)

	return nil
}

func getBaseSubSchema(d *schema.ResourceData, m interface{}) (subschema *sdk.UserSubSchema, err error) {
	var schema *sdk.UserSchema

	schema, _, err = getSupplementFromMetadata(m).GetUserSchema()
	if err != nil {
		return
	}

	subschema = getBaseProperty(schema, d.Id())
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
	schema := &sdk.UserSubSchema{
		Master: getNullableMaster(d),
		Title:  d.Get("title").(string),
		Type:   d.Get("type").(string),
		Permissions: []*sdk.UserSchemaPermission{
			{
				Action:    d.Get("permissions").(string),
				Principal: "SELF",
			},
		},
		Required: boolPtr(d.Get("required").(bool)),
	}

	_, _, err := getSupplementFromMetadata(m).UpdateBaseUserSchemaProperty(d.Get("index").(string), schema)

	return err
}
