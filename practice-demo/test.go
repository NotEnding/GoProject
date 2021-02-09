package main

import "fmt"

var Total_sum int = 0

func Sum_test(a int,b int) int {
	fmt.Printf("%d + %d = %d\n",a,b,a+b)
	Total_sum += (a+b)  // 0 += 2 + 3
	fmt.Printf("Total_sum:%d\n",Total_sum)
	return a + b
}
// 变量作用域不同


func main()  {
	var sum int
	sum = Sum_test(2,3)
	fmt.Printf("sum:%d;Total_sum:%d\n",sum,Total_sum)
}

