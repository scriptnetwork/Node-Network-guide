package tx

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/scripttoken/script/cmd/scriptcli/cmd/utils"
	"github.com/scripttoken/script/common"
	"github.com/scripttoken/script/ledger/types"
	"github.com/scripttoken/script/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	rpcc "github.com/ybbus/jsonrpc"
)

// withdrawStakeCmd represents the withdraw stake command
// Example:
//		scriptcli tx withdraw --chain="scriptnet" --source=98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5 --holder=98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5 --purpose=0 --seq=8
var withdrawStakeCmd = &cobra.Command{
	Use:     "withdraw",
	Short:   "withdraw stake to a validator or guardian",
	Example: `scriptcli tx withdraw --chain="scriptnet" --source=98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5 --holder=98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5 --purpose=0 --seq=8`,
	Run:     doWithdrawStakeCmd,
}

func doWithdrawStakeCmd(cmd *cobra.Command, args []string) {
	wallet, sourceAddress, err := walletUnlockWithPath(cmd, sourceFlag, pathFlag)
	if err != nil {
		return
	}
	defer wallet.Lock(sourceAddress)

	fee, ok := types.ParseCoinAmount(feeFlag)
	if !ok {
		utils.Error("Failed to parse fee")
	}

	source := types.TxInput{
		Address:  sourceAddress,
		Sequence: uint64(seqFlag),
	}
	holder := types.TxOutput{
		Address: common.HexToAddress(holderFlag),
	}

	withdrawStakeTx := &types.WithdrawStakeTx{
		Fee: types.Coins{
			SCPTWei: new(big.Int).SetUint64(0),
			SPAYWei: fee,
		},
		Source:  source,
		Holder:  holder,
		Purpose: purposeFlag,
	}

	sig, err := wallet.Sign(sourceAddress, withdrawStakeTx.SignBytes(chainIDFlag))
	if err != nil {
		utils.Error("Failed to sign transaction: %v\n", err)
	}
	withdrawStakeTx.SetSignature(sourceAddress, sig)

	raw, err := types.TxToBytes(withdrawStakeTx)
	if err != nil {
		utils.Error("Failed to encode transaction: %v\n", err)
	}
	signedTx := hex.EncodeToString(raw)

	client := rpcc.NewRPCClient(viper.GetString(utils.CfgRemoteRPCEndpoint))

	res, err := client.Call("script.BroadcastRawTransaction", rpc.BroadcastRawTransactionArgs{TxBytes: signedTx})
	if err != nil {
		utils.Error("Failed to broadcast transaction: %v\n", err)
	}
	if res.Error != nil {
		utils.Error("Server returned error: %v\n", res.Error)
	}
	fmt.Printf("Successfully broadcasted transaction.\n")
}

func init() {
	withdrawStakeCmd.Flags().StringVar(&chainIDFlag, "chain", "", "Chain ID")
	withdrawStakeCmd.Flags().StringVar(&sourceFlag, "source", "", "Source of the stake")
	withdrawStakeCmd.Flags().StringVar(&holderFlag, "holder", "", "Holder of the stake")
	withdrawStakeCmd.Flags().StringVar(&pathFlag, "path", "", "Wallet derivation path")
	withdrawStakeCmd.Flags().StringVar(&feeFlag, "fee", fmt.Sprintf("%dwei", types.MinimumTransactionFeeSPAYWei), "Fee")
	withdrawStakeCmd.Flags().Uint64Var(&seqFlag, "seq", 0, "Sequence number of the transaction")
	withdrawStakeCmd.Flags().Uint8Var(&purposeFlag, "purpose", 0, "Purpose of staking")
	withdrawStakeCmd.Flags().StringVar(&walletFlag, "wallet", "soft", "Wallet type (soft|nano)")

	withdrawStakeCmd.MarkFlagRequired("chain")
	withdrawStakeCmd.MarkFlagRequired("source")
	withdrawStakeCmd.MarkFlagRequired("holder")
	withdrawStakeCmd.MarkFlagRequired("seq")
}
