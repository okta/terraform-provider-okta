package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &apiServiceIntegrationResource{}
	_ resource.ResourceWithConfigure   = &apiServiceIntegrationResource{}
	_ resource.ResourceWithImportState = &apiServiceIntegrationResource{}
)

func newAPIServiceIntegrationResource() resource.Resource {
	return &apiServiceIntegrationResource{}
}

type apiServiceIntegrationResource struct {
	*config.Config
}

type GrantedScopes struct {
	Scope types.String `tfsdk:"scope"`
}

type apiServiceIntegrationResourceModel struct {
	Id            types.String    `tfsdk:"id"`
	Type          types.String    `tfsdk:"type"`
	Name          types.String    `tfsdk:"name"`
	GrantedScopes []GrantedScopes `tfsdk:"granted_scopes"`
}

func (r *apiServiceIntegrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_service_integration"
}

func (r *apiServiceIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *apiServiceIntegrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the API service integration",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the API service integration. This string is an underscore-concatenated, lowercased API service integration name.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the API service integration that corresponds with the type property.",
			},
		},
		Blocks: map[string]schema.Block{
			"granted_scopes": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"scope": schema.StringAttribute{
							Required:    true,
							Description: "The scope of the API service integration",
						},
					},
				},
				Description: "The list of Okta management scopes granted to the API Service Integration instance.",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *apiServiceIntegrationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *apiServiceIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data apiServiceIntegrationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiServiceReq := buildCreateAPIServiceIntegrationRequest(data)

	apiServiceIntegrationResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApiServiceIntegrationsAPI.CreateApiServiceIntegrationInstance(ctx).PostAPIServiceIntegrationInstanceRequest(apiServiceReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create api service integration",
			err.Error(),
		)
		return
	}

	mapAPIServiceIntegrationToState(apiServiceIntegrationResp, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *apiServiceIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data apiServiceIntegrationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getAPIServiceIntegrationResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApiServiceIntegrationsAPI.GetApiServiceIntegrationInstance(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create group owner for group "+data.Type.ValueString()+" for group owner user id: ",
			err.Error(),
		)
		return
	}
	mapAPIServiceIntegrationToState2(getAPIServiceIntegrationResp, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *apiServiceIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update Not Supported",
		"This resource cannot be updated via Terraform.",
	)
}

func (r *apiServiceIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data apiServiceIntegrationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.OktaIDaaSClient.OktaSDKClientV5().ApiServiceIntegrationsAPI.DeleteApiServiceIntegrationInstance(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to delete API Service Integration",
			err.Error(),
		)
		return
	}

}

func buildCreateAPIServiceIntegrationRequest(data apiServiceIntegrationResourceModel) v5okta.PostAPIServiceIntegrationInstanceRequest {
	var grantedScopes []string
	for _, grantedScope := range data.GrantedScopes {
		grantedScopes = append(grantedScopes, grantedScope.Scope.ValueString())
	}
	return v5okta.PostAPIServiceIntegrationInstanceRequest{
		Type:          data.Type.ValueString(),
		GrantedScopes: grantedScopes,
	}
}

func mapAPIServiceIntegrationToState(resp *v5okta.PostAPIServiceIntegrationInstance, a *apiServiceIntegrationResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	a.Type = types.StringValue(resp.GetType())
	var grantedScopes []GrantedScopes
	for _, grantedScope := range resp.GetGrantedScopes() {
		grantedScopes = append(grantedScopes, GrantedScopes{
			Scope: types.StringValue(grantedScope),
		})
	}
	a.GrantedScopes = grantedScopes
	a.Id = types.StringValue(resp.GetId())
	a.Name = types.StringValue(resp.GetName())
	return diags
}

func mapAPIServiceIntegrationToState2(resp *v5okta.APIServiceIntegrationInstance, a *apiServiceIntegrationResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	a.Type = types.StringValue(resp.GetType())
	var grantedScopes []GrantedScopes
	for _, grantedScope := range resp.GetGrantedScopes() {
		grantedScopes = append(grantedScopes, GrantedScopes{
			Scope: types.StringValue(grantedScope),
		})
	}
	a.GrantedScopes = grantedScopes
	a.Id = types.StringValue(resp.GetId())
	a.Name = types.StringValue(resp.GetName())
	return diags
}
