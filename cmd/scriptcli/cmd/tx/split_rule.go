package tx

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/scripttoken/script/cmd/scriptcli/cmd/utils"
	"github.com/scripttoken/script/common"
	"github.com/scripttoken/script/ledger/types"
	"github.com/scripttoken/script/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	rpcc "github.com/ybbus/jsonrpc"
)

// splitRuleCmd represents the split rule command
// Example:
//
//	scriptcli tx split_rule --chain="scriptnet" --from=2E833968E5bB786Ae419c4d13189fB081Cc43bab --seq=8 --resource_id=die_another_day --addresses=2E833968E5bB786Ae419c4d13189fB081Cc43bab,9F1233798E905E173560071255140b4A8aBd3Ec6 --percentages=30,30 --duration=1000
var splitRuleCmd = &cobra.Command{
	Use:     "split_rule",
	Short:   "Initiate or update a split rule",
	Example: `scriptcli tx split_rule --chain="scriptnet" --from=2E833968E5bB786Ae419c4d13189fB081Cc43bab --seq=8 --resource_id=die_another_day --addresses=2E833968E5bB786Ae419c4d13189fB081Cc43bab,9F1233798E905E173560071255140b4A8aBd3Ec6 --percentages=30,30 --duration=1000`,
	Run:     doSplitRuleCmd,
}

func doSplitRuleCmd(cmd *cobra.Command, args []string) {
	wallet, fromAddress, err := walletUnlock(cmd, fromFlag, passwordFlag)
	if err != nil {
		return
	}
	defer wallet.Lock(fromAddress)

	input := types.TxInput{
		Address:  fromAddress,
		Sequence: uint64(seqFlag),
	}

	if len(addressesFlag) != len(percentagesFlag) {
		fmt.Println("Should have the same number of addresses and percentages")
		return
	}
	var splits []types.Split
	for idx, addressStr := range addressesFlag {
		percentageStr := percentagesFlag[idx]

		address, err := hex.DecodeString(addressStr)
		if err != nil {
			fmt.Println("The address must be a hex string")
			return
		}

		percentage, err := strconv.ParseUint(percentageStr, 10, 32)
		if err != nil {
			fmt.Println(err)
			return
		}

		split := types.Split{
			Address:    common.BytesToAddress(address),
			Percentage: uint(percentage),
		}
		splits = append(splits, split)
	}

	fee, ok := types.ParseCoinAmount(feeFlag)
	if !ok {
		utils.Error("Failed to parse fee")
	}

	splitRuleTx := &types.SplitRuleTx{
		Fee: types.Coins{
			SCPTWei: new(big.Int).SetUint64(0),
			SPAYWei: fee,
		},
		ResourceID: resourceIDFlag,
		Initiator:  input,
		Duration:   durationFlag,
		Splits:     splits,
	}

	sig, err := wallet.Sign(fromAddress, splitRuleTx.SignBytes(chainIDFlag))
	if err != nil {
		utils.Error("Failed to sign transaction: %v\n", err)
	}
	splitRuleTx.SetSignature(fromAddress, sig)

	raw, err := types.TxToBytes(splitRuleTx)
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
	splitRuleCmd.Flags().StringVar(&chainIDFlag, "chain", "", "Chain ID")
	splitRuleCmd.Flags().StringVar(&fromFlag, "from", "", "Initiator's address")
	splitRuleCmd.Flags().Uint64Var(&seqFlag, "seq", 0, "Sequence number of the transaction")
	splitRuleCmd.Flags().StringVar(&feeFlag, "fee", fmt.Sprintf("%dwei", types.MinimumTransactionFeeSPAYWeiJune2021), "Fee")
	splitRuleCmd.Flags().StringVar(&resourceIDFlag, "resource_id", "", "The resourceID of interest")
	splitRuleCmd.Flags().StringSliceVar(&addressesFlag, "addresses", []string{}, "List of addresses participating in the split")
	splitRuleCmd.Flags().StringSliceVar(&percentagesFlag, "percentages", []string{}, "List of integers (between 0 and 100) representing of percentage of split")
	splitRuleCmd.Flags().Uint64Var(&durationFlag, "duration", 1000, "Reserve duration")
	splitRuleCmd.Flags().StringVar(&walletFlag, "wallet", "soft", "Wallet type (soft|nano)")
	splitRuleCmd.Flags().BoolVar(&asyncFlag, "async", false, "block until tx has been included in the blockchain")
	splitRuleCmd.Flags().StringVar(&passwordFlag, "password", "", "password to unlock the wallet")

	splitRuleCmd.MarkFlagRequired("chain")
	splitRuleCmd.MarkFlagRequired("from")
	splitRuleCmd.MarkFlagRequired("seq")
	splitRuleCmd.MarkFlagRequired("addresses")
	splitRuleCmd.MarkFlagRequired("percentages")
	splitRuleCmd.MarkFlagRequired("resource_id")
	splitRuleCmd.MarkFlagRequired("duration")
}
