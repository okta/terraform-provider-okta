package okta

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func getPolicyFactorSchema(key string) map[string]*schema.Schema {
	// These are primitives to allow defaulting. Terraform still does not support aggregate defaults.
	return map[string]*schema.Schema{
		key: &schema.Schema{
			Optional: true,
			Type:     schema.TypeMap,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enroll": &schema.Schema{
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "OPTIONAL",
						ValidateFunc: validation.StringInSlice([]string{"NOT_ALLOWED", "OPTIONAL", "REQUIRED"}, false),
						Description:  "Requirements for use-initiated enrollment.",
					},
					"consent_type": &schema.Schema{
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "NONE",
						ValidateFunc: validation.StringInSlice([]string{"NONE", "TERMS_OF_SERVICE"}, false),
						Description:  "User consent type required before enrolling in the factor: NONE or TERMS_OF_SERVICE.",
					},
				},
			},
		},
	}
}

var factorProviders = []string{
	"duo",
	"fido_u2f",
	"fido_webauthn",
	"google_otp",
	"okta_call",
	"okta_otp",
	"okta_password",
	"okta_push",
	"okta_question",
	"okta_sms",
	"rsa_token",
	"symantec_vip",
	"yubikey_token",
}

func buildFactorProviders(target map[string]*schema.Schema) map[string]*schema.Schema {
	for _, key := range factorProviders {
		sMap := getPolicyFactorSchema(key)

		for nestedKey, nestedVal := range sMap {
			target[nestedKey] = nestedVal
		}
	}

	return target
}

func resourcePolicyMfa() *schema.Resource {
	return &schema.Resource{
		Exists: resourcePolicyExists,
		Create: resourcePolicyMfaCreate,
		Read:   resourcePolicyMfaRead,
		Update: resourcePolicyMfaUpdate,
		Delete: resourcePolicyMfaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: buildPolicySchema(
			// List of factor provider above, they all follow the same schema
			buildFactorProviders(map[string]*schema.Schema{}),
		),
	}
}

func resourcePolicyMfaCreate(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultPolicy(d); err != nil {
		return err
	}

	policy := buildMfaPolicy(d, m)
	err := createPolicy(d, m, policy)
	if err != nil {
		return err
	}

	return resourcePolicyMfaRead(d, m)
}

func resourcePolicyMfaRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] List Policy %v", d.Get("name").(string))

	policy, err := getPolicy(d, m)

	if policy == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	syncFactor(d, "duo", policy.Settings.Factors.Duo)
	syncFactor(d, "fido_u2f", policy.Settings.Factors.FidoU2f)
	syncFactor(d, "fido_webauthn", policy.Settings.Factors.FidoWebauthn)
	syncFactor(d, "google_otp", policy.Settings.Factors.GoogleOtp)
	syncFactor(d, "okta_call", policy.Settings.Factors.OktaOtp)
	syncFactor(d, "okta_otp", policy.Settings.Factors.OktaOtp)
	syncFactor(d, "okta_password", policy.Settings.Factors.OktaPassword)
	syncFactor(d, "okta_push", policy.Settings.Factors.OktaPush)
	syncFactor(d, "okta_question", policy.Settings.Factors.OktaQuestion)
	syncFactor(d, "okta_sms", policy.Settings.Factors.OktaSms)
	syncFactor(d, "rsa_token", policy.Settings.Factors.RsaToken)
	syncFactor(d, "symantec_vip", policy.Settings.Factors.SymantecVip)
	syncFactor(d, "yubikey_token", policy.Settings.Factors.YubikeyToken)

	return syncPolicyFromUpstream(d, policy)
}

func resourcePolicyMfaUpdate(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultPolicy(d); err != nil {
		return err
	}

	d.Partial(true)
	policy := buildMfaPolicy(d, m)
	err := updatePolicy(d, m, policy)
	if err != nil {
		return err
	}
	d.Partial(false)

	return resourcePolicyMfaRead(d, m)
}

func resourcePolicyMfaDelete(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultPolicy(d); err != nil {
		return err
	}

	log.Printf("[INFO] Delete Policy %v", d.Get("name").(string))
	client := m.(*Config).articulateOktaClient

	_, err := client.Policies.DeletePolicy(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Policy from Okta: %v", err)
	}

	return nil
}

// create or update a password policy
func buildMfaPolicy(d *schema.ResourceData, m interface{}) *articulateOkta.Policy {
	client := getClientFromMetadata(m)

	policy := client.Policies.MfaPolicy()
	policy.Name = d.Get("name").(string)
	policy.Status = d.Get("status").(string)
	policy.Type = mfaPolicyType
	policy.Description = d.Get("description").(string)
	if priority, ok := d.GetOk("priority"); ok {
		policy.Priority = priority.(int)
	}

	policy.Settings = &articulateOkta.PolicySettings{
		Factors: &articulateOkta.Factors{
			Duo:          buildFactorProvider(d, "duo"),
			FidoU2f:      buildFactorProvider(d, "fido_u2f"),
			FidoWebauthn: buildFactorProvider(d, "fido_webauthn"),
			GoogleOtp:    buildFactorProvider(d, "google_otp"),
			OktaCall:     buildFactorProvider(d, "okta_call"),
			OktaOtp:      buildFactorProvider(d, "okta_otp"),
			OktaPassword: buildFactorProvider(d, "okta_password"),
			OktaPush:     buildFactorProvider(d, "okta_push"),
			OktaQuestion: buildFactorProvider(d, "okta_question"),
			OktaSms:      buildFactorProvider(d, "okta_sms"),
			RsaToken:     buildFactorProvider(d, "rsa_token"),
			SymantecVip:  buildFactorProvider(d, "symantec_vip"),
			YubikeyToken: buildFactorProvider(d, "yubikey_token"),
		},
	}
	policy.Conditions = &articulateOkta.PolicyConditions{
		People: getGroups(d),
	}

	return &policy
}

func buildFactorProvider(d *schema.ResourceData, key string) *articulateOkta.FactorProvider {
	consent := d.Get(fmt.Sprintf("%s.consent_type", key)).(string)
	enroll := d.Get(fmt.Sprintf("%s.enroll", key)).(string)

	if consent == "" && enroll == "" {
		return nil
	}

	provider := &articulateOkta.FactorProvider{}

	if consent != "" {
		provider.Consent = articulateOkta.Consent{Type: consent}
	}

	if enroll != "" {
		provider.Enroll = articulateOkta.Enroll{Self: enroll}
	}

	return provider
}

func syncFactor(d *schema.ResourceData, k string, f *articulateOkta.FactorProvider) {
	if f != nil {
		d.Set(fmt.Sprintf("%s.consent_type", k), f.Consent.Type)
		d.Set(fmt.Sprintf("%s.enroll", k), f.Enroll.Self)
	}
}
