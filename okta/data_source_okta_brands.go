package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBrands() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBrandsRead,
		Schema: map[string]*schema.Schema{
			"brands": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of `okta_brand` belonging to the organization",
				Elem: &schema.Resource{
					Schema: brandsDataSourceSchema,
				},
			},
		},
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
	brandResource := &schema.Resource{
		Schema: brandResourceSchema,
	}
	_ = d.Set("brands", schema.NewSet(schema.HashResource(brandResource), arr))

	return nil
}
