package main
import "fmt"

/*
结构体定义需要使用 type 和 struct 语句。struct 语句定义一个新的数据类型，结构体中有一个或多个成员。type 语句设定了结构体的名称。结构体的格式如下

type struct_variable_type struct {
   member definition
   member definition
   ...
   member definition
}

一旦定义了结构体类型，它就能用于变量的声明，语法格式如下
variable_name := structure_variable_type {value1, value2...valuen}
或
variable_name := structure_variable_type { key1: value1, key2: value2..., keyn: valuen}

 */

type Books struct {
	title string
	author string
	subject string
	book_id int
}

func main()  {
	//创建一个新的结构体
	fmt.Println(Books{"Go 语言", "www.runoob.com", "Go 语言教程", 6495407})

	// 也可以使用 key => value 格式
	fmt.Println(Books{title: "Go 语言", author: "www.runoob.com", subject: "Go 语言教程", book_id: 6495407})

	// 忽略的字段为 0 或 空
	fmt.Println(Books{title: "Go 语言", author: "www.runoob.com"})

	var book1 Books
	var book2 Books

	/* book 1 描述 */
	book1.title = "Go 语言"
	book1.author = "www.runoob.com"
	book1.subject = "Go 语言教程"
	book1.book_id = 6495407

	/* book 2 描述 */
	book2.title = "Python 教程"
	book2.author = "www.runoob.com"
	book2.subject = "Python 语言教程"
	book2.book_id = 6495700


	/* 打印 Book1 信息 */
	fmt.Printf( "Book 1 title : %s\n", book1.title)
	fmt.Printf( "Book 1 author : %s\n", book1.author)
	fmt.Printf( "Book 1 subject : %s\n", book1.subject)
	fmt.Printf( "Book 1 book_id : %d\n", book1.book_id)

	/* 打印 Book2 信息 */
	fmt.Printf( "Book 2 title : %s\n", book2.title)
	fmt.Printf( "Book 2 author : %s\n", book2.author)
	fmt.Printf( "Book 2 subject : %s\n", book2.subject)
	fmt.Printf( "Book 2 book_id : %d\n", book2.book_id)
}