package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &customizedSigninPageResource{}
	_ resource.ResourceWithConfigure   = &customizedSigninPageResource{}
	_ resource.ResourceWithImportState = &customizedSigninPageResource{}
)

func newCustomizedSigninResource() resource.Resource {
	return &customizedSigninPageResource{}
}

type customizedSigninPageResource struct {
	*config.Config
}

func (r *customizedSigninPageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customized_signin_page"
}

func (r *customizedSigninPageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	oktaMutexKV.Lock(resources.OktaIDaaSCustomizedSignInPage)
	defer oktaMutexKV.Unlock(resources.OktaIDaaSCustomizedSignInPage)

	newSchema := resourceSignInSchema
	pageContentAttribute := newSchema.Attributes["page_content"].(resourceSchema.StringAttribute)
	pageContentAttribute.Required = false
	pageContentAttribute.Optional = true
	newSchema.Attributes["page_content"] = pageContentAttribute
	widgetVersionAttribute := newSchema.Attributes["widget_version"].(resourceSchema.StringAttribute)
	widgetVersionAttribute.Required = false
	widgetVersionAttribute.Optional = true
	newSchema.Attributes["widget_version"] = widgetVersionAttribute
	newSchema.Description = "Manage the customized signin page of a brand"
	resp.Schema = newSchema
}

// Configure adds the provider configured client to the resource.
func (r *customizedSigninPageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *customizedSigninPageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	oktaMutexKV.Lock(resources.OktaIDaaSCustomizedSignInPage)
	defer oktaMutexKV.Unlock(resources.OktaIDaaSCustomizedSignInPage)

	var state signinPageModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, diags := buildSignInPageRequest(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	customizedSigninPage, _, err := r.OktaIDaaSClient.OktaSDKClientV3().CustomizationAPI.ReplaceCustomizedSignInPage(ctx, state.BrandID.ValueString()).SignInPage(reqBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update customized sign in page",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(mapSignInPageToState(ctx, customizedSigninPage, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *customizedSigninPageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state signinPageModel

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var brandID string
	if state.BrandID.ValueString() != "" {
		brandID = state.BrandID.ValueString()
	} else {
		brandID = state.ID.ValueString()
	}

	customizedSigninPage, _, err := r.OktaIDaaSClient.OktaSDKClientV3().CustomizationAPI.GetCustomizedSignInPage(ctx, brandID).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving customized signin page",
			fmt.Sprintf("Error returned: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(mapSignInPageToState(ctx, customizedSigninPage, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *customizedSigninPageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	oktaMutexKV.Lock(resources.OktaIDaaSCustomizedSignInPage)
	defer oktaMutexKV.Unlock(resources.OktaIDaaSCustomizedSignInPage)

	var state signinPageModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.OktaIDaaSClient.OktaSDKClientV3().CustomizationAPI.DeleteCustomizedSignInPage(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to delete customized signin page",
			err.Error(),
		)
		return
	}
}

func (r *customizedSigninPageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	oktaMutexKV.Lock(resources.OktaIDaaSCustomizedSignInPage)
	defer oktaMutexKV.Unlock(resources.OktaIDaaSCustomizedSignInPage)

	var state signinPageModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, diags := buildSignInPageRequest(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	customizedSigninPage, _, err := r.OktaIDaaSClient.OktaSDKClientV3().CustomizationAPI.ReplaceCustomizedSignInPage(ctx, state.BrandID.ValueString()).SignInPage(reqBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update customized sign in page",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(mapSignInPageToState(ctx, customizedSigninPage, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *customizedSigninPageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("brand_id"), req.ID)...)
}
