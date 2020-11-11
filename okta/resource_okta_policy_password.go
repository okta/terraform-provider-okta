package okta

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func resourcePolicyPassword() *schema.Resource {
	return &schema.Resource{
		Exists: resourcePolicyExists,
		Create: resourcePolicyPasswordCreate,
		Read:   resourcePolicyPasswordRead,
		Update: resourcePolicyPasswordUpdate,
		Delete: resourcePolicyPasswordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: buildPolicySchema(map[string]*schema.Schema{
			"auth_provider": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"OKTA", "ACTIVE_DIRECTORY"}, false),
				Description:  "Authentication Provider: OKTA or ACTIVE_DIRECTORY.",
				Default:      "OKTA",
			},
			"password_min_length": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Minimum password length.",
				Default:     8,
			},
			"password_min_lowercase": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1),
				Description:  "If a password must contain at least one lower case letter: 0 = no, 1 = yes. Default = 1",
				Default:      1,
			},
			"password_min_uppercase": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1),
				Description:  "If a password must contain at least one upper case letter: 0 = no, 1 = yes. Default = 1",
				Default:      1,
			},
			"password_min_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1),
				Description:  "If a password must contain at least one number: 0 = no, 1 = yes. Default = 1",
				Default:      1,
			},
			"password_min_symbol": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1),
				Description:  "If a password must contain at least one symbol (!@#$%^&*): 0 = no, 1 = yes. Default = 1",
				Default:      0,
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{statusActive, statusInactive}, false),
				Description:  "Enable or disable email password recovery: ACTIVE or INACTIVE.",
				Default:      statusActive,
			},
			"recovery_email_token": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Lifetime in minutes of the recovery email token.",
				Default:     60,
			},
			"sms_recovery": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{statusActive, statusInactive}, false),
				Description:  "Enable or disable SMS password recovery: ACTIVE or INACTIVE.",
				Default:      statusInactive,
			},
			"question_recovery": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{statusActive, statusInactive}, false),
				Description:  "Enable or disable security question password recovery: ACTIVE or INACTIVE.",
				Default:      statusActive,
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

func resourcePolicyPasswordCreate(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultPolicy(d); err != nil {
		return err
	}
	log.Printf("[INFO] Creating Policy %v", d.Get("name").(string))
	template := buildPasswordPolicy(d)
	err := createPolicy(d, m, template)
	if err != nil {
		return err
	}

	return resourcePolicyPasswordRead(d, m)
}

func resourcePolicyPasswordRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] List Policy %v", d.Get("name").(string))

	policy, err := getPolicy(d, m)

	if policy == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	// Update with upstream state when it is manually updated from Okta UI or API directly.
	// See https://github.com/terraform-providers/terraform-provider-okta/issues/61
	if policy.Conditions.AuthProvider != nil && policy.Conditions.AuthProvider.Provider != "" {
		_ = d.Set("auth_provider", policy.Conditions.AuthProvider.Provider)
	}

	if policy.Settings != nil {
		err = d.Set("password_lockout_notification_channels", convertStringSetToInterface(policy.Settings.Password.Lockout.UserLockoutNotificationChannels))
		if err != nil {
			return fmt.Errorf("error setting notification channels for resource %s: %s", d.Id(), err)
		}
		_ = d.Set("password_min_length", policy.Settings.Password.Complexity.MinLength)
		_ = d.Set("password_min_lowercase", policy.Settings.Password.Complexity.MinLowerCase)
		_ = d.Set("password_min_uppercase", policy.Settings.Password.Complexity.MinUpperCase)
		_ = d.Set("password_min_number", policy.Settings.Password.Complexity.MinNumber)
		_ = d.Set("password_min_symbol", policy.Settings.Password.Complexity.MinSymbol)
		_ = d.Set("password_exclude_username", policy.Settings.Password.Complexity.ExcludeUsername)
		_ = d.Set("password_dictionary_lookup", policy.Settings.Password.Complexity.Dictionary.Common.Exclude)
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
		_ = d.Set("skip_unlock", policy.Settings.Delegation.Options.SkipUnlock)

		valueMap := map[string]interface{}{}

		excludedAttrs := policy.Settings.Password.Complexity.ExcludeAttributes
		if len(excludedAttrs) > 0 {
			for _, v := range excludedAttrs {
				switch v {
				case "firstName":
					_ = d.Set("password_excluded_first_name", true)
				case "lastName":
					_ = d.Set("password_excluded_last_name", true)
				}
			}
		}
		err = setNonPrimitives(d, valueMap)

		if err != nil {
			return err
		}
	}

	return syncPolicyFromUpstream(d, policy)
}

func resourcePolicyPasswordUpdate(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultPolicy(d); err != nil {
		return err
	}

	log.Printf("[INFO] Update Policy %v", d.Get("name").(string))

	template := buildPasswordPolicy(d)
	err := updatePolicy(d, m, template)
	if err != nil {
		return err
	}

	return resourcePolicyPasswordRead(d, m)
}

func resourcePolicyPasswordDelete(d *schema.ResourceData, m interface{}) error {
	return deletePolicy(d, m)
}

// create or update a password policy
func buildPasswordPolicy(d *schema.ResourceData) sdk.Policy {
	template := sdk.PasswordPolicy()
	template.Name = d.Get("name").(string)
	template.Status = d.Get("status").(string)
	if description, ok := d.GetOk("description"); ok {
		template.Description = description.(string)
	}
	if priority, ok := d.GetOk("priority"); ok {
		template.Priority = int64(priority.(int))
	}
	template.Conditions = &okta.PolicyRuleConditions{
		AuthProvider: &okta.PasswordPolicyAuthenticationProviderCondition{
			Provider: d.Get("auth_provider").(string),
		},
		People: getGroups(d),
	}
	// Okta defaults
	// we add the defaults here & not in the schema map to avoid defaults appearing in the terraform plan diff
	template.Settings = &sdk.PolicySettings{
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
				OktaEmail: &sdk.PasswordPolicyRecoveryEmail{
					Properties: &sdk.PasswordPolicyRecoveryEmailProperties{
						RecoveryToken: &sdk.PasswordPolicyRecoveryEmailRecoveryToken{
							TokenLifetimeMinutes: int64(d.Get("recovery_email_token").(int)),
						},
					},
					Status: d.Get("email_recovery").(string),
				},
				OktaSms: &okta.PasswordPolicyRecoveryFactorSettings{
					Status: d.Get("sms_recovery").(string),
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
