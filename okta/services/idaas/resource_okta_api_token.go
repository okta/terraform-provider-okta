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
	Include    []IPs        `tfsdk:"include"`
	Exclude    []IPs        `tfsdk:"exclude"`
}

type apiTokenResourceModel struct {
	Id         types.String  `tfsdk:"id"`
	Name       types.String  `tfsdk:"name"`
	Network    *NetworkModel `tfsdk:"network"`
	UserId     types.String  `tfsdk:"user_id"`
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
				Description: "The ID of the user that the API token will be used.",
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
				},
				Blocks: map[string]schema.Block{
					"exclude": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"ip": schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: "The IP address the excluded zone.",
								},
							},
						},
					},
					"include": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"ip": schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: "The IP address the included zone.",
								},
							},
						},
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

	getAPITokenResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApiTokenAPI.GetApiToken(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"error in getting API token",
			err.Error(),
		)
		return
	}
	mapAPITokeToState(getAPITokenResp, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *apiTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, plan apiTokenResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	upsertAPITokenResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().ApiTokenAPI.UpsertApiToken(ctx, data.Id.ValueString()).ApiTokenUpdate(createTokenUpdate(plan, data.Created.ValueString())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"error in upserting API token",
			err.Error(),
		)
		return
	}
	mapAPITokeToState(upsertAPITokenResp, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *apiTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data apiTokenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.OktaIDaaSClient.OktaSDKClientV5().ApiTokenAPI.RevokeApiToken(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"error in revoking API token",
			err.Error(),
		)
		return
	}
}

func createTokenUpdate(data apiTokenResourceModel, created string) v5okta.ApiTokenUpdate {
	x := v5okta.ApiTokenUpdate{}

	// Set basic fields
	x.Name = data.Name.ValueStringPointer()
	x.ClientName = data.ClientName.ValueStringPointer()
	x.UserId = data.UserId.ValueStringPointer()

	// Handle Created field with proper null/error checking
	if created != "" {
		parsedTime, err := time.Parse(time.RFC3339, created)
		if err != nil {
			fmt.Printf("Warning: Could not parse created time '%s': %v\n", created, err)
		} else {
			x.Created = &parsedTime
		}
	}

	// Handle Network configuration
	if data.Network != nil {
		network := v5okta.ApiTokenNetwork{}
		network.Connection = data.Network.Connection.ValueStringPointer()
		network.Connection = data.Network.Connection.ValueStringPointer()

		// Handle Include IPs
		for _, inc := range data.Network.Include {
			if !inc.IP.IsNull() && !inc.IP.IsUnknown() {
				network.Include = append(network.Include, inc.IP.ValueString())
			}
		}

		// Handle Exclude IPs
		for _, exc := range data.Network.Exclude {
			if !exc.IP.IsNull() && !exc.IP.IsUnknown() {
				network.Exclude = append(network.Exclude, exc.IP.ValueString())
			}
		}

		x.Network = &network
	}
	return x
}

func mapAPITokeToState(resp *v5okta.ApiToken, a *apiTokenResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	a.Id = types.StringValue(resp.GetId())
	a.Name = types.StringValue(resp.GetName())
	a.UserId = types.StringValue(resp.GetUserId())
	a.ClientName = types.StringValue(resp.GetClientName())
	a.Created = types.StringValue(resp.GetCreated().Format(time.RFC3339))
	n := NetworkModel{
		Connection: types.StringValue(resp.Network.GetConnection()),
		Include:    []IPs{},
		Exclude:    []IPs{},
	}

	if resp.Network != nil {
		for _, inc := range resp.Network.Include {
			n.Include = append(n.Include, IPs{
				IP: types.StringValue(inc),
			})
		}
		for _, exc := range resp.Network.Exclude {
			n.Exclude = append(n.Exclude, IPs{
				IP: types.StringValue(exc),
			})
		}
	}

	a.Network = &n

	return diags
}
