package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type sytemLogin struct {
	LoginId string
	Passwd  string
}

func processClients(cfd net.Conn) {
	var login sytemLogin
	json.NewDecoder(cfd).Decode(&login)
	fmt.Println("loginID: ", login.LoginId, ", passwd: ", login.Passwd)

	_, err := cfd.Write([]byte("AUTH SUCCESS"))
	if err != nil {
		return
	}
}

func initServer() error {
	lfd, err := net.Listen("unix", "/tmp/unix.socket")
	if err != nil {
		return err
	}
	defer os.Remove("/tmp/unix.socket")
	fmt.Println("Listening on /tmp/unix.socket")
	for {
		fd, err := lfd.Accept()
		if err != nil {
			return err
		}
		go processClients(fd)
	}
	return nil
}

func main() {
	err := initServer()
	if err != nil {
		fmt.Println(err)
	}
}
