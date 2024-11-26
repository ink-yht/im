package main

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {

	initViperV1()

	server := InitWebServer()

	err := server.Run(":8080")
	if err != nil {
		return
	}
}

func initViperV1() {
	// 直接指定文件路径
	viper.SetConfigFile("config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
