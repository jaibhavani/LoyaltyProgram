package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
		
)
type LoyaltyPointWallet struct {
	Name         string `json:"name"`
	Password     string `json:"password"`
	PointBalance int    `json:"pointbalance"`
}

type PointsTransaction struct {
	Name            string `json:"name"`
	Entity          string `json:"entity"`
	TransactionID   string `json:"transactionid"`
	TransactionType string `json:"transactiontype"`
	LoyaltyPoints   int    `json:"loyaltypoints"`
}

// SimpleChaincode example simple Chaincode implementation
type SampleChaincode struct {
}

var logger = shim.NewLogger("mylogger")

// User defined function

// This chain code will be invoked by Airline application code where the
// arg[0] is the Name
// arg[1] is the Entity Name
// arg[2] is the TransactionsID
// arg[3] transaction type - reward
// arg[4] is the LoyaltyPoints
func addPoints(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5 walletid,entity name, transactionid, reward and points")
	}

	_, err := AddPointsToWallet(stub, args)

	if err != nil {

		fmt.Println(err)
		return nil, errors.New("Errors while Adding points to wallet for user  " + args[0])
	}

	logger.Info("Successfully Added points to user wallet " + args[0])

	return nil, nil
}


func AddPointsToWallet(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering add points to Loyalty wallet ")

	if len(args) < 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting  name, entity, transactionid, transactiontype, points")
	}

	var name = args[0]

	wallet, err := stub.GetState(name)

	if err != nil {
		return nil, errors.New("Error while getting wallet data for user " + name)
	}
	// Check if user has wallet record in the ledger. If not, then create the wallet
	if wallet == nil {
		return nil, errors.New("No wallet data exists for user " + name)
	}

	// Store the reward transaction to the ledger

	points, err := strconv.Atoi(args[4])

	if err != nil {
		return nil, errors.New("Expecting integer value for transaction points as arg 4")
	}
	var rewardTran PointsTransaction
	rewardTran.Name = name
	rewardTran.Entity = args[1]
	rewardTran.TransactionID = args[2]
	rewardTran.TransactionType = args[3]
	rewardTran.LoyaltyPoints = points

	rewardTranBytes, err := json.Marshal(rewardTran)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for reward transaction")
	}

	err = stub.PutState(args[0]+args[1]+args[2], rewardTranBytes)

	if err != nil {
		return nil, err
	}

	logger.Debug("addd reward transaction to the ledger ")

	var userWallet LoyaltyPointWallet
	err = json.Unmarshal(wallet, &userWallet)

	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of user " + name)
	}

	// Add the new points to the wallet balance

	awardPoints, err := strconv.Atoi(args[3])

	if err != nil {
		return nil, errors.New("Points awarded from entity " + args[1] + "  must be integer")
	}
	userWallet.PointBalance = userWallet.PointBalance + awardPoints

	userWalletByte, err := json.Marshal(userWallet)

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string after updating points for user " + name)
	}

	err = stub.PutState(userWallet.Name, userWalletByte)
	if err != nil {
		return nil, err
	}

	logger.Debug("Added points to wallet success fully ")

	return nil, nil

}

// Init function
func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, nil
}

// Invoke
func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	if function == "addPoints" {
		return addPoints(stub, args)
	} 

	return nil, nil
}


// Query
func (t *SampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, nil
}

// main
func main() {
	lld, _ := shim.LogLevel("DEBUG")
	fmt.Println(lld)

	logger.SetLevel(lld)
	fmt.Println(logger.IsEnabledFor(lld))

	err := shim.Start(new(SampleChaincode))
	if err != nil {
		logger.Error("Could not start SampleChaincode")
	} else {
		logger.Info("SampleChaincode successfully started")
	}
}
