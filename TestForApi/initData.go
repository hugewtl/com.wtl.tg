package main

/*
Author: zhaohu
Date:   2021-09-16 15:36:12
REVISE: 2021-10-14 16:12:23
DESC:   let test data generation much more auto
VERSION: v1.0.5
*/
import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olongfen/gen-id/generator"
)

var ( //同一主体固定当前批次请求固定值
	SameUserIDVal   = randValue(100000000, "%08v")
	SameUserNameVal = randValue(1000000, "%06v")
	SameAmountVal   = randValue(1000, "%03v")
)

func fieldDeal(appid, scene string) {
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

	//生成开户信息
	opacc := new(generator.GeneratorData)
	opacc.GeneratorPhone()
	opacc.GeneratorName()
	opacc.GeneratorIDCart(isFullAge)
	opacc.GeneratorEmail()
	opacc.GeneratorBankID()
	opacc.GeneratorAddress()

	if len(OpState) != 0 {
		//交易状态:匹配配置参数
		fieldMap["op_state"] = OpState
	} else {
		//交易状态：默认成功
		fieldMap["op_state"] = "1"
	}

	//渠道场景信息
	fieldMap["app_id"] = appid //"com.ntbank.mobileBank"
	fieldMap["scene"] = scene  //"NT10003"
	//出账人信息
	/*
		同一账号操作赋值
	*/
	if SameUserID {
		fieldMap["user_id"] = SameUserIDVal + params_mp["ReqTimes"]
	} else {
		fieldMap["user_id"] = out.PhoneNum
	}
	if SameUserName {
		fieldMap["user_name"] = SameUserNameVal + params_mp["ReqTimes"]
	} else {
		fieldMap["user_name"] = out.Name
	}
	//获取当前时间的纳秒作为操作标识唯一值
	fieldMap["op_id"] = strconv.FormatInt(time.Now().UnixNano(), 10)

	/*
		同一银行卡操作赋值
	*/

	if SameBankCard {
		fieldMap["bank_card"] = "622666020506775" + params_mp["ReqTimes"]
	} else {
		fieldMap["bank_card"] = out.BankID
	}
	//开户卡号
	if IsOpenAcc {
		reportMap["bank_card"] = opacc.BankID
		//开户备用银行卡字段
		fieldMap["bank_card_bak"] = ""
		fieldMap["bank_card_bak_account_type"] = ""
	}
	if isTradeAmount {
		/*
			同一出账账号操作赋值
		*/
		if SameOutAccount {
			fieldMap["out_account"] = "622666020506775" + params_mp["ReqTimes"]
		} else {
			fieldMap["out_account"] = out.BankID
		}
		/*
			同一入账账号操作赋值
		*/
		if SameInAccount {
			fieldMap["in_account"] = "622666020506775" + params_mp["ReqTimes"]
		} else {
			fieldMap["in_account"] = in.BankID
		}
		/*
			同名操作
		*/
		if SameOutInName {
			//出账人姓名
			//fieldMap["user_name"] = out.Name
			fieldMap["name"] = out.Name
			//入账人姓名
			fieldMap["in_name"] = out.Name
		} else {
			//出账人姓名
			//fieldMap["user_name"] = out.Name
			fieldMap["name"] = out.Name
			//入账人姓名
			fieldMap["in_name"] = in.Name
		}

		//amount金额随机数方法
		switch AmountBit {
		case 0:
			/*
				获取随机两位小数
			*/
			amt, err := strconv.Atoi(randValue(100, "%02v"))
			if err != nil {
				panic(err)
			}
			//生成2位小数
			amount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(amt)/float64(100)), 64)
			/*
				将float64转换为字符串赋值给amount
			*/
			fieldMap["amount"] = fmt.Sprintf("%v", amount)
		case 1:
			fieldMap["amount"] = "1" + randValue(10, "%01v")
		case 2:
			fieldMap["amount"] = "1" + randValue(100, "%02v")
		case 3:
			fieldMap["amount"] = "1" + randValue(1000, "%03v")
		case 4:
			fieldMap["amount"] = "1" + randValue(10000, "%04v")
		case 5:
			fieldMap["amount"] = "1" + randValue(100000, "%05v")
		case 6:
			fieldMap["amount"] = "1" + randValue(1000000, "%06v")
		case 7:
			fieldMap["amount"] = "1" + randValue(10000000, "%07v")
		case 8:
			fieldMap["amount"] = "1" + randValue(100000000, "%08v")
		case 9:
			fieldMap["amount"] = "1" + randValue(1000000000, "%09v")
		case 10:
			fieldMap["amount"] = SameAmountVal
		default:
			fieldMap["amount"] = "1" + randValue(1000, "%03v")
		}
	}

	fieldMap["bank_no"] = randValue(1000000, "%06v")
	//位置信息
	/*
		将交易时间转换为unix时间戳
	*/
	if len(TradeTime) != 0 {
		//配置时间格式
		layout := "2006-01-02 15:04:05" //固定值，go的特殊格式
		/*
			固定时区：中国
		*/
		os.Setenv("ZONEINFO", path_loc)
		loc, _ := time.LoadLocation("Asia/Shanghai")
		t, err := time.ParseInLocation(layout, TradeTime, loc)
		if err != nil {
			//打印错误信息
			fmt.Print(err)
		} else {
			fmt.Print("交易时间：")
			fmt.Print(t) //2021-09-16 01:32:20 +0800 CST
			fieldMap["trade_time"] = strconv.FormatInt(t.UnixNano()/1e6, 10)
			fmt.Printf("%s , %v , %s", " ", fieldMap["trade_time"], " ") //1631727140000
		}

	} else {
		//当前时间戳
		fieldMap["trade_time"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	}

	fieldMap["trade_no"] = randValue(100000000, "%08v")
	//位置信息
	fieldMap["overseas_trade"] = "true"
	fieldMap["id_type"] = "1"
	if SameIdNo {
		fieldMap["id_no"] = "45080419780808617" + params_mp["ReqTimes"]
	} else {
		fieldMap["id_no"] = out.IDCard
	}

	/*
		同一手机号操作赋值
	*/
	if PhoneNum {
		fieldMap["phone"] = "1851898998" + params_mp["ReqTimes"]
	} else {
		fieldMap["phone"] = out.PhoneNum
	}

	/*
		同一IP操作赋值
	*/

	if SameIp {
		if isOldVersion {
			fieldMap["client_ip"] = "56.23.52.1" + params_mp["ReqTimes"]
		} else {
			fieldMap["ip"] = "56.23.52.1" + params_mp["ReqTimes"]
		}
		fieldMap["wifi_mac"] = "00:23:89:b9:1a:fc"
	} else {
		if isOldVersion {
			fieldMap["client_ip"] = "56.23.52.1" + randValue(100, "%02v")
		} else {
			fieldMap["ip"] = "56.23.52.1" + randValue(100, "%02v")
		}
		fieldMap["wifi_mac"] = "00:23:89:b9:1a:" + randValue(100, "%02v")
	}
	/*
		同一设备操作赋值
	*/

	initParams("devices.txt")
	if SameDev {
		if isOldVersion && isApp { //兼容老平台
			//设备与采集信息
			// fieldMap["device_type"] = "Android"
			fieldMap["device_info"] = devfp[0]

			//位置信息
			fieldMap["province"] = params_mp["Province"]
			fieldMap["ip_province"] = params_mp["Province"]
			fieldMap["city"] = params_mp["City"]
			fieldMap["ip_city"] = params_mp["City"]
			fieldMap["longitude"] = params_mp["longitude"]
			fieldMap["latitude"] = params_mp["latitude"]
			// fieldMap["device_info"] = "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"
		} else {
			fieldMap["dev_fp"] = "9C613307C1AF33CC1BE626BEBE52944" + params_mp["ReqTimes"]
		}

	} else {
		if isOldVersion && isApp { //兼容老平台
			rn := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10)
			fieldMap["device_info"] = devfp[rn]
			//位置信息
			fieldMap["province"] = params_mp["Province"]
			fieldMap["ip_province"] = params_mp["Province"]
			fieldMap["city"] = params_mp["City"]
			fieldMap["ip_city"] = params_mp["City"]
			fieldMap["longitude"] = params_mp["longitude"]
			fieldMap["latitude"] = params_mp["latitude"]
			// fieldMap["device_info"] = "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"
		} else {
			fieldMap["dev_fp"] = randValue(1000000, "%06v") + "07C1AF33CC1BE626BEBE" + randValue(1000000, "%06v")
		}

	}

	//系统信息
	fieldMap["key_version"] = "01"
	fieldMap["service_type"] = "0"
	fieldMap["signature"] = ""
	if isOldVersion { //老版本字段
		fieldMap["time"] = fieldMap["trade_time"]
	}
	//线下渠道关键字段
	if isATM {
		if SameAtmNo {
			fieldMap["atm_no"] = "A-020-190012" + params_mp["ReqTimes"]
		} else {
			fieldMap["atm_no"] = "A-" + randValue(1000, "%03v") + "-" + randValue(1000000, "%06v")
		}
	}
	if isPOS {
		if SamePosNo {
			fieldMap["pos_no"] = "P-030-20190" + params_mp["ReqTimes"]
		} else {
			fieldMap["pos_no"] = "P-" + randValue(1000, "%03v") + "-" + randValue(1000000, "%06v")
		}
	}

	if isMerchant {
		if SameMechantNo {
			fieldMap["merchant_no"] = "M-010-11912" + params_mp["ReqTimes"]
		} else {
			fieldMap["merchant_no"] = "M=" + randValue(1000, "%03v") + "-" + randValue(1000000, "%06v")
		}
		if SameShopNo {
			fieldMap["shop_no"] = "S-130-12010" + params_mp["ReqTimes"]
		} else {
			fieldMap["shop_no"] = "S-" + randValue(1000, "%03v") + "-" + randValue(1000000, "%06v")
		}
	}

	//微信账号
	fieldMap["wechat_account"] = ""
}

func randValue(lens int32, fmtStr string) string {
	/*
		初始化随机值：获取lens标识位数的随机数字的字符串
	*/
	return fmt.Sprintf(fmtStr, rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(lens))
}

/*
读取配置文件，将变量->值提取到map中备用
*/
var devfp = make([]string, 20)

func initParams(filepath string) []string {
	/*
		读取配置文件，获取配置参数
	*/
	file, err := os.Open(filepath) //try to open file
	if err != nil {
		fmt.Printf("读取配置文件失败：%v\n", err)
		panic(err)
	}
	/*
		函数调用结束时关闭文件
	*/
	defer file.Close()
	/*创建文件读取对象 */
	r := bufio.NewReader(file)
	i := 0
	for {
		//按行读取，使用ReadString,解决readline读取长度过长导致字符串被截断问题
		b, err := r.ReadString('\n')
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
			/**参数配置的"="两边值，映射字段*/
			/*
				获取参数字段名称:k
			*/
			// k := strings.TrimSpace(strings.Split(str, "=")[0])
			/*
				获取对应的参数值：v
			*/
			// v := strings.TrimSpace(strings.Split(str, "=")[1])
			/*
				将参数赋值存储到map:ParamsMp
			*/
			devfp[i] = strings.TrimSpace(str)
			// fmt.Println(i, devfp[i])
			i++
		}
	}
	return devfp[:]
}

/*
读取配置文件，将变量->值提取到map中备用
*/
var Fields = make(map[string]string, 20)

func initFields(filepath string) map[string]string {
	/*
		读取配置文件，获取配置参数
	*/
	file, err := os.Open(filepath) //try to open file
	if err != nil {
		fmt.Printf("读取配置文件失败：%v\n", err)
		panic(err)
	}
	/*
		函数调用结束时关闭文件
	*/
	defer file.Close()
	/*创建文件读取对象 */
	r := bufio.NewReader(file)
	i := 0
	for {
		//按行读取，使用ReadString,解决readline读取长度过长导致字符串被截断问题
		b, err := r.ReadString('\n')

		//去掉已读取的行首位空格
		str := strings.TrimSpace(string(b))
		//如果以注释符"#"打头，不做处理
		if strings.Index(str, "#") == 0 {
			continue
		} else if len(str) != 0 {
			/**参数配置的"="两边值，映射字段*/
			/*
				获取参数字段名称:k
			*/
			k := strings.TrimSpace(strings.Split(str, "=")[0])
			/*
				获取对应的参数值：v
			*/
			v := strings.TrimSpace(strings.Split(str, "=")[1])
			/*
				将参数赋值存储到map:ParamsMp
			*/
			Fields[k] = v
			// fmt.Println(i, devfp[i])
			i++
			/*放到最后*/
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
		}
	}
	return Fields
}
