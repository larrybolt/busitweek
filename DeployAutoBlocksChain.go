package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var url = "http://localhost:7050/chaincode"
var peer = "pbft_vp0_1"
var chaincodeID string

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func main() {

	fmt.Println("URL:>", url)
	fmt.Println("\n")
	var answer string

	answer = DeployChaincode("AutoBlocks", "")
	fmt.Println(answer)
	time.Sleep(3 * time.Second)
	fmt.Println("\n")

	answer = InvokeChaincode("createPart", `"TESTPART-0004", "BLOCKCHAIN", ""`, 1)
	fmt.Println(answer)
	time.Sleep(3 * time.Second)
	fmt.Println("\n")

	answer = QueryChaincode("listParts", `""`, 2)
	fmt.Println(answer)
	time.Sleep(3 * time.Second)
	fmt.Println("\n")

	answer = QueryChaincode("getPart", `"TESTPART-0004"`, 3)
	fmt.Println(answer)
	time.Sleep(3 * time.Second)
	fmt.Println("\n")

	// answer = InvokeChaincode("updatePart", `"0009", "JellyBean", "Bounty", "KitKat", "0010"`, 8)
	// fmt.Println(answer)
	// time.Sleep(10 * time.Second)
	// fmt.Println("\n")

}

//Invoke Chaincode through REST  api
func InvokeChaincode(bcFunction string, bcParameters string, postID int) string {
	requestString := fmt.Sprintf(`{"jsonrpc": "2.0","method": "invoke","params": {"type": 1,"chaincodeID":{"name":"%s"},"ctorMsg": {"function":"%s","args":[%s]}},"id": %v}`, chaincodeID, bcFunction, bcParameters, postID)

	var jsonStr = []byte(requestString)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

//Query Chaincode through REST  api
func QueryChaincode(bcFunction string, bcParameters string, postID int) string {
	requestString := fmt.Sprintf(`{"jsonrpc": "2.0","method": "query","params": {"type": 1,"chaincodeID":{"name":"%s"},"ctorMsg": {"function":"%s","args":[%s]}},"id": %v}`, chaincodeID, bcFunction, bcParameters, postID)

	var jsonStr = []byte(requestString)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return removeEscape(body)
}

//function to remove EscapeCharacters
func removeEscape(input []byte) string {
	output := strings.Replace(string(input), "\\\"", "\"", -1)
	output = strings.Replace(output, "\\n", "", -1)
	return output
}

//Deploying chaincode
func DeployChaincode(ccDirectory string, bcParameters string) string {

	//Deploy chaincode
	requestString := fmt.Sprintf(`{"jsonrpc": "2.0","method": "deploy","params": {"type": 1,"chaincodeID":{"path":"%s"},"ctorMsg": {"function":"init","args":[%s]}},"id": 0}`, ccDirectory, bcParameters)
	var jsonStr = []byte(requestString)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	//Get new Chaincode ID
	sbody := string(body)
	i := strings.Index(sbody, "message")
	chaincodeID = sbody[i+10 : i+138]

	return sbody
}
