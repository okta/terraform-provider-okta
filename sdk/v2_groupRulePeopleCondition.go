// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type GroupRulePeopleCondition struct {
	Groups *GroupRuleGroupCondition `json:"groups,omitempty"`
	Users  *GroupRuleUserCondition  `json:"users,omitempty"`
}
