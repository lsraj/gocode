package main

import "encoding/json"
import "fmt"
import "net"
import "os"
import "os/signal"
import "syscall"
import "time"

var srvShutdown chan bool

type netServer struct {
	listener     *net.UnixListener // unix local sock
	taskChan     chan net.Conn     // unbuf channel for client connections
	sigChan      chan os.Signal    // buffered channel
	srvDone      chan bool
	numOfWorkers int
}

type cmd struct {
	Cmd string
}
type sytemLogin struct {
	LoginId string
	Passwd  string
}

// initialize signals.
func initSignals() {
	// block all async signals to this server
	signal.Ignore()
}

// initializes server:
//   1) creates listener.
//   2) creates task channel - client connections accepted on this channel.
//   3) initializes few other things such as signal channel, num of workers etc.
func initServer() (*netServer, error) {
	sockAddr, err := net.ResolveUnixAddr("unix", "/tmp/unix.socket")
	if err != nil {
		fmt.Println("initServer() failed on ResolveUnixAddr(): ", err)
		return nil, err
	}
	srvlistener, err := net.ListenUnix("unix", sockAddr)
	if err != nil {
		fmt.Println("initServer() failed on ListenUnix(): ", err)
		return nil, err
	}
	netSrv := &netServer{
		listener:     srvlistener,
		taskChan:     make(chan net.Conn),
		sigChan:      make(chan os.Signal, 1),
		srvDone:      make(chan bool),
		numOfWorkers: 3,
	}
	srvShutdown = make(chan bool)
	initSignals()
	return netSrv, nil
}

// Accept client connections. Each connection is directed onto
// a channel. Worker threads pickup and work on these connections.
func (netSrv *netServer) acceptClientReqs() {
	fmt.Println("Listening on /tmp/unix.socket")
	for {
		select {
		case <-srvShutdown:
			netSrv.srvDone <- true
			fmt.Println("Exiting server")
			return
		default:
			netSrv.listener.SetDeadline(time.Now().Add(time.Second))
			conn, err := netSrv.listener.Accept()
			if err != nil {
				fmt.Println("Accept failed : ", err)
				continue
			}
			netSrv.taskChan <- conn
		}
	}
}

func processRequest(conn net.Conn) {
	dec := json.NewDecoder(conn)
	var Cmd cmd
	dec.Decode(&Cmd)
	fmt.Println(" cmd: ", Cmd.Cmd)
	switch Cmd.Cmd {
	case "SYSLOGIN":
		var login sytemLogin
		dec.Decode(&login)
		fmt.Println("Login: ", login.LoginId, "Passwd: ", login.Passwd)
		// Call Authenticate func
	case "RESETPASSWD":
		var login sytemLogin
		dec.Decode(&login)
		fmt.Println("Login: ", login.LoginId, "Passwd: ", login.Passwd)
		// Call resetPasswd func
	}

	// reply to client
	_, err := conn.Write([]byte("AUTH SUCCESS"))
	if err != nil {
		return
	}
	conn.Close()
}

// Pickup a client connection, work on it and close connection.
func (netSrv *netServer) taskWorker(workerId int) {
	fmt.Println("taskWorker(): starting worker-", workerId)
	// for conn := range netSrv.taskChan {
	for {
		select {
		case <-srvShutdown:
			fmt.Println("taskWorker(): taskWorker-", workerId, " shutting down")
			return
		case conn := <-netSrv.taskChan:
			if conn == nil {
				fmt.Println("taskWorker(): worker-", workerId, " continue")
				continue
			}
			processRequest(conn)
			fmt.Println("taskWorker(): worker-", workerId, " processed client connection: ", conn)
		}
	}
}

// Handle SIGINT and SIGHUP for now.
func (netSrv *netServer) signalHandler() {
	signal.Notify(netSrv.sigChan, syscall.SIGINT, syscall.SIGHUP)
	sig := <-netSrv.sigChan
	fmt.Println("signalHandler() received signal", sig)
	srvShutdown <- true
	close(srvShutdown)
}

func (netSrv *netServer) cleanUp() {
	close(netSrv.taskChan)
	netSrv.listener.Close()
	close(netSrv.srvDone)
}

func main() {
	srv, err := initServer()
	if err != nil {
		fmt.Println("initServer() failed - ", err)
		return
	}
	go srv.signalHandler()

	for i := 1; i <= srv.numOfWorkers; i++ {
		go srv.taskWorker(i)
	}
	go srv.acceptClientReqs()

	// wait server to complete
	<-srv.srvDone
	srv.cleanUp()
}
