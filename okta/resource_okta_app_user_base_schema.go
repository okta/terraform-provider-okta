package okta

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func resourceAppUserBaseSchema() *schema.Resource {
	return &schema.Resource{
		Create:   resourceAppUserBaseSchemaCreate,
		Read:     resourceAppUserBaseSchemaRead,
		Update:   resourceAppUserBaseSchemaUpdate,
		Delete:   resourceAppUserBaseSchemaDelete,
		Exists:   resourceAppUserBaseSchemaExists,
		Importer: createNestedResourceImporter([]string{"app_id", "id"}),

		Schema: buildBaseUserSchema(map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func resourceAppUserBaseSchemaCreate(d *schema.ResourceData, m interface{}) error {
	if err := updateAppUserBaseSubschema(d, m); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s/%s", d.Get("app_id").(string), d.Get("index").(string)))

	return resourceAppUserBaseSchemaRead(d, m)
}

func resourceAppUserBaseSchemaExists(d *schema.ResourceData, m interface{}) (bool, error) {
	subschema, err := getAppUserBaseSubSchema(d, m)

	return subschema != nil, err
}

func resourceAppUserBaseSchemaRead(d *schema.ResourceData, m interface{}) error {
	subschema, err := getAppUserBaseSubSchema(d, m)
	if err != nil {
		return err
	} else if subschema == nil {
		d.SetId("")
		return nil
	}

	syncBaseUserSchema(d, subschema)

	return nil
}

func getAppUserBaseSubSchema(d *schema.ResourceData, m interface{}) (*sdk.UserSubSchema, error) {
	us, _, err := getSupplementFromMetadata(m).GetAppUserSchema(d.Get("app_id").(string))
	if err != nil {
		return nil, err
	}
	return getBaseProperty(us, d.Get("index").(string)), nil
}

func resourceAppUserBaseSchemaUpdate(d *schema.ResourceData, m interface{}) error {
	if err := updateAppUserBaseSubschema(d, m); err != nil {
		return err
	}

	return resourceAppUserBaseSchemaRead(d, m)
}

// can't delete Base
func resourceAppUserBaseSchemaDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

// create or modify a  subschema
func updateAppUserBaseSubschema(d *schema.ResourceData, m interface{}) error {
	subSchema := &sdk.UserSubSchema{
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

	_, _, err := getSupplementFromMetadata(m).UpdateBaseAppUserSchemaProperty(
		d.Get("index").(string),
		d.Get("app_id").(string),
		subSchema,
	)

	return err
}
