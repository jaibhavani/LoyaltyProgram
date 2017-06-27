package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/jaibhavani/LoyaltyProgram/LoyaltyPkgUtil"
)

// SimpleChaincode example simple Chaincode implementation
type SampleChaincode struct {
}

var logger = shim.NewLogger("mylogger")

func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Info(" length of arguments " + strconv.Itoa(len(args)))
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	_, err := LoyaltyPkg.CreateWallet(stub, args)

	if err != nil {

		fmt.Println(err)
		return nil, errors.New("Errors while creating  wallet for user  " + args[0])
	}

	logger.Info("Successfully created wallet for user  " + args[0])

	return nil, nil

}

// This chain code will be invoked by user from UI to redeem Loyalty points from wallet to entity specific points/cash
 //arg[0] username - User wallet id
// arg[1] entity - Points redeemed to the entity name - like  airline name, bank name, hotel  name etc
// arg[2] transactionid -- Transaction id from the entity system for points redeemed
// arg[3] transaction type - redeem
// arg[4] loyalty points to be deducted from the user wallet
// Thus function will call the corresponding entity API to redeem th points to the agreed format. Once its successful 
func redeem (stub shim.ChaincodeStubInterface, args []string) ([]byte,error){

	if len(args) != 5 {
			return nil, errors.New("Incorrect number of arguments. Expecting 5")
		}

		_, err := LoyaltyPkg.RedeemPoints(stub, args)

	if err != nil {

		fmt.Println(err)
		return nil, errors.New("Errors while creating  wallet for user  " + args[0])
	}

	logger.Info("Successfully redeemed points  for user  " + args[0])

	return nil, nil

}


// This chain code will be invoked wallet ui application code to create wallet
// arg[0] is the Name
// arg[1] is the password
// arg[2] is the points

func createWallet(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3 name, password, points")
	}

	_, err := LoyaltyPkg.CreateWallet(stub, args)

	if err != nil {

		fmt.Println(err)
		return nil, errors.New("Errors while creating  wallet for user  " + args[0])
	}

	logger.Info("Successfully created wallet for user  " + args[0])

	return nil, nil
}

// Invoke
func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	if function == "redeem" {
		return redeem(stub, args)
	} else if function == "createwallet" {
		return createWallet(stub, args)
	} 

	return nil, nil
}

// Query
func (t *SampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "read" {
		return LoyaltyPkg.GetUserLoyaltyWallet(stub, args)
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
