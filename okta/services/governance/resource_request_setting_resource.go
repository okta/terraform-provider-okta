package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &requestSettingResourceResource{}
	_ resource.ResourceWithConfigure   = &requestSettingResourceResource{}
	_ resource.ResourceWithImportState = &requestSettingResourceResource{}
)

func newRequestSettingResourceResource() resource.Resource {
	return &requestSettingResourceResource{}
}

func (r *requestSettingResourceResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r *requestSettingResourceResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

type requestSettingResourceResource struct {
	*config.Config
}

type supportedTypes struct {
	Type types.String `tfsdk:"type"`
}

type validAccessDurationSettings struct {
	MaximumDays    types.Float32    `tfsdk:"maximum_days"`
	MaximumWeeks   types.Float32    `tfsdk:"maximum_weeks"`
	MaximumHours   types.Float32    `tfsdk:"maximum_hours"`
	Required       types.Bool       `tfsdk:"required"`
	SupportedTypes []supportedTypes `tfsdk:"supported_types"`
}

type onlyFor struct {
	Type types.String `tfsdk:"type"`
}

type requestOnBehalfOfSettings struct {
	Allowed types.Bool `tfsdk:"allowed"`
	OnlyFor []onlyFor  `tfsdk:"only_for"`
}

type defaultSetting struct {
	RequestSubmissionType  types.String            `tfsdk:"request_submission_type"`
	ApprovalSequenceId     types.String            `tfsdk:"approval_sequence_id"`
	AccessDurationSettings *AccessDurationSettings `tfsdk:"access_duration_settings"`
	Error                  types.List              `tfsdk:"error"`
}

type riskSettings struct {
	DefaultSetting *defaultSetting `tfsdk:"default_setting"`
}

type requestSettingResourceResourceModel struct {
	Id                        types.String               `tfsdk:"id"`
	RequestOnBehalfOfSettings *requestOnBehalfOfSettings `tfsdk:"request_on_behalf_of_settings"`
	RiskSettings              *riskSettings              `tfsdk:"risk_settings"`
}

func (r *requestSettingResourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_setting_resource"
}

func (r *requestSettingResourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the resource in Okta ID format.",
			},
		},
		Blocks: map[string]schema.Block{
			"request_on_behalf_of_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"allowed": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Indicates that users who can request this resource could also request for another requester of the same resource",
					},
				},
				Blocks: map[string]schema.Block{
					"only_for": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed:    true,
									Optional:    true,
									Description: "Which requesters the resource requester can request on behalf of. If onlyFor is not specified then any requester may request a resource on the behalf of any other user",
								},
							},
						},
					},
				},
				Description: "Specifies if and for whom a requester may request the resource for.",
			},
			"risk_settings": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"default_setting": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"request_submission_type": schema.StringAttribute{
								Optional: true,
								Computed: true,
							},
							"error": schema.ListAttribute{
								ElementType: types.StringType,
								Computed:    true,
							},
							"approval_sequence_id": schema.StringAttribute{
								Optional:    true,
								Description: "The ID of the approval sequence.",
							},
						},
						Blocks: map[string]schema.Block{
							"access_duration_settings": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											stringvalidator.OneOf("ADMIN_FIXED_DURATION", "REQUESTER_SPECIFIED_DURATION"),
										},
									},
									"duration": schema.StringAttribute{
										Optional: true,
									},
								},
								Description: "Settings that control who may specify the access duration allowed by this request condition or risk settings, as well as what duration may be requested.",
							},
						},
						Description: "Default risk settings that are valid for an access request when a risk has been detected for the resource and requesting user.",
					},
				},
				Description: "Risk settings that are valid for an access request when a risk has been detected for the resource and requesting user.",
			},
		},
	}
}

func (r *requestSettingResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddWarning(
		"Create Not Supported",
		"This resource cannot be created via Terraform. Please import it or let Terraform read it from the existing system.",
	)
}

func (r *requestSettingResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data requestSettingResourceResourceModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	reqSettingsResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestSettingsAPI.GetRequestSettingsV2(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Request Settings",
			"Could not read request settings, unexpected error: "+err.Error(),
		)
		return
	}

	applyResourceSettingToState(&data, reqSettingsResp)

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func applyResourceSettingToState(data *requestSettingResourceResourceModel, reqSettingsResp *governance.RequestSettings) {
	data.RequestOnBehalfOfSettings = setRequesterOnBehalfSettings(reqSettingsResp.RequestOnBehalfOfSettings)
	if data.RiskSettings != nil {
		data.RiskSettings = setRiskSettings(reqSettingsResp.RiskSettings)
	}
}

func setRiskSettings(settings *governance.RiskSettingsDetails) *riskSettings {
	defaultSetting := defaultSetting{}
	if settings.DefaultSetting.RiskSettingsDefaultRestrictedDetails != nil {
		defaultSetting.RequestSubmissionType = types.StringValue(settings.DefaultSetting.RiskSettingsDefaultRestrictedDetails.GetRequestSubmissionType())
		accessDurationSettings := setAccessDurationSettings(settings.DefaultSetting.RiskSettingsDefaultRestrictedDetails.GetAccessDurationSettings())
		if accessDurationSettings != nil {
			defaultSetting.AccessDurationSettings = accessDurationSettings
		}
		var errors []attr.Value
		for _, err := range settings.DefaultSetting.RiskSettingsDefaultRestrictedDetails.GetError() {
			errors = append(errors, types.StringValue(string(err)))
		}
		errList, _ := types.ListValue(types.StringType, errors)
		defaultSetting.Error = errList

		return &riskSettings{
			DefaultSetting: &defaultSetting,
		}
	} else if settings.DefaultSetting.RiskSettingsDefaultAllowedWithOverridesDetails != nil {
		defaultSetting.RequestSubmissionType = types.StringValue(settings.DefaultSetting.RiskSettingsDefaultAllowedWithOverridesDetails.GetRequestSubmissionType())
		approvalSequenceId := settings.DefaultSetting.RiskSettingsDefaultAllowedWithOverridesDetails.GetApprovalSequenceId()
		accessDurationSettings := setAccessDurationSettings(settings.DefaultSetting.RiskSettingsDefaultAllowedWithOverridesDetails.GetAccessDurationSettings())
		defaultSetting.ApprovalSequenceId = types.StringValue(approvalSequenceId)
		if accessDurationSettings != nil {
			defaultSetting.AccessDurationSettings = accessDurationSettings
		}
		var errors []attr.Value
		for _, err := range settings.DefaultSetting.RiskSettingsDefaultAllowedWithOverridesDetails.GetError() {
			errors = append(errors, types.StringValue(string(err)))
		}
		errList, _ := types.ListValue(types.StringType, errors)
		defaultSetting.Error = errList
		defaultSetting.RequestSubmissionType = types.StringValue(settings.DefaultSetting.RiskSettingsDefaultAllowedWithOverridesDetails.GetRequestSubmissionType())
		return &riskSettings{
			DefaultSetting: &defaultSetting,
		}
	} else if settings.DefaultSetting.RiskSettingsDefaultAllowedWithNoOverridesDetails != nil {
		defaultSetting.RequestSubmissionType = types.StringValue(settings.DefaultSetting.RiskSettingsDefaultAllowedWithNoOverridesDetails.GetRequestSubmissionType())
		var errors []attr.Value
		for _, err := range settings.DefaultSetting.RiskSettingsDefaultAllowedWithNoOverridesDetails.GetError() {
			errors = append(errors, types.StringValue(string(err)))
		}
		errList, _ := types.ListValue(types.StringType, errors)
		defaultSetting.Error = errList
		return &riskSettings{
			DefaultSetting: &defaultSetting,
		}
	}
	return &riskSettings{}
}

func setRequesterOnBehalfSettings(settings *governance.RequestOnBehalfOfSettingsDetails) *requestOnBehalfOfSettings {
	var onBehalfOfSettings requestOnBehalfOfSettings
	if settings == nil {
		return nil
	}
	onBehalfOfSettings.Allowed = types.BoolValue(settings.GetAllowed())

	var requesters []onlyFor
	if len(settings.GetOnlyFor()) == 0 {
		requesters = nil
		onBehalfOfSettings.OnlyFor = requesters
		return &onBehalfOfSettings
	}
	for _, of := range settings.GetOnlyFor() {
		requesters = append(requesters, onlyFor{
			Type: types.StringValue(string(of)),
		})
	}
	onBehalfOfSettings.OnlyFor = requesters
	return &onBehalfOfSettings
}

func setValidRequesterSettings(settings []governance.ValidRequesterSetting) []supportedTypes {
	var accessSettingsType []supportedTypes
	for _, setting := range settings {
		accessSettingsType = append(accessSettingsType, supportedTypes{
			Type: types.StringValue(string(setting.GetType())),
		})
	}
	return accessSettingsType
}

func setValidAccessSettings(settings []governance.ValidAccessDetail) []supportedTypes {
	var accessSettingsType []supportedTypes
	for _, setting := range settings {
		accessSettingsType = append(accessSettingsType, supportedTypes{
			Type: types.StringValue(string(setting.GetType())),
		})
	}
	return accessSettingsType
}

func setValidAccessDurationSettings(settings governance.ValidAccessDurationSettingsDetails) *validAccessDurationSettings {
	var supportTypes []supportedTypes
	if settings.SupportedTypes != nil {
		for _, supportedType := range settings.GetSupportedTypes() {
			if supportedType.GetType() != "" {
				supportTypes = append(supportTypes, supportedTypes{
					Type: types.StringValue(string(supportedType.GetType())),
				})
			}
		}
	}
	return &validAccessDurationSettings{
		MaximumDays:    types.Float32Value(settings.GetMaximumDays()),
		MaximumHours:   types.Float32Value(settings.GetMaximumHours()),
		MaximumWeeks:   types.Float32Value(settings.GetMaximumWeeks()),
		Required:       types.BoolValue(settings.GetRequired()),
		SupportedTypes: supportTypes,
	}
}

func (r *requestSettingResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state requestSettingResourceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform plan Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	patchedRequestSettingsReq := buildResourceRequestSettingsPatchable(data)
	if patchedRequestSettingsReq == nil {
		resp.Diagnostics.AddError(
			"Error building Request Settings Patchable",
			"Could not build request settings patchable, unexpected error.",
		)
		return
	}

	// Update API call logic
	updatedResourceSettingResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestSettingsAPI.UpdateResourceRequestSettingsV2(ctx, data.Id.ValueString()).ResourceRequestSettingsPatchable(*buildResourceRequestSettingsPatchable(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error update Request Settings",
			"Could not update request settings, unexpected error: "+err.Error())
		return
	}

	applyResourceSettingToState(&data, updatedResourceSettingResp)

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func buildResourceRequestSettingsPatchable(data requestSettingResourceResourceModel) *governance.ResourceRequestSettingsPatchable {
	var patchedRiskSettings *governance.RiskSettingsPatchable // pointer, nil unless set

	if data.RiskSettings != nil && data.RiskSettings.DefaultSetting != nil {
		patchedRiskSettings = &governance.RiskSettingsPatchable{
			DefaultSetting: governance.RiskSettingsDefaultPatchable{},
		}

		reqType := data.RiskSettings.DefaultSetting.RequestSubmissionType.ValueString()
		switch reqType {
		case "RESTRICTED":
			patchedRiskSettings.DefaultSetting.RiskSettingsDefaultRestrictedPatchable = &governance.RiskSettingsDefaultRestrictedPatchable{
				RequestSubmissionType: reqType,
			}
		case "ALLOWED_WITH_OVERRIDES":
			patchedRiskSettings.DefaultSetting.RiskSettingsDefaultAllowedWithOverridesPatchable = &governance.RiskSettingsDefaultAllowedWithOverridesPatchable{
				RequestSubmissionType:  reqType,
				ApprovalSequenceId:     data.RiskSettings.DefaultSetting.ApprovalSequenceId.ValueString(),
				AccessDurationSettings: buildAccessDurationSettings(data.RiskSettings.DefaultSetting.AccessDurationSettings),
			}

		case "ALLOWED_WITH_NO_OVERRIDES":
			patchedRiskSettings.DefaultSetting.RiskSettingsDefaultAllowedWithNoOverridesPatchable = &governance.RiskSettingsDefaultAllowedWithNoOverridesPatchable{
				RequestSubmissionType: reqType,
			}
		}
	}

	var patchedRequestOnBehalfOfSettings *governance.RequestOnBehalfOfSettingsPatchable
	if data.RequestOnBehalfOfSettings != nil {
		var patchedOnlyFor []governance.RequestOnBehalfOfType
		if data.RequestOnBehalfOfSettings.OnlyFor != nil {
			for _, onlyForVal := range data.RequestOnBehalfOfSettings.OnlyFor {
				if !onlyForVal.Type.IsNull() || onlyForVal.Type.ValueString() != "" {
					patchedOnlyFor = append(patchedOnlyFor, governance.RequestOnBehalfOfType(onlyForVal.Type.ValueString()))
				}
			}
		} else {
			patchedOnlyFor = nil
		}
		if !data.RequestOnBehalfOfSettings.Allowed.IsNull() && data.RequestOnBehalfOfSettings.Allowed.ValueBool() {
			patchedRequestOnBehalfOfSettings = &governance.RequestOnBehalfOfSettingsPatchable{
				Allowed: data.RequestOnBehalfOfSettings.Allowed.ValueBool(),
				OnlyFor: patchedOnlyFor,
			}
		} else {
			patchedRequestOnBehalfOfSettings = nil
		}
	}

	resourceRequestSettingsPatchable := &governance.ResourceRequestSettingsPatchable{}
	if patchedRequestOnBehalfOfSettings != nil {
		resourceRequestSettingsPatchable.SetRequestOnBehalfOfSettings(*patchedRequestOnBehalfOfSettings)
	} else {
		resourceRequestSettingsPatchable.SetRequestOnBehalfOfSettingsNil()
	}
	if patchedRiskSettings != nil {
		resourceRequestSettingsPatchable.SetRiskSettings(*patchedRiskSettings)
	}

	return resourceRequestSettingsPatchable
}

func buildAccessDurationSettings(settings *AccessDurationSettings) governance.NullableAccessDurationSettingsPatchable {
	durationSettings := governance.NullableAccessDurationSettingsPatchable{}
	if settings != nil && settings.Type.ValueString() == "ADMIN_FIXED_DURATION" {
		fixedDuration := &governance.AccessDurationSettingsAdminFixedDuration{
			Type:     settings.Type.ValueString(),
			Duration: settings.Duration.ValueString(),
		}
		durationSettings.Set(&governance.AccessDurationSettingsPatchable{
			AccessDurationSettingsAdminFixedDuration: fixedDuration,
		})
	} else if settings != nil && settings.Type.ValueString() == "REQUESTER_SPECIFIED_DURATION" {
		requesterDuration := &governance.AccessDurationSettingsRequesterSpecifiedDuration{
			Type:            settings.Type.ValueString(),
			MaximumDuration: settings.Duration.ValueString(),
		}
		durationSettings.Set(&governance.AccessDurationSettingsPatchable{
			AccessDurationSettingsRequesterSpecifiedDuration: requesterDuration,
		})
	} else {
		durationSettings.Set(nil)
	}
	return durationSettings
}

func (r *requestSettingResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"This resource cannot be deleted via Terraform.",
	)
}
