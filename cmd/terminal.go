/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"myWorkTool/utils"

	"github.com/spf13/cobra"
)

// terminalCmd represents the terminal command
var terminalCmd = &cobra.Command{
	Use:   "term",
	Short: "go run main.go term -m 1",
	Long:  `登录主机by ssh`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("将要登录到主机: %s...\n", cfg.Host[hostId-1].Ip)
		result, err := loginBySSH()
		if err != nil {
			log.Printf("遇到错误: %s\n", err)
			return
		}
		log.Printf(result)
	},
}

func init() {
	rootCmd.AddCommand(terminalCmd)

	terminalCmd.Flags().IntVarP(&hostId, "host", "m", 0, "主机id, eg: 1")
	terminalCmd.MarkFlagRequired("host")
}

func loginBySSH() (string, error) {
	index := hostId - 1
	var (
		username     = cfg.Host[index].Auth.Username
		password     = cfg.Host[index].Auth.Passwd
		addr         = fmt.Sprintf("%s:%d", cfg.Host[index].Ip, cfg.Host[index].Port)
		googleSecret = cfg.Host[index].Auth.GoogleSecret
	)
	// 初始化
	client := utils.NewCli(username, password, addr, googleSecret)

	err := client.NewTerminal()
	if err != nil {
		log.Printf("failed to run shell,err=[%v]\n", err)
		return "", err
	}
	return "退出", nil
}
