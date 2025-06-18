package idaas

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
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
		Description: `Allows you to manage the time-based one-time password (TOTP) factors. A time-based one-time password (TOTP) is a
		temporary passcode that is generated for user authentication. Examples of TOTP include hardware authenticators and
		mobile app authenticators.
		
Once saved, the settings cannot be changed (except for the 'name' field). Any other change would force resource
recreation.`,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The TOTP name.",
			},
			"otp_length": {
				Type:        schema.TypeInt,
				Default:     6,
				Optional:    true,
				Description: "Length of the password. Default is `6`.",
				ForceNew:    true,
			},
			"hmac_algorithm": {
				Type:        schema.TypeString,
				Default:     "HMacSHA512",
				Optional:    true,
				Description: "HMAC Algorithm. Valid values: `HMacSHA1`, `HMacSHA256`, `HMacSHA512`. Default is `HMacSHA512`.",
				ForceNew:    true,
			},
			"time_step": {
				Type:        schema.TypeInt,
				Default:     15,
				Optional:    true,
				Description: "Time step in seconds. Valid values: `15`, `30`, `60`. Default is `15`.",
				ForceNew:    true,
			},
			"clock_drift_interval": {
				Type:        schema.TypeInt,
				Default:     3,
				Optional:    true,
				Description: "Clock drift interval. This setting allows you to build in tolerance for any drift between the token's current time and the server's current time. Valid values: `3`, `5`, `10`. Default is `3`.",
				ForceNew:    true,
			},
			"shared_secret_encoding": {
				Type:        schema.TypeString,
				Default:     "base32",
				Optional:    true,
				Description: "Shared secret encoding. Valid values: `base32`, `base64`, `hexadecimal`. Default is `base32`.",
				ForceNew:    true,
			},
		},
	}
}

func resourceFactorTOTPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	profile := buildTotpFactorProfile(d)
	responseProfile, _, err := getAPISupplementFromMetadata(meta).CreateHotpFactorProfile(ctx, *profile)
	if err != nil {
		return diag.Errorf("failed to create TOTP factor: %v", err)
	}
	d.SetId(responseProfile.ID)
	return resourceFactorTOTPRead(ctx, d, meta)
}

func resourceFactorTOTPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	profile := buildTotpFactorProfile(d)
	_, _, err := getAPISupplementFromMetadata(meta).UpdateHotpFactorProfile(ctx, d.Id(), *profile)
	if err != nil {
		return diag.Errorf("failed to update TOTP factor: %v", err)
	}

	return resourceFactorTOTPRead(ctx, d, meta)
}

func resourceFactorTOTPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	profile, resp, err := getAPISupplementFromMetadata(meta).GetHotpFactorProfile(ctx, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
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

func resourceFactorTOTPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// NOTE: The publicly documented DELETE /api/v1/org/factors/hotp/profiles/{id} appears to only 501 at the present time.

	resp, err := getAPISupplementFromMetadata(meta).DeleteHotpFactorProfile(ctx, d.Id())
	if err != nil {
		if resp.StatusCode == http.StatusNotImplemented {
			logger(meta).Warn("Okta API declares deletion of totp factors as not implemented")
			return nil
		}
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
