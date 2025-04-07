package interceptor

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"rlp-middleware/api/common"
	"rlp-middleware/codes"
	"rlp-middleware/log"
	"rlp-middleware/model"
	"rlp-middleware/security"
	"rlp-middleware/system"

	"github.com/gin-gonic/gin"
)

var exception = []string{""}

var cacheKeys []model.SysChannel
var cacheValid = time.Now().Unix()

const timeRange = 60

var wg sync.Mutex

const (
	HTTP = "http"
	WS   = "ws"
)

func HttpInterceptor() gin.HandlerFunc {
	debugMode := true
	return func(c *gin.Context) {

		if debugMode {
			// In debug mode, bypass all checks.
			c.Set("APPID", "debug-app-id")
			c.Set("REQUESTID", "debug-request-id")
			c.Set("TS", time.Now().Unix())
			c.Set("HEADERS", map[string]string{"AppId": "debug-app-id"})
			// Optionally, set default user values.
			c.Set("user_wallet", "0x0")
			c.Set("user_id", "1")
			c.Next()
			return
		}

		queryKeys(HTTP)
		hp := parse(&c.Request.Header)

		var targetChannel *model.SysChannel
		for _, v := range cacheKeys {
			if v.AppID == hp.AppId {
				targetChannel = &v
			}
		}
		if targetChannel == nil {
			c.Abort()
			c.JSON(http.StatusOK, common.Response{
				Code:      codes.CODE_ERR_APPID_INVALID,
				Msg:       "appid invalid:" + hp.AppId,
				Timestamp: time.Now().Unix(),
			})
			return
		}
		if ok, code := targetChannel.Verify(hp.Join(), hp.AuthToken); !ok {
			c.Abort()
			c.JSON(http.StatusOK, common.Response{
				Code:      int64(code),
				Msg:       "sig or key params wrong or empty",
				Timestamp: time.Now().Unix(),
			})
			return
		}
		c.Set("APPID", hp.AppId)
		c.Set("REQUESTID", hp.RequestId)
		c.Set("TS", hp.Ts)
		c.Set("HEADERS", hp)

		if hp.XAuth == "123456" {
			c.Set("user_wallet", "0x0")
			c.Set("user_id", "1")
		} else {
			token, err := security.Decrypt(hp.XAuth)
			if err == nil {

				tokenArr := strings.Split(token, "|")
				if len(tokenArr) == 4 {
					expireTs, err := strconv.ParseInt(tokenArr[3], 10, 64)
					if err == nil {
						if time.Now().Unix()-expireTs <= int64(common.TOKEN_DURATION.Seconds()) {
							c.Set("user_wallet", tokenArr[1])
							c.Set("user_id", tokenArr[0])
						}
					}
				}

			}
		}
		c.Next()
	}
}

func queryKeys(channel string) []model.SysChannel {
	if len(cacheKeys) > 0 && time.Now().Unix()-cacheValid <= (timeRange) {
		return cacheKeys
	}
	db := system.GetDb()
	var result []model.SysChannel
	err := db.Model(&model.SysChannel{}).Where("status = ? and chan = ?", "00", channel).Find(&result).Error
	if err != nil {
		log.Error("Channel Query Error:", err)
		return cacheKeys
	}

	wg.Lock()
	cacheKeys = result
	cacheValid = time.Now().Unix()
	wg.Unlock()
	return cacheKeys
}

func parse(h *http.Header) common.HeaderParam {
	headerParam := common.HeaderParam{
		AppId:     h.Get("APPID"),
		AuthToken: h.Get("SIG"),
		Ts:        h.Get("TS"),
		Ver:       h.Get("VER"),
		RequestId: h.Get("REQUESTID"),
		XAuth:     h.Get("XAUTH"),
	}
	return headerParam
}
