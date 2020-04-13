package configs

type C struct {
	Username    string     `yaml:"username"`
	AWSes       []CAWS     `yaml:"awses"`
	Bastions    []CBastion `yaml:"bastions"`
	Targets     []CTarget  `yaml:"targets"`
}

type CAWS struct {
	Name            string `yaml:"name"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	Region          string `yaml:"region"`
	AccountNumber   string `yaml:"account_number"`
}

type CBastion struct {
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Host     string `yaml:"host"`
}

type CTarget struct {
	Name        string `yaml:"name"`
	AWS         string `yaml:"aws"`
	Bastion     string `yaml:"bastion"`
	GroupID     string `yaml:"group_id"`
	InstanceID  string `yaml:"instance_id"`
	Description string `yaml:"description"`
	Enables     CTargetEnables `yaml:"enables"`
	SSH         CTargetSSH      `yaml:"ssh"`
	RDP         CTargetRDP     `yaml:"rdp"`
}

type CTargetEnables struct {
	 MyIP     bool `yaml:"myip"`
     Instance bool `yaml:"instance"`
     RDP      bool `yaml:"rdp"`
     SSH      bool `yaml:"ssh"`
}

type CTargetRDP struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	LocalPort string `yaml:"local_port"`
}

type CTargetSSH struct {
	Username  string `yaml:"username"`
	Direct    bool   `yaml:"direct"`
}

type STarget struct {
	Name        string
	AWS         CAWS
	Bastion     CBastion
	Username    string
	GroupID     string
	InstanceID  string
	Description string
	LocalPort   string
	Enables     CTargetEnables
	RDP         CTargetRDP
	SSH         CTargetSSH
}

