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

type OutAmTemplateParam struct {
	ID                 string
	TEMPLATE_ID        string
	PARAM_NAME         string
	PARAM_VALUE        string
	PARAM_TYPE         string
	PARAM_CHECKED_NAME string
	PARAM_DEFAULT      string
	STATES             string
	IS_INTERNAL        string
	REMARK             string
	SORT_NO            int
	CREATE_USER        string
	CREATE_TIME        time.Time
	UPDATE_USER        string
	UPDATE_TIME        time.Time
	DATA_TYPE          string
	TMPL_PARAM_ID      string
}
type AmTemplateParam struct {
	ID                 sql.NullString
	TEMPLATE_ID        sql.NullString
	PARAM_NAME         sql.NullString
	PARAM_VALUE        sql.NullString
	PARAM_TYPE         sql.NullString
	PARAM_CHECKED_NAME sql.NullString
	PARAM_DEFAULT      sql.NullString
	STATES             sql.NullString
	IS_INTERNAL        sql.NullString
	REMARK             sql.NullString
	SORT_NO            int
	CREATE_USER        sql.NullString
	CREATE_TIME        sql.NullTime
	UPDATE_USER        sql.NullString
	UPDATE_TIME        sql.NullTime
	DATA_TYPE          sql.NullString
	TMPL_PARAM_ID      sql.NullString
}

func CopyAmTemplateParam(OutAtlp OutAmTemplateParam, Atlp AmTemplateParam) OutAmTemplateParam {
	OutAtlp.ID = Atlp.ID.String
	OutAtlp.TEMPLATE_ID = Atlp.TEMPLATE_ID.String
	OutAtlp.PARAM_NAME = Atlp.PARAM_NAME.String
	OutAtlp.PARAM_VALUE = Atlp.PARAM_VALUE.String
	OutAtlp.PARAM_TYPE = Atlp.PARAM_TYPE.String
	OutAtlp.PARAM_CHECKED_NAME = Atlp.PARAM_CHECKED_NAME.String
	OutAtlp.PARAM_DEFAULT = Atlp.PARAM_DEFAULT.String
	OutAtlp.STATES = Atlp.STATES.String
	OutAtlp.IS_INTERNAL = Atlp.IS_INTERNAL.String
	OutAtlp.REMARK = Atlp.REMARK.String
	OutAtlp.SORT_NO = Atlp.SORT_NO
	OutAtlp.CREATE_USER = Atlp.CREATE_USER.String
	OutAtlp.CREATE_TIME = Atlp.CREATE_TIME.Time
	OutAtlp.UPDATE_USER = Atlp.UPDATE_USER.String
	OutAtlp.UPDATE_TIME = Atlp.UPDATE_TIME.Time
	OutAtlp.DATA_TYPE = Atlp.DATA_TYPE.String
	OutAtlp.TMPL_PARAM_ID = Atlp.TMPL_PARAM_ID.String
	return OutAtlp
}

func queryAmTemplateParam(ids string) error {
	Atlp := AmTemplateParam{}
	outAtlp := OutAmTemplateParam{}
	// 字段数据拼接：fields、fieldsTab
	var fields string
	var fieldsTab string
	for i := 0; i < reflect.TypeOf(Atlp).NumField(); i++ {
		s := reflect.TypeOf(Atlp).Field(i).Name
		if i == 0 {
			fields = fmt.Sprintf(s)
			fieldsTab = fmt.Sprintf("(" + s)
		} else if i <= reflect.TypeOf(Atlp).NumField()-1 {
			fields = fmt.Sprintf(fields + "," + s)
			if i == reflect.TypeOf(Atlp).NumField()-1 {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s + ")")
			} else {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s)
			}
		}
	}

	//解析字段数据
	colsNum := reflect.TypeOf(outAtlp).NumField()
	//查询所有关联规则的指标数据
	sql := "SELECT " + fields + " FROM AM_INDIC_TMPL_PARAM WHERE template_id IN (" + ids + ")"
	sqfName := "am_indic_tmpl_param.sql"
	sqlfile, err := os.Create(sqfName)

	rows, err := db.Query(sql, ids)
	if err != nil {
		log.Printf("query AM_INDIC_TMPL_PARAM Table err:%v \n", err)
	}
	defer rows.Close()
	// log.Println("queryTable results:")
	for rows.Next() {
		outAtlp := OutAmTemplateParam{}
		record := "NULL"
		if err = rows.Scan(&Atlp.ID, &Atlp.TEMPLATE_ID, &Atlp.PARAM_NAME, &Atlp.PARAM_VALUE, &Atlp.PARAM_TYPE, &Atlp.PARAM_CHECKED_NAME, &Atlp.PARAM_DEFAULT, &Atlp.STATES, &Atlp.IS_INTERNAL, &Atlp.REMARK, &Atlp.SORT_NO, &Atlp.CREATE_USER, &Atlp.CREATE_TIME, &Atlp.UPDATE_USER, &Atlp.UPDATE_TIME, &Atlp.DATA_TYPE, &Atlp.TMPL_PARAM_ID); err != nil {
			return err
		}
		res := CopyAmTemplateParam(outAtlp, Atlp)

		Typ := reflect.TypeOf(res)
		//字段值，&res为了获取对应字段值
		Val := reflect.ValueOf(&res).Elem()

		//遍历输出结构
		for i := 0; i < colsNum; i++ {
			// log.Printf("TYPE:%v \n", Typ.Field(i).Type.String())
			val := Val.Field(i)
			//反射获取字段类型，按类型对数据进行归类处理
			if Typ.Field(i).Type.String() == "int" {
				//int类型无引号包裹
				record = record + "," + strconv.Itoa(val.Interface().(int))
				// log.Println(record)
			}

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
					// log.Printf("长度：%v\n", clobInt64)
					if err != nil {
						log.Printf("读取Clob长度报错：%v\n", err)
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
						record = appendVals(record, clobToStr)
					} else if err == io.EOF {
						log.Printf("名单读取Clob结束：%v\n", err)
					} else {
						log.Printf("名单数据clob非空解析报错：%v\n", err)
					}
					// log.Println(record)
				}
			}
		}
		//拼接inster SQL语句，打印
		isrtSql := "insert into AM_INDIC_TMPL_PARAM " + fieldsTab + " values (" + record + ");"
		sqlfile.WriteString(isrtSql + "\n")
	}
	//关闭文件
	defer sqlfile.Close()
	return nil
}
