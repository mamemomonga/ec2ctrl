package configs

import (
	"log"
	"os"
	"io/ioutil"
	"path/filepath"
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
	homedir, err := homedir.Dir()
	if err != nil { return err }

	cf := filepath.Join(homedir,".ec2ctrl.yaml")
	if _, err := os.Stat(cf); os.IsNotExist(err) {
		log.Printf(" %s ファイルがありません\n",cf)
		return err
	}

	buf, err := ioutil.ReadFile(cf)
	if err != nil { return err }

	err = yaml.Unmarshal(buf, &t.Configs)
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

