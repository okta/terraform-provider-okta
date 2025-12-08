package governance

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-governance-sdk-golang/governance"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &requestV2Resource{}
	_ resource.ResourceWithConfigure   = &requestV2Resource{}
	_ resource.ResourceWithImportState = &requestV2Resource{}
)

func newRequestV2Resource() resource.Resource {
	return &requestV2Resource{}
}

func (r *requestV2Resource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r *requestV2Resource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

type requestV2Resource struct {
	*config.Config
}

type requested struct {
	EntryId         types.String `tfsdk:"entry_id"`
	Type            types.String `tfsdk:"type"`
	AccessScopeId   types.String `tfsdk:"access_scope_id"`
	AccessScopeType types.String `tfsdk:"access_scope_type"`
	ResourceId      types.String `tfsdk:"resource_id"`
	ResourceType    types.String `tfsdk:"resource_type"`
}

type riskAssessment struct {
	RequestSubmissionType types.String `tfsdk:"request_submission_type"`
	RiskRules             []riskRules  `tfsdk:"risk_rules"`
}

type riskRules struct {
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	ResourceName types.String `tfsdk:"resource_name"`
}

type values struct {
	value types.String `tfsdk:"value"`
}

type requestedFieldValues struct {
	Id     types.String `tfsdk:"id"`
	Label  types.String `tfsdk:"label"`
	Type   types.String `tfsdk:"type"`
	Value  types.String `tfsdk:"value"`
	Values []values     `tfsdk:"values"`
}

type requestV2ResourceModel struct {
	Id                   types.String            `tfsdk:"id"`
	Created              types.String            `tfsdk:"created"`
	CreatedBy            types.String            `tfsdk:"created_by"`
	LastUpdated          types.String            `tfsdk:"last_updated"`
	LastUpdatedBy        types.String            `tfsdk:"last_updated_by"`
	Status               types.String            `tfsdk:"status"`
	AccessDuration       types.String            `tfsdk:"access_duration"`
	Granted              types.String            `tfsdk:"granted"`
	GrantStatus          types.String            `tfsdk:"grant_status"`
	Resolved             types.String            `tfsdk:"resolved"`
	RevocationScheduled  types.String            `tfsdk:"revocation_scheduled"`
	RevocationStatus     types.String            `tfsdk:"revocation_status"`
	Revoked              types.String            `tfsdk:"revoked"`
	Requested            *requested              `tfsdk:"requested"`
	RequestedFor         *entitlementParentModel `tfsdk:"requested_for"`
	RequesterFieldValues []requestedFieldValues  `tfsdk:"requester_field_values"`
}

func (r *requestV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_v2"
}

func (r *requestV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the request.",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "The ISO 8601 formatted date and time when the resource was created.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who created the resource.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "The ISO 8601 formatted date and time when the resource was last updated.",
			},
			"last_updated_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who last updated the resource.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the request.",
			},
			"access_duration": schema.StringAttribute{
				Computed:    true,
				Description: "How long the requester retains access after their request is approved and fulfilled.",
			},
			"granted": schema.StringAttribute{
				Computed:    true,
				Description: "The date the approved access was granted. Only set if request.status is APPROVED.",
			},
			"grant_status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the granted access request.",
			},
			"resolved": schema.StringAttribute{
				Computed:    true,
				Description: "The date the request was resolved.",
			},
			"revocation_scheduled": schema.StringAttribute{
				Computed:    true,
				Description: "The date the granted access is scheduled for revocation.",
			},
			"revocation_status": schema.StringAttribute{
				Computed:    true,
				Description: "The revocation status of the request.",
			},
			"revoked": schema.StringAttribute{
				Computed:    true,
				Description: "The date the granted access was revoked.",
			},
		},
		Blocks: map[string]schema.Block{
			"requested": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"entry_id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource catalog entry.",
					},
					"type": schema.StringAttribute{
						Required:    true,
						Description: "The type of the resource.",
					},
					"access_scope_id": schema.StringAttribute{
						Computed:    true,
						Description: "The ID of the access scope associated with the resource.",
					},
					"access_scope_type": schema.StringAttribute{
						Computed:    true,
						Description: "The access scope type.",
					},
					"resource_id": schema.StringAttribute{
						Computed:    true,
						Description: "The ID of the requested resource.",
					},
					"resource_type": schema.StringAttribute{
						Computed:    true,
						Description: "The requested resource type.",
					},
				},
				Description: "A representation of a resource that can be requested for access.",
			},
			"requested_for": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the Okta user.",
					},
					"type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("OKTA_USER"),
						},
						Description: "The type of principal.",
					},
				},
				Description: "A representation of a principal.",
			},
			"requester_field_values": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Optional:    true,
							Description: "The ID of a requester field.",
						},
						"label": schema.StringAttribute{
							Optional:    true,
							Description: "A human-readable description of requester field.",
						},
						"type": schema.StringAttribute{
							Optional:    true,
							Description: "Type of value for the requester field.",
						},
						"value": schema.StringAttribute{
							Optional:    true,
							Description: "The value of requester field, which depends on the type of the field.",
						},
					},
					Blocks: map[string]schema.Block{
						"values": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"value": schema.StringAttribute{
										Optional: true,
									},
								},
							},
							Description: "The values of requester field with the type MULTISELECT. If the field type is MULTISELECT, this property is required.",
						},
					},
				},
				Description: "The requester input fields required by the approval system.",
			},
		},
	}
}

func (r *requestV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data requestV2ResourceModel

	// Read Terraform plan Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	reqCreatableResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestsAPI.CreateRequestV2(ctx).RequestCreatable2(createRequestReq(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Request V2",
			"Could not create Request V2, unexpected error: "+err.Error(),
		)
		return
	}

	applyRequestResourceToState(&data, reqCreatableResp)

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data requestV2ResourceModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getRequestV2Resp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestsAPI.GetRequestV2(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Request V2",
			"Could not reading Request V2, unexpected error: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue(getRequestV2Resp.GetId())
	data.Created = types.StringValue(getRequestV2Resp.GetCreated().Format(time.RFC3339))
	data.CreatedBy = types.StringValue(getRequestV2Resp.GetCreatedBy())
	data.LastUpdated = types.StringValue(getRequestV2Resp.GetLastUpdated().Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(getRequestV2Resp.GetLastUpdatedBy())
	data.Requested = setRequested(getRequestV2Resp.GetRequested())
	data.RequestedFor = setRequestedBy(getRequestV2Resp.GetRequestedFor())
	data.Status = types.StringValue(string(getRequestV2Resp.GetStatus()))
	data.AccessDuration = types.StringValue(getRequestV2Resp.GetAccessDuration())
	data.Granted = types.StringValue(getRequestV2Resp.GetGranted().Format(time.RFC3339))
	data.GrantStatus = types.StringValue(string(getRequestV2Resp.GetGrantStatus()))
	data.Resolved = types.StringValue(getRequestV2Resp.GetResolved().Format(time.RFC3339))
	data.RevocationStatus = types.StringValue(string(getRequestV2Resp.GetRevocationStatus()))
	data.RevocationScheduled = types.StringValue(string(getRequestV2Resp.GetRevocationScheduled().Format(time.RFC3339)))
	data.Revoked = types.StringValue(getRequestV2Resp.GetRevoked().Format(time.RFC3339))
	var requesterFieldValues []requestedFieldValues
	for _, reqValue := range getRequestV2Resp.GetRequesterFieldValues() {
		var requesterFieldValue requestedFieldValues
		requesterFieldValue.Id = types.StringValue(reqValue.GetId())
		requesterFieldValue.Label = types.StringValue(reqValue.GetLabel())
		requesterFieldValue.Type = types.StringValue(string(reqValue.GetType()))
		requesterFieldValue.Value = types.StringValue(reqValue.GetValue())
		var vals []values
		for _, val := range reqValue.Values {
			vals = append(vals, values{
				value: types.StringValue(val),
			})
		}
		requesterFieldValue.Values = vals
		requesterFieldValues = append(requesterFieldValues, requesterFieldValue)
	}
	data.RequesterFieldValues = requesterFieldValues
	var riskAssessments riskAssessment
	assessment := getRequestV2Resp.GetRiskAssessment()
	var rules []riskRules
	for _, riskRule := range assessment.GetRiskRules() {
		var rule riskRules
		rule.Name = types.StringValue(riskRule.GetName())
		rule.Description = types.StringValue(riskRule.GetDescription())
		rule.ResourceName = types.StringValue(riskRule.GetResourceName())
		rules = append(rules, rule)
	}
	riskAssessments.RequestSubmissionType = types.StringValue(string(assessment.GetRequestSubmissionType()))
	riskAssessments.RiskRules = rules

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update Not Supported",
		"This resource cannot be updated via Terraform.",
	)
}

func (r *requestV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"This resource cannot be deleted via Terraform.",
	)
}

func applyRequestResourceToState(data *requestV2ResourceModel, reqCreatableResp *governance.RequestSubmissionFull) {
	data.Id = types.StringValue(reqCreatableResp.GetId())
	data.Created = types.StringValue(reqCreatableResp.GetCreated().Format(time.RFC3339))
	data.CreatedBy = types.StringValue(reqCreatableResp.GetCreatedBy())
	data.LastUpdated = types.StringValue(reqCreatableResp.GetLastUpdated().Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(reqCreatableResp.GetLastUpdatedBy())
	data.Requested = setRequested(reqCreatableResp.GetRequested())
	data.RequestedFor = setRequestedBy(reqCreatableResp.GetRequestedFor())
	data.Status = types.StringValue(reqCreatableResp.GetStatus())
	data.AccessDuration = types.StringValue(reqCreatableResp.GetAccessDuration())
	data.Granted = types.StringValue(reqCreatableResp.GetGranted().Format(time.RFC3339))
	data.GrantStatus = types.StringValue(string(reqCreatableResp.GetGrantStatus()))
	data.Resolved = types.StringValue(reqCreatableResp.GetResolved().Format(time.RFC3339))
	data.RevocationStatus = types.StringValue(string(reqCreatableResp.GetRevocationStatus()))
	data.RevocationScheduled = types.StringValue(string(reqCreatableResp.GetRevocationScheduled().Format(time.RFC3339)))
	data.Revoked = types.StringValue(reqCreatableResp.GetRevoked().Format(time.RFC3339))
	var requesterFieldValues []requestedFieldValues
	for _, reqValue := range reqCreatableResp.GetRequesterFieldValues() {
		var requesterFieldValue requestedFieldValues
		requesterFieldValue.Id = types.StringValue(reqValue.GetId())
		requesterFieldValue.Label = types.StringValue(reqValue.GetLabel())
		requesterFieldValue.Type = types.StringValue(string(reqValue.GetType()))
		requesterFieldValue.Value = types.StringValue(reqValue.GetValue())
		var vals []values
		for _, val := range reqValue.Values {
			vals = append(vals, values{value: types.StringValue(val)})
		}
		requesterFieldValue.Values = vals
		requesterFieldValues = append(requesterFieldValues, requesterFieldValue)
	}
	data.RequesterFieldValues = requesterFieldValues
	var riskAssessments riskAssessment
	assessment := reqCreatableResp.GetRiskAssessment()
	var rules []riskRules
	for _, riskRule := range assessment.GetRiskRules() {
		var rule riskRules
		rule.Name = types.StringValue(riskRule.GetName())
		rule.Description = types.StringValue(riskRule.GetDescription())
		rule.ResourceName = types.StringValue(riskRule.GetResourceName())
		rules = append(rules, rule)
	}
	riskAssessments.RequestSubmissionType = types.StringValue(string(assessment.GetRequestSubmissionType()))
	riskAssessments.RiskRules = rules
}

func setRequestedBy(by governance.TargetPrincipal) *entitlementParentModel {
	return &entitlementParentModel{
		Type:       types.StringValue(string(by.GetType())),
		ExternalID: types.StringValue(by.GetExternalId()),
	}
}

func setRequested(getRequested governance.Requested) *requested {
	var reqResource requested
	reqResource.EntryId = types.StringValue(getRequested.GetEntryId())
	reqResource.AccessScopeId = types.StringValue(getRequested.GetAccessScopeId())
	reqResource.AccessScopeType = types.StringValue(string(getRequested.GetAccessScopeType()))
	reqResource.ResourceId = types.StringValue(getRequested.GetResourceId())
	reqResource.ResourceType = types.StringValue(string(getRequested.GetResourceType()))
	reqResource.Type = types.StringValue("CATALOG_ENTRY")
	return &reqResource
}

func createRequestReq(data requestV2ResourceModel) governance.RequestCreatable2 {
	var reqCreatable governance.RequestCreatable2
	reqCreatable.Requested = governance.RequestResourceCreatable{
		RequestResourceCatalogEntryCreatable: &governance.RequestResourceCatalogEntryCreatable{
			Type:    data.Requested.Type.ValueString(),
			EntryId: data.Requested.EntryId.ValueString(),
		},
	}
	reqCreatable.RequestedFor = governance.TargetPrincipal{
		ExternalId: data.RequestedFor.ExternalID.ValueString(),
		Type:       governance.PrincipalType(data.RequestedFor.Type.ValueString()),
	}

	var requesterFields []governance.RequestFieldValue
	for _, field := range data.RequesterFieldValues {
		var reqField governance.RequestFieldValue
		reqField.Id = field.Id.ValueString()
		reqField.Value = field.Value.ValueStringPointer()
		reqField.Label = field.Label.ValueStringPointer()
		reqFieldType := governance.RequestFieldType(field.Type.ValueString())
		reqField.Type = &reqFieldType
		var values []string
		values = append(values, reqField.Values...)
		reqField.Values = values

		requesterFields = append(requesterFields, reqField)
	}
	reqCreatable.RequesterFieldValues = requesterFields
	return reqCreatable
}
