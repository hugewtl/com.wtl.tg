package main

/*
Author: zhaohu
Date:   2021-09-16 15:36:12
REVISE: 2022-04-14 18:12:23
DESC:   for api rules test more auto
VERSION: v1.1.0
*/
import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gojsonq-2"
)

var (
	//渠道划分
	isApp         bool
	isWeb         bool
	isPOS         bool
	isATM         bool
	isCounter     bool
	isMerchant    bool
	isTradeAmount bool
	//定义报文字段存储数据结构
	fieldMap  = map[string]string{}
	reportMap = map[string]string{}
	isFullAge *bool
	//判断是否过期
	flag bool
	/*
		定义同一主体
	*/
	SameUserID     bool  //同一账号
	SameUserName   bool  //同一登录账号
	SameBankCard   bool  //同一银行卡号
	SameOutAccount bool  //同一出账账号
	SameInAccount  bool  //同一入账账号
	SameIp         bool  //同一IP
	SameDev        bool  //同一设备
	SameIdNo       bool  //同一证件号
	SameOutInName  bool  //同名转账
	SameAtmNo      bool  //同一ATM机具
	SamePosNo      bool  //同一POS机
	SameShopNo     bool  //同一商铺号
	SameMechantNo  bool  //同一商户号
	PhoneNum       bool  //同一手机号
	AmountBit      int32 //额度位数
	IsOpenAcc      bool  //是否开户操作
	NeedReport     bool  //是否上报

	OpState      string //请求状态
	TradeTime    string //敏感时间
	isOldVersion bool   //老版平台
	//配置文件string类型参数map,初始化map
	params_mp = make(map[string]string, 20)
	//定义data.zip路径位置
	path_loc string
	//请求次数
	// count int
	//接口请求地址
	apiUrl string
	repUrl string
	//定义每个场景的请求次数
	ReqTimes int = 1
	//标定当前是第几笔请求
	currentTimes int = 0
	//是否离线处理csv文件
	ifCsv bool
	//是否一次性读取scv，否则逐行读取
	ifLoadCsvAll bool
	//是否开启报文打印日志
	ifDebug bool
)

func main() {
	//获取当月最后一天:时间戳,9月：1632931200
	// fmt.Println(GetLastDateOfMonth(time.Now()).Unix())
	/*
	   base64加密,永久秘钥
	*/
	// strToSha0 := func(src string) string {
	// 	return string(base64.StdEncoding.EncodeToString([]byte(src)))
	// }(fmt.Sprintf("%v", GetLastDateOfMonth(time.Now()).Unix()+1648656000))
	// //加密LIC
	// fmt.Printf("%s: %v TDL: %v ", "lic", strToSha0, fmt.Sprintf("%v", GetLastDateOfMonth(time.Now()).Unix()+1648656000))

	/*
	   base64加密,当月秘钥
	*/
	strToSha1 := func(src string) string {
		return string(base64.StdEncoding.EncodeToString([]byte(src)))
	}(fmt.Sprintf("%v", GetLastDateOfMonth(time.Now()).Unix()))
	//加密LIC
	fmt.Printf("%s: %v TDL: %v  当前时间戳：", "lic", strToSha1, fmt.Sprintf("%v", GetLastDateOfMonth(time.Now()).Unix()))

	/*
	   base64解码：匿名函数封装base64解码逻辑
	*/
	var strToSha string
	shaToStr := func(src string) string {
		/*
		   读取lic
		*/
		lic, err := os.Open("lic") //try to open file
		if err != nil {
			//panic(err)
			fmt.Println("秘钥文件无效")
			panic(err)
		}
		/*
			函数调用结束时关闭文件
		*/
		defer lic.Close()
		//创建文件读取对象
		r := bufio.NewReader(lic)
		for {
			//按行读取
			b, _, err := r.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println("lic读取错误")
				panic(err)

			}
			//去掉已读取的行首位空格
			src := strings.TrimSpace(string(b))
			if len(src) != 0 {
				a, err := base64.StdEncoding.DecodeString(src)
				if err != nil {
					return ""
				}
				strToSha = string(a)
			}
		}
		return strToSha
	}(strToSha) //base64匿名函数解码结束

	//fmt.Println(shaToStr)

	// sha512str := func(str string) string {
	// 	Sha1Inst := sha512.New()
	// 	Sha1Inst.Write([]byte(str))
	// 	return fmt.Sprintf("%x\n\n", Sha1Inst.Sum([]byte("")))
	// }(fmt.Sprintf("%v", GetLastDateOfMonth(time.Now()).Unix()))
	// fmt.Printf("sha512str: %v\n", sha512str)

	/*
		当前时间戳与当月最后一天时间戳比较，判断是否过期
	*/
	fmt.Println(time.Now().Unix())
	unixtime, err := strconv.Atoi(shaToStr)
	if err != nil {
		fmt.Println("lic无效")
		panic(err)
	}
	if time.Now().Unix() <= int64(unixtime) {
		flag = true
	}
	if flag {
		/*
			读取配置文件，获取配置参数
		*/
		file, err := os.Open("init.conf") //try to open file
		if err != nil {
			panic(err)
		}
		/*
			函数调用结束时关闭文件
		*/
		defer file.Close()
		//创建文件读取对象
		r := bufio.NewReader(file)
		for {
			//按行读取
			b, _, err := r.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			//去掉已读取的行首位空格
			str := strings.TrimSpace(string(b))
			//如果以注释符"#"打头，不做处理
			if strings.Index(str, "#") == 0 {
				continue
			} else if len(str) != 0 {
				//参数配置的"="两边值，映射字段
				/*
					获取参数字段名称:k
				*/
				k := strings.TrimSpace(strings.Split(str, "=")[0])
				/*
					获取对应的参数值：v
				*/
				v := strings.TrimSpace(strings.Split(str, "=")[1])
				/*
					将参数赋值存储到map:params_mp
				*/
				if k == "URL" {
					//获取URL请求地址
					params_mp["URL"] = v
					apiUrl = params_mp["URL"] + "/api/risk"
					repUrl = params_mp["URL"] + "/api/logreport/businesslog"
				} else if k == "APPID_SCENES" {
					//获取渠道场景配置列表
					params_mp["APPID_SCENES"] = v
				} else if k == "TradeTime" {
					//获取配置的时间戳
					params_mp["TradeTime"] = v
					TradeTime = params_mp["TradeTime"]
				} else if k == "SameUserID" {
					//获取当前测试配置主体：是否同一账号
					params_mp["SameUserID"] = v
					if params_mp["SameUserID"] == "true" {
						SameUserID = true
					}
				} else if k == "SameUserName" {
					//获取登录账号主体配置
					params_mp["SameUserName"] = v
					if params_mp["SameUserName"] == "true" {
						SameUserName = true
					}
				} else if k == "SameBankCard" {
					//获取当前测试配置主体：是否同一银行卡号
					params_mp["SameBankCard"] = v
					if params_mp["SameBankCard"] == "true" {
						SameBankCard = true
					}
					//fmt.Print(v)
				} else if k == "SameOutAccount" {
					//获取当前测试配置主体：是否同一出账账号
					params_mp["SameOutAccount"] = v
					if params_mp["SameOutAccount"] == "true" {
						SameOutAccount = true
					}
				} else if k == "SameInAccount" {
					//获取当前测试配置主体：是否同一入账账号
					params_mp["SameInAccount"] = v
					if params_mp["SameInAccount"] == "true" {
						SameInAccount = true
					}
				} else if k == "ifDebug" {
					//获取是否打印报文请求及返回日志
					params_mp["ifDebug"] = v
					if params_mp["ifDebug"] == "true" {
						ifDebug = true
					}
				} else if k == "SameIp" {
					//获取当前测试配置主体：是否同一IP
					params_mp["SameIp"] = v
					if params_mp["SameIp"] == "true" {
						SameIp = true
					}
				} else if k == "SameDev" {
					//获取当前测试配置主体：是否同一设备
					params_mp["SameDev"] = v
					if params_mp["SameDev"] == "true" {
						SameDev = true
					}
				} else if k == "SameIdNo" {
					//获取当前测试配置主体：是否同一证件号
					params_mp["SameIdNo"] = v
					if params_mp["SameIdNo"] == "true" {
						SameIdNo = true
					}
				} else if k == "SameOutInName" {
					//获取当前测试配置主体：是否同名转账
					params_mp["SameOutInName"] = v
					if params_mp["SameOutInName"] == "true" {
						SameOutInName = true
					}
				} else if k == "PhoneNum" {
					////获取当前测试配置主体：是否同一手机号
					params_mp["PhoneNum"] = v
					if params_mp["PhoneNum"] == "true" {
						PhoneNum = true
					}
				} else if k == "OpState" {
					//获取当前测试配置主体：是否交易成功
					params_mp["OpState"] = v
					if len(params_mp["OpState"]) != 0 {
						OpState = params_mp["OpState"]
					} else {
						OpState = "1"
					}
				} else if k == "ReqTimes" {
					//获取当前测试配置主体：请求次数
					params_mp["ReqTimes"] = v
					reqtimes, err := strconv.Atoi(params_mp["ReqTimes"])
					if err != nil {
						panic(err)
					}
					ReqTimes = reqtimes
				} else if k == "PATH" {
					//data.zip路径
					params_mp["PATH"] = v
					path_loc = params_mp["PATH"]
				} else if k == "AmountBit" {
					//额度位数获取
					if len(v) != 0 {
						params_mp["AmountBit"] = v
						bit, err := strconv.ParseInt(params_mp["AmountBit"], 10, 32)
						if err != nil {
							panic(err)
						} else {
							AmountBit = int32(bit)
							//fmt.Printf("i32: %v\n", int32(bit))
						}
					} else {
						AmountBit = 11
					}

				} else if k == "Province" {
					//获取省市配置信息
					params_mp["Province"] = v
				} else if k == "City" {
					//获取省市配置信息
					params_mp["City"] = v
				} else if k == "longitude" {
					//获取经度配置信息
					params_mp["longitude"] = v
				} else if k == "latitude" {
					//获取维度配置信息
					params_mp["latitude"] = v
				} else if k == "IsOpenAcc" {
					//获取维度配置信息
					params_mp["IsOpenAcc"] = v
					if params_mp["IsOpenAcc"] == "true" {
						IsOpenAcc = true
					}
				} else if k == "NeedReport" {
					//获取维度配置信息
					params_mp["NeedReport"] = v
					if params_mp["NeedReport"] == "true" {
						NeedReport = true
					}
				} else if k == "SameAtmNo" {
					//同一ATM机具
					params_mp["SameAtmNo"] = v
					if params_mp["SameAtmNo"] == "true" {
						SameAtmNo = true
					}
				} else if k == "SamePosNo" {
					//同一POS机
					params_mp["SamePosNo"] = v
					if params_mp["SamePosNo"] == "true" {
						SamePosNo = true
					}
				} else if k == "SameMechantNo" {
					//同一商户号
					params_mp["SameMechantNo"] = v
					if params_mp["SameMechantNo"] == "true" {
						SameMechantNo = true
					}
				} else if k == "SameShopNo" {
					//同一商铺号
					params_mp["SameShopNo"] = v
					if params_mp["SameShopNo"] == "true" {
						SameShopNo = true
					}
				} else if k == "ifCsv" {
					params_mp["ifCsv"] = v
					if params_mp["ifCsv"] == "true" {
						ifCsv = true
					}
				} else if k == "CsvName" {
					params_mp["CsvName"] = v
				} else if k == "ifLoadCsvAll" {
					params_mp["ifLoadCsvAll"] = v
					if params_mp["ifLoadCsvAll"] == "true" {
						ifLoadCsvAll = true
					}
				} else if k == "dateFormt" { //获取日期时间格式
					params_mp["dateFormt"] = v
				} else if k == "isOldVersion" {
					params_mp["isOldVersion"] = v
					if params_mp["isOldVersion"] == "true" {
						isOldVersion = true
					}
				} else if k == "isApp" { //渠道划分
					params_mp["isApp"] = v
					if params_mp["isApp"] == "true" {
						isApp = true
					}
				} else if k == "isWeb" { //渠道划分
					params_mp["isWeb"] = v
					if params_mp["isWeb"] == "true" {
						isWeb = true
					}
				} else if k == "isATM" { //渠道划分
					params_mp["isATM"] = v
					if params_mp["isATM"] == "true" {
						isATM = true
					}
				} else if k == "isPOS" { //渠道划分
					params_mp["isPOS"] = v
					if params_mp["isPOS"] == "true" {
						isPOS = true
					}
				} else if k == "isCounter" { //渠道划分
					params_mp["isCounter"] = v
					if params_mp["isCounter"] == "true" {
						isCounter = true
					}
				} else if k == "isMerchant" { //渠道划分
					params_mp["isMerchant"] = v
					if params_mp["isMerchant"] == "true" {
						isMerchant = true
					}
				} else if k == "isTradeAmount" { //是否动账交易
					params_mp["isTradeAmount"] = v
					if params_mp["isTradeAmount"] == "true" {
						isTradeAmount = true
					}
				} else if k == "ifFieldsFixed" { //是否补全固定值字段
					params_mp["ifFieldsFixed"] = v
					if params_mp["ifFieldsFixed"] == "true" {
						ifFieldsFixed = true
					}
				} else if k == "separator" { //是否补全固定值字段,只能是单字符
					var runes []rune = []rune(v)
					separator = runes[0]
				} else if k == "isDate" { //是否补全固定值字段,只能是单字符
					isDate = strings.Split(v, ",")
				}

			}
		}

		//是否处理csv离线数据模式
		if !ifCsv {
			/*
			   app_id与scene参数分配映射,遍历
			*/
			for _, as := range strings.Split(params_mp["APPID_SCENES"], ";") {
				//拆分渠道和场景
				// fmt.Println(strings.Split(as, ":")[0])
				// fmt.Println(strings.Split(as, ":")[1])
				//抓取当前渠道标识
				appid := strings.Split(as, ":")[0]
				/*
					遍历当前渠道所有配置的场景
				*/
				for _, scene := range strings.Split(strings.Split(as, ":")[1], ",") {
					/*
						当前渠道，当前场景操作请求req_tims次
					*/
					for n := 1; n <= ReqTimes; n++ {
						//当前请求笔数
						currentTimes = n
						/*
							调用字段处理函数，补全渠道场景参数；业务数据自动生成
						*/
						// fieldDeal("com.ntbank.mobileBank", "NT10003")
						fieldDeal(appid, scene)
						/*
						 发送风险评估请求
						*/
						fmt.Printf("%v%v%v\n", "第", currentTimes, "次请求：")
						reqSend(fieldMap, apiUrl, appid, scene)
						/*
							发送事件上报请求
						*/
						if NeedReport {
							/*
								组事件上报请求报文
							*/
							reportMap["app_id"] = fieldMap["app_id"]
							reportMap["scene"] = fieldMap["scene"]
							reportMap["op_id"] = fieldMap["op_id"]
							reportMap["key_version"] = fieldMap["key_version"]
							reportMap["service_type"] = "1"
							reportMap["user_id"] = fieldMap["user_id"]
							reportMap["user_name"] = fieldMap["user_name"]
							if isOldVersion { //兼容老版本1.0
								reportMap["client_ip"] = fieldMap["client_ip"]
								reportMap["time"] = fieldMap["trade_time"]
							} else {
								reportMap["ip"] = fieldMap["ip"]

							}
							if IsOpenAcc {
								reportMap["bank_card_account_type"] = "2"
								reportMap["op_state"] = "1"
							} else {
								reportMap["op_state"] = "0"
							}

							//事件上报请求
							reqSend(reportMap, repUrl, appid, scene)
						}
					}
				}
			}
		}

		if ifLoadCsvAll {
			wg.Add(1)
			go func() {
				for {
					ok := readAllFromFileToJson(params_mp["CsvName"], ch1)
					if ok {
						break
					}
				}
				wg.Done()
			}()
			wg.Add(1)
			go func() {
				for {
					ok := getRowToApi("fieldFixed.conf", "dataField.conf", ch1)
					if ok {
						break
					}
				}
				wg.Done()
			}()

		} else {
			readLineFromFileToJson(params_mp["CsvName"], "dataField.conf")
		}

	} else {
		fmt.Println("已过有效期")
	} //请求结束
	wg.Wait()
}

// 请求服务json
func reqSend(fieldMp map[string]string, url, appid, scene string) {
	/*
		发送http请求api
	*/
	mjson, _ := json.Marshal(fieldMp)
	mString := string(mjson)
	/*
		输出请求报文json
	*/
	if ifDebug {
		if url[len(url)-4:] == "risk" {
			fmt.Printf("风险评估请求报文:%s\n", mString)
		} else {
			fmt.Printf("事件上报请求报文:%s\n", mString)
		}
	}

	contentType := "application/json"
	resp, err := http.Post(url, contentType, strings.NewReader(mString))
	//fmt.Print(resp)
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read from resp.Body failed, err:%v\n", err)
		return
	}
	/*
		输出返回报文json
	*/

	if url[len(url)-4:] == "risk" {
		/*
			如果是风险评估请求，处理返回报文，获取规则测试命中结果
		*/
		if ifDebug {
			fmt.Println("风险评估返回报文：" + string(body))
		}
		//解析返回报文
		ret := gojsonq.New().FromString(string(body)).Find("result.official_result.hit_rules")
		if isOldVersion {
			ret = gojsonq.New().FromString(string(body)).Find("result.risk_labels")
		}
		/*
			记录命中规则到日志文件
		*/
		// fmt.Println(ret)
		if ret != nil {
			logging_in(ret, appid, scene)
		}
		defer resp.Body.Close()
		//计数
		// count++
	} else {
		if ifDebug {
			fmt.Println("事件上报返回报文：" + string(body))
		}
	}

}

// 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth1(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth1(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// 获取某一天的0点时间
func GetZeroTime1(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}
