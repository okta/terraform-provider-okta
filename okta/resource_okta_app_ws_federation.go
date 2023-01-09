package okta

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func resourceAppWsFed() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppWsFedCreate,
		ReadContext:   resourceAppWsFedRead,
		UpdateContext: resourceAppWsFedUpdate,
		DeleteContext: resourceAppWsFedDelete,
		Importer: &schema.ResourceImporter{
			StateContext: appImporter,
		},
		// Schema: buildAppWsFedSchema(map[string]*schema.Schema{
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A display-friendly label for this app",
			},
			"site_url": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "",
				ValidateDiagFunc: stringIsURL(validURLSchemes...),
			},
			"realm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"reply_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"reply_override": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Login URL",
			},
			"name_id_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"audience_restriction": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "",
			},
			"authn_context_class_ref": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"group_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"group_value_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"username_attribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"attribute_statements": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"visibility": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Should the application icon be visible to users?",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Read:   schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
		},
	}
}

func resourceAppWsFedCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	app := buildAppWsFed(d)
	activate := d.Get("status").(string) == statusActive
	params := &query.Params{Activate: &activate}
	_, _, err := client.Application.CreateApplication(ctx, app, params)
	if err != nil {
		return diag.Errorf("failed to create WSFed application: %v", err)
	}
	d.SetId(app.Id)
	err = handleAppGroupsAndUsers(ctx, app.Id, d, m)
	if err != nil {
		return diag.Errorf("failed to handle groups and users for WSFed application: %v", err)
	}
	err = handleAppLogo(ctx, d, m, app.Id, app.Links)
	if err != nil {
		return diag.Errorf("failed to upload logo for WSFed application: %v", err)
	}
	return resourceAppWsFedRead(ctx, d, m)
}

func resourceAppWsFedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	app := okta.NewWsFederationApplication()
	err := fetchApp(ctx, d, m, app)
	if err != nil {
		return diag.Errorf("failed to get WsFed application: %v", err)
	}
	if app.Id == "" {
		d.SetId("")
		return nil
	}
	_ = d.Set("label", app.Label)
	_ = d.Set("site_url", app.Settings.App.SiteURL)
	_ = d.Set("realm", app.Settings.App.Realm)
	_ = d.Set("reply_url", app.Settings.App.WReplyURL)
	_ = d.Set("reply_override", app.Settings.App.WReplyOverride)
	_ = d.Set("name_id_format", app.Settings.App.NameIDFormat)
	_ = d.Set("audience_restriction", app.Settings.App.AudienceRestriction)
	_ = d.Set("authn_context_class_ref", app.Settings.App.AuthnContextClassRef)
	_ = d.Set("group_filter", app.Settings.App.GroupFilter)
	_ = d.Set("group_name", app.Settings.App.GroupName)
	_ = d.Set("group_value_format", app.Settings.App.GroupValueFormat)
	_ = d.Set("username_attribute", app.Settings.App.UsernameAttribute)
	_ = d.Set("attribute_statements", app.Settings.App.AttributeStatements)
	_ = d.Set("visibility", app.Visibility)

	_ = d.Set("", linksValue(app.Links, "logo", "href"))
	appRead(d, app.Name, app.Status, app.SignOnMode, app.Label, app.Accessibility, app.Visibility, app.Settings.Notes)
	err = syncGroupsAndUsers(ctx, app.Id, d, m)
	if err != nil {
		return diag.Errorf("failed to sync groups and users for WsFed application: %v", err)
	}
	return nil
}

func resourceAppWsFedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	app := buildAppWsFed(d)
	_, _, err := client.Application.UpdateApplication(ctx, d.Id(), app)
	if err != nil {
		return diag.Errorf("failed to update WSFed application: %v", err)
	}
	err = setAppStatus(ctx, d, client, app.Status)
	if err != nil {
		return diag.Errorf("failed to set WSFed application status: %v", err)
	}
	err = handleAppGroupsAndUsers(ctx, app.Id, d, m)
	if err != nil {
		return diag.Errorf("failed to handle groups and users for WSFed application: %v", err)
	}
	if d.HasChange("logo") {
		err = handleAppLogo(ctx, d, m, app.Id, app.Links)
		if err != nil {
			o, _ := d.GetChange("logo")
			_ = d.Set("logo", o)
			return diag.Errorf("failed to upload logo for WSFed application: %v", err)
		}
	}
	return resourceAppWsFedRead(ctx, d, m)
}

func resourceAppWsFedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := deleteApplication(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to delete WSFed application: %v", err)
	}
	return nil
}

func buildAppWsFed(d *schema.ResourceData) *okta.WsFederationApplication {
	// Abstracts away name and SignOnMode which are constant for this app type.
	app := okta.NewWsFederationApplication()
	app.Label = d.Get("label").(string)
	name := d.Get("preconfigured_app").(string)
	if name != "" {
		app.Name = name
		app.SignOnMode = "WS_FEDERATION" // in case pre-configured app has more than one sign-on modes
	}
	app.Settings = &okta.WsFederationApplicationSettings{
		App: &okta.WsFederationApplicationSettingsApplication{
			AttributeStatements:  d.Get("").(string),
			AudienceRestriction:  d.Get("").(string),
			AuthnContextClassRef: d.Get("").(string),
			GroupFilter:          d.Get("").(string),
			GroupName:            d.Get("").(string),
			GroupValueFormat:     d.Get("").(string),
			NameIDFormat:         d.Get("").(string),
			Realm:                d.Get("realm").(string),
			SiteURL:              d.Get("site_url").(string),
			UsernameAttribute:    d.Get("").(string),
			WReplyOverride:       d.Get("").(*bool),
			WReplyURL:            d.Get("").(string),
		},
		Notes: buildAppNotes(d),
	}
	app.Visibility = buildAppVisibility(d)
	app.Accessibility = buildAppAccessibility(d)
	app.Credentials = &okta.ApplicationCredentials{
		UserNameTemplate: buildUserNameTemplate(d),
	}

	return app
}
