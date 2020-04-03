package main

import (
	"github.com/mamemomonga/ec2ctrl/configs"
	"github.com/mamemomonga/ec2ctrl/awsi"
	"github.com/mamemomonga/ec2ctrl/ec2ctrl/buildinfo"
	"github.com/mamemomonga/ec2ctrl/commands"
//	"github.com/davecgh/go-spew/spew"
	"os"
	"fmt"
	"log"
)

var cfg *configs.Configs

func usage() {
	fmt.Println("------------------------------------------------------")
	fmt.Printf("ec2ctrl %s-%s\n",buildinfo.Version, buildinfo.Revision)
	fmt.Println("https://github.com/mamemomonga/ec2ctrl")
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

func main() {
	cfg = configs.New()
	cfg.Load()
	err := runner()
	if err != nil {
		log.Fatal(err)
	}
}

func runner() error {
	if len(os.Args) == 1 {
		usage()
		return nil
	}
	args   := os.Args[1:]

	target := args[0]
	if ! cfg.SetTarget(target) {
		usage()
		return nil
	}

	shiftArgs := func() {
		args = args[1:]
		if len(args) <= 0 {
			usage()
			return
		}
	}
	ai := awsi.New(cfg)
	cmds := commands.New(cfg, ai)

	shiftArgs()
	switch args[0] {
	case "myip":
		shiftArgs()
		if !cfg.Target.Enables.MyIP {
			usage()
			return nil
		}
		switch args[0] {
		case "set":
			return ai.MyIPSet()
		case "del":
			return ai.MyIPDel()
		}
	case "instance":
		if !cfg.Target.Enables.Instance {
			usage()
			return nil
		}
		shiftArgs()
		switch args[0] {
		case "start":
			return ai.InstanceStart()
		case "stop":
			return ai.InstanceStop()
		case "status":
			return ai.InstanceStatus()
		}
	case "ssh":
		if !cfg.Target.Enables.SSH {
			usage()
			return nil
		}
		shiftArgs()
		switch args[0] {
		case "login":
			cmds.SSHLogin(true)
			return nil
		case "show":
			cmds.SSHLogin(false)
			return nil
		}
	case "rdp":
		if !cfg.Target.Enables.RDP {
			usage()
			return nil
		}
		cmds.RDPConnect()
		return nil
	}
	return nil
}

