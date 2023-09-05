package main

import (
	"fmt"
	"strconv"
)

func main1() {
	sp := " "
	//行
	for i := 1; i <= 9; i++ {
		//列
		fmt.Println()
		for j := 1; j <= i; j++ {
			//print(i * j)
			//if j <= i {
			if len(strconv.Itoa(j*i)) == 1 {
				fmt.Printf("%d%s%d%s%d%s%s", j, "x", i, "=", j*i, sp, sp)
			} else {
				fmt.Printf("%d%s%d%s%d%s", j, "x", i, "=", j*i, sp)
			}
			//}
		}
	}
}
