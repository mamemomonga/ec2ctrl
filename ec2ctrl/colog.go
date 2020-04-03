package main

import (
	"github.com/comail/colog"
	"log"
	"os"
)

func init() {
	if os.Getenv("DEBUG") != "" {
		colog.SetMinLevel(colog.LTrace)
		colog.SetDefaultLevel(colog.LDebug)
	} else {
		colog.SetMinLevel(colog.LInfo)
		colog.SetDefaultLevel(colog.LWarning)
	}
	colog.SetFormatter(&colog.StdFormatter{
		Colors: true,
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
	})
	colog.Register()
}

