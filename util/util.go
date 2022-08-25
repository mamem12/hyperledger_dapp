package util

import (
	"hyperledger_dapp/model"
	"strconv"
)

func ConverToPositive(name, value string) (*int, error) {
	intValue, err := strconv.Atoi(value)

	if err != nil {
		return nil, model.NewCustomError(model.ConvertErrorType, name, " must be integer")
	}

	if intValue <= 0 {
		return nil, model.NewCustomError(model.ConvertErrorType, name, " must be positive")
	}

	return &intValue, nil
}
