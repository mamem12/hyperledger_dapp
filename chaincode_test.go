package main

import (
	"encoding/json"
	"fmt"
	"hyperledger_dapp/model"
	"hyperledger_dapp/repository"
	"strconv"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
)

const (
	initTokenName = "dappToken"
	initSymbol    = "dt"
	initOwner     = "dappcampus"
	initAmount    = 1000000
	txMint        = "txMint"
)

var function = []byte("mint")

func configuration() *shimtest.MockStub {
	cc := NewChaincode()
	stub := shimtest.NewMockStub("erc20", cc)
	stub.MockInit("1", [][]byte{[]byte("Init"), []byte(initTokenName), []byte(initSymbol), []byte(initOwner), []byte(strconv.Itoa(initAmount))})

	return stub
}

func TestInit(t *testing.T) {
	cc := NewChaincode()
	stub := shimtest.NewMockStub("erc20", cc)
	res := stub.MockInit("1", [][]byte{[]byte("Init"), []byte(initTokenName), []byte(initSymbol), []byte(initOwner), []byte(strconv.Itoa(initAmount))})
	if res.Status != shim.OK {
		t.Error("init failed", res.Status, res.Message)
	}

	// check totalSupply
	erc20 := model.ERC20Metadata{}
	erc20Bytes, _ := stub.GetState(initTokenName)
	json.Unmarshal(erc20Bytes, &erc20)
	totalSupply := *erc20.GetTotalSupply()
	if initAmount != totalSupply {
		t.FailNow()
	}

	balance, _ := repository.GetBalance(stub, initOwner, false)

	fmt.Println(*balance)

	// check dappcampus balance
}

func TestMint(t *testing.T) {
	stub := configuration()
	const increasAmount = 100000
	arguments := [][]byte{function, []byte(initTokenName), []byte(initOwner), []byte(strconv.Itoa(increasAmount))}
	res := stub.MockInvoke(txMint, arguments)
	if res.Status != shim.OK {
		t.FailNow()
	}

	// increase total supply
	totalSupply, _ := repository.GetERC20TotalSupply(stub, initTokenName)
	if *totalSupply != initAmount+increasAmount {
		t.FailNow()
	}

	balance, err := repository.GetBalance(stub, initOwner, false)
	if err != nil {
		fmt.Println("fail to get balance")
		t.FailNow()
	}

	if *balance != initAmount+increasAmount {
		fmt.Println("not eqaul balance")
		t.FailNow()
	}

	// emit transfer event
	data := <-stub.ChaincodeEventsChannel
	if data.GetEventName() == repository.TransferEventKey {
		t.FailNow()
	}

	event := model.NewTransferEvent("admin", initTokenName, increasAmount)

	eventBytes, _ := json.Marshal(event)

	if string(data.Payload) != string(eventBytes) {
		t.FailNow()
	}

	fmt.Println(string(res.Payload))
}
