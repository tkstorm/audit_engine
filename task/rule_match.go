package task

import (
	"fmt"
	"github.com/tkstorm/audit_engine/rabbit"
	"log"
)

const (
	RuleMatched    = 1
	RuleNotMatched = 2
)

//1与 2或
const (
	RelAnd = 1
	RelOr  = 2
)

const (
	SysPass    = 1 //系统匹配通过
	SysReject  = 2 //系统拒绝
	ObsAudit   = 3 //obs审核
	SysDefPass = 4 //系统未匹配通过
)

var AuditStatus = map[int]int{
	SysPass:    20, //规则引擎校验，自动通过
	SysReject:  21, //规则引擎校验，自动拒绝
	SysDefPass: 22, //规则全不匹配，自动通过
	ObsAudit:   30, //人工审核中
}

//规则匹配结果
type RuleMatch struct {
	RuleId      int
	FlowId      int
	RuleGo      int
	RMatch      bool
	Explain     string
	ItemMatches []ItemMatch
}

//规则项匹配结果
type ItemMatch struct {
	ItemId  int
	IMatch  bool
	Explain string
}

//bussData 转成对应项的string值
func bussDataToString(field string, bussData *rabbit.BusinessData) string {
	switch field {
	case "catId":
		return fmt.Sprintf("%d", bussData.CatId)
	case "chargePrice":
		return fmt.Sprintf("%0.4f", bussData.ChargePrice)
	case "rate":
		return fmt.Sprintf("%0.4f", bussData.Rate)
	case "virWhCode":
		return bussData.VirWhCode
	case "priceLoss":
		return "0.1"
	}

	return ""
}

//rule多条规则比较
//返回结果: 1 系统通过，2 系统驳回，3 转人工审核
func RunRuleMatch(bussData *rabbit.BusinessData, auditType *AuditType) (int, RuleMatch) {

	var rml []RuleMatch
	lenRule := len(auditType.RuleList)

	for i, rule := range auditType.RuleList {

		var iml []ItemMatch

		for _, item := range rule.ItemList {
			field := bussDataToString(item.Field, bussData)
			match := ValueCompare(field, item.Operate, item.Value)
			im := ItemMatch{
				ItemId:  item.ItemId,
				IMatch:  match,
				Explain: fmt.Sprintf(`(bussData.%v) [%v %v %v]`, item.Field, field, item.Operate, item.Value),
			}
			iml = append(iml, im)
		}

		//该条rule的验证结果
		result := RuleMatched
		switch rule.RuleRel {
		case RelAnd:
			for _, im := range iml {
				if im.IMatch {
					continue
				}
				result = RuleNotMatched
				break
			}
		case RelOr:
			for _, im := range iml {
				if im.IMatch {
					break
				}
				result = RuleNotMatched
			}
		}

		//基于规则引擎校验的结果进行进一步处理
		rml = append(rml, RuleMatch{
			RMatch:      result == RuleMatched,
			RuleId:      rule.RuleId,
			FlowId:      rule.FlowId,
			RuleGo:      rule.RuleProc,
			Explain:     fmt.Sprintf("rule items rel %d (1:and 2:or)", rule.RuleRel),
			ItemMatches: iml,
		})

		log.Printf("Rule[%d]: %+v\n", i, rml[len(rml)-1])

		if result == RuleMatched { //任一条rule通过，则进入下一步
			//1 系统通过，2 系统驳回，3 转人工审核，
			return AuditStatus[rule.RuleProc], rml[len(rml)-1]
		} else if i < lenRule-2 {
			continue
		}
	}

	//如果都不匹配，默认规则放行
	return AuditStatus[SysDefPass], RuleMatch{RuleId: 0, FlowId: 0}
}