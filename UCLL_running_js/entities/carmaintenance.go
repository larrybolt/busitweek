package entities

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"sort"
	"strconv"
	"time"
)

type Maintenance struct {
	Garage      string
	Type        string
	Date        time.Time
	Km          int
	Description string
}

type Maintenances []Maintenance

type CarMaintenance struct {
	CarID          string
	MantenanceList Maintenances
}

//The below 3 functions implement a Sort interface (order by Date)
func (m Maintenances) Len() int {
	return len(m)
}
func (m Maintenances) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
func (m Maintenances) Less(i, j int) bool {
	return m[i].Date.Before(m[j].Date)
}

//Add maintenance to a car
func (m *Maintenance) AddMaintenance(stub shim.ChaincodeStubInterface, args []string) error {
	carID := args[0]
	m.Garage = args[1]
	m.Type = args[2]
	m.Date, _ = time.Parse("02-01-2006", args[3])
	m.Km, _ = strconv.Atoi(args[4])
	m.Description = args[5]

	var cm CarMaintenance
	cmJson, _ := stub.GetState("cm-" + carID)
	if cmJson == nil {
		cm.CarID = carID
	} else {
		err := json.Unmarshal(cmJson, &cm)
		if err != nil {
			return errors.New("AddMaintence: Error in unmarshaling JSON")
		}
	}
	cm.MantenanceList = append(cm.MantenanceList, *m)
	sort.Sort(cm.MantenanceList)
	cmJsonIndent, _ := json.MarshalIndent(cm, "", "  ")
	err := stub.PutState("cm-"+carID, cmJsonIndent)
	if err != nil {
		return errors.New("AddMaintenance: Unable to PutState")
	}
	return nil
}

//Getting mainteance list for a car
func (cm *CarMaintenance) GetCarMaintenceList(stub shim.ChaincodeStubInterface, carID string) ([]byte, error) {
	cmJsonIndent, err := stub.GetState("cm-" + carID)
	if err != nil {
		return nil, err
	}
	if cmJsonIndent == nil {
		cmJsonIndent = []byte("{\"Error\":\"No maintenance track available for carID " + carID + "\"}")
	}
	fmt.Println("GetCarJSON returned:", string(cmJsonIndent))
	return cmJsonIndent, nil
}

//Loading Maintenance sample code
func (cm *CarMaintenance) LoadMaintenanceSample(stub shim.ChaincodeStubInterface) string {
	var m Maintenance

	argslist := make([][]string, 5)
	argslist[0] = []string{"0001", "Garage A", "Maintenance", "10-02-2012", "1000", "First checkup"}
	argslist[1] = []string{"0001", "Garage B", "Repair", "05-08-2012", "10256", "Broken tailight"}
	argslist[2] = []string{"0001", "Garage B", "Accident", "07-07-2012", "9600", "Dent on left side"}
	argslist[3] = []string{"0001", "Garage A", "Maintenance", "09-12-2012", "20012", "Oil Change"}
	argslist[4] = []string{"0001", "Garage C", "Repair", "08-10-2012", "15689", "Gearbox issue"}

	for _, args := range argslist {
		m.AddMaintenance(stub, args)
	}
	return "Load Maintencan samples: 5 inserted"
}
