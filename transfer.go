package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"os"
	"time"
)

func main() {
	// Read and decode the operator secret key
	operatorAccountID := hedera.AccountID{Account: 1003}
	operatorSecret, err := hedera.SecretKeyFromString(os.Getenv("OPERATOR_SECRET"))
	if err != nil {
		panic(err)
	}

	targetAccountID := hedera.AccountID{Account: 1004}
	// Read and decode target account
	/*targetAccountID, err := hedera.AccountIDFromString(os.Getenv("TARGET"))
	if err != nil {
		panic(err)
	}*/

	//
	// Connect to Hedera
	//

	client, err := hedera.Dial("testnet.hedera.com:50119")
	if err != nil {
		panic(err)
	}

	client.SetNode(hedera.AccountID{Account: 3})
	client.SetOperator(operatorAccountID, func() hedera.SecretKey {
		return operatorSecret
	})

	defer client.Close()

	//
	// Get balance for target account
	//

	balance, err := client.Account(targetAccountID).Balance().Get()
	if err != nil {
		panic(err)
	}

	fmt.Printf("account balance = %v\n", balance)

	//
	// Transfer 100 cryptos to target
	//

	nodeAccountID := hedera.AccountID{Account: 3}
	response, err := client.TransferCrypto(). //creates a transaction to transfer hbars between accounts.
		// Move 100 out of operator account
		Transfer(operatorAccountID, -100). //sets up a transfer, which pairs an account with a signed integer. In this case, the account is your account and the amount is -1. The negative number indicates that the balance of your account will be decremented by this amount.
		// And place in our new account
		Transfer(targetAccountID, 100). //creates a second transfer, pairing an account with a signed integer. In this case, the account is your friend’s account and the amount is 1. The positive number indicates that the balance of your account will be incremented by this amount. 
		//Important: the sum of all transfers contained within in a CryptoTransfer must equal zero.
		Operator(operatorAccountID). // identifies the account initiating the transaction.
		Node(nodeAccountID). //identifies the account of the Hedea node to which the transaction is being sent.
		Memo("[test] hedera-sdk-go v2").
		//add a signature based on a secret key. It is necessary to repeat this line to sign as both operator initiating the transfer transaction and account holder associated with an outgoing (negative) transfer – even though both keys are the same.
		Sign(operatorSecret). // Sign it once as operator
		Sign(operatorSecret). // And again as sender
		Execute() // executes the transaction.

	if err != nil {
		panic(err)
	}

	transactionID := response.ID //transactionID is made up of the account ID and the transaction timestamp right down to nanoseconds
	fmt.Printf("transferred; transaction AccountID: %v, transaction time: %v\n", transactionID.AccountID, transactionID.TransactionValidStart)

	//
	// Get receipt to prove we sent ok
	// Although this is not a mandatory step, it does verify that your transaction successfully reached network consensus.
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
	fmt.Printf("transaction has a successful status: %v, on contract: %v\n",receipt.Status.String(),receipt.ContractID)
	fmt.Printf("wait for 2s...\n")
	time.Sleep(2 * time.Second)

	//
	// Get balance for target account (again)
	//

	balance, err = client.Account(targetAccountID).Balance().Get()
	if err != nil {
		panic(err)
	}

	fmt.Printf("account balance = %v\n", balance)
}