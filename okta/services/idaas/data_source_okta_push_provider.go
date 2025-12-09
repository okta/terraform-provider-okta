package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &pushProviderDataSource{}

func newPushProviderDataSource() datasource.DataSource {
	return &pushProviderDataSource{}
}

func (d *pushProviderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type pushProviderDataSource struct {
	*config.Config
}

type ServiceAccountJsonDataSource struct {
	ProjectID types.String `tfsdk:"project_id"`
	FileName  types.String `tfsdk:"file_name"`
}

type FcmConfigurationDataSource struct {
	ServiceAccountJSON *ServiceAccountJsonDataSource `tfsdk:"service_account_json"`
}

type ApnsConfigurationDataSource struct {
	KeyID    types.String `tfsdk:"key_id"`
	TeamID   types.String `tfsdk:"team_id"`
	FileName types.String `tfsdk:"file_name"`
}

type ConfigurationDataSource struct {
	FcmConfiguration  *FcmConfigurationDataSource  `tfsdk:"fcm_configuration"`
	ApnsConfiguration *ApnsConfigurationDataSource `tfsdk:"apns_configuration"`
}

type pushProviderDataSourceModel struct {
	ID              types.String             `tfsdk:"id"`
	Name            types.String             `tfsdk:"name"`
	ProviderType    types.String             `tfsdk:"provider_type"`
	Configuration   *ConfigurationDataSource `tfsdk:"configuration"`
	LastUpdatedDate types.String             `tfsdk:"last_updated_date"`
}

func (d *pushProviderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_push_provider"
}

func (d *pushProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the push provider.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The display name of the push provider.",
			},
			"provider_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of push provider. Valid values are `APNS` (Apple Push Notification Service) or `FCM` (Firebase Cloud Messaging).",
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
									"project_id": schema.StringAttribute{
										Computed:    true,
										Description: "The project ID.",
									},
									"file_name": schema.StringAttribute{
										Computed:    true,
										Description: "File name for Admin Console display.",
									},
								},
							},
						},
					},
					"apns_configuration": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"key_id": schema.StringAttribute{
								Computed:    true,
								Description: "10-character Key ID obtained from the Apple developer account. Required for APNS provider type.",
							},
							"team_id": schema.StringAttribute{
								Computed:    true,
								Description: "10-character Team ID used to develop the iOS app. Required for APNS provider type.",
							},
							"file_name": schema.StringAttribute{
								Computed:    true,
								Description: "File name for Admin Console display.",
							},
						},
					},
				},
			},
		},
	}
}

func (d *pushProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data pushProviderDataSourceModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	getPushProviderResp, _, err := d.OktaIDaaSClient.OktaSDKClientV5().PushProviderAPI.GetPushProvider(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Push Provider",
			"Could not read push provider ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	if getPushProviderResp.APNSPushProvider != nil {
		data.ID = types.StringValue(getPushProviderResp.APNSPushProvider.GetId())
		data.Name = types.StringValue(getPushProviderResp.APNSPushProvider.GetName())
		data.ProviderType = types.StringValue(getPushProviderResp.APNSPushProvider.GetProviderType())
		data.LastUpdatedDate = types.StringValue(getPushProviderResp.APNSPushProvider.GetLastUpdatedDate())
		conf := &ConfigurationDataSource{}
		conf.ApnsConfiguration = &ApnsConfigurationDataSource{}
		conf.ApnsConfiguration.KeyID = types.StringValue(getPushProviderResp.APNSPushProvider.Configuration.GetKeyId())
		conf.ApnsConfiguration.TeamID = types.StringValue(getPushProviderResp.APNSPushProvider.Configuration.GetTeamId())
		conf.ApnsConfiguration.FileName = types.StringValue(getPushProviderResp.APNSPushProvider.Configuration.GetFileName())
		data.Configuration = conf
	} else if getPushProviderResp.FCMPushProvider != nil {
		data.ID = types.StringValue(getPushProviderResp.FCMPushProvider.GetId())
		data.Name = types.StringValue(getPushProviderResp.FCMPushProvider.GetName())
		data.ProviderType = types.StringValue(getPushProviderResp.FCMPushProvider.GetProviderType())
		data.LastUpdatedDate = types.StringValue(getPushProviderResp.FCMPushProvider.GetLastUpdatedDate())
		conf := &ConfigurationDataSource{}
		conf.FcmConfiguration = &FcmConfigurationDataSource{}
		conf.FcmConfiguration.ServiceAccountJSON = &ServiceAccountJsonDataSource{}
		conf.FcmConfiguration.ServiceAccountJSON.ProjectID = types.StringValue(getPushProviderResp.FCMPushProvider.Configuration.GetProjectId())
		conf.FcmConfiguration.ServiceAccountJSON.FileName = types.StringValue(getPushProviderResp.FCMPushProvider.Configuration.GetFileName())
		data.Configuration = conf
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
