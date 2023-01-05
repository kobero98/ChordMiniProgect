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
			log.Fatal("dialing:", err)
		}
		err = client.Call("ChordNode.Put", parola, reply)
		if err != nil {
			log.Fatal("arith error:", err)
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
	myNode.Name = os.Args[1]
	addr, err := net.LookupHost(myNode.Name)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	myNode.Ip = addr
	myNode.Port = 8001
	myNode.Port, _ = strconv.Atoi(os.Args[2])
	myNode.Index = calcolo_hash(myNode.Name)
	fmt.Print("ciao")
	fmt.Println(myNode)
	fmt.Print("ciao")
}
func init_map(flag int) {
	myMap = make(map[int]string)
	for i := (myPrecedente.Index + 1) * flag % 256; i != (myNode.Index+1)*flag+256*(1-flag); i = (i + 1) % (256*flag + 257*(1-flag)) {
		myMap[i] = os.Args[2]
	}
}
func (t *ChordNode) Precedente(node *Node, reply *map[int]string) error {
	//var c int = 0
	/*for i := (myPrecedente.Index + 1) % 256; i != (node.Index+1)%256; i = (i + 1) % 256 {
		(*reply)[i] = myMap[i]
		delete(myMap, i)
		c++
	}*/
	myPrecedente = *node
	//fmt.Println(myMap)
	return nil
}
func (t *ChordNode) Successivo(node *Node, reply *int) error {
	mySuccessivo = *node
	return nil
}
func comunicationToSuccessivo() {
	client, err := rpc.DialHTTP("tcp", "0.0.0.0"+":"+strconv.Itoa(mySuccessivo.Port))
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply = make(map[int]string)
	err = client.Call("ChordNode.Precedente", &myNode, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	/*for key, value := range reply {
		myMap[key] = value
	}*/
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
func main() {
	init_Node()
	myMap = make(map[int]string)
	myPrecedente, mySuccessivo = init_registration()
	if myPrecedente.Name == "" {
		//init_map(0)
		myPrecedente = myNode
		mySuccessivo = myNode

	} else {
		//init_map(1)
		comunicationToPrecedente()
		comunicationToSuccessivo()
	}
	fmt.Println(myMap)
	chord := new(ChordNode)
	rpc.Register(chord)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+os.Args[2])
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
	fmt.Println("fine programma in go")

}
