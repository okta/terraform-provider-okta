package idaas

import (
	"context"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func createOrUpdateAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, m interface{}, appId string) error {
	raw, ok := d.GetOk("authentication_policy")
	if !ok {
		return assignDefaultAuthenticationPolicy(ctx, m, appId)
	}
	policyId := raw.(string)
	_, err := GetOktaClientFromMetadata(m).Application.UpdateApplicationPolicy(ctx, appId, policyId)
	return err
}

func setAuthenticationPolicy(ctx context.Context, m interface{}, d *schema.ResourceData, links interface{}) {
	// setAuthenticationPolicy by default, switching from optional to optional and computed
	if providerIsClassicOrg(ctx, m) {
		return
	}
	accessPolicy := utils.LinksValue(links, "accessPolicy", "href")
	if accessPolicy != "" {
		d.Set("authentication_policy", path.Base(accessPolicy))
	}
}

func assignDefaultAuthenticationPolicy(ctx context.Context, m interface{}, appId string) error {
	// Apps in OIE orgs have a default authentication / access policy that is
	// type ACCESS_POLICY. Apps in classic orgs do not have an access policy
	// accessible through the public API. Only by hand in the Admin UI.
	// https://developer.okta.com/docs/reference/api/policy/#policy-object
	if providerIsClassicOrg(ctx, m) {
		return nil
	}

	policy, err := findDefaultAccessPolicy(ctx, m)
	if err != nil {
		return err
	}
	client := GetOktaClientFromMetadata(m)
	_, err = client.Application.UpdateApplicationPolicy(ctx, appId, policy.Id)
	return err
}
