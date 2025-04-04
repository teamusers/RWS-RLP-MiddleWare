package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stonksdex/externalapi/api/common"
	"github.com/stonksdex/externalapi/codes"
	"github.com/stonksdex/externalapi/log"
	"github.com/stonksdex/externalapi/model"
	"github.com/stonksdex/externalapi/system"
	"github.com/stonksdex/externalapi/thirdpart"
)

type CreateTranQuote struct {
	FromToken string `json:"from_token"`
	ToToken   string `json:"to_token"`
	Amount    uint64 `json:"amount" binding:"min=1"`
	Slippage  uint64 `json:"slippage" binding:"min=1,max=10000"`
}

type TransNotify struct {
	TxID   uint64 `json:"tx_id"`
	TxHash string `json:"tx_hash"`
}

func Quote(c *gin.Context) {
	var req CreateTranQuote
	res := common.Response{}
	res.Timestamp = time.Now().Unix()

	if err := c.ShouldBindJSON(&req); err != nil {
		res.Code = codes.CODE_ERR_REQFORMAT
		res.Msg = "invalid request" + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	okx, _ := thirdpart.GetOKXClient()
	result, err := okx.Quote(req.FromToken, req.ToToken, fmt.Sprintf("%d", req.Amount), decimal.NewFromInt(int64(req.Slippage)).Div(decimal.NewFromInt(10000)).Round(2).String())
	if err != nil {
		res.Code = codes.CODE_ERR_OKX
		res.Msg = "Quote Failed: " + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	res.Code = codes.CODE_SUCCESS
	res.Msg = "success"
	res.Data = result
	c.JSON(http.StatusOK, res)
}

func Trans(c *gin.Context) {
	var req CreateTranQuote
	res := common.Response{}
	res.Timestamp = time.Now().Unix()

	currentUser, exist := c.Get("user_id")

	if !exist {
		res.Code = codes.CODE_ERR_AUTHTOKEN_FAIL
		res.Msg = "token invalid, please relogin"
		c.JSON(http.StatusOK, res)
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		res.Code = codes.CODE_ERR_REQFORMAT
		res.Msg = "invalid request" + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	currentUserStr, _ := currentUser.(string)
	userID, err := strconv.ParseUint(currentUserStr, 10, 64)
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

	okx, _ := thirdpart.GetOKXClient()
	result, err := okx.Swap(req.FromToken, req.ToToken,
		fmt.Sprintf("%d", req.Amount),
		decimal.NewFromInt(int64(req.Slippage)).Div(decimal.NewFromInt(10000)).Round(2).String(),
		userWallet.Wallet)

	userTrans := model.TxsDiagram{
		UserID:     userID,
		Chain:      "solana",
		FromToken:  req.FromToken,
		ToToken:    req.ToToken,
		FromAmount: req.Amount,
		Slippage:   uint64(req.Slippage),
		Status:     "00",
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	defer func(e error) {
		if userTrans.Status != "00" {
			db.Save(&model.TxsLog{
				TxId:       userTrans.ID,
				Step:       1,
				ErrMsg:     e.Error(),
				CreateTime: time.Now(),
			})
		}
	}(err)

	db.Save(&userTrans)
	if err != nil {
		userTrans.Status = "90"
		res.Code = codes.CODE_ERR_OKX
		res.Msg = "Init Swap Failed: " + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	res.Code = codes.CODE_SUCCESS
	res.Msg = "success"
	res.Data = struct {
		TxID   uint64 `json:"tx_id"`
		TxData any    `json:"tx_data"`
	}{
		TxID:   userTrans.ID,
		TxData: result,
	}
	c.JSON(http.StatusOK, res)
}

func Notify(c *gin.Context) {
	var req TransNotify
	res := common.Response{}
	res.Timestamp = time.Now().Unix()

	currentUser, exist := c.Get("user_id")

	if !exist {
		res.Code = codes.CODE_ERR_AUTHTOKEN_FAIL
		res.Msg = "token invalid, please relogin"
		c.JSON(http.StatusOK, res)
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		res.Code = codes.CODE_ERR_REQFORMAT
		res.Msg = "invalid request" + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	currentUserStr, _ := currentUser.(string)
	userID, err := strconv.ParseUint(currentUserStr, 10, 64)
	if err != nil {
		res.Code = codes.CODE_ERR_REQFORMAT
		res.Msg = "token invalid, please relogin"
		c.JSON(http.StatusOK, res)
		return
	}

	if len(req.TxHash) < 60 {
		res.Code = codes.CODE_ERR_REQFORMAT
		res.Msg = "hash length invalid"
		c.JSON(http.StatusOK, res)
		return
	}

	db := system.GetDb()
	var userTrans model.TxsDiagram
	db.Model(&model.TxsDiagram{}).Where("id = ?", req.TxID).First(&userTrans)
	if userTrans.ID == 0 {
		log.Errorf("[NotifyTrans] tx_id:%d, tx_hash:%s", req.TxID, req.TxHash)
		res.Code = codes.CODE_ERR_TX
		res.Msg = "can not find tx data"
		c.JSON(http.StatusOK, res)
		return
	}

	if userTrans.UserID != userID {
		log.Errorf("[NotifyTrans] session_userid:%d, tx_userid:%d", userID, userTrans.UserID)
	}

	userTrans.TxHash = req.TxHash
	userTrans.UpdateTime = time.Now()
	userTrans.Status = "10"
	db.Save(&userTrans)

	redis := system.GetRedis()
	if redis != nil {
		log.Info("[NotifyTrans] store hash to redis: ", req.TxHash)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		redis.Set(ctx, "txnotify:"+req.TxHash, "", 5*time.Minute)
	}

	res.Code = codes.CODE_SUCCESS
	res.Msg = "success"
	c.JSON(http.StatusOK, res)
}
