# ec2ctrl

* AWS EC2の起動・停止・基本情報の参照
* 自分のIPアドレスをSecurityGroupに登録する
* SSH接続情報・RDP接続情報の表示

# 設定方法

	$ cp ec2ctrl-example.yaml ~/.ec2ctrl.yaml
	$ vim ~/.ec2ctrl.yaml

プログラムは以下の流れで設定ファイルを探します

* [BINDIR]/ec2ctrl.yaml
* [BINDIR]/ec2ctrl.yml
* [BINDIR]/.ec2ctrl.yaml
* [BINDIR]/.ec2ctrl.yml
* [BINDIR]/../ec2ctrl.yaml
* [BINDIR]/../ec2ctrl.yml
* [BINDIR]/../.ec2ctrl.yaml
* [BINDIR]/../.ec2ctrl.yml
* [HOMEDIR]/ec2ctrl.yaml
* [HOMEDIR]/ec2ctrl.yml
* [HOMEDIR]/.ec2ctrl.yaml
* [HOMEDIR]/.ec2ctrl.yml

