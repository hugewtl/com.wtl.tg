package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	//定义报文字段存储数据结构
	fieldMap  = map[string]string{}
	isFullAge *bool
	/*
		定义同一主体
	*/
	SameUserID     bool   //同一账号
	SameBankCard   bool   //同一银行卡号
	SameOutAccount bool   //同一出账账号
	SameInAccount  bool   //同一入账账号
	SameIp         bool   //同一IP
	SameDev        bool   //同一设备
	SameIdNo       bool   //同一证件号
	SameOutInName  bool   //同名转账
	OpState        string //失败请求
	TradeTime      string //敏感时间
	//配置文件string类型参数map
	params_mp = make(map[string]string, 20)
)

func main() {
	var (
		//定义每个场景的请求次数
		ReqTimes int = 1
	)

	/*
		读取配置文件，获取配置参数
	*/
	file, err := os.Open("init.conf") //try to open file
	if err != nil {
		panic(err)
	}
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
			//渠道"="两边值，映射字段
			k := strings.TrimSpace(strings.Split(str, "=")[0])
			v := strings.TrimSpace(strings.Split(str, "=")[1])
			if k == "URL" {
				params_mp["URL"] = v
			} else if k == "APPID_SCENES" {
				params_mp["APPID_SCENES"] = v
			} else if k == "TradeTime" {
				params_mp["TradeTime"] = v
				TradeTime = params_mp["TradeTime"]
			} else if k == "SameUserID" {
				params_mp["SameUserID"] = v
			} else if k == "SameBankCard" {
				params_mp["SameBankCard"] = v
				//fmt.Print(v)
			} else if k == "SameOutAccount" {
				params_mp["SameOutAccount"] = v
			} else if k == "SameInAccount" {
				params_mp["SameInAccount"] = v
			} else if k == "SameIp" {
				params_mp["SameIp"] = v
			} else if k == "SameDev" {
				params_mp["SameDev"] = v
			} else if k == "SameIdNo" {
				params_mp["SameIdNo"] = v
			} else if k == "SameOutInName" {
				params_mp["SameOutInName"] = v
			} else if k == "OpState" {
				params_mp["OpState"] = v
			} else if k == "ReqTimes" {
				params_mp["ReqTimes"] = v
				reqtimes, err := strconv.Atoi(params_mp["ReqTimes"])
				if err != nil {
					panic(err)
				}
				ReqTimes = reqtimes
			}
		}
	}
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
				/*
					分配同一主体请求
				*/
				if params_mp["SameIp"] == "true" {
					//同一IP
					SameIp = true
				}
				if params_mp["SameUserID"] == "true" {
					//同一账号操作
					SameUserID = true
				}
				if params_mp["SameDev"] == "true" {
					//同一设备操作
					SameDev = true
				}
				if params_mp["SameBankCard"] == "true" {
					//同一银行卡号
					SameBankCard = true
				}
				if params_mp["SameOutAccount"] == "true" {
					//同一出账账号
					SameOutAccount = true
				}
				if params_mp["SameInAccount"] == "true" {
					//同一入账账号
					SameInAccount = true
				}
				if params_mp["SameIdNo"] == "true" {
					//同一证件号
					SameIdNo = true
				}
				if params_mp["SameOutInName"] == "true" {
					//同名转账
					SameOutInName = true
				}
				if params_mp["OpState"] == "0" {
					//失败操作
					OpState = "0"
				}

				/*
				   调用字段处理函数，补全渠道场景参数；业务数据自动生成
				*/
				// fieldDeal("com.ntbank.mobileBank", "NT10003")
				fieldDeal(appid, scene)
				/*
					发送请求
				*/
				reqSend(params_mp["URL"])
			} //请求结束
		}

	}

}

//请求服务
func reqSend(URL string) {
	/*
		发送http请求api
	*/
	mjson, _ := json.Marshal(fieldMap)
	mString := string(mjson)
	fmt.Printf("请求报文:%s\n", mString)
	url := URL
	contentType := "application/json"
	resp, err := http.Post(url, contentType, strings.NewReader(mString))
	//fmt.Print(resp)
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read from resp.Body failed, err:%v\n", err)
		return
	}
	fmt.Println("返回报文：" + string(body))
}
