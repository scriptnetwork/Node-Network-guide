package tx

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/scripttoken/script/crypto"

	"github.com/scripttoken/script/crypto/bls"

	"github.com/scripttoken/script/cmd/scriptcli/cmd/utils"
	"github.com/scripttoken/script/common"
	"github.com/scripttoken/script/core"
	"github.com/scripttoken/script/ledger/types"
	"github.com/scripttoken/script/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	rpcc "github.com/ybbus/jsonrpc"
)

// depositStakeCmd represents the deposit stake command
// Example:
//		scriptcli tx deposit --chain="scriptnet" --source=98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5 --holder=98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5 --stake=6000000 --purpose=0 --seq=7
var depositStakeCmd = &cobra.Command{
	Use:     "deposit",
	Short:   "Deposit stake to a validator or guardian",
	Example: `scriptcli tx deposit --chain="scriptnet" --source=98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5 --holder=98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5 --stake=6000000 --purpose=0 --seq=7`,
	Run:     doDepositStakeCmd,
}

func doDepositStakeCmd(cmd *cobra.Command, args []string) {
	wallet, sourceAddress, err := walletUnlockWithPath(cmd, sourceFlag, pathFlag)
	if err != nil {
		return
	}
	defer wallet.Lock(sourceAddress)

	fee, ok := types.ParseCoinAmount(feeFlag)
	if !ok {
		utils.Error("Failed to parse fee")
	}
	stake, ok := types.ParseCoinAmount(stakeInScriptFlag)
	if !ok {
		utils.Error("Failed to parse stake")
	}
	if stake.Cmp(core.Zero) < 0 {
		utils.Error("Invalid input: stake must be positive\n")
	}

	source := types.TxInput{
		Address: sourceAddress,
		Coins: types.Coins{
			SCPTWei: stake,
			SPAYWei: new(big.Int).SetUint64(0),
		},
		Sequence: uint64(seqFlag),
	}

	depositStakeTx := &types.DepositStakeTxV2{
		Fee: types.Coins{
			SCPTWei: new(big.Int).SetUint64(0),
			SPAYWei: fee,
		},
		Source:  source,
		Purpose: purposeFlag,
	}

	// Parse holder flag.
	var holderAddress common.Address
	if purposeFlag == core.StakeForValidator {
		if len(holderFlag) != 40 && len(holderFlag) != 42 {
			utils.Error("holder must be a valid address")
		}
		holderAddress = common.HexToAddress(holderFlag)
	} else {
		if strings.HasPrefix(holderFlag, "0x") {
			holderFlag = holderFlag[2:]
		}
		if len(holderFlag) != 458 {
			utils.Error("Holder must be a valid guardian address")
		}
		guardianKeyBytes, err := hex.DecodeString(holderFlag)
		if err != nil {
			utils.Error("Failed to decode guardian address: %v\n", err)
		}
		holderAddress = common.BytesToAddress(guardianKeyBytes[:20])
		blsPubkey, err := bls.PublicKeyFromBytes(guardianKeyBytes[20:68])
		if err != nil {
			utils.Error("Failed to decode bls Pubkey: %v\n", err)
		}
		blsPop, err := bls.SignatureFromBytes(guardianKeyBytes[68:164])
		if err != nil {
			utils.Error("Failed to decode bls POP: %v\n", err)
		}
		holderSig, err := crypto.SignatureFromBytes(guardianKeyBytes[164:])
		if err != nil {
			utils.Error("Failed to decode signature: %v\n", err)
		}

		depositStakeTx.BlsPubkey = blsPubkey
		depositStakeTx.BlsPop = blsPop
		depositStakeTx.HolderSig = holderSig
	}

	depositStakeTx.Holder = types.TxOutput{
		Address: holderAddress,
	}

	sig, err := wallet.Sign(sourceAddress, depositStakeTx.SignBytes(chainIDFlag))
	if err != nil {
		utils.Error("Failed to sign transaction: %v\n", err)
	}
	depositStakeTx.SetSignature(sourceAddress, sig)

	raw, err := types.TxToBytes(depositStakeTx)
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
	depositStakeCmd.Flags().StringVar(&chainIDFlag, "chain", "", "Chain ID")
	depositStakeCmd.Flags().StringVar(&sourceFlag, "source", "", "Source of the stake")
	depositStakeCmd.Flags().StringVar(&holderFlag, "holder", "", "Holder of the stake")
	depositStakeCmd.Flags().StringVar(&pathFlag, "path", "", "Wallet derivation path")
	depositStakeCmd.Flags().StringVar(&feeFlag, "fee", fmt.Sprintf("%dwei", types.MinimumTransactionFeeSPAYWei), "Fee")
	depositStakeCmd.Flags().Uint64Var(&seqFlag, "seq", 0, "Sequence number of the transaction")
	depositStakeCmd.Flags().StringVar(&stakeInScriptFlag, "stake", "5000000", "Script amount to stake")
	depositStakeCmd.Flags().Uint8Var(&purposeFlag, "purpose", 0, "Purpose of staking")
	depositStakeCmd.Flags().StringVar(&walletFlag, "wallet", "soft", "Wallet type (soft|nano)")

	depositStakeCmd.MarkFlagRequired("chain")
	depositStakeCmd.MarkFlagRequired("source")
	depositStakeCmd.MarkFlagRequired("holder")
	depositStakeCmd.MarkFlagRequired("seq")
	depositStakeCmd.MarkFlagRequired("stake")
}
