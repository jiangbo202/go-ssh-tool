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

var source string
var target string

// scpDownCmd represents the scpDown command
var scpDownCmd = &cobra.Command{
	Use:   "down",
	Short: "例如: down -m 1 -s /root/aaa.txt",
	Long:  `从指定主机上下载文件`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("从主机: %s, 下载(scp)文件: %s \n", cfg.Host[hostId-1].Ip, source)
		result, err := ScpDownCommand(source, target)
		if err != nil {
			log.Printf("遇到错误: %s\n", err)
			return
		}
		log.Printf("执行成功, 返回信息如下: \n")
		log.Printf(result)
	},
}

func init() {
	rootCmd.AddCommand(scpDownCmd)

	scpDownCmd.Flags().StringVarP(&source, "source", "s", "", "想要下载的文件名(包括路径), eg: /tmp/111.txt")
	scpDownCmd.Flags().StringVarP(&target, "target", "t", "", "下载到目录或者重命名为, ./222.txt")
	scpDownCmd.Flags().IntVarP(&hostId, "host", "m", 0, "主机id, eg: 1")
	scpDownCmd.MarkFlagRequired("host")
	scpDownCmd.MarkFlagRequired("source")
	//scpDownCmd.MarkFlagRequired("target")  可以不传
}

func ScpDownCommand(source, target string) (string, error) {
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
	result, err := client.Scp(source, target)
	if err != nil {
		log.Printf("failed to run shell,err=[%v]\n", err)
		return "", err
	}
	return fmt.Sprintf("拷贝的大小: %d", result), nil
}
