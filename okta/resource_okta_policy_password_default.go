package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourcePolicyPasswordDefault() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyPasswordDefaultUpdate,
		ReadContext:   resourcePolicyPasswordDefaultRead,
		UpdateContext: resourcePolicyPasswordDefaultUpdate,
		DeleteContext: resourcePolicyPasswordDefaultDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				policy, err := setDefaultPolicy(ctx, d, m, sdk.PasswordPolicyType)
				if err != nil {
					return nil, err
				}
				_ = d.Set("default_auth_provider", policy.Conditions.AuthProvider.Provider)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: buildDefaultPolicySchema(map[string]*schema.Schema{
			"default_auth_provider": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Authentication Provider",
			},
			"password_min_length": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Minimum password length.",
				Default:     8,
			},
			"password_min_lowercase": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intBetween(0, 1),
				Description:      "If a password must contain at least one lower case letter: 0 = no, 1 = yes. Default = 1",
				Default:          1,
			},
			"password_min_uppercase": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intBetween(0, 1),
				Description:      "If a password must contain at least one upper case letter: 0 = no, 1 = yes. Default = 1",
				Default:          1,
			},
			"password_min_number": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intBetween(0, 1),
				Description:      "If a password must contain at least one number: 0 = no, 1 = yes. Default = 1",
				Default:          1,
			},
			"password_min_symbol": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intBetween(0, 1),
				Description:      "If a password must contain at least one symbol (!@#$%^&*): 0 = no, 1 = yes. Default = 1",
				Default:          0,
			},
			"password_exclude_username": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If the user name must be excluded from the password.",
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
				Description: "Check Passwords Against Common Password Dictionary.",
				Default:     false,
			},
			"password_max_age_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Length in days a password is valid before expiry: 0 = no limit.",
				Default:     0,
			},
			"password_expire_warn_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Length in days a user will be warned before password expiry: 0 = no warning.",
				Default:     0,
			},
			"password_min_age_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Minimum time interval in minutes between password changes: 0 = no limit.",
				Default:     0,
			},
			"password_history_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of distinct passwords that must be created before they can be reused: 0 = none.",
				Default:     0,
			},
			"password_max_lockout_attempts": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of unsuccessful login attempts allowed before lockout: 0 = no limit.",
				Default:     10,
			},
			"password_auto_unlock_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of minutes before a locked account is unlocked: 0 = no limit.",
				Default:     0,
			},
			"password_show_lockout_failures": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If a user should be informed when their account is locked.",
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
				Description: "Min length of the password recovery question answer.",
				Default:     4,
			},
			"email_recovery": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringInSlice([]string{statusActive, statusInactive}),
				Description:      "Enable or disable email password recovery: ACTIVE or INACTIVE.",
				Default:          statusActive,
			},
			"recovery_email_token": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Lifetime in minutes of the recovery email token.",
				Default:     60,
			},
			"sms_recovery": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringInSlice([]string{statusActive, statusInactive}),
				Description:      "Enable or disable SMS password recovery: ACTIVE or INACTIVE.",
				Default:          statusInactive,
			},
			"question_recovery": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringInSlice([]string{statusActive, statusInactive}),
				Description:      "Enable or disable security question password recovery: ACTIVE or INACTIVE.",
				Default:          statusActive,
			},
			"call_recovery": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringInSlice([]string{statusActive, statusInactive}),
				Description:      "Enable or disable voice call recovery: ACTIVE or INACTIVE.",
				Default:          statusInactive,
			},
			"skip_unlock": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When an Active Directory user is locked out of Okta, the Okta unlock operation should also attempt to unlock the user's Windows account.",
				Default:     false,
			},
		}),
	}
}

func resourcePolicyPasswordDefaultUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()
	if id == "" {
		policy, err := setDefaultPolicy(ctx, d, m, sdk.PasswordPolicyType)
		if err != nil {
			return diag.FromErr(err)
		}
		id = policy.Id
		_ = d.Set("default_auth_provider", policy.Conditions.AuthProvider.Provider)
	}
	_, _, err := getSupplementFromMetadata(m).UpdatePolicy(ctx, id, buildDefaultPasswordPolicy(d))
	if err != nil {
		return diag.Errorf("failed to update default password policy: %v", err)
	}
	return resourcePolicyPasswordDefaultRead(ctx, d, m)
}

func resourcePolicyPasswordDefaultRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policy, err := getPolicy(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to get default password policy: %v", err)
	}
	if policy == nil {
		return nil
	}
	err = d.Set("password_lockout_notification_channels", convertStringSetToInterface(policy.Settings.Password.Lockout.UserLockoutNotificationChannels))
	if err != nil {
		return diag.Errorf("error setting notification channels for resource %s: %v", d.Id(), err)
	}
	if policy.Settings.Password.Complexity.Dictionary != nil && policy.Settings.Password.Complexity.Dictionary.Common != nil {
		_ = d.Set("password_dictionary_lookup", policy.Settings.Password.Complexity.Dictionary.Common.Exclude)
	}
	_ = d.Set("password_min_length", policy.Settings.Password.Complexity.MinLength)
	_ = d.Set("password_min_lowercase", policy.Settings.Password.Complexity.MinLowerCase)
	_ = d.Set("password_min_uppercase", policy.Settings.Password.Complexity.MinUpperCase)
	_ = d.Set("password_min_number", policy.Settings.Password.Complexity.MinNumber)
	_ = d.Set("password_min_symbol", policy.Settings.Password.Complexity.MinSymbol)
	_ = d.Set("password_exclude_username", policy.Settings.Password.Complexity.ExcludeUsername)
	_ = d.Set("password_max_age_days", policy.Settings.Password.Age.MaxAgeDays)
	_ = d.Set("password_expire_warn_days", policy.Settings.Password.Age.ExpireWarnDays)
	_ = d.Set("password_min_age_minutes", policy.Settings.Password.Age.MinAgeMinutes)
	_ = d.Set("password_history_count", policy.Settings.Password.Age.HistoryCount)
	_ = d.Set("password_max_lockout_attempts", policy.Settings.Password.Lockout.MaxAttempts)
	_ = d.Set("password_auto_unlock_minutes", policy.Settings.Password.Lockout.AutoUnlockMinutes)
	_ = d.Set("password_show_lockout_failures", policy.Settings.Password.Lockout.ShowLockoutFailures)
	_ = d.Set("question_min_length", policy.Settings.Recovery.Factors.RecoveryQuestion.Properties.Complexity.MinLength)
	_ = d.Set("recovery_email_token", policy.Settings.Recovery.Factors.OktaEmail.Properties.RecoveryToken.TokenLifetimeMinutes)
	_ = d.Set("sms_recovery", policy.Settings.Recovery.Factors.OktaSms.Status)
	_ = d.Set("email_recovery", policy.Settings.Recovery.Factors.OktaEmail.Status)
	_ = d.Set("question_recovery", policy.Settings.Recovery.Factors.RecoveryQuestion.Status)
	_ = d.Set("call_recovery", policy.Settings.Recovery.Factors.OktaCall.Status)
	_ = d.Set("skip_unlock", policy.Settings.Delegation.Options.SkipUnlock)
	for _, v := range policy.Settings.Password.Complexity.ExcludeAttributes {
		switch v {
		case "firstName":
			_ = d.Set("password_exclude_first_name", true)
		case "lastName":
			_ = d.Set("password_exclude_last_name", true)
		}
	}
	return nil
}

// Default policy can not be removed
func resourcePolicyPasswordDefaultDelete(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}

// create or update a password policy
func buildDefaultPasswordPolicy(d *schema.ResourceData) sdk.Policy {
	policy := sdk.PasswordPolicy()
	policy.Name = d.Get("name").(string)
	policy.Status = d.Get("status").(string)
	policy.Description = d.Get("description").(string)
	policy.Priority = int64(d.Get("priority").(int))
	policy.Conditions = &okta.PolicyRuleConditions{
		AuthProvider: &okta.PasswordPolicyAuthenticationProviderCondition{
			Provider: d.Get("default_auth_provider").(string),
		},
		People: &okta.PolicyPeopleCondition{
			Groups: &okta.GroupCondition{
				Include: []string{d.Get("default_included_group_id").(string)},
			},
		},
	}
	// Okta defaults
	// we add the defaults here & not in the schema map to avoid defaults appearing in the terraform plan diff
	policy.Settings = &sdk.PolicySettings{
		Password: &sdk.PasswordPolicyPasswordSettings{
			Age: &sdk.PasswordPolicyPasswordSettingsAge{
				ExpireWarnDays: int64(d.Get("password_expire_warn_days").(int)),
				HistoryCount:   int64(d.Get("password_history_count").(int)),
				MaxAgeDays:     int64(d.Get("password_max_age_days").(int)),
				MinAgeMinutes:  int64(d.Get("password_min_age_minutes").(int)),
			},
			Complexity: &sdk.PasswordPolicyPasswordSettingsComplexity{
				Dictionary: &okta.PasswordDictionary{
					Common: &okta.PasswordDictionaryCommon{
						Exclude: boolPtr(d.Get("password_dictionary_lookup").(bool)),
					},
				},
				ExcludeAttributes: getExcludedAttrs(d.Get("password_exclude_first_name").(bool), d.Get("password_exclude_last_name").(bool)),
				ExcludeUsername:   boolPtr(d.Get("password_exclude_username").(bool)),
				MinLength:         int64(d.Get("password_min_length").(int)),
				MinLowerCase:      int64(d.Get("password_min_lowercase").(int)),
				MinNumber:         int64(d.Get("password_min_number").(int)),
				MinSymbol:         int64(d.Get("password_min_symbol").(int)),
				MinUpperCase:      int64(d.Get("password_min_uppercase").(int)),
			},
			Lockout: &sdk.PasswordPolicyPasswordSettingsLockout{
				AutoUnlockMinutes:               int64(d.Get("password_auto_unlock_minutes").(int)),
				MaxAttempts:                     int64(d.Get("password_max_lockout_attempts").(int)),
				ShowLockoutFailures:             boolPtr(d.Get("password_show_lockout_failures").(bool)),
				UserLockoutNotificationChannels: convertInterfaceToStringSet(d.Get("password_lockout_notification_channels")),
			},
		},
		Recovery: &sdk.PasswordPolicyRecoverySettings{
			Factors: &sdk.PasswordPolicyRecoveryFactors{
				OktaCall: &okta.PasswordPolicyRecoveryFactorSettings{
					Status: d.Get("call_recovery").(string),
				},
				OktaSms: &okta.PasswordPolicyRecoveryFactorSettings{
					Status: d.Get("sms_recovery").(string),
				},
				OktaEmail: &sdk.PasswordPolicyRecoveryEmail{
					Properties: &sdk.PasswordPolicyRecoveryEmailProperties{
						RecoveryToken: &sdk.PasswordPolicyRecoveryEmailRecoveryToken{
							TokenLifetimeMinutes: int64(d.Get("recovery_email_token").(int)),
						},
					},
					Status: d.Get("email_recovery").(string),
				},
				RecoveryQuestion: &sdk.PasswordPolicyRecoveryQuestion{
					Properties: &sdk.PasswordPolicyRecoveryQuestionProperties{
						Complexity: &sdk.PasswordPolicyRecoveryQuestionComplexity{
							MinLength: int64(d.Get("question_min_length").(int)),
						},
					},
					Status: d.Get("question_recovery").(string),
				},
			},
		},
		Delegation: &okta.PasswordPolicyDelegationSettings{
			Options: &okta.PasswordPolicyDelegationSettingsOptions{
				SkipUnlock: boolPtr(d.Get("skip_unlock").(bool)),
			},
		},
	}
	return policy
}
