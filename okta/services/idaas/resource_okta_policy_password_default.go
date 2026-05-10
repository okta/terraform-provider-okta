package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func resourcePolicyPasswordDefault() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyPasswordDefaultUpdate,
		ReadContext:   resourcePolicyPasswordDefaultRead,
		UpdateContext: resourcePolicyPasswordDefaultUpdate,
		DeleteContext: utils.ResourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				policy, err := setDefaultPasswordPolicyV6(ctx, d, meta)
				if err != nil {
					return nil, err
				}
				if policy.Conditions != nil && policy.Conditions.AuthProvider != nil {
					_ = d.Set("default_auth_provider", policy.Conditions.AuthProvider.GetProvider())
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Description: "Configures default password policy. This resource allows you to configure default password policy.",
		Schema: buildDefaultPolicySchema(map[string]*schema.Schema{
			"default_auth_provider": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Authentication Provider",
			},
			"password_min_length": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Minimum password length. Default is `8`.",
				Default:     8,
			},
			"password_min_lowercase": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "If a password must contain at least one lower case letter: 0 = no, 1 = yes. Default = 1",
				Default:     1,
			},
			"password_min_uppercase": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "If a password must contain at least one upper case letter: 0 = no, 1 = yes. Default = 1",
				Default:     1,
			},
			"password_min_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "If a password must contain at least one number: 0 = no, 1 = yes. Default = `1`",
				Default:     1,
			},
			"password_min_symbol": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "If a password must contain at least one symbol (!@#$%^&*): 0 = no, 1 = yes. Default = `0`",
				Default:     0,
			},
			"password_exclude_username": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If the user name must be excluded from the password. Default: `true`",
				Default:     true,
			},
			"password_exclude_first_name": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "User firstName attribute must be excluded from the password",
			},
			"password_exclude_last_name": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "User lastName attribute must be excluded from the password",
			},
			"password_dictionary_lookup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Check Passwords Against Common Password Dictionary. Default: `false`",
				Default:     false,
			},
			"password_max_age_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Length in days a password is valid before expiry: 0 = no limit. Default: `0`",
				Default:     0,
			},
			"password_expire_warn_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Length in days a user will be warned before password expiry: 0 = no warning. Default: `0`",
				Default:     0,
			},
			"password_min_age_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Minimum time interval in minutes between password changes: 0 = no limit. Default: `0`",
				Default:     0,
			},
			"password_history_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of distinct passwords that must be created before they can be reused: 0 = none. Default: `4`",
				Default:     4,
				// API documentation says default is 0 but it appears in acceptance testing on different orgs to now be 4 by default
				// historyCount -> https://developer.okta.com/docs/reference/api/policy/#age-object
			},
			"password_max_lockout_attempts": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of unsuccessful login attempts allowed before lockout: 0 = no limit. Default: `10`",
				Default:     10,
			},
			"password_auto_unlock_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of minutes before a locked account is unlocked: 0 = no limit. Default: `0`",
				Default:     0,
			},
			"password_show_lockout_failures": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If a user should be informed when their account is locked. Default: `false`",
				Default:     false,
			},
			"password_lockout_notification_channels": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Notification channels to use to notify a user when their account has been locked.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"question_min_length": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Min length of the password recovery question answer. Default: `4`",
				Default:     4,
			},
			"email_recovery": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable or disable email password recovery: ACTIVE or INACTIVE. Default: `ACTIVE`",
				Default:     StatusActive,
			},
			"recovery_email_token": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Lifetime in minutes of the recovery email token. Default: `60`",
				Default:     60,
			},
			"sms_recovery": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable or disable SMS password recovery: ACTIVE or INACTIVE. Default: `INACTIVE`",
				Default:     StatusInactive,
			},
			"question_recovery": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable or disable security question password recovery: ACTIVE or INACTIVE. Default: `ACTIVE`",
				Default:     StatusActive,
			},
			"call_recovery": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable or disable voice call recovery: ACTIVE or INACTIVE. Default: `INACTIVE`",
				Default:     StatusInactive,
			},
			"skip_unlock": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When an Active Directory user is locked out of Okta, the Okta unlock operation should also attempt to unlock the user's Windows account. Default: `false`",
				Default:     false,
			},
			"breached_password_expire_after_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of days after a breached password is detected before the user's password expires. Valid values: 0 through 10. If set to 0, expiry is immediate. Only applicable when `breached_password_logout_enabled` is `true`.",
				Default:     0,
			},
			"breached_password_logout_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If `true`, the user's sessions are terminated immediately when their credentials are detected as part of a breach. Requires `breached_password_expire_after_days` to also be configured. Default: `false`",
				Default:     false,
			},
			"breached_password_delegated_workflow_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the workflow to run when a breached password is found during a sign-in attempt.",
			},
		}),
	}
}

func resourcePolicyPasswordDefaultUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.Id() == "" {
		policy, err := setDefaultPasswordPolicyV6(ctx, d, meta)
		if err != nil {
			return diag.FromErr(err)
		}
		if policy.Conditions != nil && policy.Conditions.AuthProvider != nil {
			_ = d.Set("default_auth_provider", policy.Conditions.AuthProvider.GetProvider())
		}
	}
	if _, err := replacePolicyV6(ctx, d, meta, buildDefaultPasswordPolicy(d)); err != nil {
		return diag.Errorf("failed to update default password policy: %v", err)
	}
	return resourcePolicyPasswordDefaultRead(ctx, d, meta)
}

func resourcePolicyPasswordDefaultRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policy, err := getPolicyV6(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to get default password policy: %v", err)
	}
	if policy == nil {
		return nil
	}

	var pw *v6okta.PasswordPolicyPasswordSettings
	var factors *v6okta.PasswordPolicyRecoveryFactors
	if policy.Settings != nil {
		pw = policy.Settings.Password
		if policy.Settings.Recovery != nil {
			factors = policy.Settings.Recovery.Factors
		}
	}

	if pw != nil && pw.Lockout != nil {
		if err := d.Set("password_lockout_notification_channels", utils.ConvertStringSliceToSet(pw.Lockout.UserLockoutNotificationChannels)); err != nil {
			return diag.Errorf("error setting notification channels for resource %s: %v", d.Id(), err)
		}
	}
	if pw != nil && pw.Complexity != nil {
		if pw.Complexity.Dictionary != nil && pw.Complexity.Dictionary.Common != nil {
			_ = d.Set("password_dictionary_lookup", pw.Complexity.Dictionary.Common.GetExclude())
		}
		if pw.Complexity.MinLength != nil {
			_ = d.Set("password_min_length", int(*pw.Complexity.MinLength))
		}
		if pw.Complexity.MinLowerCase != nil {
			_ = d.Set("password_min_lowercase", int(*pw.Complexity.MinLowerCase))
		}
		if pw.Complexity.MinUpperCase != nil {
			_ = d.Set("password_min_uppercase", int(*pw.Complexity.MinUpperCase))
		}
		if pw.Complexity.MinNumber != nil {
			_ = d.Set("password_min_number", int(*pw.Complexity.MinNumber))
		}
		if pw.Complexity.MinSymbol != nil {
			_ = d.Set("password_min_symbol", int(*pw.Complexity.MinSymbol))
		}
		if pw.Complexity.ExcludeUsername != nil {
			_ = d.Set("password_exclude_username", *pw.Complexity.ExcludeUsername)
		}
	}
	if pw != nil && pw.Age != nil {
		if pw.Age.MaxAgeDays != nil {
			_ = d.Set("password_max_age_days", int(*pw.Age.MaxAgeDays))
		}
		if pw.Age.ExpireWarnDays != nil {
			_ = d.Set("password_expire_warn_days", int(*pw.Age.ExpireWarnDays))
		}
		if pw.Age.MinAgeMinutes != nil {
			_ = d.Set("password_min_age_minutes", int(*pw.Age.MinAgeMinutes))
		}
		if pw.Age.HistoryCount != nil {
			_ = d.Set("password_history_count", int(*pw.Age.HistoryCount))
		}
	}
	if pw != nil && pw.Lockout != nil {
		if pw.Lockout.MaxAttempts != nil {
			_ = d.Set("password_max_lockout_attempts", int(*pw.Lockout.MaxAttempts))
		}
		if pw.Lockout.AutoUnlockMinutes != nil {
			_ = d.Set("password_auto_unlock_minutes", int(*pw.Lockout.AutoUnlockMinutes))
		}
		if pw.Lockout.ShowLockoutFailures != nil {
			_ = d.Set("password_show_lockout_failures", *pw.Lockout.ShowLockoutFailures)
		}
	}
	if factors != nil && factors.RecoveryQuestion != nil && factors.RecoveryQuestion.Properties != nil &&
		factors.RecoveryQuestion.Properties.Complexity != nil && factors.RecoveryQuestion.Properties.Complexity.MinLength != nil {
		_ = d.Set("question_min_length", int(*factors.RecoveryQuestion.Properties.Complexity.MinLength))
	}
	if factors != nil && factors.OktaEmail != nil && factors.OktaEmail.Properties != nil &&
		factors.OktaEmail.Properties.RecoveryToken != nil && factors.OktaEmail.Properties.RecoveryToken.TokenLifetimeMinutes != nil {
		_ = d.Set("recovery_email_token", int(*factors.OktaEmail.Properties.RecoveryToken.TokenLifetimeMinutes))
	}
	if factors != nil && factors.OktaSms != nil {
		_ = d.Set("sms_recovery", factors.OktaSms.GetStatus())
	}
	if factors != nil && factors.OktaEmail != nil {
		_ = d.Set("email_recovery", factors.OktaEmail.GetStatus())
	}
	if factors != nil && factors.RecoveryQuestion != nil {
		_ = d.Set("question_recovery", factors.RecoveryQuestion.GetStatus())
	}
	if factors != nil && factors.OktaCall != nil {
		_ = d.Set("call_recovery", factors.OktaCall.GetStatus())
	}
	if policy.Settings != nil && policy.Settings.Delegation != nil && policy.Settings.Delegation.Options != nil && policy.Settings.Delegation.Options.SkipUnlock != nil {
		_ = d.Set("skip_unlock", *policy.Settings.Delegation.Options.SkipUnlock)
	}
	if pw != nil && pw.Complexity != nil {
		for _, v := range pw.Complexity.ExcludeAttributes {
			switch v {
			case "firstName":
				_ = d.Set("password_exclude_first_name", true)
			case "lastName":
				_ = d.Set("password_exclude_last_name", true)
			}
		}
	}
	if pw != nil && pw.BreachedProtection != nil {
		bp := pw.BreachedProtection
		if bp.HasLogoutEnabled() {
			_ = d.Set("breached_password_logout_enabled", bp.GetLogoutEnabled())
		}
		if bp.HasExpireAfterDays() {
			_ = d.Set("breached_password_expire_after_days", int(bp.GetExpireAfterDays()))
		}
		if bp.HasDelegatedWorkflowId() {
			_ = d.Set("breached_password_delegated_workflow_id", bp.GetDelegatedWorkflowId())
		}
	}
	return nil
}

func buildDefaultPasswordPolicy(d *schema.ResourceData) *v6okta.PasswordPolicy {
	template := v6okta.NewPasswordPolicy(d.Get("name").(string), "PASSWORD")
	// template.SetStatus(d.Get("status").(string))
	template.SetDescription(d.Get("description").(string))
	template.SetPriority(1) // default priority is 1
	authProvider := &v6okta.PasswordPolicyAuthenticationProviderCondition{}
	authProvider.SetProvider(d.Get("default_auth_provider").(string))
	template.Conditions = &v6okta.PasswordPolicyConditions{
		AuthProvider: authProvider,
		People: &v6okta.AuthenticatorEnrollmentPolicyConditionsAllOfPeople{
			Groups: &v6okta.AuthenticatorEnrollmentPolicyConditionsAllOfPeopleGroups{
				Include: []string{d.Get("default_included_group_id").(string)},
			},
		},
	}
	passwordSettings := &v6okta.PasswordPolicyPasswordSettings{
		Age: &v6okta.PasswordPolicyPasswordSettingsAge{
			ExpireWarnDays: utils.Int32Ptr(d.Get("password_expire_warn_days").(int)),
			HistoryCount:   utils.Int32Ptr(d.Get("password_history_count").(int)),
			MaxAgeDays:     utils.Int32Ptr(d.Get("password_max_age_days").(int)),
			MinAgeMinutes:  utils.Int32Ptr(d.Get("password_min_age_minutes").(int)),
		},
		Complexity: &v6okta.PasswordPolicyPasswordSettingsComplexity{
			Dictionary: &v6okta.PasswordDictionary{
				Common: &v6okta.PasswordDictionaryCommon{
					Exclude: utils.BoolPtr(d.Get("password_dictionary_lookup").(bool)),
				},
			},
			ExcludeAttributes: getExcludedAttrs(d.Get("password_exclude_first_name").(bool), d.Get("password_exclude_last_name").(bool)),
			ExcludeUsername:   utils.BoolPtr(d.Get("password_exclude_username").(bool)),
			MinLength:         utils.Int32Ptr(d.Get("password_min_length").(int)),
			MinLowerCase:      utils.Int32Ptr(d.Get("password_min_lowercase").(int)),
			MinNumber:         utils.Int32Ptr(d.Get("password_min_number").(int)),
			MinSymbol:         utils.Int32Ptr(d.Get("password_min_symbol").(int)),
			MinUpperCase:      utils.Int32Ptr(d.Get("password_min_uppercase").(int)),
		},
		Lockout: &v6okta.PasswordPolicyPasswordSettingsLockout{
			AutoUnlockMinutes:               utils.Int32Ptr(d.Get("password_auto_unlock_minutes").(int)),
			MaxAttempts:                     utils.Int32Ptr(d.Get("password_max_lockout_attempts").(int)),
			ShowLockoutFailures:             utils.BoolPtr(d.Get("password_show_lockout_failures").(bool)),
			UserLockoutNotificationChannels: utils.ConvertInterfaceToStringSet(d.Get("password_lockout_notification_channels")),
		},
	}
	logoutEnabled := d.Get("breached_password_logout_enabled").(bool)
	delegatedWorkflowId := d.Get("breached_password_delegated_workflow_id").(string)
	if logoutEnabled || delegatedWorkflowId != "" {
		bp := v6okta.NewPasswordPolicyPasswordSettingsBreachedProtection()
		bp.SetExpireAfterDays(*utils.Int32Ptr(d.Get("breached_password_expire_after_days").(int)))
		bp.SetLogoutEnabled(logoutEnabled)
		if delegatedWorkflowId != "" {
			bp.SetDelegatedWorkflowId(delegatedWorkflowId)
		}
		passwordSettings.BreachedProtection = bp
	}
	template.Settings = &v6okta.PasswordPolicySettings{
		Password: passwordSettings,
		Recovery: &v6okta.PasswordPolicyRecoverySettings{
			Factors: &v6okta.PasswordPolicyRecoveryFactors{
				OktaCall: &v6okta.PasswordPolicyRecoveryFactorSettings{
					Status: utils.StringPtr(d.Get("call_recovery").(string)),
				},
				OktaSms: &v6okta.PasswordPolicyRecoveryFactorSettings{
					Status: utils.StringPtr(d.Get("sms_recovery").(string)),
				},
				OktaEmail: &v6okta.PasswordPolicyRecoveryEmail{
					Properties: &v6okta.PasswordPolicyRecoveryEmailProperties{
						RecoveryToken: &v6okta.PasswordPolicyRecoveryEmailRecoveryToken{
							TokenLifetimeMinutes: utils.Int32Ptr(d.Get("recovery_email_token").(int)),
						},
					},
					Status: utils.StringPtr(d.Get("email_recovery").(string)),
				},
				RecoveryQuestion: &v6okta.PasswordPolicyRecoveryQuestion{
					Properties: &v6okta.PasswordPolicyRecoveryQuestionProperties{
						Complexity: &v6okta.PasswordPolicyRecoveryQuestionComplexity{
							MinLength: utils.Int32Ptr(d.Get("question_min_length").(int)),
						},
					},
					Status: utils.StringPtr(d.Get("question_recovery").(string)),
				},
			},
		},
		Delegation: &v6okta.PasswordPolicyDelegationSettings{
			Options: &v6okta.PasswordPolicyDelegationSettingsOptions{
				SkipUnlock: utils.BoolPtr(d.Get("skip_unlock").(bool)),
			},
		},
	}
	return template
}
