package commands

import (
	"github.com/mamemomonga/ec2ctrl/configs"
	"github.com/mamemomonga/ec2ctrl/awsi"
//	"log"
	"fmt"
	"os"
	"os/exec"
)

type Commands struct {
	ai      *awsi.AWSi
	configs *configs.Configs
}

func New(configs *configs.Configs, ai *awsi.AWSi) *Commands {
	t := new(Commands)
	t.ai      = ai
	t.configs = configs
	return t
}

func (t *Commands) runCommand(c string, p ...string) error {
	cmd := exec.Command(c, p...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin  = os.Stdin
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func (t *Commands) SSHLogin(run bool) {
	is := t.ai.InstanceState()
	args := []string{}
	if t.configs.Target.SSH.Direct {
		args=[]string{
			fmt.Sprintf("%s@%s",
				t.configs.Target.SSH.Username,
				is.PublicIpAddress ),
		}
	} else {
		args=[]string{
			"-o",
			fmt.Sprintf( "ProxyCommand ssh %s@%s -W %%h:%%p 2> /dev/null",
				t.configs.Target.Bastion.Username,
				t.configs.Target.Bastion.Host ),
			fmt.Sprintf( "%s@%s",
				t.configs.Target.SSH.Username,
				is.PrivateIpAddress ),
		}
	}
	if(run) {
		t.runCommand("ssh",args...)
	} else {
		if t.configs.Target.SSH.Direct {
			fmt.Printf("ssh %s\n",args[0])
		} else {
			fmt.Printf("ssh %s \"%s\" %s\n",args[0],args[1],args[2])
		}
	}
}

func (t *Commands) RDPConnect() {
	is := t.ai.InstanceState()
	args := []string{
		"-N", "-L",
		fmt.Sprintf("127.0.0.1:%s:%s:3389",
			t.configs.Target.RDP.LocalPort,
			is.PrivateIpAddress ),
		fmt.Sprintf("%s@%s",
			t.configs.Target.Bastion.Username,
			t.configs.Target.Bastion.Host ),
	}
	fmt.Println("リモートデスクトップ接続から以下の情報で接続してください")
	fmt.Printf("   [ ホスト名   ] localhost:%s\n", t.configs.Target.RDP.LocalPort )
	if t.configs.Target.RDP.Username != "" {
		fmt.Printf("   [ ユーザ名   ] %s\n", t.configs.Target.RDP.Username  )
	}
	if t.configs.Target.RDP.Password != "" {
		fmt.Printf("   [ パスワード ] %s\n", t.configs.Target.RDP.Password  )
	}
	fmt.Println("ポートフォワーディングを開始しました。CTRL+Cで切断")
	t.runCommand("ssh",args...)

}

