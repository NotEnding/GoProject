package main

import (
	"fmt"
	"./practice-demo/trans" //导入自定义的包
)

var twoPi = 2 * trans.Pi   // 调用trans下 定义的变量

func main()  {
	fmt.Printf("2*Pi= %g\n",twoPi)
}