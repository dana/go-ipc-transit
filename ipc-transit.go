package ipc-transit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var defaultTransitPath string = "/tmp/ipc_transit/"

func Send(sendMessage map[string]interface{}, qname string) error {
    var wireHeader = make(map[string]string)
	wireHeader["q"] = qname
	sendBytes, createWireHeaderErr := createWireHeader(wireHeader)
	if createWireHeaderErr != nil {
		return createWireHeaderErr
	}
	mq, err := getQueue(qname)
	if err != nil {
		return err
	}
	jsonBytes, marshalErr := json.Marshal(sendMessage)
	if marshalErr != nil {
		return marshalErr
	}
	sendBytes = append(sendBytes, jsonBytes...)
	sendErr := RawSend(sendBytes, mq)
	return sendErr
}

func RawSend(rawBytes []byte, queue MessageQueue) error {
	err := queue.Send(1, rawBytes, nil)
	return err
}

func Receive(qname string) (interface{}, error) {
	var f interface{}
	mq, err := getQueue(qname)
	if err != nil {
		return nil, err
	}
	rawBytes, receiveErr := RawReceive(mq)
	if receiveErr != nil {
		return nil, receiveErr
	}
	wireHeader, payload, parseErr := parseWireHeader(rawBytes)
	if _, ok := wireHeader["q"]; ok {
		fmt.Println("recieved from q = " + wireHeader["q"])
	}
	if parseErr != nil {
		return nil, parseErr
	}
	jsonErr := json.Unmarshal(payload, &f)
	fmt.Println(f)
	if jsonErr != nil {
		return f, jsonErr
	}

	return f, receiveErr
}
func RawReceive(queue MessageQueue) ([]byte, error) {
	rawBytes, _, err := queue.Receive(102400000, -1, nil)
	return rawBytes, err
}

//  $ cat /tmp/ipc_transit/foo
//  qid=17039435
//  qname=foo
type transitInfo struct {
	qid   int64
	qname string
}

func parseTransitFile(filePath string) (transitInfo, error) {
	info := transitInfo{0, ""}
	fi, err := os.Open(filePath)
	if err != nil {
		return info, err
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		things := strings.Split(scanner.Text(), "=")
		key := string(things[0])
		value := things[1]
		switch key {
		case "qid":
			my_qid, atoiErr := strconv.Atoi(string(value))
			if atoiErr != nil {
				return info, atoiErr
			}
			info.qid = int64(my_qid)
		case "qname":
			info.qname = string(value)
		}
	}
	if err := scanner.Err(); err != nil {
		return info, err
	}
	return info, err
}

func makeNewQueue(qname string, queuePath string) error {
	fmt.Println("makeNewQueue: " + qname)
	if _, statErr := os.Stat(defaultTransitPath); os.IsNotExist(statErr) {
		//dir does not exist
		mkdirErr := os.Mkdir(defaultTransitPath, 0777)
		if mkdirErr != nil {
			return mkdirErr
		}
	}
	fi, err := os.Create(queuePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	fi.WriteString("qname=" + qname + "\n")
	key, ftokErr := Ftok(queuePath, 100)
	fi.WriteString("qid=" + strconv.Itoa(int(key)))
	return ftokErr
}

func getQueue(qname string) (MessageQueue, error) {
	transitInfoFilePath := defaultTransitPath + qname
	if _, statErr := os.Stat(transitInfoFilePath); os.IsNotExist(statErr) {
		makeErr := makeNewQueue(qname, transitInfoFilePath)
		if makeErr != nil {
			panic(makeErr)
		}
	}
	info, err := parseTransitFile(transitInfoFilePath)
	if err != nil {
		return MessageQueue(0), err
	}
	mq, err := GetMsgQueue(info.qid, &MQFlags{
		Create: true,
		//		Create:    false,
		//		Exclusive: true,
		//		Exclusive: false,
		Perms: 0666,
	})
	if false {	//this is about always having a use of fmt.Println so I never
				//have to take it out of imports
		fmt.Println("whatevah")
	}
	return mq, err
}


func createWireHeader(headerMap map[string]string) ([]byte, error) {
    headerBytes := []byte("")
    for key, value := range headerMap {
        headerBytes = append(headerBytes, key...)
        headerBytes = append(headerBytes, "="...)
        headerBytes = append(headerBytes, value...)
        headerBytes = append(headerBytes, ","...)
    }
    if len(headerBytes) > 0 {
        headerBytes = headerBytes[:len(headerBytes)-1]
    }
    ret := []byte(strconv.Itoa(len(headerBytes)))
    ret = append(ret, ":"...)
    ret = append(ret, headerBytes...)
    return ret, nil
}

func parseWireHeader(testInput []byte) (map[string]string, []byte, error) {
    var retMap = make(map[string]string)
    testString := string(testInput)
    fullHeaderParts := strings.SplitN(testString, ":", 2)
    headerLength,atoiErr := strconv.Atoi(fullHeaderParts[0])
    if atoiErr != nil {
        return retMap, nil, atoiErr
    }
    headerString := fullHeaderParts[1][0:headerLength]
    payload := testInput[len(fullHeaderParts[0])+headerLength+1:]
    headerParts := strings.Split(headerString, ",")
    for _, part := range headerParts {
		if len(part) > 0 {
        	fields := strings.Split(part, "=")
        	key := fields[0]
        	value := fields[1]
        	retMap[key] = value
		}
    }
    return retMap, payload, nil
}


/*
Sat Nov  7 18:09:43 PST 2015
TODO:
Locking around the IPC Transit file manipulation
Nonblocking flags
Handle (and respect) the custom IPC::Transit header
	basic creating and parsing is working
Obviously turn this into a proper package
Large message handling
Remote transit
queue stats
internal(local) queues
testing with alternate directories
*/

/*
Sun Nov  8 12:24:56 PST 2015

Full example of message including header

11:d=localhost{".ipc_transit_meta":{"destination":"localhost","ttl":9,"destination_qname":"test","send_ts":1447014248},"1":2}
*/


