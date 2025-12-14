package main

import (
	"fmt"
	"gin_mall/conf"
)

func main() {
	fmt.Println("Hello World")
	conf.Init()
	fmt.Println("Hello World")
}
