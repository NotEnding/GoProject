package main

import "fmt"

func main() {
	var a int = 20
	var ip *int /*声明指针变量*/

	ip = &a /* 指针变量的存储地址 */

	fmt.Printf("a 变量的地址是：%x\n", &a)

	/* 指针变量的存储地址 */
	fmt.Printf("ip 变量储存的指针地址：%x\n", ip)

	/* 指针自身的地址 */
	fmt.Printf("ip 指针自身的地址：%x\n", &ip)

	/* 使用指针访问值 */
	fmt.Printf("*ip 变量的值:%d\n", *ip)

	var ptr *int
	fmt.Printf("空指针的值为：%x\n", ptr)
	if ptr != nil {
		fmt.Printf("ptr不是空指针")
	} else if ptr == nil {
		fmt.Printf("ptr是空指针")
	}

}
