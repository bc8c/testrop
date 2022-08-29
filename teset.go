/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

type UserRating struct{
	User string			`json:"user"`
	Average float64		`json:"average"`
	Rates []Rate		`json:"rates"`
}

type Rate struct{
	ProjectTitle string	`json:"projecttitle"`
	Score float64		`json:"score"`
}


// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) AddUser(ctx contractapi.TransactionContextInterface, username string) error {

	var user = UserRating{User:username, Average:0}
	userAsBytes, _ := json.Marshal(user)


	var err = ctx.GetStub().PutState(username, userAsBytes)

	return err
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) AddRating(ctx contractapi.TransactionContextInterface, username string, prjTitle string, prjscore string) error {
	// UserRating 조회 : username
	userAsBytes, err := ctx.GetStub().GetState(username)

	if err != nil{
		return err
	} else if userAsBytes == nil {
		return fmt.Errorf("Errror: User does not exist : " + username)
	}

	user := UserRating{}
	err = json.Unmarshal(userAsBytes, &user)

	if err != nil {
		return err
	}
	//  prjTitle,  prjscore 추가
	newRate, _ := strconv.ParseFloat(prjscore, 64)
	var rate = Rate{ProjectTitle:prjTitle, Score:newRate}
	

	// 평균값 계산
	rateCount := float64(len(user.Rates))
	user.Average = (rateCount*user.Average + newRate)/(rateCount+1)
	user.Rates = append(user.Rates, rate)


	// putstat... 
	 userAsBytes, err = json.Marshal(user)
	 if err != nil {
		return err
	 }

	 err = ctx.GetStub().PutState(username,userAsBytes)
	 if err != nil {
		return err
	 }

	 return nil
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) ReadRating(ctx contractapi.TransactionContextInterface, username string) (string, error) {
	UserAsBytes, err := ctx.GetStub().GetState(username)

	if err != nil {
		return "",err
	}

	if UserAsBytes == nil {
		return "", fmt.Errorf("user does not exist : " + username)
	}
	return string(UserAsBytes[:]), nil
}



func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
