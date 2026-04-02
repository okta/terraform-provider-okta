package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &apiServiceIntegrationDataSource{}

func newAPIServiceIntegrationDataSource() datasource.DataSource {
	return &apiServiceIntegrationDataSource{}
}

type apiServiceIntegrationDataSourceModel struct {
	Id             types.String    `tfsdk:"id"`
	Type           types.String    `tfsdk:"type"`
	Name           types.String    `tfsdk:"name"`
	ConfigGuideUrl types.String    `tfsdk:"config_guide_url"`
	Created        types.String    `tfsdk:"created"`
	CreatedBy      types.String    `tfsdk:"created_by"`
	GrantedScopes  []GrantedScopes `tfsdk:"granted_scopes"`
}

type apiServiceIntegrationDataSource struct {
	*config.Config
}

func (d *apiServiceIntegrationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_service_integration"
}

func (d *apiServiceIntegrationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *apiServiceIntegrationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the API service integration",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the API service integration. This string is an underscore-concatenated, lowercased API service integration name.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the API service integration that corresponds with the type property.",
			},
			"config_guide_url": schema.StringAttribute{
				Computed:    true,
				Description: "The URL to the API service integration configuration guide.",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the API Service Integration instance was created.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user ID of the API Service Integration instance creator.",
			},
		},
		Blocks: map[string]schema.Block{
			"granted_scopes": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"scope": schema.StringAttribute{
							Computed:    true,
							Description: "The scope of the API service integration granted.",
						},
					},
				},
				Description: "The list of Okta management scopes granted to the API Service Integration instance.",
			},
		},
	}
}

func (d *apiServiceIntegrationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data apiServiceIntegrationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getAPIServiceIntegrationResp, _, err := d.OktaIDaaSClient.OktaSDKClientV5().ApiServiceIntegrationsAPI.GetApiServiceIntegrationInstance(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create group owner for group "+data.Type.ValueString()+" for group owner user id: ",
			err.Error(),
		)
		return
	}
	data.Type = types.StringValue(getAPIServiceIntegrationResp.GetType())
	var grantedScopes []GrantedScopes
	for _, grantedScope := range getAPIServiceIntegrationResp.GetGrantedScopes() {
		grantedScopes = append(grantedScopes, GrantedScopes{
			Scope: types.StringValue(grantedScope),
		})
	}
	data.GrantedScopes = grantedScopes
	data.Id = types.StringValue(getAPIServiceIntegrationResp.GetId())
	data.Name = types.StringValue(getAPIServiceIntegrationResp.GetName())
	data.ConfigGuideUrl = types.StringValue(getAPIServiceIntegrationResp.GetConfigGuideUrl())
	data.Created = types.StringValue(getAPIServiceIntegrationResp.GetCreatedAt())
	data.CreatedBy = types.StringValue(getAPIServiceIntegrationResp.GetCreatedBy())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
