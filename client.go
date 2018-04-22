package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

type cmd struct {
	Cmd string
}
type sytemLogin struct {
	LoginId string
	Passwd  string
}

var wg sync.WaitGroup

func resetPasswd() error {
	fd, err := net.Dial("unix", "/tmp/unix.socket")
	if err != nil {
		return err
	}
	defer fd.Close()

	req := cmd{Cmd: "RESETPASSWD"}
	err = json.NewEncoder(fd).Encode(req)
	if err != nil {
		log.Fatal(err)
		return err
	}
	params := sytemLogin{LoginId: "admin@localhost", Passwd: "newadmin"}
	err = json.NewEncoder(fd).Encode(params)
	if err != nil {
		log.Fatal(err)
		return err
	}

	resp := make([]byte, 13)
	_, err = fd.Read(resp)
	fmt.Println(string(resp))
	wg.Done()
	return nil
}

func sysLogin() error {
	fd, err := net.Dial("unix", "/tmp/unix.socket")
	if err != nil {
		return err
	}
	defer fd.Close()

	req := cmd{Cmd: "SYSLOGIN"}
	err = json.NewEncoder(fd).Encode(req)
	if err != nil {
		log.Fatal(err)
		return err
	}
	params := sytemLogin{LoginId: "admin@localhost", Passwd: "admin"}
	err = json.NewEncoder(fd).Encode(params)
	if err != nil {
		log.Fatal(err)
		return err
	}

	resp := make([]byte, 13)
	_, err = fd.Read(resp)
	fmt.Println(string(resp))
	wg.Done()
	return nil
}

func main() {
	wg.Add(2)
	go sysLogin()
	go resetPasswd()
	wg.Wait()
}
