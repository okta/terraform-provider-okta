package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v3/okta"
)

func resourceEmailCustomization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEmailCustomizationCreate,
		ReadContext:   resourceEmailCustomizationRead,
		UpdateContext: resourceEmailCustomizationUpdate,
		DeleteContext: resourceEmailCustomizationDelete,
		Importer:      createNestedResourceImporter([]string{"id", "brand_id", "template_name"}),
		Schema:        emailCustomizationResourceSchema,
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

	etcr := okta.EmailCustomization{}
	if language, ok := d.GetOk("language"); ok {
		etcr.Language = language.(string)
	}
	if isDefault, ok := d.GetOk("is_default"); ok {
		etcr.IsDefault = boolPtr(isDefault.(bool))
	} else {
		etcr.IsDefault = boolPtr(false)
	}
	if subject, ok := d.GetOk("subject"); ok {
		etcr.Subject = subject.(string)
	}
	if body, ok := d.GetOk("body"); ok {
		etcr.Body = body.(string)
	}

	client := getOktaV3ClientFromMetadata(m)

	customization, _, err := client.CustomizationAPI.CreateEmailCustomization(ctx, brandID.(string), templateName.(string)).Instance(etcr).Execute()
	if err != nil {
		return diag.Errorf("failed to create email customization: %v", err)
	}

	d.SetId(customization.GetId())
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

	customization, resp, err := getOktaV3ClientFromMetadata(m).CustomizationAPI.GetEmailCustomization(ctx, etcr.brandID, etcr.templateName, d.Id()).Execute()
	if err := v3suppressErrorOn404(resp, err); err != nil {
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

	cr := okta.EmailCustomization{}
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

	customization, _, err := getOktaV3ClientFromMetadata(m).CustomizationAPI.ReplaceEmailCustomization(ctx, etcr.brandID, etcr.templateName, d.Id()).Instance(cr).Execute()
	if err != nil {
		return diag.Errorf("failed to update email customization: %v", err)
	}

	d.SetId(customization.GetId())
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

	client := getOktaV3ClientFromMetadata(m)
	// If this is the last customization template call the delete all endpoint
	// as the API doesn't allow deleting the last template explicitly should the
	// template be the default.  "Returns a 409 Conflict if the email
	// customization to be deleted is the default."
	// https://developer.okta.com/docs/reference/api/brands/#response-body-23
	// Else delete the specific customization.
	customizations, _, err := client.CustomizationAPI.ListEmailCustomizations(ctx, etcr.brandID, etcr.templateName).Execute()
	if err != nil {
		return diag.Errorf("failed to delete email customization: %v", err)
	}
	if len(customizations) == 1 {
		_, err := client.CustomizationAPI.DeleteAllCustomizations(ctx, etcr.brandID, etcr.templateName).Execute()
		if err != nil {
			return diag.Errorf("failed to delete email customization: %v", err)
		}
		return nil
	}

	_, err = client.CustomizationAPI.DeleteEmailCustomization(ctx, etcr.brandID, etcr.templateName, d.Id()).Execute()
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
