package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &sessionViolationPolicyDataSource{}

func newSessionViolationPolicyDataSource() datasource.DataSource {
	return &sessionViolationPolicyDataSource{}
}

func (d *sessionViolationPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type sessionViolationPolicyDataSource struct {
	*config.Config
}

type sessionViolationPolicyDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Status   types.String `tfsdk:"status"`
	RuleID   types.String `tfsdk:"rule_id"`
	Priority types.Int64  `tfsdk:"priority"`
	System   types.Bool   `tfsdk:"system"`
	Created  types.String `tfsdk:"created"`
}

func (d *sessionViolationPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session_violation_policy"
}

func (d *sessionViolationPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the Session Violation Detection Policy. This is a system policy that is automatically created when the Session Violation Detection feature is enabled. There is exactly one Session Violation Detection Policy per organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the Session Violation Detection Policy.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the Session Violation Detection Policy.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Status of the policy: ACTIVE or INACTIVE.",
			},
			"rule_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the modifiable policy rule (non-default). Use this for importing the policy rule resource.",
			},
			"priority": schema.Int64Attribute{
				Computed:    true,
				Description: "Priority of the policy.",
			},
			"system": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this is a system-managed policy created by Okta.",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the policy was created.",
			},
		},
	}
}

func (d *sessionViolationPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data sessionViolationPolicyDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	d.Logger.Info("reading session violation detection policy")

	policies, _, err := d.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.ListPolicies(ctx).Type_("SESSION_VIOLATION_DETECTION").Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Session Violation Detection Policy",
			"Failed to list session violation detection policies: "+err.Error(),
		)
		return
	}

	if len(policies) == 0 {
		resp.Diagnostics.AddError(
			"No Session Violation Detection Policy found",
			"Ensure the Session Violation Detection feature is enabled in your organization.",
		)
		return
	}

	// There should be exactly one Session Violation Detection Policy
	policy := policies[0]

	if policy.SessionViolationDetectionPolicy == nil {
		resp.Diagnostics.AddError(
			"Unexpected policy type",
			"Expected SessionViolationDetectionPolicy but got a different policy type.",
		)
		return
	}

	sessionViolationPolicy := policy.SessionViolationDetectionPolicy

	if sessionViolationPolicy.Id == nil {
		resp.Diagnostics.AddError(
			"Policy ID is nil",
			"The Session Violation Detection Policy ID is unexpectedly nil.",
		)
		return
	}

	data.ID = types.StringPointerValue(sessionViolationPolicy.Id)
	data.Name = types.StringValue(sessionViolationPolicy.Name)
	if sessionViolationPolicy.Status != nil {
		data.Status = types.StringValue(string(*sessionViolationPolicy.Status))
	}
	if sessionViolationPolicy.Priority != nil {
		data.Priority = types.Int64Value(int64(*sessionViolationPolicy.Priority))
	}
	if sessionViolationPolicy.System != nil {
		data.System = types.BoolPointerValue(sessionViolationPolicy.System)
	}
	if sessionViolationPolicy.Created != nil {
		data.Created = types.StringValue(sessionViolationPolicy.Created.String())
	}

	// Fetch the modifiable rule ID (non-default, priority != 99)
	policyId := *sessionViolationPolicy.Id
	rules, _, err := d.OktaIDaaSClient.OktaSDKClientV6().PolicyAPI.ListPolicyRules(ctx, policyId).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading policy rules",
			"Failed to list session violation detection policy rules: "+err.Error(),
		)
		return
	}

	for _, rule := range rules {
		if rule.SessionViolationDetectionPolicyRule != nil {
			// Skip the default rule (priority 99)
			if rule.SessionViolationDetectionPolicyRule.Priority.IsSet() {
				priority := rule.SessionViolationDetectionPolicyRule.Priority.Get()
				if priority != nil && *priority == 99 {
					continue
				}
			}
			// Found the modifiable rule
			if rule.SessionViolationDetectionPolicyRule.Id != nil {
				data.RuleID = types.StringPointerValue(rule.SessionViolationDetectionPolicyRule.Id)
				break
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
