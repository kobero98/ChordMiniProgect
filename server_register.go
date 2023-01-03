package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Node struct {
	name  string
	ip    []string
	port  int
	Index int
}

type Manager int

type appoggio struct {
	nodo   Node
	status int
}

var lista_nodi [20]appoggio
var count = 0
var request = 0

func printList() {
	for i := 0; i < count; i++ {
		fmt.Println(lista_nodi[i].nodo)
	}
}
func (t *Manager) Register(node *Node, reply *Node) error {
	fmt.Println("un nodo si é connesso", count, request)
	fmt.Println(*node)
	lista_nodi[count].nodo = *node
	lista_nodi[count].status = 1
	count = count + 1
	if count == 1 {
		reply = nil
		return nil
	}
	reply = &lista_nodi[request].nodo
	request = (request + 1) % count
	printList()
	return nil
}
func (t *Manager) Unregister(node *Node, reply *Node) error {
	fmt.Println("un nodo si é disconnesso")
	count = count - 1
	reply = nil
	return nil
}
func main() {
	fmt.Println("inizio programma in go")
	manage := new(Manager)
	rpc.Register(manage)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":8000")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
	fmt.Println("fine programma in go")
}

/*

func (t *Manager) Registretion(node *Node,reply *Node, numNod *int) error{
	lista_nodi[count]=node;
	count=count+1;
	numNod=count
	if count==1{
		reply=nil
		return nil
	}
	else{
		reply=lista_nodi[request]
		request=request+1
		return nil
	}
}

func (t *Manager) Unregistration(node *Node,cod *int){
	for i := 0; i < count; i++ {
		if lista_nodi[i].Index==node.Index {
			lista_nodi[i]=nil
			count=count-1;
			cod=0
			return nil
		}

	}
	cod=-1
	return nil
}
*/
