package awsi

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"time"
	//	"github.com/davecgh/go-spew/spew"
)

const (
	InstanceStatePending      = 0
	InstanceStateRunning      = 16
	InstanceStateShuttingDown = 32
	InstanceStateTerminated   = 48
	InstanceStateStopping     = 64
	InstanceStateStopped      = 80
)

type InstanceState struct {
	Name             string
	Description      string
	InstanceID       string
	InstanceType     string
	KeyName          string
	PrivateIpAddress string
	PublicIpAddress  string
	StateCode        int
	StateName        string
	StateJP          string
}

func (t *AWSi) InstanceState() InstanceState {
	svc := ec2.New(t.session)
	ret, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(t.configs.Target.InstanceID)},
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(ret.Reservations) == 0 {
		log.Fatal("fatal: インスタンスがありません(1)")
	}
	if len(ret.Reservations[0].Instances) == 0 {
		log.Fatal("fatal: インスタンスがありません(2)")
	}
	instance := ret.Reservations[0].Instances[0]

	stateJ := ""
	switch *instance.State.Code {
	case InstanceStatePending:
		stateJ = "保留中"
	case InstanceStateRunning:
		stateJ = "運転中"
	case InstanceStateShuttingDown:
		stateJ = "削除中"
	case InstanceStateTerminated:
		stateJ = "削除済"
	case InstanceStateStopping:
		stateJ = "停止中"
	case InstanceStateStopped:
		stateJ = "停止"
	}
	name := ""
	desc := ""
	for _, i := range instance.Tags {
		if *i.Key == "Name" {
			name = *i.Value
		}
		if *i.Key == "Description" {
			desc = *i.Value
		}
	}

	rev := func(p *string) string {
		if p == nil {
			return ""
		}
		return *p
	}

	return InstanceState{
		Name:             name,
		Description:      desc,
		InstanceID:       rev(instance.InstanceId),
		InstanceType:     rev(instance.InstanceType),
		KeyName:          rev(instance.KeyName),
		PrivateIpAddress: rev(instance.PrivateIpAddress),
		PublicIpAddress:  rev(instance.PublicIpAddress),
		StateCode:        int(*instance.State.Code),
		StateName:        rev(instance.State.Name),
		StateJP:          stateJ,
	}
}

func (t *AWSi) waitInstanceState(tp int) {
	tv := -1
	for i := 0; i < 60; i++ {
		is := t.InstanceState()
		if is.StateCode != tv {
			tv = is.StateCode
			fmt.Printf("  %s\n", is.StateJP)
		}
		if is.StateCode == tp {
			return
		}
		time.Sleep(time.Second)
	}
	log.Println("warn: タイムアウトしました")
	return
}

func (t *AWSi) InstanceStop() error {
	fmt.Printf("*** インスタンスの停止 ***\n")

	if t.InstanceState().StateCode != InstanceStateRunning {
		log.Printf("warn: インスタンスは動作していません")
		return nil
	}

	svc := ec2.New(t.session)
	_, err := svc.StopInstances(&ec2.StopInstancesInput{
		InstanceIds: []*string{aws.String(t.configs.Target.InstanceID)},
	})
	if err != nil {
		return err
	}
	fmt.Println("停止開始")
	t.waitInstanceState(InstanceStateStopped)
	return nil
}

func (t *AWSi) InstanceStart() error {
	fmt.Printf("*** インスタンスの起動 ***\n")
	if t.InstanceState().StateCode == InstanceStateRunning {
		log.Printf("warn: インスタンスは動作しています")
		return nil
	}

	svc := ec2.New(t.session)
	_, err := svc.StartInstances(&ec2.StartInstancesInput{
		InstanceIds: []*string{aws.String(t.configs.Target.InstanceID)},
	})
	if err != nil {
		return err
	}

	fmt.Println("起動開始")
	t.waitInstanceState(InstanceStateRunning)
	return nil
}

func (t *AWSi) InstanceStatus() error {
	fmt.Printf("*** インスタンスの状態 ***\n")
	is := t.InstanceState()
	fmt.Printf("  Name:             %s\n", is.Name)
	if is.Description != "" {
		fmt.Printf("  Description:      %s\n", is.Description)
	}
	fmt.Printf("  InstanceID:       %s\n", is.InstanceID)
	fmt.Printf("  InstanceType:     %s\n", is.InstanceType)
	fmt.Printf("  KeyName:          %s\n", is.KeyName)
	fmt.Printf("  PrivateIpAddress: %s\n", is.PrivateIpAddress)
	if is.PublicIpAddress != "" {
		fmt.Printf("  PublicIpAddress:  %s\n", is.PublicIpAddress)
	}
	fmt.Printf("  State:            %s(%s)\n", is.StateJP, is.StateName)
	return nil
}
