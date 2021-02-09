package main

import (
	"fmt"
	"strings"
)

func main()  {
	//对于数字类型，无需定义int 及float32、float64,系统会自动识别
	var a = 1.5
	var b = 2
	fmt.Println(a,b)

	str := "一切都是刚刚好 www\n.runnoob\n.com"  //简短声明
	fmt.Println("-------原字符串----------")
	fmt.Println(str)

	//去除空格
	str = strings.Replace(str," ","",-1)
	//去除换行符
	str = strings.Replace(str,"\n","",-1)
	fmt.Println("-------去除空格和换行后——————————")
	fmt.Println(str)

	// 声明一个变量并初始化赋值
	var a1 = "RUNOOB"
	fmt.Println(a1)

	//没有初始化，仅定义 就为零值
	var b2 int
	fmt.Println(b2)

	// bool 零值为false
	var c3 bool
	fmt.Println(c3)
}
