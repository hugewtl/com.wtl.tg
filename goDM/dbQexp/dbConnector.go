package main

/*该例程实现插入数据，修改数据，删除数据，数据查询等基本操作。*/
// 引入相关包
import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// 全局变量
var (
	db             *sql.DB
	err            error
	driverName     string
	dataSourceName string
	id             string
	hid            string
	filename       string
)

func main() {
	//读取配置文件
	initParams("conf.properties")

	//初始化日志文件
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logfile, err := os.OpenFile(ParamsMp["logfile"], os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}
	defer func() {
		logfile.Close()
	}()

	log.SetOutput(logfile)

	//数据库驱动选择：dm|mysql
	if ParamsMp["dbtype"] == "dm" {
		driverName = "dm"
		dataSourceName = "dm://" + ParamsMp["dbuser"] + ":" + escapeForSC(ParamsMp["dbpasswd"]) + "@" + ParamsMp["ip_port"] + "?clobAsString=true?schema=" + ParamsMp["schema"]
	}
	if ParamsMp["dbtype"] == "mysql" {
		driverName = "mysql"
		dataSourceName = ParamsMp["dbuser"] + ":" + escapeForSC(ParamsMp["dbpasswd"]) + "@tcp(" + ParamsMp["ip_port"] + ")/" + ParamsMp["schema"] + "?charset=utf8mb4&parseTime=true"
		// dataSourceName = "root:Mysql123..@tcp(192.168.1.198:3306)/ipa_manage_qd?charset=utf8mb4&parseTime=true"
	}

	//连接数据库
	if db, err = connect(driverName, dataSourceName); err != nil {
		log.Printf("连接 %s 失败.：%v", dataSourceName, err)
	}
	/**
	  从源端导出目标数据逻辑代码：
	*/
	if ParamsMp["database"] == "src" {
		//查询目标规则集ID
		id, err = querySetId(ParamsMp["set_name"])
		if err != nil {
			log.Printf("查询规则集ID异常: %v", err)
		}
		log.Printf("规则集ID： %v \n", id)
		//查询历史版本规则集ID
		hid, err = querHistsetId(id, ParamsMp["set_hist_version"])
		if err != nil {
			log.Printf("查询历史版本规则集ID异常，%v \n", err)
		}
		log.Printf("规则集历史版本ID： %v \n", hid)

		/*区分驱动dm|mysql匹配查询，导出规则集SQL文件*/
		if ParamsMp["dbtype"] == "mysql" {
			if ParamsMp["if_all_rule"] == "yes" {
				if err = myqueryAllAmRule(hid); err != nil {
					log.Printf("MySQL:查询历史版本规则集异常：%v \n", err)
				}
			} else {
				if err = myqueryAmRule(hid, ParamsMp["rule_no_prefix"]); err != nil {
					log.Printf("MySQL:查询历史版本规则集异常：%v \n", err)
				}
			}

		}
		if ParamsMp["dbtype"] == "dm" {

			if ParamsMp["if_all_rule"] == "yes" {
				if err = queryAllAmRule(hid); err != nil {
					log.Printf("DM:查询历史版本规则集异常：%v \n", err)
				}
			} else {
				if err = queryAmRule(hid, ParamsMp["rule_no_prefix"]); err != nil {
					log.Printf("DM:查询历史版本规则集异常：%v \n", err)
				}
			}

		}
		//查询,导出规则关联的指标，导出SQL：am_indicators
		if len(Indic_list) != 0 {
			if err := QueryIndicators(Indic_list); err != nil {
				log.Printf("查询indic err=%v\n", err)
			}
		}

		//查询,导出相关的所有风险标签数据：am_sys_params
		if len(Label_list) != 0 {
			log.Printf("关联的风险标签ID:%v \n", Label_list)
			if err := QueryAmLabel(Label_list); err != nil {
				log.Printf("查询label err=%v\n", err)
			}
		}
		//嵌套指标导出
		if len(Indic_list_nest) != 0 {
			if err = QueryNestIndicators(Indic_list_nest); err != nil {
				log.Printf("查询indic err=%v\n", err)
			}
		}

		//查询名单
		if len(Name_list) != 0 {
			log.Printf("关联的名单集ID：%v \n", Name_list)
			if err = queryAmListName(Name_list); err != nil {
				log.Printf("查询名单集 err=%v\n", err)
			}
		}

		//查询业务字段
		if len(Fields_list) != 0 {
			log.Printf("关联的字段名称：%v \n", Fields_list)
		}

		//查询指标模板
		if err = queryAmTemplate(ParamsMp["indic_template_name"]); err != nil {
			log.Printf("导出指标模板失败，err%v \n", err)
		}
		//查询指标模板参数
		if len(Tmpl_id) != 0 {
			log.Printf("指标模板ID：%v \n", Tmpl_id)
			if err = queryAmTemplateParam(Tmpl_id); err != nil {
				log.Printf("导出指标模板参数失败，err：%v \n", err)
			}
		}

		//查询方法模板
		if err = queryAmMethod(ParamsMp["indic_method_name"]); err != nil {
			log.Printf("导出指标方法失败，err%v \n", err)
		}
		//查询方法模板参数
		if len(Method_id) != 0 {
			log.Printf("方法模板ID：%v \n", Method_id)
			if err = queryAmMethodParam(Method_id); err != nil {
				log.Printf("导出指标方法模板参数失败，err：%v \n", err)
			}
		}

	}

	/**
	  如果将规则集文件导入目标规则集
	*/
	if ParamsMp["database"] == "target" {
		id, err = querySetId(ParamsMp["set_name"])
		if err != nil {
			log.Printf("查询目标规则集ID异常 %v \n", err)
		}
		log.Printf("目标规则集ID： %v \n", id)

		hid, err = querHistsetId(id, ParamsMp["set_hist_version"])
		if err != nil {
			log.Printf("查询历史版本规则集ID异常 %v \n", err)
		}
		log.Printf("目标规则集关联历史版本ID： %v \n", hid)
		//读取SQL文件，加载sql语句，插入规则表。
		if err = insertAmRule(ParamsMp["sqlfile"], hid, ParamsMp["srcHistId"]); err != nil {
			log.Printf("插入规则异常 %v \n", err)
		}
		log.Println("Rules imported succeed")
	}

	//关闭数据库连接
	if err = disconnect(); err != nil {
		log.Println(err)
		// return err
	}
	// return nil
}

/* 创建数据库连接 */
func connect(driverName string, dataSourceName string) (*sql.DB, error) {
	// var db *sql.DB
	// var err error
	if db, err = sql.Open(driverName, dataSourceName); err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Printf("connect to \"%s\" succeed.\n", dataSourceName)

	return db, nil
}

/* 关闭数据库连接 */
func disconnect() error {
	if err = db.Close(); err != nil {
		log.Printf("db close failed: %s.\n", err)
		return err
	}
	log.Println("disconnect succeed")
	return nil
}

// /* 往产品信息表插入数据 */
// func insertTable() error {
// 	var inFileName = "sanguo.txt"
// 	var sql = `INSERT INTO production.product(name,author,publisher,publishtime,
//                                 product_subcategoryid,productno,satetystocklevel,originalprice,nowprice,discount,
//                                 description,photo,type,papertotal,wordtotal,sellstarttime,sellendtime)
//                 VALUES(:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12,:13,:14,:15,:16,:17);`
// 	data, err := ioutil.ReadFile(inFileName)
// 	if err != nil {
// 		return err
// 	}
// 	t1, _ := time.Parse("2006-Jan-02", "2005-Apr-01")
// 	t2, _ := time.Parse("2006-Jan-02", "2006-Mar-20")
// 	t3, _ := time.Parse("2006-Jan-02", "1900-Jan-01")
// 	_, err = db.Exec(sql, "三国演义", "罗贯中", "中华书局", t1, 4, "9787101046121", 10, 19.0000, 15.2000,
// 		8.0,
// 		"《三国演义》是中国第一部长篇章回体小说，中国小说由短篇发展至长篇的原因与说书有关。",
// 		data, "25", 943, 93000, t2, t3)
// 	if err != nil {
// 		return err
// 	}
// 	log.Println("insertTable succeed")
// 	return nil
// }

// /* 修改产品信息表数据 */
// func updateTable() error {
// 	var sql = "UPDATE production.product SET name = :name WHERE productid = 11;"
// 	if _, err := db.Exec(sql, "三国演义（上）"); err != nil {
// 		return err
// 	}
// 	log.Println("updateTable succeed")
// 	return nil
// }

/* 查询am_rule规则表 */
// func queryTable() error {
// 	var rule_no string
// 	var rule_name string
// 	var rule_json dm.DmClob

// 	var sql = "select rule_no,rule_name,rule_json from am_rule where ruleset_history_id=(select id from am_ruleset_history where set_id=(select id from am_rule_set where set_name=?) and hist_version=?) "

// 	rows, err := db.Query(sql, "无卡支付签约（银联渠道）", 3)

// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()

// 	log.Println("queryTable results:")

// 	for rows.Next() {
// 		if err = rows.Scan(&ID, &APP_ID, &TEMPLATE_ID, &RULESET_HISTORY_ID, &RULE_NO, &RULE_NAME, &PARAMS, &RULE_JSON, &EXPRESSION, &RULE_TYPE, &LABEL_ID, &IS_WEAK, &STRATEGY, &ENABLE_TIME, &INVALID_TIME, &PRIORITY_LEVEL, &WEIGHTS, &RISK_THRESHOLD, &SCRIPT, &STATES, &REMARK, &CREATE_USER, &CREATE_TIME, &UPDATE_USER, &UPDATE_TIME, &DRL, &RULEFLOW_GROUP, &FLOW_DRL); err != nil {
// 			return err
// 		}
// 		log.Printf("insert into am_rule (\"%v\",\"%v\",\"%v\");", rule_no, rule_name, rule_json)
// 	}
// 	return nil
// }

// /* 删除产品信息表数据 */
// func deleteTable() error {
// 	var sql = "DELETE FROM production.product WHERE productid = 11;"
// 	if _, err := db.Exec(sql); err != nil {
// 		return err
// 	}
// 	log.Println("deleteTable succeed")
// 	return nil
// }

// if err = insertTable(); err != nil {
// 	log.Println(err)
// 	return
// }
// if err = updateTable(); err != nil {
// 	log.Println(err)
// 	return
// }
// if err = queryTable(); err != nil {
// 	log.Println(err)
// 	return
// }
// if err = deleteTable(); err != nil {
// 	log.Println(err)
// 	return
// }
