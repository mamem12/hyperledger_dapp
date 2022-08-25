/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"
	"hyperledger_dapp/controller"
	"hyperledger_dapp/model"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
)

// Chaincode is the definitaion of the chaincode structure
type Chaincode struct {
	Controller *controller.Controller
}

// NewChaincode is construtor function for Chaincode
func NewChaincode() *Chaincode {
	controller := controller.NewContoller()
	return &Chaincode{Controller: controller}
}

// Init is called when the chaincode is instantiated by the blockchain network.
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) sc.Response {
	_, params := stub.GetFunctionAndParameters()
	fmt.Println("Init called with params: ", params)

	return shim.Success(nil)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fnc, params := stub.GetFunctionAndParameters()

	switch fnc {
	case "init":
		return cc.init(stub, params)
	case "totalSupply":
		return cc.Controller.TotalSupply(stub, params)
	case "balanceOf":
		return cc.Controller.BalanceOf(stub, params)
	case "transfer":
		return cc.Controller.Transfer(stub, params)
	case "allowance":
		return cc.Controller.Allowance(stub, params)
	case "approve":
		return cc.Controller.Approve(stub, params)
	case "transferFrom":
		return cc.Controller.TransferFrom(stub, params)
	case "increaseAllowance":
		return cc.Controller.IncreaseAllowance(stub, params)
	case "decreaseAllowance":
		return cc.Controller.DecreaseAllowance(stub, params)
	case "approvalList":
		return cc.Controller.ApprovalList(stub, params)
	case "transferOtherToken":
		return cc.Controller.TransferOtherToken(stub, params)
	case "mint":
		return cc.Controller.Mint(stub, params)
	case "burn":
		return cc.Controller.Burn(stub, params)
	default:
		return sc.Response{Status: 404, Message: "404 Not Found", Payload: nil}
	}
}

func (cc *Chaincode) init(stub shim.ChaincodeStubInterface, params []string) sc.Response {
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
	metadata := model.NewMetadata(tokenName, symbol, owner, uint(amountUint))
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return shim.Error("failed to Marshal, error: " + err.Error())
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
