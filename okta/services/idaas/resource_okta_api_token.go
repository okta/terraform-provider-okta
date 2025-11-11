package idaas

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &apiTokenResource{}
	_ resource.ResourceWithConfigure   = &apiTokenResource{}
	_ resource.ResourceWithImportState = &apiTokenResource{}
)

func newAPITokenResource() resource.Resource {
	return &apiTokenResource{}
}

type apiTokenResource struct {
	*config.Config
}

func (r *apiTokenResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *apiTokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_token"
}

func (r *apiTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type IPs struct {
	IP types.String `tfsdk:"ip"`
}

type NetworkModel struct {
	Connection types.String `tfsdk:"connection"`
	Include    types.List   `tfsdk:"include"`
	Exclude    types.List   `tfsdk:"exclude"`
}

type apiTokenResourceModel struct {
	ID         types.String  `tfsdk:"id"`
	Name       types.String  `tfsdk:"name"`
	Network    *NetworkModel `tfsdk:"network"`
	UserID     types.String  `tfsdk:"user_id"`
	Created    types.String  `tfsdk:"created"`
	ClientName types.String  `tfsdk:"client_name"`
}

func (r *apiTokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the API token.",
			},
			"user_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The userId of the user who created the API Token.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the API token.",
			},
			"created": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Timestamp when the API token was created.",
			},
			"client_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the API token client",
			},
		},
		Blocks: map[string]schema.Block{
			"network": schema.SingleNestedBlock{
				Description: "The Network Condition of the API Token.",
				Attributes: map[string]schema.Attribute{
					"connection": schema.StringAttribute{
						Optional:    true,
						Description: "The connection type of the Network Condition.",
					},
					"exclude": schema.ListAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The IP address the excluded zone.",
						ElementType: types.StringType,
					},
					"include": schema.ListAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The IP address the included zone.",
						ElementType: types.StringType,
					},
				},
			},
		},
	}
}

func (r *apiTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddWarning(
		"Create Not Supported",
		"This resource cannot be created via Terraform.",
	)
}

func (r *apiTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data apiTokenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getAPITokenResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApiTokenAPI.GetApiToken(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"error in getting API token",
			err.Error(),
		)
		return
	}
	mapAPITokenToState(ctx, getAPITokenResp, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *apiTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, plan apiTokenResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	updateTokenReq, diags := createTokenUpdate(plan, data.Created.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	upsertAPITokenResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApiTokenAPI.UpsertApiToken(ctx, data.ID.ValueString()).ApiTokenUpdate(updateTokenReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in upserting API token",
			err.Error(),
		)
		return
	}
	mapAPITokenToState(ctx, upsertAPITokenResp, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *apiTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data apiTokenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.OktaIDaaSClient.OktaSDKClientV5().ApiTokenAPI.RevokeApiToken(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"error in revoking API token",
			err.Error(),
		)
		return
	}
}

func createTokenUpdate(data apiTokenResourceModel, created string) (v5okta.ApiTokenUpdate, diag.Diagnostics) {
	var diags diag.Diagnostics
	apiTokenUpdateRequest := v5okta.ApiTokenUpdate{}

	// Set basic fields
	apiTokenUpdateRequest.SetName(data.Name.ValueString())
	apiTokenUpdateRequest.SetClientName(data.ClientName.ValueString())
	apiTokenUpdateRequest.SetUserId(data.UserID.ValueString())

	// Handle Created field with proper null/error checking
	if created != "" {
		parsedTime, err := time.Parse(time.RFC3339, created)
		if err != nil {
			diags.AddError(" Could not parse created time", fmt.Sprintf("created time:%s", created))
		} else {
			apiTokenUpdateRequest.SetCreated(parsedTime)
		}
	}

	// Handle Network configuration
	if data.Network != nil {
		network := v5okta.ApiTokenNetwork{}
		network.SetConnection(data.Network.Connection.ValueString())

		// Handle Include IPs
		var includedZones []string
		if !data.Network.Include.IsNull() && !data.Network.Include.IsUnknown() {
			var incl []types.String
			diags := data.Network.Include.ElementsAs(context.Background(), &incl, false)
			if diags.HasError() {
				return apiTokenUpdateRequest, diags
			}
			for _, v := range incl {
				if !v.IsNull() && !v.IsUnknown() {
					includedZones = append(includedZones, v.ValueString())
				}
			}
		}
		network.SetInclude(includedZones)

		// Handle Exclude IPs
		var excludedZones []string
		if !data.Network.Exclude.IsNull() && !data.Network.Exclude.IsUnknown() {
			var excl []types.String
			diags := data.Network.Exclude.ElementsAs(context.Background(), &excl, false)
			if diags.HasError() {
				return apiTokenUpdateRequest, diags
			}
			for _, v := range excl {
				if !v.IsNull() && !v.IsUnknown() {
					excludedZones = append(excludedZones, v.ValueString())
				}
			}
		}
		network.SetExclude(excludedZones)

		apiTokenUpdateRequest.SetNetwork(network)
	}
	return apiTokenUpdateRequest, nil
}

func mapAPITokenToState(ctx context.Context, resp *v5okta.ApiToken, a *apiTokenResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	a.ID = types.StringValue(resp.GetId())
	a.Name = types.StringValue(resp.GetName())
	a.UserID = types.StringValue(resp.GetUserId())
	a.ClientName = types.StringValue(resp.GetClientName())
	a.Created = types.StringValue(resp.GetCreated().Format(time.RFC3339))
	n := NetworkModel{
		Connection: types.StringValue(resp.Network.GetConnection()),
	}

	inc, diags := types.ListValueFrom(ctx, types.StringType, resp.Network.GetInclude())
	if diags.HasError() {
		return diags
	}
	n.Include = inc

	excl, diags := types.ListValueFrom(ctx, types.StringType, resp.Network.GetExclude())
	if diags.HasError() {
		return diags
	}
	n.Exclude = excl
	a.Network = &n

	return diags
}
