package idaas

import (
	"context"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

type realmAssignmentModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Priority            types.Int32  `tfsdk:"priority"`
	Status              types.String `tfsdk:"status"`
	IsDefault           types.Bool   `tfsdk:"is_default"`
	ProfileSourceID     types.String `tfsdk:"profile_source_id"`
	RealmId             types.String `tfsdk:"realm_id"`
	ConditionExpression types.String `tfsdk:"condition_expression"`
}

type realmAssignmentResource struct {
	config *config.Config
}

func newRealmAssignmentResource() resource.Resource {
	return &realmAssignmentResource{}
}

func (r *realmAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_realm_assignment"
}

func (r *realmAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Realm Assignment ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Okta Realm Assignment.",
			},
			"priority": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The Priority of the Realm Assignment. The lower the number, the higher the priority.",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "Defines whether the Realm Assignment is active or not. Valid values: `ACTIVE` and `INACTIVE`.",
				Default:     stringdefault.StaticString("INACTIVE"),
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("ACTIVE", "INACTIVE"),
				},
			},
			"profile_source_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Profile Source.",
			},
			"realm_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Realm asscociated with the Realm Assignment.",
			},
			"condition_expression": schema.StringAttribute{
				Optional:            true,
				Description:         "Condition expression for the Realm Assignment in Okta Expression Language. Example: `user.profile.role ==\"Manager\"` or `user.profile.state.contains(\"example\")`.",
				MarkdownDescription: "Condition expression for the Realm Assignment in Okta Expression Language. Example: `user.profile.role ==\"Manager\"` or `user.profile.state.contains(\"example\")`.",
			},
			"is_default": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether the realm assignment is the default.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Description: "Creates an Okta Realm Assignment. This resource allows you to create and configure an Okta Realm Assignment.",
	}
}

func (r *realmAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.config = resourceConfiguration(req, resp)
}

func (r *realmAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state realmAssignmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRealmAssignmentRequest := v5okta.NewCreateRealmAssignmentRequest()
	createRealmAssignmentRequest.Name = state.Name.ValueStringPointer()
	createRealmAssignmentRequest.Priority = state.Priority.ValueInt32Pointer()

	actions := *v5okta.NewActionsWithDefaults()
	assignment := v5okta.NewAssignUserToRealmWithDefaults()
	assignment.RealmId = state.RealmId.ValueStringPointer()
	actions.AssignUserToRealm = assignment
	createRealmAssignmentRequest.Actions = &actions

	conditions := *v5okta.NewConditionsWithDefaults()
	conditions.ProfileSourceId = state.ProfileSourceID.ValueStringPointer()
	expression := v5okta.NewExpressionWithDefaults()
	expression.Value = state.ConditionExpression.ValueStringPointer()
	conditions.Expression = expression
	createRealmAssignmentRequest.Conditions = &conditions

	realmAssignment, response, err := r.config.OktaIDaaSClient.OktaSDKClientV5().RealmAssignmentAPI.CreateRealmAssignment(ctx).Body(*createRealmAssignmentRequest).Execute()
	if err != nil {
		body, ioErr := io.ReadAll(response.Body)
		defer response.Body.Close()
		if ioErr != nil {
			resp.Diagnostics.AddError(err.Error(), "failed to read response body")
			return
		}
		resp.Diagnostics.AddError("failed to create realm assignment:"+err.Error(), string(body))
		return
	}

	if state.Status.ValueString() == "ACTIVE" {
		response, err = r.config.OktaIDaaSClient.OktaSDKClientV5().RealmAssignmentAPI.ActivateRealmAssignment(ctx, *realmAssignment.Id).Execute()
		if err != nil {
			body, ioErr := io.ReadAll(response.Body)
			defer response.Body.Close()
			if ioErr != nil {
				resp.Diagnostics.AddError(err.Error(), "failed to read response body")
				return
			}
			resp.Diagnostics.AddError("failed to activate realm assignment:"+err.Error(), string(body))
			return
		}
		realmAssignment.Status = utils.StringPtr("ACTIVE")
	}

	resp.Diagnostics.Append(mapRealmAssignmentResourceToState(realmAssignment, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *realmAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state realmAssignmentModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	realmAssignment, _, err := r.config.OktaIDaaSClient.OktaSDKClientV5().RealmAssignmentAPI.GetRealmAssignment(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error getting realm assignment with id: %v", state.ID.ValueString()), err.Error())
		return
	}

	if realmAssignment == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(mapRealmAssignmentResourceToState(realmAssignment, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *realmAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state realmAssignmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.config.OktaIDaaSClient.OktaSDKClientV5().RealmAssignmentAPI.DeactivateRealmAssignment(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		body, ioErr := io.ReadAll(response.Body)
		defer response.Body.Close()
		if ioErr != nil {
			resp.Diagnostics.AddError(err.Error(), "failed to read response body")
			return
		}
		resp.Diagnostics.AddError("failed to deactivate realm assignment before updating:"+err.Error(), string(body))
		return
	}

	updateRealmAssignmentRequest := v5okta.NewUpdateRealmAssignmentRequest()
	updateRealmAssignmentRequest.Name = state.Name.ValueStringPointer()
	updateRealmAssignmentRequest.Priority = state.Priority.ValueInt32Pointer()

	actions := *v5okta.NewActionsWithDefaults()
	assignment := v5okta.NewAssignUserToRealmWithDefaults()
	assignment.RealmId = state.RealmId.ValueStringPointer()
	actions.AssignUserToRealm = assignment
	updateRealmAssignmentRequest.Actions = &actions

	conditions := *v5okta.NewConditionsWithDefaults()
	conditions.ProfileSourceId = state.ProfileSourceID.ValueStringPointer()
	expression := v5okta.NewExpressionWithDefaults()
	expression.Value = state.ConditionExpression.ValueStringPointer()
	conditions.Expression = expression
	updateRealmAssignmentRequest.Conditions = &conditions

	realmAssignment, response, err := r.config.OktaIDaaSClient.OktaSDKClientV5().RealmAssignmentAPI.ReplaceRealmAssignment(ctx, state.ID.ValueString()).Body(*updateRealmAssignmentRequest).Execute()
	if err != nil {
		body, ioErr := io.ReadAll(response.Body)
		defer response.Body.Close()
		if ioErr != nil {
			resp.Diagnostics.AddError(err.Error(), "failed to read response body")
			return
		}
		resp.Diagnostics.AddError("failed to update realm assignment:"+err.Error(), string(body))

		return
	}

	if state.Status.ValueString() == "ACTIVE" {
		response, err = r.config.OktaIDaaSClient.OktaSDKClientV5().RealmAssignmentAPI.ActivateRealmAssignment(ctx, *realmAssignment.Id).Execute()
		if err != nil {
			body, ioErr := io.ReadAll(response.Body)
			defer response.Body.Close()
			if ioErr != nil {
				resp.Diagnostics.AddError(err.Error(), "failed to read response body")
				return
			}
			resp.Diagnostics.AddError("failed to activate realm assignment:"+err.Error(), string(body))
			return
		}
		realmAssignment.Status = utils.StringPtr("ACTIVE")
	}

	resp.Diagnostics.Append(mapRealmAssignmentResourceToState(realmAssignment, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *realmAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state realmAssignmentModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.config.OktaIDaaSClient.OktaSDKClientV5().RealmAssignmentAPI.DeactivateRealmAssignment(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		body, ioErr := io.ReadAll(response.Body)
		defer response.Body.Close()
		if ioErr != nil {
			resp.Diagnostics.AddError(err.Error(), "failed to read response body")
			return
		}
		resp.Diagnostics.AddError("failed to deactivate realm assignment before deletion:"+err.Error(), string(body))
		return
	}

	response, err = r.config.OktaIDaaSClient.OktaSDKClientV5().RealmAssignmentAPI.DeleteRealmAssignment(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		body, ioErr := io.ReadAll(response.Body)
		defer response.Body.Close()
		if ioErr != nil {
			resp.Diagnostics.AddError(err.Error(), "failed to read response body")
			return
		}
		resp.Diagnostics.AddError("failed to delete realm assignment:"+err.Error(), string(body))
		return
	}
}

func (r *realmAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapRealmAssignmentResourceToState(realmAssignmentResource *v5okta.RealmAssignment, state *realmAssignmentModel) diag.Diagnostics {
	var diags diag.Diagnostics

	state.ID = types.StringPointerValue(realmAssignmentResource.Id)
	state.Name = types.StringPointerValue(realmAssignmentResource.Name)
	state.ConditionExpression = types.StringPointerValue(realmAssignmentResource.Conditions.Expression.Value)
	state.Priority = types.Int32PointerValue(realmAssignmentResource.Priority)
	state.Status = types.StringPointerValue(realmAssignmentResource.Status)
	state.IsDefault = types.BoolPointerValue(realmAssignmentResource.IsDefault)
	state.ProfileSourceID = types.StringPointerValue(realmAssignmentResource.Conditions.ProfileSourceId)
	state.RealmId = types.StringPointerValue(realmAssignmentResource.Actions.AssignUserToRealm.RealmId)

	return diags
}
