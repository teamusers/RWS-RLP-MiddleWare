package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stonksdex/externalapi/api/common"
	"github.com/stonksdex/externalapi/api/http/request"
	"github.com/stonksdex/externalapi/codes"
	"github.com/stonksdex/externalapi/model"
	"github.com/stonksdex/externalapi/system"

	"gorm.io/gorm"
)

func Ref(c *gin.Context) {
	res := common.Response{}
	res.Timestamp = time.Now().Unix()
	res.Code = codes.CODE_SUCCESS
	res.Msg = "success"
	res.Data = nil

	_, _ = c.Get("user_wallet")
	userId, ok := c.Get("user_id")

	if !ok {
		res.Code = codes.CODE_ERR_SECURITY
		res.Msg = "need login"
		c.JSON(http.StatusOK, res)
		return
	}

	var r model.UserRef
	db := system.GetDb()
	err := db.Model(&model.UserRef{}).Where("uw_id = ?", userId).Find(&r).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			res.Code = codes.CODE_ERR_UNKNOWN
			res.Msg = "unknown " + err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
	}
	userIdInt, _ := strconv.ParseInt(fmt.Sprintf("%v", userId), 10, 64)
	if r.ID == 0 {
		r = model.UserRef{
			UwID:       userIdInt,
			RefCode:    system.GenerateNonce(8) + system.GenerateNonce(4),
			CreateTime: time.Now(),
		}
		db.Save(&r)
	}
	res.Data = struct {
		Ref string `json:"ref"`
	}{
		Ref: r.RefCode,
	}

	c.JSON(http.StatusOK, res)
}

func RefCount(c *gin.Context) {
	res := common.Response{}
	res.Timestamp = time.Now().Unix()
	res.Code = codes.CODE_SUCCESS
	res.Msg = "success"
	res.Data = nil

	_, _ = c.Get("user_wallet")
	userId, ok := c.Get("user_id")

	if !ok {
		res.Code = codes.CODE_ERR_SECURITY
		res.Msg = "need login"
		c.JSON(http.StatusOK, res)
		return
	}

	var r model.UserRef
	db := system.GetDb()
	err := db.Model(&model.UserRef{}).Where("uw_id = ?", userId).First(&r).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			res.Code = codes.CODE_ERR_UNKNOWN
			res.Msg = "unknown " + err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
	}

	var count int64
	if r.ID == 0 {
		count = 0
	} else {
		db.Model(&model.UserWallet{}).Where("ref_id = ?", r.ID).Count(&count)
	}
	res.Data = struct {
		Count int `json:"count"`
	}{
		Count: int(count),
	}
	c.JSON(http.StatusOK, res)
}

type RefResponse struct {
	MaskAddress string    `json:"mask_address"`
	Time        time.Time `json:"time"`
	TxFlag      int       `json:"tx_flag"`
}

func RefList(c *gin.Context) {
	var req request.Page

	res := common.Response{}
	res.Timestamp = time.Now().Unix()
	res.Code = codes.CODE_SUCCESS
	res.Msg = "success"
	res.Data = nil

	if err := c.ShouldBindJSON(&req); err != nil {
		res.Code = codes.CODE_ERR_REQFORMAT
		res.Msg = "invalid request" + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	_, _ = c.Get("user_wallet")
	userId, ok := c.Get("user_id")

	if !ok {
		res.Code = codes.CODE_ERR_SECURITY
		res.Msg = "need login"
		c.JSON(http.StatusOK, res)
		return
	}

	var r model.UserRef
	db := system.GetDb()
	err := db.Model(&model.UserRef{}).Where("uw_id = ?", userId).First(&r).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			res.Code = codes.CODE_ERR_UNKNOWN
			res.Msg = "unknown " + err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
	}

	refResult := make([]RefResponse, 0)
	if r.ID > 0 {
		var wallets []model.UserWallet
		db.Model(&model.UserWallet{}).
			Where("ref_id = ?", r.ID).Order("create_time desc").
			Limit(req.Ps).Offset((req.Pn - 1) * req.Ps).Find(&wallets)

		var ids []uint64

		for _, v := range wallets {
			ids = append(ids, v.ID)
		}
		var txs []model.TxsDiagram
		if len(ids) > 0 {
			db.Model(&model.TxsDiagram{}).Where("uw_id in ?", ids).Distinct("uw_id").Find(&txs)
		}
		for _, v := range wallets {
			trade := 0
			for _, w := range txs {
				if v.ID == w.UserID {
					trade = 1
					break
				}
			}
			refResult = append(refResult, RefResponse{
				MaskAddress: formatWalletAddress(v.Wallet),
				Time:        v.CreateTime,
				TxFlag:      trade,
			})
		}
	}
	res.Data = refResult
	c.JSON(http.StatusOK, res)
}

func formatWalletAddress(address string) string {
	if len(address) <= 12 {
		return address
	}
	return address[:6] + "******" + address[len(address)-6:]
}
