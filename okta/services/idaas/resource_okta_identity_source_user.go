package idaas

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	frameworkPath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	okta "github.com/okta/okta-sdk-golang/v6/okta"

	"github.com/okta/terraform-provider-okta/okta/config"
)

// Ensure interface compliance
var (
	_ resource.Resource                = &identitySourceUserResource{}
	_ resource.ResourceWithConfigure   = &identitySourceUserResource{}
	_ resource.ResourceWithImportState = &identitySourceUserResource{}
)

// IdentitySourceUserResource defines the resource implementation.
type identitySourceUserResource struct {
	Config *config.Config
}

// identitySourceUserModel describes the resource data model.
type identitySourceUserModel struct {
	ID               types.String                         `tfsdk:"id"`
	IdentitySourceId types.String                         `tfsdk:"identity_source_id"`
	Created          types.String                         `tfsdk:"created"`
	LastUpdated      types.String                         `tfsdk:"last_updated"`
	Profile          *IdentitySourceUserModelProfileModel `tfsdk:"profile"`
}

// IdentitySourceUserModelProfileModel is the nested model for profile.
type IdentitySourceUserModelProfileModel struct {
	Email       types.String `tfsdk:"email"`
	FirstName   types.String `tfsdk:"first_name"`
	HomeAddress types.String `tfsdk:"home_address"`
	LastName    types.String `tfsdk:"last_name"`
	MobilePhone types.String `tfsdk:"mobile_phone"`
	SecondEmail types.String `tfsdk:"second_email"`
	UserName    types.String `tfsdk:"user_name"`
}

func newIdentitySourceUserResource() resource.Resource {
	return &identitySourceUserResource{}
}

func (r *identitySourceUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_source_user"
}

func (r *identitySourceUserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *identitySourceUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Okta Identity Source User.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The external ID of the user in the identity source. Used as the resource identifier.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"identity_source_id": schema.StringAttribute{
				Description: "ID of the identity source",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created": schema.StringAttribute{
				Description: "The timestamp when the user was created in the identity source",
				Computed:    true,
			},
			"last_updated": schema.StringAttribute{
				Description: "The timestamp when the user was last updated in the identity source",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"profile": schema.SingleNestedBlock{
				Description: "Profile",
				Attributes: map[string]schema.Attribute{
					"email": schema.StringAttribute{
						Description: "Email address of the user",
						Optional:    true,
					},
					"first_name": schema.StringAttribute{
						Description: "First name of the user",
						Optional:    true,
					},
					"home_address": schema.StringAttribute{
						Description: "Home address of the user",
						Optional:    true,
					},
					"last_name": schema.StringAttribute{
						Description: "Last name of the user",
						Optional:    true,
					},
					"mobile_phone": schema.StringAttribute{
						Description: "Mobile phone number of the user",
						Optional:    true,
					},
					"second_email": schema.StringAttribute{
						Description: "Alternative email address of the user",
						Optional:    true,
					},
					"user_name": schema.StringAttribute{
						Description: "Username of the user",
						Optional:    true,
					},
				},
			},
		},
	}
}

func (r *identitySourceUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: {identity_source_id}/{id}
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format: {identity_source_id}/{id}",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, frameworkPath.Root("identity_source_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, frameworkPath.Root("id"), parts[1])...)
}

func (r *identitySourceUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state identitySourceUserModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	externalId := state.ID.ValueString()
	identitySourceId := state.IdentitySourceId.ValueString()

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	result, httpResp, err := client.IdentitySourceAPI.GetIdentitySourceUser(ctx, identitySourceId, externalId).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading identity_source_user", err.Error())
		return
	}
	state.Created = types.StringValue(result.GetCreated().Format(time.RFC3339))
	state.LastUpdated = types.StringValue(result.GetLastUpdated().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *identitySourceUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan identitySourceUserModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	identitySourceId := plan.IdentitySourceId.ValueString()
	externalId := plan.ID.ValueString()

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()

	createReq := client.IdentitySourceAPI.CreateIdentitySourceUser(ctx, identitySourceId)
	body := okta.NewUserRequestSchemaWithDefaults()
	body.SetExternalId(externalId)
	if plan.Profile != nil {
		nestedProfile := okta.NewIdentitySourceUserProfileForUpsertRequiredWithDefaults()
		if !plan.Profile.Email.IsNull() && !plan.Profile.Email.IsUnknown() {
			nestedProfile.SetEmail(plan.Profile.Email.ValueString())
		}
		if !plan.Profile.FirstName.IsNull() && !plan.Profile.FirstName.IsUnknown() {
			nestedProfile.SetFirstName(plan.Profile.FirstName.ValueString())
		}
		if !plan.Profile.HomeAddress.IsNull() && !plan.Profile.HomeAddress.IsUnknown() {
			nestedProfile.SetHomeAddress(plan.Profile.HomeAddress.ValueString())
		}
		if !plan.Profile.LastName.IsNull() && !plan.Profile.LastName.IsUnknown() {
			nestedProfile.SetLastName(plan.Profile.LastName.ValueString())
		}
		if !plan.Profile.MobilePhone.IsNull() && !plan.Profile.MobilePhone.IsUnknown() {
			nestedProfile.SetMobilePhone(plan.Profile.MobilePhone.ValueString())
		}
		if !plan.Profile.SecondEmail.IsNull() && !plan.Profile.SecondEmail.IsUnknown() {
			nestedProfile.SetSecondEmail(plan.Profile.SecondEmail.ValueString())
		}
		if !plan.Profile.UserName.IsNull() && !plan.Profile.UserName.IsUnknown() {
			nestedProfile.SetUserName(plan.Profile.UserName.ValueString())
		}
		body.SetProfile(*nestedProfile)
	}
	createReq = createReq.UserRequestSchema(*body)
	_, err := createReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error creating identity_source_user", err.Error())
		return
	}
	// Fetch computed fields
	//via a follow-up Read (Create returned no body).
	readResult, _, readErr := client.IdentitySourceAPI.GetIdentitySourceUser(ctx, identitySourceId, externalId).Execute()
	if readErr == nil {
		plan.Created = types.StringValue(readResult.GetCreated().Format(time.RFC3339))
		plan.LastUpdated = types.StringValue(readResult.GetLastUpdated().Format(time.RFC3339))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *identitySourceUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan identitySourceUserModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state identitySourceUserModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	externalId := state.ID.ValueString()
	identitySourceId := state.IdentitySourceId.ValueString()

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()

	updateReq := client.IdentitySourceAPI.ReplaceExistingIdentitySourceUser(ctx, identitySourceId, externalId)
	updateBody := okta.NewUserRequestSchemaWithDefaults()
	updateBody.SetExternalId(externalId)
	if plan.Profile != nil {
		nestedProfile := okta.NewIdentitySourceUserProfileForUpsertRequiredWithDefaults()
		if !plan.Profile.Email.IsNull() && !plan.Profile.Email.IsUnknown() {
			nestedProfile.SetEmail(plan.Profile.Email.ValueString())
		}
		if !plan.Profile.FirstName.IsNull() && !plan.Profile.FirstName.IsUnknown() {
			nestedProfile.SetFirstName(plan.Profile.FirstName.ValueString())
		}
		if !plan.Profile.HomeAddress.IsNull() && !plan.Profile.HomeAddress.IsUnknown() {
			nestedProfile.SetHomeAddress(plan.Profile.HomeAddress.ValueString())
		}
		if !plan.Profile.LastName.IsNull() && !plan.Profile.LastName.IsUnknown() {
			nestedProfile.SetLastName(plan.Profile.LastName.ValueString())
		}
		if !plan.Profile.MobilePhone.IsNull() && !plan.Profile.MobilePhone.IsUnknown() {
			nestedProfile.SetMobilePhone(plan.Profile.MobilePhone.ValueString())
		}
		if !plan.Profile.SecondEmail.IsNull() && !plan.Profile.SecondEmail.IsUnknown() {
			nestedProfile.SetSecondEmail(plan.Profile.SecondEmail.ValueString())
		}
		if !plan.Profile.UserName.IsNull() && !plan.Profile.UserName.IsUnknown() {
			nestedProfile.SetUserName(plan.Profile.UserName.ValueString())
		}
		updateBody.SetProfile(*nestedProfile)
	}
	updateReq = updateReq.UserRequestSchema(*updateBody)

	result, _, err := updateReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error updating identity_source_user", err.Error())
		return
	}
	state.Created = types.StringValue(result.GetCreated().Format(time.RFC3339))
	state.LastUpdated = types.StringValue(result.GetLastUpdated().Format(time.RFC3339))
	state.Profile = plan.Profile

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *identitySourceUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state identitySourceUserModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	externalId := state.ID.ValueString()
	identitySourceId := state.IdentitySourceId.ValueString()

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	httpResp, err := client.IdentitySourceAPI.DeleteIdentitySourceUser(ctx, identitySourceId, externalId).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Error deleting identity_source_user", err.Error())
		return
	}
}

// Ensure diag is used
var _ diag.Diagnostics
