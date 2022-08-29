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

// stakeRewardDistributionCmd represents the stake reward distribution command
// Example:
//
//	scriptcli tx distribute_staking_reward --chain="privatenet" --holder=0x36A8d78C0EaD519Bd155962358A3d57A404bC20d --beneficiary=0x88884a84d980bbfb7588888126fb903486bb8888 --split_basis_point=100 --seq=8
var stakeRewardDistributionCmd = &cobra.Command{
	Use:     "distribute_staking_reward",
	Short:   "Configure the distribution of the guardian/elite edge node staking reward",
	Example: `scriptcli tx distribute_staking_reward --chain="privatenet" --holder=0x36A8d78C0EaD519Bd155962358A3d57A404bC20d --beneficiary=0x88884a84d980bbfb7588888126fb903486bb8888 --split_basis_point=100 --seq=8`,
	Run:     doStakeRewardDistributionCmd,
}

func doStakeRewardDistributionCmd(cmd *cobra.Command, args []string) {
	wallet, holderAddress, err := walletUnlockWithPath(cmd, holderFlag, pathFlag, passwordFlag)
	if err != nil {
		return
	}
	defer wallet.Lock(holderAddress)

	fee, ok := types.ParseCoinAmount(feeFlag)
	if !ok {
		utils.Error("Failed to parse fee")
	}

	holder := types.TxInput{
		Address:  holderAddress,
		Sequence: uint64(seqFlag),
	}
	beneficiary := types.TxOutput{
		Address: common.HexToAddress(beneficiaryFlag),
	}

	stakeRewardDistributionTx := &types.StakeRewardDistributionTx{
		Fee: types.Coins{
			SCPTWei: new(big.Int).SetUint64(0),
			SPAYWei: fee,
		},
		Holder:          holder,
		Beneficiary:     beneficiary,
		SplitBasisPoint: uint(splitBasisPointFlag),
		//Purpose:         purposeFlag,
	}

	sig, err := wallet.Sign(holderAddress, stakeRewardDistributionTx.SignBytes(chainIDFlag))
	if err != nil {
		utils.Error("Failed to sign transaction: %v\n", err)
	}
	stakeRewardDistributionTx.SetSignature(holderAddress, sig)

	raw, err := types.TxToBytes(stakeRewardDistributionTx)
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
	stakeRewardDistributionCmd.Flags().StringVar(&chainIDFlag, "chain", "", "Chain ID")
	stakeRewardDistributionCmd.Flags().StringVar(&holderFlag, "holder", "", "Holder of the stake")
	stakeRewardDistributionCmd.Flags().StringVar(&pathFlag, "path", "", "Wallet derivation path")
	stakeRewardDistributionCmd.Flags().StringVar(&feeFlag, "fee", fmt.Sprintf("%dwei", types.MinimumTransactionFeeSPAYWei), "Fee")
	stakeRewardDistributionCmd.Flags().Uint64Var(&seqFlag, "seq", 0, "Sequence number of the transaction")
	stakeRewardDistributionCmd.Flags().StringVar(&beneficiaryFlag, "beneficiary", "", "Address of the beneficiary")
	stakeRewardDistributionCmd.Flags().Uint64Var(&splitBasisPointFlag, "split_basis_point", 0, "fraction of the reward split in terms of basis point (1/10000). 100 basis point = 100/10000 = 1.00%")
	//stakeRewardDistributionCmd.Flags().Uint8Var(&purposeFlag, "purpose", 0, "Purpose of staking")
	stakeRewardDistributionCmd.Flags().StringVar(&walletFlag, "wallet", "soft", "Wallet type (soft|nano)")
	stakeRewardDistributionCmd.Flags().BoolVar(&asyncFlag, "async", false, "block until tx has been included in the blockchain")
	stakeRewardDistributionCmd.Flags().StringVar(&passwordFlag, "password", "", "password to unlock the wallet")

	stakeRewardDistributionCmd.MarkFlagRequired("chain")
	stakeRewardDistributionCmd.MarkFlagRequired("holder")
	stakeRewardDistributionCmd.MarkFlagRequired("seq")
}
