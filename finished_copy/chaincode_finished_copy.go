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
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

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

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func EncryptAESCFB(dst, src, key, iv []byte) error {
	aesBlockEncrypter, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(dst, src)
	return nil
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	data := []byte(`
	    {
		"generalInfo": "PID",
		"personalInfo" : "somePersonalEncrypted",
		"varianEntries":
		[
		{
		"varianNode": "var1",
		"date": "2006-01-02T15:04:05" ,
		"medicalData": "patient is healthy"	}
		]
	}
	`)
	var pt PatientTest
	errUnm := json.Unmarshal(data, &pt)
	if errUnm != nil {
		fmt.Printf("Error: %s", errUnm)
	}
	patient := &pt

	//	fmt.Printf("generalInfo: %s", patient.GeneralInfo)

	encBuf := new(bytes.Buffer)
	errrNewEnc := gob.NewEncoder(encBuf).Encode(patient)
	if errrNewEnc != nil {
		log.Fatal(errrNewEnc)
	}

	valueInit := encBuf.Bytes()

	fmt.Println("encodedPatient", valueInit)
	err := stub.PutState(args[0], valueInit)
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
	} else if function == "writeNew" {
		return t.writeNew(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// writeNew - invoke function to write new patient
func (t *SimpleChaincode) writeNew(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	data := []byte(`
			{
		"generalInfo": "PID",
		"personalInfo" : "somePersonalEncrypted",
		"varianEntries":
		[
		{
		"varianNode": "var1",
		"date": "2006-01-02T15:04:05" ,
		"medicalData": "patient is healthy"	}
		]
	}
	`)
	var key string
	var err error
	fmt.Println("running write()")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. name of the key and value to set")
	}

	key = args[0] //patientID

	fmt.Println("no record found")
	var newpt PatientTest
	errUnm2 := json.Unmarshal(data, &newpt)
	if errUnm2 != nil {
		fmt.Printf("Error: %s", errUnm2)
	}
	patientP := &newpt
	varEntry := &patientP.VarianEntries[0]
	patientP.GeneralInfo = "PID"
	patientP.PersonalInfo = "empty"
	varEntry.VarianNode = args[1]
	varEntry.Date = time.Now().String()
	varEntry.MedicalData = args[2]
	encBuf := new(bytes.Buffer)
	errrNewEnc := gob.NewEncoder(encBuf).Encode(patientP)
	if errrNewEnc != nil {
		log.Fatal(errrNewEnc)
	}

	valueInit := encBuf.Bytes()

	fmt.Println("encodedPatient", valueInit)
	err = stub.PutState(key, valueInit) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	// encode
	///valueBefCod = args[1]
	///valueBefEnc := base64.StdEncoding.EncodeToString([]byte(valueBefCod))

	//	err = stub.PutState(key, []byte(valueBefEnc)) //write the variable into the chaincode state
	//if err != nil {
	//	return nil, err
	//}

	//encryption
	///const key16 = "1234567890123456"
	//	const key161 = "6543210987654321"
	//const key24 = "123456789012345678901234"
	//const key32 = "12345678901234567890123456789012"
	///var keyforAES = key16
	//	var msg = "message"
	///var iv = []byte(keyforAES)[:aes.BlockSize] // Using IV same as key is probably bad
	///var errr error
	///fmt.Printf("!Encrypting %v %v  -> %v\n", keyforAES, []byte(iv), valueBefEnc)
	// Encrypt
	///value := make([]byte, len(valueBefEnc))
	///errr = EncryptAESCFB(value, []byte(valueBefEnc), []byte(keyforAES), iv)
	///fmt.Printf("Encrypting %v %v %s -> %v\n", keyforAES, []byte(iv), valueBefEnc, value)
	///if errr != nil {
	///panic(errr)
	///}

	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsBytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsBytes, nil
}
