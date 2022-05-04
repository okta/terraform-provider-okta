package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func dataSourceEmailCustomization() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEmailCustomizationRead,
		Schema:      emailCustomizationDataSourceSchema,
	}
}

func dataSourceEmailCustomizationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var brand *okta.Brand
	var err error
	brandID, ok := d.GetOk("brand_id")
	if ok {
		logger(m).Info("reading brand by ID", "id", brandID.(string))
		brand, _, err = getOktaClientFromMetadata(m).Brand.GetBrand(ctx, brandID.(string))
		if err != nil {
			return diag.Errorf("failed to get brand for email template: %v", err)
		}
	} else {
		return diag.Errorf("brand_id required for email template: %v", err)
	}

	templateName, ok := d.GetOk("template_name")
	if !ok {
		return diag.Errorf("template name required for email template: %v", err)
	}

	customizationId, ok := d.GetOk("customization_id")
	if !ok {
		return diag.Errorf("customization_id required for email customization: %v", err)
	}

	customization, _, err := getOktaClientFromMetadata(m).Brand.GetEmailTemplateCustomization(ctx, brandID.(string), templateName.(string), customizationId.(string))
	if err != nil {
		return diag.Errorf("failed to get email template: %v", err)
	}

	d.SetId(fmt.Sprintf("email_customization-%s-%s-%s", customization.Id, templateName, brand.Id))
	rawMap := flattenEmailCustomization(brandID.(string), templateName.(string), customization)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set email customization properties: %v", err)
	}

	return nil
}
