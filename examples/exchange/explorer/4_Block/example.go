package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gotabit/sdk-go/client/common"
	exchangeclient "github.com/gotabit/sdk-go/client/exchange"
)

func main() {
	//network := common.LoadNetwork("mainnet", "k8s")
	network := common.LoadNetwork("testnet", "k8s")
	exchangeClient, err := exchangeclient.NewExchangeClient(network.ExchangeGrpcEndpoint, common.OptionTLSCert(network.ExchangeTlsCert))
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()
	blockHeight := "5825046"
	res, err := exchangeClient.GetBlock(ctx, blockHeight)
	if err != nil {
		fmt.Println(err)
	}

	str, _ := json.MarshalIndent(res, "", " ")
	fmt.Print(string(str))
}
