package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/olongfen/gen-id/generator"
)

func fieldDeal(appid string, scene string) {
	/*
		json报文所需常用字段赋值
	*/

	//生成出账人信息
	out := new(generator.GeneratorData)
	out.GeneratorPhone()
	out.GeneratorName()
	out.GeneratorIDCart(isFullAge)
	out.GeneratorEmail()
	out.GeneratorBankID()
	out.GeneratorAddress()
	//生成入账人信息
	in := new(generator.GeneratorData)
	in.GeneratorPhone()
	in.GeneratorName()
	in.GeneratorIDCart(isFullAge)
	in.GeneratorEmail()
	in.GeneratorBankID()
	in.GeneratorAddress()
	if OpState == "1" {
		//交易状态:成功
		fieldMap["op_state"] = "1"
	} else if OpState == "0" {
		//交易状态：失败
		fieldMap["op_state"] = "0"
	}

	//渠道场景信息
	fieldMap["app_id"] = appid //"com.ntbank.mobileBank"
	fieldMap["scene"] = scene  //"NT10003"
	//出账人信息
	/*
		同一账号操作赋值
	*/
	if SameUserID {
		fieldMap["user_id"] = "SameUserId"
	} else {
		fieldMap["user_id"] = out.PhoneNum
	}
	//获取当前时间的纳秒作为操作标识唯一值
	fieldMap["op_id"] = strconv.FormatInt(time.Now().UnixNano(), 10)
	/*
		同名操作
	*/
	if SameOutInName {
		//出账人姓名
		fieldMap["user_name"] = out.Name
		fieldMap["name"] = out.Name
		//入账人姓名
		fieldMap["in_name"] = out.Name
	} else {
		//出账人姓名
		fieldMap["user_name"] = out.Name
		fieldMap["name"] = out.Name
		//入账人姓名
		fieldMap["in_name"] = in.Name
	}
	/*
		同一银行卡操作赋值
	*/
	if SameBankCard {
		fieldMap["bank_card"] = "6226660205067749"
	} else {
		fieldMap["bank_card"] = out.BankID
	}
	/*
		同一出账账号操作赋值
	*/
	if SameOutAccount {
		fieldMap["out_account"] = "6226660205067751"
	} else {
		fieldMap["out_account"] = out.BankID
	}

	fieldMap["amount"] = "11111"
	fieldMap["bank_no"] = randValue(6)
	/*
		将交易时间转换为unix时间戳
	*/
	if len(TradeTime) != 0 {
		//配置时间格式
		layout := "2006-01-02 15:04:05" //固定值，go的特殊格式
		/*
			固定时区：中国
		*/
		loc, _ := time.LoadLocation("Asia/Shanghai")
		t, err := time.ParseInLocation(layout, TradeTime, loc)
		if err != nil {
			//打印错误信息
			fmt.Print(err)
		} else {
			fmt.Print("交易时间：")
			fmt.Println(t) //2021-09-16 01:32:20 +0800 CST
			//fmt.Println(t.Unix()) //1631727140
			fieldMap["trade_time"] = strconv.FormatInt(t.Unix(), 10)
		}

	} else {
		//当前时间戳
		fieldMap["trade_time"] = strconv.FormatInt(time.Now().Unix(), 10)
	}

	fieldMap["trade_no"] = randValue(16)
	fieldMap["overseas_trade"] = "true"
	fieldMap["id_type"] = "1"
	fieldMap["id_no"] = out.IDCard
	fieldMap["phone"] = out.PhoneNum

	/*
		同一入账账号操作赋值
	*/
	if SameInAccount {
		fieldMap["in_account"] = "6226660205067752"
	} else {
		fieldMap["in_account"] = in.BankID
	}
	//设备与采集信息
	fieldMap["device_type"] = "Android"
	/*
		同一IP操作赋值
	*/
	if SameIp {
		fieldMap["ip"] = "56.23.52.43"
		fieldMap["wifi_mac"] = "00:23:89:b9:1a:fc"
	} else {
		fieldMap["ip"] = "56.23.52.1" + randValue2(2)
		fieldMap["wifi_mac"] = "00:23:89:b9:1a:" + randValue2(2)
	}
	/*
		同一设备操作赋值
	*/
	if SameDev {
		fieldMap["dev_fp"] = "9C613307C1AF33CC1BE626BEBE529444"
	} else {
		fieldMap["dev_fp"] = "9C613307C1AF33CC1BE626BEBE" + randValue(6)
	}
	//位置信息
	fieldMap["province"] = "广东省"
	fieldMap["city"] = "广州市"
	//系统信息
	fieldMap["key_version"] = "01"
	fieldMap["service_type"] = "0"
	fieldMap["signature"] = ""
}

func randValue(lens int32) string {
	/*
		初始化随机值：获取6位随机数字的字符串
	*/
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n((lens-1)<<1))
}
func randValue2(lens int32) string {
	/*
		初始化随机值：获取6位随机数字的字符串
	*/
	return fmt.Sprintf("%02v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n((lens-1)<<1))
}
