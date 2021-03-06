package LoyaltyPkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("mylogger")

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

func GetUserLoyaltyWallet(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	logger.Debug("Entering Get Loyalty wallet ")
	if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting  name")
	}

	var name = args[0]

	bytes, err := stub.GetState(name)

	if err != nil {
		return nil, errors.New("Error while getting wallet data for user " + name)
	}
	return bytes, nil

}

// This function will be invoked by Entity application code to get the point allocation transaction
// arg[0] is the UserwalletName
// arg[1] is the Entity Name
// arg[2] is the TransactionsID

func GetUserEntityPointsInWallet(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	logger.Debug("Entering Get User Entity Points In  wallet ")
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting  name, entity and transactionid ")
	}

	var name = args[0]

	bytes, err := stub.GetState(args[0]+args[1]+ args[2])

	if err != nil {
		fmt.Println(" Error while getting the data for user " + name + " and Entity " + args[1])
		return nil, errors.New("Error while getting entity transaction wallet data for user " + name)
	}
	return bytes, nil

}

// This function is called to create loyalty wallet for a user
// arg[0] name
// arg[1] password
// arg[2] points For the first time it will be 0
func CreateWallet(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering Create Loyalty wallet ")
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting  name, password, points")
	}

	var userWallet LoyaltyPointWallet
	userWallet.Name = args[0]
	userWallet.Password = args[1]
	points, err := strconv.Atoi(args[2])

	if err != nil {
		return nil, errors.New("Expecting integer value for points in CreateWallet Function")
	}
	userWallet.PointBalance = points

	userWalletByte, err := json.Marshal(userWallet)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for user wallet ")
	}
	err = stub.PutState(args[0], userWalletByte)
	if err != nil {
		return nil, err
	}

	return nil, nil

}

// Add Points to the user wallet. This function will be invoked by the MyWallet program
// arg[0] name	- User wallet name
// arg[1] entity - entity who is rewarding the loyalty points
// arg[2] transactionid - The transaction id from the entity system
// arg[3] transaction type - reward
// arg[4] loyalty points - total loyalty points awarded by the entity for the transaction
func AddPointsToWallet(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering add points to Loyalty wallet ")
	
	fmt.Println("Entering add points to Loyalty wallet ")

	if len(args) !=5 {
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

	fmt.Println("addd reward transaction to the ledger. Now adding points  ")

	var userWallet LoyaltyPointWallet
	err = json.Unmarshal(wallet, &userWallet)

	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of user " + name)
	}

	fmt.Println("*** SDS--- Now adding points ***********")

	fmt.Println(" Currnent Wallet balance " +  strconv.Itoa(userWallet.PointBalance) + " additional reward point   " + strconv.Itoa(points))

	userWallet.PointBalance = userWallet.PointBalance + points

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
// Redeem Points from user wallet. This function will be invoked by the MyWallet program
// arg[0] username - User wallet id
// arg[1] entity - Points redeemed to the entity name - like  airline name, bank name, hotel  name etc
// arg[2] transactionid -- Transaction id from the entity system for points redeemed
// arg[3] transaction type - redeem
// arg[4] loyalty points to be deducted from the user wallet
func RedeemPoints(stub shim.ChaincodeStubInterface, args []string)([]byte, error){

	logger.Debug("Entering Redeem points to Loyalty wallet ")

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
		return nil, errors.New("Errors while creating json string for redeem transaction")
	}

	err = stub.PutState(args[0]+args[1]+args[2], rewardTranBytes)

	if err != nil {
		return nil, errors.New("Error while saving the redeem transaction to the world state")
	}

	logger.Debug("addd reward transaction to the ledger ")

	var userWallet LoyaltyPointWallet
	err = json.Unmarshal(wallet, &userWallet)

	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of user " + name)
	}

	// Add the new points to the wallet balance

	userWallet.PointBalance = userWallet.PointBalance - points

	userWalletByte, err := json.Marshal(userWallet)

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string after updatredeeming points for user " + name)
	}

	err = stub.PutState(name, userWalletByte)
	if err != nil {
		return nil, err
	}

	logger.Debug("redeemed points to wallet success fully for user "+ name)

	return nil, nil


}



