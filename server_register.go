package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Manager int

type appoggio struct {
	nodo   Node
	status int
}

var lista_nodi []appoggio
var count = 0
var request = 0

func printList() {
	fmt.Println("dimensione lista", len(lista_nodi))
	for i := 0; i < len(lista_nodi); i++ {
		fmt.Println(lista_nodi[i].nodo)
	}
}
func add_elemento(n appoggio) {
	if len(lista_nodi) == 0 {
		lista_nodi = append(lista_nodi, n)
		return
	}
	var app []appoggio
	app = append(app, n)
	index := -1
	for i := 0; i < len(lista_nodi); i++ {
		if lista_nodi[i].nodo.Index > n.nodo.Index {
			index = i
			break
		}
	}
	if index == 0 {
		lista_nodi = append(app, lista_nodi...)
		fmt.Println("dimensione lista", len(lista_nodi))
		return
	}
	if index == -1 {
		lista_nodi = append(lista_nodi, app...)
		fmt.Println("dimensione lista", len(lista_nodi))
		return
	}
	app = append(app, lista_nodi[index:]...)
	lista_nodi = append(lista_nodi[:index], app...)
	fmt.Println("dimensione lista", len(lista_nodi))
	return
}
func get_PrecSucc(n Node) (Node, Node) {

	for i := 0; i < len(lista_nodi); i++ {
		if n.Index == lista_nodi[i].nodo.Index {
			if i == 0 {
				return lista_nodi[len(lista_nodi)-1].nodo, lista_nodi[(i+1)%len(lista_nodi)].nodo
			}
			fmt.Println("valore lista ", (i-1)%len(lista_nodi), (i+1)%len(lista_nodi))
			return lista_nodi[(i-1)%len(lista_nodi)].nodo, lista_nodi[(i+1)%len(lista_nodi)].nodo
		}
	}
	return n, n
}
func (t *Manager) Register(node *Node, reply *ReplyRegistration) error {
	fmt.Println("un nodo si é connesso", count, request)
	fmt.Println(*node)
	var a appoggio
	a.nodo = *node
	a.status = 1
	add_elemento(a)
	count = count + 1
	if count == 1 {
		reply = nil
		return nil
	}
	/*var rep ReplyRegistration
	rep.Precedente, rep.Successivo = get_PrecSucc(*node)
	rep.NumNod = count
	reply = &rep
	*/
	reply.Precedente, reply.Successivo = get_PrecSucc(*node)
	reply.NumNod = count
	fmt.Println(*reply)
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
	lista_nodi = make([]appoggio, 0)
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
