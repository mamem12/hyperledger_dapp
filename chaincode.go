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

// TransferEvent is the Event
type TransferEvent struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Amount    int    `json:"amount"`
}

type Approval struct {
	Spender   string `json:"spender"`
	Owner     string `json:"Owner"`
	Allowance int    `json:"allowance"`
}

// Init is called when the chaincode is instantiated by the blockchain network.
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) sc.Response {
	_, params := stub.GetFunctionAndParameters()
	fmt.Println("Init called with params: ", params)

	return shim.Success(nil)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("Invoke()", fcn, params)

	switch fcn {
	case "init":
		return cc.init(stub, params)
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
	case "approvalList":
		return cc.approvalList(stub, params)
	case "transferOtherToken":
		return cc.transferOtherToken(stub, params)
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

// totalSupply is query function
// params is tokenName
// Returns the amount of token in existence
func (cc *Chaincode) totalSupply(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// a number of params must be one
	if len(params) != 1 {
		return shim.Error("totalSupply only 1 params")
	}

	tokenName := params[0]

	// get metadata
	metadata := &Metadata{}
	metadataBytes, err := stub.GetState(tokenName)
	if err != nil {
		errMsg := fmt.Sprintf("failed to getstate from totalSupply, error : %s", err.Error())
		return shim.Error(errMsg)
	}

	err = json.Unmarshal(metadataBytes, metadata)
	if err != nil {
		errMsg := fmt.Sprintf("failed to unmarshal from totalSupply, error : %s", err.Error())
		return shim.Error(errMsg)
	}

	// convert metadata to bytes
	totalsupplyBytes, err := json.Marshal(metadata.TotalSupply)
	if err != nil {
		errMsg := fmt.Sprintf("failed to marshal from totalSupply, error : %s", err.Error())
		return shim.Error(errMsg)
	}

	fmt.Println(tokenName + "'s totalsupply is " + string(totalsupplyBytes))

	return shim.Success(totalsupplyBytes)
}

// totalSupply is query function
// params is address
// Returns the amount of tokens owned by address
func (cc *Chaincode) balanceOf(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	// a number of params must be one
	if len(params) != 1 {
		return shim.Error("balanceOf only 1 params")
	}

	address := params[0]

	// get balance
	amountBytes, err := stub.GetState(address)
	if err != nil {
		errMsg := fmt.Sprintf("failed to getstate from balanceOf, error : %s", err.Error())
		return shim.Error(errMsg)
	}

	fmt.Println(address + "'s totalsupply is " + string(amountBytes))

	return shim.Success(amountBytes)
}

// transfer is invoke fcn that moves amount token
// from the caller's address to recipient
// params - caller's address, recipient's address, amount of token
func (cc *Chaincode) transfer(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check parameter
	if len(params) != 3 {
		shim.Error("transfer only 3 params")
	}

	callerAddress, recipientAddress, transferAmount := params[0], params[1], params[2]

	transferAmountInt, err := strconv.Atoi(transferAmount)
	if err != nil {
		return shim.Error("transfer amount must be integer")
	}

	if transferAmountInt <= 0 {
		return shim.Error("transfer amount must be positive")
	}

	// get caller amount
	callerAmount, err := stub.GetState(callerAddress)
	if err != nil {
		return shim.Error("failed to GetState Error" + err.Error())
	}

	callerAmountInt, err := strconv.Atoi(string(callerAmount))
	if err != nil {
		return shim.Error("caller amount must be integer")
	}

	// get recipient amount
	recipientAmount, err := stub.GetState(recipientAddress)
	if err != nil {
		return shim.Error("failed to GetState Error" + err.Error())
	}

	if recipientAmount == nil {
		recipientAmount = []byte("0")
	}

	recipientAmountInt, err := strconv.Atoi(string(recipientAmount))
	if err != nil {
		return shim.Error("recipient amount must be integer")
	}

	// calculate amount
	callerResultAmount := callerAmountInt - transferAmountInt
	recipientResultAmount := recipientAmountInt + transferAmountInt

	// check calculate amount is positive
	if callerResultAmount < 0 {
		return shim.Error("caller's balance is not sufficient")
	}

	// save the caller & recipient amount
	err = stub.PutState(callerAddress, []byte(strconv.Itoa(callerResultAmount)))
	if err != nil {
		return shim.Error("failed to PutState of caller")
	}
	err = stub.PutState(recipientAddress, []byte(strconv.Itoa(recipientResultAmount)))
	if err != nil {
		return shim.Error("failed to PutState of recipient")
	}

	// emit transfer event
	transferEvent := TransferEvent{
		Sender:    callerAddress,
		Recipient: recipientAddress,
		Amount:    transferAmountInt,
	}

	transferEventBytes, err := json.Marshal(transferEvent)
	if err != nil {
		return shim.Error("failed to marshal transferEvent, error :" + err.Error())
	}

	err = stub.SetEvent("transferEvent", transferEventBytes)
	if err != nil {
		return shim.Error("failed to SetEvent of transferEvent, error :" + err.Error())
	}

	fmt.Println(callerAddress + " send " + transferAmount + " to " + recipientAddress)

	return shim.Success([]byte("transfer Success"))
}

// allowance is query fnc
// params - owner's address, spender's address
// return - the remaining amount of token to invoke (transferFrom)
func (cc *Chaincode) allowance(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check the number of params is 2
	if len(params) != 2 {
		return shim.Error("allowance only 2 params")
	}

	ownerAddress, spenderAddress := params[0], params[1]

	// create composite key
	approvalKey, err := stub.CreateCompositeKey("approval", []string{ownerAddress, spenderAddress})
	if err != nil {
		return shim.Error("failed CreateCompositeKey for approval")
	}

	// getstate(amount)
	approvalKeyBytes, err := stub.GetState(approvalKey)
	if err != nil {
		return shim.Error("Failed to GetState")
	}

	if approvalKeyBytes == nil {
		approvalKeyBytes = []byte("0")
	}

	return shim.Success(approvalKeyBytes)
}

// approvalList is query fnc
// params - owner's addresss
// return - approvalList by owner
func (cc *Chaincode) approvalList(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check the number of params is 1
	if len(params) != 1 {
		return shim.Error("allowance only 2 params")
	}

	ownerAddress := params[0]

	// get all approval (format is iterator)
	approvalIterator, err := stub.GetStateByPartialCompositeKey("approval", []string{ownerAddress})
	if err != nil {
		return shim.Error("failed to GetStateByPartialCompositeKey for approval iterator error :" + err.Error())
	}

	// make slice for return value
	approvalSlice := []Approval{}

	// iterator
	defer approvalIterator.Close()
	if approvalIterator.HasNext() {
		for approvalIterator.HasNext() {
			approvalKV, _ := approvalIterator.Next()
			approvalKey := approvalKV.GetKey()
			approvalValue := approvalKV.GetValue()

			// get sppender address
			_, address, err := stub.SplitCompositeKey(approvalKey)
			if err != nil {
				return shim.Error("failed to SplitCompositeKey, error :" + err.Error())
			}

			spenderAddress := address[1]

			// get amount
			amount, err := strconv.Atoi(string(approvalValue))
			if err != nil {
				return shim.Error("failed to get amount, error : " + err.Error())
			}

			// add approval result
			approvalSlice = append(approvalSlice, Approval{Spender: spenderAddress, Owner: ownerAddress, Allowance: amount})
		}
	}

	// marshal data
	response, err := json.Marshal(approvalSlice)
	if err != nil {
		return shim.Error("failed to Marshal approvalSlice, error : " + err.Error())
	}

	return shim.Success(response)
}

// approve is invoke fnc that sets amount as the allowance
// of spender over the owner tokens
// params - owner's address, spender's address, amount of token
func (cc *Chaincode) approve(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check the number of params is 3
	if len(params) != 3 {
		return shim.Error("approve only 3 params")
	}

	ownerAddress, spenderAddress, allowanceAmount := params[0], params[1], params[2]

	// check amount is integer & positive
	allowanceAmountInt, err := strconv.Atoi(allowanceAmount)

	if err != nil {
		return shim.Error("Amount is must be Integer")
	}

	if allowanceAmountInt <= 0 {
		return shim.Error("Amount is must be positive")
	}

	// create composite key for allowance - approval/{owner}/{spender}
	approvalKey, err := stub.CreateCompositeKey("approval", []string{ownerAddress, spenderAddress})
	if err != nil {
		return shim.Error("failed to create CompositeKey for approval")
	}

	// save allowance amount
	err = stub.PutState(approvalKey, []byte(allowanceAmount))
	if err != nil {
		return shim.Error("failed to PutState for approvalKey")
	}

	// emit approval event
	approvalEvent := Approval{Owner: ownerAddress, Spender: spenderAddress, Allowance: allowanceAmountInt}
	approvalBytes, err := json.Marshal(approvalEvent)
	if err != nil {
		return shim.Error("failed to Marshal for approvalEvent")
	}

	err = stub.SetEvent("approvalEvent", approvalBytes)
	if err != nil {
		return shim.Error("failed to marshal for approval SetEvent")
	}

	return shim.Success([]byte("approve success"))
}

// transferFrom is invoke fnc that moves amount of token from sender to recipient
// using allowance of sender
// params - owner's address, spender's address, recipient's address, amount of token
func (cc *Chaincode) transferFrom(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	if len(params) != 4 {
		return shim.Error("transferFrom only 4 params")
	}

	ownerAddress, spenderAddress, recipientAddress, transferAmount := params[0], params[1], params[2], params[3]

	// check amount is integer & positive
	transferAmountInt, err := strconv.Atoi(transferAmount)

	if err != nil {
		return shim.Error("Amount is must be Integer")
	}

	if transferAmountInt < 0 {
		return shim.Error("Amount is must be positive")
	}

	// get allowance
	allowanceResponse := cc.allowance(stub, []string{ownerAddress, spenderAddress})

	if allowanceResponse.GetStatus() >= 400 {
		return shim.Error("failed to get allowance, error : " + allowanceResponse.GetMessage())
	}

	// convert allowance response payload to allowance data
	allowanceInt, err := strconv.Atoi(string(allowanceResponse.GetPayload()))
	if err != nil {
		return shim.Error("allowance must be positive")
	}

	// transfer from owner to recipient
	transferResponse := cc.transfer(stub, []string{ownerAddress, recipientAddress, transferAmount})
	if transferResponse.GetStatus() >= 400 {
		return shim.Error("failed to get transfer error : " + err.Error())
	}

	// decrease allowance amount
	approveAmountInt := allowanceInt - transferAmountInt
	approveAmount := strconv.Itoa(approveAmountInt)

	// approve amount of tokens trasfered
	approveResponse := cc.approve(stub, []string{ownerAddress, spenderAddress, approveAmount})
	if approveResponse.GetStatus() >= 400 {
		return shim.Error("failed to get transfer error : " + err.Error())
	}

	return shim.Success([]byte("transferFrom success"))
}

// transferOtherToken is invoke fnc that moves amount other chaincode tokens
// from the caller's addresss to recipient
// params - chaincode name, caller's addresss, recipient's address, amount
func (cc *Chaincode) transferOtherToken(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	if len(params) != 4 {
		return shim.Error("transferFrom only 4 params")
	}

	chaincodeName, callerAddress, recipientAddress, transferAmount := params[0], params[1], params[2], params[3]

	// make arguments
	args := [][]byte{[]byte("transfer"), []byte(callerAddress), []byte(recipientAddress), []byte(transferAmount)}

	// get channel
	channel := stub.GetChannelID()

	// transfer other chaincode token
	transferResponse := stub.InvokeChaincode(chaincodeName, args, channel)

	if transferResponse.GetStatus() >= 400 {
		errMsg := fmt.Sprintf("failed to transfer %s, error :%s ", chaincodeName, transferResponse.GetMessage())
		return shim.Error(errMsg)
	}

	return shim.Success([]byte("transfer other token success"))
}

// increaseAllowance is invoke fnc that increases spender's allowance by owner
// params - owner's address, spender's addresss, amount of token
func (cc *Chaincode) increaseAllowance(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check params 3
	if len(params) != 3 {
		return shim.Error("increaseAllowance only 3 params")
	}

	ownerAddress, spenderAddress, increaseAmount := params[0], params[1], params[2]
	// check amount is integer & positive

	increaseAmountInt, err := strconv.Atoi(increaseAmount)
	if err != nil {
		return shim.Error("amount must be integer")
	}

	if increaseAmountInt < 0 {
		return shim.Error("amount must be positive")
	}

	// get allowance
	allowanceResponse := cc.allowance(stub, []string{ownerAddress, spenderAddress})
	if allowanceResponse.GetStatus() >= 400 {
		return shim.Error("failed to get allowance, error : " + allowanceResponse.GetMessage())
	}

	// convert allowance response payload to allowance data
	allowanceInt, err := strconv.Atoi(string(allowanceResponse.GetPayload()))
	if err != nil {
		return shim.Error("allowance must be positive")
	}

	// increase allowance
	resultAmountInt := allowanceInt + increaseAmountInt
	resultAmount := strconv.Itoa(resultAmountInt)

	// call approve
	approveResponse := cc.approve(stub, []string{ownerAddress, spenderAddress, resultAmount})
	if approveResponse.GetStatus() >= 400 {
		return shim.Error("failed to get allowance, error : " + approveResponse.GetMessage())
	}

	return shim.Success([]byte("increaseAllowance success"))
}

// decreaseAllowance is invoke fnc that decreases spender's allowance by owner
// params - owner's address, spender's addresss, amount of token
func (cc *Chaincode) decreaseAllowance(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check params 3
	if len(params) != 3 {
		return shim.Error("decreaseAllowance only 3 params")
	}

	ownerAddress, spenderAddress, decreaseAmount := params[0], params[1], params[2]
	// check amount is integer & positive

	decreaseAmountInt, err := strconv.Atoi(decreaseAmount)
	if err != nil {
		return shim.Error("amount must be integer")
	}

	if decreaseAmountInt < 0 {
		return shim.Error("amount must be positive")
	}

	// get allowance
	allowanceResponse := cc.allowance(stub, []string{ownerAddress, spenderAddress})
	if allowanceResponse.GetStatus() >= 400 {
		return shim.Error("failed to get allowance, error : " + allowanceResponse.GetMessage())
	}

	// convert allowance response payload to allowance data
	allowanceInt, err := strconv.Atoi(string(allowanceResponse.GetPayload()))
	if err != nil {
		return shim.Error("allowance must be positive")
	}

	// decrease allowance
	resultAmountInt := allowanceInt + decreaseAmountInt
	resultAmount := strconv.Itoa(resultAmountInt)

	// call approve
	approveResponse := cc.approve(stub, []string{ownerAddress, spenderAddress, resultAmount})
	if approveResponse.GetStatus() >= 400 {
		return shim.Error("failed to get allowance, error : " + approveResponse.GetMessage())
	}

	return shim.Success([]byte("decreaseAllowance success"))
}

func (cc *Chaincode) mint(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}

func (cc *Chaincode) burn(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}
