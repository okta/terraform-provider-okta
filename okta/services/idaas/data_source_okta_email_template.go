package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func DataSourceEmailTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEmailTemplateRead,
		Schema: utils.BuildSchema(
			map[string]*schema.Schema{
				"brand_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Brand ID",
				},
			},
			emailTemplateDataSourceSchema,
		),
		Description: "Get a single Email Template for a Brand belonging to an Okta organization.",
	}
}

func dataSourceEmailTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var brand *okta.BrandWithEmbedded
	var err error
	brandID, ok := d.GetOk("brand_id")
	if ok {
		Logger(meta).Info("reading brand by ID", "id", brandID.(string))
		brand, _, err = GetOktaV3ClientFromMetadata(meta).CustomizationAPI.GetBrand(ctx, brandID.(string)).Execute()
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

	template, _, err := GetOktaV3ClientFromMetadata(meta).CustomizationAPI.GetEmailTemplate(ctx, brandID.(string), templateName.(string)).Execute()
	if err != nil {
		return diag.Errorf("failed to get email template: %v", err)
	}

	d.SetId(fmt.Sprintf("email_template-%s-%s", templateName, brand.GetId()))
	rawMap := flattenEmailTemplate(template)
	err = utils.SetNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set email template's properties: %v", err)
	}

	return nil
}
