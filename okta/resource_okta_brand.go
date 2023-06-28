package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v3/okta"
)

func resourceBrand() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBrandCreate,
		ReadContext:   resourceBrandRead,
		UpdateContext: resourceBrandUpdate,
		DeleteContext: resourceBrandDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBrandImportStateContext,
		},
		Schema: brandResourceSchema,
	}
}

func resourceBrandCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	brandID := d.Get("brand_id").(string)

	var brand *okta.Brand
	var resp *okta.APIResponse
	var err error

	if brandID == "default" {
		brand, err = getDefaultBrand(ctx, m)
		if err != nil {
			return diag.Errorf("failed to get default brand for org: %v", err)
		}
	} else {
		// check that the brand exists, create is short circuited as a reader
		brand, resp, err = getOktaV3ClientFromMetadata(m).CustomizationApi.GetBrand(ctx, brandID).Execute()
		if err := v3suppressErrorOn404(resp, err); err != nil {
			return diag.Errorf("failed to get brand %q: %v", brandID, err)
		}
	}

	logger(m).Info("setting brand id", "id", brandID)
	d.SetId(brandID)
	rawMap := flattenBrand(brand)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set brand's properties in read: %v", err)
	}
	return nil
}

func resourceBrandRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("reading brand", "id", d.Id())
	brand, resp, err := getOktaV3ClientFromMetadata(m).CustomizationApi.GetBrand(ctx, d.Id()).Execute()
	if err := v3suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get brand: %v", err)
	}
	if brand == nil {
		d.SetId("")
		return nil
	}
	rawMap := flattenBrand(brand)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set brand's properties in read: %v", err)
	}

	return nil
}

func resourceBrandUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("updating brand", "id", d.Id())

	brandRequest := okta.BrandRequest{}
	agree, ok := d.GetOk("agree_to_custom_privacy_policy")
	if ok {
		brandRequest.AgreeToCustomPrivacyPolicy = boolPtr(agree.(bool))
	} else {
		brandRequest.AgreeToCustomPrivacyPolicy = boolPtr(false)
	}
	if val, ok := d.GetOk("custom_privacy_policy_url"); ok {
		brandRequest.CustomPrivacyPolicyUrl = stringPtr(val.(string))
	}
	if val, ok := d.GetOk("remove_powered_by_okta"); ok {
		brandRequest.RemovePoweredByOkta = boolPtr(val.(bool))
	}
	updatedBrand, _, err := getOktaV3ClientFromMetadata(m).CustomizationApi.ReplaceBrand(ctx, d.Id()).Brand(brandRequest).Execute()
	if err != nil {
		return diag.Errorf("failed to update brand: %v", err)
	}

	// NOTE: don't do a tail call on resource read, populate the result here
	rawMap := flattenBrand(updatedBrand)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set brand's properties in update: %v", err)
	}

	return nil
}

func resourceBrandDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// fake delete
	d.SetId("")
	return nil
}

func resourceBrandImportStateContext(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if _, ok := d.GetOk("brand_id"); !ok {
		_ = d.Set("brand_id", d.Id())
	}
	brand, _, err := getOktaV3ClientFromMetadata(m).CustomizationApi.GetBrand(ctx, d.Id()).Execute()
	if err != nil {
		return nil, err
	}

	d.SetId(brand.GetId())
	rawMap := flattenBrand(brand)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
