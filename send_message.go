package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
)

func readData(c net.Conn) (string, error) {
	rawData := make([]byte, 4096)
	if _, err := c.Read(rawData); err != nil {
		return "", err
	}

	data := string(rawData)

	dataLines := strings.Split(data, "\n")
	for _, v := range dataLines {
		log.Printf("S: %s", v)
	}

	return string(data), nil
}

func writeData(c net.Conn, data string) error {
	log.Printf("C: %s\n", data)
	_, err := c.Write([]byte(fmt.Sprintf("%s\r\n", data)))
	if err != nil {
		return err
	}

	return nil
}

func writeDataAndRead(c net.Conn, data string) (string, error) {
	if err := writeData(c, data); err != nil {
		return "", err
	}

	rData, err := readData(c)
	if err != nil {
		return "", err
	}
	return rData, nil
}

func writeMany(c net.Conn, data []string) error {
	for _, v := range data {
		if err := writeData(c, v); err != nil {
			return err
		}
	}
	return nil
}

func writeManyAndRead(c net.Conn, data []string) ([]string, error) {
	ret := make([]string, len(data))
	var err error
	for i, v := range data {
		ret[i], err = writeDataAndRead(c, v)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func getConfigForServer(serverName string) (tls.Config, error) {
	cert, _ := tls.LoadX509KeyPair("./cert.pem", "./privkey.pem")
	return tls.Config{
		ServerName:   serverName,
		Certificates: []tls.Certificate{cert},
	}, nil
}

func encryptedConn(unencConn net.Conn, config *tls.Config) (net.Conn, error) {
	conn := tls.Client(unencConn, config)
	err := conn.Handshake()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func getHostFromEmail(email string) (string, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("wrong email format")
	}
	return parts[1], nil
}

func sendMessage(to, from, body string) {
	log.Println("start J-M-L server")

	host, err := getHostFromEmail(to)
	if err != nil {
		log.Fatal(err)
	}

	mxrecords, err := net.LookupMX(host)
	if err != nil {
		log.Fatal(err)
	}

	serverHostMX := mxrecords[0].Host
	serverAdress := fmt.Sprintf("%s:%v", serverHostMX, defaultPort)
	log.Printf("connecting to %s\n", serverAdress)
	conn, err := net.Dial("tcp", serverAdress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())

	if _, err := readData(conn); err != nil {
		log.Fatal(err)
	}

	requests := []string{
		"EHLO smtp.staff.defernus.com",
		"STARTTLS",
	}

	if _, err := writeManyAndRead(conn, requests); err != nil {
		log.Fatal(err)
	}

	config, err := getConfigForServer(serverHostMX)
	if err != nil {
		log.Fatal(err)
	}

	conn, err = encryptedConn(conn, &config)
	if err != nil {
		log.Fatal(err)
	}

	requests = []string{
		fmt.Sprintf("MAIL FROM:<%s>", from),
		fmt.Sprintf("RCPT TO:<%s>", to),
		"DATA",
	}

	if _, err := writeManyAndRead(conn, requests); err != nil {
		log.Fatal(err)
	}

	requests = []string{
		//fmt.Sprintf("From: defernus <%s>", from),
		//fmt.Sprintf("Subject: %s", subject),
		//fmt.Sprintf("To: <%s>", to),
		body,
		".",
	}

	if err := writeMany(conn, requests); err != nil {
		log.Fatal(err)
	}

	if _, err := readData(conn); err != nil {
		log.Fatal(err)
	}

	requests = []string{
		"QUIT",
	}

	if _, err := writeManyAndRead(conn, requests); err != nil {
		log.Fatal(err)
	}

}
