package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
)

func main() {

	client, err := rpc.DialHTTP("tcp", "register:8000")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var result int
	result = 0
	var contact Node
	err = client.Call("Manager.ContactClient", &result, &contact)
	if err != nil {
		log.Fatal("arith error:", err)
	}

	client.Close()
	fmt.Println("il nodo é: ", contact)
	var reply int
	var parola string
	client1, err1 := rpc.DialHTTP("tcp", contact.Name+":"+strconv.Itoa(contact.Port))
	if err1 != nil {
		log.Fatal("dialing:", err)
	}
	x, _ := strconv.Atoi(os.Args[2])
	if x == 1 {
		parola = os.Args[1]
		fmt.Println(parola)
		err = client1.Call("ChordNode.Put", &parola, &reply)
		if err != nil {
			log.Fatal("arith error:", err)
		}
		fmt.Println("risposta: ", reply)
	} else {
		reply, _ = strconv.Atoi(os.Args[1])
		var parola2 string
		err = client1.Call("ChordNode.Get", &reply, &parola2)
		if err != nil {
			log.Fatal("arith error:", err)
		}
		fmt.Println("parola messa: ", parola2)
	}
}
