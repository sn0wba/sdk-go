package main

import (
	"context"
	"fmt"

	"github.com/gotabit/sdk-go/client/common"
	exchangeclient "github.com/gotabit/sdk-go/client/exchange"
	accountPB "github.com/gotabit/sdk-go/exchange/accounts_rpc/pb"
)

func main() {
	//network := common.LoadNetwork("mainnet", "k8s")
	network := common.LoadNetwork("testnet", "k8s")
	exchangeClient, err := exchangeclient.NewExchangeClient(network.ExchangeGrpcEndpoint, common.OptionTLSCert(network.ExchangeTlsCert))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	marketId := "0xa508cb32923323679f29a032c70342c147c17d0145625922b0ef22e955c844c0"
	subaccountId := "0xaf79152ac5df276d9a8e1e2e22822f9713474902000000000000000000000000"
	orderDirection := "buy"

	req := accountPB.SubaccountOrderSummaryRequest{
		MarketId:       marketId,
		SubaccountId:   subaccountId,
		OrderDirection: orderDirection,
	}

	res, err := exchangeClient.GetSubaccountOrderSummary(ctx, req)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("spot orders:", res.SpotOrdersTotal)
	fmt.Println("derivative orders:", res.DerivativeOrdersTotal)
}
