package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func setAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, m interface{}, appId string) error {
	raw, ok := d.GetOk("authentication_policy")
	if !ok {
		return assignDefaultAuthenticationPolicy(ctx, m, appId)
	}
	policyId := raw.(string)
	_, err := getOktaClientFromMetadata(m).Application.UpdateApplicationPolicy(ctx, appId, policyId)
	return err
}

func assignDefaultAuthenticationPolicy(ctx context.Context, m interface{}, appId string) error {
	// Inspecting the ACCESS_POLICY is OIE only https://developer.okta.com/docs/reference/api/policy/#policy-object
	// "Note: The following policy types are available only with the Identity Engine: ACCESS_POLICY or PROFILE_ENROLLMENT."
	// If the org is not OIE return early
	if config, ok := m.(*Config); ok && config.classicOrg {
		return nil
	}

	client := getOktaClientFromMetadata(m)
	qp := query.NewQueryParams()
	qp.Type = "ACCESS_POLICY"
	policies, _, err := client.Policy.ListPolicies(ctx, qp)
	if err != nil {
		return fmt.Errorf("failed delete authentication policy: %v", err)
	}

	// Assign the default policy to the app if the policy exists
	for _, p := range policies {

		v := p.(*okta.Policy)
		if v.Name == "Default Policy" && *v.System {
			_, err = client.Application.UpdateApplicationPolicy(ctx, appId, v.Id)
			return err
		}
	}

	return nil
}
