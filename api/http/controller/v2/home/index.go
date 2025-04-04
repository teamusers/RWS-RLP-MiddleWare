package home

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stonksdex/externalapi/log"
	"github.com/stonksdex/externalapi/model"
	"github.com/stonksdex/externalapi/system"
	"github.com/stonksdex/externalapi/thirdpart"
	"github.com/stonksdex/externalapi/tools"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/stonksdex/externalapi/api/common"
	"github.com/stonksdex/externalapi/codes"
)

type RankingItem struct {
	Icon    string          `json:"icon"`
	Address string          `json:"address"`
	Symbol  string          `json:"symbol"`
	Price   decimal.Decimal `json:"price"`
	// Change      decimal.Decimal `json:"change"`
	Circulation decimal.Decimal `json:"circulation"` // 流通量
	// TradeVol      decimal.Decimal `json:"trade_vol"`      //成交额
	// TradeBuyVol   decimal.Decimal `json:"trade_buy_vol"`  //买入额
	// TradeSellVol  decimal.Decimal `json:"trade_sell_vol"` //卖出额
	// TxCount       uint64          `json:"tx_count"`      //交易笔数
	// TxBuyCount    uint64          `json:"tx_buy_count"`  //买入笔数
	// TxSellCount   uint64          `json:"tx_sell_count"` //卖出笔数
	MarketCap decimal.Decimal `json:"market_cap"` //市值
	// CreateTime    string          `json:"create_time"`
	AttentionFlag int    `json:"attention_flag"` //如果用户登录，返回关注信息
	HolderCount   uint64 `json:"holder_count"`   //持币人数
	// TxAddrCount   uint64          `json:"tx_addr_count"`  //独立交易地址数
}

type HomeRequest struct {
	Chain    string `json:"chain"`
	Interval string `json:"interval"`
	Sort     string `json:"sort"`
	Page     int    `json:"pn" binding:"required,min=1"`
	PageSize int    `json:"ps" binding:"required,min=1"`
	Type     int    `json:"type"`
}

func UpdateLeaderboard(c *gin.Context) {
	var req HomeRequest
	res := common.Response{}
	res.Timestamp = time.Now().Unix()

	if err := c.ShouldBindJSON(&req); err != nil {
		res.Code = codes.CODE_ERR_REQFORMAT
		res.Msg = "invalid request"
		c.JSON(http.StatusOK, res)
		return
	}
	if req.PageSize <= 0 {
		req.PageSize = 50
	}

	currentUser, exist := c.Get("user_id")

	log.Info(currentUser, exist)
	bc, _ := thirdpart.GetBirdClient()

	tq := thirdpart.TokenListV3Request{
		SortType:     "desc",
		MinLiquidity: tools.Float64Ptr(100),
	}

	if req.Sort == "volume" {
		if req.Interval == "1h" {
			tq.SortBy = thirdpart.Volume1HUsd
		} else if req.Interval == "2h" {
			tq.SortBy = thirdpart.Volume2HUsd
		} else if req.Interval == "4h" {
			tq.SortBy = thirdpart.Volume4HUsd
		} else if req.Interval == "8h" {
			tq.SortBy = thirdpart.Volume8HUsd
		} else if req.Interval == "24h" {
			tq.SortBy = thirdpart.Volume24HUsd
		}
	} else if req.Sort == "volume_change" {
		if req.Interval == "1h" {
			tq.SortBy = thirdpart.Volume1HChangePct
		} else if req.Interval == "2h" {
			tq.SortBy = thirdpart.Volume2HChangePct
		} else if req.Interval == "4h" {
			tq.SortBy = thirdpart.Volume4HChangePct
		} else if req.Interval == "8h" {
			tq.SortBy = thirdpart.Volume8HChangePct
		} else if req.Interval == "24h" {
			tq.SortBy = thirdpart.Volume24HChangePct
		}
	} else if req.Sort == "price_change" {
		if req.Interval == "1h" {
			tq.SortBy = thirdpart.PriceChange1HPct
		} else if req.Interval == "2h" {
			tq.SortBy = thirdpart.PriceChange2HPct
		} else if req.Interval == "4h" {
			tq.SortBy = thirdpart.PriceChange4HPct
		} else if req.Interval == "8h" {
			tq.SortBy = thirdpart.PriceChange8HPct
		} else if req.Interval == "24h" {
			tq.SortBy = thirdpart.PriceChange24HPct
		}
	} else if req.Sort == "trade_count" {
		if req.Interval == "1h" {
			tq.SortBy = thirdpart.Trade1HCount
		} else if req.Interval == "2h" {
			tq.SortBy = thirdpart.Trade2HCount
		} else if req.Interval == "4h" {
			tq.SortBy = thirdpart.Trade4HCount
		} else if req.Interval == "8h" {
			tq.SortBy = thirdpart.Trade8HCount
		} else if req.Interval == "24h" {
			tq.SortBy = thirdpart.Trade24HCount
		}
	}

	if len(req.Sort) == 0 {
		tq.SortBy = thirdpart.Volume1HUsd
	}

	result, err := bc.TokenTrending(req.Chain, tq, req.Page-1, req.PageSize)
	if err != nil {
		log.Error("[Birdeyd TokenTrending] failed: ", err)
		res.Code = codes.CODE_ERR_UNKNOWN
		res.Msg = "fetch data fail"
		c.JSON(http.StatusOK, res)
		return
	}

	type includeRankingItem struct {
		RankingItem
		Decimal                uint8                  `json:"decimal"`
		Volume1hUsd            decimal.Decimal        `json:"volume_1h_usd"`
		Volume1hChangePercent  decimal.Decimal        `json:"volume_1h_change_percent"`
		Volume2hUsd            decimal.Decimal        `json:"volume_2h_usd"`
		Volume2hChangePercent  decimal.Decimal        `json:"volume_2h_change_percent"`
		Volume4hUsd            decimal.Decimal        `json:"volume_4h_usd"`
		Volume4hChangePercent  decimal.Decimal        `json:"volume_4h_change_percent"`
		Volume8hUsd            decimal.Decimal        `json:"volume_8h_usd"`
		Volume8hChangePercent  decimal.Decimal        `json:"volume_8h_change_percent"`
		Volume24hUsd           decimal.Decimal        `json:"volume_24h_usd"`
		Volume24hChangePercent decimal.Decimal        `json:"volume_24h_change_percent"`
		Trade1hCount           int64                  `json:"trade_1h_count"`
		Trade2hCount           int64                  `json:"trade_2h_count"`
		Trade4hCount           int64                  `json:"trade_4h_count"`
		Trade8hCount           int64                  `json:"trade_8h_count"`
		Trade24hCount          int64                  `json:"trade_24h_count"`
		PriceChange1hPercent   decimal.Decimal        `json:"price_change_1h_percent"`
		PriceChange2hPercent   decimal.Decimal        `json:"price_change_2h_percent"`
		PriceChange4hPercent   decimal.Decimal        `json:"price_change_4h_percent"`
		PriceChange8hPercent   decimal.Decimal        `json:"price_change_8h_percent"`
		PriceChange24hPercent  decimal.Decimal        `json:"price_change_24h_percent"`
		LastTradeUnixTime      uint64                 `json:"last_trade_unix_time"`
		Extension              map[string]interface{} `json:"extensions"`
	}

	var ranksRaw []includeRankingItem

	for _, v := range result {
		ranksRaw = append(ranksRaw, includeRankingItem{
			RankingItem: RankingItem{
				Icon:        v.LogoURI,
				Symbol:      v.Symbol,
				Address:     v.Address,
				Price:       decimal.NewFromFloat(v.Price),
				Circulation: decimal.NewFromFloat(v.Liquidity),
				MarketCap:   decimal.NewFromFloat(v.MarketCap),
				HolderCount: uint64(v.Holder),
			},
			Decimal:                v.Decimals,
			Volume1hUsd:            decimal.NewFromFloat(v.Volume1HUsd),
			Volume1hChangePercent:  decimal.NewFromFloat(v.Volume1HChangePercent),
			Volume2hUsd:            decimal.NewFromFloat(v.Volume2HUsd),
			Volume2hChangePercent:  decimal.NewFromFloat(v.Volume2HChangePercent),
			Volume4hUsd:            decimal.NewFromFloat(v.Volume4HUsd),
			Volume4hChangePercent:  decimal.NewFromFloat(v.Volume4HChangePercent),
			Volume8hUsd:            decimal.NewFromFloat(v.Volume8HUsd),
			Volume8hChangePercent:  decimal.NewFromFloat(v.Volume8HChangePercent),
			Volume24hUsd:           decimal.NewFromFloat(v.Volume24HUsd),
			Volume24hChangePercent: decimal.NewFromFloat(v.Volume24HChangePercent),
			Trade1hCount:           int64(v.Trade1HCount),
			Trade2hCount:           int64(v.Trade2HCount),
			Trade4hCount:           int64(v.Trade4HCount),
			Trade8hCount:           int64(v.Trade8HCount),
			Trade24hCount:          int64(v.Trade24HCount),
			PriceChange1hPercent:   decimal.NewFromFloat(v.PriceChange1HPercent),
			PriceChange2hPercent:   decimal.NewFromFloat(v.PriceChange2HPercent),
			PriceChange4hPercent:   decimal.NewFromFloat(v.PriceChange4HPercent),
			PriceChange8hPercent:   decimal.NewFromFloat(v.PriceChange8HPercent),
			PriceChange24hPercent:  decimal.NewFromFloat(v.PriceChange24HPercent),
			LastTradeUnixTime:      v.LastTradeUnixTime,
			Extension:              v.Extension,
		})
	}

	res.Data = ranksRaw
	c.JSON(http.StatusOK, res)

}

func getUserAttention(chain string, cas []string, currentUser any, ex bool) []model.UserAttention {
	var attObjs []model.UserAttention
	if !ex {
		return []model.UserAttention{}
	}
	currentUserStr, _ := currentUser.(string)
	userID, err := strconv.ParseInt(currentUserStr, 10, 64)
	if err == nil {
		db := system.GetDb()
		db.Model(&model.UserAttention{}).Where("uw_id = ? and flag = ? and chain = ? and ca in ?", userID, 0, chain, cas).Find(&attObjs)
	}
	return attObjs
}

func getTokenMeta(cas []string, chain string) ([]model.TokenMeta, error) {
	if len(cas) == 0 {
		return nil, errors.New("ca is empty")
	}
	db := system.GetDb()
	var tokenMetas []model.TokenMeta
	if len(cas) == 1 {
		var existMeta model.TokenMeta
		err := db.Model(&model.TokenMeta{}).Where("chain = ? and ca = ?", chain, cas[0]).First(&existMeta).Error
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) && existMeta.ID <= 0 {
			return nil, err
		}
		tokenMetas = append(tokenMetas, existMeta)
		return tokenMetas, nil

	} else {
		err := db.Model(&model.TokenMeta{}).Where("chain = ? and ca in ?", chain, cas).Scan(&tokenMetas).Error
		if err != nil {
			return nil, err
		}
		return tokenMetas, nil
	}
}
