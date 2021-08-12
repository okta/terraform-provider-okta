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
				Type:             schema.TypeInt,
				Default:          6,
				Optional:         true,
				Description:      "Factor name",
				ValidateDiagFunc: elemInSlice([]int{6, 8, 10}),
				ForceNew:         true,
			},
			"hmac_algorithm": {
				Type:             schema.TypeString,
				Default:          "HMacSHA512",
				Optional:         true,
				Description:      "Hash-based message authentication code algorithm",
				ValidateDiagFunc: elemInSlice([]string{"HMacSHA1", "HMacSHA256", "HMacSHA512"}),
				ForceNew:         true,
			},
			"time_step": {
				Type:             schema.TypeInt,
				Default:          15,
				Optional:         true,
				Description:      "Time step in seconds",
				ValidateDiagFunc: elemInSlice([]int{15, 30, 60}),
				ForceNew:         true,
			},
			"clock_drift_interval": {
				Type:             schema.TypeInt,
				Default:          3,
				Optional:         true,
				Description:      "Clock drift interval",
				ValidateDiagFunc: elemInSlice([]int{3, 5, 10}),
				ForceNew:         true,
			},
			"shared_secret_encoding": {
				Type:             schema.TypeString,
				Default:          "base32",
				Optional:         true,
				Description:      "Shared secret encoding",
				ValidateDiagFunc: elemInSlice([]string{"base32", "base64", "hexadecimal"}),
				ForceNew:         true,
			},
		},
	}
}

func resourceFactorTOTPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	profile := buildTotpFactorProfile(d)
	responseProfile, _, err := getSupplementFromMetadata(m).CreateHotpFactorProfile(ctx, *profile)
	if err != nil {
		return diag.Errorf("failed to create TOTP factor: %v", err)
	}
	d.SetId(responseProfile.ID)
	return resourceFactorTOTPRead(ctx, d, m)
}

func resourceFactorTOTPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	profile := buildTotpFactorProfile(d)
	_, _, err := getSupplementFromMetadata(m).UpdateHotpFactorProfile(ctx, d.Id(), *profile)
	if err != nil {
		return diag.Errorf("failed to update TOTP factor: %v", err)
	}

	return resourceFactorTOTPRead(ctx, d, m)
}

func resourceFactorTOTPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	profile, resp, err := getSupplementFromMetadata(m).GetHotpFactorProfile(ctx, d.Id())
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
	_, err := getSupplementFromMetadata(m).DeleteHotpFactorProfile(ctx, d.Id())
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
