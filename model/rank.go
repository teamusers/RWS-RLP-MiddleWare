package model

import (
	"github.com/shopspring/decimal"
)

type RankingItem struct {
	Icon          string          `json:"icon"`
	Address       string          `json:"address"`
	Symbol        string          `json:"symbol"`
	Price         decimal.Decimal `json:"price"`
	Change        decimal.Decimal `json:"change"`
	Circulation   decimal.Decimal `json:"circulation"`    // 流通量
	TradeVol      decimal.Decimal `json:"trade_vol"`      //成交额
	TradeBuyVol   decimal.Decimal `json:"trade_buy_vol"`  //买入额
	TradeSellVol  decimal.Decimal `json:"trade_sell_vol"` //卖出额
	TxCount       uint64          `json:"tx_count"`       //交易笔数
	TxBuyCount    uint64          `json:"tx_buy_count"`   //买入笔数
	TxSellCount   uint64          `json:"tx_sell_count"`  //卖出笔数
	MarketCap     decimal.Decimal `json:"market_cap"`     //市值
	CreateTime    string          `json:"create_time"`
	AttentionFlag int             `json:"attention_flag"` //如果用户登录，返回关注信息
	HolderCount   uint64          `json:"holder_count"`   //持币人数
	TxAddrCount   uint64          `json:"tx_addr_count"`  //独立交易地址数
}

func (h RankingItem) Validate() bool {
	return len(h.Address) > 0 && len(h.Symbol) > 0 && h.Price.Sign() > 0
}

type PriceChange struct {
	CA         string          `gorm:"column:ca" json:"ca"`
	FirstPrice decimal.Decimal `gorm:"column:firstPrice" json:"firstPrice"`
	LastPrice  decimal.Decimal `gorm:"column:lastPrice" json:"lastPrice"`
}
type LastPrice struct {
	CA    string          `gorm:"column:ca" json:"ca"`
	Price decimal.Decimal `gorm:"column:lastPrice" json:"lastPrice"`
}
type RankingItem1 struct {
	Icon        string `json:"icon"`
	Symbol      string `json:"symbol"`
	Price       string `json:"price"`
	Change      string `json:"change"`
	Circulation string `json:"circulation"`
	Volume      string `json:"volume"`
	TVL         string `json:"tvl"`
	MarketCap   string `json:"market_cap"`
	Holders     string `json:"holders"`
	CreateTime  string `json:"create_time"`
}

type RankingItemOriginal struct {
	Address        string          `gorm:"column:address" json:"address"`
	TradeCount     uint64          `gorm:"column:tradeCount" json:"tradeCount"`
	TradeBuyCount  uint64          `gorm:"column:tradeBuyCount" json:"tradeBuyCount"`
	TradeSellCount uint64          `gorm:"column:tradeSellCount" json:"tradeSellCount"`
	TradeBuy       decimal.Decimal `gorm:"column:tradeBuy" json:"tradeBuy"`
	TradeSell      decimal.Decimal `gorm:"column:tradeSell" json:"tradeSell"`
	TradeVol       decimal.Decimal `gorm:"column:tradeVol" json:"tradeVol"`
}
