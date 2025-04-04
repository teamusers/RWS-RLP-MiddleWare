package model

import "github.com/shopspring/decimal"

// Pool doris pool model
type Pool struct {
	Pair               string          `gorm:"column:pair;type:varchar;not null;primaryKey" json:"pair"` // 交易对，作为唯一键的一部分
	Token0             string          `gorm:"column:token0;type:varchar;not null;" json:"token0"`       // token0，作为唯一键的一部分
	Token1             string          `gorm:"column:token1;type:varchar" json:"token1"`
	Dex                string          `gorm:"column:dex;type:varchar" json:"dex"`
	Amount0            uint64          `gorm:"column:amount0;type:bigint;default:0" json:"amount0"`
	Amount1            uint64          `gorm:"column:amount1;type:bigint;default:0" json:"amount1"`
	CreateAt           int64           `gorm:"column:create_at;type:bigint" json:"createAt"`
	Fee                int32           `gorm:"column:fee;type:int;default:0" json:"fee"`
	Tvl                decimal.Decimal `gorm:"column:tvl;type:numeric;default:0" json:"tvl"`
	BaseToken          string          `gorm:"column:base_token;type:varchar" json:"baseToken"`
	QuoteToken         string          `gorm:"column:quote_token;type:varchar" json:"quoteToken"`
	Authority          string          `gorm:"column:authority;type:varchar" json:"authority"`
	Liq                uint64          `gorm:"column:liq;type:bigint;default:0" json:"liq"`
	StartPrice         decimal.Decimal `gorm:"column:start_price;type:numeric" json:"startPrice"`
	BaseTokenDecimals  uint8           `gorm:"-" json:"baseTokenDecimals"`
	QuoteTokenDecimals uint8           `gorm:"-" json:"quoteTokenDecimals"`
}

func (Pool) TableName() string {
	return "sol_pool"
}
