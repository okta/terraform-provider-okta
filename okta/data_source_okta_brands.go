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
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: buildSchema(
						brandDataSchema,
						map[string]*schema.Schema{
							"id": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "Brand ID",
							},
						}),
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

	d.SetId("brands") // there is only one brands list on the org
	arr := make([]map[string]interface{}, len(brands))
	for i, brand := range brands {
		rawMap := flattenBrand(brand)
		rawMap["id"] = brand.Id
		arr[i] = rawMap
	}
	_ = d.Set("brands", arr)

	return nil
}
