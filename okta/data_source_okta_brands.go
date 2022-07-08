package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBrands() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBrandsRead,
		Schema:      brandsDataSourceSchema,
	}
}

func dataSourceBrandsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	brands, _, err := getOktaClientFromMetadata(m).Brand.ListBrands(ctx)
	if err != nil {
		return diag.Errorf("failed to list brands: %v", err)
	}

	d.SetId("brands")
	arr := make([]interface{}, len(brands))
	for i, brand := range brands {
		rawMap := flattenBrand(brand)
		rawMap["id"] = brand.Id
		arr[i] = rawMap
	}
	brandDataSource := &schema.Resource{
		Schema: brandDataSourceSchema,
	}
	_ = d.Set("brands", schema.NewSet(schema.HashResource(brandDataSource), arr))

	return nil
}
