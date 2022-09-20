package main

import (
	"context"
	"encoding/json"
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
	denom := "inj"
	subaccountId := "0xaf79152ac5df276d9a8e1e2e22822f9713474902000000000000000000000000"
	transferTypes := []string{"deposit"}
	skip := uint64(0)
	limit := int32(10)

	req := accountPB.SubaccountHistoryRequest{
		Denom:         denom,
		SubaccountId:  subaccountId,
		TransferTypes: transferTypes,
		Skip:          skip,
		Limit:         limit,
	}

	res, err := exchangeClient.GetSubaccountHistory(ctx, req)
	if err != nil {
		fmt.Println(err)
	}

	str, _ := json.MarshalIndent(res, "", " ")
	fmt.Print(string(str))
}
