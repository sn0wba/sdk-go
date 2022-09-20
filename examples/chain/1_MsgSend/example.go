package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gotabit/sdk-go/client/common"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	chainclient "github.com/gotabit/sdk-go/client/chain"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

func main() {
	// network := common.LoadNetwork("mainnet")
	network := common.LoadNetwork("testnet")
	rpc, err := rpchttp.New(network.GrpcEndpoint, "/websocket")

	senderAddress, cosmosKeyring, err := chainclient.InitCosmosKeyring(
		os.Getenv("HOME")+"/.gotabitd",
		"gotabitd",
		"file",
		"user",
		"",
		"actual accuse plastic supply favorite banner trial company cloud wasp enable cactus",
		false,
	)

	if err != nil {
		panic(err)
	}

	// initialize grpc client

	clientCtx, err := chainclient.NewClientContext(
		network.ChainId,
		senderAddress.String(),
		cosmosKeyring,
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	clientCtx = clientCtx.WithNodeURI(network.GrpcEndpoint).WithClient(rpc)

	// prepare tx msg

	msg := &banktypes.MsgSend{
		FromAddress: senderAddress.String(),
		ToAddress:   "gio1rmwwe2hge0lwzcu6ftmr547xsgjewqtdhehgte",
		Amount: []sdktypes.Coin{{
			Denom: "ugtb", Amount: sdktypes.NewInt(600000000)}, // 100 ugtb
		},
	}

	chainClient, err := chainclient.NewChainClient(
		clientCtx,
		network.GrpcEndpoint,
		common.OptionGasPrices("500000000ugtb"),
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	//AsyncBroadcastMsg, SyncBroadcastMsg, QueueBroadcastMsg
	err = chainClient.QueueBroadcastMsg(msg)

	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(time.Second * 5)

	gasFee, err := chainClient.GetGasFee()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("gas fee:", gasFee, "ugtb")
}
