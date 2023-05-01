package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v3/okta"
)

func dataSourceEmailTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEmailTemplateRead,
		Schema: buildSchema(
			map[string]*schema.Schema{
				"brand_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Brand ID",
				},
			},
			emailTemplateDataSourceSchema,
		),
	}
}

func dataSourceEmailTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var brand *okta.Brand
	var err error
	brandID, ok := d.GetOk("brand_id")
	if ok {
		logger(m).Info("reading brand by ID", "id", brandID.(string))
		brand, _, err = getOktaV3ClientFromMetadata(m).CustomizationApi.GetBrand(ctx, brandID.(string)).Execute()
		if err != nil {
			return diag.Errorf("failed to get brand for email template: %v", err)
		}
	} else {
		return diag.Errorf("brand_id required for email template: %v", err)
	}

	templateName, ok := d.GetOk("name")
	if !ok {
		return diag.Errorf("name required for email template: %v", err)
	}

	template, _, err := getOktaV3ClientFromMetadata(m).CustomizationApi.GetEmailTemplate(ctx, brandID.(string), templateName.(string)).Execute()
	if err != nil {
		return diag.Errorf("failed to get email template: %v", err)
	}

	d.SetId(fmt.Sprintf("email_template-%s-%s", templateName, brand.GetId()))
	rawMap := flattenEmailTemplate(template)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set email template's properties: %v", err)
	}

	return nil
}
