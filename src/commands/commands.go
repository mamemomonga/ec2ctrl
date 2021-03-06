package commands

import (
	"os"
	"os/exec"

	"github.com/mamemomonga/ec2ctrl/src/awsi"
	"github.com/mamemomonga/ec2ctrl/src/configs"
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

