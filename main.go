package main

import (
	"fmt"
	"log"
	"net"
)

const (
	defaultPort = 25
	curentHost  = "staff.defernus.com"
)

func main() {
	log.Println("start J-M-L server")

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%v", curentHost, defaultPort))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(c)
	}
}
