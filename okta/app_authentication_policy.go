package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	// Apps in OIE orgs have a default authentication / access policy that is
	// type ACCESS_POLICY. Apps in classic orgs do not have an access policy
	// accessible through the public API. Only by hand in the Admin UI.
	// https://developer.okta.com/docs/reference/api/policy/#policy-object
	if config, ok := m.(*Config); ok && config.classicOrg {
		return nil
	}

	policy, err := findDefaultAccessPolicy(ctx, m)
	if err != nil {
		return err
	}
	client := getOktaClientFromMetadata(m)
	_, err = client.Application.UpdateApplicationPolicy(ctx, appId, policy.Id)
	return err
}
