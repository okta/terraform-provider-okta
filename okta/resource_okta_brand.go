package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
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
	brandId := d.Get("brand_id").(string)
	// check that the brand exists, create is short circuited as a reader
	brand, resp, err := getOktaClientFromMetadata(m).Brand.GetBrand(ctx, brandId)
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get brand %q: %v", brandId, err)
	}
	logger(m).Info("setting brand id", "id", brandId)
	d.SetId(brandId)
	rawMap := flattenBrand(brand)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set brand's properties in read: %v", err)
	}
	return nil
}

func resourceBrandRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("reading brand", "id", d.Id())
	brand, resp, err := getOktaClientFromMetadata(m).Brand.GetBrand(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
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

	brand := sdk.Brand{}
	agree, ok := d.GetOk("agree_to_custom_privacy_policy")
	if ok {
		brand.AgreeToCustomPrivacyPolicy = boolPtr(agree.(bool))
	} else {
		brand.AgreeToCustomPrivacyPolicy = boolPtr(false)
	}
	if val, ok := d.GetOk("custom_privacy_policy_url"); ok {
		brand.CustomPrivacyPolicyUrl = val.(string)
	}
	if val, ok := d.GetOk("remove_powered_by_okta"); ok {
		brand.RemovePoweredByOkta = boolPtr(val.(bool))
	}
	_, _, err := getOktaClientFromMetadata(m).Brand.UpdateBrand(ctx, d.Id(), brand)
	if err != nil {
		return diag.Errorf("failed to update brand: %v", err)
	}

	// NOTE: don't do a tail call on resource read, populate the result here
	rawMap := flattenBrand(&brand)
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
	brand, _, err := getOktaClientFromMetadata(m).Brand.GetBrand(ctx, d.Id())
	if err != nil {
		return nil, err
	}

	d.SetId(brand.Id)
	rawMap := flattenBrand(brand)
	err = setNonPrimitives(d, rawMap)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
