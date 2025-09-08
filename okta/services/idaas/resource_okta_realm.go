package idaas

import (
	"context"
	"io"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

type realmModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	RealmType types.String `tfsdk:"realm_type"`
	IsDefault types.Bool   `tfsdk:"is_default"`
}

type realmResource struct {
	config *config.Config
}

func newRealmResource() resource.Resource {
	return &realmResource{}
}

func (r *realmResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_realm"
}

func (r *realmResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Realm ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Okta Realm.",
			},
			"realm_type": schema.StringAttribute{
				Required:    true,
				Description: "The realm type. Valid values: `PARTNER` and `DEFAULT`",
				Validators: []validator.String{
					stringvalidator.OneOf("PARTNER", "DEFAULT"),
				},
			},
			"is_default": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether the realm is the default realm.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Description: "Creates an Okta Realm. This resource allows you to create and configure an Okta Realm.",
	}
}

func (r *realmResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.config = resourceConfiguration(req, resp)
}

func (r *realmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state realmModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRealmRequest := v5okta.NewCreateRealmRequest()
	profile := v5okta.NewRealmProfile(state.Name.ValueString())
	profile.RealmType = state.RealmType.ValueStringPointer()
	createRealmRequest.Profile = profile

	responseRealm, response, err := r.config.OktaIDaaSClient.OktaSDKClientV5().RealmAPI.CreateRealm(ctx).Body(*createRealmRequest).Execute()
	if err != nil {
		body, ioErr := io.ReadAll(response.Body)
		defer response.Body.Close()
		if ioErr != nil {
			resp.Diagnostics.AddError(err.Error(), "failed to read response body")
			return
		}
		resp.Diagnostics.AddError("failed to create realm:"+err.Error(), string(body))
		return
	}

	resp.Diagnostics.Append(mapRealmResourceToState(responseRealm, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *realmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state realmModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	realm, response, err := r.config.OktaIDaaSClient.OktaSDKClientV5().RealmAPI.GetRealm(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		body, ioErr := io.ReadAll(response.Body)
		defer response.Body.Close()
		if ioErr != nil {
			resp.Diagnostics.AddError(err.Error(), "failed to read response body")
			return
		}
		resp.Diagnostics.AddError("failed to read realm:"+err.Error(), string(body))
		return
	}

	if realm == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(mapRealmResourceToState(realm, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *realmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state realmModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRealmRequest := v5okta.NewUpdateRealmRequest()
	profile := v5okta.NewRealmProfile(state.Name.ValueString())
	profile.RealmType = state.RealmType.ValueStringPointer()
	updateRealmRequest.Profile = profile

	realm, response, err := r.config.OktaIDaaSClient.OktaSDKClientV5().RealmAPI.ReplaceRealm(ctx, state.ID.ValueString()).Body(*updateRealmRequest).Execute()
	if err != nil {
		body, ioErr := io.ReadAll(response.Body)
		defer response.Body.Close()
		if ioErr != nil {
			resp.Diagnostics.AddError(err.Error(), "failed to read response body")
			return
		}
		resp.Diagnostics.AddError("failed to update realm:"+err.Error(), string(body))
		return
	}

	resp.Diagnostics.Append(mapRealmResourceToState(realm, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *realmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state realmModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.config.OktaIDaaSClient.OktaSDKClientV5().RealmAPI.DeleteRealm(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		body, ioErr := io.ReadAll(response.Body)
		defer response.Body.Close()
		if ioErr != nil {
			resp.Diagnostics.AddError(err.Error(), "failed to read response body")
			return
		}
		resp.Diagnostics.AddError("failed to delete realm:"+err.Error(), string(body))
		return
	}
}

func (r *realmResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapRealmResourceToState(realmResource *v5okta.Realm, state *realmModel) diag.Diagnostics {
	var diags diag.Diagnostics

	state.ID = types.StringPointerValue(realmResource.Id)
	state.Name = types.StringValue(realmResource.Profile.Name)
	state.RealmType = types.StringPointerValue(realmResource.Profile.RealmType)
	state.IsDefault = types.BoolPointerValue(realmResource.IsDefault)

	return diags
}
