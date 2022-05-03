package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func dataSourceEmailTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEmailTemplatesRead,
		Schema:      emailTemplatesDataSourceSchema,
	}
}

func dataSourceEmailTemplatesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var brand *okta.Brand
	var err error
	brandID, ok := d.GetOk("brand_id")
	if ok {
		logger(m).Info("reading brand by ID", "id", brandID.(string))
		brand, _, err = getOktaClientFromMetadata(m).Brand.GetBrand(ctx, brandID.(string))
		if err != nil {
			return diag.Errorf("failed to get brand for email templates: %v", err)
		}
	}

	templates, _, err := getOktaClientFromMetadata(m).Brand.ListEmailTemplates(ctx, brand.Id, nil)
	if err != nil {
		return diag.Errorf("failed to list email templates: %v", err)
	}

	d.SetId(fmt.Sprintf("email_templates-%s", brand.Id))
	arr := make([]interface{}, len(templates))
	for i, template := range templates {
		rawMap := flattenEmailTemplate(template)
		arr[i] = rawMap
	}
	emailTemplatesDataSource := &schema.Resource{
		Schema: emailTemplateDataSourceSchema,
	}
	_ = d.Set("email_templates", schema.NewSet(schema.HashResource(emailTemplatesDataSource), arr))

	return nil
}
