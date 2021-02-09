package main

import "fmt"

func main() {
	_, b_val, c_val := numbers()
	fmt.Println(b_val, c_val)
	//舍弃 a
}

//定义一个返回数值的函数
func numbers() (int, int, string) {
	a, b, c := 1, 2, "hello world"
	return a, b, c
}