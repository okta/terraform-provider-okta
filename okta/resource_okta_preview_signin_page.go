package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &previewSigninPageResource{}
	_ resource.ResourceWithConfigure   = &previewSigninPageResource{}
	_ resource.ResourceWithImportState = &previewSigninPageResource{}
)

func NewPreviewSigninResource() resource.Resource {
	return &previewSigninPageResource{}
}

type previewSigninPageResource struct {
	*Config
}

func (r *previewSigninPageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_preview_signin_page"
}

func (r *previewSigninPageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	newSchema := resourceSignInSchema
	newSchema.Description = "Manage the preview signin page of a brand"
	resp.Schema = newSchema
}

// Configure adds the provider configured client to the resource.
func (r *previewSigninPageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.Config = p
}

func (r *previewSigninPageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state signinPageModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := buildSignInPageRequest(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to build preview sign in page",
			err.Error(),
		)
		return
	}

	previewSigninPage, _, err := r.oktaSDKClientV3.CustomizationAPI.ReplacePreviewSignInPage(ctx, state.BrandID.ValueString()).SignInPage(reqBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update preview sign in page",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(mapSignInPageToState(ctx, previewSigninPage, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *previewSigninPageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state signinPageModel

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	previewSigninPage, _, err := r.Config.oktaSDKClientV3.CustomizationAPI.GetPreviewSignInPage(ctx, state.BrandID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving preview signin page",
			fmt.Sprintf("Error returned: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(mapSignInPageToState(ctx, previewSigninPage, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *previewSigninPageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state signinPageModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.oktaSDKClientV3.CustomizationAPI.DeletePreviewSignInPage(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to delete preview signin page",
			err.Error(),
		)
		return
	}
}

func (r *previewSigninPageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state signinPageModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := buildSignInPageRequest(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to build preview sign in page",
			err.Error(),
		)
		return
	}

	previewSigninPage, _, err := r.oktaSDKClientV3.CustomizationAPI.ReplacePreviewSignInPage(ctx, state.BrandID.ValueString()).SignInPage(reqBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update preview sign in page",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(mapSignInPageToState(ctx, previewSigninPage, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *previewSigninPageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
