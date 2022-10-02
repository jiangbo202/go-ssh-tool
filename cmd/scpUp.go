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

var file string
var dir string

// scpUpCmd represents the scpUp command
var scpUpCmd = &cobra.Command{
	Use:     "up",
	Example: "例如: up -m 1 -f ./config.yaml -d /tmp",
	Short:   "例如: up -m 1 -f ./config.yaml -d /tmp",
	Long:    `上传 by scp`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("上传(scp)文件: %s 到主机: %s 的目录: %s,  \n", source, cfg.Host[hostId-1].Ip, dir)
		result, err := ScpUpCommand(file, dir)
		if err != nil {
			log.Printf("遇到错误: %s\n", err)
			return
		}
		log.Printf("执行成功, 返回信息如下: \n")
		log.Printf(result)
	},
}

func init() {
	rootCmd.AddCommand(scpUpCmd)

	scpUpCmd.Flags().StringVarP(&file, "file", "f", "", "想要上传的文件名(包括路径), eg: /tmp/111.txt")
	scpUpCmd.Flags().StringVarP(&dir, "dir", "d", "", "上传到哪个目录, eg: /tmp/")
	scpUpCmd.Flags().IntVarP(&hostId, "host", "m", 0, "主机id, eg: 1")
	scpUpCmd.MarkFlagRequired("file")
	scpUpCmd.MarkFlagRequired("host")
}

func ScpUpCommand(source, target string) (string, error) {
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
	result, err := client.Scp2(source, target)
	if err != nil {
		log.Printf("failed to run shell,err=[%v]\n", err)
		return "", err
	}
	return fmt.Sprintf("拷贝的大小: %d", result), nil
}
