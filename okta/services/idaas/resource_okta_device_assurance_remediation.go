package idaas

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-sdk-golang/v7/okta"
)

const (
	deviceAssuranceRemediationShow = "SHOW"
	deviceAssuranceRemediationHide = "HIDE"
	deviceAssuranceGraceByDateTime = "BY_DATE_TIME"
	deviceAssuranceGraceByDuration = "BY_DURATION"
)

type deviceAssuranceGracePeriodModel struct {
	Type   types.String `tfsdk:"type"`
	Expiry types.String `tfsdk:"expiry"`
}

func deviceAssuranceDisplayRemediationModeAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Description: "Whether the Sign-In Widget shows remediation instructions to users who fail the policy, `SHOW` or `HIDE`. Unset means Okta's default, `HIDE`. A `grace_period` requires `SHOW`.",
		Optional:    true,
		Validators: []validator.String{
			stringvalidator.OneOf(deviceAssuranceRemediationHide, deviceAssuranceRemediationShow),
		},
	}
}

func deviceAssuranceGracePeriodBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "Time users get to remediate before access is denied. Okta only accepts it together with `display_remediation_mode = \"SHOW\"`. Omit for no grace period.",
		Validators:  []validator.Object{deviceAssuranceGracePeriodValidator{}},
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "`BY_DURATION` counts from the first failed check, `BY_DATE_TIME` ends at a fixed instant.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(deviceAssuranceGraceByDateTime, deviceAssuranceGraceByDuration),
				},
			},
			"expiry": schema.StringAttribute{
				Description: "For `BY_DURATION`, an ISO 8601 duration in hours as Okta stores it, for example `PT168H` for 7 days. For `BY_DATE_TIME`, an RFC 3339 timestamp without fractional seconds as Okta stores it, for example `2026-12-01T00:00:00Z`.",
				Optional:    true,
			},
		},
	}
}

type deviceAssuranceGracePeriodValidator struct{}

func (deviceAssuranceGracePeriodValidator) Description(context.Context) string {
	return "grace_period needs type and expiry, and display_remediation_mode = \"SHOW\""
}

func (v deviceAssuranceGracePeriodValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (deviceAssuranceGracePeriodValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	for _, name := range []string{"type", "expiry"} {
		if value, ok := req.ConfigValue.Attributes()[name]; !ok || value.IsNull() {
			resp.Diagnostics.AddAttributeError(req.Path.AtName(name), "Missing grace period attribute", "grace_period needs both type and expiry.")
		}
	}
	var mode types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("display_remediation_mode"), &mode)...)
	if resp.Diagnostics.HasError() || mode.IsUnknown() {
		return
	}
	if mode.IsNull() || mode.ValueString() != deviceAssuranceRemediationShow {
		resp.Diagnostics.AddAttributeError(req.Path, "Grace period requires shown remediation", "Okta rejects a grace period unless display_remediation_mode is \"SHOW\". Set display_remediation_mode = \"SHOW\" or remove grace_period.")
	}
}

func applyDeviceAssuranceRemediation(policy *okta.DeviceAssurance, mode types.String, gracePeriod *deviceAssuranceGracePeriodModel) {
	if !mode.IsNull() && !mode.IsUnknown() {
		policy.SetDisplayRemediationMode(mode.ValueString())
	}
	if gracePeriod == nil {
		return
	}
	period := okta.NewGracePeriod()
	period.SetType(gracePeriod.Type.ValueString())
	period.SetExpiry(okta.StringAsGracePeriodExpiry(gracePeriod.Expiry.ValueStringPointer()))
	policy.SetGracePeriod(*period)
}

func deviceAssuranceRemediationToState(policy *okta.DeviceAssurance, prior types.String) (types.String, *deviceAssuranceGracePeriodModel) {
	mode := types.StringPointerValue(policy.DisplayRemediationMode)
	if prior.IsNull() && mode.ValueString() == deviceAssuranceRemediationHide {
		mode = types.StringNull()
	}
	period, ok := policy.GetGracePeriodOk()
	if !ok || period == nil {
		return mode, nil
	}
	return mode, &deviceAssuranceGracePeriodModel{
		Type:   types.StringPointerValue(period.Type),
		Expiry: types.StringValue(deviceAssuranceExpiryString(period.Expiry)),
	}
}

func deviceAssuranceExpiryString(expiry *okta.GracePeriodExpiry) string {
	if expiry == nil {
		return ""
	}
	if expiry.String != nil {
		return *expiry.String
	}
	if expiry.TimeTime != nil {
		return expiry.TimeTime.UTC().Format(time.RFC3339)
	}
	return ""
}

// Works around okta/okta-sdk-golang#608.
func decodeDeviceAssuranceResponse(policy *okta.ListDeviceAssurancePolicies200ResponseInner, err error) (*okta.ListDeviceAssurancePolicies200ResponseInner, error) {
	if err == nil {
		return policy, nil
	}
	var apiErr *okta.GenericOpenAPIError
	if !errors.As(err, &apiErr) || !strings.Contains(err.Error(), "oneOf(GracePeriodExpiry)") {
		return nil, err
	}
	var fields map[string]json.RawMessage
	if decodeErr := json.Unmarshal(apiErr.Body(), &fields); decodeErr != nil {
		return nil, err
	}
	var rawGracePeriod struct {
		Type   string `json:"type"`
		Expiry string `json:"expiry"`
	}
	if decodeErr := json.Unmarshal(fields["gracePeriod"], &rawGracePeriod); decodeErr != nil {
		return nil, err
	}
	delete(fields, "gracePeriod")
	withoutGracePeriod, marshalErr := json.Marshal(fields)
	if marshalErr != nil {
		return nil, err
	}
	var decoded okta.ListDeviceAssurancePolicies200ResponseInner
	if decodeErr := json.Unmarshal(withoutGracePeriod, &decoded); decodeErr != nil {
		return nil, err
	}
	base := deviceAssuranceBase(&decoded)
	if base == nil {
		return nil, err
	}
	period := okta.NewGracePeriod()
	period.SetType(rawGracePeriod.Type)
	period.SetExpiry(okta.StringAsGracePeriodExpiry(&rawGracePeriod.Expiry))
	base.SetGracePeriod(*period)
	return &decoded, nil
}

func deviceAssuranceBase(policy *okta.ListDeviceAssurancePolicies200ResponseInner) *okta.DeviceAssurance {
	switch {
	case policy.DeviceAssuranceAndroidPlatform != nil:
		return &policy.DeviceAssuranceAndroidPlatform.DeviceAssurance
	case policy.DeviceAssuranceChromeOSPlatform != nil:
		return &policy.DeviceAssuranceChromeOSPlatform.DeviceAssurance
	case policy.DeviceAssuranceIOSPlatform != nil:
		return &policy.DeviceAssuranceIOSPlatform.DeviceAssurance
	case policy.DeviceAssuranceMacOSPlatform != nil:
		return &policy.DeviceAssuranceMacOSPlatform.DeviceAssurance
	case policy.DeviceAssuranceWindowsPlatform != nil:
		return &policy.DeviceAssuranceWindowsPlatform.DeviceAssurance
	}
	return nil
}
