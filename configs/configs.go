package configs

import (
	"log"
	"os"
	"errors"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"gopkg.in/yaml.v2"
	"github.com/mitchellh/go-homedir"
//	"github.com/davecgh/go-spew/spew"
)

type Configs struct {
	Configs C
	Target  STarget
}

func New() (t *Configs) {
	t = new(Configs)
	t.Configs = C{}
	return t
}

func (t *Configs) Load() error {

	sl := []string{}
	// バイナリと同じ場所か一階層上
	{
		exe, err := os.Executable()
		if err != nil { return err }
		b, err := filepath.Abs(filepath.Join(filepath.Dir(exe),"."))
		if err != nil { return err }
		sl = append(sl,
			filepath.Join(b,"ec2ctrl.yaml"),
			filepath.Join(b,"ec2ctrl.yml"),
			filepath.Join(b,".ec2ctrl.yaml"),
			filepath.Join(b,".ec2ctrl.yml"),
			filepath.Join(b,"../ec2ctrl.yaml"),
			filepath.Join(b,"../ec2ctrl.yml"),
			filepath.Join(b,"../.ec2ctrl.yaml"),
			filepath.Join(b,"../.ec2ctrl.yml"),
		)
	}
	// ホームディレクトリ
	{
		h, err := homedir.Dir()
		if err != nil { return err }
		sl = append(sl,
			filepath.Join(h,"ec2ctrl.yaml"),
			filepath.Join(h,"ec2ctrl.yml"),
			filepath.Join(h,".ec2ctrl.yaml"),
			filepath.Join(h,".ec2ctrl.yml"),
		)
	}
	configFile := ""
	for _,i := range sl {
		if _, err := os.Stat(i); !os.IsNotExist(err) {
			log.Printf("debug: FOUND %s",i)
			configFile = i
			break
		} else {
			log.Printf("debug: NOT FOUND %s",i)
		}
	}
	if configFile == "" {
		log.Println()
		return errors.New("alert: 設定ファイルがありません")
	}

	buf, err := ioutil.ReadFile(configFile)
	if err != nil { return err }
	s := regexp.MustCompile(`\r\n|\r|\n`).ReplaceAllString(string(buf),"\n")
	err = yaml.Unmarshal([]byte(s), &t.Configs)

	if err != nil { return err }

	return nil
}

func (t *Configs) SetTarget(tn string) bool {
	found := false
	tt := CTarget{}

	for _,i := range t.Configs.Targets {
		if i.Name == tn {
			tt = i
			found = true
			break
		}
	}
	if !found {
		return false
	}

	tg := STarget{
		Name:        tt.Name,
		GroupID:     tt.GroupID,
		InstanceID:  tt.InstanceID,
		Description: tt.Description,
		Enables:     tt.Enables,
		RDP:         tt.RDP,
		SSH:         tt.SSH,
	}

	if tt.AWS != "" {
		for _,i := range t.Configs.AWSes {
			if i.Name == tt.AWS {
				tg.AWS = i
				break
			}
		}
	}

	if tt.Bastion != "" {
		for _,i := range t.Configs.Bastions {
			if i.Name == tt.Bastion {
				tg.Bastion = i
				break
			}
		}
	}

	t.Target = tg
	return true

}

