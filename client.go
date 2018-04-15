package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

var wg sync.WaitGroup

type sytemLogin struct {
	LoginId string
	Passwd  string
}

func sendReq1() error {
	fd, err := net.Dial("unix", "/tmp/unix.socket")
	if err != nil {
		return err
	}
	defer fd.Close()

	req := "SYSLOGIN1"
	err = json.NewEncoder(fd).Encode(req)
	if err != nil {
		log.Fatal(err)
		return err
	}

	loginInfo := sytemLogin{"admin@localhost1", "admin1"}
	err = json.NewEncoder(fd).Encode(loginInfo)
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

func sendReq2() error {
	fd, err := net.Dial("unix", "/tmp/unix.socket")
	if err != nil {
		return err
	}
	defer fd.Close()

	req := "SYSLOGIN2"
	err = json.NewEncoder(fd).Encode(req)
	if err != nil {
		log.Fatal(err)
		return err
	}

	loginInfo := sytemLogin{"admin@localhost2", "admin2"}
	err = json.NewEncoder(fd).Encode(loginInfo)
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

func sendReq3() error {
	fd, err := net.Dial("unix", "/tmp/unix.socket")
	if err != nil {
		return err
	}
	defer fd.Close()

	req := "SYSLOGIN3"
	err = json.NewEncoder(fd).Encode(req)
	if err != nil {
		log.Fatal(err)
		return err
	}

	loginInfo := sytemLogin{"admin@localhost3", "admin3"}
	err = json.NewEncoder(fd).Encode(loginInfo)
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
	wg.Add(3)
	go sendReq1()
	go sendReq2()
	go sendReq3()
	wg.Wait()
}
