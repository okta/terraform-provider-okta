package sdk

type GroupRulePeopleCondition struct {
	Groups *GroupRuleGroupCondition `json:"groups,omitempty"`
	Users  *GroupRuleUserCondition  `json:"users,omitempty"`
}
