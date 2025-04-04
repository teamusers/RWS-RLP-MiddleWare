package topic

import (
	"time"

	"github.com/shopspring/decimal"
)

type SolTokenFlow struct {
	Tx             string          `json:"tx" gorm:"column:tx"`
	LogIndex       string          `json:"log_index" gorm:"column:log_index"`
	TxTime         string          `json:"tx_time" gorm:"column:tx_time"`
	Pair           string          `json:"pair" gorm:"column:pair"`
	FromAddr       string          `json:"from_addr" gorm:"column:from_addr"`
	ToAddr         string          `json:"to_addr" gorm:"column:to_addr"`
	Type           string          `json:"type" gorm:"column:type"`
	Block          int64           `json:"block" gorm:"column:block"`
	Token0         string          `json:"token0" gorm:"column:token0"`
	Token1         string          `json:"token1" gorm:"column:token1"`
	Amount0In      uint64          `json:"amount0_in" gorm:"column:amount0_in"`
	Amount0Out     uint64          `json:"amount0_out" gorm:"column:amount0_out"`
	Amount1In      uint64          `json:"amount1_in" gorm:"column:amount1_in"`
	Amount1Out     uint64          `json:"amount1_out" gorm:"column:amount1_out"`
	GasPrice       int64           `json:"gas_price" gorm:"column:gas_price"`
	GasUse         int64           `json:"gas_use" gorm:"column:gas_use"`
	TxFee          uint64          `json:"tx_fee" gorm:"column:tx_fee"`
	Price          decimal.Decimal `json:"price" gorm:"column:price"`
	Payer          string          `json:"payer" gorm:"column:payer"`
	Token0Decimals uint8           `gorm:"column:token0_decimals" json:"token0_decimals"`
	TradeDirection string          `json:"trade_direction" gorm:"column:trade_direction"`
}

type SolTokenHold struct {
	Address      string             `json:"address" gorm:"column:address"`
	Amount       uint64             `json:"amount" gorm:"column:amount"`
	Decimals     uint64             `json:"decimals" gorm:"column:decimals"`
	UpdateTime   string             `json:"update_time" gorm:"column:update_time"`
	Source       TxWithTime         `json:"source" gorm:"-"`
	Trans        map[string]TxCount `json:"trans" gorm:"-"`
	LastTranTime time.Time          `json:"last_tran_time" gorm:"-"`
}

type TxWithTime struct {
	TxHash string    `json:"tx_hash"`
	Time   time.Time `json:"tx_time"`
}

type TxCount struct {
	Value decimal.Decimal `json:"value"`
	Count int             `json:"count"`
}

type SolTokenInfo struct {
	ID          uint64          `json:"id" gorm:"column:id"`
	Name        string          `json:"name" gorm:"column:name"`
	Symbol      string          `json:"symbol" gorm:"column:symbol"`
	Address     string          `json:"address" gorm:"column:address"`
	Decimals    int64           `json:"decimals" gorm:"column:decimals"`
	TotalSupply decimal.Decimal `json:"total_supply" gorm:"column:total_supply"`
	MetaUri     string          `json:"meta_uri" gorm:"column:meta_uri"`
}
