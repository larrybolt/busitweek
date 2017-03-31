package main

import (
	"AutoBlocks/entities"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
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
	return nil, nil
}

//Invoke is executed when data is stored and manipulated
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "createPart":
		var p entities.Part
		err := p.CreatePart(stub, args)
		if err != nil {
			return nil, err
		}
	case "updatePart":
		var p entities.Part
		err := p.UpdatePart(stub, args)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Invoke: Received unknonw function name")
	}
	return nil, nil
}

//Query returns a result from the database
func (t *Chaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "getPart":
		var p entities.Part
		cJsonIndent, err := p.GetPart(stub, args[0])
		return cJsonIndent, err
	case "listParts":
		var ps entities.Parts
		psJsonIndent, err := ps.ListParts(stub)
		return psJsonIndent, err
	default:
		return nil, errors.New("Invoke: Received unknown function name")
	}
	return nil, nil
}
