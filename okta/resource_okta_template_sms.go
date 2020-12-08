package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

var translationSmsResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"language": {
			Type:     schema.TypeString,
			Required: true,
		},
		"template": {
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: stringLenBetween(1, 161),
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
				Type:             schema.TypeString,
				Required:         true,
				Description:      "SMS default template",
				ValidateDiagFunc: stringLenBetween(1, 161),
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
	response, _, err := getSupplementFromMetadata(m).CreateSmsTemplate(ctx, *temp, nil)
	if err != nil {
		return diag.Errorf("failed to create SMS template: %v", err)
	}
	d.SetId(response.Id)
	return resourceTemplateSmsRead(ctx, d, m)
}

func resourceTemplateSmsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	temp, resp, err := getSupplementFromMetadata(m).GetSmsTemplate(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get SMS template: %v", err)
	}
	if temp == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("translations", flattenSmsTranslations(temp.Translations))
	return nil
}

func resourceTemplateSmsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	temp := buildSmsTemplate(d)
	_, _, err := getSupplementFromMetadata(m).UpdateSmsTemplate(ctx, d.Id(), *temp, nil)
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

func flattenSmsTranslations(temp map[string]string) *schema.Set {
	var rawSet []interface{}
	for key, val := range temp {
		rawSet = append(rawSet, map[string]interface{}{
			"language": key,
			"template": val,
		})
	}
	return schema.NewSet(schema.HashResource(translationSmsResource), rawSet)
}
