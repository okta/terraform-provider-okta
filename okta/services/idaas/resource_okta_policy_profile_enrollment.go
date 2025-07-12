package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourcePolicyProfileEnrollment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyProfileEnrollmentCreate,
		ReadContext:   resourcePolicyProfileEnrollmentRead,
		UpdateContext: resourcePolicyProfileEnrollmentUpdate,
		DeleteContext: resourcePolicyProfileEnrollmentDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Description: `Creates a Profile Enrollment Policy
		
~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.
This resource allows you to create and configure a Profile Enrollment Policy.`,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the policy",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status of the policy",
				Default:     StatusActive,
			},
		},
	}
}

func resourcePolicyProfileEnrollmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSPolicyProfileEnrollment)
	}

	policy, _, err := getAPISupplementFromMetadata(meta).CreatePolicy(ctx, buildPolicyProfileEnrollment(d))
	if err != nil {
		return diag.Errorf("failed to create profile enrollment policy: %v", err)
	}
	d.SetId(policy.Id)
	status, ok := d.GetOk("status")
	if ok && status.(string) == StatusInactive {
		_, err = getOktaClientFromMetadata(meta).Policy.DeactivatePolicy(ctx, policy.Id)
		if err != nil {
			return diag.Errorf("failed to deactivate profile enrollment policy: %v", err)
		}
	}
	return resourcePolicyProfileEnrollmentRead(ctx, d, meta)
}

func resourcePolicyProfileEnrollmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSPolicyProfileEnrollment)
	}

	policy, err := getPolicy(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to get profile enrollment policy: %v", err)
	}
	if policy == nil {
		return nil
	}
	_ = d.Set("name", policy.Name)
	_ = d.Set("status", policy.Status)
	return nil
}

func resourcePolicyProfileEnrollmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSPolicyProfileEnrollment)
	}

	_, _, err := getAPISupplementFromMetadata(meta).UpdatePolicy(ctx, d.Id(), buildPolicyProfileEnrollment(d))
	if err != nil {
		return diag.Errorf("failed to update profile enrollment policy: %v", err)
	}
	oldStatus, newStatus := d.GetChange("status")
	if oldStatus != newStatus {
		if newStatus == StatusActive {
			_, err = getOktaClientFromMetadata(meta).Policy.ActivatePolicy(ctx, d.Id())
		} else {
			_, err = getOktaClientFromMetadata(meta).Policy.DeactivatePolicy(ctx, d.Id())
		}
		if err != nil {
			return diag.Errorf("failed to change profile enrollment policy status: %v", err)
		}
	}
	return resourcePolicyProfileEnrollmentRead(ctx, d, meta)
}

func resourcePolicyProfileEnrollmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if providerIsClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(resources.OktaIDaaSPolicyProfileEnrollment)
	}

	err := deletePolicy(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to delete profile enrollment policy: %v", err)
	}
	return nil
}

// build profile enrollment policy from schema data
func buildPolicyProfileEnrollment(d *schema.ResourceData) sdk.SdkPolicy {
	policy := sdk.ProfileEnrollmentPolicy()
	policy.Name = d.Get("name").(string)
	policy.Status = d.Get("status").(string)
	return policy
}
