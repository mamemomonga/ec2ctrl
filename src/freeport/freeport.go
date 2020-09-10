package freeport

import (
	"net"
	"fmt"
	"regexp"
	"errors"
)

func SearchTCP(start int, end int) (int,error) {
	r := regexp.MustCompile("address already in use")
	for port := start; port <= end; port++ {
		addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d",port))
		if err != nil {
			return 0, err
		}
		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			if r.MatchString(err.Error()) {
				continue
			}
			return 0,err
		}
		defer l.Close()
		return l.Addr().(*net.TCPAddr).Port, nil
	}
	return 0,errors.New("all ports already in use")
}

