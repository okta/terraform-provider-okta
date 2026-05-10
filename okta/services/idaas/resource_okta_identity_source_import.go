package idaas

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	okta "github.com/okta/okta-sdk-golang/v6/okta"

	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource              = &identitySourceImportResource{}
	_ resource.ResourceWithConfigure = &identitySourceImportResource{}
)

type identitySourceImportResource struct {
	Config *config.Config
}

// --- Model types ---

type identitySourceImportModel struct {
	ID               types.String `tfsdk:"id"`
	IdentitySourceId types.String `tfsdk:"identity_source_id"`
	SessionId        types.String `tfsdk:"session_id"`
	SessionStatus    types.String `tfsdk:"session_status"`

	UpsertUsers            *identitySourceImportUpsertUsersModel            `tfsdk:"upsert_users"`
	UpsertGroups           *identitySourceImportUpsertGroupsModel           `tfsdk:"upsert_groups"`
	DeleteUsers            *identitySourceImportDeleteUsersModel            `tfsdk:"delete_users"`
	DeleteGroups           *identitySourceImportDeleteGroupsModel           `tfsdk:"delete_groups"`
	UpsertGroupMemberships *identitySourceImportUpsertGroupMembershipsModel `tfsdk:"upsert_group_memberships"`
	DeleteGroupMemberships *identitySourceImportDeleteGroupMembershipsModel `tfsdk:"delete_group_memberships"`
}

type identitySourceImportUpsertUsersModel struct {
	EntityType types.String                            `tfsdk:"entity_type"`
	Profiles   []identitySourceImportUpsertUserProfile `tfsdk:"profiles"`
}

type identitySourceImportUpsertUserProfile struct {
	ExternalId types.String                                `tfsdk:"external_id"`
	Profile    *identitySourceImportUpsertUserProfileInner `tfsdk:"profile"`
}

type identitySourceImportUpsertUserProfileInner struct {
	Email       types.String `tfsdk:"email"`
	FirstName   types.String `tfsdk:"first_name"`
	HomeAddress types.String `tfsdk:"home_address"`
	LastName    types.String `tfsdk:"last_name"`
	MobilePhone types.String `tfsdk:"mobile_phone"`
	SecondEmail types.String `tfsdk:"second_email"`
	UserName    types.String `tfsdk:"user_name"`
}

type identitySourceImportUpsertGroupsModel struct {
	Profiles []identitySourceImportUpsertGroupProfile `tfsdk:"profiles"`
}

type identitySourceImportUpsertGroupProfile struct {
	ExternalId   types.String                                 `tfsdk:"external_id"`
	GroupProfile *identitySourceImportUpsertGroupProfileInner `tfsdk:"group_profile"`
}

type identitySourceImportUpsertGroupProfileInner struct {
	Description types.String `tfsdk:"description"`
	DisplayName types.String `tfsdk:"display_name"`
}

type identitySourceImportDeleteUsersModel struct {
	EntityType types.String                            `tfsdk:"entity_type"`
	Profiles   []identitySourceImportDeleteUserProfile `tfsdk:"profiles"`
}

type identitySourceImportDeleteUserProfile struct {
	ExternalId types.String `tfsdk:"external_id"`
}

type identitySourceImportDeleteGroupsModel struct {
	ExternalIds types.List `tfsdk:"external_ids"`
}

type identitySourceImportUpsertGroupMembershipsModel struct {
	Memberships []identitySourceImportGroupMembership `tfsdk:"memberships"`
}

type identitySourceImportDeleteGroupMembershipsModel struct {
	Memberships []identitySourceImportGroupMembership `tfsdk:"memberships"`
}

type identitySourceImportGroupMembership struct {
	GroupExternalId   types.String `tfsdk:"group_external_id"`
	MemberExternalIds types.List   `tfsdk:"member_external_ids"`
}

func newIdentitySourceImportResource() resource.Resource {
	return &identitySourceImportResource{}
}

func (r *identitySourceImportResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_source_import"
}

func (r *identitySourceImportResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	cfg, ok := req.ProviderData.(*config.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *config.Config, got something else. Please report this issue to the provider developers.",
		)
		return
	}
	r.Config = cfg
}

func (r *identitySourceImportResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	membershipAttrs := map[string]schema.Attribute{
		"group_external_id": schema.StringAttribute{
			Description: "External ID of the group.",
			Optional:    true,
		},
		"member_external_ids": schema.ListAttribute{
			Description: "External IDs of the group members.",
			ElementType: types.StringType,
			Optional:    true,
		},
	}

	resp.Schema = schema.Schema{
		Description: "Runs a complete identity source import job: creates a session, uploads all staged data " +
			"(upsert/delete for users, groups, and group memberships), then triggers the import. " +
			"All upload blocks are optional; at least one should be provided.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The session ID of the triggered import job.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identity_source_id": schema.StringAttribute{
				Description: "ID of the identity source.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"session_id": schema.StringAttribute{
				Description: "The session ID created for this import job.",
				Computed:    true,
			},
			"session_status": schema.StringAttribute{
				Description: "The status of the import session after triggering.",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"upsert_users": schema.SingleNestedBlock{
				Description: "Users to create or update in Okta.",
				Attributes: map[string]schema.Attribute{
					"entity_type": schema.StringAttribute{
						Description: "Entity type. Currently only `USERS` is supported.",
						Optional:    true,
					},
				},
				Blocks: map[string]schema.Block{
					"profiles": schema.ListNestedBlock{
						Description: "User profiles to upsert.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"external_id": schema.StringAttribute{
									Description: "External ID of the user.",
									Optional:    true,
								},
							},
							Blocks: map[string]schema.Block{
								"profile": schema.SingleNestedBlock{
									Description: "User profile attributes.",
									Attributes: map[string]schema.Attribute{
										"email":        schema.StringAttribute{Description: "Email address of the user.", Optional: true},
										"first_name":   schema.StringAttribute{Description: "First name of the user.", Optional: true},
										"home_address": schema.StringAttribute{Description: "Home address of the user.", Optional: true},
										"last_name":    schema.StringAttribute{Description: "Last name of the user.", Optional: true},
										"mobile_phone": schema.StringAttribute{Description: "Mobile phone number of the user.", Optional: true},
										"second_email": schema.StringAttribute{Description: "Alternative email address of the user.", Optional: true},
										"user_name":    schema.StringAttribute{Description: "Username of the user.", Optional: true},
									},
								},
							},
						},
					},
				},
			},
			"upsert_groups": schema.SingleNestedBlock{
				Description: "Groups to create or update in Okta.",
				Blocks: map[string]schema.Block{
					"profiles": schema.ListNestedBlock{
						Description: "Group profiles to upsert.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"external_id": schema.StringAttribute{
									Description: "External ID of the group.",
									Optional:    true,
								},
							},
							Blocks: map[string]schema.Block{
								"group_profile": schema.SingleNestedBlock{
									Description: "Group profile attributes.",
									Attributes: map[string]schema.Attribute{
										"description":  schema.StringAttribute{Description: "Description of the group.", Optional: true},
										"display_name": schema.StringAttribute{Description: "Display name of the group.", Optional: true},
									},
								},
							},
						},
					},
				},
			},
			"delete_users": schema.SingleNestedBlock{
				Description: "Users to delete from Okta.",
				Attributes: map[string]schema.Attribute{
					"entity_type": schema.StringAttribute{
						Description: "Entity type. Currently only `USERS` is supported.",
						Optional:    true,
					},
				},
				Blocks: map[string]schema.Block{
					"profiles": schema.ListNestedBlock{
						Description: "User profiles to delete (by external ID).",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"external_id": schema.StringAttribute{
									Description: "External ID of the user to delete.",
									Optional:    true,
								},
							},
						},
					},
				},
			},
			"delete_groups": schema.SingleNestedBlock{
				Description: "Groups to delete from Okta.",
				Attributes: map[string]schema.Attribute{
					"external_ids": schema.ListAttribute{
						Description: "External IDs of groups to delete.",
						ElementType: types.StringType,
						Optional:    true,
					},
				},
			},
			"upsert_group_memberships": schema.SingleNestedBlock{
				Description: "Group memberships to create or update in Okta.",
				Blocks: map[string]schema.Block{
					"memberships": schema.ListNestedBlock{
						Description: "Group memberships to upsert.",
						NestedObject: schema.NestedBlockObject{
							Attributes: membershipAttrs,
						},
					},
				},
			},
			"delete_group_memberships": schema.SingleNestedBlock{
				Description: "Group memberships to delete in Okta.",
				Blocks: map[string]schema.Block{
					"memberships": schema.ListNestedBlock{
						Description: "Group memberships to delete.",
						NestedObject: schema.NestedBlockObject{
							Attributes: membershipAttrs,
						},
					},
				},
			},
		},
	}
}

// runImport creates a new session, uploads all staged data, and triggers the import.
// Returns the session ID, session status, and any diagnostics.
func (r *identitySourceImportResource) runImport(ctx context.Context, plan *identitySourceImportModel) (string, string, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	identitySourceId := plan.IdentitySourceId.ValueString()

	// 1. Create session
	session, _, err := client.IdentitySourceAPI.CreateIdentitySourceSession(ctx, identitySourceId).Execute()
	if err != nil {
		diags.AddError("Error creating identity source session", err.Error())
		return "", "", diags
	}
	sessionId := session.GetId()

	// cleanupSession deletes the session if any subsequent step fails, so a re-apply
	// isn't blocked by the "only one active session per 5 minutes" rate limit.
	cleanupSession := func() {
		_, _ = client.IdentitySourceAPI.DeleteIdentitySourceSession(ctx, identitySourceId, sessionId).Execute()
	}

	// 2. Upsert users
	if plan.UpsertUsers != nil {
		body := okta.NewBulkUpsertRequestBodyWithDefaults()
		body.SetEntityType(plan.UpsertUsers.EntityType.ValueString())
		if len(plan.UpsertUsers.Profiles) > 0 {
			var profiles []okta.BulkUpsertRequestBodyProfilesInner
			for _, item := range plan.UpsertUsers.Profiles {
				p := okta.NewBulkUpsertRequestBodyProfilesInnerWithDefaults()
				if !item.ExternalId.IsNull() && !item.ExternalId.IsUnknown() {
					p.SetExternalId(item.ExternalId.ValueString())
				}
				if item.Profile != nil {
					up := okta.NewIdentitySourceUserProfileForUpsertWithDefaults()
					if !item.Profile.Email.IsNull() && !item.Profile.Email.IsUnknown() {
						up.SetEmail(item.Profile.Email.ValueString())
					}
					if !item.Profile.FirstName.IsNull() && !item.Profile.FirstName.IsUnknown() {
						up.SetFirstName(item.Profile.FirstName.ValueString())
					}
					if !item.Profile.HomeAddress.IsNull() && !item.Profile.HomeAddress.IsUnknown() {
						up.SetHomeAddress(item.Profile.HomeAddress.ValueString())
					}
					if !item.Profile.LastName.IsNull() && !item.Profile.LastName.IsUnknown() {
						up.SetLastName(item.Profile.LastName.ValueString())
					}
					if !item.Profile.MobilePhone.IsNull() && !item.Profile.MobilePhone.IsUnknown() {
						up.SetMobilePhone(item.Profile.MobilePhone.ValueString())
					}
					if !item.Profile.SecondEmail.IsNull() && !item.Profile.SecondEmail.IsUnknown() {
						up.SetSecondEmail(item.Profile.SecondEmail.ValueString())
					}
					if !item.Profile.UserName.IsNull() && !item.Profile.UserName.IsUnknown() {
						up.SetUserName(item.Profile.UserName.ValueString())
					}
					p.SetProfile(*up)
				}
				profiles = append(profiles, *p)
			}
			body.SetProfiles(profiles)
		}
		req := client.IdentitySourceAPI.UploadIdentitySourceDataForUpsert(ctx, identitySourceId, sessionId).BulkUpsertRequestBody(*body)
		if _, err := req.Execute(); err != nil {
			cleanupSession()
			diags.AddError("Error uploading upsert users", err.Error())
			return "", "", diags
		}
	}

	// 3. Upsert groups
	if plan.UpsertGroups != nil {
		body := okta.NewBulkGroupUpsertRequestBody()
		if len(plan.UpsertGroups.Profiles) > 0 {
			var profiles []okta.BulkGroupUpsertRequestBodyProfilesInner
			for _, item := range plan.UpsertGroups.Profiles {
				p := okta.NewBulkGroupUpsertRequestBodyProfilesInnerWithDefaults()
				if !item.ExternalId.IsNull() && !item.ExternalId.IsUnknown() {
					p.SetExternalId(item.ExternalId.ValueString())
				}
				if item.GroupProfile != nil {
					gp := okta.NewIdentitySourceGroupProfileForUpsertWithDefaults()
					if !item.GroupProfile.Description.IsNull() && !item.GroupProfile.Description.IsUnknown() {
						gp.SetDescription(item.GroupProfile.Description.ValueString())
					}
					if !item.GroupProfile.DisplayName.IsNull() && !item.GroupProfile.DisplayName.IsUnknown() {
						gp.SetDisplayName(item.GroupProfile.DisplayName.ValueString())
					}
					p.SetProfile(*gp)
				}
				profiles = append(profiles, *p)
			}
			body.SetProfiles(profiles)
		}
		req := client.IdentitySourceAPI.UploadIdentitySourceGroupsForUpsert(ctx, identitySourceId, sessionId).BulkGroupUpsertRequestBody(*body)
		if _, err := req.Execute(); err != nil {
			cleanupSession()
			diags.AddError("Error uploading upsert groups", err.Error())
			return "", "", diags
		}
	}

	// 4. Delete users
	if plan.DeleteUsers != nil {
		body := okta.NewBulkDeleteRequestBodyWithDefaults()
		body.SetEntityType(plan.DeleteUsers.EntityType.ValueString())
		if len(plan.DeleteUsers.Profiles) > 0 {
			var profiles []okta.IdentitySourceUserProfileForDelete
			for _, item := range plan.DeleteUsers.Profiles {
				p := okta.NewIdentitySourceUserProfileForDeleteWithDefaults()
				if !item.ExternalId.IsNull() && !item.ExternalId.IsUnknown() {
					p.SetExternalId(item.ExternalId.ValueString())
				}
				profiles = append(profiles, *p)
			}
			body.SetProfiles(profiles)
		}
		req := client.IdentitySourceAPI.UploadIdentitySourceDataForDelete(ctx, identitySourceId, sessionId).BulkDeleteRequestBody(*body)
		if _, err := req.Execute(); err != nil {
			cleanupSession()
			diags.AddError("Error uploading delete users", err.Error())
			return "", "", diags
		}
	}

	// 5. Delete groups
	if plan.DeleteGroups != nil {
		body := okta.NewBulkGroupDeleteRequestBodyWithDefaults()
		if !plan.DeleteGroups.ExternalIds.IsNull() && !plan.DeleteGroups.ExternalIds.IsUnknown() {
			var ids []string
			for _, elem := range plan.DeleteGroups.ExternalIds.Elements() {
				if sv, ok := elem.(types.String); ok {
					ids = append(ids, sv.ValueString())
				}
			}
			body.SetExternalIds(ids)
		}
		req := client.IdentitySourceAPI.UploadIdentitySourceGroupsDataForDelete(ctx, identitySourceId, sessionId).BulkGroupDeleteRequestBody(*body)
		if _, err := req.Execute(); err != nil {
			cleanupSession()
			diags.AddError("Error uploading delete groups", err.Error())
			return "", "", diags
		}
	}

	// 6. Upsert group memberships
	if plan.UpsertGroupMemberships != nil {
		body := okta.NewBulkGroupMembershipsUpsertRequestBodyWithDefaults()
		if len(plan.UpsertGroupMemberships.Memberships) > 0 {
			var memberships []okta.IdentitySourceGroupMembershipsUpsertProfileInner
			for _, item := range plan.UpsertGroupMemberships.Memberships {
				m := okta.NewIdentitySourceGroupMembershipsUpsertProfileInnerWithDefaults()
				if !item.GroupExternalId.IsNull() && !item.GroupExternalId.IsUnknown() {
					m.SetGroupExternalId(item.GroupExternalId.ValueString())
				}
				if !item.MemberExternalIds.IsNull() && !item.MemberExternalIds.IsUnknown() {
					var ids []string
					for _, elem := range item.MemberExternalIds.Elements() {
						if sv, ok := elem.(types.String); ok {
							ids = append(ids, sv.ValueString())
						}
					}
					m.SetMemberExternalIds(ids)
				}
				memberships = append(memberships, *m)
			}
			body.SetMemberships(memberships)
		}
		req := client.IdentitySourceAPI.UploadIdentitySourceGroupMembershipsForUpsert(ctx, identitySourceId, sessionId).BulkGroupMembershipsUpsertRequestBody(*body)
		if _, err := req.Execute(); err != nil {
			cleanupSession()
			diags.AddError("Error uploading upsert group memberships", err.Error())
			return "", "", diags
		}
	}

	// 7. Delete group memberships
	if plan.DeleteGroupMemberships != nil {
		body := okta.NewBulkGroupMembershipsDeleteRequestBodyWithDefaults()
		if len(plan.DeleteGroupMemberships.Memberships) > 0 {
			var memberships []okta.IdentitySourceGroupMembershipsDeleteProfileInner
			for _, item := range plan.DeleteGroupMemberships.Memberships {
				m := okta.NewIdentitySourceGroupMembershipsDeleteProfileInnerWithDefaults()
				if !item.GroupExternalId.IsNull() && !item.GroupExternalId.IsUnknown() {
					m.SetGroupExternalId(item.GroupExternalId.ValueString())
				}
				if !item.MemberExternalIds.IsNull() && !item.MemberExternalIds.IsUnknown() {
					var ids []string
					for _, elem := range item.MemberExternalIds.Elements() {
						if sv, ok := elem.(types.String); ok {
							ids = append(ids, sv.ValueString())
						}
					}
					m.SetMemberExternalIds(ids)
				}
				memberships = append(memberships, *m)
			}
			body.SetMemberships(memberships)
		}
		req := client.IdentitySourceAPI.UploadIdentitySourceGroupMembershipsForDelete(ctx, identitySourceId, sessionId).BulkGroupMembershipsDeleteRequestBody(*body)
		if _, err := req.Execute(); err != nil {
			cleanupSession()
			diags.AddError("Error uploading delete group memberships", err.Error())
			return "", "", diags
		}
	}

	// 8. Trigger import
	result, _, err := client.IdentitySourceAPI.StartImportFromIdentitySource(ctx, identitySourceId, sessionId).Execute()
	if err != nil {
		cleanupSession()
		diags.AddError("Error triggering identity source import", err.Error())
		return "", "", diags
	}

	return sessionId, result.GetStatus(), diags
}

func (r *identitySourceImportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan identitySourceImportModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sessionId, status, diags := r.runImport(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(sessionId)
	plan.SessionId = types.StringValue(sessionId)
	plan.SessionStatus = types.StringValue(status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *identitySourceImportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state identitySourceImportModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Config.OktaIDaaSClient.OktaSDKClientV6()
	result, httpResp, err := client.IdentitySourceAPI.GetIdentitySourceSession(ctx, state.IdentitySourceId.ValueString(), state.SessionId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading identity source session", err.Error())
		return
	}

	state.SessionStatus = types.StringValue(result.GetStatus())
	// Keep session_id in sync in case it was set via import
	state.SessionId = types.StringValue(result.GetId())
	_ = result.GetLastUpdated().Format(time.RFC3339) // ensure time import is used

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *identitySourceImportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan identitySourceImportModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the Terraform resource ID from existing state
	var state identitySourceImportModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sessionId, status, diags := r.runImport(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID
	plan.SessionId = types.StringValue(sessionId)
	plan.SessionStatus = types.StringValue(status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *identitySourceImportResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"Removing this resource from configuration does not undo the import that was triggered in Okta.",
	)
}

// Ensure diag is used
var _ diag.Diagnostics
