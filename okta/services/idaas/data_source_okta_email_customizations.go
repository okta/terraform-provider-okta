package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func dataSourceEmailCustomizations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEmailCustomizationsRead,
		Schema: utils.BuildSchema(
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
			},
			emailCustomizationsDataSourceSchema,
		),
		Description: "Get the email customizations of an email template belonging to a brand in an Okta organization.",
	}
}

func dataSourceEmailCustomizationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err error
	brandID, ok := d.GetOk("brand_id")
	if !ok {
		return diag.Errorf("brand_id required for email customizations: %v", err)
	}

	templateName, ok := d.GetOk("template_name")
	if !ok {
		return diag.Errorf("template name required for email customizations: %v", err)
	}

	customizations, _, err := getOktaV3ClientFromMetadata(meta).CustomizationAPI.ListEmailCustomizations(ctx, brandID.(string), templateName.(string)).Execute()
	if err != nil {
		return diag.Errorf("failed to list email customizations: %v", err)
	}

	d.SetId(fmt.Sprintf("email_customizations-%s-%s", templateName, brandID.(string)))
	arr := make([]interface{}, len(customizations))
	for i, customization := range customizations {
		rawMap := flattenEmailCustomization(&customization)
		arr[i] = rawMap
	}

	err = d.Set("email_customizations", schema.NewSet(hashEmailCustomization, arr))
	if err != nil {
		return diag.Errorf("failed to set email customizations: %v", err)
	}

	return nil
}
