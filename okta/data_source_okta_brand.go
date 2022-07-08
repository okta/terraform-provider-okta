package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func dataSourceBrand() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBrandRead,
		Schema: buildSchema(
			map[string]*schema.Schema{
				"brand_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Brand ID",
				},
			},
			brandDataSourceSchema,
		),
	}
}

func dataSourceBrandRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var brand *okta.Brand
	var err error
	brandID, ok := d.GetOk("brand_id")
	if ok {
		logger(m).Info("reading brand by ID", "id", brandID.(string))
		brand, _, err = getOktaClientFromMetadata(m).Brand.GetBrand(ctx, brandID.(string))
		if err != nil {
			return diag.Errorf("failed to get brand: %v", err)
		}
	}
	d.SetId(brand.Id)
	rawMap := flattenBrand(brand)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set brand's properties: %v", err)
	}

	return nil
}
