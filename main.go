package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
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
	rand.Seed(time.Now().Unix())

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(c)
	}
}
