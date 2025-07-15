package menu

import (
	"fmt"
	"windowsFinder/config"
	"windowsFinder/function"
)

func init() {
	config.Banner()
}
func list() {
	fmt.Println(" Windows finder menu ")
	fmt.Println("1. 进行一些基础信息收集")
}

func Menu() {
	for {
		list()
		var choice int
		fmt.Scanf("%d", &choice)
		switch choice {
		case 1:
			function.BasicCollect()
		default:
			return
		}
	}
}
