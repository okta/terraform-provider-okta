package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/okta/terraform-provider-okta/okta/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &previewSigninPageResource{}
	_ resource.ResourceWithConfigure   = &previewSigninPageResource{}
	_ resource.ResourceWithImportState = &previewSigninPageResource{}
)

func newPreviewSigninResource() resource.Resource {
	return &previewSigninPageResource{}
}

type previewSigninPageResource struct {
	*config.Config
}

func (r *previewSigninPageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_preview_signin_page"
}

func (r *previewSigninPageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	// Clone the shared schema's Attributes map to avoid mutating the global
	// `resourceSignInSchema` which can cause concurrent map read/write panics
	// when multiple resources are initialized concurrently.
	newSchema := resourceSignInSchema
	newAttrs := make(map[string]resourceSchema.Attribute, len(resourceSignInSchema.Attributes))
	for k, v := range resourceSignInSchema.Attributes {
		newAttrs[k] = v
	}
	newSchema.Attributes = newAttrs
	newSchema.Description = "Manage the preview signin page of a brand"
	resp.Schema = newSchema
}

// Configure adds the provider configured client to the resource.
func (r *previewSigninPageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *previewSigninPageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state signinPageModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, diags := buildSignInPageRequest(ctx, state)
	if diags.HasError() {
		return
	}

	previewSigninPage, _, err := r.OktaIDaaSClient.OktaSDKClientV3().CustomizationAPI.ReplacePreviewSignInPage(ctx, state.BrandID.ValueString()).SignInPage(reqBody).Execute()
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

	previewSigninPage, _, err := r.OktaIDaaSClient.OktaSDKClientV3().CustomizationAPI.GetPreviewSignInPage(ctx, state.BrandID.ValueString()).Execute()
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

	_, err := r.OktaIDaaSClient.OktaSDKClientV3().CustomizationAPI.DeletePreviewSignInPage(ctx, state.ID.ValueString()).Execute()
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

	reqBody, diags := buildSignInPageRequest(ctx, state)
	if diags.HasError() {
		return
	}

	previewSigninPage, _, err := r.OktaIDaaSClient.OktaSDKClientV3().CustomizationAPI.ReplacePreviewSignInPage(ctx, state.BrandID.ValueString()).SignInPage(reqBody).Execute()
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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("brand_id"), req.ID)...)
}
