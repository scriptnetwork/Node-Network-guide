package types

import (
	"math/big"

	"github.com/scripttoken/script/common"
)

const (
	// DenomSCPTWei is the basic unit of script, 1 Script = 10^18 SCPTWei
	DenomSCPTWei string = "SCPTWei"

	// DenomSPAYWei is the basic unit of script, 1 Script = 10^18 SCPTWei
	DenomSPAYWei string = "SPAYWei"

	// Initial gas parameters

	// MinimumGasPrice is the minimum gas price for a smart contract transaction
	MinimumGasPrice uint64 = 1e8

	// MaximumTxGasLimit is the maximum gas limit for a smart contract transaction
	//MaximumTxGasLimit uint64 = 2e6
	MaximumTxGasLimit uint64 = 10e6

	// MinimumTransactionFeeSPAYWei specifies the minimum fee for a regular transaction
	MinimumTransactionFeeSPAYWei uint64 = 1e12

	// June 2021 gas burn adjustment

	// MinimumGasPrice is the minimum gas price for a smart contract transaction
	MinimumGasPriceJune2021 uint64 = 1e8

	// MaximumTxGasLimit is the maximum gas limit for a smart contract transaction
	MaximumTxGasLimitJune2021 uint64 = 10e6

	// MinimumTransactionFeeSPAYWei specifies the minimum fee for a regular transaction
	MinimumTransactionFeeSPAYWeiJune2021 uint64 = 1e12

	// MaxAccountsAffectedPerTx specifies the max number of accounts one transaction is allowed to modify to avoid spamming
	MaxAccountsAffectedPerTx = 512
)

const (
	// ValidatorScriptGenerationRateNumerator is used for calculating the generation rate of Script for validators
	//ValidatorScriptGenerationRateNumerator int64 = 317
	ValidatorScriptGenerationRateNumerator int64 = 0 // ZERO inflation for Script

	// ValidatorScriptGenerationRateDenominator is used for calculating the generation rate of Script for validators
	// ValidatorScriptGenerationRateNumerator / ValidatorScriptGenerationRateDenominator is the amount of SCPTWei
	// generated per existing SCPTWei per new block
	ValidatorScriptGenerationRateDenominator int64 = 1e11

	// ValidatorSPAYGenerationRateNumerator is used for calculating the generation rate of SPAY for validators
	ValidatorSPAYGenerationRateNumerator int64 = 0 // ZERO initial inflation for SPAY

	// ValidatorSPAYGenerationRateDenominator is used for calculating the generation rate of SPAY for validators
	// ValidatorSPAYGenerationRateNumerator / ValidatorSPAYGenerationRateDenominator is the amount of SPAYWei
	// generated per existing SCPTWei per new block
	ValidatorSPAYGenerationRateDenominator int64 = 1e9

	// RegularSPAYGenerationRateNumerator is used for calculating the generation rate of SPAY for other types of accounts
	//RegularSPAYGenerationRateNumerator int64 = 1900
	RegularSPAYGenerationRateNumerator int64 = 0 // ZERO initial inflation for SPAY

	// RegularSPAYGenerationRateDenominator is used for calculating the generation rate of SPAY for other types of accounts
	// RegularSPAYGenerationRateNumerator / RegularSPAYGenerationRateDenominator is the amount of SPAYWei
	// generated per existing SCPTWei per new block
	RegularSPAYGenerationRateDenominator int64 = 1e10
)

const (

	// ServiceRewardVerificationBlockDelay gives the block delay for service certificate verification
	ServiceRewardVerificationBlockDelay uint64 = 2

	// ServiceRewardFulfillmentBlockDelay gives the block delay for service reward fulfillment
	ServiceRewardFulfillmentBlockDelay uint64 = 4
)

const (

	// MaximumTargetAddressesForStakeBinding gives the maximum number of target addresses that can be associated with a bound stake
	MaximumTargetAddressesForStakeBinding uint = 1024

	// MaximumFundReserveDuration indicates the maximum duration (in terms of number of blocks) of reserving fund
	MaximumFundReserveDuration uint64 = 12 * 3600

	// MinimumFundReserveDuration indicates the minimum duration (in terms of number of blocks) of reserving fund
	MinimumFundReserveDuration uint64 = 300

	// ReservedFundFreezePeriodDuration indicates the freeze duration (in terms of number of blocks) of the reserved fund
	ReservedFundFreezePeriodDuration uint64 = 5
)

func GetMinimumGasPrice(blockHeight uint64) *big.Int {
	if blockHeight < common.HeightJune2021FeeAdjustment {
		return new(big.Int).SetUint64(MinimumGasPrice)
	}

	return new(big.Int).SetUint64(MinimumGasPriceJune2021)
}

func GetMaxGasLimit(blockHeight uint64) *big.Int {
	if blockHeight < common.HeightJune2021FeeAdjustment {
		return new(big.Int).SetUint64(MaximumTxGasLimit)
	}

	return new(big.Int).SetUint64(MaximumTxGasLimitJune2021)
}

func GetMinimumTransactionFeeSPAYWei(blockHeight uint64) *big.Int {
	if blockHeight < common.HeightJune2021FeeAdjustment {
		return new(big.Int).SetUint64(MinimumTransactionFeeSPAYWei)
	}

	return new(big.Int).SetUint64(MinimumTransactionFeeSPAYWeiJune2021)
}

// Special handling for many-to-many SendTx
func GetSendTxMinimumTransactionFeeSPAYWei(numAccountsAffected uint64, blockHeight uint64) *big.Int {
	if blockHeight < common.HeightJune2021FeeAdjustment {
		return new(big.Int).SetUint64(MinimumTransactionFeeSPAYWei) // backward compatiblity
	}

	if numAccountsAffected < 2 {
		numAccountsAffected = 2
	}

	// minSendTxFee = numAccountsAffected * MinimumTransactionFeeSPAYWeiJune2021 / 2
	minSendTxFee := big.NewInt(1).Mul(new(big.Int).SetUint64(numAccountsAffected), new(big.Int).SetUint64(MinimumTransactionFeeSPAYWeiJune2021))
	minSendTxFee = big.NewInt(1).Div(minSendTxFee, new(big.Int).SetUint64(2))

	return minSendTxFee
}
