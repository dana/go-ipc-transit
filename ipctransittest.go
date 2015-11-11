package ipctransittest

import (
	"fmt"
	"os"
	"testing"
	"reflect"
)

var test_qname string = "ipc-transit-test-queue"

func TestSendRcv(t *testing.T) {
	defer func() {
		os.Remove(defaultTransitPath + test_qname)
	}()

	// How to create this message: http://play.golang.org/p/13OSJHd5xe
	// Info about seemingly fully dynamic marshal/unmarshal: http://stackoverflow.com/questions/19482612/go-golang-array-type-inside-struct-missing-type-composite-literal
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
	if sendErr != nil {
		t.Error(sendErr)
		return
	}
	m, receiveErr := Receive(test_qname)
	if receiveErr != nil {
		t.Error(receiveErr)
		return
	}
	msg := m.(map[string]interface{})
	for k, v := range msg {
		fmt.Println(k, " -> ", reflect.TypeOf(v))
		switch vv := v.(type) {
		case string:
			if k == "Name" {
				if v != "Wednesday" {
					t.Error(receiveErr)
				}
			}
			fmt.Println(k, "is string", vv)
		case float64:
			if k == "Age" {
				if v != 6.0 {  //very strange, even though it's an int in the json, it unmarshalled as a float
					t.Error(receiveErr)
				}
			}
			fmt.Println(k, "is float64", vv)
		case map[string]interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle", vv)
		}
	}
}
