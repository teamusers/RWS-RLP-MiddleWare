package home

import (
	"net/http"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stonksdex/externalapi/log"
	"github.com/stonksdex/externalapi/model"
	"github.com/stonksdex/externalapi/thirdpart"

	"github.com/gin-gonic/gin"
	"github.com/stonksdex/externalapi/api/common"
	"github.com/stonksdex/externalapi/codes"
)

func K(c *gin.Context) {
	res := common.Response{
		Timestamp: time.Now().Unix(),
		Code:      codes.CODE_SUCCESS,
		Msg:       "success",
	}

	chain := c.Param("chain")
	ca := c.Param("ca")
	fromStr := c.Query("from")
	toStr := c.Query("to")
	resolution := c.Query("resolution")
	bird, _ := thirdpart.GetBirdClient()

	from, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		res.Code = codes.CODE_ERR_BAD_PARAMS
		res.Msg = "from is invalid"
		c.JSON(http.StatusOK, res)
		return
	}
	to, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		to = time.Now().Unix()
	}

	result, err := bird.FetchKlineData(ca, chain, resolution, "usd", from, to)
	if err != nil {
		log.Error("[Birdeye] fetch kline error: ", err)
	}
	type k struct {
		Time   uint64          `json:"time"`
		Open   decimal.Decimal `json:"open"`
		High   decimal.Decimal `json:"high"`
		Low    decimal.Decimal `json:"low"`
		Close  decimal.Decimal `json:"close"`
		Volume decimal.Decimal `json:"volume"`
	}
	kd := make([]k, 0)

	for _, v := range result {
		kd = append(kd, k{
			Time:   v.UnixTime,
			Open:   v.O,
			High:   v.H,
			Low:    v.L,
			Close:  v.C,
			Volume: v.V,
		})
	}
	res.Data = kd
	c.JSON(http.StatusOK, res)
}

func timeGroupSQL(resolution string) (string, int) {
	switch resolution {
	case "1m":
		return "(UNIX_TIMESTAMP(`time`) DIV 60)  ", 60
	case "5m":
		return "(UNIX_TIMESTAMP(`time`) DIV 300)  ", 300
	case "10m":
		return "(UNIX_TIMESTAMP(`time`) DIV 600 )  ", 600
	case "15m":
		return "(UNIX_TIMESTAMP(`time`) DIV 900 ) ", 900
	case "30m":
		return "(UNIX_TIMESTAMP(`time`) DIV 1800  ) ", 1800
	case "1h":
		return "(UNIX_TIMESTAMP(`time`) DIV 3600 ) ", 3600
	case "4h":
		return "(UNIX_TIMESTAMP(`time`) DIV 14400) ", 14400
	case "6h":
		return "(UNIX_TIMESTAMP(`time`) DIV 21600) ", 21600
	case "1d":
		return "(UNIX_TIMESTAMP(`time`) DIV 86400) ", 86400
	default:
		return "(UNIX_TIMESTAMP(`time`) DIV 900)  ", 900
	}
}

// smoothKChart 用于平滑K线数据
func smoothKChart(kCharts []model.KlineData) []model.KlineData {
	if len(kCharts) <= 1 {
		return kCharts // 数据不足，直接返回
	}

	// Reverse slices
	//reversedCharts := reverseKCharts(kCharts)

	for i := 1; i < len(kCharts); i++ {
		// Get the current K-line and the previous K-line
		current := &kCharts[i]
		previous := kCharts[i-1]

		// Smoothing: Set the opening price of the current K-line to the closing price of the previous K-line
		open := current.Open
		current.Open = previous.Close
		if current.High.Cmp(open) == -1 {
			current.High = open
		}
		if current.Low.Cmp(open) == 1 {
			current.Low = open
		}

		// If High == Low == Close, force adjustment of Open
		// if current.High.Cmp(current.Low) == 0 && current.Low.Cmp(current.Close) == 0 {
		// 	current.Open = current.High // Open 必须等于 High
		// } else {
		// 	// Other situations: Make sure Open is between Low and High
		// 	if current.Open.Cmp(current.High) >= 0 {
		// 		current.Open = current.High
		// 	} else if current.Open.Cmp(current.Low) <= 0 {
		// 		current.Open = current.Low
		// 	}
		// }
	}

	//return reverseKCharts(reversedCharts)
	return kCharts
}

func smoothKChartGPT(kCharts []model.KlineData) []model.KlineData {
	if len(kCharts) <= 1 {
		return kCharts // 数据不足，直接返回
	}

	for i := 1; i < len(kCharts); i++ {
		current := &kCharts[i]
		previous := kCharts[i-1]

		// 继承前一根K线的Close作为当前Open
		current.Open = previous.Close

		// 处理特殊情况：High == Low == Close
		if current.High.Cmp(current.Low) == 0 && current.Low.Cmp(current.Close) == 0 {
			current.Open = current.Close
		}

		// 确保 Open 在 High 和 Low 之间
		if current.Open.Cmp(current.High) > 0 {
			current.High = current.Open
		}
		if current.Open.Cmp(current.Low) < 0 {
			current.Low = current.Open
		}
	}

	return kCharts
}

// reverseKCharts
func reverseKCharts(kCharts []model.KlineData) []model.KlineData {
	n := len(kCharts)
	reversed := make([]model.KlineData, n)
	for i := 0; i < n; i++ {
		reversed[i] = kCharts[n-i-1]
	}
	return reversed
}
