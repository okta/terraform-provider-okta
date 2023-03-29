package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTheme() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceThemeRead,
		Schema: buildSchema(
			map[string]*schema.Schema{
				"brand_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Brand ID",
				},
				"theme_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Theme ID",
				},
			},
			themeDataSourceSchema,
		),
	}
}

func dataSourceThemeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	bid, ok := d.GetOk("brand_id")
	if !ok {
		return diag.Errorf("brand_id required for email template")
	}
	brandID := bid.(string)

	tid, ok := d.GetOk("theme_id")
	if !ok {
		return diag.Errorf("theme_id required for theme")
	}
	themeID := tid.(string)

	theme, _, err := getOktaV3ClientFromMetadata(m).CustomizationApi.GetBrandTheme(ctx, brandID, themeID).Execute()
	if err != nil {
		return diag.Errorf("failed to get email template: %v", err)
	}

	d.SetId(theme.GetId())
	rawMap := flattenTheme(brandID, theme)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set theme properties: %v", err)
	}

	return nil
}
