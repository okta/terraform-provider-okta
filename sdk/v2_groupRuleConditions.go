package sdk

type GroupRuleConditions struct {
	Expression *GroupRuleExpression      `json:"expression,omitempty"`
	People     *GroupRulePeopleCondition `json:"people,omitempty"`
}
