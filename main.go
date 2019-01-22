package main

import (
	"fmt"
	"os"
	"github.com/hashgraph/hedera-sdk-go"
)


func main() {
 	// Target account to get the balance for
	accountID := hedera.AccountID{Account: 1003}

	client, err := hedera.Dial("testnet.hedera.com:50119")
	if err != nil {
		panic(err)
	}

	client.SetNode(hedera.AccountID{Account: 3})
	
	client.SetOperator(accountID, func() hedera.SecretKey {
		operatorSecret, err := hedera.SecretKeyFromString(os.Getenv("SECRET_KEY"))
		if err != nil {
			panic(err)
		}

		return operatorSecret
	})

	defer client.Close()

	accountTwo := hedera.AccountID{Account: 1004}
	// Get the _answer_ for the query of getting the account balance
	balance, err := client.Account(accountTwo).Balance().Get()
	if err != nil {
		panic(err)
	}

	fmt.Printf("balance = %v tinybars\n", balance)
	//fmt.Printf("balance = %.5f hbars\n", float64(balance)/100000000.0)
	
	// testnets are throttled and may respond with an error if requests are made to the Hedera API to frequently. For this reason it is best practice to add a sleep timer inbetween requests.
	//time.Sleep(1 * time.Second)
}