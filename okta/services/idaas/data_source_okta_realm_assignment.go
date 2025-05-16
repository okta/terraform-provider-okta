package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

type realmAssignmentDataSource struct {
	config *config.Config
}

func newRealmAssignmentDataSource() datasource.DataSource {
	return &realmAssignmentDataSource{}
}

func (r *realmAssignmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_realm_assignment"
}

func (r *realmAssignmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	r.config = dataSourceConfiguration(req, resp)
}

func (r *realmAssignmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Realm Assignment ID.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the Okta Realm Assignment.",
			},
			"priority": schema.Int32Attribute{
				Computed:    true,
				Description: "The Priority of the Realm Assignment. The lower the number, the higher the priority.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Defines whether the Realm Assignment is active or not. Valid values: `ACTIVE` and `INACTIVE`.",
				Validators: []validator.String{
					stringvalidator.OneOf("ACTIVE", "INACTIVE"),
				},
			},
			"profile_source_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the Profile Source.",
			},
			"realm_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the Realm asscociated with the Realm Assignment.",
			},
			"condition_expression": schema.StringAttribute{
				Computed:            true,
				Description:         "Condition expression for the Realm Assignment in Okta Expression Language. Example: `user.profile.role ==\"Manager\"` or `user.profile.state.contains(\"example\")`.",
				MarkdownDescription: "Condition expression for the Realm Assignment in Okta Expression Language. Example: `user.profile.role ==\"Manager\"` or `user.profile.state.contains(\"example\")`.",
			},
			"is_default": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether the realm assignment is the default.",
			},
		},
		Description: "Get a realm assignment from Okta.",
	}
}

func (r *realmAssignmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state realmAssignmentModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var selectedRealmAssignment *okta.RealmAssignment
	if state.ID.ValueString() != "" {
		realmAssignment, _, err := r.config.OktaIDaaSClient.OktaSDKClientV5().RealmAssignmentAPI.GetRealmAssignment(ctx, state.ID.ValueString()).Execute()
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error getting realm assignment with id: %v", state.ID.ValueString()), err.Error())
			return
		}

		selectedRealmAssignment = realmAssignment
	} else if state.Name.ValueString() != "" {
		realmAssignments, _, err := r.config.OktaIDaaSClient.OktaSDKClientV5().RealmAssignmentAPI.ListRealmAssignments(ctx).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Error listing realm assignments", err.Error())
			return
		}
		for _, realmAssignment := range realmAssignments {
			if *realmAssignment.Name == state.Name.ValueString() {
				selectedRealmAssignment = &realmAssignment
				break
			}
		}

		if selectedRealmAssignment == nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Realm assignment with name %s not found", state.Name), "Please check the name and try again.")
			return
		}
	} else {
		resp.Diagnostics.AddError("Error reading real assignments", "Either 'id' or 'name' must be specified.")
		return
	}

	resp.Diagnostics.Append(mapRealmAssignmentResourceToState(selectedRealmAssignment, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
