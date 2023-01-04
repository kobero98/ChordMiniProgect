package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"
)

var node Node

var Addr_Server_register = "0.0.0.0"

func main() {
	//effettuo la connessione al server register
	client, err := rpc.DialHTTP("tcp", Addr_Server_register+":8000")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var node Node
	node.Name, err = os.Hostname()
	addr, err := net.LookupHost(node.Name)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	node.Ip = addr
	node.Port = 8001
	var reply Node
	fmt.Print(node)
	err = client.Call("Manager.Register", &node, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Println(reply)
	duration := time.Second
	time.Sleep(duration * 10)

	err = client.Call("Manager.Unregister", &node, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
}
