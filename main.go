package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

func main() {

	client, err := rpc.DialHTTP("tcp", "0.0.0.0:8003")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply int
	var parola string
	parola = os.Args[1]
	fmt.Println(parola)
	err = client.Call("ChordNode.Put", &parola, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Println("risposta: ", reply)
	var parola2 string
	err = client.Call("ChordNode.Get", &reply, &parola2)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Println("parola messa: ", parola2)
}
