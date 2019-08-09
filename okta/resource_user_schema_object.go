package okta

import (
	"strings"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/schema"
)

uuserSearchResource := Elem: &schema.Resource{
	Schema: userSchemaSchema,
}

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
			"base": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: userSchemaResource,
			},
			"custom": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: userSchemaResource,
			},
		},
	}
}

// Builds user schema object id, which is just a sorted list of strings
func buildUserSchemaObjectId(d *schema.ResourceData) string {
	indexList := []string{}

	if cust, ok := d.GetOk("custom"); ok {
		customList := cust.(*schema.Set).List()
		indexList = make([]string, len(customList))

		for i, item := range customList {
			mp := item.(map[string]interface{})
			indexList[i] = mp["index"].(string)
		}
	}

	return strings.Join(indexList, "/")
}

func resourceUserSchemaObjectCreate(d *schema.ResourceData, m interface{}) error {
	// No need for an actual id so just generate one
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	d.SetId(id.String())

	return resourceUserSchemaRead(d, m)
}

func resourceUserSchemaObjectRead(d *schema.ResourceData, m interface{}) error {
	schm, err := getUserSchemaObject(d, m)
	if err != nil {
		return err
	}

	d.Set("custom", flattenUserSchemaObject(schm.Definitions.Custom.Properties))

	return nil
}

func resourceUserSchemaObjectUpdate(d *schema.ResourceData, m interface{}) error {
	updateSubSchema(d.Get(""))
	return resourceUserSchemaRead(d, m)
}

func resourceUserSchemaObjectDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func flattenUserSchemaObject(props []*articulateOkta.CustomSubSchema) *schema.Set {
	propSet := schema.NewSet(schema.HashResource(userSearchResource), []interface{}{})

	for _, prop := range props {
		propSet.Add(map[string]interface{}{
			"index": prop.Index,
			"array_type": prop.Items.Type,
			"title": prop.Title,
			"type": prop.Type,
			"description": prop.Description,
			"required": prop.Required,
			"min_length": prop.MinLength,
			"max_length": prop.MaxLength,
			"enum": prop.Enum,
			"one_of": flattenOneOf(subschema.OneOf),
			"permissions": prop.Permissions,
			"master": prop.Master,
		})
	}

	return propSet
}

func getUserSchemaObject(d *schema.ResourceData, m interface{}) (schemaObj *articulateOkta.Schema, err error) {
	client := getClientFromMetadata(m)
	schemaObj, _, err = client.Schemas.GetUserSchema()
	if err != nil {
		return
	}

	return
}
