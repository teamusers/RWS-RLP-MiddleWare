package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type TokenMeta struct {
	ID                 uint64          `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Chain              string          `json:"chain" gorm:"column:chain"`
	CA                 string          `json:"ca" gorm:"column:ca"`
	Name               string          `json:"name" gorm:"column:name"`
	Symbol             string          `json:"symbol" gorm:"column:symbol"`
	Description        string          `json:"description" gorm:"column:description"`
	Logo               string          `json:"logo" gorm:"column:logo"`
	Website            string          `json:"website" gorm:"column:website"`
	XURI               string          `json:"x_uri" gorm:"column:x_uri"`
	TgURI              string          `json:"tg_uri" gorm:"column:tg_uri"`
	TotalSupply        decimal.Decimal `json:"total_supply" gorm:"column:total_supply"`
	TMintable          string          `json:"t_mintable" gorm:"column:t_mintable"`             // CHAR(1), 可能是 'Y' 或 'N'
	TTransferable      string          `json:"t_transferable" gorm:"column:t_transferable"`     // CHAR(1)
	TMetaChangable     string          `json:"t_meta_changable" gorm:"column:t_meta_changable"` // CHAR(1)
	TFreezable         string          `json:"t_freezable" gorm:"column:t_freezable"`           // CHAR(1)
	TClosable          string          `json:"t_closable" gorm:"column:t_closable"`             // CHAR(1)
	TUpgradable        string          `json:"t_upgradable" gorm:"column:t_upgradable"`         // CHAR(1)
	OpenTime           string          `json:"open_time" gorm:"column:open_time"`
	CTransferFee       decimal.Decimal `json:"c_transfer_fee" gorm:"column:c_transfer_fee"`
	MetaUri            string          `json:"meta_uri" gorm:"column:meta_uri"`
	MetaContent        string          `json:"meta_content" gorm:"column:meta_content"`
	CreateTime         time.Time       `json:"create_time" gorm:"column:create_time"`
	Price              decimal.Decimal `json:"price" gorm:"-"`
	Change5m           decimal.Decimal `json:"change_5m" gorm:"-"`
	Change1h           decimal.Decimal `json:"change_1h" gorm:"-"`
	Change4h           decimal.Decimal `json:"change_4h" gorm:"-"`
	Change24h          decimal.Decimal `json:"change_24h" gorm:"-"`
	Circulation        decimal.Decimal `json:"circulation" gorm:"-"`
	Liquidity          decimal.Decimal `json:"liquidity" gorm:"-"`
	CirculatingSupply  decimal.Decimal `json:"circulating_supply" gorm:"-"`
	HolderCount        uint64          `json:"holder_count" gorm:"-"`
	Top10HolderBalance decimal.Decimal `json:"top_10_holder_balance" gorm:"-"`
	Top10HolderPercent decimal.Decimal `json:"top_10_holder_percent" gorm:"-"`
	Top10UserBalance   decimal.Decimal `json:"top_10_user_balance" gorm:"-"`
	Top10UserPercent   decimal.Decimal `json:"top_10_user_percent" gorm:"-"`
	DevBalance         decimal.Decimal `json:"dev_balance" gorm:"-"`
	DevPercentage      decimal.Decimal `json:"dev_Percentage" gorm:"-"`
	Create             string          `json:"create" gorm:"-"`
	TagList            []string        `json:"tag_list" gorm:"-"`
}
type TokenPriceChg struct {
	StartTime  string          `json:"start_time" gorm:"column:start_time"`
	EndTime    string          `json:"end_time" gorm:"column:end_time"`
	Period     string          `json:"period" gorm:"column:period"`
	StartPrice decimal.Decimal `json:"startPrice" gorm:"column:start_price"`
	EndPrice   decimal.Decimal `json:"endPrice" gorm:"column:end_price"`
	Change     decimal.Decimal `json:"chg" gorm:"-"`
}
type TokenHolders struct {
	TotalOfTop string        `json:"total_of_top" gorm:"-"`
	RatioOfTop string        `json:"ratio_of_top" gorm:"-"`
	Holders    []TokenHolder `json:"token_holder" gorm:"-"`
}
type TokenHolder struct {
	Wallet string   `json:"wallet" gorm:"-"`
	CA     string   `json:"ca" gorm:"-"`
	Amount string   `json:"amount" gorm:"-"`
	Tags   []string `json:"tags" gorm:"-"`
}

func (TokenMeta) TableName() string {
	return "token_meta_info"
}

type TokenPair struct {
	ID         uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	TokenID    uint64    `json:"token_id" gorm:"column:token_id"`
	Pair       string    `json:"pair" gorm:"column:pair"`
	Dex        string    `json:"dex" gorm:"column:dex"`
	OpenTime   string    `json:"open_time" gorm:"column:open_time"`
	InitTvl    string    `json:"init_tvl" gorm:"column:init_tvl"`
	LpAmount   string    `json:"lp_amount" gorm:"column:lp_amount"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
}

func (TokenPair) TableName() string {
	return "token_pool_pair"
}

type KlineData struct {
	Time   string          `json:"time" gorm:"column:time"`
	Pair   string          `json:"pair" gorm:"column:pair"`
	Open   decimal.Decimal `json:"open" gorm:"column:open"`
	High   decimal.Decimal `json:"high" gorm:"column:high"`
	Low    decimal.Decimal `json:"low" gorm:"column:low"`
	Close  decimal.Decimal `json:"close" gorm:"column:close"`
	Volume decimal.Decimal `json:"volume" gorm:"column:volume"`
}

func (KlineData) TableName() string {
	return "sol_tv_kline_data_1m"
}

type TokenChgByPeriod struct {
	Period             string          `json:"period" gorm:"column:period"`
	StartPrice         decimal.Decimal `json:"start_price" gorm:"column:start_price"`
	EndPrice           decimal.Decimal `json:"end_price" gorm:"column:end_price"`
	PriceChangePercent decimal.Decimal `json:"price_change_percent"`
	TradeCount         int             `json:"trade_count" gorm:"column:tradeCount"`
	TradeBuyCount      int             `json:"trade_buy_count" gorm:"column:tradeBuyCount"`
	TradeSellCount     int             `json:"trade_sell_count" gorm:"column:tradeSellCount"`
	TradeBuy           decimal.Decimal `json:"trade_buy" gorm:"column:tradeBuy"`
	TradeSell          decimal.Decimal `json:"trade_sell" gorm:"column:tradeSell"`
	TradeVol           decimal.Decimal `json:"trade_vol" gorm:"column:tradeVol"`
}
