package idaas

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

// All profile properties here so we can do a diff against the config to see if any have changed before making the
// request or before erring due to an update on a user that is DEPROVISIONED. Since we have core user props coupled
// with group/user membership a few change requests go out in the Update function.
var profileKeys = []string{
	"city",
	"cost_center",
	"country_code",
	"custom_profile_attributes",
	"department",
	"display_name",
	"division",
	"email",
	"employee_number",
	"first_name",
	"honorific_prefix",
	"honorific_suffix",
	"last_name",
	"locale",
	"login",
	"manager",
	"manager_id",
	"middle_name",
	"mobile_phone",
	"nick_name",
	"organization",
	"postal_address",
	"preferred_language",
	"primary_phone",
	"profile_url",
	"second_email",
	"state",
	"street_address",
	"timezone",
	"title",
	"user_type",
	"zip_code",
	"password",
	"recovery_question",
	"recovery_answer",
}

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				// Supporting id and email based imports
				user, _, err := getOktaClientFromMetadata(meta).User.GetUser(ctx, d.Id())
				if err != nil {
					return nil, err
				}
				d.SetId(user.Id)
				return []*schema.ResourceData{d}, nil
			},
		},
		Description: "Creates an Okta User. This resource allows you to create and configure an Okta User.",
		Schema: map[string]*schema.Schema{
			"skip_roles": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Do not populate user roles information (prevents additional API call)",
				Deprecated:  "Because admin_roles has been removed, this attribute is a no op and will be removed",
			},
			"city": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User city",
			},
			"cost_center": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User cost center",
			},
			"country_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User country code",
			},
			"custom_profile_attributes": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true, // Required for SetNew()
				ValidateDiagFunc: stringIsJSON,
				StateFunc:        utils.NormalizeDataJSON,
				Description:      "JSON formatted custom attributes for a user. It must be JSON due to various types Okta allows.",
				DiffSuppressFunc: utils.NoChangeInObjectFromUnmarshaledJSON,
			},
			"custom_profile_attributes_to_ignore": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of custom_profile_attribute keys that should be excluded from being managed by Terraform. This is useful in situations where specific custom fields may contain sensitive information and should be managed outside of Terraform.",
			},
			"department": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User department",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User display name, suitable to show end users",
			},
			"division": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User division",
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User primary email address",
			},
			"employee_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User employee number",
			},
			"first_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User first name",
			},
			"honorific_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User honorific prefix",
			},
			"honorific_suffix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User honorific suffix",
			},
			"last_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User last name",
			},
			"locale": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User default location",
			},
			"login": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User Okta login",
			},
			"manager": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Manager of User",
			},
			"manager_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Manager ID of User",
			},
			"middle_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User middle name",
			},
			"mobile_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User mobile phone number",
			},
			"nick_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User nickname",
			},
			"organization": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User organization",
			},
			"postal_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User mailing address",
			},
			"preferred_language": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User preferred language",
			},
			"primary_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User primary phone number",
			},
			"profile_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User online profile (web page)",
			},
			"second_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User secondary email address, used for account recovery",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User state or region",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User profile property. Valid values are `ACTIVE`, `DEPROVISIONED`, `STAGED`, `SUSPENDED`. Default: `ACTIVE`",
				Default:     StatusActive,
				// ignore diff changing to ACTIVE if state is set to PROVISIONED or PASSWORD_EXPIRED
				// since this is a similar status in Okta terms
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == UserStatusProvisioned && new == StatusActive || old == UserStatusPasswordExpired && new == StatusActive
				},
			},
			"raw_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The raw status of the User in Okta - (status is mapped)",
			},
			"street_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User street address",
			},
			"timezone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User default timezone",
			},
			"title": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User title",
			},
			"user_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User employee type",
			},
			"zip_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User zipcode or postal code",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "User Password",
			},
			"expire_password_on_create": {
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      false,
				Description:  "If set to `true`, the user will have to change the password at the next login. This property will be used when user is being created and works only when `password` field is set. Default: `false`",
				RequiredWith: []string{"password"},
			},
			"password_inline_hook": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Specifies that a Password Import Inline Hook should be triggered to handle verification of the user's password the first time the user logs in. This allows an existing password to be imported into Okta directly from some other store. When updating a user with a password hook the user must be in the `STAGED` status. The `password` field should not be specified when using Password Import Inline Hook.",
				ConflictsWith: []string{"password", "password_hash"},
			},
			"old_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Old User Password. Should be only set in case the password was not changed using the provider. fter successful password change this field should be removed and `password` field should be used for further changes.",
			},
			"recovery_question": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User Password Recovery Question",
			},
			"recovery_answer": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "User Password Recovery Answer",
			},
			"realm_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The Realm ID to associate the user with",
			},
			// lintignore:S018
			"password_hash": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Description: "Specifies a hashed password to import into Okta.",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldHash, newHash := d.GetChange("password_hash")
					if oldHash != nil && newHash != nil && len(oldHash.(*schema.Set).List()) > 0 && len(newHash.(*schema.Set).List()) > 0 {
						oh := oldHash.(*schema.Set).List()[0].(map[string]interface{})
						nh := newHash.(*schema.Set).List()[0].(map[string]interface{})
						return reflect.DeepEqual(oh, nh)
					}
					return new == "" || old == new
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Description: "The algorithm used to generate the hash using the password",
							Type:        schema.TypeString,
							Required:    true,
						},
						"work_factor": {
							Description: "Governs the strength of the hash and the time required to compute it. Only required for BCRYPT algorithm",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"salt": {
							Description: "Only required for salted hashes",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"salt_order": {
							Description: "Specifies whether salt was pre- or postfixed to the password before hashing",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"value": {
							Description: "For SHA-512, SHA-256, SHA-1, MD5, This is the actual base64-encoded hash of the password (and salt, if used). This is the " +
								"Base64 encoded value of the SHA-512/SHA-256/SHA-1/MD5 digest that was computed by either pre-fixing or post-fixing the salt to the " +
								"password, depending on the saltOrder. If a salt was not used in the source system, then this should just be the the Base64 encoded " +
								"value of the password's SHA-512/SHA-256/SHA-1/MD5 digest. For BCRYPT, This is the actual radix64-encoded hashed password.",
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},

		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, v interface{}) error {
			filteredCustomAttributes := utils.ConvertInterfaceToStringSet(d.Get("custom_profile_attributes_to_ignore"))
			if len(filteredCustomAttributes) == 0 {
				return nil
			}

			oldAttrs, newAttrs := d.GetChange("custom_profile_attributes")
			var oldAttrsMap map[string]interface{}
			_ = json.Unmarshal([]byte(oldAttrs.(string)), &oldAttrsMap)
			var newAttrsMap map[string]interface{}
			_ = json.Unmarshal([]byte(newAttrs.(string)), &newAttrsMap)

			if d.Id() == "" {
				// This is a new user resource. In this case, we only have new values. We'll filter any
				// values for newly created resources as this is a rare case. If one specifies
				// `custom_profile_attributes_to_filter` and then additionally includes those fields
				// as specified in the initial resource creation, we'll simply ignore them.

				for k := range newAttrsMap {
					if utils.Contains(filteredCustomAttributes, k) {
						delete(newAttrsMap, k)
					}
				}
			} else {
				// We are updating. We've already done a read from the server so the old value will now contain
				// correct values. Thus, we update `custom_profile_attributes` with the filtered attributes
				// from the current old value.

				for k, v := range oldAttrsMap {
					if utils.Contains(filteredCustomAttributes, k) {
						newAttrsMap[k] = v
					}
				}
			}

			customProfileAttributes, _ := json.Marshal(newAttrsMap)
			d.SetNew("custom_profile_attributes", string(customProfileAttributes))

			return nil
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("creating user", "login", d.Get("login").(string))
	profile := populateUserProfile(d)
	qp := query.NewQueryParams()

	// setting activate to false on user creation will leave the user with a status of STAGED
	if d.Get("status").(string) == UserStatusStaged {
		qp = query.NewQueryParams(query.WithActivate(false))
	}

	uc := &sdk.UserCredentials{
		Password: &sdk.PasswordCredential{
			Value: d.Get("password").(string),
			Hash:  buildPasswordCredentialHash(d.Get("password_hash")),
		},
	}
	pih := d.Get("password_inline_hook").(string)
	if pih != "" {
		uc.Password = &sdk.PasswordCredential{
			Hook: &sdk.PasswordCredentialHook{
				Type: pih,
			},
		}
	}
	recoveryQuestion := d.Get("recovery_question").(string)
	recoveryAnswer := d.Get("recovery_answer").(string)
	if recoveryQuestion != "" {
		uc.RecoveryQuestion = &sdk.RecoveryQuestionCredential{
			Question: recoveryQuestion,
			Answer:   recoveryAnswer,
		}
	}

	userBody := sdk.CreateUserRequest{
		Profile:     profile,
		Credentials: uc,
	}

	if realmId, ok := d.GetOk("realm_id"); ok {
		userBody.RealmId = utils.StringPtr(realmId.(string))
	}

	client := getOktaClientFromMetadata(meta)
	user, _, err := client.User.CreateUser(ctx, userBody, qp)
	if err != nil {
		return diag.Errorf("failed to create user: %v", err)
	}
	// set the user id into state before setting roles and status in case they fail
	d.SetId(user.Id)

	// status changing can only happen after user is created as well
	if d.Get("status").(string) == UserStatusSuspended || d.Get("status").(string) == UserStatusDeprovisioned {
		err := updateUserStatus(meta, ctx, user.Id, d.Get("status").(string), client)
		if err != nil {
			return diag.Errorf("failed to update user status: %v", err)
		}
	}

	expire, ok := d.GetOk("expire_password_on_create")
	if ok && expire.(bool) {
		_, _, err = getOktaClientFromMetadata(meta).User.ExpirePassword(ctx, user.Id)
		if err != nil {
			return diag.Errorf("failed to expire user's password: %v", err)
		}
	}

	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceUserReadFilterCustomAttributes(ctx, d, meta, []string{})
}

func resourceUserReadFilterCustomAttributes(ctx context.Context, d *schema.ResourceData, meta interface{}, filteredCustomAttributes []string) diag.Diagnostics {
	logger(meta).Info("reading user", "id", d.Id())
	client := getOktaClientFromMetadata(meta)
	user, resp, err := client.User.GetUser(ctx, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get user: %v", err)
	}
	if user == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("raw_status", user.Status)
	rawMap := flattenUser(user, filteredCustomAttributes)
	err = utils.SetNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set user's properties: %v", err)
	}

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("updating user", "id", d.Id())
	status := d.Get("status").(string)
	statusChange := d.HasChange("status")

	if status == UserStatusStaged && statusChange {
		return diag.Errorf("Okta will not allow a user to be updated to STAGED. Can set to STAGED on user creation only")
	}

	// There are a few requests here so just making sure the state gets updated per successful downstream change
	userChange := hasProfileChange(d)
	realmChange := d.HasChange("realm_id")
	passwordChange := d.HasChange("password")
	passwordHashChange := d.HasChange("password_hash")
	passwordHookChange := d.HasChange("password_inline_hook")
	recoveryQuestionChange := d.HasChange("recovery_question")
	recoveryAnswerChange := d.HasChange("recovery_answer")

	client := getOktaClientFromMetadata(meta)
	if passwordChange {
		user, _, err := client.User.GetUser(ctx, d.Id())
		if err != nil {
			return diag.Errorf("failed to get user: %v", err)
		}
		if user.Status == "PROVISIONED" {
			return diag.Errorf("can not change password for provisioned user, the activation workflow should be " +
				"finished first. Please, check this diagram https://developer.okta.com/docs/reference/api/users/#user-status for more clarity.")
		}
	}

	// run the update status func first so a user that was previously deprovisioned
	// can be updated further if it's status changed in it's terraform configs
	if statusChange {
		err := updateUserStatus(meta, ctx, d.Id(), status, client)
		if err != nil {
			return diag.Errorf("failed to update user status: %v", err)
		}
		_ = d.Set("status", status)
	}

	if status == UserStatusDeprovisioned && userChange {
		return diag.Errorf("Only the status of a DEPROVISIONED user can be updated, we detected other change")
	}

	if userChange || realmChange || passwordHashChange || passwordHookChange {
		profile := populateUserProfile(d)
		userBody := sdk.User{
			Profile: profile,
		}
		if passwordHashChange {
			userBody.Credentials = &sdk.UserCredentials{
				Password: &sdk.PasswordCredential{
					Hash: buildPasswordCredentialHash(d.Get("password_hash")),
				},
			}
		}
		pih := d.Get("password_inline_hook").(string)
		if passwordHookChange && pih != "" {
			userBody.Credentials = &sdk.UserCredentials{
				Password: &sdk.PasswordCredential{
					Hook: &sdk.PasswordCredentialHook{
						Type: pih,
					},
				},
			}
		}

		if realmId, ok := d.GetOk("realm_id"); ok {
			userBody.RealmId = utils.StringPtr(realmId.(string))
		}

		_, _, err := client.User.UpdateUser(ctx, d.Id(), userBody, nil)
		if err != nil {
			return diag.Errorf("failed to update user: %v", err)
		}
	}

	if passwordChange {
		oldPassword, newPassword := d.GetChange("password")
		old, oldPasswordExist := d.GetOk("old_password")
		if oldPasswordExist {
			oldPassword = old
		}
		if oldPasswordExist {
			op := &sdk.PasswordCredential{
				Value: oldPassword.(string),
			}
			np := &sdk.PasswordCredential{
				Value: newPassword.(string),
			}
			npr := &sdk.ChangePasswordRequest{
				OldPassword: op,
				NewPassword: np,
			}
			_, _, err := client.User.ChangePassword(ctx, d.Id(), *npr, nil)
			if err != nil {
				return diag.Errorf("failed to update user's password: %v", err)
			}
		}
		if !oldPasswordExist {
			password, _ := newPassword.(string)
			user := sdk.User{
				Credentials: &sdk.UserCredentials{
					Password: &sdk.PasswordCredential{
						Value: password,
					},
				},
			}
			_, _, err := client.User.UpdateUser(ctx, d.Id(), user, nil)
			if err != nil {
				return diag.Errorf("failed to set user's password: %v", err)
			}
		}
	}

	if recoveryQuestionChange || recoveryAnswerChange {
		nuc := &sdk.UserCredentials{
			Password: &sdk.PasswordCredential{
				Value: d.Get("password").(string),
			},
			RecoveryQuestion: &sdk.RecoveryQuestionCredential{
				Question: d.Get("recovery_question").(string),
				Answer:   d.Get("recovery_answer").(string),
			},
		}
		_, _, err := client.User.ChangeRecoveryQuestion(ctx, d.Id(), *nuc)
		if err != nil {
			return diag.Errorf("failed to change user's password recovery question: %v", err)
		}
	}

	filteredCustomAttributes := utils.ConvertInterfaceToStringSet(d.Get("custom_profile_attributes_to_ignore"))

	return resourceUserReadFilterCustomAttributes(ctx, d, meta, filteredCustomAttributes)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	logger(meta).Info("deleting user", "id", d.Id())
	err := EnsureUserDelete(ctx, d.Id(), d.Get("status").(string), getOktaClientFromMetadata(meta))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func buildPasswordCredentialHash(rawPasswordHash interface{}) *sdk.PasswordCredentialHash {
	if rawPasswordHash == nil || len(rawPasswordHash.(*schema.Set).List()) == 0 {
		return nil
	}
	passwordHash := rawPasswordHash.(*schema.Set).List()
	hash := passwordHash[0].(map[string]interface{})
	wf, _ := hash["work_factor"].(int)
	h := &sdk.PasswordCredentialHash{
		Algorithm:     hash["algorithm"].(string),
		Value:         hash["value"].(string),
		WorkFactorPtr: utils.Int64Ptr(wf),
	}
	h.Salt, _ = hash["salt"].(string)
	h.SaltOrder, _ = hash["salt_order"].(string)
	return h
}

// Checks whether any profile keys have changed, this is necessary since the profile is not nested. Also, necessary
// to give a sensible user readable error when they attempt to update a DEPROVISIONED user. Previously
// this error always occurred when you set a user's status to DEPROVISIONED.
func hasProfileChange(d *schema.ResourceData) bool {
	for _, k := range profileKeys {
		if d.HasChange(k) {
			return true
		}
	}
	return false
}

func EnsureUserDelete(ctx context.Context, id, status string, client *sdk.Client) error {
	// only deprovisioned users can be deleted fully from okta
	// make two passes on the user if they aren't deprovisioned already to deprovision them first
	passes := 2
	if status == UserStatusDeprovisioned {
		passes = 1
	}
	for i := 0; i < passes; i++ {
		_, err := client.User.DeactivateOrDeleteUser(ctx, id, nil)
		if err != nil {
			return fmt.Errorf("failed to deprovision or delete user from Okta: %v", err)
		}
	}
	return nil
}

func mapStatus(currentStatus string) string {
	// PASSWORD_EXPIRED and RECOVERY are effectively ACTIVE for our purposes
	if currentStatus == UserStatusPasswordExpired || currentStatus == UserStatusRecovery {
		return StatusActive
	}
	return currentStatus
}
