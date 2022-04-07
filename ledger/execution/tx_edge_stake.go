package execution

import (
	"fmt"
	"math/big"

	"github.com/scripttoken/script/common"
	"github.com/scripttoken/script/common/result"
	"github.com/scripttoken/script/core"
	st "github.com/scripttoken/script/ledger/state"
	"github.com/scripttoken/script/ledger/types"
)

var _ TxExecutor = (*EdgeStakeTxExecutor)(nil)

// ------------------------------- Send Transaction -----------------------------------

// EdgeStakeTxExecutor implements the TxExecutor interface
type EdgeStakeTxExecutor struct {
}

// NewEdgeStakeTxExecutor creates a new instance of EdgeStakeTxExecutor
func NewEdgeStakeTxExecutor() *EdgeStakeTxExecutor {
	return &EdgeStakeTxExecutor{}
}

func (exec *EdgeStakeTxExecutor) sanityCheck(chainID string, view *st.StoreView, transaction types.Tx) result.Result {
	tx := transaction.(*types.EdgeStakeTx)

	// Validate inputs and outputs, basic
	res := validateInputsBasic(tx.Inputs)
	if res.IsError() {
		return res
	}
	res = validateOutputsBasic(tx.Outputs)
	if res.IsError() {
		return res
	}

	if len(tx.Inputs) == 0 || len(tx.Outputs) == 0 {
		return result.Error("Invalid edgeStakeTx, Inputs and/or Outputs are empty")
	}

	numAccountsAffected := uint64(len(tx.Inputs) + len(tx.Outputs))
	if numAccountsAffected > types.MaxAccountsAffectedPerTx {
		return result.Error("Trasaction modifying too many accounts. At most %v accounts are allowed per transaction",
			types.MaxAccountsAffectedPerTx)
	}

	// Get inputs
	accounts, res := getInputs(view, tx.Inputs)
	if res.IsError() {
		return res
	}

	// Get or make outputs.
	accounts, res = getOrMakeOutputs(view, accounts, tx.Outputs)
	if res.IsError() {
		return res
	}

	blockHeight := view.Height() + 1
	if blockHeight >= common.HeightEnableSmartContract {
		for _, outAcc := range accounts {
			if outAcc.IsASmartContract() {
				return result.Error(
					fmt.Sprintf("Sending SCPT/SPAY to a smart contract (%v) through a EdgeStakeTx transaction is not allowed", outAcc.Address))
			}
		}
	}

	// Validate inputs and outputs, advanced
	signBytes := tx.SignBytes(chainID)
	inTotal, res := validateInputsAdvanced(accounts, signBytes, tx.Inputs)
	if res.IsError() {
		return res
	}

	if !sanityCheckForFee(tx.Fee) {
		return result.Error("Insufficient fee. Transaction fee needs to be at least %v SPAYWei",
			types.MinimumTransactionFeeSPAYWei).WithErrorCode(result.CodeInvalidFee)
	}

	outTotal := sumOutputs(tx.Outputs)
	outPlusFees := outTotal
	outPlusFees = outTotal.Plus(tx.Fee)
	if !inTotal.IsEqual(outPlusFees) {
		return result.Error("Input total (%v) != output total + fees (%v)", inTotal, outPlusFees)
	}

	return result.OK
}

func (exec *EdgeStakeTxExecutor) process(chainID string, view *st.StoreView, transaction types.Tx) (common.Hash, result.Result) {
	tx := transaction.(*types.EdgeStakeTx)

	accounts, res := getInputs(view, tx.Inputs)
	if res.IsError() {
		return common.Hash{}, res
	}

	accounts, res = getOrMakeOutputs(view, accounts, tx.Outputs)
	if res.IsError() {
		return common.Hash{}, res
	}

	adjustByInputs(view, accounts, tx.Inputs)
	adjustByOutputs(view, accounts, tx.Outputs)

	txHash := types.TxID(chainID, tx)
	return txHash, result.OK
}

func (exec *EdgeStakeTxExecutor) getTxInfo(transaction types.Tx) *core.TxInfo {
	tx := transaction.(*types.EdgeStakeTx)
	return &core.TxInfo{
		Address:           tx.Inputs[0].Address,
		Sequence:          tx.Inputs[0].Sequence,
		EffectiveGasPrice: exec.calculateEffectiveGasPrice(transaction),
	}
}

func (exec *EdgeStakeTxExecutor) calculateEffectiveGasPrice(transaction types.Tx) *big.Int {
	tx := transaction.(*types.EdgeStakeTx)
	fee := tx.Fee
	numAccountsAffected := uint64(len(tx.Inputs) + len(tx.Outputs))
	gasUint64 := types.GasSendTxPerAccount * numAccountsAffected
	if gasUint64 < 2*types.GasSendTxPerAccount {
		gasUint64 = 2 * types.GasSendTxPerAccount // to prevent spamming with invalid transactions, e.g. empty inputs/outputs
	}
	gas := new(big.Int).SetUint64(gasUint64)
	effectiveGasPrice := new(big.Int).Div(fee.SPAYWei, gas)
	return effectiveGasPrice
}
