package model

import (
	"time"
)

const (
	TB_TxsDiagram = "txs_diagram"
	TB_TxsLog     = "txs_log"
)

type TxsDiagram struct {
	ID         uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID     uint64    `json:"uw_id" gorm:"column:uw_id"`
	Chain      string    `json:"chain" gorm:"column:chain"`
	FromToken  string    `json:"from_token" gorm:"column:from_token"`
	ToToken    string    `json:"to_token" gorm:"column:to_token"`
	FromAmount uint64    `json:"from_amount" gorm:"column:from_amount"`
	Slippage   uint64    `json:"slippage" gorm:"column:slippage"`
	Status     string    `json:"status" gorm:"column:status"`
	TxHash     string    `json:"tx_hash" gorm:"column:tx_hash"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time"`
}

func (TxsDiagram) TableName() string {
	return TB_TxsDiagram
}

type TxsLog struct {
	ID         uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TxId       uint64    `json:"tx_id" gorm:"column:tx_id"`
	Step       uint64    `json:"step" gorm:"column:step"`
	ErrMsg     string    `json:"err_msg" gorm:"column:err_msg"`
	Result     string    `json:"result" gorm:"column:result"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
}

func (TxsLog) TableName() string {
	return TB_TxsLog
}
