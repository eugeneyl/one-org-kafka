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
	// "bytes"
	"encoding/json"
	"fmt"
    "encoding/hex"
	"strconv"
	"hash/fnv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"github.com/ethereum/go-ethereum/crypto"
)

// Define the Smart Contract structure
type SmartContract struct {
}

type Token struct {
	TokenName 		string `json:tokenName`
	TokenSymbol 	string `json:tokenSymbol`
	TotalAmount 	string `json:totalAmount`
}

type Account struct {
	Address			string `json:address`
	PrivateKey 		string `json:privateKey`
	AccountValue 	string `json:accountValue`
}

type Approval struct {
	Allowance string `json:Allowance`
}

// type TransferEvent struct {
// 	AccountFrom			string `json:accountFrom`
// 	AccountTo 			string `json:accountTo`
// 	AmountTransfered	string `json:amount`	 
// }

// type ApprovalEvent struct {
// 	Owner			string `json:owner`
// 	Spender			string `json:spender`
// 	AmountApproved	string `json:amount`
// }

/*
 * Util function to generate wallet for user registration
 */
func generateKey() (string, string){
	key, _ := crypto.GenerateKey()
	address := crypto.PubkeyToAddress(key.PublicKey).Hex()
	privateKey := hex.EncodeToString(key.D.Bytes())
	return address, privateKey
}

func authenticateAccount(accountAsByte []byte, privateKey string) bool {
	account := Account{}
	err := json.Unmarshal(accountAsByte, &account)
	if err != nil {
		fmt.Println(err)
	}
	if account.PrivateKey == hash(privateKey) {
		return true
	} else {
		return false
	}	
}

func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return strconv.Itoa(int(h.Sum32()))
}


/*
 * The Init method is called when the Smart Contract "token" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "token"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	switch function {
	case "initToken":
		return s.initLedger(APIstub, args)
	case "queryTotalAmount":
		return s.queryTotalAmount(APIstub)
	case "queryTokenName":
		return s.queryTokenName(APIstub)
	case "queryTokenSymbol":
		return s.queryTokenSymbol(APIstub)
	case "queryReserve":
		return s.queryReserve(APIstub)
	case "createAccount":
		return s.createAccount(APIstub, args)
	case "balanceOf":
		return s.queryValue(APIstub, args)
	case "transfer":
		return s.transfer(APIstub, args)
	case "buyToken":
		newArgs := []string{"admin",args[1],args[0],args[2]}
		return s.transfer(APIstub, newArgs)
	case "sellToken":
		newArgs := []string{args[0],args[1],"admin",args[2]}
		return s.transfer(APIstub, newArgs)
	case "mintToken":
		return s.mintTokens(APIstub, args)
	case "burnToken":
		return s.burnTokens(APIstub, args)
	case "approve":
		return s.approve(APIstub, args)
	case "allowance":
		return s.queryAllowance(APIstub, args)
	case "increaseAllowance":
		return s.increaseAllowance(APIstub,  args)
	case "decreaseAllowance":
		return s.decreaseAllowance(APIstub, args)
	case "transferFrom":
		return s.transferFrom(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * Initledger method will create the token and create an admin who will be the owner of the coin. This
 * function can only be called once.
*/
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	checkBytes, _ := APIstub.GetState("token")
	if len(string(checkBytes)) > 0 {
		return shim.Error("Token already instantiated.")
	}

	//Initialise an admin account
	adminAdd, adminKey := generateKey()
	admin := Account{Address: adminAdd, PrivateKey: hash(adminKey), AccountValue:args[2] }
	adminAsBytes, _ := json.Marshal(admin)
	APIstub.PutState("admin", adminAsBytes)

	token := Token{TokenName: args[0], TokenSymbol: args[1], TotalAmount: args[2]}
	tokenAsBytes, _ := json.Marshal(token)
	APIstub.PutState("token", tokenAsBytes)
	fmt.Println(string(`{"Address":"` + adminAdd + `","PrivateKey":"` + adminKey + `"}`))

	return shim.Success([]byte(`{"Address":"` + adminAdd + `","PrivateKey":"` + adminKey + `"}`))
}

func (s *SmartContract) queryTotalAmount(APIstub shim.ChaincodeStubInterface) sc.Response {

	tokenAsBytes, _ := APIstub.GetState("token")
	token := Token{}
	err := json.Unmarshal(tokenAsBytes, &token)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(token.TotalAmount)	
	return shim.Success([]byte(token.TotalAmount))
}

func (s *SmartContract) queryTokenName(APIstub shim.ChaincodeStubInterface) sc.Response {

	tokenAsBytes, _ := APIstub.GetState("token")
	token := Token{}
	err := json.Unmarshal(tokenAsBytes, &token)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(token.TokenName)	
	return shim.Success([]byte(token.TokenName))
}

func (s *SmartContract) queryTokenSymbol(APIstub shim.ChaincodeStubInterface) sc.Response {

	tokenAsBytes, _ := APIstub.GetState("token")
	token := Token{}
	err := json.Unmarshal(tokenAsBytes, &token)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(token.TokenSymbol)	
	return shim.Success([]byte(token.TokenSymbol))
}

func (s *SmartContract) queryReserve(APIstub shim.ChaincodeStubInterface) sc.Response {

	adminAsBytes, _ := APIstub.GetState("admin")
	admin := Account{}
	err := json.Unmarshal(adminAsBytes, &admin)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(admin.AccountValue)	
	return shim.Success([]byte(admin.AccountValue))
}

func (s *SmartContract) createAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	adminAsBytes, _ := APIstub.GetState("admin")
	if !authenticateAccount(adminAsBytes, args[0]) {
		return shim.Error("Please use the key of an Admin to create an account")
	}

	//Reduce remaining amount
	admin := Account{}
	err := json.Unmarshal(adminAsBytes, &admin)
	if err != nil {
		fmt.Println(err)
	}
	reserve, _ := strconv.Atoi(admin.AccountValue)
	accValue, _ := strconv.Atoi(args[1])
	if reserve < accValue {
		return shim.Error("Insufficient funds in reserve, cannot create new account")
	}
	admin.AccountValue = strconv.Itoa(reserve-accValue)
	adminAsBytes,_ = json.Marshal(admin)
	APIstub.PutState("admin", adminAsBytes)

	address, privateKey := generateKey()
	account := Account{Address: address, PrivateKey: hash(privateKey), AccountValue: args[1]}
	accountAsBytes, _ := json.Marshal(account)
	APIstub.PutState(address, accountAsBytes)

	fmt.Println("New account made")
	fmt.Println("Address: " + address)
	fmt.Println("Private Key: " + privateKey)

	return shim.Success([]byte(`{"Address":"` + address + `","PrivateKey":"` + privateKey + `"}`))
}

func (s *SmartContract) queryValue(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	//TODO:Check for invalid account
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	accountAsBytes, _ := APIstub.GetState(args[0])
	target := Account{}
	if len(accountAsByte == 0) {
		return shim.Error("Account not found")
	}
	err := json.Unmarshal(accountAsBytes, &target)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(target.AccountValue)
	return shim.Success([]byte(target.AccountValue))
}

func (s *SmartContract) transfer(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	//TODO:Check for invalid account
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	senderAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Unable to get sender account")
	}

	if len(senderAsByte == 0) {
		return shim.Error("Sender account not found")
	}

	if !authenticateAccount(senderAsBytes, args[1]) {
		return shim.Error("Invalid private key")
	}

	receiverAsBytes, err := APIstub.GetState(args[2])
	if err != nil {
		return shim.Error("Unable to get receiver account")
	}

	if len(receiverAsByte == 0) {
		return shim.Error("Receiver account not found")
	}

	sender := Account{}
	receiver := Account{}
	transferAmt,_ := strconv.Atoi(args[3])
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
	APIstub.PutState(args[2], receiverAsBytes)

	fmt.Println("Transfer success. Sender new balance: " + sender.AccountValue + 
	", Receiver new balance: " + receiver.AccountValue)
	return shim.Success(nil)
}

func (s *SmartContract) mintTokens(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	adminAsBytes, _ := APIstub.GetState("admin")
	if !authenticateAccount(adminAsBytes, args[0]) {
		return shim.Error("Please use the key of an Admin")
	}
	admin := Account{}
	json.Unmarshal(adminAsBytes, &admin)
	oldValue, _ := strconv.Atoi(admin.AccountValue)
	addOn, _ := strconv.Atoi(args[1])
	admin.AccountValue = strconv.Itoa(oldValue + addOn)
	adminAsBytes, _ = json.Marshal(admin)
	APIstub.PutState("admin", adminAsBytes)
	tokenAsBytes, _ := APIstub.GetState("token")
	token := Token{}
	err := json.Unmarshal(tokenAsBytes, &token)
	if err != nil {
		fmt.Println(err)
	}
	oldValue, _ = strconv.Atoi(token.TotalAmount)
	token.TotalAmount = strconv.Itoa(oldValue + addOn)
	tokenAsBytes, _ = json.Marshal(token)
	APIstub.PutState("token", tokenAsBytes)
	fmt.Println("New Total Token: " + token.TotalAmount)
	return shim.Success(nil)
}

func (s *SmartContract) burnTokens(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	adminAsBytes, _ := APIstub.GetState("admin")
	if !authenticateAccount(adminAsBytes, args[0]) {
		return shim.Error("Please use the key of an Admin")
	}
	admin := Account{}
	json.Unmarshal(adminAsBytes, &admin)
	oldValue, _ := strconv.Atoi(admin.AccountValue)
	addOn, _ := strconv.Atoi(args[1])
	newValue := oldValue - addOn
	if newValue < 0 {
		return shim.Error("Insufficient tokens to burn.")
	}
	admin.AccountValue = strconv.Itoa(oldValue - addOn)
	adminAsBytes, _ = json.Marshal(admin)
	APIstub.PutState("admin", adminAsBytes)
	tokenAsBytes, _ := APIstub.GetState("token")
	token := Token{}
	err := json.Unmarshal(tokenAsBytes, &token)
	if err != nil {
		fmt.Println(err)
	}
	oldValue, _ = strconv.Atoi(token.TotalAmount)
	token.TotalAmount = strconv.Itoa(oldValue - addOn)
	tokenAsBytes, _ = json.Marshal(token)
	APIstub.PutState("token", tokenAsBytes)
	fmt.Println("New Total Token: " + token.TotalAmount)
	return shim.Success(nil)
}

func (s *SmartContract) approve(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	ownerAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Unable to get sender account")
	}

	if len(ownerAsBytes == 0) {
		return shim.Error("Owner account not found")
	}

	if !authenticateAccount(ownerAsBytes, args[1]) {
		return shim.Error("Invalid private key")
	}

	senderAsBytes, err := APIstub.GetState(args[2])
	if err != nil {
		return shim.Error("Unable to get sender account")
	}

	if len(senderAsBytes == 0) {
		return shim.Error("Sender account not found")
	}
	approval := Approval{Allowance: args[3]}
	approvalAsBytes, _ := json.Marshal(approval)
	APIstub.PutState(args[0]+"-"+args[2], approvalAsBytes)
	fmt.Println(string(approvalAsBytes))

	return shim.Success(nil)
}

func (s *SmartContract) queryAllowance(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	key := args[0] + "-" + args[1]
	approvalAsBytes, _ := APIstub.GetState(key)
	approval := Approval{}
	err := json.Unmarshal(approvalAsBytes, &approval)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(approval.Allowance)
	return shim.Success([]byte(approval.Allowance))
}

func (s *SmartContract) increaseAllowance(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	ownerAsBytes, _ := APIstub.GetState(args[0])
	if !authenticateAccount(ownerAsBytes, args[1]) {
		return shim.Error("Invalid private key")
	}

	key := args[0] + "-" + args[2]
	approvalAsBytes, _ := APIstub.GetState(key)
	if len(string(approvalAsBytes)) == 0 {
		return shim.Error("Allowance does not exist.")
	}

	approval := Approval{}
	err := json.Unmarshal(approvalAsBytes, &approval)
	if err != nil {
		fmt.Println(err)
	}
	oldValue, _ := strconv.Atoi(approval.Allowance)
	change, _ := strconv.Atoi(args[3])
	approval.Allowance = strconv.Itoa(oldValue + change)
	approvalAsBytes, _ = json.Marshal(approval)
	APIstub.PutState(key, approvalAsBytes)


	fmt.Println("New allowance: " + approval.Allowance)
	return shim.Success([]byte(approval.Allowance))
}

func (s *SmartContract) decreaseAllowance(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	ownerAsBytes, _ := APIstub.GetState(args[0])
	if !authenticateAccount(ownerAsBytes, args[1]) {
		return shim.Error("Invalid private key")
	}

	key := args[0] + "-" + args[2]
	approvalAsBytes, _ := APIstub.GetState(key)
	if len(string(approvalAsBytes)) == 0 {
		return shim.Error("Allowance does not exist.")
	}

	approval := Approval{}
	err := json.Unmarshal(approvalAsBytes, &approval)
	if err != nil {
		fmt.Println(err)
	}
	oldValue, _ := strconv.Atoi(approval.Allowance)
	change, _ := strconv.Atoi(args[3])
	if oldValue - change > 0 {
		approval.Allowance = strconv.Itoa(oldValue - change)
	} else {
		approval.Allowance = strconv.Itoa(0)
	}
	approvalAsBytes, _ = json.Marshal(approval)
	APIstub.PutState(key, approvalAsBytes)


	fmt.Println("New allowance: " + approval.Allowance)
	return shim.Success([]byte(approval.Allowance))
}

func (s *SmartContract) transferFrom(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	spenderAsBytes, _ := APIstub.GetState(args[0])
	if !authenticateAccount(spenderAsBytes, args[1]) {
		return shim.Error("Invalid private key")
	}

	key := args[2] + "-" + args[0]
	approvalAsBytes, _ := APIstub.GetState(key)
	if len(string(approvalAsBytes)) == 0 {
		return shim.Error("Allowance does not exist.")
	}
	
	approval := Approval{}
	err := json.Unmarshal(approvalAsBytes, &approval)
	if err != nil {
		fmt.Println(err)
	}
	oldValue, _ := strconv.Atoi(approval.Allowance)
	transferAmt,_ := strconv.Atoi(args[4])
	if oldValue < transferAmt {
		return shim.Error("Not enough allowance to complete transaction")
	} 
	approval.Allowance = strconv.Itoa(oldValue - transferAmt)
	approvalAsBytes, _ = json.Marshal(approval)
	APIstub.PutState(key, approvalAsBytes)


	senderAsBytes, err := APIstub.GetState(args[2])
	if err != nil {
		return shim.Error("Unable to get sender account")
	}

	receiverAsBytes, err := APIstub.GetState(args[3])
	if err != nil {
		return shim.Error("Unable to get receiver account")
	}

	sender := Account{}
	receiver := Account{}
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

	receiverAmt = receiverAmt + transferAmt
	receiver.AccountValue = strconv.Itoa(receiverAmt)
	receiverAsBytes,_ = json.Marshal(receiver)

	APIstub.PutState(args[0], senderAsBytes)
	APIstub.PutState(args[2], receiverAsBytes)
	APIstub.PutState(key, approvalAsBytes)

	fmt.Println(string(senderAsBytes))
	fmt.Println(string(receiverAsBytes))
	fmt.Println(string(approvalAsBytes))

	return shim.Success([]byte(approval.Allowance))
}


// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
