package tx

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/scripttoken/script/cmd/scriptcli/cmd/utils"
	"github.com/scripttoken/script/common"
	"github.com/scripttoken/script/ledger/types"
	"github.com/scripttoken/script/rpc"
	wtypes "github.com/scripttoken/script/wallet/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ybbus/jsonrpc"
	rpcc "github.com/ybbus/jsonrpc"
)

// edgeStakeCmd represents the edgeStake command
// Example:
//		scriptcli tx edgeStake --chain="scriptnet" --from=98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5 --to=98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5 --script=10 --spay=9 --seq=1
//		scriptcli tx edgeStake --chain="scriptnet" --path "m/44'/60'/0'/0/0" --to=98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5 --script=10 --spay=9 --seq=1 --wallet=trezor
//		scriptcli tx edgeStake --chain="scriptnet" --path "m/44'/60'/0'/0" --to=98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5 --script=10 --spay=9 --seq=1 --wallet=nano
var edgeStakeCmd = &cobra.Command{
	Use:     "edgeStake",
	Short:   "EdgeStake tokens",
	Example: `scriptcli tx edgeStake --chain="scriptnet" --from=98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5 --to=98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5 --script=10 --spay=9 --seq=1`,
	Run:     doEdgeStakeCmd,
}

func doEdgeStakeCmd(cmd *cobra.Command, args []string) {
	walletType := getWalletType(cmd)
	if walletType == wtypes.WalletTypeSoft && len(fromFlag) == 0 {
		utils.Error("The from address cannot be empty") // we don't need to specify the "from address" for hardware wallets
		return
	}

	if len(toFlag) == 0 {
		utils.Error("The to address cannot be empty")
		return
	}
	if fromFlag == toFlag {
		utils.Error("The from and to address cannot be identical")
		return
	}

	wallet, fromAddress, err := walletUnlockWithPath(cmd, fromFlag, pathFlag)
	if err != nil || wallet == nil {
		return
	}
	defer wallet.Lock(fromAddress)

	script, ok := types.ParseCoinAmount(scriptAmountFlag)
	if !ok {
		utils.Error("Failed to parse script amount")
	}
	spay, ok := types.ParseCoinAmount(spayAmountFlag)
	if !ok {
		utils.Error("Failed to parse spay amount")
	}
	fee, ok := types.ParseCoinAmount(feeFlag)
	if !ok {
		utils.Error("Failed to parse fee")
	}
	inputs := []types.TxInput{{
		Address: fromAddress,
		Coins: types.Coins{
			SPAYWei: new(big.Int).Add(spay, fee),
			SCPTWei: script,
		},
		Sequence: uint64(seqFlag),
	}}
	outputs := []types.TxOutput{{
		Address: common.HexToAddress(toFlag),
		Coins: types.Coins{
			SPAYWei: spay,
			SCPTWei: script,
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

	sig, err := wallet.Sign(fromAddress, edgeStakeTx.SignBytes(chainIDFlag))
	if err != nil {
		utils.Error("Failed to sign transaction: %v\n", err)
	}
	edgeStakeTx.SetSignature(fromAddress, sig)

	raw, err := types.TxToBytes(edgeStakeTx)
	if err != nil {
		utils.Error("Failed to encode transaction: %v\n", err)
	}
	signedTx := hex.EncodeToString(raw)

	client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

	var res *jsonrpc.RPCResponse
	if asyncFlag {
		res, err = client.Call("script.BroadcastRawTransactionAsync", rpc.BroadcastRawTransactionArgs{TxBytes: signedTx})
	} else {
		res, err = client.Call("script.BroadcastRawTransaction", rpc.BroadcastRawTransactionArgs{TxBytes: signedTx})
	}

	if err != nil {
		utils.Error("Failed to broadcast transaction: %v\n", err)
	}
	if res.Error != nil {
		utils.Error("Server returned error: %v\n", res.Error)
	}
	result := &rpc.BroadcastRawTransactionResult{}
	err = res.GetObject(result)
	if err != nil {
		utils.Error("Failed to parse server response: %v\n", err)
	}
	formatted, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		utils.Error("Failed to parse server response: %v\n", err)
	}
	fmt.Printf("Successfully broadcasted transaction:\n%s\n", formatted)
}

func init() {
	edgeStakeCmd.Flags().StringVar(&chainIDFlag, "chain", "", "Chain ID")
	edgeStakeCmd.Flags().StringVar(&fromFlag, "from", "", "Address to edgeStake from")
	edgeStakeCmd.Flags().StringVar(&toFlag, "to", "", "Address to edgeStake to")
	edgeStakeCmd.Flags().StringVar(&pathFlag, "path", "", "Wallet derivation path")
	edgeStakeCmd.Flags().Uint64Var(&seqFlag, "seq", 0, "Sequence number of the transaction")
	edgeStakeCmd.Flags().StringVar(&scriptAmountFlag, "script", "0", "Script amount")
	edgeStakeCmd.Flags().StringVar(&spayAmountFlag, "spay", "0", "Script amount")
	edgeStakeCmd.Flags().StringVar(&feeFlag, "fee", fmt.Sprintf("%dwei", types.MinimumTransactionFeeSPAYWei), "Fee")
	edgeStakeCmd.Flags().StringVar(&walletFlag, "wallet", "soft", "Wallet type (soft|nano|trezor)")
	edgeStakeCmd.Flags().BoolVar(&asyncFlag, "async", false, "block until tx has been included in the blockchain")

	edgeStakeCmd.MarkFlagRequired("chain")
	//edgeStakeCmd.MarkFlagRequired("from")
	edgeStakeCmd.MarkFlagRequired("to")
	edgeStakeCmd.MarkFlagRequired("seq")
}
