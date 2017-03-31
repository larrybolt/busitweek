package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var url = "http://localhost:7050/chaincode"
var peer = "pbft_vp0_1"
var chaincodeID string

func main() {

	fmt.Println("URL:>", url)
	fmt.Println("\n")
	var answer string

	answer = DeployChaincode("UCLL", "")
	fmt.Println(answer)
	time.Sleep(5 * time.Second)
	fmt.Println("\n")

	 answer = InvokeChaincode("createPart", `"0010", "Lay's", "Mars", "Cote d'or", "0009"`, 6)
	 fmt.Println(answer)
	 time.Sleep(10 * time.Second)
	 fmt.Println("\n")

	answer = InvokeChaincode("createPart", `"0009", "aaa", "bbb", "ccc", "0001"`, 6)
	fmt.Println(answer)
	time.Sleep(10 * time.Second)
	fmt.Println("\n")

	answer = QueryChaincode("getPart", `"0010"`, 7)
	fmt.Println(answer)
	time.Sleep(5 * time.Second)
	fmt.Println("\n")

	// answer = InvokeChaincode("updatePart", `"0009", "JellyBean", "Bounty", "KitKat", "0010"`, 8)
	// fmt.Println(answer)
	// time.Sleep(10 * time.Second)
	// fmt.Println("\n")

	answer = QueryChaincode("listParts", `""`, 9)
	fmt.Println(answer)
	time.Sleep(5 * time.Second)
	fmt.Println("\n")

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

	//Copy source code
	cmdStr := "docker cp $GOPATH/src/" + ccDirectory + " " + peer + ":go/src/"
	out, _ := exec.Command("/bin/sh", "-c", cmdStr).Output()
	fmt.Printf("%s", out)

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
