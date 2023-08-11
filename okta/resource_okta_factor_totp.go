package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceFactorTOTP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFactorTOTPCreate,
		ReadContext:   resourceFactorTOTPRead,
		UpdateContext: resourceFactorTOTPUpdate,
		DeleteContext: resourceFactorTOTPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Factor name",
			},
			"otp_length": {
				Type:        schema.TypeInt,
				Default:     6,
				Optional:    true,
				Description: "Factor name",
				ForceNew:    true,
			},
			"hmac_algorithm": {
				Type:        schema.TypeString,
				Default:     "HMacSHA512",
				Optional:    true,
				Description: "Hash-based message authentication code algorithm",
				ForceNew:    true,
			},
			"time_step": {
				Type:        schema.TypeInt,
				Default:     15,
				Optional:    true,
				Description: "Time step in seconds",
				ForceNew:    true,
			},
			"clock_drift_interval": {
				Type:        schema.TypeInt,
				Default:     3,
				Optional:    true,
				Description: "Clock drift interval",
				ForceNew:    true,
			},
			"shared_secret_encoding": {
				Type:        schema.TypeString,
				Default:     "base32",
				Optional:    true,
				Description: "Shared secret encoding",
				ForceNew:    true,
			},
		},
	}
}

func resourceFactorTOTPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	profile := buildTotpFactorProfile(d)
	responseProfile, _, err := getAPISupplementFromMetadata(m).CreateHotpFactorProfile(ctx, *profile)
	if err != nil {
		return diag.Errorf("failed to create TOTP factor: %v", err)
	}
	d.SetId(responseProfile.ID)
	return resourceFactorTOTPRead(ctx, d, m)
}

func resourceFactorTOTPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	profile := buildTotpFactorProfile(d)
	_, _, err := getAPISupplementFromMetadata(m).UpdateHotpFactorProfile(ctx, d.Id(), *profile)
	if err != nil {
		return diag.Errorf("failed to update TOTP factor: %v", err)
	}

	return resourceFactorTOTPRead(ctx, d, m)
}

func resourceFactorTOTPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	profile, resp, err := getAPISupplementFromMetadata(m).GetHotpFactorProfile(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get TOTP factor: %v", err)
	}
	if profile == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", profile.Name)
	_ = d.Set("otp_length", profile.Settings.OtpLength)
	_ = d.Set("time_step", profile.Settings.TimeStep)
	_ = d.Set("clock_drift_interval", profile.Settings.AcceptableAdjacentIntervals)
	_ = d.Set("shared_secret_encoding", profile.Settings.Encoding)
	_ = d.Set("hmac_algorithm", profile.Settings.HmacAlgorithm)
	return nil
}

func resourceFactorTOTPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// NOTE: The publicly documented DELETE /api/v1/org/factors/hotp/profiles/{id} appears to only 501 at the present time.

	_, err := getAPISupplementFromMetadata(m).DeleteHotpFactorProfile(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete TOTP factor: %v", err)
	}
	return nil
}

func buildTotpFactorProfile(d *schema.ResourceData) *sdk.HotpFactorProfile {
	return &sdk.HotpFactorProfile{
		Name: d.Get("name").(string),
		Settings: sdk.HotpFactorProfileSettings{
			TimeBased:                   true,
			OtpLength:                   d.Get("otp_length").(int),
			TimeStep:                    d.Get("time_step").(int),
			AcceptableAdjacentIntervals: d.Get("clock_drift_interval").(int),
			Encoding:                    d.Get("shared_secret_encoding").(string),
			HmacAlgorithm:               d.Get("hmac_algorithm").(string),
		},
	}
}
