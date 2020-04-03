package main

import (
	"github.com/mamemomonga/ec2ctrl/configs"
//	"github.com/davecgh/go-spew/spew"
	"log"
)

var cfg *configs.Configs

func main() {
	cfg = configs.New()

	if err := cfg.Load(); err != nil {
		log.Fatal(err)
	}

	if err := runner(); err != nil {
		log.Fatal(err)
	}
}

