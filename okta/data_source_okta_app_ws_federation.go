package okta

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func dataSourceAppWsFed() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppWsFedRead,
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
			"site_url": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Launch URL for the Web Application",
				ValidateDiagFunc: stringIsURL(validURLSchemes...),
			},
			"realm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The trust realm for the Web Application",
			},
			"reply_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ReplyTo URL to which responses are directed",
			},
			"reply_override": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable web application to override ReplyTo URL with reply param",
			},
			"name_id_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name ID Format",
			},
			"audience_restriction": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The assertion containing a bearer subject confirmation MUST contain an Audience Restriction including the service provider's unique identifier as an Audience",
			},
			"authn_context_class_ref": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the Authentication Context for the issued SAML Assertion",
			},
			"group_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An expression that will be used to filter groups. If the Okta group name matches the expression, the group name will be included in the SAML Assertion Attribute Statement",
			},
			"group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the SAML attribute name for a user's group memberships",
			},
			"group_value_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the SAML assertion attribute value for filtered groups",
			},
			"username_attribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies additional username attribute statements to include in the SAML Assertion",
			},
			"attribute_statements": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Defines custom SAML attribute statements",
			},
			"visibility": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Application icon visibility to users",
			},
		}),
	}
}

func dataSourceAppWsFedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	filters, err := getAppFilters(d)
	if err != nil {
		return diag.Errorf("invalid WsFed app filters: %v", err)
	}
	var app *okta.WsFederationApplication
	if filters.ID != "" {
		respApp, _, err := getOktaClientFromMetadata(m).Application.GetApplication(ctx, filters.ID, okta.NewWsFederationApplication(), nil)
		if err != nil {
			return diag.Errorf("failed get app by ID: %v", err)
		}
		app = respApp.(*okta.WsFederationApplication)
	} else {
		re := getOktaClientFromMetadata(m).GetRequestExecutor()
		qp := &query.Params{Limit: 1, Filter: filters.Status, Q: filters.getQ()}
		req, err := re.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/apps%s", qp.String()), nil)
		if err != nil {
			return diag.Errorf("failed to list WsFed apps: %v", err)
		}
		var appList []*okta.WsFederationApplication
		_, err = re.Do(ctx, req, &appList)
		if err != nil {
			return diag.Errorf("failed to list WsFed apps: %v", err)
		}
		if len(appList) < 1 {
			return diag.Errorf("no WsFed application found with provided filter: %s", filters)
		}
		if filters.Label != "" && appList[0].Label != filters.Label {
			return diag.Errorf("no WsFed application found with the provided label: %s", filters.Label)
		}
		logger(m).Info("found multiple WsFed applications with the criteria supplied, using the first one, sorted by creation date")
		app = appList[0]
	}
	err = setAppUsersIDsAndGroupsIDs(ctx, d, getOktaClientFromMetadata(m), app.Id)
	if err != nil {
		return diag.Errorf("failed to list WsFed app groups and users: %v", err)
	}
	d.SetId(app.Id)
	_ = d.Set("label", app.Label)
	_ = d.Set("name", app.Name)
	_ = d.Set("status", app.Status)
	p, _ := json.Marshal(app.Links)
	_ = d.Set("links", string(p))
	return nil
}
