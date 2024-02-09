package okta

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v4/okta"
)

var brandDataSourceSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The ID of the Brand",
	},
	"name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Brand name",
	},
	"custom_privacy_policy_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Custom privacy policy URL",
	},
	"links": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Link relations for this object - JSON HAL - Discoverable resources related to the brand",
	},
	"remove_powered_by_okta": {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: `Removes "Powered by Okta" from the Okta-hosted sign-in page and "© 2021 Okta, Inc." from the Okta End-User Dashboard`,
	},
}

var brandsDataSourceSchema = map[string]*schema.Schema{
	"brands": {
		Type:        schema.TypeSet,
		Computed:    true,
		Description: "List of `okta_brand` belonging to the organization",
		Elem: &schema.Resource{
			Schema: brandDataSourceSchema,
		},
	},
}

func flattenBrand(brand *okta.BrandWithEmbedded) map[string]interface{} {
	attrs := map[string]interface{}{}
	attrs["id"] = brand.GetId()
	attrs["custom_privacy_policy_url"] = ""
	if brand.GetCustomPrivacyPolicyUrl() != "" {
		attrs["custom_privacy_policy_url"] = brand.GetCustomPrivacyPolicyUrl()
	}
	attrs["name"] = brand.GetName()
	links, _ := json.Marshal(brand.GetLinks())
	attrs["links"] = string(links)
	attrs["remove_powered_by_okta"] = false
	if brand.RemovePoweredByOkta != nil {
		attrs["remove_powered_by_okta"] = brand.GetRemovePoweredByOkta()
	}

	return attrs
}

func getDefaultBrand(ctx context.Context, m interface{}) (*okta.BrandWithEmbedded, error) {
	brands, _, err := getOktaV3ClientFromMetadata(m).CustomizationAPI.ListBrands(ctx).Execute()
	if err != nil {
		return nil, err
	}

	for _, brand := range brands {
		if brand.GetIsDefault() {
			return &brand, nil
		}
	}

	return nil, errors.New("failed to get default brand")
}
