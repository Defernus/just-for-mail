package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func write(conn net.Conn, str string) error {
	log.Printf("S: %s", str)
	if _, err := conn.Write([]byte(fmt.Sprintf("%s\r\n", str))); err != nil {
		return err
	}
	return nil
}

func read(conn net.Conn) (string, error) {
	rawResult := make([]byte, 4096)
	n, err := conn.Read(rawResult)
	if err != nil {
		return "", err
	}
	result := strings.TrimSpace(string(rawResult[:n]))
	log.Printf("C: %s", result)
	return result, nil
}

func sendGreeting(conn net.Conn) error {
	if err := write(conn, fmt.Sprintf("220 %s Simple Mail Transfer Service Ready", curentHost)); err != nil {
		return err
	}
	return nil
}

func readRequest(conn net.Conn) (*SMTPRequest, error) {
	netData, err := read(conn)
	if err != nil {
		return nil, err
	}

	return ParseRequest(netData), nil
}

func handleConnection(conn net.Conn) {
	log.Printf("Serving %s\n", conn.RemoteAddr().String())

	if err := sendGreeting(conn); err != nil {
		log.Fatal(err)
	}

	handler := NewRequestHandler()

	isDataReceiving := false

	for {
		if isDataReceiving {
			netData, err := read(conn)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("C: %s", netData)

			response, action := handler.HandleData(netData)
			switch action {
			case dataActionClose:
				conn.Close()
				return
			case dataActionDataEnd:
				isDataReceiving = false
				if err := write(conn, response); err != nil {
					log.Fatal(err)
				}
			}
		} else {
			request, err := readRequest(conn)
			if err != nil {
				if err.Error() == "EOF" {
					log.Printf("Conection closed by %s\n", conn.RemoteAddr().String())
					return
				}
				log.Fatal(err)
			}
			response, action := handler.HandleRequest(request)
			if err := write(conn, response); err != nil {
				log.Fatal(err)
			}
			switch action {
			case requestActionClose:
				conn.Close()
				return
			case requestActionData:
				isDataReceiving = true
			}
		}
	}

}
