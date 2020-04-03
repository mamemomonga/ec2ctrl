package main

import (
	"github.com/mamemomonga/ec2ctrl/ec2ctrl/buildinfo"
//	"github.com/davecgh/go-spew/spew"
	"os"
	"fmt"
)

func usage() {
	fmt.Println("------------------------------------------------------")
	fmt.Printf("  ec2ctrl %s-%s\n",buildinfo.Version, buildinfo.Revision)
	fmt.Println("  https://github.com/mamemomonga/ec2ctrl/")
	fmt.Println("------------------------------------------------------")
	fmt.Printf("USAGE: %s [COMMANDS]\n", os.Args[0])
	fmt.Printf("COMMANDS:\n")
	for _,i := range cfg.Configs.Targets {
		fmt.Println()
		fmt.Printf(" %s\n",i.Description)
		if i.Enables.MyIP {
			fmt.Printf("  %s myip set        自分のIPアドレスを許可設定\n" ,i.Name)
			fmt.Printf("  %s myip del        自分のIPアドレスを許可削除\n" ,i.Name)
		}
		if i.Enables.Instance {
			fmt.Printf("  %s instance status 状態を表示\n"  ,i.Name)
			fmt.Printf("  %s instance start  起動\n" ,i.Name)
			fmt.Printf("  %s instance stop   停止\n"  ,i.Name)
		}
		if i.Enables.SSH {
			fmt.Printf("  %s ssh login       SSHログイン\n"  ,i.Name)
			fmt.Printf("  %s ssh show        SSHコマンドの表示\n"  ,i.Name)
		}
		if i.Enables.RDP {
			fmt.Printf("  %s rdp             RDPトンネル開始と接続情報の表示\n"  ,i.Name)
		}
	}
	fmt.Println()
	os.Exit(1)
}
