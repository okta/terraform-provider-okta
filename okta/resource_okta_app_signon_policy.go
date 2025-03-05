package okta

import (
	"context"
	"errors"
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
		Description: ` Manages a sign-on policy.
		
~> **WARNING:** This feature is only available as a part of the Okta Identity
Engine (OIE) and ***is not*** compatible with Classic orgs. Authentication
policies for applications in a Classic org can only be modified in the Admin UI,
there isn't a public API for this. Therefore the Okta Terraform Provider does
not support this resource for Classic orgs. [Contact
support](mailto:dev-inquiries@okta.com) for further information.
This resource allows you to create and configure a sign-on policy for the
application. Inside the product a sign-on policy is referenced as an
_authentication policy_, in the public API the policy is of type
['ACCESS_POLICY'](https://developer.okta.com/docs/reference/api/policy/#policy-object).
A newly created app's sign-on policy will always contain the default
authentication policy unless one is assigned via 'authentication_policy' in the
app resource. At the API level the default policy has system property value of
true.
~> **WARNING:** When this policy is destroyed any other applications that
associate the policy as their authentication policy will be reassigned to the
default/system access policy.`,
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
				Description: "If false, the default rule of the policy is set access to `DENY`. Otherwise default behavior of the default rule is to leave access at `ALLOW`.  **WARNING** setting this attribute to false changes policy rule's default behavior. Use at your own risk. This is only applied during creation and does not affect import or update.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"default_rule_id": schema.StringAttribute{
				Description: "Default rule (system=true) id of the policy",
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
	if frameworkIsClassicOrg(ctx, r.Config) {
		resp.Diagnostics.Append(frameworkResourceOIEOnlyFeatureError(appSignOnPolicy)...)
		return
	}
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

	resp.Diagnostics.Append(r.mapAccessPolicyToState(ctx, accessPolicy, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.CatchAll.ValueBool() {
		defaultRule, err := r.findDefaultPolicyRuleResponse(ctx, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to find default policy rule",
				err.Error(),
			)
			return
		}
		if actions, ok := defaultRule.AccessPolicyRule.GetActionsOk(); ok {
			if _, ok := actions.GetAppSignOnOk(); ok {
				defaultRule.AccessPolicyRule.Actions.AppSignOn.SetAccess("DENY")
			}
		}
		defaultRule.AccessPolicyRule.Actions.AppSignOn.SetAccess("DENY")

		_, _, err = r.oktaSDKClientV5.PolicyAPI.ReplacePolicyRule(ctx, state.ID.ValueString(), state.DefaultRuleID.ValueString()).PolicyRule(*defaultRule).Execute()
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
	if frameworkIsClassicOrg(ctx, r.Config) {
		resp.Diagnostics.Append(frameworkResourceOIEOnlyFeatureError(appSignOnPolicy)...)
		return
	}
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

	resp.Diagnostics.Append(r.mapAccessPolicyToState(ctx, accessPolicy, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *appSignOnPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if frameworkIsClassicOrg(ctx, r.Config) {
		resp.Diagnostics.Append(frameworkResourceOIEOnlyFeatureError(appSignOnPolicy)...)
		return
	}

	var state appSignOnPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
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
		if a.GetActualInstance() == nil {
			// nil bumper
			continue
		}
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

	resp.Diagnostics.Append(r.mapAccessPolicyToState(ctx, accessPolicy, &state)...)
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

func (r *appSignOnPolicyResource) mapAccessPolicyToState(ctx context.Context, data *okta.ListPolicies200ResponseInner, state *appSignOnPolicyResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if data.AccessPolicy == nil {
		diags.AddError("Empty response", "Access policy")
		return diags
	}
	state.ID = types.StringPointerValue(data.AccessPolicy.Id)
	state.Name = types.StringPointerValue(data.AccessPolicy.Name)
	state.Description = types.StringPointerValue(data.AccessPolicy.Description)

	defaultRule, err := r.findDefaultPolicyRuleResponse(ctx, state.ID.ValueString())
	if err != nil {
		diags.AddError(
			"failed to get default access policy rule",
			err.Error(),
		)
		return diags
	}
	defaultRuleID := defaultRule.AccessPolicyRule.GetId()
	state.DefaultRuleID = types.StringValue(defaultRuleID)

	return diags
}

// findDefaultPolicyRuleResponse find the default policy rule from the list and return
// it. Default rule is the first to be marked system.
func (r *appSignOnPolicyResource) findDefaultPolicyRuleResponse(ctx context.Context, accessPolicyId string) (*okta.ListPolicyRules200ResponseInner, error) {
	rules, _, err := r.oktaSDKClientV5.PolicyAPI.ListPolicyRules(ctx, accessPolicyId).Execute()
	if err != nil {
		return nil, err
	}
	for _, rule := range rules {
		if rule.AccessPolicyRule == nil {
			continue
		}
		if system, ok := rule.AccessPolicyRule.GetSystemOk(); ok {
			if *system {
				return &rule, nil
			}
		}
	}
	return nil, errors.New("policy does not have a default (system) access policy rule")
}
