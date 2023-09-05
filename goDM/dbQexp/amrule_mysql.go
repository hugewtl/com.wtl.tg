package main

import (
	"bufio"
	"database/sql"
	"dm"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

/*
定义原规则表结构体，字段列表（目标数据提取）
*/
type myAmRule struct {
	ID                 sql.NullString
	APP_ID             sql.NullString
	TEMPLATE_ID        sql.NullString
	RULESET_HISTORY_ID sql.NullString
	RULE_NO            sql.NullString
	RULE_NAME          sql.NullString
	PARAMS             *dm.DmClob
	RULE_JSON          *dm.DmBlob
	EXPRESSION         sql.NullString
	RULE_TYPE          sql.NullString
	LABEL_ID           sql.NullString
	IS_WEAK            sql.NullString
	STRATEGY           sql.NullString
	ENABLE_TIME        sql.NullTime
	INVALID_TIME       sql.NullTime
	PRIORITY_LEVEL     sql.NullString
	WEIGHTS            sql.NullString
	RISK_THRESHOLD     sql.NullString
	SCRIPT             sql.NullString
	STATES             sql.NullString
	REMARK             sql.NullString
	CREATE_USER        sql.NullString
	CREATE_TIME        sql.NullTime
	UPDATE_USER        sql.NullString
	UPDATE_TIME        sql.NullTime
	DRL                *dm.DmClob
	RULEFLOW_GROUP     sql.NullString
	FLOW_DRL           *dm.DmClob
}

/*
定义表字段输出结构
*/
type myoutAmRule struct {
	ID                 string
	APP_ID             string
	TEMPLATE_ID        string
	RULESET_HISTORY_ID string
	RULE_NO            string
	RULE_NAME          string
	PARAMS             *dm.DmClob
	RULE_JSON          *dm.DmBlob
	EXPRESSION         string
	RULE_TYPE          string
	LABEL_ID           string
	IS_WEAK            string
	STRATEGY           string
	ENABLE_TIME        time.Time
	INVALID_TIME       time.Time
	PRIORITY_LEVEL     string
	WEIGHTS            string
	RISK_THRESHOLD     string
	SCRIPT             string
	STATES             string
	REMARK             string
	CREATE_USER        string
	CREATE_TIME        time.Time
	UPDATE_USER        string
	UPDATE_TIME        time.Time
	DRL                *dm.DmClob
	RULEFLOW_GROUP     string
	FLOW_DRL           *dm.DmClob
}

/*
输出结构与原表结构数据字段映射
*/
func mycopyFieldsVal(outR myoutAmRule, aR myAmRule) myoutAmRule {
	outR.ID = aR.ID.String
	outR.APP_ID = aR.APP_ID.String
	outR.TEMPLATE_ID = aR.TEMPLATE_ID.String
	outR.RULESET_HISTORY_ID = aR.RULESET_HISTORY_ID.String
	outR.RULE_NO = aR.RULE_NO.String
	outR.RULE_NAME = aR.RULE_NAME.String
	outR.PARAMS = aR.PARAMS
	outR.RULE_JSON = aR.RULE_JSON
	outR.EXPRESSION = aR.EXPRESSION.String
	outR.RULE_TYPE = aR.RULE_TYPE.String
	outR.LABEL_ID = aR.LABEL_ID.String
	outR.IS_WEAK = aR.IS_WEAK.String
	outR.STRATEGY = aR.STRATEGY.String
	outR.ENABLE_TIME = aR.ENABLE_TIME.Time
	outR.INVALID_TIME = aR.INVALID_TIME.Time
	outR.PRIORITY_LEVEL = aR.PRIORITY_LEVEL.String
	outR.WEIGHTS = aR.WEIGHTS.String
	outR.RISK_THRESHOLD = aR.RISK_THRESHOLD.String
	outR.SCRIPT = aR.SCRIPT.String
	outR.STATES = aR.STATES.String
	outR.REMARK = aR.REMARK.String
	outR.CREATE_USER = aR.CREATE_USER.String
	outR.CREATE_TIME = aR.CREATE_TIME.Time
	outR.UPDATE_USER = aR.UPDATE_USER.String
	outR.UPDATE_TIME = aR.UPDATE_TIME.Time
	outR.DRL = aR.DRL
	outR.RULEFLOW_GROUP = aR.RULEFLOW_GROUP.String
	outR.FLOW_DRL = aR.FLOW_DRL
	return outR
}

/*查询整个规则集规则*/
func myqueryAllAmRule(setHistId string) error {
	r := myAmRule{}
	outRule := myoutAmRule{}
	// 字段数据拼接：fields、fieldsTab
	var fields string
	var fieldsTab string
	for i := 0; i < reflect.TypeOf(r).NumField(); i++ {
		s := reflect.TypeOf(r).Field(i).Name
		if i == 0 {
			fields = fmt.Sprintf(s)
			fieldsTab = fmt.Sprintf("(" + s)
		} else if i <= reflect.TypeOf(r).NumField()-1 {
			fields = fmt.Sprintf(fields + "," + s)
			if i == reflect.TypeOf(r).NumField()-1 {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s + ")")
			} else {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s)
			}
		}
	}
	// log.Println(fields)
	// log.Println(fieldsTab)
	var sql = "SELECT " + fields + " FROM AM_RULE WHERE RULESET_HISTORY_ID=?"
	// log.Println(sql)
	rows, err := db.Query(sql, setHistId)
	if err != nil {
		return err
	}
	defer rows.Close()

	//生成sqlfile文件
	sqfName := "am_rule_" + ParamsMp["set_name"] + "_" + time.Time.Format(time.Now(), "20060102150405") + ".sql"
	sqlfile, err := os.Create(sqfName)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	//解析字段数据
	colsNum := reflect.TypeOf(outRule).NumField()

	for rows.Next() {
		outRule := myoutAmRule{}
		record := "NULL"
		//遍历字段查询数据
		if err = rows.Scan(&r.ID, &r.APP_ID, &r.TEMPLATE_ID, &r.RULESET_HISTORY_ID, &r.RULE_NO, &r.RULE_NAME, &r.PARAMS, &r.RULE_JSON, &r.EXPRESSION, &r.RULE_TYPE, &r.LABEL_ID, &r.IS_WEAK, &r.STRATEGY, &r.ENABLE_TIME, &r.INVALID_TIME, &r.PRIORITY_LEVEL, &r.WEIGHTS, &r.RISK_THRESHOLD, &r.SCRIPT, &r.STATES, &r.REMARK, &r.CREATE_USER, &r.CREATE_TIME, &r.UPDATE_USER, &r.UPDATE_TIME, &r.DRL, &r.RULEFLOW_GROUP, &r.FLOW_DRL); err != nil {
			return err
		}
		//将查出的一行记录映射到outRule
		res := mycopyFieldsVal(outRule, r)
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
					record = myappendVals(record, "NULL")
					// log.Println(record)
				} else {
					// log.Printf("STRING:%v \n", val.String())
					// record = record + val.String()
					record = myappendVals(record, val.String())
					// log.Println(record)
				}
			}

			if Typ.Field(i).Type.String() == "time.Time" {
				if val.Interface().(time.Time).Format("2006-01-02 15:04:05.000") == "0001-01-01 00:00:00.000" {
					// log.Printf("TIME:%v \n", "NULL")
					// record = record + "'NULL'"
					record = myappendVals(record, "NULL")
					// log.Println(record)
				} else {
					// log.Printf("TIME:%v \n", val.Interface().(time.Time).Format("2006-01-02 15:04:05.000"))
					// record = record + val.Interface().(time.Time).Format("2006-01-02 15:04:05.000")
					record = myappendVals(record, val.Interface().(time.Time).Format("2006-01-02 15:04:05.000"))
					// log.Println(record)
				}
			}
			if Typ.Field(i).Type.String() == "*dm.DmClob" {
				if val.Interface().(*dm.DmClob) == nil {
					// log.Printf("CLOB:%v", "NULL")
					// record = record + "'NULL'"
					record = myappendVals(record, "NULL")
					// log.Println(record)
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
					if err == nil {
						// log.Print(clobToStr)
						record = myappendVals(record, fmt.Sprintf("%v", clobToStr))
					} else {
						log.Println(err)
					}
				}
			}

			if Typ.Field(i).Type.String() == "*dm.DmBlob" {
				if val.Interface().(*dm.DmBlob) == nil {
					// log.Printf("CLOB:%v", "NULL")
					// record = record + "'NULL'"
					record = myappendVals(record, "NULL")
					// log.Println(record)
				} else {
					//读取blob==>字符串
					blob := *val.Interface().(*dm.DmBlob)
					//获取blob对象长度int64==>32
					blobint64, err := blob.GetLength()
					strInt64 := strconv.FormatInt(blobint64, 10)
					blobint32, err := strconv.Atoi(strInt64)
					//定义[]byte数组，用于存放读取到的blob字节
					var dst []byte
					//初始化容量为blob字节长度
					dst = make([]byte, blobint32)
					posInt, err := blob.Read(dst)
					if err != nil {
						log.Println(posInt, err)
					}

					if err == nil {
						// log.Print(clobToStr)
						record = myappendVals(record, fmt.Sprintf("%v", string(dst[:])))
					} else {
						log.Println(err)
					}

					// log.Println(record)
				}
			}
		}
		//拼接inster SQL语句，打印
		isrtSql := "insert into am_rule " + fieldsTab + " values (" + record + ");"
		sqlfile.WriteString(isrtSql + "\n")

	}
	//关闭文件
	defer sqlfile.Close()
	return nil
}

/*
定义目标规则获取方法: 依据rule_no模糊匹配，ruleset_histroy_id精确匹配
*/
func myqueryAmRule(setHistId string, ruleNos string) error {
	r := myAmRule{}
	outRule := myoutAmRule{}
	// 字段数据拼接：fields、fieldsTab
	var fields string
	var fieldsTab string
	for i := 0; i < reflect.TypeOf(r).NumField(); i++ {
		s := reflect.TypeOf(r).Field(i).Name
		if i == 0 {
			fields = fmt.Sprintf(s)
			fieldsTab = fmt.Sprintf("(" + s)
		} else if i <= reflect.TypeOf(r).NumField()-1 {
			fields = fmt.Sprintf(fields + "," + s)
			if i == reflect.TypeOf(r).NumField()-1 {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s + ")")
			} else {
				fieldsTab = fmt.Sprintf(fieldsTab + "," + s)
			}
		}
	}
	// log.Println(fields)
	// log.Println(fieldsTab)
	var sql = "SELECT " + fields + " FROM AM_RULE WHERE RULESET_HISTORY_ID=? AND RULE_NO LIKE ?"
	// log.Println(sql)
	rows, err := db.Query(sql, setHistId, ruleNos)
	if err != nil {
		return err
	}
	defer rows.Close()

	//生成sqlfile文件
	sqfName := "am_rule_" + ParamsMp["set_name"] + "_" + time.Time.Format(time.Now(), "20060102150405") + ".sql"
	sqlfile, err := os.Create(sqfName)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	//解析字段数据
	colsNum := reflect.TypeOf(outRule).NumField()

	for rows.Next() {
		outRule := myoutAmRule{}
		record := "NULL"
		//遍历字段查询数据
		if err = rows.Scan(&r.ID, &r.APP_ID, &r.TEMPLATE_ID, &r.RULESET_HISTORY_ID, &r.RULE_NO, &r.RULE_NAME, &r.PARAMS, &r.RULE_JSON, &r.EXPRESSION, &r.RULE_TYPE, &r.LABEL_ID, &r.IS_WEAK, &r.STRATEGY, &r.ENABLE_TIME, &r.INVALID_TIME, &r.PRIORITY_LEVEL, &r.WEIGHTS, &r.RISK_THRESHOLD, &r.SCRIPT, &r.STATES, &r.REMARK, &r.CREATE_USER, &r.CREATE_TIME, &r.UPDATE_USER, &r.UPDATE_TIME, &r.DRL, &r.RULEFLOW_GROUP, &r.FLOW_DRL); err != nil {
			return err
		}
		//将查出的一行记录映射到outRule
		res := mycopyFieldsVal(outRule, r)
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
					record = myappendVals(record, "NULL")
					// log.Println(record)
				} else {
					// log.Printf("STRING:%v \n", val.String())
					// record = record + val.String()
					record = myappendVals(record, val.String())
					// log.Println(record)
				}
			}

			if Typ.Field(i).Type.String() == "time.Time" {
				if val.Interface().(time.Time).Format("2006-01-02 15:04:05.000") == "0001-01-01 00:00:00.000" {
					// log.Printf("TIME:%v \n", "NULL")
					// record = record + "'NULL'"
					record = myappendVals(record, "NULL")
					// log.Println(record)
				} else {
					// log.Printf("TIME:%v \n", val.Interface().(time.Time).Format("2006-01-02 15:04:05.000"))
					// record = record + val.Interface().(time.Time).Format("2006-01-02 15:04:05.000")
					record = myappendVals(record, val.Interface().(time.Time).Format("2006-01-02 15:04:05.000"))
					// log.Println(record)
				}
			}
			if Typ.Field(i).Type.String() == "*dm.DmClob" {
				if val.Interface().(*dm.DmClob) == nil {
					// log.Printf("CLOB:%v", "NULL")
					// record = record + "'NULL'"
					record = myappendVals(record, "NULL")
					// log.Println(record)
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
					if err == nil {
						// log.Print(clobToStr)
						record = myappendVals(record, fmt.Sprintf("%v", clobToStr))
					} else {
						log.Println(err)
					}
				}
			}

			if Typ.Field(i).Type.String() == "*dm.DmBlob" {
				if val.Interface().(*dm.DmBlob) == nil {
					// log.Printf("CLOB:%v", "NULL")
					// record = record + "'NULL'"
					record = myappendVals(record, "NULL")
					// log.Println(record)
				} else {
					//读取blob==>字符串
					blob := *val.Interface().(*dm.DmBlob)
					//获取blob对象长度int64==>32
					blobint64, err := blob.GetLength()
					strInt64 := strconv.FormatInt(blobint64, 10)
					blobint32, err := strconv.Atoi(strInt64)
					//定义[]byte数组，用于存放读取到的blob字节
					var dst []byte
					//初始化容量为blob字节长度
					dst = make([]byte, blobint32)
					posInt, err := blob.Read(dst)
					if err != nil {
						log.Println(posInt, err)
					}

					if err == nil {
						// log.Print(clobToStr)
						record = myappendVals(record, fmt.Sprintf("%v", string(dst[:])))
					} else {
						log.Println(err)
					}

					// log.Println(record)
				}
			}
		}
		//拼接inster SQL语句，打印
		isrtSql := "insert into am_rule " + fieldsTab + " values (" + record + ");"
		sqlfile.WriteString(isrtSql + "\n")

	}
	//关闭文件
	defer sqlfile.Close()
	return nil
}

// 将字段拼接成字符串，处理字段为空输出逻辑+字段拼接逻辑
func myappendVals(rec string, val string) string {
	if val != "NULL" {
		val = "'" + val + "'"
	}
	if rec != "NULL" {
		rec = rec + "," + val
	} else {
		rec = val
	}
	return rec
}

/*
查询规则集历史版本ID
*/
func myquerHistsetId(setId string, histVersion string) (id string, err error) {
	var sql = "SELECT ID FROM AM_RULESET_HISTORY WHERE SET_ID=? AND HIST_VERSION=?"
	rows, err := db.Query(sql, setId, histVersion)
	if err != nil {
		return id, err
	}
	defer rows.Close()
	// log.Println("queryTable results:")
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return id, err
		}
		// log.Printf(id)
	}
	return id, nil
}

/*
查询规则集ID
*/
func myquerySetId(setName string) (id string, err error) {
	var sql = "SELECT ID FROM AM_RULE_SET WHERE SET_NAME=? AND STATES=1"
	rows, err := db.Query(sql, setName)
	if err != nil {
		return id, err
	}
	defer rows.Close()
	// log.Println("queryTable results:")
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return id, err
		}
	}
	return id, err
}

/* 读取SQL文件，往规则表插入新增规则数据 */
func myinsertAmRule(sqlfile string, targetHistId string, srcHistId string) error {
	infile, err := os.Open(sqlfile)
	if err != nil {
		return err
	}
	defer infile.Close()
	recordSql := bufio.NewReader(infile)

	for { //遍历SQL文件，安行读取SQL，执行插入
		recordSQL, err := recordSql.ReadString('\n')
		if err == io.EOF {
			log.Println("Read SQL file complete")
			break
		}
		if err != nil {
			return err
		}
		//执行insert SQL
		/*执行SQL前，替换SQL中的规则集history id*/
		recordSQL = strings.Replace(recordSQL, srcHistId, targetHistId, -1)
		//后门参数：如果ID冲突，需要将SQL中所有ID手动替换成"abcdefghijklmnopqrstuvwxyz"，程序自动检索替换成uuid
		recordSQL = strings.Replace(recordSQL, "abcdefghijklmnopqrstuvwxyz", uuid.NewV4().String(), -1)
		_, err = db.Exec(recordSQL)
		if err != nil {
			return err
		}
	}

	fmt.Println("Insert AmRule succeed")
	return nil
}
