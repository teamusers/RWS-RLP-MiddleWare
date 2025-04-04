package model

import (
	"github.com/shopspring/decimal"
)

type StatWalletToken struct {
	ID             uint64          `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Chain          string          `json:"chain" gorm:"column:chain"`
	Token          string          `json:"token" gorm:"column:token"`
	Wallet         string          `json:"wallet" gorm:"column:wallet"`
	Balance        uint64          `json:"balance" gorm:"column:balance"`
	CountSell      int             `json:"count_sell" gorm:"column:count_sell"`
	CountBuy       int             `json:"count_buy" gorm:"column:count_buy"`
	TotalSell      uint64          `json:"total_sell" gorm:"column:total_sell"`
	TotalBuy       uint64          `json:"total_buy" gorm:"column:total_buy"`
	FirstTx        string          `json:"first_tx" gorm:"column:first_tx"`
	FirstTs        string          `json:"first_ts" gorm:"column:first_ts"`
	LastTx         string          `json:"last_tx" gorm:"column:last_tx"`
	LastTs         string          `json:"last_ts" gorm:"column:last_ts"`
	LastDirect     string          `json:"last_direct" gorm:"column:last_direct"`
	TotalSellValue decimal.Decimal `json:"total_sell_value" gorm:"column:total_sell_value"`
	TotalBuyValue  decimal.Decimal `json:"total_buy_value" gorm:"column:total_buy_value"`
}

func (StatWalletToken) TableName() string {
	return "stat_wallet_token"
}
