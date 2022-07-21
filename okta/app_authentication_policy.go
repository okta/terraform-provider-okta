package okta

import (
	"context"
	"errors"
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
	client := getOktaClientFromMetadata(m)
	qp := query.NewQueryParams()
	qp.Type = "ACCESS_POLICY"
	policies, _, err := client.Policy.ListPolicies(ctx, qp)
	if err != nil {
		return fmt.Errorf("failed delete authentication policy: %v", err)
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
		return errors.New("no default policy found")
	}
	_, err = client.Application.UpdateApplicationPolicy(ctx, appId, defaultPolicy.Id)
	return err
}
