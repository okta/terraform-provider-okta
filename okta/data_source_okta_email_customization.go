package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEmailCustomization() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEmailCustomizationRead,
		Schema: buildSchema(
			map[string]*schema.Schema{
				"brand_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Brand ID",
				},
				"template_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Template Name",
				},
				"customization_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The ID of the customization",
				},
			},
			emailCustomizationDataSourceSchema,
		),
	}
}

func dataSourceEmailCustomizationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	brandID, ok := d.GetOk("brand_id")
	if !ok {
		return diag.Errorf("brand_id required for email template")
	}

	templateName, ok := d.GetOk("template_name")
	if !ok {
		return diag.Errorf("template name required for email template")
	}

	customizationId, ok := d.GetOk("customization_id")
	if !ok {
		return diag.Errorf("customization_id required for email customization")
	}

	customization, _, err := getOktaV3ClientFromMetadata(m).CustomizationApi.GetEmailCustomization(ctx, brandID.(string), templateName.(string), customizationId.(string)).Execute()
	if err != nil {
		return diag.Errorf("failed to get email template: %v", err)
	}

	d.SetId(fmt.Sprintf("email_customization-%s-%s-%s", customization.GetId(), templateName.(string), brandID.(string)))
	rawMap := flattenEmailCustomization(customization)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set email customization properties: %v", err)
	}

	return nil
}
