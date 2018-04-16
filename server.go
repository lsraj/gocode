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
	listener      net.Listener
	clientReqChan chan net.Conn
	serverDone    chan bool
}

func initServer() (*pxmServer, error) {
	listener, err := net.Listen("unix", "/tmp/unix.socket")
	if err != nil {
		return nil, err
	}
	reqChan := make(chan net.Conn)
	done := make(chan bool)
	pxms := &pxmServer{listener, reqChan, done}
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
		pxms.clientReqChan <- conn
	}
	close(pxms.clientReqChan)

}

func (pxms *pxmServer) processClients() {
	for conn := range pxms.clientReqChan {
		dec := json.NewDecoder(conn)
		var Cmd string
		dec.Decode(&Cmd)
		fmt.Println(" cmd: ", Cmd)

		var login sytemLogin
		dec.Decode(&login)
		fmt.Println("loginID: ", login.LoginId, ", passwd: ", login.Passwd)

		_, err := conn.Write([]byte("AUTH SUCCESS"))
		if err != nil {
			return
		}
		conn.Close()
	}
	pxms.serverDone <- true

}

func main() {
	srv, err := initServer()
	if err != nil {
		fmt.Println("initServer() failed - ", err)
		return
	}
	go srv.Start()
	go srv.processClients()
	<-srv.serverDone
}
