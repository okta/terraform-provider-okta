package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

var translationSmsResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"language": {
			Type:     schema.TypeString,
			Required: true,
		},
		"template": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 161),
		},
	},
}

func resourceTemplateSms() *schema.Resource {
	return &schema.Resource{
		Create: resourceTemplateSmsCreate,
		Exists: resourceTemplateSmsExists,
		Read:   resourceTemplateSmsRead,
		Update: resourceTemplateSmsUpdate,
		Delete: resourceTemplateSmsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SMS template type",
			},
			"template": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "SMS default template",
				ValidateFunc: validation.StringLenBetween(1, 161),
			},
			"translations": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     translationSmsResource,
			},
		},
	}
}

func buildSmsTemplate(d *schema.ResourceData) *sdk.SmsTemplate {
	trans := map[string]string{}
	rawTransList := d.Get("translations").(*schema.Set)

	for _, val := range rawTransList.List() {
		rawTrans := val.(map[string]interface{})
		trans[rawTrans["language"].(string)] = rawTrans["template"].(string)
	}

	return &sdk.SmsTemplate{
		Name:         "Custom",
		Type:         d.Get("type").(string),
		Translations: trans,
		Template:     d.Get("template").(string),
	}
}

func flattenSmsTranlations(temp map[string]string) *schema.Set {
	rawSet := []interface{}{}

	for key, val := range temp {
		rawSet = append(rawSet, map[string]interface{}{
			"language": key,
			"template": val,
		})
	}

	return schema.NewSet(schema.HashResource(translationSmsResource), rawSet)
}

func resourceTemplateSmsCreate(d *schema.ResourceData, m interface{}) error {
	temp := buildSmsTemplate(d)
	response, _, err := getSupplementFromMetadata(m).CreateSmsTemplate(*temp, nil)
	if err != nil {
		return err
	}

	d.SetId(response.Id)

	return resourceTemplateSmsRead(d, m)
}

func resourceTemplateSmsExists(d *schema.ResourceData, m interface{}) (bool, error) {
	temp, resp, err := getSupplementFromMetadata(m).GetSmsTemplate(d.Id())

	return temp != nil && !is404(resp.StatusCode), err
}

func resourceTemplateSmsRead(d *schema.ResourceData, m interface{}) error {
	temp, resp, err := getSupplementFromMetadata(m).GetSmsTemplate(d.Id())

	if resp != nil && is404(resp.StatusCode) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	_ = d.Set("translations", flattenSmsTranlations(temp.Translations))

	return nil
}

func resourceTemplateSmsUpdate(d *schema.ResourceData, m interface{}) error {
	temp := buildSmsTemplate(d)
	_, _, err := getSupplementFromMetadata(m).UpdateSmsTemplate(d.Id(), *temp, nil)
	if err != nil {
		return err
	}

	return resourceTemplateSmsRead(d, m)
}

func resourceTemplateSmsDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getSupplementFromMetadata(m).DeleteSmsTemplate(d.Id())

	return err
}
