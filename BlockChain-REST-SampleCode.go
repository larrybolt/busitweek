package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//var url = "http://172.18.0.2:7050/chaincode"
var url = "http://localhost:7050/chaincode"

//var peer = "noops_vp0_1"
var peer = "pbft_vp0_1"

var chaincodeID string

//var chaincodeID = "bee3e6b93a92d5b93c6daf985b4a7efee81c2b8d3a872f10ca0d4d40da8563eb7c23fd128f4cc6be59ae28e5c7941894611df37e4ab101e5630df56e531e5033"

func main() {

	fmt.Println("URL:>", url)
	fmt.Println("\n")
	var answer string

	answer = DeployChaincode("UCLL", "")
	fmt.Println(answer)
	time.Sleep(10 * time.Second)
	fmt.Println("\n")

	/*
		answer = QueryChaincode("hello", `"0001"`, 3)
		fmt.Println(answer)
		time.Sleep(5 * time.Second)
		fmt.Println("\n")
	*/

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
	fmt.Println(chaincodeID)

	return sbody
}
