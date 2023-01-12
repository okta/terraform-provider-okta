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
		Schema: buildSchema(skipUsersAndGroupsSchema, map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This label displays under the app on your home page",
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
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable web application to override ReplyTo URL with reply param",
			},
			"name_id_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name ID Format",
			},
			"audience_restriction": {
				Type:        schema.TypeString,
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
			"auto_submit_toolbar": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Display auto submit toolbar",
			},
			"hide_ios": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Do not display application icon on mobile app",
			},
			"hide_web": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Do not display application icon to users",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     true,
				Description: "Activation status of the application",
			},
		}),
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

	Visibility := d.Get("visibility").(bool)
	_ = d.Set("visibility", &Visibility)

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
	app := okta.NewWsFederationApplication()
	app.Label = d.Get("label").(string)
	WReplyOverride := d.Get("reply_override").(bool)
	app.Settings = &okta.WsFederationApplicationSettings{
		App: &okta.WsFederationApplicationSettingsApplication{
			AttributeStatements:  d.Get("attribute_statements").(string),
			AudienceRestriction:  d.Get("audience_restriction").(string),
			AuthnContextClassRef: d.Get("authn_context_class_ref").(string),
			GroupFilter:          d.Get("group_filter").(string),
			GroupName:            d.Get("group_name").(string),
			GroupValueFormat:     d.Get("group_value_format").(string),
			NameIDFormat:         d.Get("name_id_format").(string),
			Realm:                d.Get("realm").(string),
			SiteURL:              d.Get("site_url").(string),
			UsernameAttribute:    d.Get("username_attribute").(string),
			WReplyOverride:       &WReplyOverride,
			WReplyURL:            d.Get("reply_url").(string),
		},
		Notes: buildAppNotes(d),
	}
	app.Visibility = buildAppVisibility(d)
	return app
}
