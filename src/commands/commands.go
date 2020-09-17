package commands

import (
	"fmt"
	"github.com/mamemomonga/ec2ctrl/src/awsi"
	"github.com/mamemomonga/ec2ctrl/src/configs"
	"github.com/mamemomonga/ec2ctrl/src/freeport"
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

func (t *Commands) RDPConnect() {

	port,err := freeport.SearchTCP(33891,33900)
	if err != nil {
		log.Fatal(err)
	}

	is := t.ai.InstanceState()
	args := []string{
		"-N", "-L",
		fmt.Sprintf("127.0.0.1:%d:%s:3389",
			port,
			is.PrivateIpAddress),
		fmt.Sprintf("%s@%s",
			t.configs.Target.Bastion.Username,
			t.configs.Target.Bastion.Host),
	}
	fmt.Println("リモートデスクトップ接続から以下の情報で接続してください")
	fmt.Printf("   [ ホスト名   ] localhost:%d\n", port)
	if t.configs.Target.RDP.Username != "" {
		fmt.Printf("   [ ユーザ名   ] %s\n", t.configs.Target.RDP.Username)
	}
	if t.configs.Target.RDP.Password != "" {
		fmt.Printf("   [ パスワード ] %s\n", t.configs.Target.RDP.Password)
	}
	fmt.Println("ポートフォワーディングを開始しました。CTRL+Cで切断")
	t.runCommand("ssh", args...)

}
