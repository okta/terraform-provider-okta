package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceEmailCustomization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEmailCustomizationCreate,
		ReadContext:   resourceEmailCustomizationRead,
		UpdateContext: resourceEmailCustomizationUpdate,
		DeleteContext: resourceEmailCustomizationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: emailCustomizationResourceSchema,
	}
}

func resourceEmailCustomizationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	brandID, ok := d.GetOk("brand_id")
	if !ok {
		return diag.Errorf("brand_id required to create email customization")
	}

	templateName, ok := d.GetOk("template_name")
	if !ok {
		return diag.Errorf("template name required to create email customization")
	}

	etcr := okta.EmailTemplateCustomizationRequest{}
	if language, ok := d.GetOk("language"); ok {
		etcr.Language = language.(string)
	}
	if isDefault, ok := d.GetOk("is_default"); ok {
		etcr.IsDefault = boolPtr(isDefault.(bool))
	}
	if subject, ok := d.GetOk("subject"); ok {
		etcr.Subject = subject.(string)
	}
	if body, ok := d.GetOk("body"); ok {
		etcr.Body = body.(string)
	}

	customization, _, err := getOktaClientFromMetadata(m).Brand.CreateEmailTemplateCustomization(ctx, brandID.(string), templateName.(string), etcr)
	if err != nil {
		return diag.Errorf("failed to create email customization: %v", err)
	}

	d.SetId(customization.Id)
	rawMap := flattenEmailCustomization(customization)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set new email customization properties: %v", err)
	}

	return nil
}

func resourceEmailCustomizationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	etcr, diagErr := etcrValues("read", d)
	if diagErr != nil {
		return diagErr
	}

	customization, resp, err := getOktaClientFromMetadata(m).Brand.GetEmailTemplateCustomization(ctx, etcr.brandID, etcr.templateName, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get email customization: %v", err)
	}
	if customization == nil {
		d.SetId("")
		return nil
	}

	rawMap := flattenEmailCustomization(customization)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set email customization properties: %v", err)
	}

	return nil
}

func resourceEmailCustomizationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	etcr, diagErr := etcrValues("update", d)
	if diagErr != nil {
		return diagErr
	}

	cr := okta.EmailTemplateCustomizationRequest{}
	if language, ok := d.GetOk("language"); ok {
		cr.Language = language.(string)
	}
	if isDefault, ok := d.GetOk("is_default"); ok {
		cr.IsDefault = boolPtr(isDefault.(bool))
	}
	if subject, ok := d.GetOk("subject"); ok {
		cr.Subject = subject.(string)
	}
	if body, ok := d.GetOk("body"); ok {
		cr.Body = body.(string)
	}

	customization, _, err := getOktaClientFromMetadata(m).Brand.UpdateEmailTemplateCustomization(ctx, etcr.brandID, etcr.templateName, d.Id(), cr)
	if err != nil {
		return diag.Errorf("failed to update email customization: %v", err)
	}

	d.SetId(customization.Id)
	rawMap := flattenEmailCustomization(customization)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set email customization properties: %v", err)
	}

	return nil
}

func resourceEmailCustomizationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	etcr, diagErr := etcrValues("delete", d)
	if diagErr != nil {
		return diagErr
	}

	_, err := getOktaClientFromMetadata(m).Brand.DeleteEmailTemplateCustomization(ctx, etcr.brandID, etcr.templateName, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete email customization: %v", err)
	}

	return nil
}

type etcrHelper struct {
	brandID      string
	templateName string
}

func etcrValues(action string, d *schema.ResourceData) (*etcrHelper, diag.Diagnostics) {
	brandID, ok := d.GetOk("brand_id")
	if !ok {
		return nil, diag.Errorf("brand_id required to %s email customization", action)
	}

	templateName, ok := d.GetOk("template_name")
	if !ok {
		return nil, diag.Errorf("template name required to %s email customization", action)
	}

	return &etcrHelper{
		brandID:      brandID.(string),
		templateName: templateName.(string),
	}, nil
}
