package home

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stonksdex/externalapi/log"
	"github.com/stonksdex/externalapi/model"
	"github.com/stonksdex/externalapi/system"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/stonksdex/externalapi/api/common"
	"github.com/stonksdex/externalapi/codes"
	"github.com/stonksdex/externalapi/thirdpart"
)

func TokenHolders(c *gin.Context) {
	res := common.Response{
		Timestamp: time.Now().Unix(),
		Code:      codes.CODE_SUCCESS,
		Msg:       "success",
	}

	_ = c.Param("chain")
	ca := c.Param("ca")
	pnStr := c.Query("pn")
	psStr := c.Query("ps")

	pn, err := strconv.ParseInt(pnStr, 10, 64)
	if err != nil || pn < 1 {
		pn = 1
	}
	ps, err := strconv.ParseInt(psStr, 10, 64)
	if err != nil || ps < 1 {
		ps = 1
	}

	b, _ := thirdpart.GetBirdClient()

	result, err := b.TokenHolder("solana", ca, int(pn-1), int(ps))
	if err != nil {
		log.Error("get token holders error: ", err)
	}

	type BirdTokenHold struct {
		Address      string `json:"address" gorm:"column:address"`
		Amount       string `json:"amount" gorm:"column:amount"`
		Decimals     uint64 `json:"decimals" gorm:"column:decimals"`
		TokenAccount string `json:"token_account"`
	}
	retData := make([]BirdTokenHold, 0)

	if len(result) > 0 {
		for _, v := range result {

			retData = append(retData, BirdTokenHold{
				Address:      v.Owner,
				Amount:       v.Amount,
				Decimals:     uint64(v.Decimals),
				TokenAccount: v.TokenAccount,
			})
		}
	}
	res.Data = retData
	c.JSON(http.StatusOK, res)
}

var pMap = []string{"5m", "1h", "4h", "1d"}

func TokenInfoV2(c *gin.Context) {
	res := common.Response{
		Timestamp: time.Now().Unix(),
		Code:      codes.CODE_SUCCESS,
		Msg:       "success",
	}

	chain := c.Param("chain")
	ca := c.Param("ca")

	//log.Infof("TokenInfoChange: %v ", chg)

	// Token info
	mdb := system.GetDb()
	var existMeta model.TokenMeta
	err := mdb.Model(&model.TokenMeta{}).Where("chain = ? and ca = ?", chain, ca).First(&existMeta).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("Unexpected:", err)
	}

	if existMeta.ID == 0 {
		// need insert data
		result, err := thirdpart.GetMetaInfoFromGoPlus(ca)
		if err != nil {
			log.Error("GOPLUS getMetaInfoFromGoPlus: ", err)
			res.Code = codes.CODE_ERR_REMOTE
			res.Msg = "error to get token info: " + err.Error()
			c.JSON(http.StatusOK, res)
			return
		}

		goplusToken := result.Result[ca]
		existMeta, existPairs := thirdpart.CopyMetaIntoObject(goplusToken)
		existMeta.CA = ca
		existMeta.Chain = chain

		logo, website, xUri, tgUri := checkLogAndUri(existMeta.MetaContent)
		if len(existMeta.Logo) == 0 {
			existMeta.Logo = logo
		}
		if len(existMeta.Website) == 0 {
			existMeta.Website = website
		}
		if len(existMeta.XURI) == 0 {
			existMeta.XURI = xUri
		}
		if len(existMeta.TgURI) == 0 {
			existMeta.TgURI = tgUri
		}

		err = mdb.Save(&existMeta).Error
		for _, vs := range existPairs {
			mdb.Save(&vs)
		}
		log.Error(err)
	}
	view := GetTokenOverView(ca, chain)
	if view != nil {
		existMeta.Price = view.Price
		existMeta.TotalSupply = view.TotalSupply
		existMeta.Circulation = view.MarketCap
		existMeta.Liquidity = view.Liquidity
		existMeta.CirculatingSupply = view.CirculatingSupply
		existMeta.HolderCount = uint64(view.Holder)
		existMeta.Change5m = view.PriceChange1hPercent
		existMeta.Change1h = view.PriceChange1hPercent
		existMeta.Change4h = view.PriceChange4hPercent
		existMeta.Change24h = view.PriceChange24hPercent
		if len(existMeta.Logo) == 0 {
			existMeta.Logo = view.LogoURI
		}
		if len(existMeta.Website) == 0 {
			existMeta.Website = view.Extensions.Website
		}
		if len(existMeta.XURI) == 0 {
			existMeta.XURI = view.Extensions.Twitter
		}
		if len(existMeta.TgURI) == 0 {
			existMeta.TgURI = view.Extensions.Telegram
		}
		if len(existMeta.Description) == 0 {
			existMeta.Description = view.Extensions.Description
		}

	}
	tokenSecurity := GetTokenSecurity(ca, chain)
	if tokenSecurity != nil {
		existMeta.Top10HolderBalance = tokenSecurity.Top10HolderBalance
		if tokenSecurity.Top10HolderPercent.Sign() > 0 {
			existMeta.Top10HolderBalance = tokenSecurity.Top10HolderBalance.Mul(decimal.NewFromInt(100))
		}
		existMeta.Top10HolderPercent = tokenSecurity.Top10HolderPercent
		existMeta.Top10UserBalance = tokenSecurity.Top10UserBalance
		existMeta.Top10UserPercent = tokenSecurity.Top10UserPercent
		if tokenSecurity.CreatorBalance != nil {
			existMeta.DevBalance = *tokenSecurity.CreatorBalance
			existMeta.DevPercentage = *tokenSecurity.CreatorPercentage
		}
		existMeta.Create = tokenSecurity.CreatorAddress
	}
	existMeta.TagList = []string{}

	res.Data = existMeta
	c.JSON(http.StatusOK, res)
}
func TokenChgV2(c *gin.Context) {
	res := common.Response{
		Timestamp: time.Now().Unix(),
		Code:      codes.CODE_SUCCESS,
		Msg:       "success",
	}

	chain := c.Param("chain")
	ca := c.Param("ca")
	if len(ca) <= 0 || len(chain) <= 0 {
		res.Code = codes.CODE_ERR_REMOTE
		res.Msg = "error to get token info:  chain or ca is empty"
		c.JSON(http.StatusOK, res)
		return
	}
	v := GetTokenOverView(ca, chain)
	var periods = []string{"30m", "1h", "12h", "1d"}
	periodRes := make([]model.TokenChgByPeriod, 0)
	for i, period := range periods {
		if i == 0 {
			periodRes = append(periodRes, model.TokenChgByPeriod{
				Period:             period,
				StartPrice:         v.History30mPrice,
				EndPrice:           v.Price,
				PriceChangePercent: v.PriceChange30mPercent,
				TradeCount:         v.Trade30m,
				TradeBuyCount:      v.Buy30m,
				TradeSellCount:     v.Sell30m,
				TradeBuy:           v.VBuy30m,
				TradeSell:          v.VSell30m,
				TradeVol:           v.V30m,
			})
		}
		if i == 1 {
			periodRes = append(periodRes, model.TokenChgByPeriod{
				Period:             period,
				StartPrice:         v.History1hPrice,
				EndPrice:           v.Price,
				PriceChangePercent: v.PriceChange1hPercent,
				TradeCount:         v.Trade1h,
				TradeBuyCount:      v.Buy1h,
				TradeSellCount:     v.Sell1h,
				TradeBuy:           v.VBuy1h,
				TradeSell:          v.VSell1h,
				TradeVol:           v.V1h,
			})
		}
		if i == 2 {
			periodRes = append(periodRes, model.TokenChgByPeriod{
				Period:             period,
				StartPrice:         v.History4hPrice,
				EndPrice:           v.Price,
				PriceChangePercent: v.PriceChange4hPercent,
				TradeCount:         v.Trade4h,
				TradeBuyCount:      v.Buy4h,
				TradeSellCount:     v.Sell4h,
				TradeBuy:           v.VBuy4h,
				TradeSell:          v.VSell4h,
				TradeVol:           v.V4h,
			})
		}
		if i == 3 {
			periodRes = append(periodRes, model.TokenChgByPeriod{
				Period:             period,
				StartPrice:         v.History24hPrice,
				EndPrice:           v.Price,
				PriceChangePercent: v.PriceChange24hPercent,
				TradeCount:         v.Trade24h,
				TradeBuyCount:      v.Buy24h,
				TradeSellCount:     v.Sell24h,
				TradeBuy:           v.VBuy24h,
				TradeSell:          v.VSell24h,
				TradeVol:           v.V24h,
			})
		}
	}

	res.Data = periodRes
	c.JSON(http.StatusOK, res)

}

func GetTokenOverView(ca, chain string) *thirdpart.BeTokenOverview {
	bird, _ := thirdpart.GetBirdClient()
	overview, err := bird.TokenOverview(ca, chain)
	if err != nil {
		log.Error("GetTokenOverView: ", err)
	}
	return overview
}
func GetTokenSecurity(ca, chain string) *thirdpart.BeTokenSecurity {
	bird, _ := thirdpart.GetBirdClient()
	overview, err := bird.TokenSecurity(ca, chain)
	if err != nil {
		log.Error("GetTokenOverView: ", err)
	}
	return overview
}

func checkLogAndUri(content string) (string, string, string, string) {
	mapMeta := make(map[string]interface{})
	err := json.Unmarshal([]byte(content), &mapMeta)
	if err != nil {
		return "", "", "", ""
	}
	logo, website, xUri, tgUri := "", "", "", ""
	if err == nil && len(mapMeta) > 0 {
		v, ex := mapMeta["image"]
		if ex {
			logo = v.(string)
		}
		v, ex = mapMeta["website"]
		if ex {
			website = v.(string)
		}
		v, ex = mapMeta["twitter"]
		if ex {
			xUri = v.(string)
		}
		v, ex = mapMeta["telegram"]
		if ex {
			tgUri = v.(string)
		}
	}
	return logo, website, xUri, tgUri
}

func TokenNewList(c *gin.Context) {
	res := common.Response{
		Timestamp: time.Now().Unix(),
		Code:      codes.CODE_SUCCESS,
		Msg:       "success",
	}

	_ = c.Param("chain")
	memeplat := c.Query("memeplat")
	memeplatInclude := true
	if memeplat == "0" {
		memeplatInclude = false
	}

	b, _ := thirdpart.GetBirdClient()

	result, err := b.TokenNewList("solana", uint64(time.Now().Unix()), memeplatInclude, 20)
	if err != nil {
		log.Error("get token holders error: ", err)
	}

	res.Data = result
	c.JSON(http.StatusOK, res)
}
