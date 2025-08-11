package governance

import (
	"context"
	"example.com/aditya-okta/okta-ig-sdk-golang/governance"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/okta/terraform-provider-okta/okta/config"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*requestSettingResourceResource)(nil)

func NewRequestSettingResourceResource() resource.Resource {
	return &requestSettingResourceResource{}
}

type requestSettingResourceResource struct {
	*config.Config
}

type SupportedTypes struct {
	Type types.String `tfsdk:"type"`
}

type ValidAccessDurationSettings struct {
	MaximumDays    types.Float32    `tfsdk:"maximum_days"`
	MaximumWeeks   types.Float32    `tfsdk:"maximum_weeks"`
	MaximumHours   types.Float32    `tfsdk:"maximum_hours"`
	Required       types.Bool       `tfsdk:"required"`
	SupportedTypes []SupportedTypes `tfsdk:"supported_types"`
}

type ValidSettings struct {
	Supported []SupportedTypes `tfsdk:"supported"`
}

type RequestOnBehalfOfSettings struct {
	Allowed types.Bool `tfsdk:"allowed"`
	OnlyFor types.List `tfsdk:"only_for"`
}

type DefaultSetting struct {
	RequestSubmissionType  types.String           `tfsdk:"request_submission_type"`
	ApprovalSequenceId     types.String           `tfsdk:"approval_sequence_id"`
	AccessDurationSettings AccessDurationSettings `tfsdk:"access_duration"`
	Error                  types.List             `tfsdk:"error"`
}

type RiskSettings struct {
	DefaultSetting DefaultSetting `tfsdk:"default_setting"`
}

type requestSettingResourceResourceModel struct {
	ResourceId                     types.String                `tfsdk:"resource_id"`
	ValidAccessDurationSettings    ValidAccessDurationSettings `tfsdk:"valid_access_duration_settings"`
	ValidAccessScopeSettings       []SupportedTypes            `tfsdk:"valid_access_scope_settings"`
	ValidRequesterSettings         []SupportedTypes            `tfsdk:"valid_requester_settings"`
	RequestOnBehalfOfSettings      RequestOnBehalfOfSettings   `tfsdk:"request_on_behalf_of_settings"`
	RiskSettings                   RiskSettings                `tfsdk:"risk_settings"`
	ValidRequestOnBehalfOfSettings []SupportedTypes            `tfsdk:"valid_request_on_behalf_of_settings"`
	ValidRiskSettings              ValidSettings               `tfsdk:"valid_risk_settings"`
}

func (r *requestSettingResourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_setting_resource"
}

func (r *requestSettingResourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource_id": schema.StringAttribute{
				Required: true,
			},
		},
		Blocks: map[string]schema.Block{
			"valid_access_duration_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"maximum_days": schema.Float32Attribute{
						Computed: true,
					},
					"maximum_weeks": schema.Float32Attribute{
						Computed: true,
					},
					"maximum_hours": schema.Float32Attribute{
						Computed: true,
					},
					"required": schema.BoolAttribute{
						Computed: true,
					},
				},
				Blocks: map[string]schema.Block{
					"supported_types": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.OneOf("ADMIN_FIXED_DURATION", "REQUESTER_SPECIFIED_DURATION"),
									},
								},
							},
						},
					},
				},
			},
			"valid_access_scope_settings": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"ENTITLEMENT_BUNDLES",
									"GROUPS",
									"RESOURCE_DEFAULT",
								),
							},
						},
					},
				},
				Description: "Access scope settings eligible to be added to a request condition.",
			},
			"valid_requester_settings": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"EVERYONE",
									"GROUPS",
									"TEAMS",
								),
							},
						},
					},
				},
				Description: "Access scope settings eligible to be added to a request condition.",
			},
			"requester_on_behalf_of_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"allowed": schema.BoolAttribute{
						Computed: true,
					},
					"only_for": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			"risk_settings": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"default_setting": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"request_submission_type": schema.StringAttribute{
								Optional: true,
							},
							"error": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
							},
							"approval_sequence_id": schema.StringAttribute{
								Optional: true,
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
							},
						},
					},
				},
			},
			"valid_request_on_behalf_of_settings": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"EVERYONE",
									"DIRECT_REPORT",
								),
							},
						},
					},
				},
			},
			"valid_risk_settings": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"supported_types": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.OneOf(
											"DEFAULT_SETTING",
										),
									},
								},
							},
						},
						Description: "Access scope settings eligible to be added to a request condition.",
					},
				},
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

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	reqSettingsResp, _, err := r.OktaGovernanceClient.OktaIGSDKClientV5().RequestSettingsAPI.GetRequestSettingsV2(ctx, data.ResourceId.ValueString()).Execute()
	if err != nil {
		return
	}

	data.ValidAccessDurationSettings = setValidAccessDurationSettings(reqSettingsResp.ValidAccessDurationSettings)
	data.ValidAccessScopeSettings = setValidAccessSettings(reqSettingsResp.ValidAccessScopeSettings)
	data.ValidRequesterSettings = setValidRequesterSettings(reqSettingsResp.ValidRequesterSettings)
	data.RequestOnBehalfOfSettings = setRequesterOnBehalfSettings(reqSettingsResp.RequestOnBehalfOfSettings)
	data.RiskSettings = setRiskSettings(reqSettingsResp.RiskSettings)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func setRiskSettings(settings *governance.RiskSettingsDetails) RiskSettings {

}

func setRequesterOnBehalfSettings(settings *governance.RequestOnBehalfOfSettingsDetails) RequestOnBehalfOfSettings {
	var onBehalfOfSettings RequestOnBehalfOfSettings
	onBehalfOfSettings.Allowed = types.BoolValue(settings.GetAllowed())

	var reqExpVals []attr.Value
	for _, reqExp := range settings.GetOnlyFor() {
		reqExpVals = append(reqExpVals, types.StringValue(string(reqExp)))
	}
	onlyFor, _ := types.ListValue(types.StringType, reqExpVals)
	onBehalfOfSettings.OnlyFor = onlyFor
	return onBehalfOfSettings
}

func setValidRequesterSettings(settings []governance.ValidRequesterSetting) []SupportedTypes {
	var accessSettingsType []SupportedTypes
	for _, setting := range settings {
		accessSettingsType = append(accessSettingsType, SupportedTypes{
			Type: types.StringValue(string(setting.GetType())),
		})
	}
	return accessSettingsType
}

func setValidAccessSettings(settings []governance.ValidAccessDetail) []SupportedTypes {
	var accessSettingsType []SupportedTypes
	for _, setting := range settings {
		accessSettingsType = append(accessSettingsType, SupportedTypes{
			Type: types.StringValue(string(setting.GetType())),
		})
	}
	return accessSettingsType
}

func setValidAccessDurationSettings(settings governance.ValidAccessDurationSettingsDetails) ValidAccessDurationSettings {
	var supportedTypes []SupportedTypes
	for _, supportedType := range settings.SupportedTypes {
		supportedTypes = append(supportedTypes, SupportedTypes{
			Type: types.StringValue(string(supportedType.GetType())),
		})
	}
	return ValidAccessDurationSettings{
		MaximumDays:    types.Float32Value(settings.GetMaximumDays()),
		MaximumHours:   types.Float32Value(settings.GetMaximumHours()),
		MaximumWeeks:   types.Float32Value(settings.GetMaximumWeeks()),
		Required:       types.BoolValue(settings.GetRequired()),
		SupportedTypes: supportedTypes,
	}
}

func (r *requestSettingResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data requestSettingResourceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestSettingResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data requestSettingResourceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
}
