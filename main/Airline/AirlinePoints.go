package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/util"
	
		
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
// arg[5] is the chaincode hash value for MyWallet
func addPoints(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5 walletid,entity name, transactionid, reward,  points & hash value of wallet ")
	}

	fmt.Println(" Calling AddPoints  ")
	var chainCodeToCall = args[5]
	f := "addpointstowallet"
	
	invokeArgs := util.ToChaincodeArgs(f, args[0], args[1], args[2], args[3], args[4])
	userWalletByte, err := stub.InvokeChaincode(chainCodeToCall, invokeArgs)
	if err != nil {
		errStr := fmt.Sprintf("Failed to invoke chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return nil, errors.New(errStr)
	}
	
	fmt.Println("Successfully Added points to user wallet " + args[0])

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
// arg[3] is the chaincode hash value for MyWallet

func getAirlinePointTransactionByUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	logger.Debug("Entering getAirlinePointTransactionByUser ")
	if len(args) !=4 {
		return nil, errors.New("Incorrect number of arguments. Expecting  4")
	}

	
	fmt.Println(" Calling MyWallet Chain code to get user transaction   ")
	var chainCodeToCall = args[3]
	f := "query"
	
	invokeArgs := util.ToChaincodeArgs(f, args[0], args[1], args[2], args[3])
	userAirlineWalletTranByte, err := stub.InvokeChaincode(chainCodeToCall, invokeArgs)
	if err != nil {
		errStr := fmt.Sprintf("Failed to invoke chaincode. Got error: %s", err.Error())
		fmt.Printf(errStr)
		return nil, errors.New(errStr)
	}

	fmt.Println(" Success in getting the user entity transaction ")
	return userAirlineWalletTranByte, nil

}

// Query
func (t *SampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "read" {
		return getAirlinePointTransactionByUser(stub, args)
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
