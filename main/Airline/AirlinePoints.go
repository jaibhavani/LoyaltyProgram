package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/jaibhavani/LoyaltyProgram/LoyaltyPkgUtil"
	
)

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

	_, err := LoyaltyPkg.AddPointsToWallet(stub, args)

	if err != nil {

		fmt.Println(err)
		return nil, errors.New("Errors while Adding points to wallet for user  " + args[0])
	}

	logger.Info("Successfully Added points to user wallet " + args[0])

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
