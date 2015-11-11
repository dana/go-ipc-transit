package main

import (
	"fmt"
	"github.com/dana/go-ipc-transit"
)

func main() {
	fmt.Println("howdy")
    sendMessage := map[string]interface{}{
        "Name": "Wednesday",
        "Age":  6,
        "Parents": map[string]interface{}{
            "bee": "boo",
            "foo": map[string]interface{}{
                "hi": []string{"a", "b"},
            },
        },
    }
    sendErr := Send(sendMessage, test_qname)
	fmt.Println(sendErr)
}
