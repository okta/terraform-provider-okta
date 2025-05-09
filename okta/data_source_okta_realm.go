package okta

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-sdk-golang/v5/okta"
)

type realmDataSource struct {
	config *Config
}

func NewRealmDataSource() datasource.DataSource {
	return &realmDataSource{}
}

func (r *realmDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_realm"
}

func (r *realmDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	r.config = dataSourceConfiguration(req, resp)
}

func (r *realmDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The id of the Okta Realm.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("name"),
					}...),
				},
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The name of the Okta Realm.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("id"),
					}...),
				},
			},
			"realm_type": schema.StringAttribute{
				Optional:    true,
				Description: "The realm type. Valid values: `PARTNER` and `DEFAULT`",
			},
			"is_default": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether the realm is the default realm.",
			},
		},
		Description: "Get a realm from Okta.",
	}
}

func (r *realmDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state realmModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var selectedRealm *okta.Realm
	if state.ID.ValueString() != "" {
		realm, response, err := r.config.oktaSDKClientV5.RealmAPI.GetRealm(ctx, state.ID.ValueString()).Execute()
		if err != nil {
			body, ioErr := io.ReadAll(response.Body)
			defer response.Body.Close()
			if ioErr != nil {
				resp.Diagnostics.AddError(err.Error(), "failed to read response body")
				return
			}
			resp.Diagnostics.AddError("failed to read realm:"+err.Error(), string(body))
			return
		}
		selectedRealm = realm
	} else if state.Name.ValueString() != "" {
		searchString := fmt.Sprintf(`profile.name eq "%s"`, state.Name.ValueString())

		const retryCount = 3
		for range retryCount {
			realms, response, err := r.config.oktaSDKClientV5.RealmAPI.ListRealms(ctx).Search(searchString).Execute()
			if err != nil {
				body, ioErr := io.ReadAll(response.Body)
				defer response.Body.Close()
				if ioErr != nil {
					resp.Diagnostics.AddError(err.Error(), "failed to read response body")
					return
				}
				resp.Diagnostics.AddError("failed to list realms:"+err.Error(), string(body))
				return
			}
			if len(realms) == 0 {
				resp.Diagnostics.AddWarning("Realm not found", fmt.Sprintf("No realm found with name %s. Retrying...", state.Name.ValueString()))
				time.Sleep(time.Second)
				continue
			}

			if len(realms) != 1 {
				resp.Diagnostics.AddError("Multiple realms found", fmt.Sprintf("Found %d realms with name %s. Please specify a unique name.", len(realms), state.Name.ValueString()))
				return
			}

			selectedRealm = &realms[0]
			break
		}
		if selectedRealm == nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Realm with name %s not found", state.Name), "Please check the name and try again.")
			return
		}
	} else {
		resp.Diagnostics.AddError("Error reading realm", "Either 'id' or 'name' must be specified.")
		return
	}

	state.ID = types.StringPointerValue(selectedRealm.Id)
	state.Name = types.StringValue(selectedRealm.Profile.Name)
	state.RealmType = types.StringPointerValue(selectedRealm.Profile.RealmType)
	state.IsDefault = types.BoolPointerValue(selectedRealm.IsDefault)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
