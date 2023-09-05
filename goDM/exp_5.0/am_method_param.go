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

type OutAmMethodParam struct {
	ID           string
	MID          string
	PARAM_NAME   string
	PARAM_SERIAL string
	PARAM_SRC    *dm.DmClob
	PARAM_TYPE   string
	REMARK       string
	CREATE_USER  string
	CREATE_TIME  string
	UPDATE_USER  string
	UPDATE_TIME  string
	PARAM_FIELD  string
}
type AmMethodParam struct {
	ID           sql.NullString
	MID          sql.NullString
	PARAM_NAME   sql.NullString
	PARAM_SERIAL sql.NullString
	PARAM_SRC    *dm.DmClob
	PARAM_TYPE   sql.NullString
	REMARK       sql.NullString
	CREATE_USER  sql.NullString
	CREATE_TIME  sql.NullString
	UPDATE_USER  sql.NullString
	UPDATE_TIME  sql.NullString
	PARAM_FIELD  sql.NullString
}

func CopyAmMethodParam(OutAmtp OutAmMethodParam, Amtp AmMethodParam) OutAmMethodParam {
	OutAmtp.ID = Amtp.ID.String
	OutAmtp.MID = Amtp.MID.String
	OutAmtp.PARAM_NAME = Amtp.PARAM_NAME.String
	OutAmtp.PARAM_SERIAL = Amtp.PARAM_SERIAL.String
	OutAmtp.PARAM_SRC = Amtp.PARAM_SRC
	OutAmtp.PARAM_TYPE = Amtp.PARAM_TYPE.String
	OutAmtp.REMARK = Amtp.REMARK.String
	OutAmtp.CREATE_USER = Amtp.CREATE_USER.String
	OutAmtp.CREATE_TIME = Amtp.CREATE_TIME.String
	OutAmtp.UPDATE_USER = Amtp.UPDATE_USER.String
	OutAmtp.UPDATE_TIME = Amtp.UPDATE_TIME.String
	OutAmtp.PARAM_FIELD = Amtp.PARAM_FIELD.String
	return OutAmtp
}

func queryAmMethodParam(ids string) error {
	Amtp := AmMethodParam{}
	OutAmtp := OutAmMethodParam{}
	// 字段数据拼接：fields、fieldsTab
	var fields string
	var fieldsTab string
	for i := 0; i < reflect.TypeOf(Amtp).NumField(); i++ {
		s := reflect.TypeOf(Amtp).Field(i).Name
		if i == 0 {
			fields = fmt.Sprintf(s)
			fieldsTab = fmt.Sprintf("(" + s)
		} else if i <= reflect.TypeOf(Amtp).NumField()-1 {
			fields = fmt.Sprintf(fields + "," + s)
			if i == reflect.TypeOf(Amtp).NumField()-1 {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s + ")")
			} else {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s)
			}
		}
	}

	//解析字段数据
	colsNum := reflect.TypeOf(OutAmtp).NumField()
	//查询所有关联规则的指标数据
	sql := "SELECT " + fields + " FROM AM_METHOD_PARAM WHERE MID IN (" + ids + ")"
	sqfName := "am_method_param.sql"
	sqlfile, err := os.Create(sqfName)

	rows, err := db.Query(sql, ids)
	if err != nil {
		fmt.Printf("query AM_METHOD_PARAM Table err:%v \n", err)
	}
	defer rows.Close()
	// fmt.Println("queryTable results:")
	for rows.Next() {
		OutAmtp := OutAmMethodParam{}
		record := "NULL"
		if err = rows.Scan(&Amtp.ID, &Amtp.MID, &Amtp.PARAM_NAME, &Amtp.PARAM_SERIAL, &Amtp.PARAM_SRC, &Amtp.PARAM_TYPE, &Amtp.REMARK, &Amtp.CREATE_USER, &Amtp.CREATE_TIME, &Amtp.UPDATE_USER, &Amtp.UPDATE_TIME, &Amtp.PARAM_FIELD); err != nil {
			return err
		}
		res := CopyAmMethodParam(OutAmtp, Amtp)

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
						fmt.Printf("指标方法参数读取Clob结束：%v\n", err)
					} else {
						fmt.Printf("指标方法参数数据clob非空解析报错：%v\n", err)
					}

					// fmt.Println(record)
				}
			}
		}
		//拼接inster SQL语句，打印
		isrtSql := "insert into AM_METHOD_PARAM " + fieldsTab + " values (" + record + ");"
		sqlfile.WriteString(isrtSql + "\n")
	}
	//关闭文件
	defer sqlfile.Close()
	return nil
}
