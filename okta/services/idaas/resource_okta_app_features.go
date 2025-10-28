package idaas

import (
	"context"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &appFeatures{}
	_ resource.ResourceWithConfigure   = &appFeatures{}
	_ resource.ResourceWithImportState = &appFeatures{}
)

var _ resource.Resource = &appFeatures{}

type appFeatures struct {
	*config.Config
}

type Lifecycle struct {
	Status types.String `tfsdk:"status"`
}

type Create struct {
	LifecycleCreate Lifecycle `tfsdk:"lifecycle_create"`
}

type Passsword struct {
	Change types.String `tfsdk:"change"`
	Seed   types.String `tfsdk:"seed"`
	Status types.String `tfsdk:"status"`
}

type Profile struct {
	Status types.String `tfsdk:"status"`
}

type Update struct {
	LifecycleDelete *Lifecycle `tfsdk:"lifecycle_delete"`
	Password        *Passsword `tfsdk:"password"`
	Profile         *Profile   `tfsdk:"profile"`
}

type UserCreateAndMatch struct {
	AllowPartialMatch       types.Bool   `tfsdk:"allow_partial_match"`
	AutoActivateNewUsers    types.Bool   `tfsdk:"auto_activate_new_users"`
	AutoConfirmExactMatch   types.Bool   `tfsdk:"autoconfirm_exact_match"`
	AutoConfirmNewUsers     types.Bool   `tfsdk:"autoconfirm_new_users"`
	AutoConfirmPartialMatch types.Bool   `tfsdk:"autoconfirm_partial_match"`
	ExactMatchCriteria      types.String `tfsdk:"exact_match_criteria"`
}

type ImportRules struct {
	UserCreatAndMatch *UserCreateAndMatch `tfsdk:"user_create_and_match"`
}

type Username struct {
	UsernameFormat     types.String `tfsdk:"username_format"`
	UsernameExpression types.String `tfsdk:"username_expression"`
}

type Import struct {
	Expression types.String `tfsdk:"expression"`
	Timezone   types.String `tfsdk:"timezone"`
}
type Schedule struct {
	FullImport        *Import      `tfsdk:"full_import"`
	IncrementalImport *Import      `tfsdk:"incremental_import"`
	Status            types.String `tfsdk:"status"`
}

type ImportSettings struct {
	Username *Username `tfsdk:"username"`
	Schedule *Schedule `tfsdk:"schedule"`
}

type Capabilities struct {
	Create         *Create         `tfsdk:"create"`
	Update         *Update         `tfsdk:"update"`
	ImportRules    *ImportRules    `tfsdk:"import_rules"`
	ImportSettings *ImportSettings `tfsdk:"import_settings"`
}

type appFeaturesModel struct {
	Id           types.String  `tfsdk:"id"`
	AppId        types.String  `tfsdk:"app_id"`
	Description  types.String  `tfsdk:"description"`
	Name         types.String  `tfsdk:"name"`
	Status       types.String  `tfsdk:"status"`
	Capabilities *Capabilities `tfsdk:"capabilities"`
}

func newAppFeaturesResource() resource.Resource {
	return &appFeatures{}
}

func (r *appFeatures) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_features"
}

func (r *appFeatures) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *appFeatures) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, tfpath.Root("id"), req, resp)
}

func (r *appFeatures) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"app_id": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"status": schema.StringAttribute{
				Optional: true,
			},
		},
		Blocks: map[string]schema.Block{
			"capabilities": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{},
				Blocks: map[string]schema.Block{
					"create": schema.SingleNestedBlock{
						Blocks: map[string]schema.Block{
							"lifecycle_create": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"status": schema.StringAttribute{
										Optional: true,
									},
								},
							},
						},
					},
					"update": schema.SingleNestedBlock{
						Blocks: map[string]schema.Block{
							"lifecycle_delete": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"status": schema.StringAttribute{
										Optional: true,
									},
								},
							},
							"password": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"change": schema.StringAttribute{
										Optional: true,
									},
									"seed": schema.StringAttribute{
										Optional: true,
									},
									"status": schema.StringAttribute{
										Optional: true,
									},
								},
							},
							"profile": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"status": schema.StringAttribute{
										Optional: true,
									},
								},
							},
						},
					},
					"import_rules": schema.SingleNestedBlock{
						Blocks: map[string]schema.Block{
							"user_create_and_match": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"exact_match_criteria": schema.StringAttribute{
										Optional: true,
									},
									"allow_partial_match": schema.BoolAttribute{
										Optional: true,
									},
									"auto_activate_new_users": schema.BoolAttribute{
										Optional: true,
									},
									"autoconfirm_exact_match": schema.BoolAttribute{
										Optional: true,
									},
									"autoconfirm_new_users": schema.BoolAttribute{
										Optional: true,
									},
									"autoconfirm_partial_match": schema.BoolAttribute{
										Optional: true,
									},
								},
							},
						},
					},
					"import_settings": schema.SingleNestedBlock{
						Blocks: map[string]schema.Block{
							"username": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"username_format": schema.StringAttribute{
										Optional: true,
									},
									"username_expression": schema.StringAttribute{
										Optional: true,
									},
								},
							},
							"schedule": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"status": schema.StringAttribute{
										Optional: true,
									},
								},
								Blocks: map[string]schema.Block{
									"full_import": schema.SingleNestedBlock{
										Attributes: map[string]schema.Attribute{
											"expression": schema.StringAttribute{
												Optional: true,
											},
											"timezone": schema.StringAttribute{
												Optional: true,
											},
										},
									},
									"incremental_import": schema.SingleNestedBlock{
										Attributes: map[string]schema.Attribute{
											"expression": schema.StringAttribute{
												Optional: true,
											},
											"timezone": schema.StringAttribute{
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *appFeatures) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data appFeaturesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	updateAppFeatureResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApplicationFeaturesAPI.UpdateFeatureForApplication(ctx, data.AppId.ValueString(), data.Name.ValueString()).UpdateFeatureForApplicationRequest(buildUpdateAppFeature(data)).Execute()
	if err != nil {
		return
	}

	updateAppFeatureState(&data, updateAppFeatureResp)

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func updateAppFeatureState(data *appFeaturesModel, updateAppFeatureResp *v5okta.ListFeaturesForApplication200ResponseInner) {
	data.Capabilities.ImportRules.UserCreatAndMatch.ExactMatchCriteria = types.StringValue(updateAppFeatureResp.InboundProvisioningApplicationFeature.Capabilities.ImportRules.UserCreateAndMatch.GetExactMatchCriteria())
	data.Capabilities.ImportRules.UserCreatAndMatch.AllowPartialMatch = types.BoolValue(updateAppFeatureResp.InboundProvisioningApplicationFeature.Capabilities.ImportRules.UserCreateAndMatch.GetAllowPartialMatch())
	data.Capabilities.ImportRules.UserCreatAndMatch.AutoActivateNewUsers = types.BoolValue(updateAppFeatureResp.InboundProvisioningApplicationFeature.Capabilities.ImportRules.UserCreateAndMatch.GetAutoActivateNewUsers())
	data.Capabilities.ImportRules.UserCreatAndMatch.AutoConfirmExactMatch = types.BoolValue(updateAppFeatureResp.InboundProvisioningApplicationFeature.Capabilities.ImportRules.UserCreateAndMatch.GetAutoConfirmExactMatch())
	data.Capabilities.ImportRules.UserCreatAndMatch.AutoConfirmNewUsers = types.BoolValue(updateAppFeatureResp.InboundProvisioningApplicationFeature.Capabilities.ImportRules.UserCreateAndMatch.GetAutoConfirmNewUsers())
	data.Capabilities.ImportRules.UserCreatAndMatch.AutoConfirmPartialMatch = types.BoolValue(updateAppFeatureResp.InboundProvisioningApplicationFeature.Capabilities.ImportRules.UserCreateAndMatch.GetAutoConfirmPartialMatch())
	data.Capabilities.ImportSettings.Username.UsernameFormat = types.StringValue(updateAppFeatureResp.InboundProvisioningApplicationFeature.Capabilities.ImportSettings.Username.GetUsernameFormat())
	data.Capabilities.ImportSettings.Username.UsernameExpression = types.StringValue(updateAppFeatureResp.InboundProvisioningApplicationFeature.Capabilities.ImportSettings.Username.GetUserNameExpression())

	data.Capabilities.Create.LifecycleCreate.Status = types.StringValue(updateAppFeatureResp.UserProvisioningApplicationFeature.Capabilities.Create.LifecycleCreate.GetStatus())
	data.Capabilities.Update.LifecycleDelete.Status = types.StringValue(updateAppFeatureResp.UserProvisioningApplicationFeature.Capabilities.Update.LifecycleDeactivate.GetStatus())

	data.Capabilities.Update.Password.Change = types.StringValue(updateAppFeatureResp.UserProvisioningApplicationFeature.Capabilities.Update.Password.GetChange())
	data.Capabilities.Update.Password.Seed = types.StringValue(updateAppFeatureResp.UserProvisioningApplicationFeature.Capabilities.Update.Password.GetSeed())
	data.Capabilities.Update.Password.Status = types.StringValue(updateAppFeatureResp.UserProvisioningApplicationFeature.Capabilities.Update.Password.GetStatus())

	data.Capabilities.Update.Profile.Status = types.StringValue(updateAppFeatureResp.UserProvisioningApplicationFeature.Capabilities.Update.Profile.GetStatus())
	data.Description = types.StringValue(updateAppFeatureResp.InboundProvisioningApplicationFeature.GetDescription())
	data.Status = types.StringValue(updateAppFeatureResp.UserProvisioningApplicationFeature.GetStatus())
	data.Name = types.StringValue(updateAppFeatureResp.UserProvisioningApplicationFeature.GetName())
	data.Id = types.StringValue(data.AppId.String() + "/" + data.Name.ValueString())
}

func buildUpdateAppFeature(data appFeaturesModel) v5okta.UpdateFeatureForApplicationRequest {
	var updateFeatureForApplicationRequest v5okta.UpdateFeatureForApplicationRequest
	updateFeatureForApplicationRequest.CapabilitiesObject = &v5okta.CapabilitiesObject{}
	if data.Capabilities != nil && data.Capabilities.Create != nil {
		updateFeatureForApplicationRequest.CapabilitiesObject.Create = &v5okta.CapabilitiesCreateObject{
			LifecycleCreate: &v5okta.LifecycleCreateSettingObject{
				Status: data.Capabilities.Create.LifecycleCreate.Status.ValueStringPointer(),
			},
		}
	}

	if data.Capabilities != nil && data.Capabilities.Update != nil {
		if data.Capabilities.Update.LifecycleDelete != nil {
			updateFeatureForApplicationRequest.CapabilitiesObject.Update = &v5okta.CapabilitiesUpdateObject{
				LifecycleDeactivate: &v5okta.LifecycleDeactivateSettingObject{
					Status: data.Capabilities.Update.LifecycleDelete.Status.ValueStringPointer(),
				},
			}
		}

		if data.Capabilities.Update.Password != nil {
			updateFeatureForApplicationRequest.CapabilitiesObject.Update.Password = &v5okta.PasswordSettingObject{
				Change: data.Capabilities.Update.Password.Change.ValueStringPointer(),
				Seed:   data.Capabilities.Update.Password.Seed.ValueStringPointer(),
				Status: data.Capabilities.Update.Password.Status.ValueStringPointer(),
			}
		}

		if data.Capabilities.Update.Profile != nil {
			updateFeatureForApplicationRequest.CapabilitiesObject.Update.Profile = &v5okta.ProfileSettingObject{
				Status: data.Capabilities.Update.Profile.Status.ValueStringPointer(),
			}
		}
	}

	if data.Capabilities != nil && data.Capabilities.ImportRules != nil {
		if data.Capabilities.ImportRules.UserCreatAndMatch != nil {
			if data.Capabilities.ImportRules.UserCreatAndMatch.AllowPartialMatch != types.BoolNull() {
				updateFeatureForApplicationRequest.CapabilitiesInboundProvisioningObject.ImportRules.UserCreateAndMatch.AllowPartialMatch = data.Capabilities.ImportRules.UserCreatAndMatch.AllowPartialMatch.ValueBoolPointer()
			}
			if data.Capabilities.ImportRules.UserCreatAndMatch.AutoActivateNewUsers != types.BoolNull() {
				updateFeatureForApplicationRequest.CapabilitiesInboundProvisioningObject.ImportRules.UserCreateAndMatch.AutoActivateNewUsers = data.Capabilities.ImportRules.UserCreatAndMatch.AutoActivateNewUsers.ValueBoolPointer()
			}
			if data.Capabilities.ImportRules.UserCreatAndMatch.AutoConfirmExactMatch != types.BoolNull() {
				updateFeatureForApplicationRequest.CapabilitiesInboundProvisioningObject.ImportRules.UserCreateAndMatch.AutoConfirmExactMatch = data.Capabilities.ImportRules.UserCreatAndMatch.AutoConfirmExactMatch.ValueBoolPointer()
			}
			if data.Capabilities.ImportRules.UserCreatAndMatch.AutoConfirmNewUsers != types.BoolNull() {
				updateFeatureForApplicationRequest.CapabilitiesInboundProvisioningObject.ImportRules.UserCreateAndMatch.AutoConfirmNewUsers = data.Capabilities.ImportRules.UserCreatAndMatch.AutoConfirmNewUsers.ValueBoolPointer()
			}
			if data.Capabilities.ImportRules.UserCreatAndMatch.AutoConfirmPartialMatch != types.BoolNull() {
				updateFeatureForApplicationRequest.CapabilitiesInboundProvisioningObject.ImportRules.UserCreateAndMatch.AutoConfirmPartialMatch = data.Capabilities.ImportRules.UserCreatAndMatch.AutoConfirmPartialMatch.ValueBoolPointer()
			}
			if data.Capabilities.ImportRules.UserCreatAndMatch.ExactMatchCriteria != types.StringNull() {
				updateFeatureForApplicationRequest.CapabilitiesInboundProvisioningObject.ImportRules.UserCreateAndMatch.ExactMatchCriteria = data.Capabilities.ImportRules.UserCreatAndMatch.ExactMatchCriteria.ValueStringPointer()
			}
		}
	}

	if data.Capabilities != nil && data.Capabilities.ImportSettings != nil {
		if data.Capabilities.ImportSettings.Username != nil {
			updateFeatureForApplicationRequest.CapabilitiesInboundProvisioningObject.ImportSettings.Username = &v5okta.ImportUsernameObject{
				UsernameFormat:     data.Capabilities.ImportSettings.Username.UsernameFormat.ValueString(),
				UserNameExpression: data.Capabilities.ImportSettings.Username.UsernameExpression.ValueStringPointer(),
			}
		}

		if data.Capabilities.ImportSettings.Schedule != nil {
			if data.Capabilities.ImportSettings.Schedule.FullImport != nil {
				updateFeatureForApplicationRequest.CapabilitiesInboundProvisioningObject.ImportSettings.Schedule.FullImport = &v5okta.ImportScheduleObjectFullImport{
					Expression: data.Capabilities.ImportSettings.Schedule.FullImport.Expression.ValueString(),
					Timezone:   data.Capabilities.ImportSettings.Schedule.FullImport.Timezone.ValueStringPointer(),
				}
			}

			if data.Capabilities.ImportSettings.Schedule.IncrementalImport != nil {
				updateFeatureForApplicationRequest.CapabilitiesInboundProvisioningObject.ImportSettings.Schedule.IncrementalImport = &v5okta.ImportScheduleObjectIncrementalImport{
					Expression: data.Capabilities.ImportSettings.Schedule.IncrementalImport.Expression.ValueString(),
					Timezone:   data.Capabilities.ImportSettings.Schedule.IncrementalImport.Timezone.ValueStringPointer(),
				}
			}

			updateFeatureForApplicationRequest.CapabilitiesInboundProvisioningObject.ImportSettings.Schedule.Status = data.Capabilities.ImportSettings.Schedule.Status.ValueStringPointer()
		}
	}

	return updateFeatureForApplicationRequest
}

func (r *appFeatures) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data appFeaturesModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getAppFeatureResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApplicationFeaturesAPI.GetFeatureForApplication(ctx, data.Id.ValueString(), data.Name.ValueString()).Execute()
	if err != nil {
		return
	}

	updateAppFeatureState(&data, getAppFeatureResp)

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *appFeatures) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data appFeaturesModel

	// Read Terraform plan Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	updateAppFeatureResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApplicationFeaturesAPI.UpdateFeatureForApplication(ctx, data.Id.ValueString(), data.Name.ValueString()).UpdateFeatureForApplicationRequest(buildUpdateAppFeature(data)).Execute()
	if err != nil {
		return
	}

	updateAppFeatureState(&data, updateAppFeatureResp)
	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *appFeatures) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"This resource cannot be deleted via Terraform.",
	)
}
