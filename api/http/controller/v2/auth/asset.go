package auth

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stonksdex/externalapi/api/common"
	"github.com/stonksdex/externalapi/codes"
	"github.com/stonksdex/externalapi/log"
	"github.com/stonksdex/externalapi/model"
	"github.com/stonksdex/externalapi/system"
	"github.com/stonksdex/externalapi/thirdpart"
)

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

	be, err := thirdpart.GetBirdClient()
	if err != nil {
		res.Code = codes.CODE_ERR_UNKNOWN
		res.Msg = "service unavailable " + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	assets, err := be.WalletPortfolio("solana", userWallet.Wallet)
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

func AssetTokenTrans(c *gin.Context) {
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
	if len(token) == 0 {
		res.Code = codes.CODE_ERR_REQFORMAT
		res.Msg = "please pass token address"
		c.JSON(http.StatusOK, res)
		return
	}
	pn, _ := strconv.ParseInt(c.Query("pn"), 10, 64)
	ps, _ := strconv.ParseInt(c.Query("ps"), 10, 64)
	if pn < 1 {
		pn = 1
	}
	if ps < 1 {
		ps = 1
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

	log.Info("query user wallet: ", userWallet.Wallet)
	be, err := thirdpart.GetBirdClient()
	tokenTrades, err := be.TokenTxsByTokenAddress(token, "solana", "swap", "", userWallet.Wallet, "", 0, 0, int(pn-1), int(ps))

	res.Code = codes.CODE_SUCCESS
	res.Msg = "success"
	res.Data = tokenTrades
	c.JSON(http.StatusOK, res)
}
