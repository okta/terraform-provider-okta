package okta

import (
	"context"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceAppSignOnPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppSignOnPolicyCreate,
		ReadContext:   resourceAppSignOnPolicyRead,
		UpdateContext: resourceAppSignOnPolicyUpdate,
		DeleteContext: resourceAppSignOnPolicyDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy Name",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy Description",
			},
		},
	}
}

func buildAccessPoilicy(d *schema.ResourceData) okta.Policies {
	accessPolicy := okta.NewAccessPolicy()
	accessPolicy.Name = d.Get("name").(string)
	accessPolicy.Description = d.Get("description").(string)
	return accessPolicy
}

func resourceAppSignOnPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(m) {
		return resourceOIEOnlyFeatureError(appSignOnPolicy)
	}

	logger(m).Info("creating authentication policy", "name", d.Get("name").(string))
	var policy okta.Policies
	policy = buildAccessPoilicy(d)

	oktaClient := getOktaClientFromMetadata(m)

	responsePolicy, _, err := oktaClient.Policy.CreatePolicy(ctx, policy, nil)
	if err != nil {
		return diag.Errorf("failed to create authentication policy: %v", err)
	}
	id := responsePolicy.(*okta.AccessPolicy).Id
	d.SetId(id)
	return resourceAppSignOnPolicyRead(ctx, d, m)
}

func resourceAppSignOnPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(m) {
		return resourceOIEOnlyFeatureError(appSignOnPolicy)
	}

	logger(m).Info("reading authentication policy", "id", d.Id(), "name", d.Get("name").(string))
	policy := &okta.Policy{}
	authenticationPolicy, resp, err := getOktaClientFromMetadata(m).Policy.GetPolicy(ctx, d.Id(), policy, nil)
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get authentication policy: %v", err)
	}
	if authenticationPolicy == nil {
		d.SetId("")
		return nil
	}
	policyFromServer := authenticationPolicy.(*okta.Policy)
	d.SetId(policyFromServer.Id)
	d.Set("name", policyFromServer.Name)
	d.Set("description", policyFromServer.Description)
	return nil
}

func resourceAppSignOnPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(m) {
		return resourceOIEOnlyFeatureError(appSignOnPolicy)
	}

	logger(m).Info("updating authentication policy", "id", d.Id(), "name", d.Get("name").(string))
	policyToUpdate := buildAccessPoilicy(d)
	_, _, err := getOktaClientFromMetadata(m).Policy.UpdatePolicy(ctx, d.Id(), policyToUpdate)
	if err != nil {
		return diag.Errorf("failed to update authentication policy: %v", err)
	}
	return resourceAppSignOnPolicyRead(ctx, d, m)
}

// resourceAppSignOnPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
func resourceAppSignOnPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(m) {
		return resourceOIEOnlyFeatureError(appSignOnPolicy)
	}

	// 1. find the default app policy
	// 2. assign default policy to all apps whose authentication policy is the policy about to be deleted
	// 3. delete the policy
	defaultPolicy, err := findDefaultAccessPolicy(ctx, m)
	if err != nil {
		return diag.Errorf("Error finding default access policy: %v", err)
	}

	client := getOktaClientFromMetadata(m)
	apps, err := listApps(ctx, client, nil, defaultPaginationLimit)
	if err != nil {
		return diag.Errorf("failed to list apps in preparation to delete authentication policy: %v", err)
	}

	// assign the default app policy to all clients using the current policy
	for _, app := range apps {
		accessPolicy := linksValue(app.Links, "accessPolicy", "href")
		// ignore apps that don't have an access policy, typically Classic org apps.
		if accessPolicy == "" {
			continue
		}
		// app uses this policy as its access policy, change that back to using the default policy
		if path.Base(accessPolicy) == d.Id() {
			// update the app with the default policy, ignore errors
			_, _ = client.Application.UpdateApplicationPolicy(ctx, app.Id, defaultPolicy.Id)
		}
	}

	// delete will error out if the policy is still associated with apps
	_, err = client.Policy.DeletePolicy(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed delete authentication policy: %v", err)
	}

	return nil
}
