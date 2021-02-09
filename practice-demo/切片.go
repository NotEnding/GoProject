package main

import "fmt"

func main() {
	//var numbers = make([]int, 3, 5)

	var numbers []int  //空切片
	printSlice(numbers)

	if numbers == nil{
		fmt.Printf("切片是空的！！！")
	}

}

//输出数组info
func printSlice(x []int) {
	fmt.Printf("len=%d cap=%d slice=%v\n", len(x), cap(x), x)
}
