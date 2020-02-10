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

type OrderResponse struct {
	Orders []Order `json:"items"`
	Number int     `json:"total_count"`
}

type Order struct {
	ObjectType         string          `json:docType`
	EntityID           int             `json:"entity_id"`
	CreatedAt          string          `json:"created_at"`
	CustomerEmail      string          `json:"customer_email"`
	CustomerFirstname  string          `json:"customer_firstname"`
	CustomerID         int             `json:"customer_id"`
	CustomerLastname   string          `json:"customer_lastname"`
	GlobalCurrencyCode string          `json:"global_currency_code"`
	GrandTotal         float32         `json:"grand_total"`
	IncrementId        string          `json:"increment_id"`
	ShippingAmount     float32         `json:"shipping_amount"`
	State              string          `json:"state"`
	Status             string          `json:"status"`
	Subtotal           float32         `json:"subtotal"`
	totalorderCount    int             `json:"total_order_count"`
	UpdatedAt          string          `json:"updated_at"`
	orders             []order         `json:"orders"`
	BillingAddress     Address         `json:"billing_address"`
	Payment            Payment         `json:"payment"`
	StatusHistories    []StatusHistory `json:"status_histories"`
}

type order struct {
	CreatedAt   string  `json:"created_at"`
	Description string  `json:"description"`
	orderID     int     `json:"order_id"`
	Name        string  `json:"name"`
	Price       float32 `json:"price"`
	ProductID   int     `json:"product_id"`
	QtyInvoiced int     `json:"qty_invoiced"`
	QntOrdered  int     `json:"qty_ordered"`
	Sku         string  `json:"sku"`
	StoreId     int     `json:"store_id"`
	UpdatedAt   string  `json:"updated_at"`
}

type Address struct {
	AddressType string   `json:"address_type"`
	City        string   `json:"city"`
	Company     string   `json:"company"`
	CountryID   string   `json:"country_id"`
	Email       string   `json:"email"`
	EntityID    int      `json:"entity_id"`
	Firstname   string   `json:"firstname"`
	Lastname    string   `json:"lastname"`
	Postcode    string   `json:"postcode"`
	Region      string   `json:"region"`
	RegionCode  string   `json:"region_code"`
	Street      []string `json:"street"`
	Telephone   string   `json:"telephone"`
}

type Payment struct {
	AccountStatus         string   `json:"account_status"`
	AdditionalInformation []string `json:"additional_information"`
	AmountOrdered         float32  `json:"amount_ordered"`
	AmountPaid            float32  `json:"amount_paid"`
	EntityID              int      `json:"entity_id"`
	Method                string   `json:"method"`
	ShippingAmount        float32  `json:"shipping_amount"`
}

type StatusHistory struct {
	Comment    string `json:"comment"`
	CreatedAt  string `json:"created_at"`
	EntityID   int    `json:"entity_id"`
	EntityName string `json:"entity_name"`
	Status     string `json:"status"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initOrders" {
		return s.initOrders(APIstub)
	} else if function == "queryOrder" {
		return s.queryOrder(APIstub, args)
	} else if function == "queryAllOrders" {
		return s.queryAllOrders(APIstub)
	} else if function == "createOrder" {
		return s.createOrder(APIstub, args)
	} else if function == "editOrder" {
		return s.editOrder(APIstub, args)
	} else if function == "deleteOrder" {
		return s.deleteOrder(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) initOrders(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Calling magento orders api")

	//For testing, only take in the latest 100 orders due to server restriction.
	//Todo: Load all orders into ledger
	request, _ := http.NewRequest(`GET`, apiHost+`orders?searchCriteria[pageSize]=100&searchCriteria[currentPage]=1&searchCriteria[sortOrders][0][field]=entity_id&searchCriteria[sortOrders][0][direction]=desc`, nil)
	request.Header.Set("Authorization", magentoAccessToken)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return shim.Error(err.Error())
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// fmt.Println(string(data))
		var responseObject OrderResponse
		json.Unmarshal(data, &responseObject)
		for i := 0; i < len(responseObject.Orders); i++ {
			fmt.Println(responseObject.Orders[i].EntityID)
			responseObject.Orders[i].ObjectType = "orders"
			orderAsBytes, _ := json.Marshal(responseObject.Orders[i])
			key := "O" + strconv.Itoa(responseObject.Orders[i].EntityID)
			APIstub.PutState(key, orderAsBytes)
		}
	}
	return shim.Success(nil)
}

func (s *SmartContract) queryOrder(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := "O" + args[0]
	orderAsBytes, err := APIstub.GetState(key)
	if err != nil {
		fmt.Println(err)
		return shim.Error("Unable to query order")
	}
	if len(orderAsBytes) == 0 {
		fmt.Println("Invalid entity ID")
		return shim.Error("Invalid entity ID")
	}

	fmt.Println(string(orderAsBytes))

	return shim.Success(orderAsBytes)
}

func (s *SmartContract) queryAllOrders(APIstub shim.ChaincodeStubInterface) sc.Response {
	startKey := "O0"
	endKey := "O999999"

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

	fmt.Printf("- queryAllOrders:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) createOrder(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var entityID = args[0]
	request, _ := http.NewRequest(`GET`, apiHost+`orders/`+entityID, nil)
	request.Header.Set("Authorization", magentoAccessToken)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return shim.Error(err.Error())
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var order Order
		json.Unmarshal(data, &order)
		if strconv.Itoa(order.EntityID) != entityID {
			fmt.Println("Invalid order ID")
			return shim.Error("Invalid order ID")
		}
		order.ObjectType = "orders"
		key := "O" + strconv.Itoa(order.EntityID)
		orderAsBytes, _ := json.Marshal(order)
		APIstub.PutState(key, orderAsBytes)
	}
	return shim.Success(nil)
}

func (s *SmartContract) editOrder(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var entityID = args[0]
	key := "O" + entityID
	tmpAsBytes, err := APIstub.GetState(key)
	if err != nil {
		fmt.Println(err)
		return shim.Error("Invalid order ID")
	}
	if len(tmpAsBytes) == 0 {
		return shim.Error("Order ID not found")
	}
	request, _ := http.NewRequest(`GET`, apiHost+`orders/`+entityID, nil)
	request.Header.Set("Authorization", magentoAccessToken)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return shim.Error(err.Error())
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var order Order
		json.Unmarshal(data, &order)
		if strconv.Itoa(order.EntityID) != entityID {
			fmt.Println("Invalid order ID")
			return shim.Error("Invalid order ID")
		}
		order.ObjectType = "orders"
		orderAsBytes, _ := json.Marshal(order)
		APIstub.PutState(key, orderAsBytes)
	}
	return shim.Success(nil)
}

func (s *SmartContract) deleteOrder(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	key := "O" + args[0]

	orderAsBytes, err := APIstub.GetState(key)
	if len(orderAsBytes) == 0 {
		fmt.Println("Invalid order ID")
		return shim.Error("Invalid order ID")
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
