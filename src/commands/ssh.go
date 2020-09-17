package commands

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"strings"
)

type sshProxyOptOutput struct {
	option     string
	connection string
}


func (t *Commands) sshProxyOpt() sshProxyOptOutput {
	is := t.ai.InstanceState()
	if t.configs.Target.SSH.Direct {
		return sshProxyOptOutput{
			connection: fmt.Sprintf("%s@%s",
				t.configs.Target.SSH.Username,
				is.PublicIpAddress),
		}
	} else {
		return sshProxyOptOutput{
			option: fmt.Sprintf("ProxyCommand ssh %s@%s -W %%h:%%p 2> /dev/null",
				t.configs.Target.Bastion.Username,
				t.configs.Target.Bastion.Host),
			connection: fmt.Sprintf("%s@%s",
				t.configs.Target.SSH.Username,
				is.PrivateIpAddress),
		}
	}
}

func (t *Commands) sshCommand(args []string) []string {
	opts := []string{}
	cmds := []string{}

	// ハイフンをコマンドに置換
	{
		f := false
		for _,v := range args {
			if v == "-" {
				f = true
				continue
			}
			if f {
				cmds = append(cmds,v)
			} else {
				opts = append(opts,v)
			}
		}
	}

	ag := []string{}
	sl := t.sshProxyOpt()
	if sl.option != "" {
		ag = append(ag, "-o", sl.option)
	}
	ag = append(ag, opts...)
	ag = append(ag, sl.connection)
	ag = append(ag, cmds...)
	return ag
}


func (t *Commands) SSHLogin(args []string) {
	ag := t.sshCommand(args)
	log.Printf("debug: args %s", spew.Sdump(ag))
	t.runCommand("ssh", ag...)
}

func (t *Commands) SSHLoginCmdShow(args []string) {
	ag := t.sshCommand(args)
	an := []string{}

	for _,s := range ag {
		if strings.Index(s," ") > 0 {
			an = append(an,`'`+s+`'`)
		} else {
			an = append(an,s)
		}
	}

	fmt.Printf("ssh %s\n", strings.Join(an," "))
}

