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
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

const (
	apiHost            = "Magento API URL"
	magentoAccessToken = "Magento bearer token"
	verifyCodeToken    = "Magento basic token"
)

// Define the Smart Contract structure
type SmartContract struct {
}

type ProductResponse struct {
	Items  []Item `json:"items"`
	Number int    `json:"total_count"`
}

//Item is the structure of the products stored in API
type Item struct {
	ID            int     `json:"id"`
	Sku           string  `json:"sku"`
	Name          string  `json:"name"`
	AtributeSetID int     `json:"attribute_set_id"`
	Price         float32 `json:"price"`
	Status        int     `json:"status"`
	Visibility    int     `json:"visibility"`
	TypeID        string  `json:"type_id"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	// ProductLinks     string `json:"product_links"`
	// TierPrices       string `json:"tier_prices"`
	CustomAttributes []Attribute `json:"custom_attributes"`
}

//Attribute refers to the custom attributes of the product
type Attribute struct {
	AttributeCode string `json:"attribute_code"`
	Value         string `json:"value"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initProducts" {
		return s.initProducts(APIstub)
	} else if function == "createProduct" {
		return s.createProduct(APIstub, args)
	} else if function == "queryProduct" {
		return s.queryProduct(APIstub, args)
	} else if function == "queryAllProducts" {
		return s.queryAllProducts(APIstub)
	} else if function == "editProduct" {
		return s.editProduct(APIstub, args)
	} else if function == "deleteProduct" {
		return s.deleteProduct(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) initProducts(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Calling magento product api")
	request, _ := http.NewRequest(`GET`, apiHost+`products?searchCriteria={"sortOrders":[{"field":"id","direction":"asc"}]}`, nil)
	request.Header.Set("Authorization", magentoAccessToken)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return shim.Error(err.Error())
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// fmt.Println(string(data))
		var responseObject ProductResponse
		json.Unmarshal(data, &responseObject)
		for i := 0; i < len(responseObject.Items); i++ {
			fmt.Println(responseObject.Items[i].Sku)
			itemAsBytes, _ := json.Marshal(responseObject.Items[i])
			key := "P" + strconv.Itoa(responseObject.Items[i].ID)
			APIstub.PutState(key, itemAsBytes)
		}
	}
	return shim.Success(nil)
}

func (s *SmartContract) createProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var sku = args[0]
	sku = strings.Replace(sku, " ", "+", -1)
	request, _ := http.NewRequest(`GET`, apiHost+`products/`+sku, nil)
	request.Header.Set("Authorization", magentoAccessToken)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return shim.Error(err.Error())
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var item Item
		json.Unmarshal(data, &item)
		if item.Sku != args[0] {
			fmt.Println("Invalid SKU")
			return shim.Error("Invalid SKU")
		}
		// if len(item.Sku) == 0 {
		// 	return shim.Error("Invalid SKU")
		// }
		key := "P" + strconv.Itoa(item.ID)
		itemAsBytes, _ := json.Marshal(item)
		APIstub.PutState(key, itemAsBytes)
	}
	return shim.Success(nil)
}

func (s *SmartContract) queryProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := "P" + args[0]
	itemAsBytes, err := APIstub.GetState(key)
	if err != nil {
		fmt.Println(err)
		return shim.Error("Unable to query product")
	}
	fmt.Println(len(key))
	fmt.Println(len(itemAsBytes))
	if len(itemAsBytes) == 0 {
		fmt.Println("Invalid product ID")
		return shim.Error("Invalid product ID")
	}

	fmt.Println(string(itemAsBytes))

	return shim.Success(itemAsBytes)
}

func (s *SmartContract) queryAllProducts(APIstub shim.ChaincodeStubInterface) sc.Response {
	startKey := "P0"
	endKey := "P999999"

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
		buffer.WriteString(`{"Key":`)
		buffer.WriteString(`"`)
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString(`"`)

		buffer.WriteString(`, "Record":`)
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString(`}`)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString(`]`)

	fmt.Printf("- queryAllProducts:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) editProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	sku := args[0]
	sku = strings.Replace(sku, " ", "+", -1)
	request, _ := http.NewRequest(`GET`, apiHost+`products/`+sku, nil)
	request.Header.Set("Authorization", magentoAccessToken)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return shim.Error(err.Error())
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var item Item
		json.Unmarshal(data, &item)
		if item.Sku != args[0] {
			fmt.Println("Invalid SKU")
			return shim.Error("Invalid SKU")
		}
		key := "P" + strconv.Itoa(item.ID)
		itemAsBytes, _ := json.Marshal(item)
		APIstub.PutState(key, itemAsBytes)
	}
	return shim.Success(nil)
}

func (s *SmartContract) deleteProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	key := "P" + args[0]

	itemAsBytes, err := APIstub.GetState(key)
	if len(itemAsBytes) == 0 {
		fmt.Println("Invalid product ID")
		return shim.Error("Invalid product ID")
	}
	err = APIstub.DelState(key)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
