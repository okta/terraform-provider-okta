package okta

import (
	"context"
	"errors"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
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

func buildAppSignOnPoilicy(d *schema.ResourceData) *okta.AccessPolicy {
	return &okta.AccessPolicy{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        "ACCESS_POLICY",
	}
}

func resourceAppSignOnPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("creating authentication policy", "name", d.Get("name").(string))
	policyToCreate := buildAppSignOnPoilicy(d)

	oktaClient := getOktaClientFromMetadata(m)

	responsePolicy, _, err := oktaClient.Policy.CreatePolicy(ctx, policyToCreate, nil)
	if err != nil {
		return diag.Errorf("failed to create authentication policy: %v", err)
	}
	id := responsePolicy.(*okta.AccessPolicy).Id
	d.SetId(id)
	return resourceAppSignOnPolicyRead(ctx, d, m)
}

func resourceAppSignOnPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	logger(m).Info("updating authentication policy", "id", d.Id(), "name", d.Get("name").(string))
	policyToUpdate := buildAppSignOnPoilicy(d)
	_, _, err := getOktaClientFromMetadata(m).Policy.UpdatePolicy(ctx, d.Id(), policyToUpdate)
	if err != nil {
		return diag.Errorf("failed to update authentication policy: %v", err)
	}
	return resourceAppSignOnPolicyRead(ctx, d, m)
}

func resourceAppSignOnPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	/**
		1. find the default app policy
		2. assign the default policy to all apps using the current policy (the one to delete)
		3. delete the policy
	**/

	client := getOktaClientFromMetadata(m)
	qp := query.NewQueryParams()
	qp.Type = "ACCESS_POLICY"
	policies, _, err := client.Policy.ListPolicies(ctx, qp)
	if err != nil {
		return diag.Errorf("failed delete authentication policy: %v", err)
	}

	// find the default policy
	var defaultPolicy *okta.Policy
	for _, p := range policies {

		v := p.(*okta.Policy)
		if v.Name == "Default Policy" && *v.System {
			defaultPolicy = v
		}

	}
	if defaultPolicy == nil {
		return diag.Errorf("failed delete authentication policy: %v", errors.New("no default policy found"))
	}

	clients, listErr := listApplications(ctx, client)
	if listErr != nil {
		// delete the policy right away

		_, err = client.Policy.DeletePolicy(ctx, d.Id())
		if err != nil {
			return diag.Errorf("failed delete authentication policy: %v", err)
		}
		return nil
	} else {
		// first unassign the policy and delete it then
		// assign the default app policy to all clients using the current policy
		for _, c := range clients {
			app := c.(*okta.Application)
			accessPolicy := linksValue(app.Links, "accessPolicy", "href")
			if accessPolicy == "" {
				return diag.Errorf("app does not support sign-on policy or this feature is not available")
			}
			if path.Base(accessPolicy) == d.Id() {
				// check if client still exists
				a := okta.NewApplication()
				_, _, getErr := client.Application.GetApplication(ctx, app.Id, a, nil)
				if getErr == nil {
					// only perform the update if the app exists
					_, updateErr := client.Application.UpdateApplicationPolicy(ctx, app.Id, defaultPolicy.Id)
					if updateErr != nil {
						return diag.Errorf("failed to assign default policy '%v' to app %v: %v", defaultPolicy.Id, app.Id, updateErr)
					}
				}

			}
		}

		_, err = client.Policy.DeletePolicy(ctx, d.Id())
		if err != nil {
			return diag.Errorf("failed delete authentication policy: %v", err)
		}
		return nil
	}
}

func listApplications(ctx context.Context, client *okta.Client) ([]okta.App, error) {
	var resClients []okta.App

	clients, resp, err := client.Application.ListApplications(ctx, &query.Params{Limit: defaultPaginationLimit})
	if err != nil {
		return nil, err
	}
	for {
		resClients = append(resClients, clients...)
		if resp.HasNextPage() {
			resp, err = resp.Next(ctx, &clients)
			if err != nil {
				return nil, err
			}
			continue
		} else {
			break
		}
	}
	return resClients, nil
}
