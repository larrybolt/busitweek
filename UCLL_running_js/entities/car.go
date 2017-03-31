package entities

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/robertkrimen/otto"
)

//Insert an new car in the database + update index (CarId Separated by ,)
func (c *Car) ProxyQuery(stub shim.ChaincodeStubInterface, args []string) error {

	fmt.Println("ProxyQuery:", args[0])
	err := stub.PutState('state', '')
	if err != nil {
		return err
	}
	//Update Car index
	return nil
}

//Get a car based on its ID
func (c *Car) ProxyInvoke(stub shim.ChaincodeStubInterface, carID string) ([]byte, error) {
	cJsonIndent, err := stub.GetState(carID)
	if err != nil {
		return nil, err
	}
	if cJsonIndent == nil {
		cJsonIndent = []byte("{\"Error\":\"Car with ID" + carID + " not found\"}")
	}
	fmt.Println("GetCarJSON returned:", string(cJsonIndent))
	return cJsonIndent, nil
}

//List all cars in the database
func (cs *Cars) ListCars(stub shim.ChaincodeStubInterface) ([]byte, error) {
	idxCarsByte, _ := stub.GetState("idx_Cars")
	carIDs := strings.Split(string(idxCarsByte), ",")
	carList := "{\"Cars\":"
	for i, carID := range carIDs {
		if i != 0 {
			carList = carList + ","
		}
		cJsonIndent, _ := stub.GetState(carID)
		carList = carList + string(cJsonIndent)
	}
	carList = carList + "\n}"
	return []byte(carList), nil
}

//Load Cars sample data
func (cs *Cars) LoadSample(stub shim.ChaincodeStubInterface) string {
	var c Car
	vm := otto.New()
	vm.Run(`abc = 2 + 2;`)
	argslist := make([][]string, 6)
	argslist[0] = []string{"0001", "Renault", "Megane", "1600D", "2012"}
	argslist[1] = []string{"0002", "Mercedes", "C-Class", "220", "2014"}
	argslist[2] = []string{"0003", "Ford", "Focus", "1.8 16V", "2005"}
	argslist[3] = []string{"0004", "Renault", "Clio", "1200cc", "2014"}
	argslist[4] = []string{"0005", "Opel", "Astra", "1.9CDTI", "2011"}
	argslist[5] = []string{"0006", "Opel", "Astra", "2.0", "2010"}

	for _, args := range argslist {
		c.CreateCar(stub, args)
	}
	return "Load Car samples: 6 inserted"
}
