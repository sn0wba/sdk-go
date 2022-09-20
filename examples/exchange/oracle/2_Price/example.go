package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gotabit/sdk-go/client/common"
	exchangeclient "github.com/gotabit/sdk-go/client/exchange"
)

func main() {
	network := common.LoadNetwork("mainnet", "lb")
	exchangeClient, err := exchangeclient.NewExchangeClient(network.ExchangeGrpcEndpoint, common.OptionTLSCert(network.ExchangeTlsCert))
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()
	baseSymbol := "BTC"
	quoteSymbol := "USDT"
	oracleType := "BandIBC"
	oracleScaleFactor := uint32(6)
	res, err := exchangeClient.GetPrice(ctx, baseSymbol, quoteSymbol, oracleType, oracleScaleFactor)
	if err != nil {
		fmt.Println(err)
	}

	str, _ := json.MarshalIndent(res, "", " ")
	fmt.Print(string(str))
}
