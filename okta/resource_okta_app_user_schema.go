package okta

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
)

func resourceAppUserSchema() *schema.Resource {
	return &schema.Resource{
		Create:   resourceAppUserSchemaCreate,
		Read:     resourceAppUserSchemaRead,
		Update:   resourceAppUserSchemaUpdate,
		Delete:   resourceAppUserSchemaDelete,
		Exists:   resourceAppUserSchemaExists,
		Importer: createNestedResourceImporter([]string{"app_id", "id"}),

		Schema: buildCustomUserSchema(map[string]*schema.Schema{
			"app_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func resourceAppUserSchemaCreate(d *schema.ResourceData, m interface{}) error {
	if err := updateAppUserSubschema(d, m); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s/%s", d.Get("app_id").(string), d.Get("index").(string)))

	return resourceAppUserSchemaRead(d, m)
}

func resourceAppUserSchemaExists(d *schema.ResourceData, m interface{}) (bool, error) {
	subschema, err := getAppUserSubSchema(d, m)

	return subschema != nil, err
}

func resourceAppUserSchemaRead(d *schema.ResourceData, m interface{}) error {
	subschema, err := getAppUserSubSchema(d, m)
	if err != nil {
		return err
	} else if subschema == nil {
		d.SetId("")
		return nil
	}

	return syncUserSchema(d, subschema)
}

func resourceAppUserSchemaUpdate(d *schema.ResourceData, m interface{}) error {
	if err := updateAppUserSubschema(d, m); err != nil {
		return err
	}

	return resourceAppUserSchemaRead(d, m)
}

func resourceAppUserSchemaDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getSupplementFromMetadata(m).DeleteAppUserSchemaProperty(d.Get("index").(string), d.Get("app_id").(string))

	return err
}

func getAppUserSubSchema(d *schema.ResourceData, m interface{}) (subschema *sdk.UserSubSchema, err error) {
	var schema *sdk.UserSchema

	schema, _, err = getSupplementFromMetadata(m).GetAppUserSchema(d.Get("app_id").(string))
	if err != nil {
		return
	}

	subschema = getCustomProperty(schema, d.Get("index").(string))
	return
}

func updateAppUserSubschema(d *schema.ResourceData, m interface{}) error {
	_, _, err := getSupplementFromMetadata(m).UpdateCustomAppUserSchemaProperty(
		d.Get("index").(string),
		d.Get("app_id").(string),
		getUserSubSchema(d),
	)

	return err
}
