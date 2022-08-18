/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
)

// Chaincode stores a value
type Chaincode struct {
}

type Metadata struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Owner       string `json:"owner"`
	TotalSupply uint64 `json:"totalSupply"`
}

// Init is called when the chaincode is instantiated by the blockchain network.
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) sc.Response {
	_, params := stub.GetFunctionAndParameters()
	fmt.Println("Init called with params: ", params)
	if len(params) != 4 {
		return shim.Error("incorrect number of parameter")
	}

	tokenName, symbol, owner, amount := params[0], params[1], params[2], params[3]

	// check amount is unsigned int
	amountUint, err := strconv.ParseUint(string(amount), 10, 64)
	if err != nil {
		return shim.Error("amount must be a number or amount cannot be negative")
	}

	// tokenName & symbol & owner cannot be empty
	if len(tokenName) == 0 || len(symbol) == 0 || len(owner) == 0 {
		return shim.Error("tokenName or symbol or owner cannot be emtpy")
	}

	// make metadata
	metadata := &Metadata{Name: tokenName, Symbol: symbol, Owner: owner, TotalSupply: amountUint}
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return shim.Error("failed to Marshal erc20, error: " + err.Error())
	}

	// save token meta data
	err = stub.PutState(tokenName, metadataBytes)
	if err != nil {
		return shim.Error("failed to PutState, error: " + err.Error())
	}

	// save owner balance
	err = stub.PutState(owner, []byte(amount))
	if err != nil {
		return shim.Error("failed to PutState, error: " + err.Error())
	}

	return shim.Success(nil)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("Invoke()", fcn, params)

	switch fcn {
	case "totalSupply":
		return cc.totalSupply(stub, params)
	case "balanceOf":
		return cc.balanceOf(stub, params)
	case "transfer":
		return cc.transfer(stub, params)
	case "allowance":
		return cc.allowance(stub, params)
	case "approve":
		return cc.approve(stub, params)
	case "transferFrom":
		return cc.transferFrom(stub, params)
	case "increaseAllowance":
		return cc.increaseAllowance(stub, params)
	case "decreaseAllowance":
		return cc.decreaseAllowance(stub, params)
	case "mint":
		return cc.mint(stub, params)
	case "burn":
		return cc.burn(stub, params)
	default:
		return sc.Response{Status: 404, Message: "404 Not Found", Payload: nil}
	}
}

func (cc *Chaincode) totalSupply(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	// a number of params must be one
	if len(params) != 1 {
		return shim.Error("incorrect number of parameter")
	}

	tokenName := params[0]

	// Get ERC20 Metadata
	metadata := Metadata{}
	metadataBytes, err := stub.GetState(tokenName)
	if err != nil {
		return shim.Error("failed to GetState, error: " + err.Error())
	}
	err = json.Unmarshal(metadataBytes, &metadata)
	if err != nil {
		return shim.Error("failed to Unmarshal, error: " + err.Error())
	}

	// Convert TotalSupply to Bytes
	totalSupplyBytes, err := json.Marshal(metadata.TotalSupply)
	if err != nil {
		return shim.Error("failed to Marshal totalSupply, error: " + err.Error())
	}
	fmt.Println(tokenName + "'s totalSupply is " + string(totalSupplyBytes))

	return shim.Success(totalSupplyBytes)
}

func (cc *Chaincode) balanceOf(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success([]byte("balanceOf call !!"))
}

func (cc *Chaincode) transfer(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success([]byte("transfer call!!!!!"))
}

func (cc *Chaincode) allowance(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}

func (cc *Chaincode) approve(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}

func (cc *Chaincode) transferFrom(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}

func (cc *Chaincode) increaseAllowance(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}

func (cc *Chaincode) decreaseAllowance(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}

func (cc *Chaincode) mint(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}

func (cc *Chaincode) burn(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}
