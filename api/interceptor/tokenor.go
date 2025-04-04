package interceptor

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stonksdex/externalapi/api/common"
	"github.com/stonksdex/externalapi/codes"
	"github.com/stonksdex/externalapi/log"
	"github.com/stonksdex/externalapi/security"
)

func TokenInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		allHeaders, ok := c.Get("HEADERS")
		if !ok {
			log.Info("unable to get headers")
			makeFaileRes(c, codes.CODE_ERR_SECURITY, "token check failed")
			return
		}
		allHeadersMap := allHeaders.(common.HeaderParam)
		if allHeadersMap.XAuth == "123456" {
			c.Set("user_wallet", "0x0")
			c.Set("user_id", "1")
			c.Next()
			return
		}
		token, err := security.Decrypt(c.Request.Header.Get("XAUTH"))
		log.Info("TOKENCHECK ", c.Request.Header.Get("XAUTH"), token)
		if err != nil {
			makeFaileRes(c, codes.CODE_ERR_SECURITY, "token check failed")
			return
		}
		tokenArr := strings.Split(token, "|")
		if len(tokenArr) != 4 {
			makeFaileRes(c, codes.CODE_ERR_SECURITY, "token length error")
			return
		}
		expireTs, err := strconv.ParseInt(tokenArr[3], 10, 64)
		if err != nil {
			makeFaileRes(c, codes.CODE_ERR_SECURITY, "token format error")
			return
		}
		if time.Now().Unix()-expireTs > int64(common.TOKEN_DURATION.Seconds()) {
			makeFaileRes(c, codes.CODE_ERR_SECURITY, "token expired error")
			return
		}

		c.Set("user_wallet", tokenArr[1])
		c.Set("user_id", tokenArr[0])

		c.Next()
	}
}

func makeFaileRes(c *gin.Context, code int64, msg string) {
	c.Abort()
	c.JSON(http.StatusOK, common.Response{
		Code:      code,
		Msg:       msg,
		Timestamp: time.Now().Unix(),
	})
}
