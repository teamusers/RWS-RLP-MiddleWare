package model

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/mr-tron/base58"
)

type UserWallet struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Wallet     string    `gorm:"column:wallet;type:varchar(255);not null" json:"wallet"`
	Chain      string    `gorm:"column:chain" json:"chain"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	MainID     uint64    `gorm:"column:main_id" json:"main_id"`
	RefID      uint64    `gorm:"column:ref_id" json:"ref_id"`
}

func (UserWallet) TableName() string {
	return "user_wallet"
}

type AuthMessage struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	AuthKey    string    `gorm:"column:auth_key;type:varchar(255);not null" json:"auth_key"`
	AuthMsg    string    `gorm:"column:auth_msg" json:"auth_msg"`
	Nonce      string    `gorm:"type:varchar(255);not null" json:"nonce"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	ExpireTime time.Time `gorm:"column:expire_time" json:"expire_time"`
}

func (AuthMessage) TableName() string {
	return "user_auth_msg"
}

func (auth AuthMessage) ComputeAuthDigest(base64Sig string) bool {
	data := auth.Format()

	publicKey, err := base58.Decode(auth.AuthKey)
	if err != nil {
		log.Println(err)
	}

	signature, err := base64.StdEncoding.DecodeString(base64Sig)
	if err != nil {
		log.Println(err)
	}

	return ed25519.Verify(publicKey, []byte(data), signature)
}

func (auth AuthMessage) Format() string {
	data := fmt.Sprintf("Wallet:%s\n|Message:%s\n|Nonce:%s\n",
		auth.AuthKey,
		auth.AuthMsg,
		auth.Nonce,
	)
	return data
}

type UserRef struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UwID       int64     `gorm:"column:uw_id;type:varchar(255);not null" json:"uw_id"`
	RefCode    string    `gorm:"column:ref_code" json:"ref_code"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
}

func (UserRef) TableName() string {
	return "user_ref_key"
}

type DailyCheck struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UwID      int64     `gorm:"column:uw_id;type:varchar(255);not null" json:"uw_id"`
	CheckDate string    `gorm:"column:check_date" json:"check_date"`
	CheckTime time.Time `gorm:"column:check_time" json:"check_time"`
}

func (DailyCheck) TableName() string {
	return "daily_checkin"
}

type UserAttention struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UwID       int64     `gorm:"column:uw_id;type:varchar(255);not null" json:"uw_id"`
	Chain      string    `gorm:"column:chain" json:"chain"`
	Ca         string    `gorm:"column:ca" json:"ca"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
	Flag       int       `gorm:"column:flag" json:"flag"`
}

func (UserAttention) TableName() string {
	return "user_attention"
}
