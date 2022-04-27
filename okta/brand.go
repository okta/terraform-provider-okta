package okta

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

var brandDataSchema = map[string]*schema.Schema{
	"custom_privacy_policy_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Custom privacy policy URL",
	},
	"links": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Link relations for this object - JSON HAL - Discoverable resources related to the app",
	},
	"remove_powered_by_okta": {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: `Removes "Powered by Okta" from the Okta-hosted sign-in page and "© 2021 Okta, Inc." from the Okta End-User Dashboard`,
	},
}

func flattenBrand(brand *okta.Brand) map[string]interface{} {
	attrs := map[string]interface{}{}

	attrs["custom_privacy_policy_url"] = brand.CustomPrivacyPolicyUrl
	links, _ := json.Marshal(brand.Links)
	attrs["links"] = string(links)
	attrs["remove_powered_by_okta"] = *brand.RemovePoweredByOkta

	return attrs
}
