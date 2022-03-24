package okta

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourcePolicyProfileEnrollmentApps() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyProfileEnrollmentAppsCreate,
		ReadContext:   resourcePolicyProfileEnrollmentAppsRead,
		UpdateContext: resourcePolicyProfileEnrollmentAppsUpdate,
		DeleteContext: resourcePolicyProfileEnrollmentAppsDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
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

func resourcePolicyProfileEnrollmentAppsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setDefaultPolicyID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := d.Get("policy_id").(string)
	apps := convertInterfaceToStringSetNullable(d.Get("apps"))
	client := getSupplementFromMetadata(m)

	for i := range apps {
		body := sdk.AddAppToEnrollmentPolicyRequest{
			ResourceType: "APP",
			ResourceId:   apps[i],
		}
		_, _, err := client.AddAppToEnrollmentPolicy(ctx, policyID, body)
		if err != nil {
			return diag.Errorf("failed to add an app to the policy, %v", err)
		}
	}
	d.SetId(policyID)
	return nil
}

func resourcePolicyProfileEnrollmentAppsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setDefaultPolicyID(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	apps, err := listPolicyEnrollmentAppIDs(ctx, getSupplementFromMetadata(m), d.Id())
	if err != nil {
		return diag.Errorf("failed to get list of enrollment policy apps: %v", err)
	}
	_ = d.Set("policy_id", d.Id())
	_ = d.Set("apps", convertStringSliceToSetNullable(apps))
	return nil
}

func resourcePolicyProfileEnrollmentAppsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	oldApps, newApps := d.GetChange("apps")
	oldSet := oldApps.(*schema.Set)
	newSet := newApps.(*schema.Set)
	appsToAdd := convertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	appsToRemove := convertInterfaceArrToStringArr(oldSet.Difference(newSet).List())

	client := getSupplementFromMetadata(m)
	policyID := d.Get("policy_id").(string)

	for i := range appsToAdd {
		body := sdk.AddAppToEnrollmentPolicyRequest{
			ResourceType: "APP",
			ResourceId:   appsToAdd[i],
		}
		_, _, err := client.AddAppToEnrollmentPolicy(ctx, policyID, body)
		if err != nil {
			return diag.Errorf("failed to add an app to the policy, %v", err)
		}
	}

	defaultPolicyID := d.Get("default_policy_id").(string)

	for i := range appsToRemove {
		body := sdk.AddAppToEnrollmentPolicyRequest{
			ResourceType: "APP",
			ResourceId:   appsToRemove[i],
		}
		_, _, err := client.AddAppToEnrollmentPolicy(ctx, defaultPolicyID, body)
		if err != nil {
			return diag.Errorf("failed to add reassign app to the default policy, %v", err)
		}
	}

	return nil
}

func resourcePolicyProfileEnrollmentAppsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	defaultPolicyID := d.Get("default_policy_id").(string)
	apps := convertInterfaceToStringSetNullable(d.Get("apps"))
	client := getSupplementFromMetadata(m)

	for i := range apps {
		body := sdk.AddAppToEnrollmentPolicyRequest{
			ResourceType: "APP",
			ResourceId:   apps[i],
		}
		_, _, err := client.AddAppToEnrollmentPolicy(ctx, defaultPolicyID, body)
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
		var nextApps []*okta.Application
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

func setDefaultPolicyID(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	policy, err := findPolicy(ctx, m, "Default Policy", sdk.ProfileEnrollmentPolicyType)
	if err != nil {
		return err
	}
	policyID := d.Get("policy_id").(string)
	if policyID == policy.Id {
		return errors.New("default enrollment policy cannot be used here, since it is used as a policy to re-assign apps to when they are unassigned from this one")
	}
	_ = d.Set("default_policy_id", policy.Id)
	return nil
}
