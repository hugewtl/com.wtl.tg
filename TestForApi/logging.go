package main

/*
Author: zhaohu
Date:   2021-09-17 17:07:12
REVISE: 2022-04-14 18:12:23
DESC:   for api rules test logging hit rules
VERSION:v1.1.0
*/
import (
	"fmt"
	"os"

	"github.com/go-logging-master"
)

var log = logging.MustGetLogger("example")
var log1 = logging.MustGetLogger("example")

var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02 15:04:05.000} 命中规则: %{color:reset} %{message}`,
)

// type Password string

// func (p Password) Redacted() interface{} {
// 	return logging.Redact(string(p))
// }

func logging_in(info, appid, scene interface{}) {

	logf := "log.txt"
	logFile, err := os.OpenFile(logf, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	//函数调用结束时关闭文件
	defer logFile.Close()
	backend1 := logging.NewLogBackend(logFile, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.INFO, "")

	logging.SetBackend(backend1Leveled, backend2Formatter)

	// log.Debugf("debug %s", Password("secret"))
	// fmt.Printf("%T", info)
	/*
		处理interface{}类型的数组info  []interface{}
	*/
	for _, v := range info.([]interface{}) {
		// fmt.Printf("命中规则：%v\n", v)
		// fmt.Printf("%T\n", v)
		// log.Info(v)
		for k, r := range v.(map[string]interface{}) {
			// fmt.Printf("%v-%v\n", k, r)
			// if k == "rule_no" {
			// 	log.Info(r)
			// }
			/*新平台返回报文命中规则解析*/
			if k == "rule_name" {
				log.Info(appid, scene, r)
			}
			// fmt.Println(k)
			/**兼容老版本命中规则报文解析*/
			if isOldVersion {
				if k == "label_hit_rules" {
					for _, oldR := range r.([]interface{}) {
						for oldRk, oldRr := range oldR.(map[string]interface{}) {
							if oldRk == "rule_name" {
								log.Info(appid, scene, oldRr)
							}
						}
					}
				}
			} //兼容老版本命中规则解析结束

		}
	}
	//log.Info(info)
	// log.Notice("notice")
	// log.Warning("warning")
	// log.Error("Err")
	// log.Critical("critical")
}
func logging_in1(info interface{}) {
	logf := "log_csv.txt"
	logFile, err := os.OpenFile(logf, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	//函数调用结束时关闭文件
	defer logFile.Close()
	backend1 := logging.NewLogBackend(logFile, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.INFO, "")

	logging.SetBackend(backend1Leveled, backend2Formatter)

	// log.Debugf("debug %s", Password("secret"))
	// fmt.Printf("%T", info)
	/*
		处理interface{}类型的数组info  []interface{}
	*/
	for _, v := range info.([]interface{}) {
		// fmt.Printf("命中规则：%v\n", v)
		// fmt.Printf("%T\n", v)
		// log.Info(v)
		for k, r := range v.(map[string]interface{}) {
			// fmt.Printf("%v-%v\n", k, r)
			// if k == "rule_no" {
			// 	log.Info(r)
			// }
			if k == "rule_name" {
				log1.Info(r)
			}
		}
	}
	//log.Info(info)
	// log.Notice("notice")
	// log.Warning("warning")
	// log.Error("Err")
	// log.Critical("critical")
}
