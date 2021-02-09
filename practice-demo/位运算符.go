package main

import "fmt"

func main()  {
	var a int = 60   /* 60 = 0011 1100 */
	var b int =  13  /* 13 = 0000 1101 */
	var c int = 0

	c = a & b   /* 12 = 0000 1100 */
	fmt.Printf("第一行 c 的值为：%d \n",c)

	c = a | b   /* 61 = 0011 1101 */
	fmt.Printf("第二行 c 的值为：%d \n",c)

	c = a ^ b   /* 49 = 0011 0001  对应的二进制位相异为1，相同为0 */
	fmt.Printf("第三行 c 的值为：%d \n",c)

	c = a << 2  /* 240 = 1111 0000 */
	fmt.Printf("第四行 c 的值为：%d \n",c)

	c = b >> 2  /* 3 = 0000 0011 */
	fmt.Printf("第五行 c 的值为：%d \n",c)
}