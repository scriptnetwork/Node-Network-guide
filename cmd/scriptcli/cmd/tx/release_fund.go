package tx

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/scripttoken/script/cmd/scriptcli/cmd/utils"
	"github.com/scripttoken/script/ledger/types"
	"github.com/scripttoken/script/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	rpcc "github.com/ybbus/jsonrpc"
)

// releaseFundCmd represents the release fund command
// Example:
//
//	scriptcli tx release --chain="privatenet" --from=2E833968E5bB786Ae419c4d13189fB081Cc43bab  --reserve_seq=8 --seq=8
var releaseFundCmd = &cobra.Command{
	Use:     "release",
	Short:   "Release fund",
	Example: `scriptcli tx release --chain="privatenet" --from=2E833968E5bB786Ae419c4d13189fB081Cc43bab  --reserve_seq=8 --seq=8`,
	Run:     doReleaseFundCmd,
}

func doReleaseFundCmd(cmd *cobra.Command, args []string) {
	wallet, fromAddress, err := walletUnlock(cmd, fromFlag, passwordFlag)
	if err != nil {
		return
	}
	defer wallet.Lock(fromAddress)

	input := types.TxInput{
		Address:  fromAddress,
		Sequence: uint64(seqFlag),
	}

	spay, ok := types.ParseCoinAmount(feeFlag)
	if !ok {
		utils.Error("Failed to parse spay amount")
	}
	releaseFundTx := &types.ReleaseFundTx{
		Fee: types.Coins{
			SCPTWei: new(big.Int).SetUint64(0),
			SPAYWei: spay,
		},
		Source:          input,
		ReserveSequence: reserveSeqFlag,
	}

	sig, err := wallet.Sign(fromAddress, releaseFundTx.SignBytes(chainIDFlag))
	if err != nil {
		utils.Error("Failed to sign transaction: %v\n", err)
	}
	releaseFundTx.SetSignature(fromAddress, sig)

	raw, err := types.TxToBytes(releaseFundTx)
	if err != nil {
		utils.Error("Failed to encode transaction: %v\n", err)
	}
	signedTx := hex.EncodeToString(raw)

	client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

	var res *rpcc.RPCResponse
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
	fmt.Printf("Successfully broadcasted transaction.\n")
}

func init() {
	releaseFundCmd.Flags().StringVar(&chainIDFlag, "chain", "", "Chain ID")
	releaseFundCmd.Flags().StringVar(&fromFlag, "from", "", "Reserve owner's address")
	releaseFundCmd.Flags().Uint64Var(&seqFlag, "seq", 0, "Sequence number of the transaction")
	releaseFundCmd.Flags().StringVar(&feeFlag, "fee", fmt.Sprintf("%dwei", types.MinimumTransactionFeeSPAYWeiJune2021), "Fee")
	releaseFundCmd.Flags().Uint64Var(&reserveSeqFlag, "reserve_seq", 1000, "Reserve sequence")
	releaseFundCmd.Flags().StringVar(&walletFlag, "wallet", "soft", "Wallet type (soft|nano)")
	releaseFundCmd.Flags().BoolVar(&asyncFlag, "async", false, "block until tx has been included in the blockchain")
	releaseFundCmd.Flags().StringVar(&passwordFlag, "password", "", "password to unlock the wallet")

	releaseFundCmd.MarkFlagRequired("chain")
	releaseFundCmd.MarkFlagRequired("from")
	releaseFundCmd.MarkFlagRequired("seq")
	releaseFundCmd.MarkFlagRequired("reserve_seq")
	releaseFundCmd.MarkFlagRequired("resource_id")

}
