package main

import "wxSendMessage/defs/utils"

func init() {
	// 初始化配置
	utils.ReadConfigurationFile()
}

func main() {
	utils.SendMessage()
}
