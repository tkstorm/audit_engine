package mydb

//规则类型
type RuleType struct {
	id          int
	title       string
	audit_mark  string
	description string
}

//规则
type Rule struct {
	ruleItems []RuleItem
}

//规则项
type RuleItem struct {
}

//审核类型模板
func GetRuleTypes() {

}

//审核规则(1:通过2:驳回3:人工审核)
func GetRules() {

}