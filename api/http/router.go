package http

import (
	v1 "lbe/api/http/controllers/v1"

	user "lbe/api/http/controllers/v1/user"
	"lbe/api/interceptor"

	"github.com/gin-gonic/gin"
)

func Routers(e *gin.RouterGroup) {

	v1Group := e.Group("/v1")
	v1Group.POST("/auth", v1.AuthHandler)
	usersGroup := v1Group.Group("/user", interceptor.HttpInterceptor())
	{
		// The endpoints below will all require a valid access token.
		usersGroup.POST("/login", user.Login)
		usersGroup.POST("/register/verify", user.VerifyUserExistence)
		usersGroup.POST("/register", user.CreateUser)

		usersGroup.GET("/member/:external_id", user.GetMemberProfile)
		usersGroup.PUT("/member/update/:external_id", user.UpdateMemberProfile)
		//archive not ready yet for RLP - SessionM API
		//PUT - LBE-11 - api/v1/member/archive - withdraw member profile (active_status=0, previous email=current email, email=null)

		//PUT - LBE-5 - api/v1/user/pin - burn PIN update
		usersGroup.PUT("/pin", user.UpdateBurnPin)
		//POST - LBE-6 - api/v1/user/gr - GR user's profile verification
		usersGroup.POST("/gr", user.VerifyGrExistence)
		//POST - LBE-7 - api/v1/user/gr-cms - GR user's profile pushed by CMS
		usersGroup.POST("/gr-cms", user.VerifyGrCmsExistence)
		//GET - LBE-8 - api/v1/user/gr-reg - verify GR user's profile pushed by CMS

		//View Transaction - Timeline APIs
		//View Store Transaction - SM.TransactionsDomain.API  - post /api/1.0/transactions/info/get_store_transactions
		//API - create/update Transaction - post /api/1.0/transactions/info/get_transaction - Get the Transaction ID and generate a new one in LBE

	}

}
