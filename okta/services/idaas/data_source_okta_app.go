package idaas

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func dataSourceApp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppRead,
		Schema: utils.BuildSchema(skipUsersAndGroupsSchema, map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"label", "label_prefix"},
				Description:   "Id of application to retrieve, conflicts with label and label_prefix.",
			},
			"label": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label_prefix"},
				Description: `The label of the app to retrieve, conflicts with
				label_prefix and id. Label uses the ?q=<label> query parameter exposed by
				Okta's List Apps API. The API will search both name and label using that
				query. Therefore similarly named and labeled apps may be returned in the query
				and have the unintended result of associating the wrong app with this data
				source. See:
				https://developer.okta.com/docs/reference/api/apps/#list-applications`,
			},
			"label_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label"},
				Description: `Label prefix of the app to retrieve, conflicts with label and id. This will tell the
				provider to do a starts with query as opposed to an equals query.`,
			},
			"active_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Search only ACTIVE applications.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of application.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of application.",
			},
			"links": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Discoverable resources related to the app",
			},
			"groups": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Groups associated with the application",
				Deprecated:  "The `groups` field is now deprecated for the data source `okta_app`, please replace all uses of this with: `okta_app_group_assignments`",
			},
			"users": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Users associated with the application",
				Deprecated:  "The `users` field is now deprecated for the data source `okta_app`, please replace all uses of this with: `okta_app_user_assignments`",
			},
			"authentication_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the app's authentication policy",
			},
		}),
		Description: "Get an application of any kind from Okta.",
	}
}

func dataSourceAppRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	filters, err := getAppFilters(d)
	if err != nil {
		return diag.Errorf("invalid app filters: %v", err)
	}
	var app *sdk.Application
	if filters.ID != "" {
		respApp, _, err := getOktaClientFromMetadata(meta).Application.GetApplication(ctx, filters.ID, sdk.NewApplication(), nil)
		if err != nil {
			return diag.Errorf("failed get app by ID: %v", err)
		}
		app = respApp.(*sdk.Application)
	} else {
		appList, err := ListAppsV2(ctx, getOktaClientFromMetadata(meta), filters, 1)
		if err != nil {
			return diag.Errorf("failed to list apps: %v", err)
		}
		if len(appList) < 1 {
			return diag.Errorf("no application found with the provided filter: %+v", filters)
		}

		// Okta API for list apps uses a starts with query on label and name
		// which could result in multiple matches on the data source's "label"
		// property.  We need to further inspect for an exact label match.
		if filters.Label != "" {
			// guard on nils, an app is always set
			app = appList[0]
			for i, _app := range appList {
				if _app.Label == filters.Label {
					app = appList[i]
					break
				}
			}
		} else {
			if len(appList) > 1 {
				logger(meta).Info("found multiple applications with the criteria supplied, using the first one, sorted by creation date")
			}
			app = appList[0]
		}
	}
	d.SetId(app.Id)
	_ = d.Set("label", app.Label)
	_ = d.Set("name", app.Name)
	_ = d.Set("status", app.Status)
	p, _ := json.Marshal(app.Links)
	_ = d.Set("links", string(p))
	setAuthenticationPolicy(ctx, meta, d, app.Links)
	return nil
}
