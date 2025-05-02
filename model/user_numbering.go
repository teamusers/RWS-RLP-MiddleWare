package model

import (
	"gorm.io/gorm"
)

type RLPUserNumbering struct {
	ID            uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Year          int64  `gorm:"column:year" json:"year"`
	Month         int64  `gorm:"column:month" json:"month"`
	Day           int64  `gorm:"column:day" json:"day"`
	RLP_ID        string `gorm:"column:rlp_id;unique" json:"rlp_id"`
	RLP_NO        string `gorm:"column:rlp_no;unique" json:"rlp_no"`
	RLPIDEndingNO int    `gorm:"column:rlp_id_ending_no" json:"rlp_id_ending_no"`
}

func MigrateRLPUserNumbering(db *gorm.DB) error {
	return db.AutoMigrate(&RLPUserNumbering{})
}
