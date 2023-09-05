package main

/*
Author: zhaohu
Date:   2021-09-16 15:36:12
REVISE: 2022-04-14 18:12:23
REVISE: 2023-08-16 13:57:23
DESC:   for api rules test more auto
VERSION: v1.1.1
*/
import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gojsonq-2"
)

var (
	//字段数量
	fieldNum int = 0
	//定义报文字段存储数据结构,数组要做初始化
	fieldData = make([]string, 50)
	// //字段分割符
	separator rune
	//时间日期所在列
	isDate = make([]string, 10)
	//补充固定值字段
	ifFieldsFixed bool
	//配置文件string类型参数map,初始化map
	params_mpcsv = make(map[string]string, 25)
	// 	//定义channel,接收map：params_mp
	ch1 = make(chan []string, 1000)
	wg  sync.WaitGroup
)

/*
从csv文件读取每行数据，返回[]string
*/
func readLineFromFileToJson(datafileName string, dataFieldPath string) {
	/*
		拷贝数组，留存；同时获取数组元素个数，这样可以随时访问fieldData1
	*/
	fieldData1 := getFieldsData(dataFieldPath)
	//准备读取文件
	fs, err := os.Open(datafileName)
	if err != nil {
		log.Fatalf("can not open the file, err is %+v", err)
	}
	defer fs.Close()
	r := csv.NewReader(fs)
	//针对大文件，一行一行的读取文件
	i := 0
	for {
		row, err := r.Read()
		r.Comma = separator
		/*
			golang 的csv reader 会检查每一行的字段数量，如果不等于 FieldsPerRecord 就会抛出该错误
			√ 为负值时，不检查
			√ 0，这是默认值。检查第一行的字段数量，然后赋值给 FieldsPerRecord
			√ 正值
		*/
		r.FieldsPerRecord = -1
		if err != nil && err != io.EOF {
			log.Fatalf("can not read, err is %+v", err)
		}
		if err == io.EOF {
			break
		}
		//对表头+数据映射到map中；让数据个数匹配字段
		for i, k := range row {
			if i == 3 {
				params_mpcsv[fieldData1[i]] = strconv.FormatInt(CsvPreTreatment(params_mp["dateFormt"], k), 10)
			} else {
				params_mpcsv[fieldData1[i]] = k

			}
		}
		fmt.Printf("逐行读取第 %d 行数据：\n", i+1)
		reqSendCsv(params_mpcsv, apiUrl)
		// in_ch <- params_mp
		i++
	}

}

/*
从csv文件读取每行数据，返回[]string
*/
func readAllFromFileToJson(datafileName string, ch chan []string) bool {
	/*
		拷贝数组，留存；同时获取数组元素个数，这样可以随时访问fieldData1
	*/

	// ch = make(chan []string, 1000)
	//准备读取文件
	fs, err := os.Open(datafileName)
	if err != nil {
		log.Fatalf("can not open the file, err is %+v", err)
	}
	defer fs.Close()
	r := csv.NewReader(fs)
	//针对大文件，一行一行的读取文件
	/*
		定义reader的字段分隔符
	*/
	r.Comma = separator
	defer close(ch)

	for {
		row, err := r.Read()

		/*
			golang 的csv reader 会检查每一行的字段数量，如果不等于 FieldsPerRecord 就会抛出该错误
			√ 为负值时，不检查
			√ 0，这是默认值。检查第一行的字段数量，然后赋值给 FieldsPerRecord
			√ 正值
		*/
		r.FieldsPerRecord = -1
		if err != nil && err != io.EOF {
			log.Fatalf("can not read, err is %+v", err)
		}
		if err == io.EOF {
			break
		}
		ch <- row
	}
	return true
}

func getRowToApi(dataFieldGDPath string, dataFieldPath string, ch chan []string) bool {
	fieldData1 := getFieldsData(dataFieldPath)
	fieldData2 := initFields(dataFieldGDPath)
	i := 0
	//对表头+数据映射到map中；让数据个数匹配字段
	for {
		for i, k := range <-ch {
			/*
				数据中时间字段需要特殊转换，需要预先知道时间字段是哪列，如何优化？
			*/

			if ifDate(i, isDate) { //时间日期字段
				params_mpcsv[fieldData1[i]] = strconv.FormatInt(CsvPreTreatment(params_mp["dateFormt"], k), 10)
			} else {
				params_mpcsv[fieldData1[i]] = k
			}
		}
		/*
			补充固定值类数据
		*/
		if ifFieldsFixed {
			/*
				获取固定值字段
			*/
			for k, v := range fieldData2 {
				params_mpcsv[k] = v
			}
		}
		fmt.Printf("请求处理第 %d 行数据：\n", i+1)
		reqSendCsv(params_mpcsv, apiUrl)
		// in_ch <- params_mp
		i++
		if len(ch) == 0 {
			break
		}
	}

	return true
}

// 判断数字是否在字符串数组中
func ifDate(i int, dateStr []string) bool {
	for _, v := range dateStr {
		if strconv.Itoa(i) == v {
			return true
		}
	}
	return false
}

// 请求服务,发送json请求
func reqSendCsv(out_ch map[string]string, url string) {
	/*
		发送http请求api
	*/
	mjson, _ := json.Marshal(out_ch)
	mString := string(mjson)
	/*
		输出请求报文json
	*/
	// ifDebug := true
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
	defer resp.Body.Close()
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
		/*
			记录命中规则到日志文件
		*/
		if ret != nil {
			logging_in1(ret)
		}

		//计数
		// count++
	} else {
		if ifDebug {
			fmt.Println("事件上报返回报文：" + string(body))
		}
	}

}

// 将配置文件中的字段读取到数组标记:动态生成文件表头字段
func getFieldsData(dataFiledspath string) []string {
	/*
		读取配置文件，获取配置参数
	*/
	file, err := os.Open(dataFiledspath) //try to open file
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
			//将字段标记对应的列值做类型转换为数组下标
			i, err := strconv.Atoi(v)

			if err != nil {
				panic(err)
			}
			fieldData[i] = k
			//全局变量统计字段数量
			fieldNum++
			// fmt.Println(i, fieldData[i])
			/*放到最后*/
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
		}
	}
	return fieldData
}

// 对csv数据trade_time预处理，转换为时间戳
func CsvPreTreatment(dateformat string, datetime string) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")             //设置时区
	tt, _ := time.ParseInLocation(dateformat, datetime, loc) //2006-01-02 15:04:05是转换的格式
	return tt.UnixNano() / 1e6
	// strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
}
