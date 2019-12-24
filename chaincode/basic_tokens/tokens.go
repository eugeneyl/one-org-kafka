/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
type Account struct {
	AccountId    string `json:accountId`
	AccountValue string `json:accountValue`
}

func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return strconv.Itoa(int(h.Sum32()))
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryValueByID" {
		return s.queryValueByID(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "queryValue" {
		return s.queryValue(APIstub, args)
	} else if function == "createAccount" {
		return s.createAccount(APIstub, args)
	} else if function == "transferFrom" {
		return s.transferFrom(APIstub, args)
	} else if function == "queryAllAccounts" {
		return s.queryAllAccounts(APIstub)
	} else if function == "queryTotalAmount" {
		return s.queryTotalAmount(APIstub)
	} else if function == "addTokens" {
		return s.addTokens(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryValueByID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	fmt.Println(hash(args[0]))
	accountAsBytes, _ := APIstub.GetState(hash(args[0]))
	fmt.Println(string(accountAsBytes))
	var target Account
	err := json.Unmarshal(accountAsBytes, &target)
	if err != nil {
		fmt.Println(err)
	}
	return shim.Success([]byte("Remaining Balance: " + target.AccountValue))
}

func (s *SmartContract) queryValue(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	fmt.Println(hash(args[0]))
	accountAsBytes, _ := APIstub.GetState(args[0])
	fmt.Println(string(accountAsBytes))
	var target Account
	err := json.Unmarshal(accountAsBytes, &target)
	if err != nil {
		fmt.Println(err)
	}
	return shim.Success([]byte("Remaining Balance: " + target.AccountValue))
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {

	accounts := []Account{
		Account{AccountId: "peer1_org1_master", AccountValue: "1000"},
		Account{AccountId: "peer2_org1_master", AccountValue: "1000"},
		Account{AccountId: "peer3_org1_master", AccountValue: "1000"},
	}

	i := 0
	for i < len(accounts) {
		fmt.Println("i is ", i)
		accountAsBytes, _ := json.Marshal(accounts[i])
		APIstub.PutState(hash(accounts[i].AccountId), accountAsBytes)
		// fmt.Println(hash(accounts[i].AccountId))
		// fmt.Println(accountAsBytes)
		fmt.Println("Added", accounts[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var account = Account{AccountId: args[0], AccountValue: args[1]}

	accountAsBytes, _ := json.Marshal(account)
	APIstub.PutState(hash(args[0]), accountAsBytes)

	fmt.Println("New account made")
	fmt.Println("Address: " + hash(args[0]))
	fmt.Println(string(accountAsBytes))

	return shim.Success([]byte(hash(args[0])))
}

func (s *SmartContract) transferFrom(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	senderAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Unable to get sender account")
	}

	receiverAsBytes, err := APIstub.GetState(args[1])
	if err != nil {
		return shim.Error("Unable to get receiver account")
	}

	fmt.Println("Able to transfer.")

	sender := Account{}
	receiver := Account{}
	transferAmt,_ := strconv.Atoi(args[2])
	json.Unmarshal(senderAsBytes, &sender)
	json.Unmarshal(receiverAsBytes, &receiver)
	senderAmt,_ := strconv.Atoi(sender.AccountValue)
	receiverAmt,_ := strconv.Atoi(receiver.AccountValue)
	if senderAmt < transferAmt {
		return shim.Error("Insufficient amount for transfer.")
	}

	senderAmt = senderAmt - transferAmt
	sender.AccountValue = strconv.Itoa(senderAmt)
	senderAsBytes,_ = json.Marshal(sender)
	APIstub.PutState(args[0], senderAsBytes)
	
	receiverAmt = receiverAmt + transferAmt
	receiver.AccountValue = strconv.Itoa(receiverAmt)
	receiverAsBytes,_ = json.Marshal(receiver)
	APIstub.PutState(args[1], receiverAsBytes)

	// orderAsBytes, _ := APIstub.GetState(args[0])
	// order := Order{}

	// json.Unmarshal(orderAsBytes, &order)
	// order.Status = args[1]

	// orderAsBytes, _ = json.Marshal(order)
	// APIstub.PutState(args[0], orderAsBytes)
	fmt.Println("Transfer success. Sender new balance: " + sender.AccountValue + 
	", Receiver new balance: " + receiver.AccountValue)

	return shim.Success([]byte("Transfer success. Sender new balance: " + sender.AccountValue + 
		", Receiver new balance: " + receiver.AccountValue))
}



func (s *SmartContract) queryAllAccounts(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "9999999999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		// fmt.Println(queryResponse.Value)
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllAccounts:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryTotalAmount(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "9999999999"
	totalAmount := 0

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		accountAsBytes := queryResponse.Value
		tempAcc := Account{}
		err1 := json.Unmarshal(accountAsBytes, &tempAcc)
		if err1 != nil {
			fmt.Println(err)
		}
		tempValue,_ := strconv.Atoi(tempAcc.AccountValue)
		totalAmount += tempValue
		fmt.Println(strconv.Itoa(totalAmount))
	}
	return shim.Success([]byte("Total Amount: " + strconv.Itoa(totalAmount)))
}

func (s *SmartContract) addTokens(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	accountAsBytes, _ := APIstub.GetState(args[0])
	account := Account{}
	json.Unmarshal(accountAsBytes, &account)
	oldValue, _ := strconv.Atoi(account.AccountValue)
	addOn, _ := strconv.Atoi(args[1])
	account.AccountValue = strconv.Itoa(oldValue + addOn)
	accountAsBytes, _ = json.Marshal(account)
	APIstub.PutState(args[0], accountAsBytes)
	return shim.Success([]byte("New value = " + account.AccountValue))

}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
