package main

import (
	"encoding/json"
	"fmt"
	"github.com/dana/go-ipc-transit"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: strans <qname> <message>")
		os.Exit(1)
	}
	destQname := os.Args[1]
	rawJSON := os.Args[2]
	var sendMessage map[string]interface{}
	jsonErr := json.Unmarshal([]byte(rawJSON), &sendMessage)
	if jsonErr != nil {
		fmt.Println("invalid json")
		os.Exit(1)
	}
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
	sendErr := ipctransit.Send(sendMessage, destQname)
	if sendErr != nil {
		fmt.Println(sendErr)
		os.Exit(1)
	}
	os.Exit(0)
}
