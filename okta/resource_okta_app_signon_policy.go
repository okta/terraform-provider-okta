package okta

import (
	"context"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceAppSignOnPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppSignOnPolicyCreate,
		ReadContext:   resourceAppSignOnPolicyRead,
		UpdateContext: resourceAppSignOnPolicyUpdate,
		DeleteContext: resourceAppSignOnPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: ` Manages a sign-on policy.
		
~> **WARNING:** This feature is only available as a part of the Okta Identity
Engine (OIE) and ***is not*** compatible with Classic orgs. Authentication
policies for applications in a Classic org can only be modified in the Admin UI,
there isn't a public API for this. Therefore the Okta Terraform Provider does
not support this resource for Classic orgs. [Contact
support](mailto:dev-inquiries@okta.com) for further information.
This resource allows you to create and configure a sign-on policy for the
application. Inside the product a sign-on policy is referenced as an
_authentication policy_, in the public API the policy is of type
['ACCESS_POLICY'](https://developer.okta.com/docs/reference/api/policy/#policy-object).
A newly created app's sign-on policy will always contain the default
authentication policy unless one is assigned via 'authentication_policy' in the
app resource. At the API level the default policy has system property value of
true.
~> **WARNING:** When this policy is destroyed any other applications that
associate the policy as their authentication policy will be reassigned to the
default/system access policy.`,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Description of the policy.",
			},
		},
	}
}

func buildAccessPolicy(d *schema.ResourceData) sdk.Policies {
	accessPolicy := sdk.NewAccessPolicy()
	accessPolicy.Name = d.Get("name").(string)
	accessPolicy.Description = d.Get("description").(string)
	return accessPolicy
}

func resourceAppSignOnPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, m) {
		return resourceOIEOnlyFeatureError(appSignOnPolicy)
	}

	logger(m).Info("creating authentication policy", "name", d.Get("name").(string))
	policy := buildAccessPolicy(d)
	oktaClient := getOktaClientFromMetadata(m)

	responsePolicy, _, err := oktaClient.Policy.CreatePolicy(ctx, policy, nil)
	if err != nil {
		return diag.Errorf("failed to create authentication policy: %v", err)
	}
	id := responsePolicy.(*sdk.AccessPolicy).Id
	d.SetId(id)
	return resourceAppSignOnPolicyRead(ctx, d, m)
}

func resourceAppSignOnPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, m) {
		return resourceOIEOnlyFeatureError(appSignOnPolicy)
	}

	logger(m).Info("reading authentication policy", "id", d.Id(), "name", d.Get("name").(string))
	policy := &sdk.Policy{}
	authenticationPolicy, resp, err := getOktaClientFromMetadata(m).Policy.GetPolicy(ctx, d.Id(), policy, nil)
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get authentication policy: %v", err)
	}
	if authenticationPolicy == nil {
		d.SetId("")
		return nil
	}
	policyFromServer := authenticationPolicy.(*sdk.Policy)
	d.SetId(policyFromServer.Id)
	d.Set("name", policyFromServer.Name)
	d.Set("description", policyFromServer.Description)
	return nil
}

func resourceAppSignOnPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, m) {
		return resourceOIEOnlyFeatureError(appSignOnPolicy)
	}

	logger(m).Info("updating authentication policy", "id", d.Id(), "name", d.Get("name").(string))
	policyToUpdate := buildAccessPolicy(d)
	_, _, err := getOktaClientFromMetadata(m).Policy.UpdatePolicy(ctx, d.Id(), policyToUpdate)
	if err != nil {
		return diag.Errorf("failed to update authentication policy: %v", err)
	}
	return resourceAppSignOnPolicyRead(ctx, d, m)
}

// resourceAppSignOnPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
func resourceAppSignOnPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, m) {
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
