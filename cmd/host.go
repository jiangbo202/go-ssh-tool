/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
    Use:   "host",
    Short: "列出所有主机",
    Long:  `列出所有主机`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("主机列表如下:")
        for index, host := range cfg.Host {
            fmt.Printf("序号: %d, 主机: %s\n", index+1, host.Ip)
        }
    },
}

func init() {
    rootCmd.AddCommand(hostCmd)
}
