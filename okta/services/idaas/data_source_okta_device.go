package idaas

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &deviceDataSource{}

func newDeviceDataSource() datasource.DataSource {
	return &deviceDataSource{}
}

type deviceDataSource struct {
	*config.Config
}

func (d *deviceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

func (d *deviceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type deviceProfileModel struct {
	DisplayName        types.String `tfsdk:"display_name"`
	Platform           types.String `tfsdk:"platform"`
	Registered         types.Bool   `tfsdk:"registered"`
	DiskEncryptionType types.String `tfsdk:"disk_encryption_type"`
}

type deviceResourceDisplayNameModel struct {
	Sensitive types.Bool   `tfsdk:"sensitive"`
	Value     types.String `tfsdk:"value"`
}

type deviceDataSourceModel struct {
	Id                  types.String                    `tfsdk:"id"`
	Status              types.String                    `tfsdk:"status"`
	ResourceType        types.String                    `tfsdk:"resource_type"`
	Created             types.String                    `tfsdk:"created"`
	LastUpdated         types.String                    `tfsdk:"last_updated"`
	ResourceAlternateId types.String                    `tfsdk:"resource_alternate_id"`
	ResourceId          types.String                    `tfsdk:"resource_id"`
	Profile             *deviceProfileModel             `tfsdk:"profile"`
	ResourceDisplayName *deviceResourceDisplayNameModel `tfsdk:"resource_display_name"`
}

func (d *deviceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the device.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the device.",
			},
			"resource_type": schema.StringAttribute{
				Computed:    true,
				Description: "The resource type of the device.",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "The creation timestamp of the device.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "The last updated timestamp of the device.",
			},
			"resource_alternate_id": schema.StringAttribute{
				Computed:    true,
				Description: "The alternate ID of the device resource.",
			},
			"resource_id": schema.StringAttribute{
				Computed:    true,
				Description: "Alternate key for the id.",
			},
		},
		Blocks: map[string]schema.Block{
			"profile": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"display_name": schema.StringAttribute{
						Computed:    true,
						Description: "The display name of the device.",
					},
					"platform": schema.StringAttribute{
						Computed:    true,
						Description: "The platform of the device.",
					},
					"registered": schema.BoolAttribute{
						Computed:    true,
						Description: "Indicates if the device is registered at Okta.",
					},
					"disk_encryption_type": schema.StringAttribute{
						Computed:    true,
						Description: "The disk encryption type of the device.",
					},
				},
			},
			"resource_display_name": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"sensitive": schema.BoolAttribute{
						Computed:    true,
						Description: "Indicates if the resource display name is sensitive.",
					},
					"value": schema.StringAttribute{
						Computed:    true,
						Description: "The value of the resource display name.",
					},
				},
			},
		},
	}
}

func (d *deviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data deviceDataSourceModel

	// Read the state from the Terraform configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getDeviceResp, _, err := d.OktaIDaaSClient.OktaSDKClientV5().DeviceAPI.GetDevice(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading device",
			"Could not read device with Id "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue(getDeviceResp.GetId())
	data.ResourceType = types.StringValue(getDeviceResp.GetResourceType())
	data.Status = types.StringValue(getDeviceResp.GetStatus())
	data.Created = types.StringValue(getDeviceResp.GetCreated().Format(time.RFC3339))
	data.LastUpdated = types.StringValue(getDeviceResp.GetLastUpdated().Format(time.RFC3339))
	data.ResourceAlternateId = types.StringValue(getDeviceResp.GetResourceAlternateId())
	data.ResourceId = types.StringValue(getDeviceResp.GetResourceId())
	data.Profile = &deviceProfileModel{}
	profile := getDeviceResp.GetProfile()
	data.Profile.DisplayName = types.StringValue(profile.GetDisplayName())
	data.Profile.Platform = types.StringValue(profile.GetPlatform())
	data.Profile.Registered = types.BoolValue(profile.GetRegistered())
	data.Profile.DiskEncryptionType = types.StringValue(profile.GetDiskEncryptionType())
	data.ResourceDisplayName = &deviceResourceDisplayNameModel{}
	resourceDisplayName := getDeviceResp.GetResourceDisplayName()
	data.ResourceDisplayName.Sensitive = types.BoolValue(resourceDisplayName.GetSensitive())
	data.ResourceDisplayName.Value = types.StringValue(resourceDisplayName.GetValue())

	// Save updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
