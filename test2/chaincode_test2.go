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

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	valueInit := args[1]

	fmt.Println("encodedPatient", valueInit)
	err := stub.PutState(args[0], []byte(valueInit))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "writeNew" {
		exist, errwriteNew := t.read(stub, args)
		if errwriteNew != nil {
			return nil, errwriteNew
		}
		if exist != nil {
			return nil, nil
		} else {
			return t.writeNew(stub, args)
		}

	} else if function == "addData" {
		exist, errAddData := t.read(stub, args)
		if errAddData != nil {
			return nil, errAddData
		}
		return t.addData(stub, exist, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// writeNew - invoke function to write new patient
func (t *SimpleChaincode) addData(stub shim.ChaincodeStubInterface, prev []byte, args []string) ([]byte, error) {

	var key, valueUpd, toAppend string
	var err error
	fmt.Println("running addNew()")
	key = args[0] //patientID
	fmt.Printf("adding data")
	toAppend = args[1]
	prev = append(prev, toAppend...)
	fmt.Printf("encodedPatient", valueUpd)
	err = stub.PutState(key, []byte(prev)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) writeNew(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var key, val string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //patientID
	val = args[1]

	err = stub.PutState(key, []byte(val)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	key = args[0]
	valAsBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsBytes, nil
}
