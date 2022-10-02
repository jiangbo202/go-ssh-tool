/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	cmd2 "myWorkTool/cmd"
	"myWorkTool/utils"
)

func main() {

	hostList := utils.InitConfig("config.yaml")
	cmd2.Execute(hostList)

}
