package awsi

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"io/ioutil"
	"log"
	"net/http"
//	"github.com/davecgh/go-spew/spew"
)

func (t *AWSi) myIPDeleteCurrent(groupid, mmyip, mdesc string) (bool, error) {
	svc := ec2.New(t.session)

	// 現在の設定を取得
	current, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		GroupIds: aws.StringSlice([]string{groupid}),
	})
	if err != nil {
		return false, err
	}

	// すでに存在するIPアドレス・詳細を含む項目を削除する
	for _,sg := range current.SecurityGroups {
		for _, ipp := range sg.IpPermissions {
			for _, ipr := range ipp.IpRanges {
				flag := false
				ipranges := []*ec2.IpRange{}
				// 自分のIPを含んでいたら消す
				if mmyip == *ipr.CidrIp {
					flag = true
					ipranges = []*ec2.IpRange{{
						CidrIp: aws.String(mmyip),
					}}
				}
				// 自分の設定した項目なら消す
				if mdesc == *ipr.Description {
					flag = true
					ipranges = []*ec2.IpRange{{
						CidrIp:      ipr.CidrIp,
						Description: aws.String(mdesc),
					}}
				}
				// 削除項目なし
				if !flag {
					continue
				}
				// 削除実行
				_, err := svc.RevokeSecurityGroupIngress(&ec2.RevokeSecurityGroupIngressInput{
					GroupId: aws.String(groupid),
					IpPermissions: []*ec2.IpPermission{{
						IpProtocol: aws.String("tcp"),
						FromPort:   aws.Int64(22),
						ToPort:     aws.Int64(22),
						IpRanges:   ipranges,
					}},
				})
				if err != nil {
					return false, err
				}
				return true, nil
			}
		}
	}
	return false, nil
}

func (t *AWSi) myIPAdd(groupid, mmyip, mdesc string) error {
	svc := ec2.New(t.session)

	// 新しい設定
	_, err := svc.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String(groupid),
		IpPermissions: []*ec2.IpPermission{{
			FromPort:   aws.Int64(22),
			ToPort:     aws.Int64(22),
			IpProtocol: aws.String("tcp"),
			IpRanges: []*ec2.IpRange{{
				CidrIp:      aws.String(mmyip),
				Description: aws.String(mdesc),
			}},
		}},
	})
	if err != nil {
		return err
	}
	return nil
}

func (t *AWSi) MyIPSet() error {
	fmt.Printf("*** 許可IPアドレスの追加・更新 ***\n")
	myip := t.getMyIP()

	mmyip := fmt.Sprintf("%s/32", myip)
	mdesc := fmt.Sprintf("rmt-%s", t.configs.Configs.Username)

	fmt.Printf("  SecurityGroupID: %s\n", t.configs.Target.GroupID)
	fmt.Printf("  許可IP:          %s\n", mmyip)
	fmt.Printf("  Description:     %s\n", mdesc)

	deleted, err := t.myIPDeleteCurrent(t.configs.Target.GroupID, mmyip, mdesc)
	if err != nil {
		return err
	}
	if deleted {
		fmt.Printf("以前の設定を削除しました\n")
	}

	err = t.myIPAdd(t.configs.Target.GroupID, mmyip, mdesc)
	if err != nil {
		return err
	}
	fmt.Printf("設定しました\n")

	return nil
}

func (t *AWSi) MyIPDel() error {
	fmt.Printf("*** 許可IPアドレスの削除 ***\n")

	mdesc := fmt.Sprintf("rmt-%s", t.configs.Configs.Username)
	fmt.Printf("   SecurityGroupID: %s\n", t.configs.Target.GroupID)
	fmt.Printf("   Description:     %s\n", mdesc)

	deleted, err := t.myIPDeleteCurrent(t.configs.Target.GroupID, "", mdesc)
	if err != nil {
		return err
	}
	if deleted {
		fmt.Printf("以前の設定を削除しました\n")
	}

	return nil
}

func (t *AWSi) getMyIP() string {
	type httpBin struct {
		Origin string `json:"origin"`
	}
	res, err := http.Get("http://httpbin.org/ip")
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatal(fmt.Sprintf("Status Code: %d", res.StatusCode))
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var b httpBin
	err = json.Unmarshal([]byte(body), &b)
	if err != nil {
		log.Fatal(err)
	}
	return string(b.Origin)
}
