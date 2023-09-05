package main

import (
	"log"

	"github.com/tidwall/gjson"
)

var (
	//条件列表
	ConditionList = "conditionList"
	//条件清单
	ConditionItems = "conditionItems"
	//提取主体：名单、指标、业务字段
	Name_list       string
	Indic_list      string
	Indic_list_nest string
	Fields_list     string
	//风险标签
	Label_list string
	//模板名称
	Tmpl_name []string
	Tmpl_id   string
	//方法名称
	Method_name []string
	Method_id   string
	//提取indictors_param：过滤条件、方法；模块->提取parameters
	FilterCondition = "filterCondition"
	Methods         = "methods"
	//参数：可能存在指标数据的主体
	Parameters = "parameters"
	Value      = "value"
)

/*解析规则中的rule_json，提取主体：字段、指标、名单*/
func parseRuleJson(rule_json string) error {
	data := gjson.Get(rule_json, ConditionList).Array()
	for _, ardt := range data {
		indicDt := gjson.Get(ardt.Raw, ConditionItems).Array()
		for _, itermDt := range indicDt {
			// log.Println(itermDt.Raw)
			fieldSource := gjson.Get(itermDt.Raw, "fieldSource").String()
			fieldValue := gjson.Get(itermDt.Raw, "fieldValue").String()
			//指标提取
			if fieldSource == "CALCULATION_VAR" {
				// log.Println(gjson.Get(itermDt.Raw, "fieldValue").String())
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
				// log.Println(fieldValue)
			}
			//业务字段提取
			if fieldSource == "BUSINESS_VAR" {
				if Fields_list == "" {
					Fields_list = "'" + fieldValue + "'"
				} else {
					Fields_list = appendValsSingle(Fields_list, fieldValue)
				}
				// log.Println(fieldValue)
			}
		}
	}

	return nil
}

/*解析指标中的indic_params，提取主体：嵌套指标（来源：方法methods、过滤条件filterCondition）、模板、模板参数*/
func parseIndicParam(indicators_param string) error {
	//过滤条件提取指标数据
	// indicators_param = "{\"clazz\":\"一段时间内主体操作次数统计\",\"id\":\"一段时间内主体操作次数统计\",\"filterCondition\":[{\"cnName\":\"条件类型\",\"enName\":\"JUDGE\",\"id\":\"ab0a5098006a469f\",\"parameters\":[{\"cnName\":\"过滤对象\",\"src\":\"inputOrSrcList\",\"dataType\":\"String\",\"enName\":\"filterObject\",\"index\":1,\"value\":{\"name\":\"设备当前交易时间和上一次交易时间差值(新)\",\"id\":\"ec2849e770854c9c81cac30d3fce3492\",\"value\":null,\"srcType\":\"CALCULATION_VAR\"}},{\"cnName\":\"运算符\",\"src\":\"下拉框\",\"dataType\":\"String\",\"enName\":\"operationalCharacter\",\"index\":2,\"value\":{\"name\":\"小于等于\",\"id\":\"6a5f9893ef784189920aa2ba18298eae\",\"value\":\"<=\",\"srcType\":\"COMPARATOR_VAR\"}},{\"cnName\":\"值\",\"src\":\"inputOrSrcList\",\"dataType\":\"String\",\"enName\":\"filterObjectValue\",\"index\":3,\"value\":{\"dataType\":\"String\",\"value\":\"30000\",\"id\":\"30000\",\"name\":\"30000\",\"srcType\":\"INPUT\"}}],\"filterSign\":\"A\"}],\"logicalRelationship\":\"and\",\"customizeConfig\":\"\",\"methods\":[{\"cnName\":\"一段时间内主体操作次数统计\",\"enName\":\"一段时间内主体操作次数统计\",\"parameters\":[{\"cnName\":\"时间间隔\",\"enName\":\"interval\",\"index\":0,\"comment\":\"\",\"value\":{\"id\":\"7\",\"name\":\"7\",\"value\":\"7\",\"dataType\":\"String\",\"srcType\":\"INPUT\",\"selectT\":\"7\",\"selectV\":\"7\"},\"src\":\"31\",\"dataType\":\"int\"},{\"cnName\":\"时间类型\",\"enName\":\"timeType\",\"index\":1,\"comment\":\"\",\"value\":{\"name\":\"天\",\"id\":\"40423bcb3ef54d68b47337317b4650a2\",\"value\":\"day\",\"srcType\":\"SYS_VAR\",\"selectT\":\"天\",\"selectV\":\"40423bcb3ef54d68b47337317b4650a2\"},\"src\":\"下拉框\",\"dataType\":\"String\"},{\"cnName\":\"主体\",\"enName\":\"mainPart\",\"index\":2,\"comment\":\"\",\"value\":{\"isEnumeration\":\"0\",\"dataType\":\"String\",\"name\":\"设备指纹\",\"id\":\"db69056dd9fa441e926dafb7e3426549\",\"value\":\"dev_fp\",\"srcType\":\"BUSINESS_VAR\",\"selectT\":\"设备指纹\",\"selectV\":\"db69056dd9fa441e926dafb7e3426549\"},\"src\":\"下拉框\",\"dataType\":\"String\"},{\"cnName\":\"操作状态\",\"enName\":\"opState\",\"index\":3,\"comment\":\"\",\"value\":{\"name\":\"成功\",\"id\":\"3ba4fc6b89a2411d867a47721186d451\",\"value\":\"1\",\"srcType\":\"SYS_VAR\",\"selectT\":\"成功\",\"selectV\":\"3ba4fc6b89a2411d867a47721186d451\"},\"src\":\"下拉框\",\"dataType\":\"int\"}],\"channelTrade\":[{\"cnName\":\"适用渠道\",\"enName\":\"channelType\",\"index\":0,\"type\":null,\"value\":{\"name\":\"当前渠道\",\"id\":\"1a95072964984fb280cd7812cb386a41\",\"value\":\"CURRENT\",\"srcType\":\"SYS_VAR\"},\"dataType\":\"String\"},{\"cnName\":\"应用标识\",\"enName\":\"appid\",\"index\":1,\"type\":null,\"value\":{\"name\":\"\",\"id\":\"\",\"value\":\"\",\"srcType\":\"APP_VAR\"},\"dataType\":\"String\"},{\"cnName\":\"适用交易\",\"enName\":\"tradeType\",\"index\":2,\"type\":null,\"value\":{\"name\":\"当前交易\",\"id\":\"7a9f41b09c854a3385db999632b7e8e9\",\"value\":\"CURRENT\",\"srcType\":\"TRADETYPE_VAR\"},\"dataType\":\"String\"},{\"cnName\":\"交易类型\",\"enName\":\"opType\",\"index\":3,\"type\":null,\"value\":{\"name\":\"\",\"id\":\"\",\"value\":\"\",\"srcType\":\"TRADETYPE_VAR\"},\"dataType\":\"String\"}]}]}"
	fcdata := gjson.Get(indicators_param, FilterCondition).Array()
	//方法中提取指标数据
	mtdata := gjson.Get(indicators_param, Methods).Array()
	/*解析过滤条件模块*/
	if len(fcdata) == 0 {
		log.Printf("FilterCondition中无嵌套指标！ \n")
	} else {
		for _, fcdt := range fcdata {
			//过滤条件中参数获取:parameters
			fcParams := gjson.Get(fcdt.Raw, Parameters).Array()
			// log.Printf("fcParams=%v \n", fcParams)
			for _, fcParam := range fcParams {
				// log.Printf("fcParam=%v \n", fcParam.Raw)
				fcIterm := gjson.Get(fcParam.Raw, Value)
				srcTyp := gjson.Get(fcIterm.Raw, "srcType").String()
				srcId := gjson.Get(fcIterm.Raw, "id").String()
				//业务字段名称
				srcValue := gjson.Get(fcIterm.Raw, "value").String()
				// log.Printf("srcTyp=%v \n", srcTyp)
				// log.Printf("srcId=%v \n", srcId)
				/*提取指标类型数据*/
				if srcTyp == "CALCULATION_VAR" {
					// log.Printf("filterCondition: srcTyp=%v,srcId=%v \n", srcTyp, srcId)
					if Indic_list_nest == "" {
						Indic_list_nest = "'" + srcId + "'"
					} else {
						Indic_list_nest = appendValsSingle(Indic_list_nest, srcId)
						log.Printf("查询到过滤条件中嵌套指标：%v \n", srcId)
					}

					// if err = QueryIndicators(srcId); err != nil {
					// 	log.Printf("查询导出过滤条件中的嵌套指标失败：%v \n", err)
					// } else {
					// 	log.Printf("导出过滤条件中嵌套指标成功！：%v \n", srcId)
					// }
				}

				//名单提取
				if srcTyp == "LIST_LIB" {
					if Name_list == "" {
						Name_list = "'" + srcId + "'"
					} else {
						Name_list = appendValsSingle(Name_list, srcId)
						log.Printf("查询到过滤条件中嵌套名单：%v \n", srcId)
					}
					// log.Println(fieldValue)
				}
				//业务字段提取
				if srcTyp == "BUSINESS_VAR" {
					if Fields_list == "" {
						Fields_list = "'" + srcValue + "'"
					} else {
						Fields_list = appendValsSingle(Fields_list, srcValue)
						log.Printf("查询到过滤条件中嵌套业务字段ID=%v,名称=%v： \n", srcId, srcValue)
					}
					// log.Println(fieldValue)
				}

			}
		}
	}

	/*解析方法模块:methods*/
	if len(mtdata) == 0 {
		log.Printf("Method中无嵌套指标！ \n")
	} else {
		// log.Printf("Method:mtdata=%v", mtdata)
		for _, mthdata := range mtdata {
			//方法中参数获取:parameters
			mtParams := gjson.Get(mthdata.Raw, Parameters).Array()
			// log.Printf("Method:mtParams=%v", mtParams)
			for _, mtParam := range mtParams {
				mtIterm := gjson.Get(mtParam.Raw, Value)
				srcTyp := gjson.Get(mtIterm.Raw, "srcType").String()
				srcId := gjson.Get(mtIterm.Raw, "id").String()
				//业务字段名称
				srcValue := gjson.Get(mtIterm.Raw, "value").String()
				/*提取指标类型数据*/
				if srcTyp == "CALCULATION_VAR" {
					// log.Printf("Method:srcTyp=%v,srcId=%v \n", srcTyp, srcId)
					if Indic_list_nest == "" {
						Indic_list_nest = "'" + srcId + "'"
					} else {
						Indic_list_nest = appendValsSingle(Indic_list_nest, srcId)
						log.Printf("查询到方法中嵌套指标：%v \n", srcId)
					}

					// if err = QueryIndicators(srcId); err != nil {
					// 	log.Printf("查询导出方法中的嵌套指标失败：%v \n", err)
					// } else {
					// 	log.Printf("导出方法中嵌套指标成功！：%v \n", srcId)
					// }
				}
				//名单提取
				if srcTyp == "LIST_LIB" {
					if Name_list == "" {
						Name_list = "'" + srcId + "'"
					} else {
						Name_list = appendValsSingle(Name_list, srcId)
						log.Printf("查询到方法中嵌套名单：%v \n", srcId)
					}
					// log.Println(fieldValue)
				}
				//业务字段提取
				if srcTyp == "BUSINESS_VAR" {
					if Fields_list == "" {
						Fields_list = "'" + srcValue + "'"
					} else {
						Fields_list = appendValsSingle(Fields_list, srcValue)
						log.Printf("查询到方法中嵌套业务字段ID=%v,名称=%v \n", srcId, srcValue)
					}
					// log.Println(srcValue)
				}
			}
		}
	}

	//返回值：无报错
	return nil
}
