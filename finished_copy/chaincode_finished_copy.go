/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
type PatientTest struct {
	GeneralInfo   string `json:"generalInfo"`
	PersonalInfo  string `json:"personalInfo"`
	VarianEntries []struct {
		Date        string `json:"date"`
		MedicalData string `json:"medicalData"`
		VarianNode  string `json:"varianNode"`
	} `json:"varianEntries"`
}

const (
	columnAccountID   = "Account"
	columnContactInfo = "ContactInfo"
)

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	key := args[0]
	err := stub.CreateTable(key, []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: columnAccountID, Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: columnContactInfo, Type: shim.ColumnDefinition_STRING, Key: false},
	})
	//, {&shim.ColumnDefinition.Name: "date", &ColumnDefinition.Type: shim.ColumnDefinition_String, &ColumnDefinition.Key: true}) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}
func (t *SimpleChaincode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var key string
	var err error
	fmt.Println("running write()")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. name of the table,varNode and medData ")
	}
	key = args[0]
	a := args[1]
	b := args[2]

	ok, err := stub.InsertRow(key, shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: a}},
			&shim.Column{Value: &shim.Column_String_{String_: b}}},
	})
	if !ok && err == nil {

		return nil, errors.New("new error")
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var key string

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	accountID := args[1]

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: accountID}}
	columns = append(columns, col1)
	fmt.Println(stub.GetRow(key, columns))
	return nil, nil
}
