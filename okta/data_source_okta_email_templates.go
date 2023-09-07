package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v3/okta"
)

func dataSourceEmailTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEmailTemplatesRead,
		Schema: buildSchema(
			map[string]*schema.Schema{
				"brand_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Brand ID",
				},
			},
			emailTemplatesDataSourceSchema,
		),
	}
}

func dataSourceEmailTemplatesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var brand *okta.BrandWithEmbedded
	var err error
	brandID, ok := d.GetOk("brand_id")
	if ok {
		logger(m).Info("reading brand by ID", "id", brandID.(string))
		brand, _, err = getOktaV3ClientFromMetadata(m).CustomizationApi.GetBrand(ctx, brandID.(string)).Execute()
		if err != nil {
			return diag.Errorf("failed to get brand for email templates: %v", err)
		}
	} else {
		return diag.Errorf("brand_id required for email templates: %v", err)
	}

	templates, err := collectEmailTempates(ctx, getOktaV3ClientFromMetadata(m), brand)
	if err != nil {
		return diag.Errorf("failed to list email templates: %v", err)
	}

	d.SetId(fmt.Sprintf("email_templates-%s", brand.GetId()))
	arr := make([]interface{}, len(templates))
	for i, template := range templates {
		rawMap := flattenEmailTemplate(&template)
		arr[i] = rawMap
	}
	emailTemplatesDataSource := &schema.Resource{
		Schema: emailTemplateDataSourceSchema,
	}
	_ = d.Set("email_templates", schema.NewSet(schema.HashResource(emailTemplatesDataSource), arr))

	return nil
}

func collectEmailTempates(ctx context.Context, client *okta.APIClient, brand *okta.BrandWithEmbedded) ([]okta.EmailTemplate, error) {
	templates, resp, err := client.CustomizationApi.ListEmailTemplates(ctx, brand.GetId()).Limit(int32(defaultPaginationLimit)).Execute()
	if err != nil {
		return nil, err
	}
	for resp.HasNextPage() {
		var nextTemplates []okta.EmailTemplate
		resp, err = resp.Next(&nextTemplates)
		if err != nil {
			return nil, err
		}
		templates = append(templates, nextTemplates...)
	}
	return templates, nil
}
