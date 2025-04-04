package svs

import (
	"errors"
	"time"

	"github.com/stonksdex/externalapi/log"
	"github.com/stonksdex/externalapi/model"
	"github.com/stonksdex/externalapi/system"
	"gorm.io/gorm"
)

type TokenPairRes struct {
	ID         int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TokenID    int64     `json:"token_id" gorm:"column:token_id"`
	Pair       string    `json:"pair" gorm:"column:pair"`
	Dex        string    `json:"dex" gorm:"column:dex"`
	OpenTime   string    `json:"open_time" gorm:"column:open_time"`
	InitTvl    string    `json:"init_tvl" gorm:"column:init_tvl"`
	LpAmount   string    `json:"lp_amount" gorm:"column:lp_amount"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
}

func PoolsInMysql(ca string, chain string) ([]TokenPairRes, error) {
	db := system.GetDb()
	var indb []TokenPairRes
	// todo 这里的sql需要进行转义
	result := db.Table(model.TokenMeta{}.TableName()).
		Joins("inner join token_pool_pair tpp on token_meta_info.id = tpp.token_id").
		Select("tpp.*,token_meta_info.ca ca").
		Where("token_meta_info.ca = ?", ca).
		Order("init_tvl desc").Find(&indb)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Infof("Error:getting ca:%s, %v", ca, result.Error)
		}
		return []TokenPairRes{}, result.Error
	}
	if len(indb) > 0 {

	}
	return indb, nil
}
