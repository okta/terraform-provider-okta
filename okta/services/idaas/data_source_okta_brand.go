package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func DataSourceBrand() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBrandRead,
		Schema: utils.BuildSchema(
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

func dataSourceBrandRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var brand *okta.BrandWithEmbedded
	var err error
	brandID := d.Get("brand_id").(string)

	if brandID == "default" {
		brand, err = getDefaultBrand(ctx, meta)
		if err != nil {
			return diag.Errorf("failed to get default brand for org: %v", err)
		}
	} else {
		Logger(meta).Info("reading brand by ID", "id", brandID)
		brand, _, err = GetOktaV3ClientFromMetadata(meta).CustomizationAPI.GetBrand(ctx, brandID).Execute()
		if err != nil {
			return diag.Errorf("failed to get brand: %v", err)
		}
	}

	d.SetId(brand.GetId())
	rawMap := flattenBrand(brand)
	err = utils.SetNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set brand's properties: %v", err)
	}

	return nil
}
