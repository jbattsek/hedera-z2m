package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"os"
	"time"
)

func main() {
	//
	// Generate keys
	//

	// Read and decode the operator secret key
	operatorSecret, err := hedera.SecretKeyFromString(os.Getenv("SECRET_KEY"))
	if err != nil {
		panic(err)
	}

	// Generate a new keypair for the new account
	secret, _ := hedera.GenerateSecretKey()
	public := secret.Public()

	fmt.Printf("secret = %v\n", secret)
	fmt.Printf("public = %v\n", public)

	//
	// Connect to Hedera
	//

	client, err := hedera.Dial("testnet.hedera.com:50119")
	if err != nil {
		panic(err)
	}

	defer client.Close()

	//
	// Send transaction to create account
	//

	nodeAccountID := hedera.AccountID{Account: 3}
	operatorAccountID := hedera.AccountID{Account: 1003}
	response, err := client.CreateAccount().
		Key(public).
		InitialBalance(0).
		Operator(operatorAccountID).
		Node(nodeAccountID).
		Memo("[test] hedera-sdk-go v2").
		Sign(operatorSecret).
		Execute()

	if err != nil {
		panic(err)
	}

	transactionID := response.ID
	fmt.Printf("created account; transaction = %v\n", transactionID)

	//
	// Get receipt to prove we created it ok
	//

	fmt.Printf("wait for 2s...\n")
	time.Sleep(2 * time.Second)

	receipt, err := client.Transaction(*transactionID).Receipt().Get()
	if err != nil {
		panic(err)
	}

	if receipt.Status != hedera.StatusSuccess {
		panic(fmt.Errorf("transaction has a non-successful status: %v", receipt.Status.String()))
	}

	fmt.Printf("account = %v\n", *receipt.AccountID)
}
