package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-sdk-golang/v4/okta"
)

func NewTrustedServerResource() resource.Resource {
	return &trustedServerResource{}
}

type trustedServerResource struct {
	*Config
}

type trustedServerResourceModel struct {
	ID          types.String `tfsdk:"id"`
	AuthSeverID types.String `tfsdk:"auth_server_id"`
	Trusted     types.Set    `tfsdk:"trusted"`
}

func (r *trustedServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trusted_server"
}

func (r *trustedServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Associated (Trusted) authorization servers allow you to designate a trusted authorization server that you associate with another authorization server.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Resource id",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"auth_server_id": schema.StringAttribute{
				Description: "Authorization server ID",
				Required:    true,
			},
			"trusted": schema.SetAttribute{
				Description: "A list of the authorization server IDs user want to trust",
				ElementType: types.StringType,
				Required:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *trustedServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *trustedServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state trustedServerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	elements := make([]types.String, 0, len(state.Trusted.Elements()))
	resp.Diagnostics.Append(state.Trusted.ElementsAs(ctx, &elements, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	convertedElements := make([]string, 0)
	for _, v := range elements {
		convertedElements = append(convertedElements, v.ValueString())
	}
	listAuthorizationServer, _, err := r.oktaSDKClientV3.AuthorizationServerAssocAPI.CreateAssociatedServers(ctx, state.AuthSeverID.ValueString()).AssociatedServerMediated(okta.AssociatedServerMediated{Trusted: convertedElements}).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create trusted servers",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(mapTrustedServersToState(listAuthorizationServer, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *trustedServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state trustedServerResourceModel

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	listAuthorizationServer, _, err := r.Config.oktaSDKClientV3.AuthorizationServerAssocAPI.ListAssociatedServersByTrustedType(ctx, state.AuthSeverID.ValueString()).Trusted(true).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving list trusted server",
			fmt.Sprintf("Error returned: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(mapTrustedServersToState(listAuthorizationServer, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *trustedServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state trustedServerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	elements := make([]types.String, 0, len(state.Trusted.Elements()))
	resp.Diagnostics.Append(state.Trusted.ElementsAs(ctx, &elements, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	convertedElements := make([]string, 0)
	for _, v := range elements {
		convertedElements = append(convertedElements, v.ValueString())
	}

	for _, trustedServerID := range convertedElements {
		_, err := r.oktaSDKClientV3.AuthorizationServerAssocAPI.DeleteAssociatedServer(ctx, state.AuthSeverID.ValueString(), trustedServerID).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed to delete trusted server %v", trustedServerID),
				err.Error(),
			)
			return
		}
	}
}

func (r *trustedServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan trustedServerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state trustedServerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	oldTrustedIDs := make([]types.String, 0, len(state.Trusted.Elements()))
	resp.Diagnostics.Append(state.Trusted.ElementsAs(ctx, &oldTrustedIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	convertedOldTrustedIDs := make([]string, 0)
	for _, v := range oldTrustedIDs {
		convertedOldTrustedIDs = append(convertedOldTrustedIDs, v.ValueString())
	}

	newTrustedIDs := make([]types.String, 0, len(plan.Trusted.Elements()))
	resp.Diagnostics.Append(plan.Trusted.ElementsAs(ctx, &newTrustedIDs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	convertedNewTrustedIDs := make([]string, 0)
	for _, v := range newTrustedIDs {
		convertedNewTrustedIDs = append(convertedNewTrustedIDs, v.ValueString())
	}

	_, toDelete, toAdd := Intersection(convertedOldTrustedIDs, convertedNewTrustedIDs)
	var err error
	listAuthorizationServer := make([]okta.AuthorizationServer, 0)
	if len(toAdd) > 0 {
		listAuthorizationServer, _, err = r.oktaSDKClientV3.AuthorizationServerAssocAPI.CreateAssociatedServers(ctx, state.AuthSeverID.ValueString()).AssociatedServerMediated(okta.AssociatedServerMediated{Trusted: toAdd}).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to update trusted servers",
				err.Error(),
			)
			return
		}
	}

	for _, trustedServerID := range toDelete {
		_, err := r.oktaSDKClientV3.AuthorizationServerAssocAPI.DeleteAssociatedServer(ctx, state.AuthSeverID.ValueString(), trustedServerID).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed to delete trusted server %v", trustedServerID),
				err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(mapTrustedServersToState(listAuthorizationServer, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func mapTrustedServersToState(data []okta.AuthorizationServer, state *trustedServerResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	state.ID = types.StringValue(state.AuthSeverID.ValueString())
	return diags
}
