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

type OutAmMethod struct {
	ID          string
	GID         string
	FUN_NAME    string
	REMARK      string
	CREATE_TIME time.Time
	CREATE_USER string
	DATA_TYPE   string
	MTHOD_NAME  string
}

type AmMethod struct {
	ID          sql.NullString
	GID         sql.NullString
	FUN_NAME    sql.NullString
	REMARK      sql.NullString
	CREATE_TIME sql.NullTime
	CREATE_USER sql.NullString
	DATA_TYPE   sql.NullString
	MTHOD_NAME  sql.NullString
}

func CopyAmMethod(OutAmt OutAmMethod, Amt AmMethod) OutAmMethod {
	OutAmt.ID = Amt.ID.String
	OutAmt.GID = Amt.GID.String
	OutAmt.FUN_NAME = Amt.FUN_NAME.String
	OutAmt.REMARK = Amt.REMARK.String
	OutAmt.CREATE_TIME = Amt.CREATE_TIME.Time
	OutAmt.CREATE_USER = Amt.CREATE_USER.String
	OutAmt.DATA_TYPE = Amt.DATA_TYPE.String
	OutAmt.MTHOD_NAME = Amt.MTHOD_NAME.String
	return OutAmt
}

func queryAmMethod(mName string) error {
	Amt := AmMethod{}
	OutAmt := OutAmMethod{}
	// 字段数据拼接：fields、fieldsTab
	var fields string
	var fieldsTab string
	for i := 0; i < reflect.TypeOf(Amt).NumField(); i++ {
		s := reflect.TypeOf(Amt).Field(i).Name
		if i == 0 {
			fields = fmt.Sprintf(s)
			fieldsTab = fmt.Sprintf("(" + s)
		} else if i <= reflect.TypeOf(Amt).NumField()-1 {
			fields = fmt.Sprintf(fields + "," + s)
			if i == reflect.TypeOf(Amt).NumField()-1 {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s + ")")
			} else {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s)
			}
		}
	}

	//解析字段数据
	colsNum := reflect.TypeOf(OutAmt).NumField()
	//查询所有关联规则的指标数据
	sql := "SELECT " + fields + " FROM AM_METHOD WHERE MTHOD_NAME IN (" + mName + ")"
	sqfName := "am_method.sql"
	sqlfile, err := os.Create(sqfName)

	rows, err := db.Query(sql, mName)
	if err != nil {
		log.Printf("query AM_METHOD Table err:%v \n", err)
	}
	defer rows.Close()
	// log.Println("queryTable results:")
	for rows.Next() {
		OutAmt := OutAmMethod{}
		record := "NULL"
		if err = rows.Scan(&Amt.ID, &Amt.GID, &Amt.FUN_NAME, &Amt.REMARK, &Amt.CREATE_TIME, &Amt.CREATE_USER, &Amt.DATA_TYPE, &Amt.MTHOD_NAME); err != nil {
			return err
		}
		res := CopyAmMethod(OutAmt, Amt)

		Typ := reflect.TypeOf(res)
		//字段值，&res为了获取对应字段值
		Val := reflect.ValueOf(&res).Elem()

		//提取方法模板ID
		if len(Method_id) == 0 {
			Method_id = "'" + Amt.ID.String + "'"
		} else {
			Method_id = appendVals(Method_id, Amt.ID.String)
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
						fmt.Printf("指标方法读取Clob结束：%v\n", err)
					} else {
						fmt.Printf("指标方法数据clob非空解析报错：%v\n", err)
					}

					// log.Println(record)
				}
			}
		}
		//拼接inster SQL语句，打印
		isrtSql := "insert into AM_METHOD " + fieldsTab + " values (" + record + ");"
		sqlfile.WriteString(isrtSql + "\n")
	}
	//关闭文件
	defer sqlfile.Close()
	return nil
}
