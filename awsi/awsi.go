package awsi

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/mamemomonga/ec2ctrl/configs"
    "log"
	"os"
)

type AWSi struct {
	session *session.Session
	configs *configs.Configs
}

func New(configs *configs.Configs) *AWSi {
	t := new(AWSi)

	t.configs = configs

	os.Setenv("AWS_ACCESS_KEY_ID",     t.configs.Target.AWS.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", t.configs.Target.AWS.SecretAccessKey)

    sess, err := session.NewSession(&aws.Config{ Region: aws.String(t.configs.Target.AWS.Region) })
	if err != nil {
		log.Println(err)
		log.Fatal("alert: セッション開始に失敗しました。")
	}
	t.session = sess
	t.checkMatchAccountNumber()
	return t
}

func (t *AWSi) checkMatchAccountNumber() {
	svc := sts.New(t.session)
	r,err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Println(err)
		log.Fatal("alert: アカウント番号取得に失敗しました。認証設定を確認してください")
	}
	if *r.Account != t.configs.Target.AWS.AccountNumber {
		log.Printf("Target:%s Result:%s \n", t.configs.Target.AWS.AccountNumber,  *r.Account)
		log.Fatal("alert: アカウント番号が合致しません。認証設定を確認してください")
	}
}

