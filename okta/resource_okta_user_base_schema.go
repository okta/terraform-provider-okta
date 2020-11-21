package okta

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
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
		Schema: userBaseSchemaSchema,
	}
}

func resourceUserBaseSchemaCreate(d *schema.ResourceData, m interface{}) error {
	schemaURL, err := getUserTypeSchemaURL(m, d.Get("user_type").(string))
	if err != nil {
		return err
	}

	if err := updateBaseSubschema(schemaURL, d, m); err != nil {
		return err
	}
	d.SetId(d.Get("index").(string))

	return resourceUserBaseSchemaRead(d, m)
}

func resourceUserBaseSchemaExists(d *schema.ResourceData, m interface{}) (bool, error) {
	schemaURL, err := getUserTypeSchemaURL(m, d.Get("user_type").(string))
	if err != nil {
		return false, err
	}
	subschema, err := getBaseSubSchema(schemaURL, d, m)

	return subschema != nil, err
}

func resourceUserBaseSchemaRead(d *schema.ResourceData, m interface{}) error {
	schemaURL, err := getUserTypeSchemaURL(m, d.Get("user_type").(string))
	if err != nil {
		return err
	}
	subschema, err := getBaseSubSchema(schemaURL, d, m)
	if err != nil {
		return err
	} else if subschema == nil {
		d.SetId("")
		return nil
	}

	syncBaseUserSchema(d, subschema)

	return nil
}

func getBaseSubSchema(schemaURL string, d *schema.ResourceData, m interface{}) (*sdk.UserSubSchema, error) {
	s, _, err := getSupplementFromMetadata(m).GetUserSchema(schemaURL)
	if err != nil {
		return nil, err
	}
	return getBaseProperty(s, d.Id()), err
}

func resourceUserBaseSchemaUpdate(d *schema.ResourceData, m interface{}) error {
	schemaURL, err := getUserTypeSchemaURL(m, d.Get("user_type").(string))
	if err != nil {
		return err
	}
	if err := updateBaseSubschema(schemaURL, d, m); err != nil {
		return err
	}
	return resourceUserBaseSchemaRead(d, m)
}

// can't delete Base
func resourceUserBaseSchemaDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

// create or modify a  subschema
func updateBaseSubschema(schemaURL string, d *schema.ResourceData, m interface{}) error {
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

	_, _, err := getSupplementFromMetadata(m).UpdateBaseUserSchemaProperty(schemaURL, d.Get("index").(string), subSchema)

	return err
}
