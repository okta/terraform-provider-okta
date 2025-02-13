package okta

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func resourcePolicyProfileEnrollmentApps() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyProfileEnrollmentAppsCreate,
		ReadContext:   resourcePolicyProfileEnrollmentAppsRead,
		UpdateContext: resourcePolicyProfileEnrollmentAppsUpdate,
		DeleteContext: resourcePolicyProfileEnrollmentAppsDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Description: `Manages Profile Enrollment Policy Apps
~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.
This resource allows you to manage the apps in the Profile Enrollment Policy. 
**Important Notes:** 
 - Default Enrollment Policy can not be used in this resource since it is used as a policy to re-assign apps to when they are unassigned from this one.
 - When re-assigning the app to another policy, please use 'depends_on' in the policy to which the app will be assigned. This is necessary to avoid 
  unexpected behavior, since if the app is unassigned from the policy it is just assigned to the 'Default' one.`,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the enrollment policy.",
			},
			"apps": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of app IDs to be added to this policy",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"default_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the Default Enrollment Policy. This policy is used as a policy to re-assign apps to when they are unassigned from this one",
			},
		},
	}
}

func resourcePolicyProfileEnrollmentAppsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(policyProfileEnrollmentApps)
	}

	err := setDefaultProfileEnrollmentPolicyID(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := d.Get("policy_id").(string)
	apps := convertInterfaceToStringSetNullable(d.Get("apps"))
	client := getOktaClientFromMetadata(meta)

	for i := range apps {
		_, err := client.Application.UpdateApplicationPolicy(ctx, apps[i], policyID)
		if err != nil {
			return diag.Errorf("failed to add an app to the policy, %v", err)
		}
	}
	d.SetId(policyID)
	return nil
}

func resourcePolicyProfileEnrollmentAppsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(policyProfileEnrollmentApps)
	}

	err := setDefaultProfileEnrollmentPolicyID(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	apps, err := listPolicyEnrollmentAppIDs(ctx, getAPISupplementFromMetadata(meta), d.Id())
	if err != nil {
		return diag.Errorf("failed to get list of enrollment policy apps: %v", err)
	}
	_ = d.Set("policy_id", d.Id())
	_ = d.Set("apps", convertStringSliceToSetNullable(apps))
	return nil
}

func resourcePolicyProfileEnrollmentAppsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(policyProfileEnrollmentApps)
	}

	oldApps, newApps := d.GetChange("apps")
	oldSet := oldApps.(*schema.Set)
	newSet := newApps.(*schema.Set)
	appsToAdd := convertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	appsToRemove := convertInterfaceArrToStringArr(oldSet.Difference(newSet).List())

	client := getOktaClientFromMetadata(meta)
	policyID := d.Get("policy_id").(string)

	for i := range appsToAdd {
		_, err := client.Application.UpdateApplicationPolicy(ctx, appsToAdd[i], policyID)
		if err != nil {
			return diag.Errorf("failed to add an app to the policy, %v", err)
		}
	}

	defaultPolicyID := d.Get("default_policy_id").(string)

	for i := range appsToRemove {
		_, err := client.Application.UpdateApplicationPolicy(ctx, appsToRemove[i], defaultPolicyID)
		if err != nil {
			return diag.Errorf("failed to add reassign app to the default policy, %v", err)
		}
	}

	return nil
}

func resourcePolicyProfileEnrollmentAppsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, meta) {
		return resourceOIEOnlyFeatureError(policyProfileEnrollmentApps)
	}

	defaultPolicyID := d.Get("default_policy_id").(string)
	apps := convertInterfaceToStringSetNullable(d.Get("apps"))
	client := getOktaClientFromMetadata(meta)

	for i := range apps {
		_, err := client.Application.UpdateApplicationPolicy(ctx, apps[i], defaultPolicyID)
		if err != nil {
			return diag.Errorf("failed to add reassign app to the default policy, %v", err)
		}
	}

	return nil
}

func listPolicyEnrollmentAppIDs(ctx context.Context, client *sdk.APISupplement, policyID string) ([]string, error) {
	apps, resp, err := client.ListEnrollmentPolicyApps(ctx, policyID, &query.Params{Limit: defaultPaginationLimit})
	if err != nil {
		return nil, err
	}
	appIDs := make([]string, len(apps))
	for i := range apps {
		appIDs[i] = apps[i].Id
	}
	for resp.HasNextPage() {
		var nextApps []*sdk.Application
		resp, err = resp.Next(ctx, &nextApps)
		if err != nil {
			return nil, err
		}
		for i := range nextApps {
			appIDs = append(appIDs, nextApps[i].Id)
		}
	}
	return appIDs, nil
}

func setDefaultProfileEnrollmentPolicyID(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	policies, err := findSystemPolicyByType(ctx, meta, sdk.ProfileEnrollmentPolicyType)
	if err != nil {
		return err
	}
	var policy *sdk.Policy
	for _, p := range policies {
		if strings.Contains(p.Name, "Default") {
			policy = p
			break
		}
	}
	if policy == nil {
		return errors.New("cannot find default PROFILE_ENROLLMENT policy")
	}
	policyID := d.Get("policy_id").(string)
	if policyID == policy.Id {
		return errors.New("default enrollment policy cannot be used here, since it is used as a policy to re-assign apps to when they are unassigned from this one")
	}
	_ = d.Set("default_policy_id", policy.Id)
	return nil
}
