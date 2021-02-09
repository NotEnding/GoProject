package main

import (
	"log"
	"time"
)

var myMap = make(map[int64]int64)

// 计算阶乘
func factorial(n int64) int64 {
	if n == 0 {
		return 1
	}
	var i int64 = 1
	var res int64 = 1
	for ; i <= n; i++ {
		res *= i
	}
	return res
}

func putMap(n int64) {
	myMap[n] = factorial(n)
}

func main() {
	for i := 1; i <= 20; i++ {
		go putMap(int64(i))
	}
	time.Sleep(3 * time.Second)
	// 遍历map
	for i, v := range myMap {
		log.Println(i, v)
	}
}