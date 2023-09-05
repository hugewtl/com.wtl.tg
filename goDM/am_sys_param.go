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

type OutAmLabel struct {
	ID          string
	GROUP_ID    string
	PARAM_NAME  string
	PARAM_VALUE string
	TYPE        string
	TYPE_VALUE  string
	IS_INTERNAL string
	REMARK      string
	CREATE_USER string
	CREATE_TIME time.Time
	UPDATE_USER string
	UPDATE_TIME time.Time
}

type AmLabel struct {
	ID          sql.NullString
	GROUP_ID    sql.NullString
	PARAM_NAME  sql.NullString
	PARAM_VALUE sql.NullString
	TYPE        sql.NullString
	TYPE_VALUE  sql.NullString
	IS_INTERNAL sql.NullString
	REMARK      sql.NullString
	CREATE_USER sql.NullString
	CREATE_TIME sql.NullTime
	UPDATE_USER sql.NullString
	UPDATE_TIME sql.NullTime
}

func CopyAmLabel(outAl OutAmLabel, Al AmLabel) OutAmLabel {
	outAl.ID = Al.ID.String
	outAl.GROUP_ID = Al.GROUP_ID.String
	outAl.PARAM_NAME = Al.PARAM_NAME.String
	outAl.PARAM_VALUE = Al.PARAM_NAME.String
	outAl.TYPE = Al.TYPE.String
	outAl.TYPE_VALUE = Al.TYPE_VALUE.String
	outAl.IS_INTERNAL = Al.IS_INTERNAL.String
	outAl.REMARK = Al.REMARK.String
	outAl.CREATE_USER = Al.CREATE_USER.String
	outAl.CREATE_TIME = Al.CREATE_TIME.Time
	outAl.UPDATE_USER = Al.UPDATE_USER.String
	outAl.UPDATE_TIME = Al.UPDATE_TIME.Time
	return outAl
}

func QueryAmLabel(ids string) error {
	al := AmLabel{}
	outAl := OutAmLabel{}
	// 字段数据拼接：fields、fieldsTab
	var fields string
	var fieldsTab string
	for i := 0; i < reflect.TypeOf(al).NumField(); i++ {
		s := reflect.TypeOf(al).Field(i).Name
		if i == 0 {
			fields = fmt.Sprintf(s)
			fieldsTab = fmt.Sprintf("(" + s)
		} else if i <= reflect.TypeOf(al).NumField()-1 {
			fields = fmt.Sprintf(fields + "," + s)
			if i == reflect.TypeOf(al).NumField()-1 {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s + ")")
			} else {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s)
			}
		}
	}

	//解析字段数据
	colsNum := reflect.TypeOf(outAl).NumField()
	//查询所有关联规则的指标数据

	sql := "SELECT " + fields + " FROM AM_SYS_PARAMS WHERE ID IN (" + ids + ")"
	sqfName := "am_sys_params_" + ParamsMp["set_name"] + "_" + time.Time.Format(time.Now(), "20060102150405") + ".sql"
	sqlfile, err := os.Create(sqfName)
	rows, err := db.Query(sql, ids)
	if err != nil {
		log.Printf("query sysParam Table err:%v \n", err)
	}
	defer rows.Close()
	// log.Println("queryTable results:")
	for rows.Next() {
		outAl := OutAmLabel{}
		record := "NULL"
		if err = rows.Scan(&al.ID, &al.GROUP_ID, &al.PARAM_NAME, &al.PARAM_VALUE, &al.TYPE, &al.TYPE_VALUE, &al.IS_INTERNAL, &al.REMARK, &al.CREATE_USER, &al.CREATE_TIME, &al.UPDATE_USER, &al.UPDATE_TIME); err != nil {
			return err
		}
		res := CopyAmLabel(outAl, al)

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
						record = appendVals(record, clobToStr)
					} else if err == io.EOF {
						log.Printf("标签读取Clob结束：%v\n", err)
					} else {
						log.Printf("标签数据clob非空解析报错：%v\n", err)
					}

					// log.Println(record)
				}
			}
		}
		//拼接inster SQL语句，打印
		isrtSql := "insert into am_sys_params " + fieldsTab + " values (" + record + ");"
		sqlfile.WriteString(isrtSql + "\n")
	}
	//关闭文件
	defer sqlfile.Close()
	return nil
}
