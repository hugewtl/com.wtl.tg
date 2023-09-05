package main

import (
	"github.com/tidwall/gjson"
)

var (
	//条件列表
	ConditionList = "conditionList"
	//条件清单
	ConditionItems = "conditionItems"
	//提取主体：名单、指标、业务字段
	Name_list   string
	Indic_list  string
	Fields_list string
	//风险标签
	Label_list string
	//模板名称
	Tmpl_name []string
	Tmpl_id   string
	//方法名称
	Method_name []string
	Method_id   string
)

/*解析规则中的rule_json，提取主体：字段、指标、名单*/
func parseRuleJson(rule_json string) error {
	data := gjson.Get(rule_json, ConditionList).Array()
	for _, ardt := range data {
		indicDt := gjson.Get(ardt.Raw, ConditionItems).Array()
		for _, itermDt := range indicDt {
			// fmt.Println(itermDt.Raw)
			fieldSource := gjson.Get(itermDt.Raw, "fieldSource").String()
			fieldValue := gjson.Get(itermDt.Raw, "fieldValue").String()
			//指标提取
			if fieldSource == "CALCULATION_VAR" {
				// fmt.Println(gjson.Get(itermDt.Raw, "fieldValue").String())
				if Indic_list == "" {
					Indic_list = "'" + fieldValue + "'"
				} else {
					Indic_list = appendValsSingle(Indic_list, fieldValue)
				}

			}
			//名单提取
			if fieldSource == "LIST_LIB" {
				if Name_list == "" {
					Name_list = "'" + fieldValue + "'"
				} else {
					Name_list = appendValsSingle(Name_list, fieldValue)
				}
				// fmt.Println(fieldValue)
			}
			//业务字段提取
			if fieldSource == "BUSINESS_VAR" {
				if Fields_list == "" {
					Fields_list = "'" + fieldValue + "'"
				} else {
					Fields_list = appendValsSingle(Fields_list, fieldValue)
				}
				// fmt.Println(fieldValue)
			}
		}
	}

	return nil
}

/*解析指标中的indic_params，提取主体：模板、模板参数*/
