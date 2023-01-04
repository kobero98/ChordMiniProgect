package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

type ChordNode int

var myMap map[int]string
var FingerTable []Node //ancora non utilizzata
var mySuccessivo Node
var myPrecedente Node

func calcolo_hash(text string) int {
	hash := md5.Sum([]byte(text))
	var test byte
	test = 0
	for i := 0; i < 8; i++ {
		test = hash[i] ^ test
	}
	return int(test)
}

func checkMyKey(key int) bool {
	_, ok := myMap[key]
	return ok
}

func (t *ChordNode) succ(key *int, reply *Node) {
	return
}
func (t *ChordNode) get() {
	return
}
func (t *ChordNode) put(parola *string, reply *int) {
	key := calcolo_hash(*parola)
	if checkMyKey(key) {
		myMap[key] = *parola
		//metti parola nella mia lista
	} else {
		client, err := rpc.DialHTTP("tcp", mySuccessivo.Ip[0]+strconv.Itoa(mySuccessivo.Port))
		if err != nil {
			log.Fatal("dialing:", err)
		}
		var reply ReplyRegistration
		err = client.Call("Manager.Register", &myNode, &reply)
		if err != nil {
			log.Fatal("arith error:", err)
		}
	}
	*reply = 1
	return
}

var Addr_Server_register = "0.0.0.0"

var myNode Node

func init_registration() (Node, Node) {
	//effettuo la connessione al server register
	client, err := rpc.DialHTTP("tcp", Addr_Server_register+":8000")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply ReplyRegistration
	err = client.Call("Manager.Register", &myNode, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Println(reply)
	client.Close()
	return reply.Precedente, reply.Successivo
}

func init_FingerTable(node *Node) {
	return
}

func init_Node() {
	var err error
	myNode.Name, err = os.Hostname()
	addr, err := net.LookupHost(myNode.Name)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	myNode.Ip = addr
	myNode.Port = 8001
	myNode.Index = calcolo_hash(myNode.Name)
	fmt.Print("ciao")
	fmt.Println(myNode)
	fmt.Print("ciao")
}
func init_map() {
	for i := myPrecedente.Index + 1; i <= myNode.Index; i = (i + 1) % 255 {
		myMap[i] = ""
	}
}
func main() {
	init_Node()
	//myPrecedente, mySuccessivo := init_registration()
	init_map()
}
