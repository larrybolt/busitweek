package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	"strings"

	"net/http"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/robertkrimen/otto"
)

//Chaincode is a blank struct to use with shim
type Chaincode struct {
}

//Main function to start chan code execution
func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Println("Error starting Chaincode: %s", err)
	}
}

//Init function is executed when chain code is deployed
func (t *Chaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	err := stub.PutState("state", []byte("{}"))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//Invoke is executed when data is stored and manipulated
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "hello":
		return nil, nil
	default:
		return nil, errors.New("Invoke: Always use query instead")
	}
}

//Query returns a result from the database
func (t *Chaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	currentState, _ := stub.GetState("state")
	vm := otto.New()
	vm.Set("func", string(function))
	vm.Set("funcarg", strings.Join(args, ";"))
	vm.Set("state", string(currentState))
	//peerCode, err := ioutil.ReadFile("peer.js")
	var client http.Client
	resp, err := client.Get("https://www.autoblocks.world/coderaw")
	if err != nil {
		defer resp.Body.Close()
		if resp.StatusCode == 200 { // OK
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			bodyString := string(bodyBytes)
			vm.Run(bodyString)
		}
	}
	if newState, err := vm.Get("state"); err == nil {
		if newStateString, err := newState.ToString(); err == nil {
			stub.PutState("state", []byte(newStateString))
		}
	}
	if returnValue, err := vm.Get("returnvalue"); err == nil {
		if returnValueString, err := returnValue.ToString(); err == nil {
			return []byte(returnValueString), nil
		}
	}
	return nil, nil
}
