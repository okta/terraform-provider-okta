package idaas

import (
	"context"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func createOrUpdateAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, m interface{}, appId string) error {
	// Check if authentication policy operations should be skipped
	// Only check if the field exists in the schema to avoid errors
	if _, exists := d.GetOk("skip_authentication_policy"); exists {
		if skip, ok := d.GetOk("skip_authentication_policy"); ok && skip.(bool) {
			return nil
		}
	}

	raw, ok := d.GetOk("authentication_policy")
	if !ok {
		return assignDefaultAuthenticationPolicy(ctx, m, appId)
	}
	policyId := raw.(string)
	_, err := getOktaClientFromMetadata(m).Application.UpdateApplicationPolicy(ctx, appId, policyId)
	return err
}

func setAuthenticationPolicy(ctx context.Context, m interface{}, d *schema.ResourceData, links interface{}) {
	// Check if authentication policy operations should be skipped
	// Only check if the field exists in the schema to avoid errors
	if _, exists := d.GetOk("skip_authentication_policy"); exists {
		if skip, ok := d.GetOk("skip_authentication_policy"); ok && skip.(bool) {
			return
		}
	}
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
	client := getOktaClientFromMetadata(m)
	_, err = client.Application.UpdateApplicationPolicy(ctx, appId, policy.Id)
	return err
}
