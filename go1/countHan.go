package main

import (
	"fmt"
	"unicode"
)

func main() {
	//统计出一个字符串中，汉字的数量
	s1 := "12312字符串中的汉字huge"
	s2 := []rune(s1)
	l := len(s2)
	count := 0
	//方法1：使用range遍历
	for k, v := range s2 {
		//unicode.Is()
		if unicode.Is(unicode.Han, v) {
			fmt.Printf("(%d,%c)\n", k, v)
			count++
		}
	}
	//使用哑元变量去掉索引key
	for _, v := range s2 {
		//unicode.Is()
		if unicode.Is(unicode.Han, v) {
			fmt.Printf("%c\n", v)
			count++
		}
	}

	//方法2：使用数组方法遍历
	for i := 0; i < l; i++ {
		//unicode.Is()
		if unicode.Is(unicode.Han, s2[i]) {
			fmt.Printf("%c\n", s2[i])
			count++
		}
	}
	fmt.Print(count)
	//fmt.Printf("%c\n", s2)
	//fmt.Printf("%v\n", s1)
	//fmt.Printf("%v\n", s2)
}
