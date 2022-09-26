package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

//ref madflojo.medium.com/keeping-tcp-connections-alive-in-golang-801a78b7cf1
func server(addr *net.TCPAddr) error {

	// Start TCP Listener
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return fmt.Errorf("Unable to start listener: %v", err)
	}

	// Wait for new connections and send them to reader()
	c, err := l.AcceptTCP()
	if err != nil {
		return fmt.Errorf("Listener returned: %v", err)
	}

	// Enable Keepalives
	err = c.SetKeepAlive(false)
	if err != nil {
		return fmt.Errorf("Unable to set keepalive: %v", err)
	}
	for true {
		msg, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			return fmt.Errorf("Unable to read from client: %v", err)
		}
		fmt.Println("received: " + strings.TrimSuffix(msg, "\n"))
	}
	return nil
}

func client(addr *net.TCPAddr) error {
	// Open TCP Connection
	c, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return fmt.Errorf("Unable to dial to server: %v", err)
	}

	err = c.SetKeepAlive(false)
	if err != nil {
		return fmt.Errorf("Unable to set keepalive: %v", err)
	}
	for true {
		time.Sleep(1 * time.Second)
		msg := "ping"
		fmt.Println("send: " + msg)
		_, err = fmt.Fprintf(c, msg+"\n")
		if err != nil {
			return fmt.Errorf("Unable to send msg: %v", err)
		}
	}
	return nil
}

func main() {
	kind := os.Args[1]
	addr := os.Args[2]

	// Resolve TCP Address
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic("Unable to resolve IP")
	}
	if kind == "s" {
		err = server(tcpAddr)

	} else if kind == "c" {
		err = client(tcpAddr)
	}
	if err != nil {
		panic(err)
	}
}
