package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
)

var translationResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"language": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"subject": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"template": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	},
}

func resourceTemplateEmail() *schema.Resource {
	return &schema.Resource{
		Create: resourceTemplateEmailCreate,
		Exists: resourceTemplateEmailExists,
		Read:   resourceTemplateEmailRead,
		Update: resourceTemplateEmailUpdate,
		Delete: resourceTemplateEmailDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"default_language": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "en",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Email template type",
			},
			"translations": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     translationResource,
			},
		},
	}
}

func buildEmailTemplate(d *schema.ResourceData) *sdk.EmailTemplate {
	trans := map[string]*sdk.EmailTranslation{}
	rawTransList := d.Get("translations").(*schema.Set)

	for _, val := range rawTransList.List() {
		rawTrans := val.(map[string]interface{})
		trans[rawTrans["language"].(string)] = &sdk.EmailTranslation{
			Subject:  rawTrans["subject"].(string),
			Template: rawTrans["template"].(string),
		}
	}
	defaultLang := d.Get("default_language").(string)

	return &sdk.EmailTemplate{
		DefaultLanguage: defaultLang,
		Name:            "Custom",
		Type:            d.Get("type").(string),
		Translations:    trans,
		Subject:         trans[defaultLang].Subject,
		Template:        trans[defaultLang].Template,
	}
}

func flattenEmailTranlations(temp map[string]*sdk.EmailTranslation) *schema.Set {
	rawSet := []interface{}{}

	for key, val := range temp {
		rawSet = append(rawSet, map[string]interface{}{
			"language": key,
			"subject":  val.Subject,
			"template": val.Template,
		})
	}

	return schema.NewSet(schema.HashResource(translationResource), rawSet)
}

func resourceTemplateEmailCreate(d *schema.ResourceData, m interface{}) error {
	temp := buildEmailTemplate(d)
	id := d.Get("type").(string)
	_, _, err := getSupplementFromMetadata(m).CreateEmailTemplate(id, *temp, nil)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceTemplateEmailRead(d, m)
}

func resourceTemplateEmailExists(d *schema.ResourceData, m interface{}) (bool, error) {
	temp, resp, err := getSupplementFromMetadata(m).GetEmailTemplate(d.Id())

	return temp != nil && !is404(resp.StatusCode), err
}

func resourceTemplateEmailRead(d *schema.ResourceData, m interface{}) error {
	temp, resp, err := getSupplementFromMetadata(m).GetEmailTemplate(d.Id())

	if is404(resp.StatusCode) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	d.Set("translations", flattenEmailTranlations(temp.Translations))
	d.Set("default_language", temp.DefaultLanguage)

	return nil
}

func resourceTemplateEmailUpdate(d *schema.ResourceData, m interface{}) error {
	temp := buildEmailTemplate(d)
	_, _, err := getSupplementFromMetadata(m).UpdateEmailTemplate(d.Id(), *temp, nil)
	if err != nil {
		return err
	}

	return resourceTemplateEmailRead(d, m)
}

func resourceTemplateEmailDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getSupplementFromMetadata(m).DeleteEmailTemplate(d.Id())

	return err
}
