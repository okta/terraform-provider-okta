package idaas

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

func newFeaturesDataSource() datasource.DataSource {
	return &FeaturesDataSource{}
}

type FeaturesDataSource struct {
	*config.Config
}

type FeaturesDataSourceModel struct {
	ID        types.String        `tfsdk:"id"`
	Label     types.String        `tfsdk:"label"`
	Substring types.String        `tfsdk:"substring"`
	Features  []OktaFeaturesModel `tfsdk:"features"`
}

type OktaFeaturesModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Status      types.String `tfsdk:"status"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Stage       types.Object `tfsdk:"stage"`
}

type OktaFeaturesStageModel struct {
	State types.String `tfsdk:"state"`
	Value types.String `tfsdk:"value"`
}

func (d *FeaturesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_features"
}

func (d *FeaturesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get a list of features from Okta.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Generated ID",
				Computed:    true,
			},
			"label": schema.StringAttribute{
				Optional:    true,
				Description: "Searches for features whose label or name property matches this value exactly. Case sensitive",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("substring"),
					}...),
				},
			},
			"substring": schema.StringAttribute{
				Optional:    true,
				Description: "Searches for features whose label or name property substring match this value. Case sensitive",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("label"),
					}...),
				},
			},
			"features": schema.ListAttribute{
				Description: "The list of features that match the search criteria.",
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":          types.StringType,
						"name":        types.StringType,
						"status":      types.StringType,
						"type":        types.StringType,
						"description": types.StringType,
						"stage": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"state": types.StringType,
								"value": types.StringType,
							},
						},
					},
				},
			},
		},
	}
}

func (d *FeaturesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *FeaturesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state FeaturesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	featureList, _, err := d.OktaIDaaSClient.OktaSDKClientV5().FeatureAPI.ListFeatures(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Unable to list Okta Features", fmt.Sprintf("Error retrieving features: %s", err.Error()))
		return
	}
	state.ID = types.StringValue(uuid.New().String())
	for _, feature := range featureList {
		if !state.Label.IsNull() && feature.GetName() != state.Label.ValueString() {
			continue
		}
		if !state.Substring.IsNull() && !strings.Contains(feature.GetName(), state.Substring.ValueString()) {
			continue
		}
		featureStageValue := map[string]attr.Value{
			"state": types.StringPointerValue(feature.GetStage().State),
			"value": types.StringPointerValue(feature.GetStage().Value),
		}
		featureStageTypes := map[string]attr.Type{
			"state": types.StringType,
			"value": types.StringType,
		}
		featureStage, _ := types.ObjectValue(featureStageTypes, featureStageValue)

		state.Features = append(state.Features, OktaFeaturesModel{
			ID:          types.StringPointerValue(feature.Id),
			Name:        types.StringPointerValue(feature.Name),
			Status:      types.StringPointerValue(feature.Status),
			Type:        types.StringPointerValue(feature.Type),
			Description: types.StringPointerValue(feature.Description),
			Stage:       featureStage,
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
