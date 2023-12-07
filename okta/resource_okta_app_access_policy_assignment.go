package okta

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/okta/okta-sdk-golang/v3/okta"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &appAccessPolicyAssignmentResource{}
	_ resource.ResourceWithConfigure   = &appAccessPolicyAssignmentResource{}
	_ resource.ResourceWithImportState = &appAccessPolicyAssignmentResource{}
)

func NewAppAccessPolicyAssignmentResource() resource.Resource {
	return &appAccessPolicyAssignmentResource{}
}

type appAccessPolicyAssignmentResource struct {
	*Config
}

type appAccessPolicyAssignmentResourceModel struct {
	ID       types.String `tfsdk:"id"`
	AppID    types.String `tfsdk:"app_id"`
	PolicyID types.String `tfsdk:"policy_id"`
}

func (r *appAccessPolicyAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_access_policy_assignment"
}

func (r *appAccessPolicyAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages assignment of Access Policy to an Application",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the resource. This ID is simply the application ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_id": schema.StringAttribute{
				Description: "The application ID; this value is immutable and can not be updated.",
				Required:    true,
			},
			"policy_id": schema.StringAttribute{
				Description: "The access policy ID.",
				Required:    true,
			},
		},
	}
}

func (r *appAccessPolicyAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *appAccessPolicyAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan appAccessPolicyAssignmentResourceModel

	// read TF plan data into model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// find the app
	appInnerResp, err := r.findAppSDKInnerResponse(ctx, plan.AppID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("create failed to find app %q for policy assignment", plan.AppID.ValueString()),
			err.Error(),
		)
		return
	}
	appID, err := concreteAppID(appInnerResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"issue with inner app response",
			err.Error(),
		)
		return
	}

	// assign policy to app
	policyID := plan.PolicyID.ValueString()
	_, err = r.oktaSDKClientV3.ApplicationPoliciesAPI.AssignApplicationPolicy(ctx, appID, policyID).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("couldn't assign policy %q to app %q", policyID, appID),
			err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(appID)
	plan.AppID = types.StringValue(appID)
	plan.PolicyID = types.StringValue(policyID)

	// Save resource representation into the state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *appAccessPolicyAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state appAccessPolicyAssignmentResourceModel

	// read TF plan data into model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// find the app
	appInnerResp, err := r.findAppSDKInnerResponse(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("read failed to find app %q for policy assignment", state.ID.ValueString()),
			err.Error(),
		)
		return
	}

	// policy is in the app links
	appLinks, err := concreteAppLinks(appInnerResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"issue with app links",
			err.Error(),
		)
		return
	}
	_url, err := url.Parse(appLinks.AccessPolicy.Href)
	if err != nil {
		resp.Diagnostics.AddError(
			"issue with access policy URL",
			err.Error(),
		)
		return
	}

	// If the assigned policy changed outside of terraform this will cause
	// change detection to occur
	policyID := strings.TrimPrefix(_url.EscapedPath(), "/api/v1/policies/")
	state.AppID = state.ID // read might be called by import so be sure to explicitly set AppID
	state.PolicyID = types.StringValue(policyID)

	// Set state with fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *appAccessPolicyAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state appAccessPolicyAssignmentResourceModel

	// read TF plan data into model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !plan.AppID.Equal(state.AppID) {
		resp.Diagnostics.AddError(
			"Application ID is immutable",
			"Application ID can not be changed in the configuration once it has been established during create.",
		)
		return
	}

	// find the app
	_, err := r.findAppSDKInnerResponse(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("update failed to find app %q for policy assignment", state.AppID.ValueString()),
			err.Error(),
		)
		return
	}

	if !plan.PolicyID.Equal(state.PolicyID) {
		// policy id has changed in the config, update
		appID := plan.AppID.ValueString()
		_, err = r.oktaSDKClientV3.ApplicationPoliciesAPI.AssignApplicationPolicy(ctx, appID, plan.PolicyID.ValueString()).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("couldn't re-assign policy %q to app %q", plan.PolicyID.ValueString(), appID),
				err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *appAccessPolicyAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Warn(ctx, "True delete for an okta_app_access_policy_assignment is a no-op as this resource will not delete an app or a policy. Additionally there is not an API endpoint to remove an access policy from an app, only update.")
}

func (r *appAccessPolicyAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *appAccessPolicyAssignmentResource) findAppSDKInnerResponse(ctx context.Context, appID string) (*okta.ListApplications200ResponseInner, error) {
	appInnerResp, _, err := r.oktaSDKClientV3.ApplicationAPI.GetApplication(ctx, appID).Execute()
	return appInnerResp, err
}

func concreteAppID(src *okta.ListApplications200ResponseInner) (id string, err error) {
	if src == nil {
		return "", errors.New("list application inner response is nil")
	}
	app := src.GetActualInstance()
	if app == nil {
		return "", errors.New("okta list applications response does not contain a concrete app")
	}

	switch v := app.(type) {
	case *okta.AutoLoginApplication:
		id = v.GetId()
	case *okta.BasicAuthApplication:
		id = v.GetId()
	case *okta.BookmarkApplication:
		id = v.GetId()
	case *okta.BrowserPluginApplication:
		id = v.GetId()
	case *okta.OpenIdConnectApplication:
		id = v.GetId()
	case *okta.SamlApplication:
		id = v.GetId()
	case *okta.SecurePasswordStoreApplication:
		id = v.GetId()
	case *okta.WsFederationApplication:
		id = v.GetId()
	}

	if id == "" {
		err = fmt.Errorf("list application inner response does not contain a concrete app type %T", src)
	}

	return
}

func concreteAppLinks(src *okta.ListApplications200ResponseInner) (links *okta.ApplicationLinks, err error) {
	if src == nil {
		return nil, errors.New("list application inner response is nil")
	}
	app := src.GetActualInstance()
	if app == nil {
		return nil, errors.New("okta list applications response does not contain a concrete app")
	}

	switch v := app.(type) {
	case *okta.AutoLoginApplication:
		links, _ = v.GetLinksOk()
	case *okta.BasicAuthApplication:
		links, _ = v.GetLinksOk()
	case *okta.BookmarkApplication:
		links, _ = v.GetLinksOk()
	case *okta.BrowserPluginApplication:
		links, _ = v.GetLinksOk()
	case *okta.OpenIdConnectApplication:
		links, _ = v.GetLinksOk()
	case *okta.SamlApplication:
		links, _ = v.GetLinksOk()
	case *okta.SecurePasswordStoreApplication:
		links, _ = v.GetLinksOk()
	case *okta.WsFederationApplication:
		links, _ = v.GetLinksOk()
	}

	if links == nil {
		err = fmt.Errorf("list application inner response does not contain a concrete app type %T", src)
	}

	return
}
