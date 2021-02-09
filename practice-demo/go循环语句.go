package main

import "fmt"

func main() {
	for i := 1; i <= 9; i++ { // i 控制行，以及计算的最大值
		for j := 1; j <= i; j++ { // j 控制每行的计算个数
			fmt.Printf("%d*%d=%d ", j, i, j*i)
		}
		fmt.Println("")
	}
}

// 九九乘法表
