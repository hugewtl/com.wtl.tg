package main

import (
	"database/sql"
	"dm"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
)

type AmIndicator struct {
	ID               sql.NullString
	GROUP_ID         sql.NullString
	INDIC_NAME       sql.NullString
	INDIC_TYPE       sql.NullString
	INDIC_MODE       sql.NullString
	INDIC_NO         sql.NullString
	TMPL_GROUP_ID    sql.NullString
	TMPL_ID          sql.NullString
	EXPRESSION       sql.NullString
	INDICATORS_PARAM *dm.DmClob
	STATES           sql.NullString
	IS_INTERNAL      sql.NullString
	REMARK           sql.NullString
	CREATE_USER      sql.NullString
	CREATE_TIME      sql.NullTime
	UPDATE_USER      sql.NullString
	UPDATE_TIME      sql.NullTime
}

type OutAmIndicator struct {
	ID               string
	GROUP_ID         string
	INDIC_NAME       string
	INDIC_TYPE       string
	INDIC_MODE       string
	INDIC_NO         string
	TMPL_GROUP_ID    string
	TMPL_ID          string
	EXPRESSION       string
	INDICATORS_PARAM *dm.DmClob
	STATES           string
	IS_INTERNAL      string
	REMARK           string
	CREATE_USER      string
	CREATE_TIME      time.Time
	UPDATE_USER      string
	UPDATE_TIME      time.Time
}

func CopyAmIndicator(outAi OutAmIndicator, Ai AmIndicator) OutAmIndicator {
	outAi.ID = Ai.ID.String
	outAi.GROUP_ID = Ai.GROUP_ID.String
	outAi.INDIC_NAME = Ai.INDIC_NAME.String
	outAi.INDIC_TYPE = Ai.INDIC_TYPE.String
	outAi.INDIC_MODE = Ai.INDIC_MODE.String
	outAi.INDIC_NO = Ai.INDIC_NO.String
	outAi.TMPL_GROUP_ID = Ai.TMPL_GROUP_ID.String
	outAi.TMPL_ID = Ai.TMPL_ID.String
	outAi.EXPRESSION = Ai.EXPRESSION.String
	outAi.INDICATORS_PARAM = Ai.INDICATORS_PARAM
	outAi.STATES = Ai.STATES.String
	outAi.IS_INTERNAL = Ai.IS_INTERNAL.String
	outAi.REMARK = Ai.REMARK.String
	outAi.CREATE_USER = Ai.CREATE_USER.String
	outAi.CREATE_TIME = Ai.CREATE_TIME.Time
	outAi.UPDATE_USER = Ai.UPDATE_USER.String
	outAi.UPDATE_TIME = Ai.UPDATE_TIME.Time
	return outAi
}

func QueryIndicators(ids string) error {
	ai := AmIndicator{}
	outAi := OutAmIndicator{}
	// 字段数据拼接：fields、fieldsTab
	var fields string
	var fieldsTab string
	for i := 0; i < reflect.TypeOf(ai).NumField(); i++ {
		s := reflect.TypeOf(ai).Field(i).Name
		if i == 0 {
			fields = fmt.Sprintf(s)
			fieldsTab = fmt.Sprintf("(" + s)
		} else if i <= reflect.TypeOf(ai).NumField()-1 {
			fields = fmt.Sprintf(fields + "," + s)
			if i == reflect.TypeOf(ai).NumField()-1 {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s + ")")
			} else {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s)
			}
		}
	}

	//解析字段数据
	colsNum := reflect.TypeOf(outAi).NumField()
	//查询所有关联规则的指标数据
	sql := "SELECT " + fields + " FROM AM_INDICATORS WHERE ID IN (" + ids + ")"
	sqfName := "am_indicators_" + ParamsMp["set_name"] + "_" + time.Time.Format(time.Now(), "20060102150405") + ".sql"
	sqlfile, err := os.Create(sqfName)
	rows, err := db.Query(sql, ids)
	if err != nil {
		log.Printf("query Indicators Table err:%v \n", err)
	}
	defer rows.Close()
	// log.Println("queryTable results:")
	for rows.Next() {
		outAi := OutAmIndicator{}
		record := "NULL"
		if err = rows.Scan(&ai.ID, &ai.GROUP_ID, &ai.INDIC_NAME, &ai.INDIC_TYPE, &ai.INDIC_MODE, &ai.INDIC_NO, &ai.TMPL_GROUP_ID, &ai.TMPL_ID, &ai.EXPRESSION, &ai.INDICATORS_PARAM, &ai.STATES, &ai.IS_INTERNAL, &ai.REMARK, &ai.CREATE_USER, &ai.CREATE_TIME, &ai.UPDATE_USER, &ai.UPDATE_TIME); err != nil {
			return err
		}
		res := CopyAmIndicator(outAi, ai)

		Typ := reflect.TypeOf(res)
		//字段值，&res为了获取对应字段值
		Val := reflect.ValueOf(&res).Elem()

		//遍历输出结构
		for i := 0; i < colsNum; i++ {
			// log.Printf("TYPE:%v \n", Typ.Field(i).Type.String())
			val := Val.Field(i)
			//反射获取字段类型，按类型对数据进行归类处理
			if Typ.Field(i).Type.String() == "string" {
				if val.String() == "" {
					// log.Printf("STRING:%v \n", "NULL")
					// record = record + "'NULL'"
					record = appendVals(record, "NULL")
					// log.Println(record)
				} else {
					// log.Printf("STRING:%v \n", val.String())
					// record = record + val.String()
					record = appendVals(record, val.String())
					// log.Println(record)
				}
			}
			if Typ.Field(i).Type.String() == "time.Time" {
				if val.Interface().(time.Time).Format("2006-01-02 15:04:05.000") == "0001-01-01 00:00:00.000" {
					// log.Printf("TIME:%v \n", "NULL")
					// record = record + "'NULL'"
					record = appendVals(record, "NULL")
					// log.Println(record)
				} else {
					// log.Printf("TIME:%v \n", val.Interface().(time.Time).Format("2006-01-02 15:04:05.000"))
					// record = record + val.Interface().(time.Time).Format("2006-01-02 15:04:05.000")
					record = appendVals(record, val.Interface().(time.Time).Format("2006-01-02 15:04:05.000"))
					// log.Println(record)
				}
			}
			if Typ.Field(i).Type.String() == "*dm.DmClob" {
				if val.Interface().(*dm.DmClob) == nil {
					record = appendVals(record, "NULL")
				} else {
					//读取clob为字符串
					clob := *val.Interface().(*dm.DmClob)
					clobInt64, err := clob.GetLength()
					if err != nil {
						log.Println(err)
					}
					//将int64转换int32
					strInt64 := strconv.FormatInt(clobInt64, 10)
					Int32, err := strconv.Atoi(strInt64)
					// log.Printf("Int32:%v \n", Int32)
					clobToStr, err := clob.ReadString(1, Int32)

					if clobInt64 == 0 {
						record = appendVals(record, "NULL")
					}
					if err == nil {
						//遍历indicator_params，解析主体：导出嵌套指标---字段、名单集
						if reflect.TypeOf(res).Field(i).Name == "INDICATORS_PARAM" {
							err = parseIndicParam(clobToStr)
							if err != nil {
								log.Printf("解析INDICATORS_PARAM出错：err %v \n", err)
							}
						}
						record = appendVals(record, clobToStr)
					} else if err == io.EOF {
						log.Printf("指标读取Clob结束：%v\n", err)
					} else {
						log.Printf("指标数据clob非空解析报错：%v\n", err)
					}
					// log.Println(record)
				}
			}
		}
		//拼接inster SQL语句，打印
		isrtSql := "insert into am_indicators " + fieldsTab + " values (" + record + ");"
		sqlfile.WriteString(isrtSql + "\n")
	}
	//关闭文件
	defer sqlfile.Close()
	return nil
}

func QueryNestIndicators(ids string) error {
	ai := AmIndicator{}
	outAi := OutAmIndicator{}
	// 字段数据拼接：fields、fieldsTab
	var fields string
	var fieldsTab string
	for i := 0; i < reflect.TypeOf(ai).NumField(); i++ {
		s := reflect.TypeOf(ai).Field(i).Name
		if i == 0 {
			fields = fmt.Sprintf(s)
			fieldsTab = fmt.Sprintf("(" + s)
		} else if i <= reflect.TypeOf(ai).NumField()-1 {
			fields = fmt.Sprintf(fields + "," + s)
			if i == reflect.TypeOf(ai).NumField()-1 {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s + ")")
			} else {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s)
			}
		}
	}

	//解析字段数据
	colsNum := reflect.TypeOf(outAi).NumField()
	//查询所有关联规则的指标数据
	sql := "SELECT " + fields + " FROM AM_INDICATORS WHERE ID IN (" + ids + ")"
	sqfName := "am_indicators_" + ParamsMp["set_name"] + "_Nest.sql"
	sqlfile, err := os.Create(sqfName)
	rows, err := db.Query(sql, ids)
	if err != nil {
		log.Printf("query Indicators Table err:%v \n", err)
	}
	defer rows.Close()
	// log.Println("queryTable results:")
	for rows.Next() {
		outAi := OutAmIndicator{}
		record := "NULL"
		if err = rows.Scan(&ai.ID, &ai.GROUP_ID, &ai.INDIC_NAME, &ai.INDIC_TYPE, &ai.INDIC_MODE, &ai.INDIC_NO, &ai.TMPL_GROUP_ID, &ai.TMPL_ID, &ai.EXPRESSION, &ai.INDICATORS_PARAM, &ai.STATES, &ai.IS_INTERNAL, &ai.REMARK, &ai.CREATE_USER, &ai.CREATE_TIME, &ai.UPDATE_USER, &ai.UPDATE_TIME); err != nil {
			return err
		}
		res := CopyAmIndicator(outAi, ai)

		Typ := reflect.TypeOf(res)
		//字段值，&res为了获取对应字段值
		Val := reflect.ValueOf(&res).Elem()

		//遍历输出结构
		for i := 0; i < colsNum; i++ {
			// log.Printf("TYPE:%v \n", Typ.Field(i).Type.String())
			val := Val.Field(i)
			//反射获取字段类型，按类型对数据进行归类处理
			if Typ.Field(i).Type.String() == "string" {
				if val.String() == "" {
					// log.Printf("STRING:%v \n", "NULL")
					// record = record + "'NULL'"
					record = appendVals(record, "NULL")
					// log.Println(record)
				} else {
					// log.Printf("STRING:%v \n", val.String())
					// record = record + val.String()
					record = appendVals(record, val.String())
					// log.Println(record)
				}
			}
			if Typ.Field(i).Type.String() == "time.Time" {
				if val.Interface().(time.Time).Format("2006-01-02 15:04:05.000") == "0001-01-01 00:00:00.000" {
					// log.Printf("TIME:%v \n", "NULL")
					// record = record + "'NULL'"
					record = appendVals(record, "NULL")
					// log.Println(record)
				} else {
					// log.Printf("TIME:%v \n", val.Interface().(time.Time).Format("2006-01-02 15:04:05.000"))
					// record = record + val.Interface().(time.Time).Format("2006-01-02 15:04:05.000")
					record = appendVals(record, val.Interface().(time.Time).Format("2006-01-02 15:04:05.000"))
					// log.Println(record)
				}
			}
			if Typ.Field(i).Type.String() == "*dm.DmClob" {
				if val.Interface().(*dm.DmClob) == nil {
					record = appendVals(record, "NULL")
				} else {
					//读取clob为字符串
					clob := *val.Interface().(*dm.DmClob)
					clobInt64, err := clob.GetLength()
					if err != nil {
						log.Println(err)
					}
					//将int64转换int32
					strInt64 := strconv.FormatInt(clobInt64, 10)
					Int32, err := strconv.Atoi(strInt64)
					// log.Printf("Int32:%v \n", Int32)
					clobToStr, err := clob.ReadString(1, Int32)

					if clobInt64 == 0 {
						record = appendVals(record, "NULL")
					}
					if err == nil {
						//遍历indicator_params，解析主体：导出嵌套指标---字段、名单集
						if reflect.TypeOf(res).Field(i).Name == "INDICATORS_PARAM" {
							err = parseIndicParam(clobToStr)
							if err != nil {
								log.Printf("解析INDICATORS_PARAM出错：err %v \n", err)
							}
						}
						record = appendVals(record, clobToStr)
					} else if err == io.EOF {
						log.Printf("指标读取Clob结束：%v\n", err)
					} else {
						log.Printf("指标数据clob非空解析报错：%v\n", err)
					}
					// log.Println(record)
				}
			}
		}
		//拼接inster SQL语句，打印
		isrtSql := "insert into am_indicators " + fieldsTab + " values (" + record + ");"
		sqlfile.WriteString(isrtSql + "\n")
	}
	//关闭文件
	defer sqlfile.Close()
	return nil
}
