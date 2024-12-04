package okta

import (
	"context"
	"path"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-sdk-golang/v5/okta"
)

func NewAppSignOnPolicyResource() resource.Resource {
	return &appSignOnPolicyResource{}
}

type appSignOnPolicyResource struct {
	*Config
}

type appSignOnPolicyResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	CatchAll      types.Bool   `tfsdk:"catch_all"`
	DefaultRuleID types.String `tfsdk:"default_rule_id"`
}

func (r *appSignOnPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_signon_policy"
}

func (r *appSignOnPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Policy id",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the policy.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the policy.",
				Required:    true,
			},
			"catch_all": schema.BoolAttribute{
				Description: "Default rules of the policy set to `DENY` or not. If `false`, it is set to `DENY`. Note that this is only apply during creation, so any import or update will not work",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"default_rule_id": schema.StringAttribute{
				Description: "Default rules id of the policy",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *appSignOnPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *appSignOnPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state appSignOnPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accessPolicy, _, err := r.oktaSDKClientV5.PolicyAPI.CreatePolicy(ctx).Policy(buildV5AccessPolicy(state)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create access policy",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(mapAccessPolicyToState(accessPolicy, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODU
	rules, _, err := r.oktaSDKClientV5.PolicyAPI.ListPolicyRules(ctx, accessPolicy.AccessPolicy.GetId()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to get default access policy rule",
			err.Error(),
		)
		return
	}
	if len(rules) != 1 {
		resp.Diagnostics.AddError(
			"find more than one default access policy rule",
			"",
		)
		return
	}
	defaultRuleID := rules[0].GetActualInstance().(OktaPolicyRule).GetId()
	state.DefaultRuleID = types.StringValue(defaultRuleID)

	// TODU
	if !state.CatchAll.ValueBool() {
		_, _, err = r.oktaSDKClientV5.PolicyAPI.ReplacePolicyRule(ctx, accessPolicy.AccessPolicy.GetId(), defaultRuleID).PolicyRule(buildV5AccessPolicyRule()).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to update access policy default rule to DENY",
				err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *appSignOnPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state appSignOnPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accessPolicy, _, err := r.oktaSDKClientV5.PolicyAPI.GetPolicy(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read access policy",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(mapAccessPolicyToState(accessPolicy, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *appSignOnPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state appSignOnPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if frameworkIsClassicOrg(ctx, r.Config) {
		resp.Diagnostics.Append(frameworkResourceOIEOnlyFeatureError(appSignOnPolicy)...)
		return
	}

	// 1. find the default app policy
	// 2. assign default policy to all apps whose authentication policy is the policy about to be deleted
	// 3. delete the policy

	defaultPolicy, err := frameworkFindDefaultAccessPolicy(ctx, r.Config)
	if err != nil {
		resp.Diagnostics.AddError(
			"error finding default access policy: %v",
			err.Error(),
		)
		return
	}

	apps, err := frameworkListApps(ctx, r.Config, nil, defaultPaginationLimit)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to list apps in preparation to delete authentication policy: %v",
			err.Error(),
		)
		return
	}

	for _, a := range apps {
		app := a.GetActualInstance().(OktaApp)
		accessPolicy := app.GetLinks().AccessPolicy.GetHref()
		// ignore apps that don't have an access policy, typically Classic org apps.
		if accessPolicy == "" {
			continue
		}
		// app uses this policy as its access policy, change that back to using the default policy
		if path.Base(accessPolicy) == state.ID.ValueString() {
			// update the app with the default policy, ignore errors
			dp := defaultPolicy.GetActualInstance().(OktaPolicy)
			r.oktaSDKClientV5.ApplicationPoliciesAPI.AssignApplicationPolicy(ctx, app.GetId(), dp.GetId()).Execute()
		}
	}

	_, err = r.oktaSDKClientV5.PolicyAPI.DeletePolicy(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to delete access policy",
			err.Error(),
		)
		return
	}
}

func (r *appSignOnPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if frameworkIsClassicOrg(ctx, r.Config) {
		resp.Diagnostics.Append(frameworkResourceOIEOnlyFeatureError(appSignOnPolicy)...)
		return
	}

	var state appSignOnPolicyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accessPolicy, _, err := r.oktaSDKClientV5.PolicyAPI.ReplacePolicy(ctx, state.ID.ValueString()).Policy(buildV5AccessPolicy(state)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update access policy",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(mapAccessPolicyToState(accessPolicy, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *appSignOnPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, tfpath.Root("id"), req, resp)
}

func buildV5AccessPolicy(model appSignOnPolicyResourceModel) okta.ListPolicies200ResponseInner {
	accessPolicy := &okta.AccessPolicy{}
	accessPolicy.SetType("ACCESS_POLICY")
	accessPolicy.SetName(model.Name.ValueString())
	accessPolicy.SetDescription(model.Description.ValueString())
	return okta.ListPolicies200ResponseInner{AccessPolicy: accessPolicy}
}

func mapAccessPolicyToState(data *okta.ListPolicies200ResponseInner, state *appSignOnPolicyResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if data.AccessPolicy == nil {
		diags.AddError("Empty response", "Access policy")
		return diags
	}
	state.ID = types.StringPointerValue(data.AccessPolicy.Id)
	state.Name = types.StringPointerValue(data.AccessPolicy.Name)
	state.Description = types.StringPointerValue(data.AccessPolicy.Description)
	return diags
}

// TODU
func buildV5AccessPolicyRule() okta.ListPolicyRules200ResponseInner {
	accessPolicyRule := &okta.AccessPolicyRule{}
	return okta.ListPolicyRules200ResponseInner{AccessPolicyRule: accessPolicyRule}
}

// TODU double check crud and rerun the test
