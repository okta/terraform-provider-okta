package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v3/okta"
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
		Description: "Get a single Brand from Okta.",
	}
}

func dataSourceBrandRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var brand *okta.BrandWithEmbedded
	var err error
	brandID := d.Get("brand_id").(string)

	if brandID == "default" {
		brand, err = getDefaultBrand(ctx, m)
		if err != nil {
			return diag.Errorf("failed to get default brand for org: %v", err)
		}
	} else {
		logger(m).Info("reading brand by ID", "id", brandID)
		brand, _, err = getOktaV3ClientFromMetadata(m).CustomizationApi.GetBrand(ctx, brandID).Execute()
		if err != nil {
			return diag.Errorf("failed to get brand: %v", err)
		}
	}

	d.SetId(brand.GetId())
	rawMap := flattenBrand(brand)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set brand's properties: %v", err)
	}

	return nil
}
