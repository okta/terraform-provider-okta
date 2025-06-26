package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func resourceEmailCustomization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEmailCustomizationCreate,
		ReadContext:   resourceEmailCustomizationRead,
		UpdateContext: resourceEmailCustomizationUpdate,
		DeleteContext: resourceEmailCustomizationDelete,
		Importer:      utils.CreateNestedResourceImporter([]string{"id", "brand_id", "template_name"}),
		Description: `Create an email customization of an email template belonging to a brand in an Okta organization.
		Use this resource to create an [email
		customization](https://developer.okta.com/docs/reference/api/brands/#create-email-customization)
		of an email template belonging to a brand in an Okta organization.
		~> Okta's public API is strict regarding the behavior of the 'is_default'
		property in [an email
		customization](https://developer.okta.com/docs/reference/api/brands/#email-customization).
		Make use of 'depends_on' meta argument to ensure the provider navigates email customization
		language versions seamlessly. Have all secondary customizations depend on the primary
		customization that is marked default. See [Example Usage](#example-usage).
		~> Caveats for [creating an email
		customization](https://developer.okta.com/docs/reference/api/brands/#response-body-19).
		If this is the first customization being created for the email template, and
		'is_default' is not set for the customization in its resource configuration, the
		API will respond with the created customization marked as default. The API will
		400 if the language parameter is not one of the supported languages or the body
		parameter does not contain a required variable reference. The API will error 409
		if 'is_default' is true and a default customization exists. The API will 404 for
		an invalid 'brand_id' or 'template_name'.
		~> Caveats for [updating an email
		customization](https://developer.okta.com/docs/reference/api/brands/#response-body-22).
		If the 'is_default' parameter is true, the previous default email customization
		has its 'is_default' set to false (see previous note about mitigating this with
		'depends_on' meta argument). The API will 409 if there’s already another email
		customization for the specified language or the 'is_default' parameter is false
		and the email customization being updated is the default. The API will 400 if
		the language parameter is not one of the supported locales or the body parameter
		does not contain a required variable reference.  The API will 404 for an invalid
		'brand_id' or 'template_name'.`,
		Schema: emailCustomizationResourceSchema,
	}
}

func resourceEmailCustomizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		etcr.IsDefault = utils.BoolPtr(isDefault.(bool))
	} else {
		etcr.IsDefault = utils.BoolPtr(false)
	}
	if subject, ok := d.GetOk("subject"); ok {
		etcr.Subject = subject.(string)
	}
	if body, ok := d.GetOk("body"); ok {
		etcr.Body = body.(string)
	}

	client := getOktaV3ClientFromMetadata(meta)

	customization, _, err := client.CustomizationAPI.CreateEmailCustomization(ctx, brandID.(string), templateName.(string)).Instance(etcr).Execute()
	if err != nil {
		return diag.Errorf("failed to create email customization: %v", err)
	}

	d.SetId(customization.GetId())
	rawMap := flattenEmailCustomization(customization)
	err = utils.SetNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set new email customization properties: %v", err)
	}

	return nil
}

func resourceEmailCustomizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	etcr, diagErr := etcrValues("read", d)
	if diagErr != nil {
		return diagErr
	}

	customization, resp, err := getOktaV3ClientFromMetadata(meta).CustomizationAPI.GetEmailCustomization(ctx, etcr.brandID, etcr.templateName, d.Id()).Execute()
	if err := utils.SuppressErrorOn404_V3(resp, err); err != nil {
		return diag.Errorf("failed to get email customization: %v", err)
	}
	if customization == nil {
		d.SetId("")
		return nil
	}

	rawMap := flattenEmailCustomization(customization)
	err = utils.SetNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set email customization properties: %v", err)
	}

	return nil
}

func resourceEmailCustomizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	etcr, diagErr := etcrValues("update", d)
	if diagErr != nil {
		return diagErr
	}

	cr := okta.EmailCustomization{}
	if language, ok := d.GetOk("language"); ok {
		cr.Language = language.(string)
	}
	if isDefault, ok := d.GetOk("is_default"); ok {
		cr.IsDefault = utils.BoolPtr(isDefault.(bool))
	}
	if subject, ok := d.GetOk("subject"); ok {
		cr.Subject = subject.(string)
	}
	if body, ok := d.GetOk("body"); ok {
		cr.Body = body.(string)
	}

	customization, _, err := getOktaV3ClientFromMetadata(meta).CustomizationAPI.ReplaceEmailCustomization(ctx, etcr.brandID, etcr.templateName, d.Id()).Instance(cr).Execute()
	if err != nil {
		return diag.Errorf("failed to update email customization: %v", err)
	}

	d.SetId(customization.GetId())
	rawMap := flattenEmailCustomization(customization)
	err = utils.SetNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set email customization properties: %v", err)
	}

	return nil
}

func resourceEmailCustomizationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	etcr, diagErr := etcrValues("delete", d)
	if diagErr != nil {
		return diagErr
	}

	client := getOktaV3ClientFromMetadata(meta)
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
