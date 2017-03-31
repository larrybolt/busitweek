package entities

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strings"
)

//Part entity
type Part struct {
	Id            string
	Manufacturer  string
	PartIds       string
	Specification string
	Notes         string
}

type Parts []Part

func (p *Part) CreatePart(stub shim.ChaincodeStubInterface, args []string) error {
	p.Id = args[0]
	p.Manufacturer = args[1]
	p.Specification = args[2]
	p.Notes = args[3]
	p.PartIds = args[4]

	cJsonIndent, _ := json.MarshalIndent(p, "", " ")
	fmt.Println("CreatePart:", string(cJsonIndent))
	err := stub.PutState(p.Id, cJsonIndent)
	if err != nil {
		return err
	}
	//Update Part index
	idxPartsByte, _ := stub.GetState("idx_Parts")
	if idxPartsByte == nil {
		err := stub.PutState("idx_Parts", []byte(args[0]))
		if err != nil {
			return err
		}
		return nil
	} else {
		idxPartsByte = []byte(string(idxPartsByte) + "," + args[0])
		err := stub.PutState("idx_Parts", idxPartsByte)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

//Get a Part based on its ID
func (p *Part) GetPart(stub shim.ChaincodeStubInterface, Id string) ([]byte, error) {
	cJsonIndent, err := stub.GetState(Id)
	if err != nil {
		return nil, err
	}
	if cJsonIndent == nil {
		cJsonIndent = []byte("{\"Error\":\"Part with id" + Id + " not found\"}")
	}
	fmt.Println("GetPartJSON returned:", string(cJsonIndent))
	return cJsonIndent, nil
}

//List all Parts in the database
func (ps *Parts) ListParts(stub shim.ChaincodeStubInterface) ([]byte, error) {
	idxPartsByte, _ := stub.GetState("idx_Parts")
	fmt.Println(idxPartsByte)
	PartIDs := strings.Split(string(idxPartsByte), ",")
	fmt.Println(PartIDs)
	PartList := "{\"Parts\":["
	for i, PartID := range PartIDs {
		if i != 0 {
			PartList = PartList + ","
		}
		cJsonIndent, _ := stub.GetState(PartID)
		PartList = PartList + string(cJsonIndent)
		fmt.Println(PartList)
	}
	PartList = PartList + "]}"
	return []byte(PartList), nil
}

// //Update a Part in the database
func (p *Part) UpdatePart(stub shim.ChaincodeStubInterface, args []string) error {
	p.Id = args[0]
	p.Manufacturer = args[1]
	p.Specification = args[2]
	p.Notes = args[3]
	p.PartIds = args[4]

	cJsonIndent, err := stub.GetState(p.Id)
	if err != nil {
		return err
	}
	if cJsonIndent != nil {
		cJsonIndent, _ := json.MarshalIndent(p, "", " ")
		fmt.Println("CreatePart:", string(cJsonIndent))
		err := stub.PutState(p.Id, cJsonIndent)
		if err != nil {
			return err
		}
		//Update Part index
		idxPartsByte, _ := stub.GetState("idx_Parts")
		if idxPartsByte == nil {
			err := stub.PutState("idx_Parts", []byte(p.Id))
			if err != nil {
				return err
			}
			return nil
		} else {
			idxPartsByte = []byte(string(idxPartsByte) + "," + p.Id)
			err := stub.PutState("idx_Parts", idxPartsByte)
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	}
	fmt.Println("No entry found with ID:", string(p.Id))
	return nil
}
