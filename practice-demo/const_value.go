package main

import "unsafe"

const (
	a = "nothing"
	b = len(a)
	c = unsafe.Sizeof(a)
)

func main()  {
	println(a,b,c)
}
