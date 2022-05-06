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
		Schema:      emailCustomizationsDataSourceSchema,
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

	customizations, _, err := getOktaClientFromMetadata(m).Brand.ListEmailTemplateCustomizations(ctx, brandID.(string), templateName.(string))
	if err != nil {
		return diag.Errorf("failed to list email customizations: %v", err)
	}

	d.SetId(fmt.Sprintf("email_customizations-%s-%s", templateName, brandID.(string)))
	arr := make([]interface{}, len(customizations))
	for i, customization := range customizations {
		rawMap := flattenEmailCustomization(brandID.(string), templateName.(string), true, customization)
		arr[i] = rawMap
	}

	err = d.Set("email_customizations", schema.NewSet(hashEmailCustomization, arr))
	if err != nil {
		return diag.Errorf("failed to set email customizations: %v", err)
	}

	return nil
}
