package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		LeverageMax:         sdkmath.LegacyNewDec(10),
		EpochLength:         (int64)(1),
		MaxOpenPositions:    (int64)(9999),
		PoolOpenThreshold:   sdkmath.LegacyNewDecWithPrec(2, 1),  // 0.2
		SafetyFactor:        sdkmath.LegacyNewDecWithPrec(11, 1), // 1.1
		WhitelistingEnabled: false,
		FallbackEnabled:     true,
		NumberPerBlock:      (int64)(1000),
		EnabledPools:        []uint64(nil),
		ExitBuffer:          sdkmath.LegacyMustNewDecFromStr("0.05"),
		StopLossEnabled:     true,
		LiabilitiesFactor:   sdkmath.LegacyMustNewDecFromStr("1.0"),
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.LeverageMax.IsNil() {
		return fmt.Errorf("leverage max must be not nil")
	}
	if !p.LeverageMax.GT(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("leverage max must be greater than 1: %s", p.LeverageMax.String())
	}
	if p.EpochLength <= 0 {
		return fmt.Errorf("epoch length should be positive: %d", p.EpochLength)
	}
	if p.PoolOpenThreshold.IsNil() {
		return fmt.Errorf("pool open threshold must be not nil")
	}
	if !p.PoolOpenThreshold.IsPositive() {
		return fmt.Errorf("pool open threshold must be positive: %s", p.PoolOpenThreshold.String())
	}
	if p.SafetyFactor.IsNil() {
		return fmt.Errorf("safety factor must be not nil")
	}
	if !p.SafetyFactor.IsPositive() {
		return fmt.Errorf("safety factor must be positive: %s", p.SafetyFactor.String())
	}
	if p.NumberPerBlock < 0 {
		return fmt.Errorf("number of positions per block must be positive: %d", p.NumberPerBlock)
	}

	if p.NumberPerBlock > MaxPageLimit {
		return fmt.Errorf("number of positions per block should not exceed page limit: %d, number of positions: %d", MaxPageLimit, p.NumberPerBlock)
	}

	if containsDuplicates(p.EnabledPools) {
		return fmt.Errorf("array must not contain duplicate values")
	}

	if p.ExitBuffer.IsNil() {
		return fmt.Errorf("exit buffer must be not nil")
	}

	if p.LiabilitiesFactor.IsNil() {
		return fmt.Errorf("liabilities factor must be not nil")
	}

	if !p.LiabilitiesFactor.IsPositive() {
		return fmt.Errorf("liabilities factor must be positive: %s", p.LiabilitiesFactor.String())
	}

	if p.LiabilitiesFactor.GT(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("liabilities factor must be less than or equal to 1: %s", p.LiabilitiesFactor.String())
	}

	return nil
}

func containsDuplicates(arr []uint64) bool {
	valueMap := make(map[uint64]struct{})
	for _, num := range arr {
		if _, exists := valueMap[num]; exists {
			return true
		}
		valueMap[num] = struct{}{}
	}
	return false
}

func (p Params) GetBigDecSafetyFactor() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.SafetyFactor)
}

func (p Params) GetBigDecPoolOpenThreshold() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.PoolOpenThreshold)
}

func (p Params) GetBigDecExitBuffer() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.ExitBuffer)
}
