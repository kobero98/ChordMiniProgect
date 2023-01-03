package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/rpc"
	"os"
)

var server_register = "0.0.0.0"
var s [128]Node

func Looking(s string) {
	fmt.Println("fase di ricerca ancora non introdotta")
}

func init() {
	client, err := rpc.DialHTTP("tcp", server_register+":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	client.Call()

}

func main() {
	fmt.Println("Hello Word")
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Hostname: %s", hostname)
	s := hostname
	h := md5.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	fmt.Println(s)
	fmt.Println(bs)
	fmt.Printf("%x\n", bs)
}
