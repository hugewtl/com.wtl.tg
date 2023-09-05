package readFile

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var (
	// fieldName string
	rt       string
	fieldMap = map[string]string{}
	// mString   string
)

//读取字段文件：字段源数据
func readFile() {
	file, err := os.Open("./column")
	count := 0

	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n') //注意是字符
		if err == io.EOF {
			if len(line) != 0 {
				//fmt.Println(line)
				count++
			}
			fmt.Println("字段列表加载完成")
			break
		}
		if err != nil {
			fmt.Println("read file failed, err:", err)
			return
		}
		count++
		//log.Print(rand.Intn(3))
		//给字段赋值：字段+随机数字
		fieldMap[line] = randValue(line)
		//fmt.Print(line)
	}
	//字段总数
	fmt.Println(count)
	//fmt.Print(fieldMap)
}

//对每个字段生成随机值
func randValue(fieldName string) string {
	num := 0
	rand.Seed(time.Now().UnixNano())
	num = rand.Intn(100000)
	rt = strconv.Itoa(num) + "-" + fieldName
	//fmt.Print(rt)
	return rt
}
