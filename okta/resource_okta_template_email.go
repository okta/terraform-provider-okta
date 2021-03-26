package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

var translationResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"language": {
			Type:     schema.TypeString,
			Required: true,
		},
		"subject": {
			Type:     schema.TypeString,
			Required: true,
		},
		"template": {
			Type:     schema.TypeString,
			Required: true,
		},
	},
}

func resourceTemplateEmail() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTemplateEmailCreate,
		ReadContext:   resourceTemplateEmailRead,
		UpdateContext: resourceTemplateEmailUpdate,
		DeleteContext: resourceTemplateEmailDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"default_language": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "en",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Email template type",
				ForceNew:    true,
			},
			"translations": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     translationResource,
			},
		},
	}
}

func resourceTemplateEmailCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	temp := buildEmailTemplate(d)
	id := d.Get("type").(string)
	_, _, err := getSupplementFromMetadata(m).CreateEmailTemplate(ctx, *temp, nil)
	if err != nil {
		return diag.Errorf("failed to create email template: %v", err)
	}
	d.SetId(id)
	return resourceTemplateEmailRead(ctx, d, m)
}

func resourceTemplateEmailRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	temp, resp, err := getSupplementFromMetadata(m).GetEmailTemplate(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get email template: %v", err)
	}
	if temp == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("translations", flattenEmailTranslations(temp.Translations))
	_ = d.Set("default_language", temp.DefaultLanguage)
	return nil
}

func resourceTemplateEmailUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	temp := buildEmailTemplate(d)
	_, _, err := getSupplementFromMetadata(m).UpdateEmailTemplate(ctx, d.Id(), *temp, nil)
	if err != nil {
		return diag.Errorf("failed to update email template: %v", err)
	}
	return resourceTemplateEmailRead(ctx, d, m)
}

func resourceTemplateEmailDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, err := getSupplementFromMetadata(m).DeleteEmailTemplate(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete email template: %v", err)
	}
	return nil
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

func flattenEmailTranslations(temp map[string]*sdk.EmailTranslation) *schema.Set {
	var rawSet []interface{}
	for key, val := range temp {
		rawSet = append(rawSet, map[string]interface{}{
			"language": key,
			"subject":  val.Subject,
			"template": val.Template,
		})
	}
	return schema.NewSet(schema.HashResource(translationResource), rawSet)
}
