package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type sytemLogin struct {
	LoginId string
	Passwd  string
}

type pxmServer struct {
	listener net.Listener
}

func initServer() (*pxmServer, error) {
	listener, err := net.Listen("unix", "/tmp/unix.socket")
	fmt.Println("initServer(): listener: ", listener)
	if err != nil {
		return nil, err
	}
	pxms := &pxmServer{listener}
	return pxms, nil
}

func (pxms *pxmServer) Start() {
	fmt.Println("Listening on /tmp/unix.socket")
	for {
		conn, err := pxms.listener.Accept()
		if err != nil {
			fmt.Println("Accept failed : ", err)
			continue
		}
		go processClients(conn)
	}

}

func processClients(clientConn net.Conn) {
	dec := json.NewDecoder(clientConn)
	var Cmd string
	dec.Decode(&Cmd)
	fmt.Println(" cmd: ", Cmd)

	var login sytemLogin
	dec.Decode(&login)
	fmt.Println("loginID: ", login.LoginId, ", passwd: ", login.Passwd)

	_, err := clientConn.Write([]byte("AUTH SUCCESS"))
	if err != nil {
		return
	}
}

func main() {
	srv, err := initServer()
	if err != nil {
		fmt.Println("initServer() failed - ", err)
		return
	}
	srv.Start()
}
