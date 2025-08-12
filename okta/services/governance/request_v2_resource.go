package governance

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*requestV2Resource)(nil)

func NewRequestV2Resource() resource.Resource {
	return &requestV2Resource{}
}

type requestV2Resource struct{}

type requested struct {
	EntryId types.String `type:"entry_id"`
	Type    types.String `type:"type"`
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

type requestedFieldValues struct {
	Id     types.String `tfsdk:"id"`
	Label  types.String `tfsdk:"label"`
	Type   types.String `tfsdk:"type"`
	Value  types.String `tfsdk:"value"`
	Values types.List   `tfsdk:"values"`
}

type requestV2ResourceModel struct {
	Id                   types.String           `tfsdk:"id"`
	Created              types.String           `tfsdk:"created"`
	CreatedBy            types.String           `tfsdk:"created_by"`
	LastUpdated          types.String           `tfsdk:"last_updated"`
	LastUpdatedBy        types.String           `tfsdk:"last_updated_by"`
	Status               types.String           `tfsdk:"status"`
	AccessDuration       types.String           `tfsdk:"access_duration"`
	Granted              types.String           `tfsdk:"granted"`
	GrantStatus          types.String           `tfsdk:"grant_status"`
	Resolved             types.String           `tfsdk:"resolved"`
	RevocationScheduled  types.String           `tfsdk:"revocation_scheduled"`
	RevocationStatus     types.String           `tfsdk:"revocation_status"`
	Revoked              types.String           `tfsdk:"revoked"`
	RiskAssessment       riskAssessment         `tfsdk:"risk_assessment"`
	Requested            requested              `tfsdk:"requested"`
	RequestedFor         entitlementParentModel `tfsdk:"requested_for"`
	RequestedBy          entitlementParentModel `tfsdk:"requested_by"`
	RequesterFieldValues []requestedFieldValues `tfsdk:"requester_field_values"`
}

func (r *requestV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_v2"
}

func (r *requestV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"created_by": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"last_updated_by": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"access_duration": schema.StringAttribute{
				Computed: true,
			},
			"granted": schema.StringAttribute{
				Computed: true,
			},
			"grant_status": schema.StringAttribute{
				Computed: true,
			},
			"resolved": schema.StringAttribute{
				Computed: true,
			},
			"revocation_scheduled": schema.StringAttribute{
				Computed: true,
			},
			"revocation_status": schema.StringAttribute{
				Computed: true,
			},
			"revoked": schema.StringAttribute{
				Computed: true,
			},
			"risk_assessment": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"requested": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"entry_id": schema.StringAttribute{
						Required: true,
					},
					"type": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"requested_for": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Required: true,
					},
					"type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("OKTA_USER"),
						},
					},
				},
			},
			"requested_by": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"external_id": schema.StringAttribute{
						Required: true,
					},
					"type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("OKTA_USER"),
						},
					},
				},
			},
			"requester_field_values": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required: true,
						},
						"label": schema.StringAttribute{
							Optional: true,
						},
						"type": schema.StringAttribute{
							Optional: true,
						},
						"value": schema.StringAttribute{
							Optional: true,
						},
						"values": schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"risk_assessment": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"request_submission_type": schema.StringAttribute{
						Optional: true,
					},
				},
				Blocks: map[string]schema.Block{
					"risk_rules": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Optional: true,
								},
								"description": schema.StringAttribute{
									Optional: true,
								},
								"resource_name": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *requestV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data requestV2ResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic

	// Example data value setting
	data.Id = types.StringValue("example-id")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data requestV2ResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data requestV2ResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data requestV2ResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
}
