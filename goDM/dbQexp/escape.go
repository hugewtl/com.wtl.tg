package main

import (
	"fmt"
	"strings"
)

/*
针对特殊字符进行转义：将映射关系存入map[string]string
*/
var escapeMp = make(map[string]string, 15)

/*
转义字符串中特殊字符
! : %21
@ : %40
# : %23
$ : %24
% : %25
^ : %5e
& : %26
* : %2a
( : %28
) : %29
_ : %5f
+ : %2b
= : %3d
*/
func escapeForSC(str string) string {
	/*
		% 比较特殊，放在第一位，如果原字符串包含%,则进行转义
	*/
	specialChar := []string{"%", "!", "@", "#", "$", "^", "&", "*", "(", ")", "_", "+", "="}
	if strings.ContainsAny(str, "%!@#$^&*()_+=") {
		for _, Sc := range specialChar {
			str = strings.ReplaceAll(str, Sc, "%"+fmt.Sprintf("%x", Sc))
			// log.Println(Sc)
		}
	}
	return str
}

// func main() {
// 	log.Println(escapeForSC("Cq!(myg#1#&2+3@5"))
// }
