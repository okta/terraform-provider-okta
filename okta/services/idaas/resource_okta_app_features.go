package idaas

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	LifecycleCreate *Lifecycle `tfsdk:"lifecycle_create"`
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
	LifecycleDeactivate *Lifecycle `tfsdk:"lifecycle_deactivate"`
	Password            *Passsword `tfsdk:"password"`
	Profile             *Profile   `tfsdk:"profile"`
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
	ID           types.String  `tfsdk:"id"`
	AppID        types.String  `tfsdk:"app_id"`
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
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected format: resource_id/sequence_id",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("app_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[1])...)
}

func (r *appFeatures) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "`id` used to specify the app feature ID. Its a combination of `app_id` and `name` separated by a forward slash (/).",
			},
			"app_id": schema.StringAttribute{
				Required:    true,
				Description: "`app_id` used to specify the app ID",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the feature.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Key name of the feature.",
				Validators: []validator.String{
					stringvalidator.OneOf("USER_PROVISIONING", "INBOUND_PROVISIONING"),
				},
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "Setting status.",
				Validators: []validator.String{
					stringvalidator.OneOf("DISABLED", "ENABLED"),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"capabilities": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"create": schema.SingleNestedBlock{
						Blocks: map[string]schema.Block{
							"lifecycle_create": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"status": schema.StringAttribute{
										Optional:    true,
										Description: "Setting status.",
										Validators: []validator.String{
											stringvalidator.OneOf("DISABLED", "ENABLED"),
										},
									},
								},
								Description: "Determines whether to update a user in the app when a user in Okta is updated.",
							},
						},
					},
					"update": schema.SingleNestedBlock{
						Blocks: map[string]schema.Block{
							"lifecycle_deactivate": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"status": schema.StringAttribute{
										Optional:    true,
										Description: "Setting status.",
										Validators: []validator.String{
											stringvalidator.OneOf("DISABLED", "ENABLED"),
										},
									},
								},
								Description: "Determines whether deprovisioning occurs when the app is unassigned.",
							},
							"password": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"change": schema.StringAttribute{
										Optional:    true,
										Description: "Determines whether a change in a user's password also updates the user's password in the app.",
										Validators: []validator.String{
											stringvalidator.OneOf("CHANGE", "KEEP_EXISTING"),
										},
									},
									"seed": schema.StringAttribute{
										Optional:    true,
										Description: "Determines whether the generated password is the user's Okta password or a randomly generated password.",
										Validators: []validator.String{
											stringvalidator.OneOf("OKTA", "RANDOM"),
										},
									},
									"status": schema.StringAttribute{
										Optional:    true,
										Description: "Setting status.",
										Validators: []validator.String{
											stringvalidator.OneOf("DISABLED", "ENABLED"),
										},
									},
								},
								Description: "Determines whether Okta creates and pushes a password in the app for each assigned user.",
							},
							"profile": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"status": schema.StringAttribute{
										Optional:    true,
										Description: "Setting status.",
										Validators: []validator.String{
											stringvalidator.OneOf("DISABLED", "ENABLED"),
										},
									},
								},
								Description: "Determines whether updates to a user's profile are pushed to the app.",
							},
						},
					},
					"import_rules": schema.SingleNestedBlock{
						Blocks: map[string]schema.Block{
							"user_create_and_match": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"exact_match_criteria": schema.StringAttribute{
										Optional:    true,
										Description: "Determines the attribute to match users.",
									},
									"allow_partial_match": schema.BoolAttribute{
										Optional:    true,
										Description: "Allows user import upon partial matching. Partial matching occurs when the first and last names of an imported user match those of an existing Okta user, even if the username or email attributes don't match.",
									},
									"auto_activate_new_users": schema.BoolAttribute{
										Optional:    true,
										Description: "If set to true, imported new users are automatically activated.",
									},
									"autoconfirm_exact_match": schema.BoolAttribute{
										Optional:    true,
										Description: "If set to true, exact-matched users are automatically confirmed on activation. If set to false, exact-matched users need to be confirmed manually.",
									},
									"autoconfirm_new_users": schema.BoolAttribute{
										Optional:    true,
										Description: "If set to true, imported new users are automatically confirmed on activation. This doesn't apply to imported users that already exist in Okta.",
									},
									"autoconfirm_partial_match": schema.BoolAttribute{
										Optional:    true,
										Description: "If set to true, partially matched users are automatically confirmed on activation. If set to false, partially matched users need to be confirmed manually.",
									},
								},
								Description: "Rules for matching and creating users.",
							},
						},
					},
					"import_settings": schema.SingleNestedBlock{
						Blocks: map[string]schema.Block{
							"username": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"username_format": schema.StringAttribute{
										Optional:    true,
										Description: "Determines the username format when users sign in to Okta.",
									},
									"username_expression": schema.StringAttribute{
										Optional:    true,
										Description: "For usernameFormat=CUSTOM, specifies the Okta Expression Language statement for a username format that imported users use to sign in to Okta.",
									},
								},
							},
							"schedule": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"status": schema.StringAttribute{
										Optional:    true,
										Description: "Setting status.",
										Validators: []validator.String{
											stringvalidator.OneOf("DISABLED", "ENABLED"),
										},
									},
								},
								Blocks: map[string]schema.Block{
									"full_import": schema.SingleNestedBlock{
										Attributes: map[string]schema.Attribute{
											"expression": schema.StringAttribute{
												Optional:    true,
												Description: "The import schedule in UNIX cron format.",
											},
											"timezone": schema.StringAttribute{
												Optional:    true,
												Description: "The import schedule time zone in Internet Assigned Numbers Authority (IANA) time zone name format.",
											},
										},
									},
									"incremental_import": schema.SingleNestedBlock{
										Attributes: map[string]schema.Attribute{
											"expression": schema.StringAttribute{
												Optional:    true,
												Description: "The import schedule in UNIX cron format.",
											},
											"timezone": schema.StringAttribute{
												Optional:    true,
												Description: "The import schedule time zone in Internet Assigned Numbers Authority (IANA) time zone name format.",
											},
										},
										Description: "Determines the incremental import schedule.",
									},
								},
							},
						},
						Description: "Defines import settings.",
					},
				},
				Description: "Defines the configurations for the USER_PROVISIONING/INBOUND_PROVISIONING feature.",
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
	updateAppFeatureResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApplicationFeaturesAPI.UpdateFeatureForApplication(ctx, data.AppID.ValueString(), data.Name.ValueString()).UpdateFeatureForApplicationRequest(buildUpdateAppFeature(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating app features",
			"Could not create app feature, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(updateAppFeatureState(&data, updateAppFeatureResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func updateAppFeatureState(data *appFeaturesModel, updateAppFeatureResp *v5okta.ListFeaturesForApplication200ResponseInner) diag.Diagnostics {
	var diags diag.Diagnostics

	// Initialize Capabilities if nil
	if data.Capabilities == nil {
		data.Capabilities = &Capabilities{}
	}

	// Handle InboundProvisioningApplicationFeature
	if updateAppFeatureResp.InboundProvisioningApplicationFeature != nil {
		diags.Append(updateInboundProvisioningFeature(data, updateAppFeatureResp.InboundProvisioningApplicationFeature)...)
	}

	// Handle UserProvisioningApplicationFeature
	if updateAppFeatureResp.UserProvisioningApplicationFeature != nil {
		diags.Append(updateUserProvisioningFeature(data, updateAppFeatureResp.UserProvisioningApplicationFeature)...)
	}

	// Set ID
	data.ID = types.StringValue(data.AppID.ValueString() + "/" + data.Name.ValueString())

	return diags
}

func updateInboundProvisioningFeature(data *appFeaturesModel, feature *v5okta.InboundProvisioningApplicationFeature) diag.Diagnostics {
	var diags diag.Diagnostics

	// Initialize ImportRules if nil
	if data.Capabilities.ImportRules == nil {
		data.Capabilities.ImportRules = &ImportRules{}
	}

	// Handle ImportRules.UserCreateAndMatch
	if feature.Capabilities != nil &&
		&feature.Capabilities.ImportRules != nil &&
		feature.Capabilities.ImportRules.UserCreateAndMatch != nil {

		data.Capabilities.ImportRules.UserCreatAndMatch = &UserCreateAndMatch{
			ExactMatchCriteria:      types.StringValue(feature.Capabilities.ImportRules.UserCreateAndMatch.GetExactMatchCriteria()),
			AllowPartialMatch:       types.BoolValue(feature.Capabilities.ImportRules.UserCreateAndMatch.GetAllowPartialMatch()),
			AutoActivateNewUsers:    types.BoolValue(feature.Capabilities.ImportRules.UserCreateAndMatch.GetAutoActivateNewUsers()),
			AutoConfirmExactMatch:   types.BoolValue(feature.Capabilities.ImportRules.UserCreateAndMatch.GetAutoConfirmExactMatch()),
			AutoConfirmNewUsers:     types.BoolValue(feature.Capabilities.ImportRules.UserCreateAndMatch.GetAutoConfirmNewUsers()),
			AutoConfirmPartialMatch: types.BoolValue(feature.Capabilities.ImportRules.UserCreateAndMatch.GetAutoConfirmPartialMatch()),
		}
	}

	// Initialize ImportSettings if nil
	if data.Capabilities.ImportSettings == nil {
		data.Capabilities.ImportSettings = &ImportSettings{}
	}

	if feature.Capabilities != nil {
		// Handle ImportSettings.Username
		if &feature.Capabilities.ImportSettings != nil && feature.Capabilities.ImportSettings.Username != nil {
			data.Capabilities.ImportSettings.Username = &Username{
				UsernameFormat:     types.StringValue(feature.Capabilities.ImportSettings.Username.GetUsernameFormat()),
				UsernameExpression: types.StringValue(feature.Capabilities.ImportSettings.Username.GetUserNameExpression()),
			}
		}
		// Handle ImportSettings.schedule
		if &feature.Capabilities.ImportSettings != nil &&
			feature.Capabilities.ImportSettings.Schedule != nil {
			data.Capabilities.ImportSettings.Schedule = &Schedule{
				Status: types.StringValue(feature.Capabilities.ImportSettings.Schedule.GetStatus()),
				FullImport: &Import{
					Expression: types.StringValue(feature.Capabilities.ImportSettings.Schedule.FullImport.GetExpression()),
					Timezone:   types.StringValue(feature.Capabilities.ImportSettings.Schedule.FullImport.GetTimezone()),
				},
				IncrementalImport: &Import{
					Expression: types.StringValue(feature.Capabilities.ImportSettings.Schedule.IncrementalImport.GetExpression()),
					Timezone:   types.StringValue(feature.Capabilities.ImportSettings.Schedule.IncrementalImport.GetTimezone()),
				},
			}
		}
	}

	// Set description
	data.Description = types.StringValue(feature.GetDescription())
	data.Status = types.StringValue(feature.GetStatus())
	data.Name = types.StringValue(feature.GetName())

	return diags
}

func updateUserProvisioningFeature(data *appFeaturesModel, feature *v5okta.UserProvisioningApplicationFeature) diag.Diagnostics {
	var diags diag.Diagnostics

	// Initialize Create if nil
	if data.Capabilities.Create == nil {
		data.Capabilities.Create = &Create{}
	}

	// Handle Create.LifecycleCreate
	if feature.Capabilities != nil &&
		feature.Capabilities.Create != nil &&
		feature.Capabilities.Create.LifecycleCreate != nil {

		data.Capabilities.Create.LifecycleCreate = &Lifecycle{
			Status: types.StringValue(feature.Capabilities.Create.LifecycleCreate.GetStatus()),
		}
	}

	// Initialize Update if nil
	if data.Capabilities.Update == nil {
		data.Capabilities.Update = &Update{}
	}

	// Handle Update capabilities
	if feature.Capabilities != nil && feature.Capabilities.Update != nil {
		// LifecycleDeactivate
		if feature.Capabilities.Update.LifecycleDeactivate != nil {
			data.Capabilities.Update.LifecycleDeactivate = &Lifecycle{
				Status: types.StringValue(feature.Capabilities.Update.LifecycleDeactivate.GetStatus()),
			}
		}

		// Password
		if feature.Capabilities.Update.Password != nil {
			data.Capabilities.Update.Password = &Passsword{
				Change: types.StringValue(feature.Capabilities.Update.Password.GetChange()),
				Seed:   types.StringValue(feature.Capabilities.Update.Password.GetSeed()),
				Status: types.StringValue(feature.Capabilities.Update.Password.GetStatus()),
			}
		}

		// Profile
		if feature.Capabilities.Update.Profile != nil {
			if data.Capabilities.Update.Profile == nil {
				data.Capabilities.Update.Profile = &Profile{}
			}
			data.Capabilities.Update.Profile.Status = types.StringValue(feature.Capabilities.Update.Profile.GetStatus())
		}
	}

	// Set status and name
	data.Status = types.StringValue(feature.GetStatus())
	data.Name = types.StringValue(feature.GetName())

	return diags
}

func buildUpdateAppFeature(data appFeaturesModel) v5okta.UpdateFeatureForApplicationRequest {
	var req v5okta.UpdateFeatureForApplicationRequest
	req.CapabilitiesObject = &v5okta.CapabilitiesObject{}

	// --- CREATE CAPABILITY ---
	if data.Capabilities != nil && data.Capabilities.Create != nil && data.Capabilities.Create.LifecycleCreate != nil {
		req.CapabilitiesObject.Create = &v5okta.CapabilitiesCreateObject{
			LifecycleCreate: &v5okta.LifecycleCreateSettingObject{
				Status: data.Capabilities.Create.LifecycleCreate.Status.ValueStringPointer(),
			},
		}
	}

	// --- UPDATE CAPABILITY ---
	if data.Capabilities != nil && data.Capabilities.Update != nil {
		req.CapabilitiesObject.Update = &v5okta.CapabilitiesUpdateObject{}

		if data.Capabilities.Update.LifecycleDeactivate != nil {
			req.CapabilitiesObject.Update.LifecycleDeactivate = &v5okta.LifecycleDeactivateSettingObject{
				Status: data.Capabilities.Update.LifecycleDeactivate.Status.ValueStringPointer(),
			}
		}

		if data.Capabilities.Update.Password != nil {
			req.CapabilitiesObject.Update.Password = &v5okta.PasswordSettingObject{
				Status: data.Capabilities.Update.Password.Status.ValueStringPointer(),
				Seed:   data.Capabilities.Update.Password.Seed.ValueStringPointer(),
				Change: data.Capabilities.Update.Password.Change.ValueStringPointer(),
			}
		}

		if data.Capabilities.Update.Profile != nil {
			req.CapabilitiesObject.Update.Profile = &v5okta.ProfileSettingObject{
				Status: data.Capabilities.Update.Profile.Status.ValueStringPointer(),
			}
		}
	}

	// --- INBOUND PROVISIONING (ONLY IF SET) ---
	if data.Capabilities != nil && (data.Capabilities.ImportRules != nil || data.Capabilities.ImportSettings != nil) {
		req.CapabilitiesInboundProvisioningObject = &v5okta.CapabilitiesInboundProvisioningObject{}

		if data.Capabilities.ImportRules != nil && data.Capabilities.ImportRules.UserCreatAndMatch != nil {
			userRules := data.Capabilities.ImportRules.UserCreatAndMatch
			req.CapabilitiesInboundProvisioningObject.ImportRules = v5okta.CapabilitiesImportRulesObject{
				UserCreateAndMatch: &v5okta.CapabilitiesImportRulesUserCreateAndMatchObject{
					AllowPartialMatch:       userRules.AllowPartialMatch.ValueBoolPointer(),
					AutoActivateNewUsers:    userRules.AutoActivateNewUsers.ValueBoolPointer(),
					AutoConfirmExactMatch:   userRules.AutoConfirmExactMatch.ValueBoolPointer(),
					AutoConfirmNewUsers:     userRules.AutoConfirmNewUsers.ValueBoolPointer(),
					AutoConfirmPartialMatch: userRules.AutoConfirmPartialMatch.ValueBoolPointer(),
					ExactMatchCriteria:      userRules.ExactMatchCriteria.ValueStringPointer(),
				},
			}
		}

		if data.Capabilities.ImportSettings != nil {
			req.CapabilitiesInboundProvisioningObject.ImportSettings = v5okta.CapabilitiesImportSettingsObject{}

			if data.Capabilities.ImportSettings.Username != nil {
				req.CapabilitiesInboundProvisioningObject.ImportSettings.Username = &v5okta.ImportUsernameObject{
					UsernameFormat:     data.Capabilities.ImportSettings.Username.UsernameFormat.ValueString(),
					UserNameExpression: data.Capabilities.ImportSettings.Username.UsernameExpression.ValueStringPointer(),
				}
			}

			if data.Capabilities.ImportSettings.Schedule != nil {
				sched := data.Capabilities.ImportSettings.Schedule
				schedObj := &v5okta.ImportScheduleObject{
					Status: sched.Status.ValueStringPointer(),
				}

				if sched.FullImport != nil {
					schedObj.FullImport = &v5okta.ImportScheduleObjectFullImport{
						Expression: sched.FullImport.Expression.ValueString(),
						Timezone:   sched.FullImport.Timezone.ValueStringPointer(),
					}
				}
				if sched.IncrementalImport != nil {
					schedObj.IncrementalImport = &v5okta.ImportScheduleObjectIncrementalImport{
						Expression: sched.IncrementalImport.Expression.ValueString(),
						Timezone:   sched.IncrementalImport.Timezone.ValueStringPointer(),
					}
				}
				req.CapabilitiesInboundProvisioningObject.ImportSettings.Schedule = schedObj
			}
		}
	}

	return req
}

func (r *appFeatures) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data appFeaturesModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getAppFeatureResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApplicationFeaturesAPI.GetFeatureForApplication(ctx, data.AppID.ValueString(), data.Name.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading app features",
			"Could not read app feature, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(updateAppFeatureState(&data, getAppFeatureResp)...)
	if resp.Diagnostics.HasError() {
		return
	}

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
	updateAppFeatureResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApplicationFeaturesAPI.UpdateFeatureForApplication(ctx, data.AppID.ValueString(), data.Name.ValueString()).UpdateFeatureForApplicationRequest(buildUpdateAppFeature(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating app features",
			"Could not update app feature, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(updateAppFeatureState(&data, updateAppFeatureResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *appFeatures) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"This resource cannot be deleted via Terraform.",
	)
}
