package home

import (
	"github.com/stonksdex/externalapi/thirdpart"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stonksdex/externalapi/api/common"
	"github.com/stonksdex/externalapi/codes"
)

type SolPairFlow struct {
	Tx                  string          `json:"tx" gorm:"column:tx"`
	TxTime              string          `json:"tx_time" gorm:"column:tx_time"`
	Pair                string          `json:"pair" gorm:"column:pair"`
	Block               uint64          `json:"block" gorm:"column:block"`
	Token0              string          `json:"token0" gorm:"column:token0"`
	Token1              string          `json:"token1" gorm:"column:token1"`
	Symbol0             string          `json:"symbol0" gorm:"-"`
	Symbol1             string          `json:"symbol1" gorm:"-"`
	Amount0             decimal.Decimal `json:"amount0" gorm:"column:amount0"`
	Amount1             decimal.Decimal `json:"amount1" gorm:"column:amount1"`
	Liquidity           decimal.Decimal `json:"liquidity" gorm:"column:liquidity"`
	AddLiquidityAddress string          `json:"add_liquidity_address" gorm:"column:add_liquidity_address"`
	Type                string          `json:"type" gorm:"column:type"`
	Dex                 string          `json:"dex" gorm:"column:dex"`
	Price               decimal.Decimal `json:"price" gorm:"column:price"`
	Value               decimal.Decimal `json:"value" gorm:"-"`
}

func PairFlowV2(c *gin.Context) {
	res := common.Response{
		Timestamp: time.Now().Unix(),
		Code:      codes.CODE_SUCCESS,
		Msg:       "success",
	}

	chain := c.Param("chain")
	ca := c.Param("ca")
	pnStr := c.Query("pn")
	psStr := c.Query("ps")
	beforeStr := c.Query("should_before")

	pn, err1 := strconv.ParseInt(pnStr, 10, 64)
	ps, err2 := strconv.ParseInt(psStr, 10, 64)
	before, err3 := strconv.ParseInt(beforeStr, 10, 64)
	if err1 != nil {
		pn = 1
	}
	if err2 != nil {
		ps = 50
	}
	if err3 != nil {
		before = time.Now().Unix()
	}
	offset := (pn - 1) * ps
	flow := getPairFlow(ca, chain, offset, ps, before)
	if len(flow) <= 0 {
		res.Data = []struct{}{}
		c.JSON(http.StatusOK, res)
		return
	}
	rea := make([]SolPairFlow, 0)
	for _, v := range flow {
		txType := "MINT"
		if v.TxType == "remove" {
			txType = "BURN"
		}
		token0, token1 := v.Tokens[0], v.Tokens[1]
		format := time.Unix(int64(v.BlockUnixTime), 0).Format("2006-01-02 15:04:05")
		rea = append(rea, SolPairFlow{
			Tx:                  v.TxHash,
			TxTime:              format,
			Pair:                v.PoolId,
			Block:               0,
			Token0:              token0.Address,
			Token1:              token1.Address,
			Amount0:             token0.UiAmount,
			Amount1:             token1.UiAmount,
			Symbol0:             token0.Symbol,
			Symbol1:             token1.Symbol,
			Liquidity:           decimal.NewFromInt(0),
			AddLiquidityAddress: v.Owner,
			Type:                txType,
			Dex:                 v.Source,
			Price:               v.PricePair,
			Value:               v.VolumeUsd,
		})

	}
	res.Data = rea
	c.JSON(http.StatusOK, res)
}

func getPairFlow(ca string, chain string, offset int64, ps int64, before int64) []thirdpart.BeTokenTrade {
	b, _ := thirdpart.GetBirdClient()
	adds, err := b.TokenTxsByTokenAddress(ca, "solana", "add", "", "", "", before, 0, int(offset), int(ps))
	removes, err1 := b.TokenTxsByTokenAddress(ca, "solana", "remove", "", "", "", before, 0, int(offset), int(ps))
	if err != nil && err1 != nil {
		return []thirdpart.BeTokenTrade{}
	}
	addsLen := len(adds)
	removesLen := len(removes)
	temp := make([]thirdpart.BeTokenTrade, 0)
	temp = append(temp, adds...)
	temp = append(temp, removes...)
	// There is enough data to directly return the first n pieces of data merged
	if addsLen > 0 && removesLen > 0 && addsLen == removesLen {
		res := make([]thirdpart.BeTokenTrade, 0)
		res = temp
		sort.SliceStable(res, func(i, j int) bool {
			return res[i].BlockUnixTime > res[j].BlockUnixTime // 降序排列
		})
		if len(res) > 10 && addsLen < 10 {
			res = res[:10]
		}
		if len(res) > addsLen {
			res = res[:addsLen]
		}
		return res
	}
	// Insufficient data and all data are all returned
	if addsLen > 0 && removesLen > 0 && addsLen != removesLen && addsLen+removesLen < 10 {
		res := make([]thirdpart.BeTokenTrade, 0)
		res = temp
		sort.SliceStable(res, func(i, j int) bool {
			return res[i].BlockUnixTime > res[j].BlockUnixTime // 降序排列
		})

		return res
	}
	//  v  If there is not enough data, then all the merged data will be returned
	if addsLen > 0 && removesLen > 0 && addsLen != removesLen && addsLen+removesLen > 10 {
		shouldLen := int(math.Min(float64(addsLen), float64(removesLen)))
		var targetTime int
		if shouldLen == addsLen {
			targetTime = adds[shouldLen-1].BlockUnixTime
		}
		if shouldLen == removesLen {
			targetTime = removes[shouldLen-1].BlockUnixTime
		}
		res := make([]thirdpart.BeTokenTrade, 0)
		for _, add := range adds {
			if add.BlockUnixTime >= targetTime {
				res = append(res, add)
			}

		}
		for _, add := range removes {
			if add.BlockUnixTime >= targetTime {
				res = append(res, add)
			}
		}
		sort.SliceStable(res, func(i, j int) bool {
			return res[i].BlockUnixTime > res[j].BlockUnixTime // 降序排列
		})
		if len(res) >= 10 {
			res = res[:10]
			return res
		}
		return res

	}
	sort.SliceStable(temp, func(i, j int) bool {
		return temp[i].BlockUnixTime > temp[j].BlockUnixTime // 降序排列
	})
	return temp
}
