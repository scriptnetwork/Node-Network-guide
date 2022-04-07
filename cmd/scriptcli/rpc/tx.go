package rpc

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/spf13/viper"
	rpcc "github.com/ybbus/jsonrpc"

	"github.com/scripttoken/script/cmd/scriptcli/cmd/utils"
	"github.com/scripttoken/script/common"
	"github.com/scripttoken/script/core"
	"github.com/scripttoken/script/ledger/types"
	trpc "github.com/scripttoken/script/rpc"
)

// ------------------------------- SendTx -----------------------------------

type SendArgs struct {
	ChainID  string `json:"chain_id"`
	From     string `json:"from"`
	To       string `json:"to"`
	SCPTWei  string `json:"SCPTWei"`
	SPAYWei  string `json:"SPAYWei"`
	Fee      string `json:"fee"`
	Sequence string `json:"sequence"`
	Async    bool   `json:"async"`
}

type SendResult struct {
	TxHash string            `json:"hash"`
	Block  *core.BlockHeader `json:"block",rlp:"nil"`
}

func (t *scriptcliRPCService) Send(args *SendArgs, result *SendResult) (err error) {
	if len(args.From) == 0 || len(args.To) == 0 {
		return fmt.Errorf("The from and to address cannot be empty")
	}
	if args.From == args.To {
		return fmt.Errorf("The from and to address cannot be identical")
	}

	from := common.HexToAddress(args.From)
	to := common.HexToAddress(args.To)
	SCPTWei, ok := new(big.Int).SetString(args.SCPTWei, 10)
	if !ok {
		return fmt.Errorf("Failed to parse SCPTWei: %v", args.SCPTWei)
	}
	SPAYWei, ok := new(big.Int).SetString(args.SPAYWei, 10)
	if !ok {
		return fmt.Errorf("Failed to parse SPAYWei: %v", args.SPAYWei)
	}
	fee, ok := new(big.Int).SetString(args.Fee, 10)
	if !ok {
		return fmt.Errorf("Failed to parse fee: %v", args.Fee)
	}
	sequence, err := strconv.ParseUint(args.Sequence, 10, 64)
	if err != nil {
		return err
	}

	if !t.wallet.IsUnlocked(from) {
		return fmt.Errorf("The from address %v has not been unlocked yet", from.Hex())
	}

	inputs := []types.TxInput{{
		Address: from,
		Coins: types.Coins{
			SPAYWei: new(big.Int).Add(SPAYWei, fee),
			SCPTWei: SCPTWei,
		},
		Sequence: sequence,
	}}
	outputs := []types.TxOutput{{
		Address: to,
		Coins: types.Coins{
			SPAYWei: SPAYWei,
			SCPTWei: SCPTWei,
		},
	}}
	sendTx := &types.SendTx{
		Fee: types.Coins{
			SCPTWei: new(big.Int).SetUint64(0),
			SPAYWei: fee,
		},
		Inputs:  inputs,
		Outputs: outputs,
	}

	signBytes := sendTx.SignBytes(args.ChainID)
	sig, err := t.wallet.Sign(from, signBytes)
	if err != nil {
		utils.Error("Failed to sign transaction: %v\n", err)
	}
	sendTx.SetSignature(from, sig)

	raw, err := types.TxToBytes(sendTx)
	if err != nil {
		utils.Error("Failed to encode transaction: %v\n", err)
	}
	signedTx := hex.EncodeToString(raw)

	client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

	rpcMethod := "script.BroadcastRawTransaction"
	if args.Async {
		rpcMethod = "script.BroadcastRawTransactionAsync"
	}
	res, err := client.Call(rpcMethod, trpc.BroadcastRawTransactionArgs{TxBytes: signedTx})
	if err != nil {
		return err
	}
	if res.Error != nil {
		return fmt.Errorf("Server returned error: %v", res.Error)
	}
	trpcResult := &trpc.BroadcastRawTransactionResult{}
	err = res.GetObject(trpcResult)
	if err != nil {
		return fmt.Errorf("Failed to parse Script node response: %v", err)
	}

	result.TxHash = trpcResult.TxHash
	result.Block = trpcResult.Block

	return nil
}

func (t *scriptcliRPCService) EdgeStake(args *SendArgs, result *SendResult) (err error) {
	if len(args.From) == 0 || len(args.To) == 0 {
		return fmt.Errorf("The from and to address cannot be empty")
	}
	if args.From == args.To {
		return fmt.Errorf("The from and to address cannot be identical")
	}

	from := common.HexToAddress(args.From)
	to := common.HexToAddress(args.To)
	SCPTWei, ok := new(big.Int).SetString(args.SCPTWei, 10)
	if !ok {
		return fmt.Errorf("Failed to parse SCPTWei: %v", args.SCPTWei)
	}
	SPAYWei, ok := new(big.Int).SetString(args.SPAYWei, 10)
	if !ok {
		return fmt.Errorf("Failed to parse SPAYWei: %v", args.SPAYWei)
	}
	fee, ok := new(big.Int).SetString(args.Fee, 10)
	if !ok {
		return fmt.Errorf("Failed to parse fee: %v", args.Fee)
	}
	sequence, err := strconv.ParseUint(args.Sequence, 10, 64)
	if err != nil {
		return err
	}

	if !t.wallet.IsUnlocked(from) {
		return fmt.Errorf("The from address %v has not been unlocked yet", from.Hex())
	}

	inputs := []types.TxInput{{
		Address: from,
		Coins: types.Coins{
			SPAYWei: new(big.Int).Add(SPAYWei, fee),
			SCPTWei: SCPTWei,
		},
		Sequence: sequence,
	}}
	outputs := []types.TxOutput{{
		Address: to,
		Coins: types.Coins{
			SPAYWei: SPAYWei,
			SCPTWei: SCPTWei,
		},
	}}
	edgeStakeTx := &types.EdgeStakeTx{
		Fee: types.Coins{
			SCPTWei: new(big.Int).SetUint64(0),
			SPAYWei: fee,
		},
		Inputs:  inputs,
		Outputs: outputs,
	}

	signBytes := edgeStakeTx.SignBytes(args.ChainID)
	sig, err := t.wallet.Sign(from, signBytes)
	if err != nil {
		utils.Error("Failed to sign transaction: %v\n", err)
	}
	edgeStakeTx.SetSignature(from, sig)

	raw, err := types.TxToBytes(edgeStakeTx)
	if err != nil {
		utils.Error("Failed to encode transaction: %v\n", err)
	}
	signedTx := hex.EncodeToString(raw)

	client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

	rpcMethod := "script.BroadcastRawTransaction"
	if args.Async {
		rpcMethod = "script.BroadcastRawTransactionAsync"
	}
	res, err := client.Call(rpcMethod, trpc.BroadcastRawTransactionArgs{TxBytes: signedTx})
	if err != nil {
		return err
	}
	if res.Error != nil {
		return fmt.Errorf("Server returned error: %v", res.Error)
	}
	trpcResult := &trpc.BroadcastRawTransactionResult{}
	err = res.GetObject(trpcResult)
	if err != nil {
		return fmt.Errorf("Failed to parse Script node response: %v", err)
	}

	result.TxHash = trpcResult.TxHash
	result.Block = trpcResult.Block

	return nil
}
