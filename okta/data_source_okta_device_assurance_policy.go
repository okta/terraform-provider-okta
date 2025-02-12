package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-sdk-golang/v5/okta"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &deviceAssurancePolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceAssurancePolicyDataSource{}
)

func NewDeviceAssurancePolicyDataSource() datasource.DataSource {
	return &deviceAssurancePolicyDataSource{}
}

type deviceAssurancePolicyDataSource struct {
	config *Config
}

func (d *deviceAssurancePolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_assurance_policy"
}

type deviceAssurancePolicyModel struct {
	ID                       types.String          `tfsdk:"id"`
	Name                     types.String          `tfsdk:"name"`
	Platform                 types.String          `tfsdk:"platform"`
	DiskEncryptionType       types.Object          `tfsdk:"disk_encryption_type"`
	Jaibreak                 types.Bool            `tfsdk:"jailbreak"`
	OSVersion                types.Object          `tfsdk:"os_version"`
	ScreenlockType           types.Object          `tfsdk:"screenlock_type"`
	SecureHardwarePresent    types.Bool            `tfsdk:"secure_hardware_present"`
	ThirdPartySignalProvider types.Object          `tfsdk:"third_party_signal_provider"`
	OSVersionConstraint      []OsVersionConstraint `tfsdk:"os_version_constraint"`
}

type OsVersionConstraint struct {
	MajorVersionConstraint    types.String `tfsdk:"major_version_constraint"`
	DynamicVersionRequirement types.Object `tfsdk:"dynamic_version_requirement"`
}

type OktaDeviceAssurancePolicy interface {
	GetId() string
	GetName() string
	GetPlatform() string
}

func (d *deviceAssurancePolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get a policy assurance from Okta.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the user type to retrieve, conflicts with `name`.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("name"),
					}...),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of user type to retrieve, conflicts with `id`.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("id"),
					}...),
				},
			},
			"platform": schema.StringAttribute{
				Description: "Policy device assurance platform",
				Computed:    true,
			},
			"disk_encryption_type": schema.ObjectAttribute{
				Description: "List of disk encryption type, can be `FULL`, `USER`",
				Computed:    true,
				AttributeTypes: map[string]attr.Type{
					"include": types.SetType{
						ElemType: types.StringType,
					},
				},
			},
			"jailbreak": schema.BoolAttribute{
				Description: "Is the device jailbroken in the device assurance policy.",
				Computed:    true,
			},
			"os_version": schema.ObjectAttribute{
				Description: "Minimum os version of the device in the device assurance policy.",
				Computed:    true,
				AttributeTypes: map[string]attr.Type{
					"minimum": types.StringType,
					"dynamic_version_requirement": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"type":                       types.StringType,
							"distance_from_latest_major": types.Int64Type,
							"latest_security_patch":      types.BoolType,
						},
					},
				},
			},
			"screenlock_type": schema.ObjectAttribute{
				Description: "List of screenlock type, can be `BIOMETRIC` or `BIOMETRIC, PASSCODE`",
				Computed:    true,
				AttributeTypes: map[string]attr.Type{
					"include": types.SetType{
						ElemType: types.StringType,
					},
				},
			},
			"secure_hardware_present": schema.BoolAttribute{
				Description: "Indicates if the device contains a secure hardware functionality",
				Optional:    true,
			},
			"third_party_signal_provider": schema.ObjectAttribute{
				Description: "Indicates if the device contains a secure hardware functionality",
				Optional:    true,
				AttributeTypes: map[string]attr.Type{
					"dtc": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"allow_screen_lock": types.BoolType,
							"browser_version": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"minimum": types.StringType,
								},
							},
							"built_in_dns_client_enabled":       types.BoolType,
							"chrome_remote_desktop_app_blocked": types.BoolType,
							"device_enrollment_domain":          types.StringType,
							"disk_encrypted":                    types.BoolType,
							"key_trust_level":                   types.StringType,
							"managed_device":                    types.BoolType,
							"os_firewall":                       types.BoolType,
							"os_version": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"minimum": types.StringType,
								},
							},
							"password_protection_warning_trigger": types.StringType,
							"realtime_url_check_mode":             types.BoolType,
							"safe_browsing_protection_level":      types.StringType,
							"screen_lock_secured":                 types.BoolType,
							"site_isolation_enabled":              types.BoolType,
							"crowd_strike_agent_id":               types.StringType,
							"crowd_strike_customer_id":            types.StringType,
							"third_party_blocking_enabled":        types.BoolType,
							"windows_machine_domain":              types.StringType,
							"windows_user_domain":                 types.StringType,
						},
					},
				},
			},
			"os_version_constraint": schema.ListAttribute{
				Description: "The list of os version constraints.",
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"major_version_constraint": types.StringType,
						"dynamic_version_requirement": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"type":                       types.StringType,
								"distance_from_latest_major": types.Int64Type,
								"latest_security_patch":      types.BoolType,
							},
						},
					},
				},
			},
		},
	}
}

func (d *deviceAssurancePolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.config = dataSourceConfiguration(req, resp)
}

func (d *deviceAssurancePolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data deviceAssurancePolicyModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var devicePolicyAssuranceResp *okta.ListDeviceAssurancePolicies200ResponseInner
	if data.ID.ValueString() != "" {
		devicePolicyAssuranceResp, _, err = d.config.oktaSDKClientV5.DeviceAssuranceAPI.GetDeviceAssurancePolicy(ctx, data.ID.ValueString()).Execute()
	} else {
		devicePolicyAssuranceResp, err = findDeviceAssurancePolicyByName(ctx, d.config.oktaSDKClientV5, data.Name.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to get device assurance policy",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(mapDeviceAssurancePolicyDatasource(devicePolicyAssuranceResp, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapDeviceAssurancePolicyDatasource(data *okta.ListDeviceAssurancePolicies200ResponseInner, state *deviceAssurancePolicyModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if data == nil {
		diags.AddError("Empty response", "device assurance object")
		return diags
	}

	dataInstance := data.GetActualInstance()
	sharedData := dataInstance.(OktaDeviceAssurancePolicy)
	switch v := dataInstance.(type) {
	case *okta.DeviceAssuranceAndroidPlatform:
		if det, ok := v.GetDiskEncryptionTypeOk(); ok {
			include := make([]attr.Value, 0)
			for _, i := range det.Include {
				include = append(include, types.StringValue(i))
			}
			includeValue, diag := types.SetValue(types.StringType, include)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			detValueMap := map[string]attr.Value{
				"include": includeValue,
			}
			detTypesMap := map[string]attr.Type{
				"include": types.SetType{
					ElemType: types.StringType,
				},
			}
			diskEncryptionType, diag := types.ObjectValue(detTypesMap, detValueMap)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			state.DiskEncryptionType = diskEncryptionType
		}
		if osv, ok := v.GetOsVersionOk(); ok {
			dvrTypesMap := map[string]attr.Type{
				"type":                       types.StringType,
				"distance_from_latest_major": types.Int64Type,
				"latest_security_patch":      types.BoolType,
			}
			dynamicVersionRequirement := types.ObjectValueMust(dvrTypesMap, map[string]attr.Value{
				"type":                       types.StringPointerValue(nil),
				"distance_from_latest_major": types.Int64PointerValue(nil),
				"latest_security_patch":      types.BoolPointerValue(nil),
			})
			osvValueMap := map[string]attr.Value{
				"minimum":                     types.StringPointerValue(osv.Minimum),
				"dynamic_version_requirement": dynamicVersionRequirement,
			}
			osvTypesMap := map[string]attr.Type{
				"minimum": types.StringType,
				"dynamic_version_requirement": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type":                       types.StringType,
						"distance_from_latest_major": types.Int64Type,
						"latest_security_patch":      types.BoolType,
					},
				},
			}
			if dvr, ok := osv.GetDynamicVersionRequirementOk(); ok {
				dvrValueMap := map[string]attr.Value{
					"type":                       types.StringPointerValue(dvr.Type),
					"distance_from_latest_major": types.Int64Value(int64(dvr.GetDistanceFromLatestMajor())),
					"latest_security_patch":      types.BoolPointerValue(dvr.LatestSecurityPatch),
				}
				dynamicVersionRequirement, diag := types.ObjectValue(dvrTypesMap, dvrValueMap)
				if diag != nil {
					diags.Append(diag...)
					return diags
				}
				osvValueMap["dynamic_version_requirement"] = dynamicVersionRequirement
			}
			osVersion, diag := types.ObjectValue(osvTypesMap, osvValueMap)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			state.OSVersion = osVersion
		}
		if slt, ok := v.GetScreenLockTypeOk(); ok {
			include := make([]attr.Value, 0)
			for _, i := range slt.Include {
				include = append(include, types.StringValue(i))
			}
			includeValue, diag := types.SetValue(types.StringType, include)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			sltValueMap := map[string]attr.Value{
				"include": includeValue,
			}
			sltTypesMap := map[string]attr.Type{
				"include": types.SetType{
					ElemType: types.StringType,
				},
			}
			screenlockType, diag := types.ObjectValue(sltTypesMap, sltValueMap)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			state.ScreenlockType = screenlockType
		}
		state.Jaibreak = types.BoolPointerValue(v.Jailbreak)
		state.SecureHardwarePresent = types.BoolPointerValue(v.SecureHardwarePresent)
	case *okta.DeviceAssuranceChromeOSPlatform:
		if tpsp, ok := v.GetThirdPartySignalProvidersOk(); ok {
			tpspTypesMap, dtcTypesMap, bvTypesMap, osvTypesMap, tpspValueMap, dtcValueMap, bvValueMap, osvValueMap := initEmptyTpsp()
			if dtc, ok := tpsp.GetDtcOk(); ok {
				dtcValueMap["allow_screen_lock"] = types.BoolPointerValue(dtc.AllowScreenLock)
				dtcValueMap["built_in_dns_client_enabled"] = types.BoolPointerValue(dtc.BuiltInDnsClientEnabled)
				dtcValueMap["chrome_remote_desktop_app_blocked"] = types.BoolPointerValue(dtc.ChromeRemoteDesktopAppBlocked)
				dtcValueMap["device_enrollment_domain"] = types.StringPointerValue(dtc.DeviceEnrollmentDomain)
				dtcValueMap["disk_encrypted"] = types.BoolPointerValue(dtc.DiskEncrypted)
				dtcValueMap["key_trust_level"] = types.StringPointerValue(dtc.KeyTrustLevel)
				dtcValueMap["managed_device"] = types.BoolPointerValue(dtc.ManagedDevice)
				dtcValueMap["os_firewall"] = types.BoolPointerValue(dtc.OsFirewall)
				dtcValueMap["password_protection_warning_trigger"] = types.StringPointerValue(dtc.PasswordProtectionWarningTrigger)
				dtcValueMap["realtime_url_check_mode"] = types.BoolPointerValue(dtc.RealtimeUrlCheckMode)
				dtcValueMap["safe_browsing_protection_level"] = types.StringPointerValue(dtc.SafeBrowsingProtectionLevel)
				dtcValueMap["screen_lock_secured"] = types.BoolPointerValue(dtc.ScreenLockSecured)
				dtcValueMap["site_isolation_enabled"] = types.BoolPointerValue(dtc.SiteIsolationEnabled)
				if bv, ok := dtc.GetBrowserVersionOk(); ok {
					bvValueMap["minimum"] = types.StringPointerValue(bv.Minimum)
					browserVersion, diag := types.ObjectValue(bvTypesMap, bvValueMap)
					if diag != nil {
						diags.Append(diag...)
						return diags
					}
					dtcValueMap["browser_version"] = browserVersion
				}
				if osv, ok := dtc.GetOsVersionOk(); ok {
					osvValueMap["minimum"] = types.StringPointerValue(osv.Minimum)
					osVersion, diag := types.ObjectValue(osvTypesMap, osvValueMap)
					if diag != nil {
						diags.Append(diag...)
						return diags
					}
					dtcValueMap["os_version"] = osVersion
				}
				dtc, diag := types.ObjectValue(dtcTypesMap, dtcValueMap)
				if diag != nil {
					diags.Append(diag...)
					return diags
				}
				tpspValueMap["dtc"] = dtc
			}
			thirdPartySignalProvider, diag := types.ObjectValue(tpspTypesMap, tpspValueMap)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			state.ThirdPartySignalProvider = thirdPartySignalProvider
		}
	case *okta.DeviceAssuranceIOSPlatform:
		if osv, ok := v.GetOsVersionOk(); ok {
			dvrTypesMap := map[string]attr.Type{
				"type":                       types.StringType,
				"distance_from_latest_major": types.Int64Type,
				"latest_security_patch":      types.BoolType,
			}
			dynamicVersionRequirement := types.ObjectValueMust(dvrTypesMap, map[string]attr.Value{
				"type":                       types.StringPointerValue(nil),
				"distance_from_latest_major": types.Int64PointerValue(nil),
				"latest_security_patch":      types.BoolPointerValue(nil),
			})
			osvValueMap := map[string]attr.Value{
				"minimum":                     types.StringPointerValue(osv.Minimum),
				"dynamic_version_requirement": dynamicVersionRequirement,
			}
			osvTypesMap := map[string]attr.Type{
				"minimum": types.StringType,
				"dynamic_version_requirement": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type":                       types.StringType,
						"distance_from_latest_major": types.Int64Type,
						"latest_security_patch":      types.BoolType,
					},
				},
			}
			if dvr, ok := osv.GetDynamicVersionRequirementOk(); ok {
				dvrValueMap := map[string]attr.Value{
					"type":                       types.StringPointerValue(dvr.Type),
					"distance_from_latest_major": types.Int64Value(int64(dvr.GetDistanceFromLatestMajor())),
					"latest_security_patch":      types.BoolPointerValue(dvr.LatestSecurityPatch),
				}
				dynamicVersionRequirement, diag := types.ObjectValue(dvrTypesMap, dvrValueMap)
				if diag != nil {
					diags.Append(diag...)
					return diags
				}
				osvValueMap["dynamic_version_requirement"] = dynamicVersionRequirement
			}
			osVersion, diag := types.ObjectValue(osvTypesMap, osvValueMap)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			state.OSVersion = osVersion
		}
		if slt, ok := v.GetScreenLockTypeOk(); ok {
			include := make([]attr.Value, 0)
			for _, i := range slt.Include {
				include = append(include, types.StringValue(i))
			}
			includeValue, diag := types.SetValue(types.StringType, include)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			sltValueMap := map[string]attr.Value{
				"include": includeValue,
			}
			sltTypesMap := map[string]attr.Type{
				"include": types.SetType{
					ElemType: types.StringType,
				},
			}
			screenlockType, diag := types.ObjectValue(sltTypesMap, sltValueMap)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			state.ScreenlockType = screenlockType
		}
		state.Jaibreak = types.BoolPointerValue(v.Jailbreak)
	case *okta.DeviceAssuranceMacOSPlatform:
		if det, ok := v.GetDiskEncryptionTypeOk(); ok {
			include := make([]attr.Value, 0)
			for _, i := range det.Include {
				include = append(include, types.StringValue(i))
			}
			includeValue, diag := types.SetValue(types.StringType, include)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			detValueMap := map[string]attr.Value{
				"include": includeValue,
			}
			detTypesMap := map[string]attr.Type{
				"include": types.SetType{
					ElemType: types.StringType,
				},
			}
			diskEncryptionType, diag := types.ObjectValue(detTypesMap, detValueMap)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			state.DiskEncryptionType = diskEncryptionType
		}
		if osv, ok := v.GetOsVersionOk(); ok {
			dvrTypesMap := map[string]attr.Type{
				"type":                       types.StringType,
				"distance_from_latest_major": types.Int64Type,
				"latest_security_patch":      types.BoolType,
			}
			dynamicVersionRequirement := types.ObjectValueMust(dvrTypesMap, map[string]attr.Value{
				"type":                       types.StringPointerValue(nil),
				"distance_from_latest_major": types.Int64PointerValue(nil),
				"latest_security_patch":      types.BoolPointerValue(nil),
			})
			osvValueMap := map[string]attr.Value{
				"minimum":                     types.StringPointerValue(osv.Minimum),
				"dynamic_version_requirement": dynamicVersionRequirement,
			}
			osvTypesMap := map[string]attr.Type{
				"minimum": types.StringType,
				"dynamic_version_requirement": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type":                       types.StringType,
						"distance_from_latest_major": types.Int64Type,
						"latest_security_patch":      types.BoolType,
					},
				},
			}
			if dvr, ok := osv.GetDynamicVersionRequirementOk(); ok {
				dvrValueMap := map[string]attr.Value{
					"type":                       types.StringPointerValue(dvr.Type),
					"distance_from_latest_major": types.Int64Value(int64(dvr.GetDistanceFromLatestMajor())),
					"latest_security_patch":      types.BoolPointerValue(dvr.LatestSecurityPatch),
				}
				dynamicVersionRequirement, diag := types.ObjectValue(dvrTypesMap, dvrValueMap)
				if diag != nil {
					diags.Append(diag...)
					return diags
				}
				osvValueMap["dynamic_version_requirement"] = dynamicVersionRequirement
			}
			osVersion, diag := types.ObjectValue(osvTypesMap, osvValueMap)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			state.OSVersion = osVersion
		}
		if slt, ok := v.GetScreenLockTypeOk(); ok {
			include := make([]attr.Value, 0)
			for _, i := range slt.Include {
				include = append(include, types.StringValue(i))
			}
			includeValue, diag := types.SetValue(types.StringType, include)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			sltValueMap := map[string]attr.Value{
				"include": includeValue,
			}
			sltTypesMap := map[string]attr.Type{
				"include": types.SetType{
					ElemType: types.StringType,
				},
			}
			screenlockType, diag := types.ObjectValue(sltTypesMap, sltValueMap)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			state.ScreenlockType = screenlockType
		}
		if tpsp, ok := v.GetThirdPartySignalProvidersOk(); ok {
			tpspTypesMap, dtcTypesMap, bvTypesMap, osvTypesMap, tpspValueMap, dtcValueMap, bvValueMap, osvValueMap := initEmptyTpsp()
			if dtc, ok := tpsp.GetDtcOk(); ok {
				dtcValueMap["built_in_dns_client_enabled"] = types.BoolPointerValue(dtc.BuiltInDnsClientEnabled)
				dtcValueMap["chrome_remote_desktop_app_blocked"] = types.BoolPointerValue(dtc.ChromeRemoteDesktopAppBlocked)
				dtcValueMap["device_enrollment_domain"] = types.StringPointerValue(dtc.DeviceEnrollmentDomain)
				dtcValueMap["disk_encrypted"] = types.BoolPointerValue(dtc.DiskEncrypted)
				dtcValueMap["key_trust_level"] = types.StringPointerValue(dtc.KeyTrustLevel)
				dtcValueMap["os_firewall"] = types.BoolPointerValue(dtc.OsFirewall)
				dtcValueMap["password_protection_warning_trigger"] = types.StringPointerValue(dtc.PasswordProtectionWarningTrigger)
				dtcValueMap["realtime_url_check_mode"] = types.BoolPointerValue(dtc.RealtimeUrlCheckMode)
				dtcValueMap["safe_browsing_protection_level"] = types.StringPointerValue(dtc.SafeBrowsingProtectionLevel)
				dtcValueMap["screen_lock_secured"] = types.BoolPointerValue(dtc.ScreenLockSecured)
				dtcValueMap["site_isolation_enabled"] = types.BoolPointerValue(dtc.SiteIsolationEnabled)
				if bv, ok := dtc.GetBrowserVersionOk(); ok {
					bvValueMap["minimum"] = types.StringPointerValue(bv.Minimum)
					browserVersion, diag := types.ObjectValue(bvTypesMap, bvValueMap)
					if diag != nil {
						diags.Append(diag...)
						return diags
					}
					dtcValueMap["browser_version"] = browserVersion
				}
				if osv, ok := dtc.GetOsVersionOk(); ok {
					osvValueMap["minimum"] = types.StringPointerValue(osv.Minimum)
					osVersion, diag := types.ObjectValue(osvTypesMap, osvValueMap)
					if diag != nil {
						diags.Append(diag...)
						return diags
					}
					dtcValueMap["os_version"] = osVersion
				}
				dtc, diag := types.ObjectValue(dtcTypesMap, dtcValueMap)
				if diag != nil {
					diags.Append(diag...)
					return diags
				}
				tpspValueMap["dtc"] = dtc
			}
			thirdPartySignalProvider, diag := types.ObjectValue(tpspTypesMap, tpspValueMap)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			state.ThirdPartySignalProvider = thirdPartySignalProvider
		}
		state.SecureHardwarePresent = types.BoolPointerValue(v.SecureHardwarePresent)
	case *okta.DeviceAssuranceWindowsPlatform:
		if det, ok := v.GetDiskEncryptionTypeOk(); ok {
			include := make([]attr.Value, 0)
			for _, i := range det.Include {
				include = append(include, types.StringValue(i))
			}
			includeValue, diag := types.SetValue(types.StringType, include)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			detValueMap := map[string]attr.Value{
				"include": includeValue,
			}
			detTypesMap := map[string]attr.Type{
				"include": types.SetType{
					ElemType: types.StringType,
				},
			}
			diskEncryptionType, diag := types.ObjectValue(detTypesMap, detValueMap)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			state.DiskEncryptionType = diskEncryptionType
		}
		if osv, ok := v.GetOsVersionOk(); ok {
			dvrTypesMap := map[string]attr.Type{
				"type":                       types.StringType,
				"distance_from_latest_major": types.Int64Type,
				"latest_security_patch":      types.BoolType,
			}
			dynamicVersionRequirement := types.ObjectValueMust(dvrTypesMap, map[string]attr.Value{
				"type":                       types.StringPointerValue(nil),
				"distance_from_latest_major": types.Int64PointerValue(nil),
				"latest_security_patch":      types.BoolPointerValue(nil),
			})
			osvValueMap := map[string]attr.Value{
				"minimum":                     types.StringPointerValue(osv.Minimum),
				"dynamic_version_requirement": dynamicVersionRequirement,
			}
			osvTypesMap := map[string]attr.Type{
				"minimum": types.StringType,
				"dynamic_version_requirement": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type":                       types.StringType,
						"distance_from_latest_major": types.Int64Type,
						"latest_security_patch":      types.BoolType,
					},
				},
			}
			osVersion, diag := types.ObjectValue(osvTypesMap, osvValueMap)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			state.OSVersion = osVersion
		}
		if slt, ok := v.GetScreenLockTypeOk(); ok {
			include := make([]attr.Value, 0)
			for _, i := range slt.Include {
				include = append(include, types.StringValue(i))
			}
			includeValue, diag := types.SetValue(types.StringType, include)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			sltValueMap := map[string]attr.Value{
				"include": includeValue,
			}
			sltTypesMap := map[string]attr.Type{
				"include": types.SetType{
					ElemType: types.StringType,
				},
			}
			screenlockType, diag := types.ObjectValue(sltTypesMap, sltValueMap)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			state.ScreenlockType = screenlockType
		}
		if tpsp, ok := v.GetThirdPartySignalProvidersOk(); ok {
			tpspTypesMap, dtcTypesMap, bvTypesMap, osvTypesMap, tpspValueMap, dtcValueMap, bvValueMap, osvValueMap := initEmptyTpsp()
			if dtc, ok := tpsp.GetDtcOk(); ok {
				dtcValueMap["built_in_dns_client_enabled"] = types.BoolPointerValue(dtc.BuiltInDnsClientEnabled)
				dtcValueMap["chrome_remote_desktop_app_blocked"] = types.BoolPointerValue(dtc.ChromeRemoteDesktopAppBlocked)
				dtcValueMap["device_enrollment_domain"] = types.StringPointerValue(dtc.DeviceEnrollmentDomain)
				dtcValueMap["disk_encrypted"] = types.BoolPointerValue(dtc.DiskEncrypted)
				dtcValueMap["key_trust_level"] = types.StringPointerValue(dtc.KeyTrustLevel)
				dtcValueMap["os_firewall"] = types.BoolPointerValue(dtc.OsFirewall)
				dtcValueMap["password_protection_warning_trigger"] = types.StringPointerValue(dtc.PasswordProtectionWarningTrigger)
				dtcValueMap["realtime_url_check_mode"] = types.BoolPointerValue(dtc.RealtimeUrlCheckMode)
				dtcValueMap["safe_browsing_protection_level"] = types.StringPointerValue(dtc.SafeBrowsingProtectionLevel)
				dtcValueMap["screen_lock_secured"] = types.BoolPointerValue(dtc.ScreenLockSecured)
				dtcValueMap["site_isolation_enabled"] = types.BoolPointerValue(dtc.SiteIsolationEnabled)
				dtcValueMap["crowd_strike_agent_id"] = types.StringPointerValue(dtc.CrowdStrikeAgentId)
				dtcValueMap["crowd_strike_customer_id"] = types.StringPointerValue(dtc.CrowdStrikeCustomerId)
				dtcValueMap["third_party_blocking_enabled"] = types.BoolPointerValue(dtc.ThirdPartyBlockingEnabled)
				dtcValueMap["windows_machine_domain"] = types.StringPointerValue(dtc.WindowsMachineDomain)
				dtcValueMap["windows_user_domain"] = types.StringPointerValue(dtc.WindowsUserDomain)
				if bv, ok := dtc.GetBrowserVersionOk(); ok {
					bvValueMap["minimum"] = types.StringPointerValue(bv.Minimum)
					browserVersion, diag := types.ObjectValue(bvTypesMap, bvValueMap)
					if diag != nil {
						diags.Append(diag...)
						return diags
					}
					dtcValueMap["browser_version"] = browserVersion
				}
				if osv, ok := dtc.GetOsVersionOk(); ok {
					osvValueMap["minimum"] = types.StringPointerValue(osv.Minimum)
					osVersion, diag := types.ObjectValue(osvTypesMap, osvValueMap)
					if diag != nil {
						diags.Append(diag...)
						return diags
					}
					dtcValueMap["os_version"] = osVersion
				}
				dtc, diag := types.ObjectValue(dtcTypesMap, dtcValueMap)
				if diag != nil {
					diags.Append(diag...)
					return diags
				}
				tpspValueMap["dtc"] = dtc
			}
			thirdPartySignalProvider, diag := types.ObjectValue(tpspTypesMap, tpspValueMap)
			if diag != nil {
				diags.Append(diag...)
				return diags
			}
			state.ThirdPartySignalProvider = thirdPartySignalProvider
		}
		if osvc, ok := v.GetOsVersionConstraintsOk(); ok {
			for _, o := range osvc {
				dvrTypesMap := map[string]attr.Type{
					"type":                       types.StringType,
					"distance_from_latest_major": types.Int64Type,
					"latest_security_patch":      types.BoolType,
				}
				dvrValueMap := map[string]attr.Value{
					"type":                       types.StringPointerValue(nil),
					"distance_from_latest_major": types.Int64PointerValue(nil),
					"latest_security_patch":      types.BoolPointerValue(nil),
				}
				if o.DynamicVersionRequirement != nil {
					dvrValueMap["type"] = types.StringPointerValue(o.DynamicVersionRequirement.Type)
					dvrValueMap["distance_from_latest_major"] = types.Int64Value(int64(o.DynamicVersionRequirement.GetDistanceFromLatestMajor()))
					dvrValueMap["latest_security_patch"] = types.BoolPointerValue(o.DynamicVersionRequirement.LatestSecurityPatch)
				}
				dynamicVersionRequirement, diag := types.ObjectValue(dvrTypesMap, dvrValueMap)
				if diag != nil {
					diags = append(diags, diag...)
				}
				state.OSVersionConstraint = append(state.OSVersionConstraint, OsVersionConstraint{
					MajorVersionConstraint:    types.StringValue(o.MajorVersionConstraint),
					DynamicVersionRequirement: dynamicVersionRequirement,
				})
			}
		}
		state.SecureHardwarePresent = types.BoolPointerValue(v.SecureHardwarePresent)
	}
	state.ID = types.StringValue(sharedData.GetId())
	state.Name = types.StringValue(sharedData.GetName())
	state.Platform = types.StringValue(sharedData.GetPlatform())
	return diags
}

func initEmptyTpsp() (tpspTypesMap, dtcTypesMap, bvTypesMap, osvTypesMap map[string]attr.Type, tpspValueMap, dtcValueMap, bvValueMap, osvValueMap map[string]attr.Value) {
	bvValueMap = map[string]attr.Value{
		"minimum": types.StringPointerValue(nil),
	}
	bvTypesMap = map[string]attr.Type{
		"minimum": types.StringType,
	}
	bv := types.ObjectValueMust(bvTypesMap, bvValueMap)
	osvValueMap = map[string]attr.Value{
		"minimum": types.StringPointerValue(nil),
	}
	osvTypesMap = map[string]attr.Type{
		"minimum": types.StringType,
	}
	osv := types.ObjectValueMust(osvTypesMap, osvValueMap)
	dtcValueMap = map[string]attr.Value{
		"allow_screen_lock":                   types.BoolPointerValue(nil),
		"browser_version":                     bv,
		"built_in_dns_client_enabled":         types.BoolPointerValue(nil),
		"chrome_remote_desktop_app_blocked":   types.BoolPointerValue(nil),
		"device_enrollment_domain":            types.StringPointerValue(nil),
		"disk_encrypted":                      types.BoolPointerValue(nil),
		"key_trust_level":                     types.StringPointerValue(nil),
		"managed_device":                      types.BoolPointerValue(nil),
		"os_firewall":                         types.BoolPointerValue(nil),
		"os_version":                          osv,
		"password_protection_warning_trigger": types.StringPointerValue(nil),
		"realtime_url_check_mode":             types.BoolPointerValue(nil),
		"safe_browsing_protection_level":      types.StringPointerValue(nil),
		"screen_lock_secured":                 types.BoolPointerValue(nil),
		"site_isolation_enabled":              types.BoolPointerValue(nil),
		"crowd_strike_agent_id":               types.StringPointerValue(nil),
		"crowd_strike_customer_id":            types.StringPointerValue(nil),
		"third_party_blocking_enabled":        types.BoolPointerValue(nil),
		"windows_machine_domain":              types.StringPointerValue(nil),
		"windows_user_domain":                 types.StringPointerValue(nil),
	}
	dtcTypesMap = map[string]attr.Type{
		"allow_screen_lock": types.BoolType,
		"browser_version": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"minimum": types.StringType,
			},
		},
		"built_in_dns_client_enabled":       types.BoolType,
		"chrome_remote_desktop_app_blocked": types.BoolType,
		"device_enrollment_domain":          types.StringType,
		"disk_encrypted":                    types.BoolType,
		"key_trust_level":                   types.StringType,
		"managed_device":                    types.BoolType,
		"os_firewall":                       types.BoolType,
		"os_version": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"minimum": types.StringType,
			},
		},
		"password_protection_warning_trigger": types.StringType,
		"realtime_url_check_mode":             types.BoolType,
		"safe_browsing_protection_level":      types.StringType,
		"screen_lock_secured":                 types.BoolType,
		"site_isolation_enabled":              types.BoolType,
		"crowd_strike_agent_id":               types.StringType,
		"crowd_strike_customer_id":            types.StringType,
		"third_party_blocking_enabled":        types.BoolType,
		"windows_machine_domain":              types.StringType,
		"windows_user_domain":                 types.StringType,
	}
	dtc := types.ObjectValueMust(dtcTypesMap, dtcValueMap)
	tpspValueMap = map[string]attr.Value{
		"dtc": dtc,
	}
	tpspTypesMap = map[string]attr.Type{
		"dtc": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"allow_screen_lock": types.BoolType,
				"browser_version": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"minimum": types.StringType,
					},
				},
				"built_in_dns_client_enabled":       types.BoolType,
				"chrome_remote_desktop_app_blocked": types.BoolType,
				"device_enrollment_domain":          types.StringType,
				"disk_encrypted":                    types.BoolType,
				"key_trust_level":                   types.StringType,
				"managed_device":                    types.BoolType,
				"os_firewall":                       types.BoolType,
				"os_version": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"minimum": types.StringType,
					},
				},
				"password_protection_warning_trigger": types.StringType,
				"realtime_url_check_mode":             types.BoolType,
				"safe_browsing_protection_level":      types.StringType,
				"screen_lock_secured":                 types.BoolType,
				"site_isolation_enabled":              types.BoolType,
				"crowd_strike_agent_id":               types.StringType,
				"crowd_strike_customer_id":            types.StringType,
				"third_party_blocking_enabled":        types.BoolType,
				"windows_machine_domain":              types.StringType,
				"windows_user_domain":                 types.StringType,
			},
		},
	}
	return tpspTypesMap, dtcTypesMap, bvTypesMap, osvTypesMap, tpspValueMap, dtcValueMap, bvValueMap, osvValueMap
}

func findDeviceAssurancePolicyByName(ctx context.Context, client *okta.APIClient, name string) (*okta.ListDeviceAssurancePolicies200ResponseInner, error) {
	var res *okta.ListDeviceAssurancePolicies200ResponseInner
	dapListResp, _, err := client.DeviceAssuranceAPI.ListDeviceAssurancePolicies(ctx).Execute()
	if err != nil {
		return nil, err
	}
	for _, dap := range dapListResp {
		data := dap.GetActualInstance().(OktaDeviceAssurancePolicy)
		if strings.EqualFold(name, data.GetName()) {
			res = &dap
			return res, nil
		}
	}
	return nil, fmt.Errorf("user type '%s' does not exist", name)
}
