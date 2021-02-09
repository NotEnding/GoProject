//package main
//
//import "fmt"
//
//func main() {
//	if num := 10; num%2 == 0 { // 判断数字是否为偶数
//		fmt.Println(num, "是偶数")
//	} else {
//		fmt.Println(num, "是奇数")
//	}
//}

/*
if statement; condition {
}
*/


package main

import "fmt"
func main() {
	if num := 9; num < 0 {
		fmt.Println(num, "is negative")
	} else if num < 10 {
		fmt.Println(num, "has 1 digit")
	} else {
		fmt.Println(num, "has multiple digits")
	}
}
