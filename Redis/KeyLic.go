package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func ifDateValid() bool {
	flag := false
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
	/*
		当前时间戳与当月最后一天时间戳比较，判断是否过期
	*/
	log.Printf("当前时间戳：%v \n", time.Now().Unix())
	unixtime, err := strconv.Atoi(shaToStr)
	if err != nil {
		fmt.Println("lic无效")
		panic(err)
	}
	if time.Now().Unix() <= int64(unixtime) {
		flag = true
	}
	// if flag {
	// 	fmt.Println("Ok")
	// }
	return flag
}

// func main() {
// 	// base64KeyGenerate()

// 	ifDateValid()

// 	// // 生成密钥对，保存到文件
// 	// GenerateRSAKey(2048)
// 	// dateValid := time.Now().AddDate(0, 0, 90).Format("200601021504")
// 	// message := []byte(dateValid)
// 	// // message1 := []byte("202412112")
// 	// cipherText := RSA_Encrypt(message, "public.pem")
// 	// fmt.Println("加密后为：", string(cipherText))
// 	// // 解密
// 	// // dateNow := time.Now().Format("200601021504")
// 	// // cipherText := []byte("202312211612")
// 	// plainText := RSA_Decrypt(cipherText, "private.pem")
// 	// fmt.Println("解密后为：", string(plainText))
// }

// func base64KeyGenerate() {
// 	//获取一年后的时间戳：1726910232
// 	dateValidRange := time.Now().AddDate(0, 0, 365).Unix()
// 	// fmt.Println(GetLastDateOfMonth(time.Now().AddDate(0, 0, 365)).Unix())
// 	fmt.Println(dateValidRange)
// 	// /*
// 	//    base64加密,永久秘钥
// 	// */
// 	// strToSha0 := func(src string) string {
// 	// 	return string(base64.StdEncoding.EncodeToString([]byte(src)))
// 	// }(fmt.Sprintf("%v", GetLastDateOfMonth(time.Now()).Unix()+1648656000))
// 	// // //加密LIC
// 	// fmt.Printf("%s: %v TDL: %v ", "lic", strToSha0, fmt.Sprintf("%v", GetLastDateOfMonth(time.Now()).Unix()+1648656000))
// 	/*
// 	   base64加密,一年有效期密钥
// 	*/
// 	strToSha1 := func(src string) string {
// 		return string(base64.StdEncoding.EncodeToString([]byte(src)))
// 	}(fmt.Sprintf("%v", dateValidRange))
// 	//加密LIC
// 	fmt.Printf("%s: %v TDL: %v  当前时间戳：", "lic", strToSha1, fmt.Sprintf("%v", dateValidRange))

// 	// /*
// 	//    base64加密,当月秘钥
// 	// */
// 	// strToSha1 := func(src string) string {
// 	// 	return string(base64.StdEncoding.EncodeToString([]byte(src)))
// 	// }(fmt.Sprintf("%v", GetLastDateOfMonth(time.Now()).Unix()))
// 	// //加密LIC
// 	// fmt.Printf("%s: %v TDL: %v  当前时间戳：", "lic", strToSha1, fmt.Sprintf("%v", GetLastDateOfMonth(time.Now()).Unix()))

// }

// // 生成RSA私钥和公钥，保存到文件中
// func GenerateRSAKey(bits int) {
// 	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
// 	//Reader是一个全局、共享的密码用强随机数生成器
// 	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//保存私钥
// 	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
// 	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
// 	//使用pem格式对x509输出的内容进行编码
// 	//创建文件保存私钥
// 	privateFile, err := os.Create("private.pem")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer privateFile.Close()
// 	//构建一个pem.Block结构体对象
// 	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
// 	//将数据保存到文件
// 	pem.Encode(privateFile, &privateBlock)
// 	//保存公钥
// 	//获取公钥的数据
// 	publicKey := privateKey.PublicKey
// 	//X509对公钥编码
// 	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//pem格式编码
// 	//创建用于保存公钥的文件
// 	publicFile, err := os.Create("public.pem")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer publicFile.Close()
// 	//创建一个pem.Block结构体对象
// 	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
// 	//保存到文件
// 	pem.Encode(publicFile, &publicBlock)
// }

// // RSA加密
// func RSA_Encrypt(plainText []byte, path string) []byte {
// 	//打开文件
// 	file, err := os.Open(path)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()
// 	//读取文件的内容
// 	info, _ := file.Stat()
// 	buf := make([]byte, info.Size())
// 	file.Read(buf)
// 	//pem解码
// 	block, _ := pem.Decode(buf)
// 	//x509解码
// 	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//类型断言
// 	publicKey := publicKeyInterface.(*rsa.PublicKey)
// 	//对明文进行加密
// 	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//返回密文
// 	return cipherText
// }

// // RSA解密
// func RSA_Decrypt(cipherText []byte, path string) []byte {
// 	//打开文件
// 	file, err := os.Open(path)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()
// 	//获取文件内容
// 	info, _ := file.Stat()
// 	buf := make([]byte, info.Size())
// 	file.Read(buf)
// 	//pem解码
// 	block, _ := pem.Decode(buf)
// 	//X509解码
// 	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//对密文进行解密
// 	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
// 	//返回明文
// 	return plainText
// }

// 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}
