package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

var wg sync.WaitGroup

type sytemLogin struct {
	LoginId string
	Passwd  string
}

func sendReq() error {
	fd, err := net.Dial("unix", "/tmp/unix.socket")
	if err != nil {
		return err
	}
	defer fd.Close()
	loginInfo := sytemLogin{"admin@localhost", "admin"}
	json.NewEncoder(fd).Encode(loginInfo)

	resp := make([]byte, 13)
	_, err = fd.Read(resp)
	fmt.Println(string(resp))
	wg.Done()
	return nil
}

func main() {
	wg.Add(5)
	go sendReq()
	go sendReq()
	go sendReq()
	go sendReq()
	go sendReq()
	wg.Wait()
}
