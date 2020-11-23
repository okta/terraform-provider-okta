package okta

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				resourceIndex := d.Id()
				resourceUserType := "default"
				if strings.Contains(d.Id(), ".") {
					resourceUserType = strings.Split(d.Id(), ".")[0]
					resourceIndex = strings.Split(d.Id(), ".")[1]
				}

				d.SetId(resourceIndex)
				_ = d.Set("index", resourceIndex)
				_ = d.Set("user_type", resourceUserType)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: userSchemaSchema,
	}
}

func resourceUserSchemaCreate(d *schema.ResourceData, m interface{}) error {
	schemaUrl, err := getUserTypeSchemaUrl(m, d.Get("user_type").(string))
	if err != nil {
		return err
	}

	if err := updateSubschema(schemaUrl, d, m); err != nil {
		return err
	}
	d.SetId(d.Get("index").(string))

	return resourceUserSchemaRead(d, m)
}

func resourceUserSchemaExists(d *schema.ResourceData, m interface{}) (bool, error) {
	schemaUrl, err := getUserTypeSchemaUrl(m, d.Get("user_type").(string))
	if err != nil {
		return false, err
	}
	subschema, err := getSubSchema(schemaUrl, d, m)

	return subschema != nil, err
}

func resourceUserSchemaRead(d *schema.ResourceData, m interface{}) error {
	schemaUrl, err := getUserTypeSchemaUrl(m, d.Get("user_type").(string))
	if err != nil {
		return err
	}

	subschema, err := getSubSchema(schemaUrl, d, m)
	if err != nil {
		return err
	} else if subschema == nil {
		d.SetId("")
		return nil
	}

	return syncUserSchema(d, subschema)
}

func getSubSchema(schemaUrl string, d *schema.ResourceData, m interface{}) (*sdk.UserSubSchema, error) {
	s, _, err := getSupplementFromMetadata(m).GetUserSchema(schemaUrl)
	if err != nil {
		return nil, err
	}
	return getCustomProperty(s, d.Id()), err
}

func resourceUserSchemaUpdate(d *schema.ResourceData, m interface{}) error {
	schemaUrl, err := getUserTypeSchemaUrl(m, d.Get("user_type").(string))
	if err != nil {
		return err
	}
	if err := updateSubschema(schemaUrl, d, m); err != nil {
		return err
	}
	return resourceUserSchemaRead(d, m)
}

func resourceUserSchemaDelete(d *schema.ResourceData, m interface{}) error {
	schemaUrl, err := getUserTypeSchemaUrl(m, d.Get("user_type").(string))
	if err != nil {
		return err
	}
	_, err = getSupplementFromMetadata(m).DeleteUserSchemaProperty(schemaUrl, d.Id())
	return err
}

// create or modify a custom subschema
func updateSubschema(schemaUrl string, d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	_, _, err := client.UpdateCustomUserSchemaProperty(schemaUrl, d.Get("index").(string), getUserSubSchema(d))

	return err
}
