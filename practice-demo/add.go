//当前程序的包名
package main

// 为 fmt 起别名为 fmt2
import fmt2 "fmt"

/* 省略调用（不建议使用）
调用的时候只需要 println(),不需要 fmt.println()
*/
import . "fmt"

// 常量定义
const PI = 3.14

// 全局变量的声明和赋值
var name = "go developer"

// 一般类型声明
type newType int

// 结构的声明
type gopher struct {

}


// 接口的声明
type golang interface {

}

// 由main函数作为程序的入口点启动

func main(){
	println("code change the world")
}

// 函数名首字母大写 即为public方法