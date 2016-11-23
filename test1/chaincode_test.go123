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
// create account->add private data->add permissions->addnewDataItem-> read
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"time"
	"unicode/utf8"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type MyDataItem struct {
	DoctorId   string `json:"doctorId"`
	DataLink   string `json:"dataLink"`
	TimestampD string `json:"timestampD"`
	HashData   uint32 `json:"hashData"`
}

type MyPermission struct {
	DoctorId       string `json:"doctorId"`
	Right          string `json:"roght"`
	TimestampP     string `json:"TimestampP"`
	HashPermission uint32 `json:"hashPermission"`
}

type MyRecord struct {
	PatientId   string         `json:"patientId"`
	PrivateData string         `json:"privateData"`
	DataItems   []MyDataItem   `json:"dataItems"`
	Permissions []MyPermission `json:"permissions"`
}

// verify if there a permission *permission* given to a *doctorId*
func (rec *MyRecord) VerifyPermission(doctorId string, permission string) bool {
	for i := 0; i < len(rec.DataItems); i++ {
		if rec.Permissions[i].DoctorId == doctorId {
			if rec.Permissions[i].Right == permission {
				fmt.Println(rec.Permissions[i].Right)
				return true
			} else {
				fmt.Println("access denied: not enough rights")
			}
		} else {
			fmt.Println("access denied: no permissions found for ", doctorId)
		}
	}
	return false
}

func (rec *MyRecord) AddDataItem(dataitem MyDataItem) []MyDataItem {
	rec.DataItems = append(rec.DataItems, dataitem)
	return rec.DataItems
}

func (rec *MyRecord) AddPermission(permission MyPermission) []MyPermission {
	//check it there is a permision for this doc
	//then if yes, if it the same
	//then later add data period
	rec.Permissions = append(rec.Permissions, permission)
	return rec.Permissions
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

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
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	fmt.Println("chaincode initialized")
	err := stub.PutState(args[0], []byte("initialized"))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "createPatient" {
		exist, errcreatePatient := t.read(stub, args)
		if errcreatePatient != nil {
			return nil, errcreatePatient
		}
		if exist != nil { //there is patient with this id
			return nil, nil
		} else {
			return t.createPatient(stub, args)
		}

	} else if function == "addPermission" {
		exist, erraddPermission := t.read(stub, args)
		if erraddPermission != nil {
			return nil, erraddPermission
		}
		exPtn, size := utf8.DecodeRune(exist)
		fmt.Println(size)
		var pnt MyRecord
		err := json.Unmarshal([]byte(string(exPtn)), &pnt)
		if err != nil {
			return nil, err
		}
		return t.addPermission(stub, pnt, args)
	} else if function == "addDataItem" {
		exist, erraddDataItem := t.read(stub, args)
		if erraddDataItem != nil {
			return nil, erraddDataItem
		}
		exPtn, size := utf8.DecodeRune(exist)
		fmt.Println(size)
		var pnt MyRecord
		err := json.Unmarshal([]byte(string(exPtn)), &pnt)
		if err != nil {
			return nil, err
		}
		if pnt.VerifyPermission(args[1], "write") != true {
			fmt.Println("error...")
		}
		return t.addDataItem(stub, pnt, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {
		exist, errRead := t.read(stub, args)
		if errRead != nil {
			return nil, errRead
		}
		//exPtn, size:= utf8.DecodeRune(exist)
		//var pnt MyRecord
		//err := json.Unmarshal(exPtn, &pnt)
		//if pnt.VerifyPermission(*doctorID!!!*, "read")!=true{
		//	fmt.Println("error...")
		//}
		return exist, errRead
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) createPatient(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	itemsData := []MyDataItem{}
	itemsPerm := []MyPermission{}
	pid := args[1]
	prdata := "someData"

	patient := MyRecord{pid, prdata, itemsData, itemsPerm}
	initVal, errMarsh := json.Marshal(patient)
	if errMarsh != nil {
		return nil, errMarsh
	}

	fmt.Println("patient with the following id is created:", pid)
	errinit := stub.PutState(pid, []byte(initVal))
	if errinit != nil {
		return nil, errinit
	}

	return nil, nil
}

// addPermission - invoke to add data (by doctor)
func (t *SimpleChaincode) addPermission(stub *shim.ChaincodeStub, pnt MyRecord, args []string) ([]byte, error) {

	item1 := MyPermission{DoctorId: args[1], Right: args[2], TimestampP: time.Now().String(), HashPermission: hash(args[2])}
	var key string
	fmt.Println("running addPermission()")
	key = args[0] //patientID
	pnt.AddPermission(item1)
	addedPrm, errMarsh := json.Marshal(pnt)
	if errMarsh != nil {
		return nil, errMarsh
	}

	fmt.Println("permission is added")
	errPerm := stub.PutState(key, []byte(addedPrm))
	if errPerm != nil {
		return nil, errPerm
	}

	return nil, nil
}

// addDataItem - invoke to add data (by doctor)
func (t *SimpleChaincode) addDataItem(stub *shim.ChaincodeStub, pnt MyRecord, args []string) ([]byte, error) {

	item1 := MyDataItem{DoctorId: args[1], DataLink: args[2], TimestampD: time.Now().String(), HashData: hash(args[2])}
	var key string
	var err error
	fmt.Println("running addDataItem()")
	key = args[0] //patientID
	pnt.AddDataItem(item1)
	addedVal, errMarsh := json.Marshal(pnt)
	if errMarsh != nil {
		return nil, errMarsh
	}

	fmt.Println("dataItem is added")
	err = stub.PutState(key, []byte(addedVal))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
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
