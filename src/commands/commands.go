package commands

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/mamemomonga/ec2ctrl/src/awsi"
	"github.com/mamemomonga/ec2ctrl/src/configs"
	"log"
	"os"
	"os/exec"
)

type Commands struct {
	ai      *awsi.AWSi
	configs *configs.Configs
}

func New(configs *configs.Configs, ai *awsi.AWSi) *Commands {
	t := new(Commands)
	t.ai = ai
	t.configs = configs
	return t
}

func (t *Commands) runCommand(c string, p ...string) error {
	cmd := exec.Command(c, p...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

type SSHLoginCmdS struct {
	option     string
	connection string
}

func (t *Commands) SSHLoginCmd() SSHLoginCmdS {
	is := t.ai.InstanceState()
	if t.configs.Target.SSH.Direct {
		return SSHLoginCmdS{
			connection: fmt.Sprintf("%s@%s",
				t.configs.Target.SSH.Username,
				is.PublicIpAddress),
		}
	} else {
		return SSHLoginCmdS{
			option: fmt.Sprintf("ProxyCommand ssh %s@%s -W %%h:%%p 2> /dev/null",
				t.configs.Target.Bastion.Username,
				t.configs.Target.Bastion.Host),
			connection: fmt.Sprintf("%s@%s",
				t.configs.Target.SSH.Username,
				is.PrivateIpAddress),
		}
	}
}

func (t *Commands) SSHLogin(args []string) {

	opts := []string{}
	cmds := []string{}
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
	sl := t.SSHLoginCmd()
	if sl.option != "" {
		ag = append(ag, "-o", sl.option)
	}
	ag = append(ag, opts...)
	ag = append(ag, sl.connection)
	ag = append(ag, cmds...)
	log.Printf("debug: args %s", spew.Sdump(ag))
	t.runCommand("ssh", ag...)
}

func (t *Commands) RDPConnect() {
	is := t.ai.InstanceState()
	args := []string{
		"-N", "-L",
		fmt.Sprintf("127.0.0.1:%s:%s:3389",
			t.configs.Target.RDP.LocalPort,
			is.PrivateIpAddress),
		fmt.Sprintf("%s@%s",
			t.configs.Target.Bastion.Username,
			t.configs.Target.Bastion.Host),
	}
	fmt.Println("リモートデスクトップ接続から以下の情報で接続してください")
	fmt.Printf("   [ ホスト名   ] localhost:%s\n", t.configs.Target.RDP.LocalPort)
	if t.configs.Target.RDP.Username != "" {
		fmt.Printf("   [ ユーザ名   ] %s\n", t.configs.Target.RDP.Username)
	}
	if t.configs.Target.RDP.Password != "" {
		fmt.Printf("   [ パスワード ] %s\n", t.configs.Target.RDP.Password)
	}
	fmt.Println("ポートフォワーディングを開始しました。CTRL+Cで切断")
	t.runCommand("ssh", args...)

}
