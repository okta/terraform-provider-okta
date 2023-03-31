package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEmailCustomizations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEmailCustomizationsRead,
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
			},
			emailCustomizationsDataSourceSchema,
		),
	}
}

func dataSourceEmailCustomizationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error
	brandID, ok := d.GetOk("brand_id")
	if !ok {
		return diag.Errorf("brand_id required for email customizations: %v", err)
	}

	templateName, ok := d.GetOk("template_name")
	if !ok {
		return diag.Errorf("template name required for email customizations: %v", err)
	}

	customizations, _, err := getOktaV3ClientFromMetadata(m).CustomizationApi.ListEmailCustomizations(ctx, brandID.(string), templateName.(string)).Execute()
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
