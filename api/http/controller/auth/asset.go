package auth

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stonksdex/externalapi/api/common"
	"github.com/stonksdex/externalapi/codes"
	"github.com/stonksdex/externalapi/model"
	"github.com/stonksdex/externalapi/system"
	"github.com/stonksdex/externalapi/thirdpart"
)

type pnl struct {
	Value   string `json:"value"`
	Percent string `json:"percent"`
}
type AssetBoard struct {
	TotalValue  string `json:"total_value"`
	TotalIncome pnl    `json:"total_income"`
	TodayPAL    pnl    `json:"today_pal"`
	SevenPAL    pnl    `json:"7_pal"`
	ThirtyPAL   pnl    `json:"30_pal"`
}

func AssetView(c *gin.Context) {
	res := common.Response{}
	res.Timestamp = time.Now().Unix()

	currentUser, exist := c.Get("user_id")

	if !exist {
		res.Code = codes.CODE_ERR_AUTHTOKEN_FAIL
		res.Msg = "token invalid, please relogin"
		c.JSON(http.StatusOK, res)
		return
	}
	currentUserStr, _ := currentUser.(string)
	userID, err := strconv.ParseInt(currentUserStr, 10, 64)
	if err != nil {
		res.Code = codes.CODE_ERR_REQFORMAT
		res.Msg = "token invalid, please relogin"
		c.JSON(http.StatusOK, res)
		return
	}

	db := system.GetDb()
	var userWallet model.UserWallet
	db.Model(&model.UserWallet{}).Where("id = ?", userID).First(&userWallet)
	if userWallet.ID == 0 {
		res.Code = codes.CODE_ERR_TX
		res.Msg = "need login"
		c.JSON(http.StatusOK, res)
		return
	}

	okx, err := thirdpart.GetOKXClient()
	if err != nil {
		res.Code = codes.CODE_ERR_UNKNOWN
		res.Msg = "service unavailable " + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	assets, err := okx.WalletValueSpec(userWallet.Wallet, "501", "0", true)
	if err != nil {
		res.Code = codes.CODE_ERR_UNKNOWN
		res.Msg = "unable retrieve assets " + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	totalValue := "0"
	if len(assets) > 0 {
		totalValue = assets[0].TotalValue
	}

	res.Code = codes.CODE_SUCCESS
	res.Msg = "success"
	res.Data = AssetBoard{
		TotalValue: totalValue,
		TotalIncome: pnl{
			Value:   "1",
			Percent: "0.88",
		},
		TodayPAL: pnl{
			Value:   "1",
			Percent: "-1",
		},
		ThirtyPAL: pnl{
			Value:   "2",
			Percent: "1300",
		},
		SevenPAL: pnl{
			Value:   "2",
			Percent: "-230",
		},
	}
	c.JSON(http.StatusOK, res)
}

func AssetList(c *gin.Context) {
	res := common.Response{}
	res.Timestamp = time.Now().Unix()

	currentUser, exist := c.Get("user_id")

	if !exist {
		res.Code = codes.CODE_ERR_AUTHTOKEN_FAIL
		res.Msg = "token invalid, please relogin"
		c.JSON(http.StatusOK, res)
		return
	}
	currentUserStr, _ := currentUser.(string)
	userID, err := strconv.ParseInt(currentUserStr, 10, 64)
	if err != nil {
		res.Code = codes.CODE_ERR_REQFORMAT
		res.Msg = "token invalid, please relogin"
		c.JSON(http.StatusOK, res)
		return
	}

	db := system.GetDb()
	var userWallet model.UserWallet
	db.Model(&model.UserWallet{}).Where("id = ?", userID).First(&userWallet)
	if userWallet.ID == 0 {
		res.Code = codes.CODE_ERR_TX
		res.Msg = "need login"
		c.JSON(http.StatusOK, res)
		return
	}

	okx, err := thirdpart.GetOKXClient()
	if err != nil {
		res.Code = codes.CODE_ERR_UNKNOWN
		res.Msg = "service unavailable " + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	assets, err := okx.WalletAssets(userWallet.Wallet, "501", "0")
	if err != nil {
		res.Code = codes.CODE_ERR_UNKNOWN
		res.Msg = "unable retrieve assets " + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	res.Code = codes.CODE_SUCCESS
	res.Msg = "success"
	res.Data = assets
	c.JSON(http.StatusOK, res)
}

func AssetTrans(c *gin.Context) {
	res := common.Response{}
	res.Timestamp = time.Now().Unix()

	currentUser, exist := c.Get("user_id")

	if !exist {
		res.Code = codes.CODE_ERR_AUTHTOKEN_FAIL
		res.Msg = "token invalid, please relogin"
		c.JSON(http.StatusOK, res)
		return
	}
	currentUserStr, _ := currentUser.(string)
	userID, err := strconv.ParseInt(currentUserStr, 10, 64)
	if err != nil {
		res.Code = codes.CODE_ERR_REQFORMAT
		res.Msg = "token invalid, please relogin"
		c.JSON(http.StatusOK, res)
		return
	}
	token := c.Query("token")

	db := system.GetDb()
	var userWallet model.UserWallet
	db.Model(&model.UserWallet{}).Where("id = ?", userID).First(&userWallet)
	if userWallet.ID == 0 {
		res.Code = codes.CODE_ERR_TX
		res.Msg = "need login"
		c.JSON(http.StatusOK, res)
		return
	}

	okx, err := thirdpart.GetOKXClient()
	if err != nil {
		res.Code = codes.CODE_ERR_UNKNOWN
		res.Msg = "service unavailable " + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	assets, err := okx.WalletTransactions(userWallet.Wallet, "501", token, 0, 0, 0)
	if err != nil {
		res.Code = codes.CODE_ERR_UNKNOWN
		res.Msg = "unable retrieve assets " + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	res.Code = codes.CODE_SUCCESS
	res.Msg = "success"
	res.Data = assets
	c.JSON(http.StatusOK, res)
}
