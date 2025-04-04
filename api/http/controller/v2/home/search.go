package home

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stonksdex/externalapi/api/common"
	"github.com/stonksdex/externalapi/codes"
	"github.com/stonksdex/externalapi/log"
	"github.com/stonksdex/externalapi/thirdpart"
)

func Search(c *gin.Context) {
	res := common.Response{}
	res.Timestamp = time.Now().Unix()
	res.Code = codes.CODE_SUCCESS
	res.Msg = "success"

	searchKey := c.Param("key")
	sizeStr := c.DefaultQuery("size", "50")
	searchType := c.Query("type")
	if len(searchType) == 0 || (searchType != "symbol" && searchType != "name" && searchType != "address") {
		searchType = "symbol"
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 || size > 20 {
		size = 20
	}

	if len(searchKey) == 0 {
		c.JSON(http.StatusOK, res)
		return
	}

	b, _ := thirdpart.GetBirdClient()

	ts := thirdpart.TokenSearchRequest{
		Chain:       "solana",
		Keyword:     searchKey,
		Target:      "token",
		SearchMode:  "fuzzy",
		SearchBy:    searchType,
		SortBy:      "marketcap",
		SortType:    "desc",
		VerifyToken: "true",
	}
	tokenResult, marketResult, err := b.TokenSearch(ts, 0, size)
	if err != nil {
		log.Error("TokenSearch Error: ", searchKey, size, err)
	}

	type searchResult struct {
		Chain                        string          `json:"chain"`
		Type                         string          `json:"type"`
		CA                           string          `json:"ca"`
		Name                         string          `json:"name"`
		Symbol                       string          `json:"symbol"`
		Logo                         string          `json:"logo"`
		Liquidity                    decimal.Decimal `json:"liquidity"`
		TotalSupply                  decimal.Decimal `json:"total_supply"`
		FDV                          decimal.Decimal `json:"fdv"`
		MarketCap                    decimal.Decimal `json:"market_cap"`
		Price                        decimal.Decimal `json:"price"`
		PriceChange24HPercent        decimal.Decimal `json:"price_change_24h_percent"`
		TradeBuy24H                  int             `json:"trade_buy_24h"`
		TradeBuy24HChangePercent     decimal.Decimal `json:"trade_buy_24h_change_percent"`
		TradeSell24H                 int             `json:"trade_sell_24h"`
		TradeSell24HChangePercent    decimal.Decimal `json:"trade_sell_24h_change_percent"`
		Trade24H                     int             `json:"trade_24h"`
		Trade24HChangePercent        decimal.Decimal `json:"trade_24h_change_percent"`
		UniqueWallet24H              int             `json:"unique_wallet_24h"`
		UniqueWallet24HChangePercent decimal.Decimal `json:"unique_wallet_24h_change_percent"`
		Volume24HUsd                 decimal.Decimal `json:"volume_24h_usd"`
		Volume24HChangePercent       decimal.Decimal `json:"volume_24h_change_percent"`
		LastTradeUnixTime            uint64          `json:"last_trade_unix_time"`
		Source                       string          `json:"source"`
	}
	result := make([]searchResult, 0)
	for _, t := range tokenResult {
		result = append(result, searchResult{
			Chain:                        t.Network,
			Type:                         "token",
			CA:                           t.Address,
			Name:                         t.Name,
			Symbol:                       t.Symbol,
			Logo:                         t.LogoURI,
			Liquidity:                    t.Liquidity,
			TotalSupply:                  t.Supply,
			FDV:                          t.FDV,
			MarketCap:                    t.MarketCap,
			Price:                        t.Price,
			PriceChange24HPercent:        t.PriceChange24hPercent,
			TradeBuy24H:                  t.Buy24h,
			TradeBuy24HChangePercent:     t.Buy24hChangePercent,
			TradeSell24H:                 t.Sell24h,
			TradeSell24HChangePercent:    t.Sell24hChangePercent,
			Trade24H:                     t.Trade24h,
			Trade24HChangePercent:        t.Trade24hChangePercent,
			UniqueWallet24H:              t.UniqueWallet24h,
			UniqueWallet24HChangePercent: t.UniqueWallet24hChangePercent,
			Volume24HUsd:                 t.Volume24hUSD,
			Volume24HChangePercent:       t.Volume24hChangePercent,
			LastTradeUnixTime:            uint64(t.LastTradeUnixTime),
		})
	}
	for _, m := range marketResult {
		result = append(result, searchResult{
			Chain:                        m.Network,
			Type:                         "market",
			CA:                           m.Address,
			Liquidity:                    m.Liquidity,
			UniqueWallet24H:              m.UniqueWallet24h,
			UniqueWallet24HChangePercent: m.UniqueWallet24hChangePercent,
			Trade24H:                     m.Trade24h,
			TradeBuy24HChangePercent:     m.Trade24hChangePercent,
			Volume24HUsd:                 m.Volume24hUSD,
			LastTradeUnixTime:            uint64(m.LastTradeUnixTime),
			Source:                       m.Source,
		})
	}
	res.Data = result

	c.JSON(http.StatusOK, res)
}
