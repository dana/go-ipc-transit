package main

import (
	"fmt"
	"github.com/dana/go-ipc-transit"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: rtrans <qname>")
		os.Exit(1)
	}
	srcQname := os.Args[1]
	/*
	       fmt.Println(sendMessage)
	   	sendMessage := map[string]interface{}{
	   		"Name": "Wednesday",
	   		"Age":  6,
	   		"Parents": map[string]interface{}{
	   			"bee": "boo",
	   			"foo": map[string]interface{}{
	   				"hi": []string{"a", "b"},
	   			},
	   		},
	   	} */
	message, recvErr := ipctransit.Receive(srcQname)
	if recvErr != nil {
		fmt.Println(recvErr)
		os.Exit(1)
	}
	fmt.Println(message)
	os.Exit(0)
}
