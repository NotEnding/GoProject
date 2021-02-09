package main

import (
	"fmt"
)

func main() {
	var a int = 21
	var b int = 10

	if a == b {
		fmt.Println("第一行：a 等于 b")
	} else {
		fmt.Println("第一行：a 不等于 b")
	}

	if a < b {
		fmt.Println("第二行：a 小于 b")
	} else {
		fmt.Println("第二行：a大于b")
	}

	if a <= b {
		fmt.Println("第三行：a 小于等于 b")
	}

	if a >= b {
		fmt.Println("第四行：a 大于等于 b")
	}

}
