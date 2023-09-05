package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// 将参数存储到map中
var ParamsMp = make(map[string]string, 15)

/*
读取配置文件，将变量->值提取到map中备用
*/
func initParams(filepath string) map[string]string {
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
			ParamsMp[k] = v
		}
	}
	return ParamsMp
}

/*
传入目标参数，遍历参数文件查到参数值
*/
// func readParams(param string) string {
// 	/*
// 		读取配置文件，获取配置参数
// 	*/
// 	file, err := os.Open("conf.properties") //try to open file
// 	if err != nil {
// 		fmt.Printf("读取配置文件失败：%v\n", err)
// 		panic(err)
// 	}
// 	/*
// 		函数调用结束时关闭文件
// 	*/
// 	defer file.Close()
// 	/*创建文件读取对象 */
// 	r := bufio.NewReader(file)
// 	for {
// 		//按行读取
// 		b, _, err := r.ReadLine()
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			panic(err)
// 		}
// 		//去掉已读取的行首位空格
// 		str := strings.TrimSpace(string(b))
// 		//如果以注释符"#"打头，不做处理
// 		if strings.Index(str, "#") == 0 {
// 			continue
// 		} else if len(str) != 0 {
// 			//参数配置的"="两边值，映射字段
// 			/*
// 				获取参数字段名称:k
// 			*/
// 			k := strings.TrimSpace(strings.Split(str, "=")[0])
// 			/*
// 				获取对应的参数值：v
// 			*/
// 			v := strings.TrimSpace(strings.Split(str, "=")[1])
// 			/*
// 				将参数赋值存储到map:params_mp
// 			*/
// 			if k == param {
// 				return v
// 			}
// 		}
// 	}

// 	return os.DevNull
// }
