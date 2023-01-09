package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
)

func main() {

	client, err := rpc.DialHTTP("tcp", "0.0.0.0:8003")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply int
	var parola string
	x, _ := strconv.Atoi(os.Args[2])
	if x == 1 {
		parola = os.Args[1]
		fmt.Println(parola)
		err = client.Call("ChordNode.Put", &parola, &reply)
		if err != nil {
			log.Fatal("arith error:", err)
		}
		fmt.Println("risposta: ", reply)
	} else {
		reply, _ = strconv.Atoi(os.Args[1])
		var parola2 string
		err = client.Call("ChordNode.Get", &reply, &parola2)
		if err != nil {
			log.Fatal("arith error:", err)
		}
		fmt.Println("parola messa: ", parola2)
	}
}
