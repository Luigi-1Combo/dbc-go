package common

import (
	"github.com/gagliardetto/solana-go"
	"lukechampine.com/uint128"
)

type BaseFeeConfig struct {
	CliffFeeNumerator uint64
	PeriodFrequency   uint64
	ReductionFactor   uint64
	NumberOfPeriod    uint16
	FeeSchedulerMode  uint8
	Padding0          [5]uint8
}

type DynamicFeeConfig struct {
	Initialized              uint8
	Padding                  [7]uint8
	MaxVolatilityAccumulator uint32
	VariableFeeControl       uint32
	BinStep                  uint16
	FilterPeriod             uint16
	DecayPeriod              uint16
	ReductionFactor          uint16
	Padding2                 [8]uint8
	BinStepU128              uint128.Uint128
}

type PoolFeesConfig struct {
	BaseFee            BaseFeeConfig
	DynamicFee         DynamicFeeConfig
	Padding0           [5]uint64
	Padding1           [6]uint8
	ProtocolFeePercent uint8
	ReferralFeePercent uint8
}

// LiquidityDistributionConfig represents the liquidity distribution configuration
type LiquidityDistributionConfig struct {
	SqrtPrice uint128.Uint128
	Liquidity uint128.Uint128
}

type LockedVestingConfig struct {
	AmountPerPeriod                uint64
	CliffDurationFromMigrationTime uint64
	Frequency                      uint64
	NumberOfPeriod                 uint64
	CliffUnlockAmount              uint64
	Padding                        uint64
}

// PoolConfig represents the pool configuration data structure
type PoolConfig struct {
	QuoteMint                   solana.PublicKey
	FeeClaimer                  solana.PublicKey
	LeftoverReceiver            solana.PublicKey
	PoolFees                    PoolFeesConfig
	CollectFeeMode              uint8
	MigrationOption             uint8
	ActivationType              uint8
	TokenDecimal                uint8
	Version                     uint8
	TokenType                   uint8
	QuoteTokenFlag              uint8
	PartnerLockedLpPercentage   uint8
	PartnerLpPercentage         uint8
	CreatorLockedLpPercentage   uint8
	CreatorLpPercentage         uint8
	MigrationFeeOption          uint8
	FixedTokenSupplyFlag        uint8
	CreatorTradingFeePercentage uint8
	Padding0                    [2]uint8
	Padding1                    [8]uint8
	SwapBaseAmount              uint64
	MigrationQuoteThreshold     uint64
	MigrationBaseThreshold      uint64
	MigrationSqrtPrice          uint128.Uint128
	LockedVestingConfig         LockedVestingConfig
	PreMigrationTokenSupply     uint64
	PostMigrationTokenSupply    uint64
	Padding2                    [2]uint128.Uint128
	SqrtStartPrice              uint128.Uint128
	Curve                       [20]LiquidityDistributionConfig
}
