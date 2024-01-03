package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

var translationSmsResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"language": {
			Type:     schema.TypeString,
			Required: true,
		},
		"template": {
			Type:     schema.TypeString,
			Required: true,
		},
	},
}

func resourceTemplateSms() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTemplateSmsCreate,
		ReadContext:   resourceTemplateSmsRead,
		UpdateContext: resourceTemplateSmsUpdate,
		DeleteContext: resourceTemplateSmsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SMS template type",
			},
			"template": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SMS default template",
			},
			"translations": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     translationSmsResource,
			},
		},
	}
}

func resourceTemplateSmsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	temp := buildSmsTemplate(d)
	response, _, err := getOktaClientFromMetadata(m).SmsTemplate.CreateSmsTemplate(ctx, *temp)
	if err != nil {
		return diag.Errorf("failed to create SMS template: %v", err)
	}
	d.SetId(response.Id)
	return resourceTemplateSmsRead(ctx, d, m)
}

func resourceTemplateSmsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	temp, resp, err := getOktaClientFromMetadata(m).SmsTemplate.GetSmsTemplate(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get SMS template: %v", err)
	}
	if temp == nil {
		d.SetId("")
		return nil
	}
	if temp.Translations != nil {
		_ = d.Set("translations", flattenSmsTranslations(*temp.Translations))
	}
	return nil
}

func resourceTemplateSmsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	temp := buildSmsTemplate(d)
	_, _, err := getOktaClientFromMetadata(m).SmsTemplate.UpdateSmsTemplate(ctx, d.Id(), *temp)
	if err != nil {
		return diag.Errorf("failed to update SMS template: %v", err)
	}
	return resourceTemplateSmsRead(ctx, d, m)
}

func resourceTemplateSmsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, err := getOktaClientFromMetadata(m).SmsTemplate.DeleteSmsTemplate(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete SMS template: %v", err)
	}
	return nil
}

func buildSmsTemplate(d *schema.ResourceData) *sdk.SmsTemplate {
	trans := make(sdk.SmsTemplateTranslations)
	rawTransList := d.Get("translations").(*schema.Set).List()

	for _, val := range rawTransList {
		rawTrans := val.(map[string]interface{})
		trans[rawTrans["language"].(string)] = rawTrans["template"]
	}

	return &sdk.SmsTemplate{
		Name:         "Custom",
		Type:         d.Get("type").(string),
		Translations: &trans,
		Template:     d.Get("template").(string),
	}
}

func flattenSmsTranslations(temp sdk.SmsTemplateTranslations) *schema.Set {
	var rawSet []interface{}
	for key, val := range map[string]interface{}(temp) {
		rawSet = append(rawSet, map[string]interface{}{
			"language": key,
			"template": val,
		})
	}
	return schema.NewSet(schema.HashResource(translationSmsResource), rawSet)
}
