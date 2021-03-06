package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	//	"github.com/davecgh/go-spew/spew"
	"github.com/mamemomonga/ec2ctrl/src/awsi"
	"github.com/mamemomonga/ec2ctrl/src/commands"
	"github.com/mamemomonga/ec2ctrl/src/configs"
	"github.com/mamemomonga/ec2ctrl/src/ec2ctrl/buildinfo"
	"github.com/spf13/cobra"
)

func main() {
	var cfg = configs.New()

	if err := cfg.Load(); err != nil {
		log.Fatal(err)
	}

	r := NewRunner(cfg)
	r.Cobra().Execute()
}

/* --------------------------------------- */

type Runner struct {
	cfg *configs.Configs
}

func NewRunner(c *configs.Configs) *Runner {
	t := new(Runner)
	t.cfg = c
	return t
}

func (t *Runner) Cobra() *cobra.Command {

	c := []*cobra.Command{}
	for _, i := range t.cfg.Configs.Targets {
		c = append(c, t.subCommands(i))
	}

	r := &cobra.Command{Use: os.Args[0]}
	r.AddCommand(c...)
	r.SetUsageTemplate(t.templateUsage())
	r.SetHelpTemplate(t.templateHelp())
	return r
}

func (t *Runner) subCommands(i configs.CTarget) *cobra.Command {
	c := &cobra.Command{
		Use:   i.Name,
		Short: i.Description,
	}
	if i.Enables.MyIP {
		c0 := &cobra.Command{
			Use:   "myip",
			Short: "パブリックIPアドレス",
		}
		c1 := &cobra.Command{
			Use:   "set",
			Short: "設定",
			Run:   func(cmd *cobra.Command, args []string) { t.action(i.Name, "MyIPSet") },
		}
		c2 := &cobra.Command{
			Use:   "del",
			Short: "削除",
			Run:   func(cmd *cobra.Command, args []string) { t.action(i.Name, "MyIPDel") },
		}
		c0.AddCommand(c1, c2)
		c.AddCommand(c0)
	}
	if i.Enables.Instance {
		c0 := &cobra.Command{
			Use:   "instance",
			Short: "EC2インスタンス",
		}
		c1 := &cobra.Command{
			Use:   "status",
			Short: "状態",
			Run:   func(cmd *cobra.Command, args []string) { t.action(i.Name, "InstanceStatus") },
		}
		c2 := &cobra.Command{
			Use:   "start",
			Short: "起動",
			Run:   func(cmd *cobra.Command, args []string) { t.action(i.Name, "InstanceStart") },
		}
		c3 := &cobra.Command{
			Use:   "stop",
			Short: "停止",
			Run:   func(cmd *cobra.Command, args []string) { t.action(i.Name, "InstanceStop") },
		}
		c0.AddCommand(c1, c2, c3)
		c.AddCommand(c0)
	}
	if i.Enables.SSH {
		c0 := &cobra.Command{
			Use:   "ssh",
			Short: "SSH関連操作",
		}
		c1 := &cobra.Command{
			Use:   "con",
			Short: "SSH接続",
			Long: "引数はそのままSSHコマンドに引き継がれます。ハイフンがあるとその先の引数はコマンドとして解釈されます。",
			DisableFlagParsing: true,
			Run: func(cmd *cobra.Command, args []string) {
				t.cfg.SetTarget(i.Name)
				ai := awsi.New(t.cfg)
				commands.New(t.cfg, ai).SSHLogin(args)
			},
		}
		c2 := &cobra.Command{
			Use:  "show",
			Short: "SSHコマンドを表示します",
			Run: func(cmd *cobra.Command, args []string) {
				t.cfg.SetTarget(i.Name)
				ai := awsi.New(t.cfg)
				commands.New(t.cfg, ai).SSHLoginCmdShow(args)
			},
		}
		c0.AddCommand(c1, c2)
		c.AddCommand(c0)
	}
	if i.Enables.RDP {
		c0 := &cobra.Command{
			Use:   "rdp",
			Short: "RDP",
			Run:   func(cmd *cobra.Command, args []string) { t.action(i.Name, "RDP") },
		}
		c.AddCommand(c0)
	}
	return c
}

func (t *Runner) action(target string, action string) {
	t.cfg.SetTarget(target)
	ai := awsi.New(t.cfg)
	cm := commands.New(t.cfg, ai)

	err := error(nil)

	switch action {
	case "MyIPSet":
		err = ai.MyIPSet()
	case "MyIPDel":
		err = ai.MyIPDel()
	case "InstanceStatus":
		err = ai.InstanceStatus()
	case "InstanceStop":
		err = ai.InstanceStop()
	case "InstanceStart":
		err = ai.InstanceStart()
	case "RDP":
		cm.RDPConnect()

	}

	if err != nil {
		log.Fatal(err)
	}

}

func (t *Runner) templateUsage() string {
	return `使い方:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

同名コマンド:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

使用例:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

コマンド一覧:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

フラグ:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

全体フラグ:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

追加ヘルプ:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
"{{.CommandPath}} [command] --help" を実行すればコマンドの詳細が確認できます{{end}}
`
}

func (t *Runner) templateHelp() string {
	s := fmt.Sprintf(`
====================================
ec2ctrl EC2コントロールツール %s-%s {{with (or .Long .Short)}}
  {{. | trimTrailingWhitespaces}}{{end}}
====================================
{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}
`,buildinfo.Version, buildinfo.Revision)
	return strings.TrimLeft(s,"\n")
}

