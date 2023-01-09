package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net"
	"net/http"
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
func checkMyKey2(key int) bool {
	//forse bastava fare se minore del myIndex.precede
	if myNode.Index-myPrecedente.Index == 0 {
		return true
	}
	if myNode.Index-myPrecedente.Index > 0 {
		return key > myPrecedente.Index && key <= myNode.Index
	}
	return key > myPrecedente.Index || key <= myNode.Index
}
func checkMyKey(key int) bool {
	_, ok := myMap[key]
	return ok
}

func (t *ChordNode) Remove(key *int, reply *string) error {
	fmt.Println("mi hanno contattato per rimuove la kiave: ", *key)
	fmt.Println("io mi occupo di: ", myPrecedente.Index, myNode.Index)
	if checkMyKey2(*key) {
		str := myMap[*key]
		*reply = str
		return nil
	} else {
		client, err := rpc.DialHTTP("tcp", mySuccessivo.Ip[0]+":"+strconv.Itoa(mySuccessivo.Port))
		if err != nil {
			log.Fatal("dialing:", err)
		}
		err = client.Call("ChordNode.Get", key, reply)
		if err != nil {
			log.Fatal("arith error:", err)
		}
	}
	return nil
}

func (t *ChordNode) Get(key *int, reply *string) error {
	fmt.Println("mi hanno contattato per la kiave: ", *key)
	fmt.Println("io mi occupo di: ", myPrecedente.Index, myNode.Index)
	if checkMyKey2(*key) {
		str, test := myMap[*key]
		if test == false {
			str = "NOVALUE"
			return nil
		}
		*reply = str
		return nil
	} else {
		client, err := rpc.DialHTTP("tcp", mySuccessivo.Ip[0]+":"+strconv.Itoa(mySuccessivo.Port))
		if err != nil {
			log.Fatal("dialing:", err)
		}
		err = client.Call("ChordNode.Get", key, reply)
		if err != nil {
			log.Fatal("arith error:", err)
		}
	}
	return nil
}
func (t *ChordNode) Put(parola *string, reply *int) error {
	fmt.Println("mi hanno contattato per la parola: ", *parola)
	key := calcolo_hash(*parola)
	fmt.Println("la chiave  é ", key)
	fmt.Println("io mi occupo di: ", myPrecedente.Index, myNode.Index)
	if checkMyKey2(key) {
		myMap[key] = *parola
	} else {
		client, err := rpc.DialHTTP("tcp", mySuccessivo.Ip[0]+":"+strconv.Itoa(mySuccessivo.Port))
		if err != nil {
			myPrecedente, mySuccessivo = init_registration()
			//log.Fatal("dialing:", err)
			fmt.Println("ri ottentimento del precedente e del successivo")
			return err
		}
		err = client.Call("ChordNode.Put", parola, reply)
		if err != nil {
			myPrecedente, mySuccessivo = init_registration()
			//log.Fatal("arith error:", err)
			fmt.Println("ri ottentimento del precedente e del successivo")
			return err
		}
	}
	*reply = key
	return nil
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
		log.Fatal("errore nel ottenere l'indirizzo ip dell'host:", err)
	}
	myNode.Name = os.Args[1]
	myNode.Ip = addr
	// provvisorio
	myNode.Ip = make([]string, 1)
	myNode.Ip[0] = "127.0.0.1"
	fmt.Println(myNode.Ip)

	myNode.Port = 8001
	myNode.Port, _ = strconv.Atoi(os.Args[2])
	myNode.Index = calcolo_hash(myNode.Name)
	fmt.Println(myNode)
}
func (t *ChordNode) Precedente(node *Node, reply *int) error {
	myPrecedente = *node
	*reply = 1
	return nil
}
func (t *ChordNode) Successivo(node *Node, reply *int) error {
	mySuccessivo = *node
	return nil
}
func (t *ChordNode) HeartBit(answer *int, reply *int) error {
	*reply = 1
	return nil
}
func comunicationToSuccessivo() {
	client, err := rpc.DialHTTP("tcp", "0.0.0.0"+":"+strconv.Itoa(mySuccessivo.Port))
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply int
	err = client.Call("ChordNode.Precedente", &myNode, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	client.Close()
}
func comunicationToPrecedente() {
	client, err := rpc.DialHTTP("tcp", "0.0.0.0"+":"+strconv.Itoa(myPrecedente.Port))
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply int
	err = client.Call("ChordNode.Successivo", &myNode, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	client.Close()
}
func ChangeStatus() {
	client, err := rpc.DialHTTP("tcp", "0.0.0.0"+":8000")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply int
	err = client.Call("Manager.ChangeStatus", &myNode, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	client.Close()
}
func main() {
	init_Node()
	myMap = make(map[int]string)
	myPrecedente, mySuccessivo = init_registration()
	if myPrecedente.Name == "" && myPrecedente.Port == 0 {
		myPrecedente = myNode
		mySuccessivo = myNode

	} else {
		comunicationToPrecedente()
		comunicationToSuccessivo()
	}
	chord := new(ChordNode)
	rpc.Register(chord)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+os.Args[2])
	if e != nil {
		log.Fatal("listen error:", e)
	}
	ChangeStatus()
	http.Serve(l, nil)
	fmt.Println("fine programma in go")
}
