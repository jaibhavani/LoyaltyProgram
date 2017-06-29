package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/jaibhavani/LoyaltyProgram/LoyaltyPkgUtil"
	
		
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

	fmt.Println(" Calling AddPoints ")
	userWalletByte, err := LoyaltyPkgUtil.AddPointsToWallet(stub, args)

	if err != nil {

		fmt.Println(err)
		return nil, errors.New("Errors while Adding points to wallet for user  " + args[0])
	}

	fmt.Println("Successfully Added points to user wallet " + args[0])

	return userWalletByte, nil
}


func AddPointsToWallet(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering add points to Loyalty wallet ")

	if len(args) < 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting  name, entity, transactionid, transactiontype, points")
	}

	var name = args[0]
	fmt.Println("Get the wallet for user " + name)

	wallet, err := stub.GetState(name)

	if err != nil {
		fmt.Println("No wallet data found for user " + name)
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

	fmt.Println(" marshalling reward transaction to bytes " )

	rewardTranBytes, err := json.Marshal(rewardTran)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for reward transaction")
	}

	fmt.Println(" Storing the reward tran to state with ID " + args[0]+args[1]+args[2])
	err = stub.PutState(args[0]+args[1]+args[2], rewardTranBytes)

	if err != nil {
		return nil, err
	}

	fmt.Println("addd reward transaction to the ledger ")

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

	fmt.Println(" Currnent Wallet balance " +  strconv.Itoa(userWallet.PointBalance) + " additional reward point   " + strconv.Itoa(awardPoints))

	userWallet.PointBalance = userWallet.PointBalance + awardPoints

	userWalletByte, err := json.Marshal(userWallet)

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string after updating points for user " + name)
	}

	err = stub.PutState(name, userWalletByte)
	if err != nil {
		return nil, err
	}

	fmt.Println("Added points to wallet success fully ")

	return userWalletByte, nil

}

// Init function
func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	return nil, nil
}

// Invoke
func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	if function == "addPoints" {
		fmt.Println(" Calling addPoints function from Invoke ")
		return addPoints(stub, args)
	} 

	return nil, nil
}

// This chain code will be invoked by Airline application code to get the point allocation transaction
// arg[0] is the UserwalletName
// arg[1] is the Entity Name
// arg[2] is the TransactionsID

func GetAirlinePointTransactionByUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	logger.Debug("Entering Get Loyalty wallet ")
	if len(args) < 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting  3")
	}

	var name = args[0]
	var entity = args[1]
	var transaction = args[2]

	bytes, err := stub.GetState(name+entity+transaction)

	if err != nil {
		return nil, errors.New("Error while getting Airline point transaction data for user " + name)
	}
	return bytes, nil

}

// Query
func (t *SampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "read" {
		return GetAirlinePointTransactionByUser(stub, args)
	}
	fmt.Println(" Invalid function passed to Query function " + function)

	return nil, errors.New(" Invalid function passed to Query function " + function)
	
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
