package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func ResourcePolicyPassword() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyPasswordCreate,
		ReadContext:   resourcePolicyPasswordRead,
		UpdateContext: resourcePolicyPasswordUpdate,
		DeleteContext: resourcePolicyPasswordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Creates a Password Policy. This resource allows you to create and configure a Password Policy.",
		Schema: buildPolicySchema(map[string]*schema.Schema{
			"auth_provider": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Authentication Provider: `OKTA`, `ACTIVE_DIRECTORY` or `LDAP`. Default: `OKTA`",
				Default:     "OKTA",
			},
			"password_min_length": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Minimum password length. Default: `8`",
				Default:     8,
			},
			"password_min_lowercase": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "If a password must contain at least one lower case letter: 0 = no, 1 = yes. Default: `1`",
				Default:     1,
			},
			"password_min_uppercase": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "If a password must contain at least one upper case letter: 0 = no, 1 = yes. Default: `1`",
				Default:     1,
			},
			"password_min_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "If a password must contain at least one number: 0 = no, 1 = yes. Default: `1`",
				Default:     1,
			},
			"password_min_symbol": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "If a password must contain at least one symbol (!@#$%^&*): 0 = no, 1 = yes. Default: `0`",
				Default:     0,
			},
			"password_exclude_username": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If the user name must be excluded from the password. Default: `true`",
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
				Description: "Number of distinct passwords that must be created before they can be reused: 0 = none. Default: `0`",
				Default:     0,
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
				Description: "Enable or disable email password recovery: `ACTIVE` or `INACTIVE`. Default: `ACTIVE`",
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
				Description: "Enable or disable SMS password recovery: `ACTIVE` or `INACTIVE`. Default: `INACTIVE`",
				Default:     StatusInactive,
			},
			"question_recovery": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable or disable security question password recovery: `ACTIVE` or `INACTIVE`. Default: `ACTIVE`",
				Default:     StatusActive,
			},
			"call_recovery": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable or disable voice call recovery: `ACTIVE` or `INACTIVE`. Default: `INACTIVE`",
				Default:     StatusInactive,
			},
			"skip_unlock": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When an Active Directory user is locked out of Okta, the Okta unlock operation should also attempt to unlock the user's Windows account. Default: `false`",
				Default:     false,
			},
		}),
	}
}

func resourcePolicyPasswordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	template := buildPasswordPolicy(d)
	err := createPolicy(ctx, d, meta, template)
	if err != nil {
		return diag.Errorf("failed to create password policy: %v", err)
	}
	return resourcePolicyPasswordRead(ctx, d, meta)
}

func resourcePolicyPasswordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policy, err := getPolicy(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to get password policy: %v", err)
	}
	if policy == nil {
		return nil
	}

	// Update with upstream state when it is manually updated from Okta UI or API directly.
	// See https://github.com/okta/terraform-provider-okta/issues/61
	if policy.Conditions.AuthProvider != nil && policy.Conditions.AuthProvider.Provider != "" {
		_ = d.Set("auth_provider", policy.Conditions.AuthProvider.Provider)
	}

	if policy.Settings != nil {
		err = d.Set("password_lockout_notification_channels", utils.ConvertStringSliceToSet(policy.Settings.Password.Lockout.UserLockoutNotificationChannels))
		if err != nil {
			return diag.Errorf("error setting notification channels for resource %s: %v", d.Id(), err)
		}
		if policy.Settings.Password.Complexity.MinLengthPtr != nil {
			_ = d.Set("password_min_length", policy.Settings.Password.Complexity.MinLengthPtr)
		}
		if policy.Settings.Password.Complexity.MinLowerCasePtr != nil {
			_ = d.Set("password_min_lowercase", policy.Settings.Password.Complexity.MinLowerCasePtr)
		}
		if policy.Settings.Password.Complexity.MinUpperCasePtr != nil {
			_ = d.Set("password_min_uppercase", policy.Settings.Password.Complexity.MinUpperCasePtr)
		}
		if policy.Settings.Password.Complexity.MinNumberPtr != nil {
			_ = d.Set("password_min_number", policy.Settings.Password.Complexity.MinNumberPtr)
		}
		if policy.Settings.Password.Complexity.MinSymbolPtr != nil {
			_ = d.Set("password_min_symbol", policy.Settings.Password.Complexity.MinSymbolPtr)
		}
		_ = d.Set("password_exclude_username", policy.Settings.Password.Complexity.ExcludeUsername)
		if policy.Settings.Password.Complexity.Dictionary != nil && policy.Settings.Password.Complexity.Dictionary.Common != nil {
			_ = d.Set("password_dictionary_lookup", policy.Settings.Password.Complexity.Dictionary.Common.Exclude)
		}
		if policy.Settings.Password.Age.MaxAgeDaysPtr != nil {
			_ = d.Set("password_max_age_days", policy.Settings.Password.Age.MaxAgeDaysPtr)
		}
		if policy.Settings.Password.Age.ExpireWarnDaysPtr != nil {
			_ = d.Set("password_expire_warn_days", policy.Settings.Password.Age.ExpireWarnDaysPtr)
		}
		if policy.Settings.Password.Age.MinAgeMinutesPtr != nil {
			_ = d.Set("password_min_age_minutes", policy.Settings.Password.Age.MinAgeMinutesPtr)
		}
		if policy.Settings.Password.Age.HistoryCountPtr != nil {
			_ = d.Set("password_history_count", policy.Settings.Password.Age.HistoryCountPtr)
		}
		if policy.Settings.Password.Lockout.MaxAttemptsPtr != nil {
			_ = d.Set("password_max_lockout_attempts", policy.Settings.Password.Lockout.MaxAttemptsPtr)
		}
		if policy.Settings.Password.Lockout.AutoUnlockMinutesPtr != nil {
			_ = d.Set("password_auto_unlock_minutes", policy.Settings.Password.Lockout.AutoUnlockMinutesPtr)
		}
		_ = d.Set("password_show_lockout_failures", policy.Settings.Password.Lockout.ShowLockoutFailures)
		if policy.Settings.Recovery.Factors.RecoveryQuestion.Properties.Complexity.MinLengthPtr != nil {
			_ = d.Set("question_min_length", policy.Settings.Recovery.Factors.RecoveryQuestion.Properties.Complexity.MinLengthPtr)
		}
		if policy.Settings.Recovery.Factors.OktaEmail.Properties.RecoveryToken.TokenLifetimeMinutesPtr != nil {
			_ = d.Set("recovery_email_token", policy.Settings.Recovery.Factors.OktaEmail.Properties.RecoveryToken.TokenLifetimeMinutesPtr)
		}
		_ = d.Set("sms_recovery", policy.Settings.Recovery.Factors.OktaSms.Status)
		_ = d.Set("email_recovery", policy.Settings.Recovery.Factors.OktaEmail.Status)
		_ = d.Set("question_recovery", policy.Settings.Recovery.Factors.RecoveryQuestion.Status)
		_ = d.Set("call_recovery", policy.Settings.Recovery.Factors.OktaCall.Status)
		_ = d.Set("skip_unlock", policy.Settings.Delegation.Options.SkipUnlock)

		excludedAttrs := policy.Settings.Password.Complexity.ExcludeAttributes
		if len(excludedAttrs) > 0 {
			for _, v := range excludedAttrs {
				switch v {
				case "firstName":
					_ = d.Set("password_exclude_first_name", true)
				case "lastName":
					_ = d.Set("password_exclude_last_name", true)
				}
			}
		}
	}
	err = syncPolicyFromUpstream(d, policy)
	if err != nil {
		return diag.Errorf("failed to set password policy: %v", err)
	}
	return nil
}

func resourcePolicyPasswordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	template := buildPasswordPolicy(d)
	err := updatePolicy(ctx, d, meta, template)
	if err != nil {
		return diag.Errorf("failed to update password policy: %v", err)
	}
	return resourcePolicyPasswordRead(ctx, d, meta)
}

func resourcePolicyPasswordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := deletePolicy(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to delete password policy: %v", err)
	}
	return nil
}

// create or update a password policy
func buildPasswordPolicy(d *schema.ResourceData) sdk.SdkPolicy {
	template := sdk.PasswordPolicy()
	template.Name = d.Get("name").(string)
	template.Status = d.Get("status").(string)
	if description, ok := d.GetOk("description"); ok {
		template.Description = description.(string)
	}
	if priority, ok := d.GetOk("priority"); ok {
		template.PriorityPtr = utils.Int64Ptr(priority.(int))
	}
	template.Conditions = &sdk.PolicyRuleConditions{
		AuthProvider: &sdk.PasswordPolicyAuthenticationProviderCondition{
			Provider: d.Get("auth_provider").(string),
		},
		People: getGroups(d),
	}
	// Okta defaults
	// we add the defaults here & not in the schema map to avoid defaults appearing in the terraform plan diff
	template.Settings = &sdk.SdkPolicySettings{
		Password: &sdk.PasswordPolicyPasswordSettings{
			Age: &sdk.PasswordPolicyPasswordSettingsAge{
				ExpireWarnDaysPtr: utils.Int64Ptr(d.Get("password_expire_warn_days").(int)),
				HistoryCountPtr:   utils.Int64Ptr(d.Get("password_history_count").(int)),
				MaxAgeDaysPtr:     utils.Int64Ptr(d.Get("password_max_age_days").(int)),
				MinAgeMinutesPtr:  utils.Int64Ptr(d.Get("password_min_age_minutes").(int)),
			},
			Complexity: &sdk.PasswordPolicyPasswordSettingsComplexity{
				Dictionary: &sdk.PasswordDictionary{
					Common: &sdk.PasswordDictionaryCommon{
						Exclude: utils.BoolPtr(d.Get("password_dictionary_lookup").(bool)),
					},
				},
				ExcludeAttributes: getExcludedAttrs(d.Get("password_exclude_first_name").(bool), d.Get("password_exclude_last_name").(bool)),
				ExcludeUsername:   utils.BoolPtr(d.Get("password_exclude_username").(bool)),
				MinLengthPtr:      utils.Int64Ptr(d.Get("password_min_length").(int)),
				MinLowerCasePtr:   utils.Int64Ptr(d.Get("password_min_lowercase").(int)),
				MinNumberPtr:      utils.Int64Ptr(d.Get("password_min_number").(int)),
				MinSymbolPtr:      utils.Int64Ptr(d.Get("password_min_symbol").(int)),
				MinUpperCasePtr:   utils.Int64Ptr(d.Get("password_min_uppercase").(int)),
			},
			Lockout: &sdk.PasswordPolicyPasswordSettingsLockout{
				AutoUnlockMinutesPtr:            utils.Int64Ptr(d.Get("password_auto_unlock_minutes").(int)),
				MaxAttemptsPtr:                  utils.Int64Ptr(d.Get("password_max_lockout_attempts").(int)),
				ShowLockoutFailures:             utils.BoolPtr(d.Get("password_show_lockout_failures").(bool)),
				UserLockoutNotificationChannels: utils.ConvertInterfaceToStringSet(d.Get("password_lockout_notification_channels")),
			},
		},
		Recovery: &sdk.PasswordPolicyRecoverySettings{
			Factors: &sdk.PasswordPolicyRecoveryFactors{
				OktaCall: &sdk.PasswordPolicyRecoveryFactorSettings{
					Status: d.Get("call_recovery").(string),
				},
				OktaSms: &sdk.PasswordPolicyRecoveryFactorSettings{
					Status: d.Get("sms_recovery").(string),
				},
				OktaEmail: &sdk.PasswordPolicyRecoveryEmail{
					Properties: &sdk.PasswordPolicyRecoveryEmailProperties{
						RecoveryToken: &sdk.PasswordPolicyRecoveryEmailRecoveryToken{
							TokenLifetimeMinutesPtr: utils.Int64Ptr(d.Get("recovery_email_token").(int)),
						},
					},
					Status: d.Get("email_recovery").(string),
				},
				RecoveryQuestion: &sdk.PasswordPolicyRecoveryQuestion{
					Properties: &sdk.PasswordPolicyRecoveryQuestionProperties{
						Complexity: &sdk.PasswordPolicyRecoveryQuestionComplexity{
							MinLengthPtr: utils.Int64Ptr(d.Get("question_min_length").(int)),
						},
					},
					Status: d.Get("question_recovery").(string),
				},
			},
		},
		Delegation: &sdk.PasswordPolicyDelegationSettings{
			Options: &sdk.PasswordPolicyDelegationSettingsOptions{
				SkipUnlock: utils.BoolPtr(d.Get("skip_unlock").(bool)),
			},
		},
	}
	return template
}

func getExcludedAttrs(excludeFirstName, excludeLastName bool) []string {
	var excludedAttrs []string
	if excludeFirstName {
		excludedAttrs = append(excludedAttrs, "firstName")
	}
	if excludeLastName {
		excludedAttrs = append(excludedAttrs, "lastName")
	}
	return excludedAttrs
}
