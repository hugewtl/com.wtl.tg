package main

import (
	"database/sql"
	"dm"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"time"
)

// 定义名单集表
type AmListName struct {
	ID           sql.NullString
	LIST_NAME    sql.NullString
	BUSINESS_VAR sql.NullString
	APP_ID       sql.NullString
	LIST_TYPE    sql.NullString
	STATES       sql.NullString
	DESCRIPTION  sql.NullString
	CREATE_USER  sql.NullString
	CREATE_TIME  sql.NullTime
	UPDATE_USER  sql.NullString
	UPDATE_TIME  sql.NullTime
	TRADE_TYPE   *dm.DmClob
	APPLY_SCOPE  sql.NullString
}

type OutAmListName struct {
	ID           string
	LIST_NAME    string
	BUSINESS_VAR string
	APP_ID       string
	LIST_TYPE    string
	STATES       string
	DESCRIPTION  string
	CREATE_USER  string
	CREATE_TIME  time.Time
	UPDATE_USER  string
	UPDATE_TIME  time.Time
	TRADE_TYPE   *dm.DmClob
	APPLY_SCOPE  string
}

func CopyAmListName(outAln OutAmListName, aln AmListName) OutAmListName {
	outAln.ID = aln.ID.String
	outAln.LIST_NAME = aln.LIST_NAME.String
	outAln.BUSINESS_VAR = aln.BUSINESS_VAR.String
	outAln.APP_ID = aln.APP_ID.String
	outAln.LIST_TYPE = aln.LIST_TYPE.String
	outAln.STATES = aln.STATES.String
	outAln.DESCRIPTION = aln.DESCRIPTION.String
	outAln.CREATE_USER = aln.CREATE_USER.String
	outAln.CREATE_TIME = aln.CREATE_TIME.Time
	outAln.UPDATE_USER = aln.UPDATE_USER.String
	outAln.UPDATE_TIME = aln.UPDATE_TIME.Time
	outAln.TRADE_TYPE = aln.TRADE_TYPE
	outAln.APPLY_SCOPE = aln.APPLY_SCOPE.String
	return outAln
}

func queryAmListName(ids string) error {
	aln := AmListName{}
	outAln := OutAmListName{}
	// 字段数据拼接：fields、fieldsTab
	var fields string
	var fieldsTab string
	for i := 0; i < reflect.TypeOf(aln).NumField(); i++ {
		s := reflect.TypeOf(aln).Field(i).Name
		if i == 0 {
			fields = fmt.Sprintf(s)
			fieldsTab = fmt.Sprintf("(" + s)
		} else if i <= reflect.TypeOf(aln).NumField()-1 {
			fields = fmt.Sprintf(fields + "," + s)
			if i == reflect.TypeOf(aln).NumField()-1 {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s + ")")
			} else {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s)
			}
		}
	}

	//解析字段数据
	colsNum := reflect.TypeOf(outAln).NumField()
	//查询所有关联规则的指标数据
	sql := "SELECT " + fields + " FROM AM_LIST_NAME WHERE ID IN (" + ids + ")"
	sqfName := "am_list_name_" + ParamsMp["set_name"] + "_" + time.Time.Format(time.Now(), "20060102150405") + ".sql"
	sqlfile, err := os.Create(sqfName)

	rows, err := db.Query(sql, ids)
	if err != nil {
		fmt.Printf("query listName Table err:%v \n", err)
	}
	defer rows.Close()
	// fmt.Println("queryTable results:")
	for rows.Next() {
		outAln := OutAmListName{}
		record := "NULL"
		if err = rows.Scan(&aln.ID, &aln.LIST_NAME, &aln.BUSINESS_VAR, &aln.APP_ID, &aln.LIST_TYPE, &aln.STATES, &aln.DESCRIPTION, &aln.CREATE_USER, &aln.CREATE_TIME, &aln.UPDATE_USER, &aln.UPDATE_TIME, &aln.TRADE_TYPE, &aln.APPLY_SCOPE); err != nil {
			return err
		}
		res := CopyAmListName(outAln, aln)

		Typ := reflect.TypeOf(res)
		//字段值，&res为了获取对应字段值
		Val := reflect.ValueOf(&res).Elem()

		//遍历输出结构
		for i := 0; i < colsNum; i++ {
			// fmt.Printf("TYPE:%v \n", Typ.Field(i).Type.String())
			val := Val.Field(i)
			//反射获取字段类型，按类型对数据进行归类处理
			if Typ.Field(i).Type.String() == "string" {
				if val.String() == "" {
					// fmt.Printf("STRING:%v \n", "NULL")
					// record = record + "'NULL'"
					record = appendVals(record, "NULL")
					// fmt.Println(record)
				} else {
					// fmt.Printf("STRING:%v \n", val.String())
					// record = record + val.String()
					record = appendVals(record, val.String())
					// fmt.Println(record)
				}
			}
			if Typ.Field(i).Type.String() == "time.Time" {
				if val.Interface().(time.Time).Format("2006-01-02 15:04:05.000") == "0001-01-01 00:00:00.000" {
					// fmt.Printf("TIME:%v \n", "NULL")
					// record = record + "'NULL'"
					record = appendVals(record, "NULL")
					// fmt.Println(record)
				} else {
					// fmt.Printf("TIME:%v \n", val.Interface().(time.Time).Format("2006-01-02 15:04:05.000"))
					// record = record + val.Interface().(time.Time).Format("2006-01-02 15:04:05.000")
					record = appendVals(record, val.Interface().(time.Time).Format("2006-01-02 15:04:05.000"))
					// fmt.Println(record)
				}
			}
			if Typ.Field(i).Type.String() == "*dm.DmClob" {
				if val.Interface().(*dm.DmClob) == nil {
					record = appendVals(record, "NULL")
				} else {
					//读取clob为字符串
					clob := *val.Interface().(*dm.DmClob)
					clobInt64, err := clob.GetLength()
					// fmt.Printf("长度：%v\n", clobInt64)
					if err != nil {
						fmt.Printf("读取Clob长度报错：%v\n", err)
					}
					//将int64转换int32
					strInt64 := strconv.FormatInt(clobInt64, 10)
					Int32, err := strconv.Atoi(strInt64)
					// fmt.Printf("Int32:%v \n", Int32)
					clobToStr, err := clob.ReadString(1, Int32)
					if clobInt64 == 0 {
						record = appendVals(record, "NULL")
					}
					if err == nil {
						record = appendVals(record, clobToStr)
					} else if err == io.EOF {
						fmt.Printf("名单读取Clob结束：%v\n", err)
					} else {
						fmt.Printf("名单数据clob非空解析报错：%v\n", err)
					}

					// fmt.Println(record)
				}
			}
		}
		//拼接inster SQL语句，打印
		isrtSql := "insert into AM_LIST_NAME " + fieldsTab + " values (" + record + ");"
		sqlfile.WriteString(isrtSql + "\n")
	}
	//关闭文件
	defer sqlfile.Close()
	return nil
}
