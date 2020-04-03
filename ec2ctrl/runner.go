package main

import (
	"github.com/mamemomonga/ec2ctrl/awsi"
	"github.com/mamemomonga/ec2ctrl/commands"
//	"github.com/davecgh/go-spew/spew"
	"os"
)

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

