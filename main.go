package main

import (
	"fmt"
	"time"
	"github.com/hashgraph/hedera-sdk-go"
)


func main() {
 	 
 	//Establish connection to the Hedera node
 	client := hedera.Dial("testnet.hedera.com:50123")
	if err != nil {
	    panic(err)
	}
	//Defer the disconnection of the connection to guarantee a clean disconnect from the node.
	defer client.Close()
	// 0.0.1003
	myAccount := hedera.NewAccountID(0, 0, 1003)

	//GetAccountBalance constructs the request; adding .Answer() executes the request. 

	myBalance, err := client.GetAccountBalance(myAccount).Answer()
 	//Don’t forget error-handling.
 	if err != nil {
    	panic(err)
	}

	fmt.Printf("Your balance: %v \n", myBalance)
	// testnets are throttled and may respond with an error if requests are made to the Hedera API to frequently. For this reason it is best practice to add a sleep timer inbetween requests.
	//time.Sleep(1 * time.Second)
}