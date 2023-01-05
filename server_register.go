package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"time"
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
func remove_elemento(n appoggio) int {
	if len(lista_nodi) == 0 {
		return 0
	}
	for i := 0; i < len(lista_nodi); i++ {
		if lista_nodi[i].nodo.Index == n.nodo.Index {
			lista_nodi = append(lista_nodi[:i], lista_nodi[i+1:]...)
			return 0
		}
	}
	return -1
}
func add_elemento(n appoggio) int {
	if len(lista_nodi) == 0 {
		lista_nodi = append(lista_nodi, n)
		return 0
	}
	var app []appoggio
	app = append(app, n)
	index := -1
	for i := 0; i < len(lista_nodi); i++ {
		if lista_nodi[i].nodo.Index == n.nodo.Index {
			if lista_nodi[i].nodo.Name == n.nodo.Name && lista_nodi[i].nodo.Port == n.nodo.Port && lista_nodi[i].nodo.Ip[0] == n.nodo.Ip[0] {
				return -1
			}
			return -2
		}
		if lista_nodi[i].nodo.Index > n.nodo.Index {
			index = i
			break
		}
	}
	if index == 0 {
		lista_nodi = append(app, lista_nodi...)
		fmt.Println("dimensione lista", len(lista_nodi))
		return 0
	}
	if index == -1 {
		lista_nodi = append(lista_nodi, app...)
		fmt.Println("dimensione lista", len(lista_nodi))
		return 0
	}
	app = append(app, lista_nodi[index:]...)
	lista_nodi = append(lista_nodi[:index], app...)
	fmt.Println("dimensione lista", len(lista_nodi))
	return 0
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
func (t *Manager) ChangeStatus(node *Node, reply *int) error {
	fmt.Println("ricevuto un cambio di stato dal nodo: ", node.Index)
	for i := 0; i < len(lista_nodi); i++ {
		if lista_nodi[i].nodo.Index == node.Index {
			lista_nodi[i].status = 1
			*reply = 0
			return nil
		}
	}
	*reply = -1
	return nil
}
func (t *Manager) Register(node *Node, reply *ReplyRegistration) error {
	fmt.Println("un nodo si é connesso", count, request)
	fmt.Println(*node)
	var a appoggio
	a.nodo = *node
	a.status = 0
	contr := add_elemento(a)
	if contr == -2 {
		reply = nil
		return nil
	}
	count = count + 1 + contr
	if count == 1 {
		reply = nil
		return nil
	}
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
func heartBit() {
	for {
		for i := 0; i < len(lista_nodi); i++ {
			if lista_nodi[i].status == 1 {
				client, err := rpc.DialHTTP("tcp", lista_nodi[i].nodo.Ip[0]+":"+strconv.Itoa(lista_nodi[i].nodo.Port))
				if err != nil {
					fmt.Println("Elemento rimosso 1")
					remove_elemento(lista_nodi[i])
					continue
				}
				var reply int
				var answer int
				err = client.Call("ChordNode.HeartBit", &answer, &reply)
				if err != nil {
					fmt.Println("Elemento rimosso 2")
					remove_elemento(lista_nodi[i])
					client.Close()
					continue
				}
				client.Close()
			}
		}
		time.Sleep(10 * time.Second)
	}
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
	go heartBit()
	http.Serve(l, nil)
	fmt.Println("fine programma in go")

}
