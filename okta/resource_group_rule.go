package okta

import (
	"net/http"

	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupRuleCreate,
		Exists: resourceGroupRuleExists,
		Read:   resourceGroupRuleRead,
		Update: resourceGroupRuleUpdate,
		Delete: resourceGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: addPeopleAssignments(map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"group_assignments": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     schema.TypeString,
			},
			"expression_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"expression_value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"status": statusSchema,
		}),
	}
}

func buildGroupRule(d *schema.ResourceData) *okta.GroupRule {
	return &okta.GroupRule{
		Actions: &okta.GroupRuleAction{
			AssignUserToGroups: &okta.GroupRuleGroupAssignment{
				GroupIds: convertInterfaceToStringSet(d.Get("group_assignments")),
			},
		},
		Conditions: &okta.GroupRuleConditions{
			Expression: &okta.GroupRuleExpression{
				Type:  d.Get("expression_type").(string),
				Value: d.Get("expression_value").(string),
			},
			People: getPeopleConditions(d),
		},
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
		Type:   "group_rule",
	}
}

func resourceGroupRuleCreate(d *schema.ResourceData, m interface{}) error {
	groupRule := buildGroupRule(d)
	responseGroupRule, _, err := getOktaClientFromMetadata(m).Group.CreateRule(*groupRule)
	if err != nil {
		return err
	}

	d.SetId(responseGroupRule.Id)

	return resourceGroupRuleRead(d, m)
}

func resourceGroupRuleExists(d *schema.ResourceData, m interface{}) (bool, error) {
	g, err := fetchGroupRule(d, m)

	return err == nil && g != nil, err
}

func resourceGroupRuleRead(d *schema.ResourceData, m interface{}) error {
	g, err := fetchGroupRule(d, m)
	if err != nil {
		return err
	}

	d.Set("name", g.Name)
	d.Set("type", g.Type)
	d.Set("status", g.Status)
	d.Set("expression_type", g.Status)
	d.Set("expression_value", g.Status)

	err = setPeopleAssignments(d, g.Conditions.People)
	if err != nil {
		return err
	}

	return setNonPrimitives(d, map[string]interface{}{
		"group_assignments": g.Actions.AssignUserToGroups.GroupIds,
	})
}

func resourceGroupRuleUpdate(d *schema.ResourceData, m interface{}) error {
	rule := buildGroupRule(d)
	_, _, err := getOktaClientFromMetadata(m).Group.UpdateRule(d.Id(), *rule)
	if err != nil {
		return err
	}

	return resourceGroupRuleRead(d, m)
}

func resourceGroupRuleDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getOktaClientFromMetadata(m).Group.DeleteRule(d.Id(), &query.Params{})

	return err
}

func fetchGroupRule(d *schema.ResourceData, m interface{}) (*okta.GroupRule, error) {
	g, resp, err := getOktaClientFromMetadata(m).Group.GetRule(d.Id())

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	return g, err
}
