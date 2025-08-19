package governance

import (
	"context"
	"example.com/aditya-okta/okta-ig-sdk-golang/governance"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &requestTypeResource{}
	_ resource.ResourceWithConfigure   = &requestTypeResource{}
	_ resource.ResourceWithImportState = &requestTypeResource{}
)

func newRequestTypeResource() resource.Resource {
	return &requestTypeResource{}
}

func (r *requestTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_type"
}

func (r *requestTypeResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

func (r *requestTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type requestTypeResource struct {
	*config.Config
}

type requestTypeTargetResourceModel struct {
	ResourceId []string `tfsdk:"resource_id"`
}

type requestTypeResourceSettingsModel struct {
	Types           string                           `tfsdk:"types"`
	TargetResources []requestTypeTargetResourceModel `tfsdk:"target_resources"`
}

type approvalsModel struct {
	ApproverMemberOf []requestMemberOf     `tfsdk:"approver_member_of"`
	ApproverType     types.String          `tfsdk:"approver_type"`
	ApproverUserId   types.String          `tfsdk:"approver_user_id"`
	ApproverFields   []requesterFieldModel `tfsdk:"approver_fields"`
	Description      types.String          `tfsdk:"description"`
}

type requestTypeApproverSettingsModel struct {
	Approvals []approvalsModel `tfsdk:"approvals"`
	Type      string           `tfsdk:"type"`
}

type requestTypeResourceModel struct {
	Id               types.String                      `tfsdk:"id"`
	Name             types.String                      `tfsdk:"name"`
	OwnerID          types.String                      `tfsdk:"owner_id"`
	AccessDuration   types.String                      `tfsdk:"access_duration"`
	Description      types.String                      `tfsdk:"description"`
	Status           types.String                      `tfsdk:"status"`
	RequestSettings  *requestSettingsModel             `tfsdk:"request_settings"`
	ResourceSettings *requestTypeResourceSettingsModel `tfsdk:"resource_settings"`
	ApprovalSettings *requestTypeApproverSettingsModel `tfsdk:"approval_settings"`
}

type requestMemberOf struct {
	GroupId types.String `tfsdk:"group_id"`
}

type requestSettingsModel struct {
	RequesterMemberOf []requestMemberOf     `tfsdk:"requester_member_of"`
	Type              types.String          `tfsdk:"type"`
	RequesterFields   []requesterFieldModel `tfsdk:"requester_fields"`
}

type requesterFieldModel struct {
	Id       types.String           `tfsdk:"id"`
	Prompt   types.String           `tfsdk:"prompt"`
	Type     types.String           `tfsdk:"type"`
	Required types.Bool             `tfsdk:"required"`
	Options  []requesterOptionModel `tfsdk:"options"`
}

type requesterOptionModel struct {
	Value types.String `tfsdk:"value"`
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
				Computed: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"status": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"request_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("EVERYONE", "MEMBER_OF"),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"requester_fields": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed: true,
								},
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
					"requester_member_of": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"group_id": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
			"resource_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"types": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("GROUPS", "APPS", "ENTITLEMENT_BUNDLES"),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"target_resources": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"resource_id": schema.ListAttribute{
									Required:    true,
									ElementType: types.StringType,
									Validators: []validator.List{
										listvalidator.SizeAtMost(1),
									},
								},
							},
						},
						Validators: []validator.List{
							listvalidator.SizeAtMost(1),
						},
					},
				},
			},
			"approval_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("SERIAL", "NONE"),
						},
					},
				},
				Blocks: map[string]schema.Block{
					"approvals": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"approver_type": schema.StringAttribute{
									Required: true,
								},
								"approver_user_id": schema.StringAttribute{
									Optional: true,
								},
								"description": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
							},
							Blocks: map[string]schema.Block{
								"approver_member_of": schema.SingleNestedBlock{
									Attributes: map[string]schema.Attribute{
										"group_id": schema.StringAttribute{
											Optional: true,
										},
									},
								},
								"approver_fields": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Computed: true,
											},
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
			},
		},
	}
}

func (r *requestTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data requestTypeResourceModel

	// Read Terraform plan Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic

	requestTypeResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestTypesAPI.CreateRequestType(ctx).RequestTypeCreatable(*createRequest(data)).Execute() // Panic likely here
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Request type",
			"Could not creating Request type, unexpected error: "+err.Error(),
		)
		return
	}
	// Example Data value setting
	data.Id = types.StringValue(requestTypeResp.GetId())
	data.Name = types.StringValue(requestTypeResp.GetName())
	data.OwnerID = types.StringValue(requestTypeResp.GetOwnerId())
	data.AccessDuration = types.StringValue(requestTypeResp.GetAccessDuration())
	data.Description = types.StringValue(requestTypeResp.GetDescription())
	data.Status = types.StringValue(string(requestTypeResp.GetStatus()))
	data.ResourceSettings = setResourceSettings(requestTypeResp.ResourceSettings)
	data.RequestSettings = setRequestSettings(requestTypeResp.RequestSettings)
	data.ApprovalSettings = setApproverSettings(requestTypeResp.ApprovalSettings)
	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data requestTypeResourceModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	readRequestTypeResp, _, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestTypesAPI.GetRequestType(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Request type",
			"Could not reading Request type, unexpected error: "+err.Error(),
		)
		return
	}
	data.Id = types.StringValue(readRequestTypeResp.GetId())
	data.Name = types.StringValue(readRequestTypeResp.GetName())
	data.OwnerID = types.StringValue(readRequestTypeResp.GetOwnerId())
	data.AccessDuration = types.StringValue(readRequestTypeResp.GetAccessDuration())
	data.Description = types.StringValue(readRequestTypeResp.GetDescription())
	data.Status = types.StringValue(string(readRequestTypeResp.GetStatus()))
	data.ResourceSettings = setResourceSettings(readRequestTypeResp.ResourceSettings)
	data.RequestSettings = setRequestSettings(readRequestTypeResp.RequestSettings)
	data.ApprovalSettings = setApproverSettings(readRequestTypeResp.ApprovalSettings)

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update Not Supported",
		"No fields are updatable for this resource. Terraform will retain the existing state.",
	)
}

func (r *requestTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data requestTypeResourceModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	_, err := r.OktaGovernanceClient.OktaGovernanceSDKClient().RequestTypesAPI.DeleteRequestType(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Request type",
			"Could not delete Request type, unexpected error: "+err.Error(),
		)
		return
	}
}

func setRequestSettings(settings governance.RequestTypeRequestSettingsReadable) *requestSettingsModel {
	var result requestSettingsModel
	if settings.RequestTypeRequesterMemberOf != nil {
		var requesterFields []requesterFieldModel
		for _, rf := range settings.RequestTypeRequesterMemberOf.RequesterFields {
			var requestedField requesterFieldModel
			if rf.FieldText != nil {
				requestedField.Id = types.StringValue(rf.FieldText.GetId())
				requestedField.Type = types.StringValue(string(rf.FieldText.GetType()))
				requestedField.Prompt = types.StringValue(rf.FieldText.GetPrompt())
				requestedField.Required = types.BoolValue(rf.FieldText.GetRequired())
			} else if rf.FieldSelect != nil {
				requestedField.Id = types.StringValue(rf.FieldSelect.GetId())
				requestedField.Type = types.StringValue(string(rf.FieldSelect.GetType()))
				requestedField.Prompt = types.StringValue(rf.FieldSelect.GetPrompt())
				requestedField.Required = types.BoolValue(rf.FieldSelect.GetRequired())
				var option []requesterOptionModel
				for _, val := range rf.FieldSelect.GetOptions() {
					option = append(option, requesterOptionModel{Value: types.StringValue(val.GetValue())})
				}
				requestedField.Options = option
			} else if rf.FieldDate != nil {
				requestedField.Id = types.StringValue(rf.FieldDate.GetId())
				requestedField.Type = types.StringValue(string(rf.FieldDate.GetType()))
				requestedField.Prompt = types.StringValue(rf.FieldDate.GetPrompt())
				requestedField.Required = types.BoolValue(rf.FieldDate.GetRequired())
			}
			requesterFields = append(requesterFields, requestedField)
		}
		var ids []requestMemberOf
		result.RequesterFields = requesterFields
		result.Type = types.StringValue(settings.RequestTypeRequesterMemberOf.GetType())
		for _, id := range settings.RequestTypeRequesterMemberOf.RequesterMemberOf {
			ids = append(ids, requestMemberOf{
				GroupId: types.StringValue(id),
			})
		}
		result.RequesterMemberOf = ids
		return &result
	} else if settings.RequestTypeRequesterCustom != nil {
		result.Type = types.StringValue(settings.RequestTypeRequesterCustom.GetType())
		return &result
	} else if settings.RequestTypeRequesterEveryone != nil {
		var requesterFields []requesterFieldModel
		for _, rf := range settings.RequestTypeRequesterEveryone.RequesterFields {
			var requestedField requesterFieldModel
			if rf.FieldText != nil {
				requestedField.Id = types.StringValue(rf.FieldText.GetId())
				requestedField.Type = types.StringValue(string(rf.FieldText.GetType()))
				requestedField.Prompt = types.StringValue(rf.FieldText.GetPrompt())
				requestedField.Required = types.BoolValue(rf.FieldText.GetRequired())
			} else if rf.FieldSelect != nil {
				requestedField.Id = types.StringValue(rf.FieldSelect.GetId())
				requestedField.Type = types.StringValue(string(rf.FieldSelect.GetType()))
				requestedField.Prompt = types.StringValue(rf.FieldSelect.GetPrompt())
				requestedField.Required = types.BoolValue(rf.FieldSelect.GetRequired())
				var option []requesterOptionModel
				for _, val := range rf.FieldSelect.GetOptions() {
					option = append(option, requesterOptionModel{Value: types.StringValue(val.GetValue())})
				}
				requestedField.Options = option
			} else if rf.FieldDate != nil {
				requestedField.Id = types.StringValue(rf.FieldDate.GetId())
				requestedField.Type = types.StringValue(string(rf.FieldDate.GetType()))
				requestedField.Prompt = types.StringValue(rf.FieldDate.GetPrompt())
				requestedField.Required = types.BoolValue(rf.FieldDate.GetRequired())
			}
			requesterFields = append(requesterFields, requestedField)
		}
		result.RequesterFields = requesterFields
		result.Type = types.StringValue(settings.RequestTypeRequesterEveryone.GetType())
		return &result
	}
	return nil
}

func setResourceSettings(settings governance.RequestTypeResourceSettingsReadable) *requestTypeResourceSettingsModel {
	if settings.RequestTypeResourceSettingsApps != nil {
		var resourceIds []string
		for _, tr := range settings.RequestTypeResourceSettingsApps.TargetResources {
			resourceIds = append(resourceIds, tr.ResourceId)
		}
		return &requestTypeResourceSettingsModel{
			TargetResources: []requestTypeTargetResourceModel{
				{
					ResourceId: resourceIds,
				},
			},
			Types: "APPS",
		}
	} else if settings.RequestTypeResourceSettingsEntitlementBundles != nil {
		var resourceIds []string
		for _, tr := range settings.RequestTypeResourceSettingsEntitlementBundles.TargetResources {
			resourceIds = append(resourceIds, tr.ResourceId)
		}
		return &requestTypeResourceSettingsModel{
			TargetResources: []requestTypeTargetResourceModel{
				{
					ResourceId: resourceIds,
				},
			},
			Types: "ENTITLEMENT_BUNDLES",
		}
	} else if settings.RequestTypeResourceSettingsGroups != nil {
		var resourceIds []string
		for _, tr := range settings.RequestTypeResourceSettingsGroups.TargetResources {
			resourceIds = append(resourceIds, tr.ResourceId)
		}
		return &requestTypeResourceSettingsModel{
			TargetResources: []requestTypeTargetResourceModel{
				{
					ResourceId: resourceIds,
				},
			},
			Types: "GROUPS",
		}
	}
	return nil
}

func setApproverSettings(settings governance.RequestTypeApprovalSettingsReadable) *requestTypeApproverSettingsModel {
	var approverSettings requestTypeApproverSettingsModel
	if settings.RequestTypeApprovalSettingsSerial != nil {
		approverSettings.Type = settings.RequestTypeApprovalSettingsSerial.GetType()
		var approvals []approvalsModel
		for _, as := range settings.RequestTypeApprovalSettingsSerial.Approvals {
			var approval approvalsModel
			if as.RequestTypeApprovalMemberOf != nil {
				approval.ApproverType = types.StringValue(as.RequestTypeApprovalMemberOf.GetApproverType())
				approval.Description = types.StringValue(as.RequestTypeApprovalMemberOf.GetDescription())
				var approvalMemberOf []requestMemberOf
				for _, m := range as.RequestTypeApprovalMemberOf.GetApproverMemberOf() {
					approvalMemberOf = append(approvalMemberOf, requestMemberOf{GroupId: types.StringValue(m)})
				}
				approval.ApproverMemberOf = approvalMemberOf

				var requesterFields []requesterFieldModel
				for _, rf := range as.RequestTypeApprovalMemberOf.GetApproverFields() {
					var requestedField requesterFieldModel
					if rf.FieldText != nil {
						requestedField.Id = types.StringValue(rf.FieldText.GetId())
						requestedField.Type = types.StringValue(string(rf.FieldText.GetType()))
						requestedField.Prompt = types.StringValue(rf.FieldText.GetPrompt())
						requestedField.Required = types.BoolValue(rf.FieldText.GetRequired())
					} else if rf.FieldSelect != nil {
						requestedField.Id = types.StringValue(rf.FieldSelect.GetId())
						requestedField.Type = types.StringValue(string(rf.FieldSelect.GetType()))
						requestedField.Prompt = types.StringValue(rf.FieldSelect.GetPrompt())
						requestedField.Required = types.BoolValue(rf.FieldSelect.GetRequired())
						var option []requesterOptionModel
						for _, val := range rf.FieldSelect.GetOptions() {
							option = append(option, requesterOptionModel{Value: types.StringValue(val.GetValue())})
						}
						requestedField.Options = option
					} else if rf.FieldDate != nil {
						requestedField.Id = types.StringValue(rf.FieldDate.GetId())
						requestedField.Type = types.StringValue(string(rf.FieldDate.GetType()))
						requestedField.Prompt = types.StringValue(rf.FieldDate.GetPrompt())
						requestedField.Required = types.BoolValue(rf.FieldDate.GetRequired())
					}
					requesterFields = append(requesterFields, requestedField)
				}
				approval.ApproverFields = requesterFields
				approvals = append(approvals, approval)
			} else if as.RequestTypeApprovalResourceOwner != nil {
				approval.ApproverType = types.StringValue(as.RequestTypeApprovalResourceOwner.GetApproverType())
				approval.Description = types.StringValue(as.RequestTypeApprovalResourceOwner.GetDescription())
				approvals = append(approvals, approval)
			} else if as.RequestTypeApprovalUser != nil {
				approval.ApproverType = types.StringValue(as.RequestTypeApprovalUser.GetApproverType())
				approval.Description = types.StringValue(as.RequestTypeApprovalUser.GetDescription())
				approval.ApproverUserId = types.StringValue(as.RequestTypeApprovalUser.GetApproverUserId())
				var requesterFields []requesterFieldModel
				for _, rf := range as.RequestTypeApprovalMemberOf.GetApproverFields() {
					var requestedField requesterFieldModel
					if rf.FieldText != nil {
						requestedField.Id = types.StringValue(rf.FieldText.GetId())
						requestedField.Type = types.StringValue(string(rf.FieldText.GetType()))
						requestedField.Prompt = types.StringValue(rf.FieldText.GetPrompt())
						requestedField.Required = types.BoolValue(rf.FieldText.GetRequired())
					} else if rf.FieldSelect != nil {
						requestedField.Id = types.StringValue(rf.FieldSelect.GetId())
						requestedField.Type = types.StringValue(string(rf.FieldSelect.GetType()))
						requestedField.Prompt = types.StringValue(rf.FieldSelect.GetPrompt())
						requestedField.Required = types.BoolValue(rf.FieldSelect.GetRequired())
						var option []requesterOptionModel
						for _, val := range rf.FieldSelect.GetOptions() {
							option = append(option, requesterOptionModel{Value: types.StringValue(val.GetValue())})
						}
						requestedField.Options = option
					} else if rf.FieldDate != nil {
						requestedField.Id = types.StringValue(rf.FieldDate.GetId())
						requestedField.Type = types.StringValue(string(rf.FieldDate.GetType()))
						requestedField.Prompt = types.StringValue(rf.FieldDate.GetPrompt())
						requestedField.Required = types.BoolValue(rf.FieldDate.GetRequired())
					}
					requesterFields = append(requesterFields, requestedField)
				}
				approval.ApproverFields = requesterFields
				approvals = append(approvals, approval)
			} else if as.RequestTypeApprovalManager != nil {
				approval.ApproverType = types.StringValue(as.RequestTypeApprovalManager.GetApproverType())
				approval.Description = types.StringValue(as.RequestTypeApprovalManager.GetDescription())
				var requesterFields []requesterFieldModel
				for _, rf := range as.RequestTypeApprovalManager.GetApproverFields() {
					var requestedField requesterFieldModel
					if rf.FieldText != nil {
						requestedField.Id = types.StringValue(rf.FieldText.GetId())
						requestedField.Type = types.StringValue(string(rf.FieldText.GetType()))
						requestedField.Prompt = types.StringValue(rf.FieldText.GetPrompt())
						requestedField.Required = types.BoolValue(rf.FieldText.GetRequired())
					} else if rf.FieldSelect != nil {
						requestedField.Id = types.StringValue(rf.FieldSelect.GetId())
						requestedField.Type = types.StringValue(string(rf.FieldSelect.GetType()))
						requestedField.Prompt = types.StringValue(rf.FieldSelect.GetPrompt())
						requestedField.Required = types.BoolValue(rf.FieldSelect.GetRequired())
						var option []requesterOptionModel
						for _, val := range rf.FieldSelect.GetOptions() {
							option = append(option, requesterOptionModel{Value: types.StringValue(val.GetValue())})
						}
						requestedField.Options = option
					} else if rf.FieldDate != nil {
						requestedField.Id = types.StringValue(rf.FieldDate.GetId())
						requestedField.Type = types.StringValue(string(rf.FieldDate.GetType()))
						requestedField.Prompt = types.StringValue(rf.FieldDate.GetPrompt())
						requestedField.Required = types.BoolValue(rf.FieldDate.GetRequired())
					}
					requesterFields = append(requesterFields, requestedField)
				}
				approval.ApproverFields = requesterFields
				approvals = append(approvals, approval)
			}
		}
		approverSettings.Approvals = approvals
		return &approverSettings
	} else if settings.RequestTypeApprovalSettingsNone != nil {
		approverSettings.Type = settings.RequestTypeApprovalSettingsNone.GetType()
		return &approverSettings
	} else if settings.RequestTypeApprovalSettingsCustom != nil {
		approverSettings.Type = settings.RequestTypeApprovalSettingsCustom.GetType()
		return &approverSettings
	}
	return nil
}

func createRequest(data requestTypeResourceModel) *governance.RequestTypeCreatable {
	var requestType governance.RequestTypeCreatable
	requestType.SetName(data.Name.ValueString())
	requestType.SetOwnerId(data.OwnerID.ValueString())
	if data.AccessDuration.ValueString() != "" {
		requestType.SetAccessDuration(data.AccessDuration.ValueString())
	}
	requestType.SetDescription(data.Description.ValueString())
	if data.Status.ValueString() != "" {
		requestType.SetStatus((governance.RequestTypeCreatableStatus)(data.Status.ValueString()))
	}
	requestSettings := createRequestSettings(data.RequestSettings)
	if requestSettings != nil {
		requestType.SetRequestSettings(*requestSettings)
	} else {
		requestType.SetRequestSettings(governance.RequestTypeRequestSettingsMutable{
			RequestTypeRequesterEveryoneWritable: &governance.RequestTypeRequesterEveryoneWritable{
				Type:            "EVERYONE",
				RequesterFields: []governance.FieldWritable{},
			},
		})
	}
	requestType.SetResourceSettings(*createResourceSettings(data.ResourceSettings))
	requestType.SetApprovalSettings(*createApprovalSettings(data.ApprovalSettings))
	return &requestType
}

func createApprovalSettings(settings *requestTypeApproverSettingsModel) *governance.RequestTypeApprovalSettingsMutable {
	if settings.Type == "NONE" {
		return nil
	} else if settings.Type == "SERIAL" {
		var approvalSettings governance.RequestTypeApprovalSettingsMutable
		approvalSettings.RequestTypeApprovalSettingsSerialWritable = &governance.RequestTypeApprovalSettingsSerialWritable{}
		approvalSettings.RequestTypeApprovalSettingsSerialWritable.Type = settings.Type
		var approvals []governance.RequestTypeApprovalWritable
		for _, app := range settings.Approvals {
			if app.ApproverType.ValueString() == "MEMBER_OF" {
				var approval governance.RequestTypeApprovalMemberOfWritable
				approval.ApproverType = app.ApproverType.ValueString()
				if !app.Description.IsNull() && app.Description.ValueString() != "" {
					approval.Description = app.Description.ValueStringPointer()
				}
				for _, member := range app.ApproverMemberOf {
					approval.ApproverMemberOf = append(approval.ApproverMemberOf, member.GroupId.ValueString())
				}
				var approverFields []governance.FieldWritable
				getApproverFields(app, approverFields)
				approval.ApproverFields = approverFields
				approvals = append(approvals, governance.RequestTypeApprovalWritable{
					RequestTypeApprovalMemberOfWritable: &approval,
				})
			} else if app.ApproverType.ValueString() == "RESOURCE_OWNER" {
				var approval governance.RequestTypeApprovalResourceOwnerWritable
				approval.ApproverType = app.ApproverType.ValueString()
				if !app.Description.IsNull() && app.Description.ValueString() != "" {
					approval.Description = app.Description.ValueStringPointer()
				}
				approvals = append(approvals, governance.RequestTypeApprovalWritable{
					RequestTypeApprovalResourceOwnerWritable: &approval,
				})
			} else if app.ApproverType.ValueString() == "MANAGER" {
				var approval governance.RequestTypeApprovalManagerWritable
				approval.ApproverType = app.ApproverType.ValueString()
				if !app.Description.IsNull() && app.Description.ValueString() != "" {
					approval.Description = app.Description.ValueStringPointer()
				}
				var approverFields []governance.FieldWritable
				getApproverFields(app, approverFields)
				approval.ApproverFields = approverFields
				approvals = append(approvals, governance.RequestTypeApprovalWritable{
					RequestTypeApprovalManagerWritable: &approval,
				})
			} else if app.ApproverType.ValueString() == "USER" {
				var approval governance.RequestTypeApprovalUserWritable
				approval.ApproverType = app.ApproverType.ValueString()
				if !app.Description.IsNull() && app.Description.ValueString() != "" {
					approval.Description = app.Description.ValueStringPointer()
				}
				approval.ApproverUserId = app.ApproverUserId.ValueString()
				var approverFields []governance.FieldWritable
				getApproverFields(app, approverFields)
				approval.ApproverFields = approverFields
				approvals = append(approvals, governance.RequestTypeApprovalWritable{
					RequestTypeApprovalUserWritable: &approval,
				})
			}
		}
		approvalSettings.RequestTypeApprovalSettingsSerialWritable.Approvals = approvals
		return &approvalSettings
	}
	return nil
}

func getApproverFields(app approvalsModel, fields []governance.FieldWritable) {
	for _, field := range app.ApproverFields {
		if field.Type.ValueString() == "SELECT" {
			var selectWritable *governance.FieldSelectWritable
			selectWritable.Prompt = field.Prompt.ValueString()
			selectWritable.Required = field.Required.ValueBoolPointer()
			selectWritable.Type = governance.FieldSelectType(field.Type.ValueString())
			var options []governance.FieldOption
			for _, option := range field.Options {
				op := option.Value.String()
				options = append(options, governance.FieldOption{Value: op})
			}
			selectWritable.Options = options
			fields = append(fields, governance.FieldWritable{FieldSelectWritable: selectWritable})
		} else if field.Type.ValueString() == "TEXT" {
			var textWritable *governance.FieldTextWritable
			textWritable.Prompt = field.Prompt.ValueString()
			textWritable.Required = field.Required.ValueBoolPointer()
			textWritable.Type = governance.FieldTextType(field.Type.ValueString())
			fields = append(fields, governance.FieldWritable{FieldTextWritable: textWritable})
		} else if field.Type.ValueString() == "DATE-TIME" {
			var dateWritable *governance.FieldDateWritable
			dateWritable.Prompt = field.Prompt.ValueString()
			dateWritable.Required = field.Required.ValueBoolPointer()
			dateWritable.Type = governance.FieldDateTimeType(field.Type.ValueString())
			fields = append(fields, governance.FieldWritable{FieldDateWritable: dateWritable})
		}
	}
}

func createResourceSettings(settings *requestTypeResourceSettingsModel) *governance.RequestTypeResourceSettingsMutable {
	if settings.Types == "APPS" {
		var resourceSettings governance.RequestTypeResourceSettingsApps
		resourceSettings.Type = "APPS"
		var targetResourcesId []governance.OktaApplicationResource
		for _, t := range settings.TargetResources {
			for _, tr := range t.ResourceId {
				targetResourcesId = append(targetResourcesId, governance.OktaApplicationResource{
					ResourceId: tr,
				})
			}
			resourceSettings.TargetResources = targetResourcesId
		}

		return &governance.RequestTypeResourceSettingsMutable{RequestTypeResourceSettingsApps: &resourceSettings}
	} else if settings.Types == "ENTITLEMENT_BUNDLES" {
		var resourceSettings governance.RequestTypeResourceSettingsEntitlementBundles
		resourceSettings.Type = "ENTITLEMENT_BUNDLES"
		var targetResourcesId []governance.OktaEntitlementBundleResource
		for _, t := range settings.TargetResources {
			for _, tr := range t.ResourceId {
				targetResourcesId = append(targetResourcesId, governance.OktaEntitlementBundleResource{
					ResourceId: tr,
				})
			}
			resourceSettings.TargetResources = targetResourcesId
		}

		return &governance.RequestTypeResourceSettingsMutable{RequestTypeResourceSettingsEntitlementBundles: &resourceSettings}
	} else if settings.Types == "GROUPS" {
		var resourceSettings governance.RequestTypeResourceSettingsGroups
		resourceSettings.Type = "GROUPS"
		var targetResourcesId []governance.OktaGroupResource
		for _, t := range settings.TargetResources {
			for _, tr := range t.ResourceId {
				targetResourcesId = append(targetResourcesId, governance.OktaGroupResource{
					ResourceId: tr,
				})
			}
			resourceSettings.TargetResources = targetResourcesId
		}

		return &governance.RequestTypeResourceSettingsMutable{RequestTypeResourceSettingsGroups: &resourceSettings}
	}
	return nil
}

func createRequestSettings(settings *requestSettingsModel) *governance.RequestTypeRequestSettingsMutable {
	if settings != nil {
		if settings.Type != types.StringNull() && settings.Type.ValueString() == "EVERYONE" {
			requestTypeRequester := &governance.RequestTypeRequesterEveryoneWritable{
				RequesterFields: []governance.FieldWritable{},
			}
			requestTypeRequester.Type = settings.Type.ValueString()
			var requesterFields []governance.FieldWritable
			getRequesterFields(*settings, requesterFields)
			return &governance.RequestTypeRequestSettingsMutable{
				RequestTypeRequesterEveryoneWritable: requestTypeRequester,
			}
		} else if settings.Type != types.StringNull() && settings.Type.ValueString() == "MEMBER_OF" {
			requestTypeRequester := &governance.RequestTypeRequesterMemberOfWritable{
				RequesterFields: []governance.FieldWritable{},
			}
			requestTypeRequester.Type = settings.Type.ValueString()
			var groups []string
			for _, group := range settings.RequesterMemberOf {
				groups = append(groups, group.GroupId.ValueString())
			}
			requestTypeRequester.RequesterMemberOf = groups
			var requesterFields []governance.FieldWritable
			getRequesterFields(*settings, requesterFields)
			return &governance.RequestTypeRequestSettingsMutable{
				RequestTypeRequesterMemberOfWritable: requestTypeRequester,
			}
		}
	}
	return nil
}

func getRequesterFields(settings requestSettingsModel, requesterFields []governance.FieldWritable) {
	for _, field := range settings.RequesterFields {
		if field.Type.ValueString() == "SELECT" {
			var selectWritable *governance.FieldSelectWritable
			selectWritable.Prompt = field.Prompt.ValueString()
			selectWritable.Required = field.Required.ValueBoolPointer()
			selectWritable.Type = governance.FieldSelectType(field.Type.ValueString())
			var options []governance.FieldOption
			for _, option := range field.Options {
				op := option.Value.String()
				options = append(options, governance.FieldOption{Value: op})
			}
			selectWritable.Options = options
			requesterFields = append(requesterFields, governance.FieldWritable{FieldSelectWritable: selectWritable})
		} else if field.Type.ValueString() == "TEXT" {
			var textWritable *governance.FieldTextWritable
			textWritable.Prompt = field.Prompt.ValueString()
			textWritable.Required = field.Required.ValueBoolPointer()
			textWritable.Type = governance.FieldTextType(field.Type.ValueString())
			requesterFields = append(requesterFields, governance.FieldWritable{FieldTextWritable: textWritable})
		} else if field.Type.ValueString() == "DATE-TIME" {
			var dateWritable *governance.FieldDateWritable
			dateWritable.Prompt = field.Prompt.ValueString()
			dateWritable.Required = field.Required.ValueBoolPointer()
			dateWritable.Type = governance.FieldDateTimeType(field.Type.ValueString())
			requesterFields = append(requesterFields, governance.FieldWritable{FieldDateWritable: dateWritable})
		}
	}
}
