package okta

import (
	"runtime"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

const (
	browserPlugin       = "BROWSER_PLUGIN"
	autoLogin           = "AUTO_LOGIN"
	securePasswordStore = "SECURE_PASSWORD_STORE"
	swaPlugin 			= "PLUGIN"
	swaThreeField 		= "THREE_FIELD"
	swaNoPlugin 		= "NO_PLUGIN"
	swaCustom 			= "CUSTOM"
)

// Logical SWA Types:
// Plugin: BROWSER_PLUGIN
// 3 Field: BROWSER_PLUGIN
// No Plugin: SECURE_PASSWORD_STORE
// Custom SWA: AUTO_LOGIN

// Credentials schema:
// EDIT_USERNAME_AND_PASSWORD: user editable username and password
// ADMIN_SETS_CREDENTIALS: Admin sets username and password
// EDIT_PASSWORD_ONLY: User editable password
// EXTERNAL_PASSWORD_SYNC: scheme for a SWA application with a username template
// SHARED_USERNAME_AND_PASSWORD: scheme for a SWA application with a username and password
func resourceSwaApp() *schema.Resource {
	return &schema.Resource{
		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			return nil
		},
		Create: resourceSwaAppCreate,
		Read:   resourceSwaAppRead,
		Update: resourceSwaAppUpdate,
		Delete: resourceSwaAppDelete,
		Exists: resourceSwaAppExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of preexisting SWA application.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice([]string{swaCustom, swaPlugin, swaNoPlugin, swaThreeField}),
				ForceNew: true,
				Description: "SWA App type.",
			},
			"sign_on_mode": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sign on mode of application.",
			},
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Application label.",
			},
			"plugin_settings": &schema.Schema{
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"three_field_settings", "no_plugin_settings", "custom_settings"},
				Elem: map[string]*schema.Schema{
					"button_field": &schema.Schema{
						Type:        schema.TypeString,
						Required:    true,
						Description: "Login button field",
					},
					"password_field": &schema.Schema{
						Type:        schema.TypeString,
						Required:    true,
						Description: "Login password field",
					},
					"username_field": &schema.Schema{
						Type:        schema.TypeString,
						Required:    true,
						Description: "Login username field",
					},
					"url": &schema.Schema{
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Login URL",
						ValidateFunc: validateIsURL,
					},
					"url_regex": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "A regex that further restricts URL to the specified regex",
					},
				},
			},
			"three_field_settings": &schema.Schema{
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"plugin_settings", "no_plugin_settings", "custom_settings"},
				Elem: map[string]*schema.Schema{
					"button_selector": &schema.Schema{
						Type:        schema.TypeString,
						Required:    true,
						Description: "Login button field CSS selector",
					},
					"password_selector": &schema.Schema{
						Type:        schema.TypeString,
						Required:    true,
						Description: "Login password field CSS selector",
					},
					"username_selector": &schema.Schema{
						Type:        schema.TypeString,
						Required:    true,
						Description: "Login username field CSS selector",
					},
					"extra_field_selector": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Extra field CSS selector",
					},
					"extra_field_value": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Value for extra form field",
					},
					"url": &schema.Schema{
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Login URL",
						ValidateFunc: validateIsURL,
					},
					"url_regex": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "A regex that further restricts URL to the specified regex",
					},
				},
			},
			"no_plugin_settings": &schema.Schema{
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"three_field_settings", "plugin_settings", "custom_settings"},
				Elem: map[string]*schema.Schema{
					"password_field": &schema.Schema{
						Type:        schema.TypeString,
						Required:    true,
						Description: "Login password field",
					},
					"username_field": &schema.Schema{
						Type:        schema.TypeString,
						Required:    true,
						Description: "Login username field",
					},
					"url": &schema.Schema{
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Login URL",
						ValidateFunc: validateIsURL,
					},
					"optional_field1": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name of optional param in the login form",
					},
					"optional_field1_value": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name of optional value in login form",
					},
					"optional_field2": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name of optional param in the login form",
					},
					"optional_field2_value": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name of optional value in login form",
					},
					"optional_field3": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name of optional param in the login form",
					},
					"optional_field3_value": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name of optional value in login form",
					},
				},
			},
			"custom_settings": &schema.Schema{
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"three_field_settings", "no_plugin_settings", "plugin_settings"},
				Elem: map[string]*schema.Schema{
					"url": &schema.Schema{
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Login URL",
						ValidateFunc: validateIsURL,
					},
					"redirect_url": &schema.Schema{
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Post login redirect URL",
						ValidateFunc: validateIsURL,
					},
				},
			},
			"status": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ACTIVE",
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
				Description:  "Status of application.",
			},
			"sign_on_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sign on URL",
			},
			"redirect_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Post login redirect URL",
			},
			"accessibility_self_service": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable self service",
			},
			"accessibility_error_redirect_url": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Custom error page URL",
				ValidateFunc: validateIsURL,
			},
			"auto_submit_toolbar": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Display auto submit toolbar",
			},
			"hide_ios": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Do not display application icon on mobile app",
			},
			"hide_web": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Do not display application icon to users",
			},
			"app_links_login": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Display app login link",
			},
			"credentials_scheme": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"EDIT_USERNAME_AND_PASSWORD",
						"ADMIN_SETS_CREDENTIALS",
						"EDIT_PASSWORD_ONLY",
						"EXTERNAL_PASSWORD_SYNC",
						"SHARED_USERNAME_AND_PASSWORD",
					},
					false,
				),
				Description: "Application credentials scheme",
			},
			"user_name_template": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Username template",
			},
			"user_name_template_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"BUILT_IN"}, false),
				Description:  "Username template type",
			},
		},
	}
}

func resourceSwaAppExists(d *schema.ResourceData, m interface{}) (bool, error) {
	app := okta.NewApplication()
	err := fetchApp(d, m, app)

	// Not sure if a non-nil app with an empty ID is possible but checking to avoid false positives.
	return app != nil && app.Id != "", err
}

func resourceSwaAppCreate(d *schema.ResourceData, m interface{}) error {
	runtime.Breakpoint()
	client := getOktaClientFromMetadata(m)
	app := buildApp(d, m)
	activate := d.Get("status").(string) == "ACTIVE"
	params := &query.Params{Activate: &activate}
	_, _, err := client.Application.CreateApplication(app, params)

	if err != nil {
		return err
	}

	d.SetId(app.Id)

	return resourceSwaAppRead(d, m)
}

func resourceSwaAppRead(d *schema.ResourceData, m interface{}) error {
	d.Get()
	app := okta.NewSwaApplication()
	err := fetchApp(d, m, app)

	if err != nil {
		return err
	}

	d.Set("name", app.Name)
	d.Set("status", app.Status)
	d.Set("sign_on_mode", app.SignOnMode)
	d.Set("label", app.Label)
	d.Set("auto_submit_toolbar", *app.Visibility.AutoSubmitToolbar)
	d.Set("hide_ios", *app.Visibility.Hide.IOS)
	d.Set("hide_web", *app.Visibility.Hide.Web)
	d.Set("plugin_settings.button_field", app.Settings.App.ButtonField)
	d.Set("plugin_settings.password_field", app.Settings.App.PasswordField)
	d.Set("plugin_settings.username_field", app.Settings.App.UsernameField)
	d.Set("plugin_settings.url", app.Settings.App.Url)
	d.Set("plugin_settings.url_regex", app.Settings.App.LoginUrlRegex)
			"three_field_settings": &schema.Schema{
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"plugin_settings", "no_plugin_settings", "custom_settings"},
				Elem: map[string]*schema.Schema{
					"button_selector": &schema.Schema{
						Type:        schema.TypeString,
						Required:    true,
						Description: "Login button field CSS selector",
					},
					"password_selector": &schema.Schema{
						Type:        schema.TypeString,
						Required:    true,
						Description: "Login password field CSS selector",
					},
					"username_selector": &schema.Schema{
						Type:        schema.TypeString,
						Required:    true,
						Description: "Login username field CSS selector",
					},
					"extra_field_selector": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Extra field CSS selector",
					},
					"extra_field_value": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Value for extra form field",
					},
					"url": &schema.Schema{
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Login URL",
						ValidateFunc: validateIsURL,
					},
					"url_regex": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "A regex that further restricts URL to the specified regex",
					},
				},
			},
			"no_plugin_settings": &schema.Schema{
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"three_field_settings", "plugin_settings", "custom_settings"},
				Elem: map[string]*schema.Schema{
					"password_field": &schema.Schema{
						Type:        schema.TypeString,
						Required:    true,
						Description: "Login password field",
					},
					"username_field": &schema.Schema{
						Type:        schema.TypeString,
						Required:    true,
						Description: "Login username field",
					},
					"url": &schema.Schema{
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Login URL",
						ValidateFunc: validateIsURL,
					},
					"optional_field1": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name of optional param in the login form",
					},
					"optional_field1_value": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name of optional value in login form",
					},
					"optional_field2": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name of optional param in the login form",
					},
					"optional_field2_value": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name of optional value in login form",
					},
					"optional_field3": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name of optional param in the login form",
					},
					"optional_field3_value": &schema.Schema{
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name of optional value in login form",
					},
				},
			},
			"custom_settings": &schema.Schema{
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"three_field_settings", "no_plugin_settings", "plugin_settings"},
				Elem: map[string]*schema.Schema{
					"url": &schema.Schema{
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Login URL",
						ValidateFunc: validateIsURL,
					},
					"redirect_url": &schema.Schema{
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Post login redirect URL",
						ValidateFunc: validateIsURL,
					},
				},
			},
			"status": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ACTIVE",
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
				Description:  "Status of application.",
			},
			"sign_on_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sign on URL",
			},
			"redirect_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Post login redirect URL",
			},
			"accessibility_self_service": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable self service",
			},
			"accessibility_error_redirect_url": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Custom error page URL",
				ValidateFunc: validateIsURL,
			},
			"auto_submit_toolbar": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Display auto submit toolbar",
			},
			"hide_ios": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Do not display application icon on mobile app",
			},
			"hide_web": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Do not display application icon to users",
			},
			"app_links_login": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Display app login link",
			},
			"credentials_scheme": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"EDIT_USERNAME_AND_PASSWORD",
						"ADMIN_SETS_CREDENTIALS",
						"EDIT_PASSWORD_ONLY",
						"EXTERNAL_PASSWORD_SYNC",
						"SHARED_USERNAME_AND_PASSWORD",
					},
					false,
				),
				Description: "Application credentials scheme",
			},
			"user_name_template": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Username template",
			},
			"user_name_template_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"BUILT_IN"}, false),
				Description:  "Username template type",
			},

	return nil
}

func resourceSwaAppUpdate(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	app := buildApp(d, m)
	_, _, err := client.Application.UpdateApplication(d.Id(), app)

	if err != nil {
		return err
	}

	desiredStatus := d.Get("status").(string)
	err = setAppStatus(d, client, app.Status, desiredStatus)

	if err != nil {
		return err
	}

	return resourceSwaAppRead(d, m)
}

func resourceSwaAppDelete(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	_, err := client.Application.DeactivateApplication(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Application.DeleteApplication(d.Id())

	return err
}

func buildSwaApp(d *schema.ResourceData, m interface{}) *okta.SwaApplication {
	// Abstracts away name and SignOnMode which are constant for this app type.
	app := okta.NewSwaApplication()
	app.Label = d.Get("label").(string)
	app.Name = d.Get("name").(string)
	app.SignOnMode = getSignOnMode(d)

	autoSubmit := d.Get("auto_submit_toolbar").(bool)
	hideMobile := d.Get("hide_ios").(bool)
	hideWeb := d.Get("hide_web").(bool)
	app.Settings = okta.NewSwaApplicationSettings()
	app.Visibility = &okta.ApplicationVisibility{
		AutoSubmitToolbar: &autoSubmit,
		Hide: &okta.ApplicationVisibilityHide{
			IOS: &hideMobile,
			Web: &hideWeb,
		},
	}

	return app
}

func getSignOnMode(d *schema.ResourceData) string {
	if _, ok := d.GetOk("custom_settings.url"); ok {
		return autoLogin
	} else if _, ok := d.GetOk("no_plugin_settings.url"); ok {
		return securePasswordStore
	}

	return browserPlugin
}
