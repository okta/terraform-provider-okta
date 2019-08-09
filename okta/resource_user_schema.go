package okta

import (
	"fmt"
	"time"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
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
		DeprecationMessage: "This resource is now deprecated, please use okta_user_schema_object",
		Schema:             userSchemaSchema,
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
		// See https://github.com/articulate/terraform-provider-okta/issues/144
		// Occassionally a schema prop would be created and the read would not find it.
		// There appears to be a delay of availability on the Okta side, thus the backoff.
		fmt.Println("Could not find an existence subschema property, backing off and retrying. Known timing issue")
		time.Sleep(time.Second * 3)

		subschema, err = getSubSchema(d, m)
		if err != nil {
			return err
		} else if subschema == nil {
			return fmt.Errorf("Okta did not return a subschema for \"%s\"", d.Id())
		}
	}

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

func getSubSchema(d *schema.ResourceData, m interface{}) (subschema *articulateOkta.CustomSubSchema, err error) {
	var schema *articulateOkta.Schema
	id := d.Id()

	client := getClientFromMetadata(m)
	schema, _, err = client.Schemas.GetUserSchema()
	if err != nil {
		return
	}

	for _, part := range schema.Definitions.Custom.Properties {
		if part.Index == id {
			subschema = &part
			return
		}
	}

	return
}

func resourceUserSchemaUpdate(d *schema.ResourceData, m interface{}) error {
	if err := updateSubschema(d, m); err != nil {
		return err
	}

	return resourceUserSchemaRead(d, m)
}

func resourceUserSchemaDelete(d *schema.ResourceData, m interface{}) error {
	client := getClientFromMetadata(m)
	_, _, err := client.Schemas.DeleteUserCustomSubSchema(d.Id())

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
