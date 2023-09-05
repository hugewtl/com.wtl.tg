package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func main1() {
	if LicIsExp() {
		fmt.Println("没过期")
	}
}

/*
判断license是否过期
*/
func LicIsExp() bool {
	var flag bool
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
					return "解密错误"
				}
				strToSha = string(a)
			}
		}
		return strToSha
	}(strToSha) //base64匿名函数解码结束

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
	return flag
}

// 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}
