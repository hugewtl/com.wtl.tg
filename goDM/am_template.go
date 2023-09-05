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

type OutAmTemplate struct {
	ID            string
	TEMP_NAME     string
	TEMP_PATH     string
	STATES        string
	REMARK        string
	CREATE_USER   string
	CREATE_TIME   time.Time
	UPDATE_USER   string
	UPDATE_TIME   time.Time
	IS_INTERNAL   string
	TMPL_GROUP_ID string
	DATA_TYPE     string
}

type AmTemplate struct {
	ID            sql.NullString
	TEMP_NAME     sql.NullString
	TEMP_PATH     sql.NullString
	STATES        sql.NullString
	REMARK        sql.NullString
	CREATE_USER   sql.NullString
	CREATE_TIME   sql.NullTime
	UPDATE_USER   sql.NullString
	UPDATE_TIME   sql.NullTime
	IS_INTERNAL   sql.NullString
	TMPL_GROUP_ID sql.NullString
	DATA_TYPE     sql.NullString
}

func CopyAmTemplate(outAtl OutAmTemplate, Atl AmTemplate) OutAmTemplate {
	outAtl.ID = Atl.ID.String
	outAtl.TEMP_NAME = Atl.TEMP_NAME.String
	outAtl.TEMP_PATH = Atl.TEMP_PATH.String
	outAtl.STATES = Atl.STATES.String
	outAtl.REMARK = Atl.REMARK.String
	outAtl.CREATE_USER = Atl.CREATE_USER.String
	outAtl.CREATE_TIME = Atl.CREATE_TIME.Time
	outAtl.UPDATE_USER = Atl.UPDATE_USER.String
	outAtl.UPDATE_TIME = Atl.UPDATE_TIME.Time
	outAtl.IS_INTERNAL = Atl.IS_INTERNAL.String
	outAtl.TMPL_GROUP_ID = Atl.TMPL_GROUP_ID.String
	outAtl.DATA_TYPE = Atl.DATA_TYPE.String
	return outAtl
}

func queryAmTemplate(tName string) error {
	Atl := AmTemplate{}
	outAtl := OutAmTemplate{}
	// 字段数据拼接：fields、fieldsTab
	var fields string
	var fieldsTab string
	for i := 0; i < reflect.TypeOf(Atl).NumField(); i++ {
		s := reflect.TypeOf(Atl).Field(i).Name
		if i == 0 {
			fields = fmt.Sprintf(s)
			fieldsTab = fmt.Sprintf("(" + s)
		} else if i <= reflect.TypeOf(Atl).NumField()-1 {
			fields = fmt.Sprintf(fields + "," + s)
			if i == reflect.TypeOf(Atl).NumField()-1 {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s + ")")
			} else {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s)
			}
		}
	}

	//解析字段数据
	colsNum := reflect.TypeOf(outAtl).NumField()
	//查询所有关联规则的指标数据
	sql := "SELECT " + fields + " FROM AM_INDIC_TMPL WHERE temp_name IN (" + tName + ")"
	sqfName := "am_indic_tmpl.sql"
	sqlfile, err := os.Create(sqfName)

	rows, err := db.Query(sql, tName)
	if err != nil {
		log.Printf("query AM_INDIC_TMPL Table err:%v \n", err)
	}
	defer rows.Close()
	// log.Println("queryTable results:")
	for rows.Next() {
		outAtl := OutAmTemplate{}
		record := "NULL"
		if err = rows.Scan(&Atl.ID, &Atl.TEMP_NAME, &Atl.TEMP_PATH, &Atl.STATES, &Atl.REMARK, &Atl.CREATE_USER, &Atl.CREATE_TIME, &Atl.UPDATE_USER, &Atl.UPDATE_TIME, &Atl.IS_INTERNAL, &Atl.TMPL_GROUP_ID, &Atl.DATA_TYPE); err != nil {
			return err
		}
		res := CopyAmTemplate(outAtl, Atl)

		Typ := reflect.TypeOf(res)
		//字段值，&res为了获取对应字段值
		Val := reflect.ValueOf(&res).Elem()
		//提取模板ID
		if len(Tmpl_id) == 0 {
			Tmpl_id = "'" + Atl.ID.String + "'"
		} else {
			Tmpl_id = appendVals(Tmpl_id, Atl.ID.String)
		}

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
					// log.Printf("长度：%v\n", clobInt64)
					if err != nil {
						fmt.Printf("读取Clob长度报错：%v\n", err)
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
						fmt.Printf("模板读取Clob结束：%v\n", err)
					} else {
						fmt.Printf("模板数据clob非空解析报错：%v\n", err)
					}
					// log.Println(record)
				}
			}
		}
		//拼接inster SQL语句，打印
		isrtSql := "insert into AM_INDIC_TMPL " + fieldsTab + " values (" + record + ");"
		sqlfile.WriteString(isrtSql + "\n")
	}
	//关闭文件
	defer sqlfile.Close()
	return nil
}
