package governance

import (
	"context"
	"example.com/aditya-okta/okta-ig-sdk-golang/oktaInternalGovernance"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ resource.Resource = &requestTypeResource{}

func NewRequestTypeResource() resource.Resource {
	return &requestTypeResource{}
}

type requestTypeResource struct {
	*config.Config
}

type requestTypeResourceModel struct {
	Id              types.String           `tfsdk:"id"`
	Name            types.String           `tfsdk:"name"`
	OwnerID         types.String           `tfsdk:"owner_id"`
	AccessDuration  types.String           `tfsdk:"access_duration"`
	Description     types.String           `tfsdk:"description"`
	Status          types.String           `tfsdk:"status"`
	RequestSettings []requestSettingsModel `tfsdk:"request_settings"`
}

type requestSettingsModel struct {
	RequesterMemberOf types.String          `tfsdk:"requester_member_of"`
	Type              types.String          `tfsdk:"type"`
	RequesterFields   []requesterFieldModel `tfsdk:"requester_fields"`
}

type requesterFieldModel struct {
	Prompt   types.String           `tfsdk:"prompt"`
	Type     types.String           `tfsdk:"type"`
	Required types.Bool             `tfsdk:"required"`
	Options  []requesterOptionModel `tfsdk:"options"`
}

type requesterOptionModel struct {
	Value types.String `tfsdk:"value"`
}

func (r *requestTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_type"
}

func (r *requestTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"owner_id": schema.StringAttribute{
				Required: true,
			},
			"access_duration": schema.StringAttribute{
				Optional: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"status": schema.StringAttribute{
				Optional: true,
			},
		},
		Blocks: map[string]schema.Block{
			"request_settings": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"requester_member_of": schema.StringAttribute{
							Optional: true,
						},
						"type": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("EVERYONE", "MEMBER_OF"),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"requester_fields": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"prompt": schema.StringAttribute{
										Optional: true,
									},
									"type": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											stringvalidator.OneOf("DATE-TIME", "TEXT", "SELECT"),
										},
									},
									"required": schema.BoolAttribute{
										Optional: true,
									},
								},
								Blocks: map[string]schema.Block{
									"options": schema.ListNestedBlock{
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"value": schema.StringAttribute{
													Optional: true,
												},
											},
										},
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtMost(5),
							},
						},
					},
				},
			},
		},
	}
}

func (r *requestTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data requestTypeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	r.OktaGovernanceClient.OktaIGSDKClientV5().RequestTypesAPI.CreateRequestType(ctx).RequestTypeCreatable(createRequest(data)).Execute()
	// Example data value setting
	data.Id = types.StringValue("example-id")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func createRequest(data requestTypeResourceModel) oktaInternalGovernance.RequestTypeCreatable {
	return oktaInternalGovernance.RequestTypeCreatable{
		Name:            data.Name.ValueString(),
		OwnerId:         data.OwnerID.ValueString(),
		AccessDuration:  *oktaInternalGovernance.NewNullableString(data.AccessDuration.ValueStringPointer()),
		Description:     data.Description.ValueStringPointer(),
		Status:          (*oktaInternalGovernance.RequestTypeCreatableStatus)(data.Status.ValueStringPointer()),
		RequestSettings: createRequestSettings(data.RequestSettings),
	}
}

func createRequestSettings(settings []requestSettingsModel) *oktaInternalGovernance.RequestTypeRequestSettingsMutable {
	//if len(settings) == 0 {
	//	return nil
	//}
	//
	//requestSettings := make([]oktaInternalGovernance.RequestTypeRequestSettingsMutable, len(settings))
	//for i, setting := range settings {
	//	requestSettings[i] = oktaInternalGovernance.RequestTypeRequestSettingsMutable{
	//		RequesterMemberOf: setting.RequesterMemberOf.ValueStringPointer(),
	//		Type:              (*oktaInternalGovernance.RequestTypeRequestSettingsMutableType)(setting.Type.ValueStringPointer()),
	//		RequesterFields:   createRequesterFields(setting.RequesterFields),
	//	}
	//}
	//
	//return &oktaInternalGovernance.RequestTypeRequestSettingsMutable{
	//	RequesterFields: requestSettings,
	//}
	return nil
}

func (r *requestTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data requestTypeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data requestTypeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data requestTypeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
}
