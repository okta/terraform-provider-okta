package okta

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func dataSourceApp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppRead,
		Schema: buildSchema(skipUsersAndGroupsSchema, map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"label", "label_prefix"},
			},
			"label": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label_prefix"},
			},
			"label_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id", "label"},
			},
			"active_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Search only ACTIVE applications.",
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
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
		}),
	}
}

func dataSourceAppRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	filters, err := getAppFilters(d)
	if err != nil {
		return diag.Errorf("invalid app filters: %v", err)
	}
	var app *sdk.Application
	if filters.ID != "" {
		respApp, _, err := getOktaClientFromMetadata(m).Application.GetApplication(ctx, filters.ID, sdk.NewApplication(), nil)
		if err != nil {
			return diag.Errorf("failed get app by ID: %v", err)
		}
		app = respApp.(*sdk.Application)
	} else {
		appList, err := listApps(ctx, getOktaClientFromMetadata(m), filters, 1)
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
				logger(m).Info("found multiple applications with the criteria supplied, using the first one, sorted by creation date")
			}
			app = appList[0]
		}
	}
	err = setAppUsersIDsAndGroupsIDs(ctx, d, getOktaClientFromMetadata(m), app.Id)
	if err != nil {
		return diag.Errorf("failed to list app's groups and users: %v", err)
	}
	d.SetId(app.Id)
	_ = d.Set("label", app.Label)
	_ = d.Set("name", app.Name)
	_ = d.Set("status", app.Status)
	p, _ := json.Marshal(app.Links)
	_ = d.Set("links", string(p))
	return nil
}
