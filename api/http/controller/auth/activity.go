package auth

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stonksdex/externalapi/api/common"
	"github.com/stonksdex/externalapi/codes"
	"github.com/stonksdex/externalapi/model"
	"github.com/stonksdex/externalapi/system"
)

var signInLocks sync.Map

func DailyCheckin(c *gin.Context) {
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

	_, loaded := signInLocks.LoadOrStore(userID, struct{}{})
	defer func() {
		signInLocks.Delete(userID)
	}()
	if loaded {
		res.Code = codes.CODE_ERR_PROCESSING
		res.Msg = "checking, please do not repeat the operation"
		c.JSON(http.StatusOK, res)
		return
	}
	currentDate := time.Now().Format("20060102")
	var checkinObj model.DailyCheck
	db := system.GetDb()
	db.Model(&model.DailyCheck{}).Where("check_date = ? and uw_id = ?", currentDate, userID).First(&checkinObj)

	if checkinObj.ID > 0 {
		res.Code = codes.CODE_ERR_REPEAT
		res.Msg = "already finished daily check-in"
		c.JSON(http.StatusOK, res)
		return
	}
	checkinObj = model.DailyCheck{
		UwID:      userID,
		CheckDate: currentDate,
		CheckTime: time.Now(),
	}
	err = db.Save(&checkinObj).Error
	if err != nil {
		res.Code = codes.CODE_ERR_UNKNOWN
		res.Msg = "checkin record save error " + err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	res.Code = codes.CODE_SUCCESS
	res.Msg = "success"
	c.JSON(http.StatusOK, res)
}

func DailyCheckinRecord(c *gin.Context) {
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
	var result []model.DailyCheck

	db.Model(&model.DailyCheck{}).Where("uw_id = ?", userID).Order("check_date desc").Find(&result)

	res.Code = codes.CODE_SUCCESS
	res.Msg = "success"
	res.Data = result

	c.JSON(http.StatusOK, res)
}
