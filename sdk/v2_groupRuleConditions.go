// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type GroupRuleConditions struct {
	Expression *GroupRuleExpression      `json:"expression,omitempty"`
	People     *GroupRulePeopleCondition `json:"people,omitempty"`
}
