package main

import (
	"fmt"
	"gin_mall/conf"
	"gin_mall/routes"
)

func main() {
	fmt.Println("Start")
	conf.Init()
	r := routes.NewRouter()
	_ = r.Run(conf.HttpPort)
	fmt.Println("end")
}
