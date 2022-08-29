package common

// HeightEnableValidatorReward specifies the minimal block height to enable the validtor SPAY reward
const HeightEnableValidatorReward uint64 = 1 // approximate time: 2pm January 14th, 2020 PST

// HeightEnableScript2 specifies the minimal block height to enable the Script2.0 feature.
const HeightEnableScript2 uint64 = 1 // approximate time: 12pm May 27th, 2020 PDT

// HeightLowerGNStakeThresholdTo1000 specifies the minimal block height to lower the GN Stake Threshold to 1,000 SCRIPT
const HeightLowerGNStakeThresholdTo1000 uint64 = 1 // approximate time: 12pm Dec 10th, 2020 PST

// HeightEnableSmartContract specifies the minimal block height to eanble the Turing-complete smart contract support
const HeightEnableSmartContract uint64 = 1 // approximate time: 12pm Dec 10th, 2020 PST

// HeightSampleStakingReward specifies the block heigth to enable sampling of staking reward
const HeightSampleStakingReward uint64 = 1 // approximate time: 7pm Mar 10th, 2021 PST

// HeightJune2021FeeAdjustment specifies the block heigth to enable transaction fee burning adjustment
const HeightJune2021FeeAdjustment uint64 = 1 // approximate time: 12pm June 11, 2021 PT

// HeightEnableScript3 specifies the minimal block height to enable the Script3.0 feature.
const HeightEnableScript3 uint64 = 1 // approximate time: 12pm June 30, 2021 PT

// HeightRPCCompatibility specifies the block height to enable Ethereum compatible RPC support
const HeightRPCCompatibility uint64 = 1 // approximate time: 12pm July 30, 2021 PT

// HeightTxWrapperExtension specifies the block height to extend the Tx Wrapper
const HeightTxWrapperExtension uint64 = 1

// HeightSupportScriptTokenInSmartContract specifies the block height to support Script in smart contracts
const HeightSupportScriptTokenInSmartContract uint64 = 1 // approximate time: 5pm Dec 4, 2021 PT

// HeightValidatorStakeChangedTo200K specifies the block height to lower the validator stake to 200,000 Script
const HeightValidatorStakeChangedTo200K uint64 = 1 // approximate time: 12pm Mar 14, 2022 PT

// HeightSupportWrappedScript specifies the block height to support wrapped Script
const HeightSupportWrappedScript uint64 = 1

// CheckpointInterval defines the interval between checkpoints.
const CheckpointInterval = int64(100)

// IsCheckPointHeight returns if a block height is a checkpoint.
func IsCheckPointHeight(height uint64) bool {
	return height%uint64(CheckpointInterval) == 1
}

// LastCheckPointHeight returns the height of the last checkpoint
func LastCheckPointHeight(height uint64) uint64 {
	multiple := height / uint64(CheckpointInterval)
	lastCheckpointHeight := uint64(CheckpointInterval)*multiple + 1
	return lastCheckpointHeight
}
