/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"myWorkTool/utils"
)

var (
	command string
)

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:     "exe",
	Example: "exe -m 1 -c \"cat /etc/redhat-release\"",
	//Args:    cobra.ExactArgs(1),  这个是直接参数
	Short: "ssh指定机器执行一个命令，并返回",
	Long:  `ssh指定机器执行一个命令，并返回`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("对主机: %s, 执行命令:%s \n", cfg.Host[hostId-1].Ip, command)
		result, err := ExeCommand(hostId, command)
		if err != nil {
			log.Printf("遇到错误: %s\n", err)
			return
		}
		log.Printf("执行成功, 返回信息如下: \n")
		log.Printf(result)
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)

	executeCmd.Flags().StringVarP(&command, "command", "c", "echo ok", "想要执行的命令, eg: \"cat /etc/redhat-release\"")
	executeCmd.PersistentFlags().IntVarP(&hostId, "host", "m", 0, "主机id, eg: 1")
	executeCmd.MarkFlagRequired("host")
	executeCmd.MarkFlagRequired("command")
}

func ExeCommand(hostId int, command string) (string, error) {
	index := hostId - 1
	var (
		username     = cfg.Host[index].Auth.Username
		password     = cfg.Host[index].Auth.Passwd
		addr         = fmt.Sprintf("%s:%d", cfg.Host[index].Ip, cfg.Host[index].Port)
		googleSecret = cfg.Host[index].Auth.GoogleSecret
	)
	// 初始化
	client := utils.NewCli(username, password, addr, googleSecret)

	// ssh 并执行命令
	result, err := client.Run(command)
	if err != nil {
		log.Printf("failed to run shell,err=[%v]\n", err)
		return "", err
	}
	return result, nil
}
