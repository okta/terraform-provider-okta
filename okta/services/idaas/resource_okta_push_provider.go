package idaas

import (
	"context"

	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/okta/terraform-provider-okta/okta/config"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-sdk-golang/v5/okta"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ resource.Resource                = &pushProviderResource{}
	_ resource.ResourceWithConfigure   = &pushProviderResource{}
	_ resource.ResourceWithImportState = &pushProviderResource{}
)

// pushProviderResource defines the resource implementation
type pushProviderResource struct {
	*config.Config
}

func newPushProvidersResource() resource.Resource {
	return &pushProviderResource{}
}

func (r *pushProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, tfpath.Root("id"), req, resp)
}

type Configuration struct {
	FcmConfiguration       *fcmServiceAccountJsonModel `tfsdk:"fcm_configuration"`
	ApnsConfigurationModel *apnsConfigurationModel     `tfsdk:"apns_configuration"`
}

// pushProviderResourceModel describes the resource data model
type pushProviderResourceModel struct {
	ID              types.String   `tfsdk:"id"`
	Name            types.String   `tfsdk:"name"`
	ProviderType    types.String   `tfsdk:"provider_type"`
	Configuration   *Configuration `tfsdk:"configuration"`
	LastUpdatedDate types.String   `tfsdk:"last_updated_date"`
}

// pushProviderConfigurationModel describes the configuration block
type apnsConfigurationModel struct {
	// APNS Configuration
	KeyID           types.String `tfsdk:"key_id"`
	TeamID          types.String `tfsdk:"team_id"`
	TokenSigningKey types.String `tfsdk:"token_signing_key"`
	FileName        types.String `tfsdk:"file_name"`
}

type ServiceAccount struct {
	Type                    types.String `tfsdk:"type"`
	ProjectID               types.String `tfsdk:"project_id"`
	PrivateKeyID            types.String `tfsdk:"private_key_id"`
	PrivateKey              types.String `tfsdk:"private_key"`
	ClientEmail             types.String `tfsdk:"client_email"`
	ClientId                types.String `tfsdk:"client_id"`
	AuthURI                 types.String `tfsdk:"auth_uri"`
	TokenURI                types.String `tfsdk:"token_uri"`
	AuthProviderX509CertURL types.String `tfsdk:"auth_provider_x509_cert_url"`
	ClientX509CertURL       types.String `tfsdk:"client_x509_cert_url"`
	FileName                types.String `tfsdk:"file_name"`
}

// fcmServiceAccountJsonModel describes the FCM service account JSON structure
type fcmServiceAccountJsonModel struct {
	ServiceAccountJson *ServiceAccount `tfsdk:"service_account_json"`
}

// NewPushProviderResource is a helper function to simplify the provider implementation
func NewPushProviderResource() resource.Resource {
	return &pushProviderResource{}
}

// Metadata returns the resource type name
func (r *pushProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_push_provider"
}

// Configure adds the provider configured client to the resource
func (r *pushProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

// Schema defines the schema for the resource
func (r *pushProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Okta Push Provider. This resource allows you to create and manage push providers for mobile push notifications.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the push provider.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the push provider.",
			},
			"provider_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of push provider. Valid values are `APNS` (Apple Push Notification Service) or `FCM` (Firebase Cloud Messaging).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"last_updated_date": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the push provider was last modified.",
			},
		},
		Blocks: map[string]schema.Block{
			"configuration": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"fcm_configuration": schema.SingleNestedBlock{
						Blocks: map[string]schema.Block{
							"service_account_json": schema.SingleNestedBlock{
								Description: "JSON containing the private service account key and service account details. Required for FCM provider type.",
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Optional:    true,
										Description: "The type of the service account.",
									},
									"project_id": schema.StringAttribute{
										Optional:    true,
										Description: "The project ID.",
									},
									"private_key_id": schema.StringAttribute{
										Optional:    true,
										Description: "The private key ID.",
									},
									"private_key": schema.StringAttribute{
										Optional:    true,
										Sensitive:   true,
										Description: "The private key.",
									},
									"client_email": schema.StringAttribute{
										Optional:    true,
										Description: "The client email.",
									},
									"client_id": schema.StringAttribute{
										Optional:    true,
										Description: "The client ID.",
									},
									"auth_uri": schema.StringAttribute{
										Optional:    true,
										Description: "The auth URI.",
									},
									"token_uri": schema.StringAttribute{
										Optional:    true,
										Description: "The token URI.",
									},
									"auth_provider_x509_cert_url": schema.StringAttribute{
										Optional:    true,
										Description: "The auth provider x509 cert URL.",
									},
									"client_x509_cert_url": schema.StringAttribute{
										Optional:    true,
										Description: "The client x509 cert URL.",
									},
									"file_name": schema.StringAttribute{
										Optional:    true,
										Description: "File name for Admin Console display.",
									},
								},
							},
						},
					},
					"apns_configuration": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"key_id": schema.StringAttribute{
								Optional:    true,
								Description: "10-character Key ID obtained from the Apple developer account. Required for APNS provider type.",
							},
							"team_id": schema.StringAttribute{
								Optional:    true,
								Description: "10-character Team ID used to develop the iOS app. Required for APNS provider type.",
							},
							"token_signing_key": schema.StringAttribute{
								Optional:    true,
								Sensitive:   true,
								Description: "APNs private authentication token signing key. Required for APNS provider type.",
							},
							"file_name": schema.StringAttribute{
								Optional:    true,
								Description: "File name for Admin Console display.",
							},
						},
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state
func (r *pushProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan pushProviderResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the push provider
	createPushProviderResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().PushProviderAPI.CreatePushProvider(ctx).PushProvider(createPushProviderReq(plan)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Push Provider",
			"Could not create push provider, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(applyPushProviderToState(createPushProviderResp, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func applyPushProviderToState(resp *okta.ListPushProviders200ResponseInner, plan *pushProviderResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if resp.APNSPushProvider != nil {
		plan.ID = types.StringValue(resp.APNSPushProvider.GetId())
		plan.Name = types.StringValue(resp.APNSPushProvider.GetName())
		plan.ProviderType = types.StringValue(resp.APNSPushProvider.GetProviderType())
		plan.Configuration.ApnsConfigurationModel.KeyID = types.StringValue(resp.APNSPushProvider.Configuration.GetKeyId())
		plan.Configuration.ApnsConfigurationModel.TeamID = types.StringValue(resp.APNSPushProvider.Configuration.GetTeamId())
		if resp.APNSPushProvider.Configuration.GetFileName() != "" {
			plan.Configuration.ApnsConfigurationModel.FileName = types.StringValue(resp.APNSPushProvider.Configuration.GetFileName())
		}
		if resp.APNSPushProvider.GetLastUpdatedDate() != "" {
			plan.LastUpdatedDate = types.StringValue(resp.APNSPushProvider.GetLastUpdatedDate())
		}
	} else if resp.FCMPushProvider != nil {
		plan.ID = types.StringValue(resp.FCMPushProvider.GetId())
		plan.Name = types.StringValue(resp.FCMPushProvider.GetName())
		plan.ProviderType = types.StringValue(resp.FCMPushProvider.GetProviderType())
		if resp.FCMPushProvider.GetLastUpdatedDate() != "" {
			plan.LastUpdatedDate = types.StringValue(resp.FCMPushProvider.GetLastUpdatedDate())
		}
	}
	return diags
}

func createPushProviderReq(plan pushProviderResourceModel) okta.ListPushProviders200ResponseInner {
	if plan.ProviderType.ValueString() == "APNS" {
		provider := okta.PushProvider{
			Name:         plan.Name.ValueStringPointer(),
			ProviderType: plan.ProviderType.ValueStringPointer(),
		}
		apnsConfig := okta.APNSConfiguration{
			FileName:        plan.Configuration.ApnsConfigurationModel.FileName.ValueStringPointer(),
			KeyId:           plan.Configuration.ApnsConfigurationModel.KeyID.ValueStringPointer(),
			TeamId:          plan.Configuration.ApnsConfigurationModel.TeamID.ValueStringPointer(),
			TokenSigningKey: plan.Configuration.ApnsConfigurationModel.TokenSigningKey.ValueStringPointer(),
		}
		return okta.APNSPushProviderAsListPushProviders200ResponseInner(&okta.APNSPushProvider{
			PushProvider:  provider,
			Configuration: &apnsConfig,
		})
	} else if plan.ProviderType.ValueString() == "FCM" {
		provider := okta.PushProvider{
			Name:         plan.Name.ValueStringPointer(),
			ProviderType: plan.ProviderType.ValueStringPointer(),
		}

		data := map[string]interface{}{
			"type":                        plan.Configuration.FcmConfiguration.ServiceAccountJson.Type.ValueString(),
			"project_id":                  plan.Configuration.FcmConfiguration.ServiceAccountJson.ProjectID.ValueString(),
			"private_key":                 plan.Configuration.FcmConfiguration.ServiceAccountJson.PrivateKey.ValueString(),
			"private_key_id":              plan.Configuration.FcmConfiguration.ServiceAccountJson.PrivateKeyID.ValueString(),
			"client_email":                plan.Configuration.FcmConfiguration.ServiceAccountJson.ClientEmail.ValueString(),
			"client_id":                   plan.Configuration.FcmConfiguration.ServiceAccountJson.ClientId.ValueString(),
			"auth_uri":                    plan.Configuration.FcmConfiguration.ServiceAccountJson.AuthURI.ValueString(),
			"token_uri":                   plan.Configuration.FcmConfiguration.ServiceAccountJson.TokenURI.ValueString(),
			"auth_provider_x509_cert_url": plan.Configuration.FcmConfiguration.ServiceAccountJson.AuthProviderX509CertURL.ValueString(),
			"client_x509_cert_url":        plan.Configuration.FcmConfiguration.ServiceAccountJson.ClientX509CertURL.ValueString(),
		}

		fcmConfig := okta.FCMConfiguration{
			ServiceAccountJson: data,
		}
		return okta.FCMPushProviderAsListPushProviders200ResponseInner(&okta.FCMPushProvider{
			PushProvider:  provider,
			Configuration: &fcmConfig,
		})
	}
	return okta.ListPushProviders200ResponseInner{}
}

// Read refreshes the Terraform state with the latest data
func (r *pushProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state pushProviderResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed push provider value from Okta
	getPushProviderResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().PushProviderAPI.GetPushProvider(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Push Provider",
			"Could not read push provider ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(applyPushProviderToState(getPushProviderResp, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success
func (r *pushProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan pushProviderResourceModel
	var state pushProviderResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the push provider
	pushProvider, _, err := r.OktaIDaaSClient.OktaSDKClientV5().PushProviderAPI.ReplacePushProvider(ctx, state.ID.ValueString()).PushProvider(createPushProviderReq(plan)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Push Provider",
			"Could not update push provider, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(applyPushProviderToState(pushProvider, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success
func (r *pushProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state pushProviderResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing push provider
	_, err := r.OktaIDaaSClient.OktaSDKClientV5().PushProviderAPI.DeletePushProvider(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Push Provider",
			"Could not delete push provider, unexpected error: "+err.Error(),
		)
		return
	}
}
