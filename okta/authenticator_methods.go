package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func createPhoneMethods(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	_, ok := d.GetOk("phone_methods")
	if ok {
		return configurePhoneMethods(ctx, d, m)
	}
	return nil
}

func updatePhoneMethods(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	old, new := d.GetChange("phone_methods")
	if !old.(*schema.Set).Equal(new) {
		// if len is 0 assume not using provider to control methods / enable both the default (UI does not allow both disabled)
		if new.(*schema.Set).Len() == 0 {
			for _, v := range []string{"sms", "voice"} {
				if err := getSupplementFromMetadata(m).ActivateAuthenticatorMethod(ctx, d.Id(), v); err != nil {
					return err
				}
			}
		} else {
			return configurePhoneMethods(ctx, d, m)
		}
	}
	return nil
}

func configurePhoneMethods(ctx context.Context, d *schema.ResourceData, m interface{}) error { // diag.Diagnostics {
	if methods := d.Get("phone_methods"); methods != nil {
		for _, v := range []string{"sms", "voice"} {
			if methods.(*schema.Set).Contains(v) {
				if err := getSupplementFromMetadata(m).ActivateAuthenticatorMethod(ctx, d.Id(), v); err != nil {
					return err
				}
			} else {
				if err := getSupplementFromMetadata(m).DeactivateAuthenticatorMethod(ctx, d.Id(), v); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func readPhoneMethods(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	methods := d.Get("phone_methods")
	if methods != nil {
		for _, v := range []string{"sms", "voice"} {
			status, err := getSupplementFromMetadata(m).GetAuthenticatorMethodStatus(ctx, d.Id(), v)
			if err != nil {
				return err
			}
			if status == "ACTIVE" {
				methods.(*schema.Set).Add(v)
			} else {
				methods.(*schema.Set).Remove(v)
			}
		}
		_ = d.Set("phone_methods", methods)
	}
	return nil
}
