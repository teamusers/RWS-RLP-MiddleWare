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
	memberGroup := v1Group.Group("/member", interceptor.HttpInterceptor())
	{

		memberGroup.GET("/:external_id", v1.GetMemberProfile)
		memberGroup.PUT("/pin", v1.UpdateBurnPin)
		memberGroup.PUT("/update/:external_id", v1.UpdateMemberProfile)

		//archive not ready yet for RLP - SessionM API
		//PUT - LBE-11 - api/v1/member/archive - withdraw member profile (active_status=0, previous email=current email, email=null)
	}
	usersGroup := v1Group.Group("/user", interceptor.HttpInterceptor())
	{
		// The endpoints below will all require a valid access token.
		usersGroup.POST("/login", user.Login)
		usersGroup.POST("/register/verify", user.VerifyUserExistence)
		usersGroup.POST("/register", user.CreateUser)

		//POST - LBE-6 - api/v1/user/gr - GR user's profile verification
		usersGroup.POST("/gr", user.VerifyGrExistence)
		//POST - LBE-7 - api/v1/user/gr-cms - GR user's profile pushed by CMS
		usersGroup.POST("/gr-cms", user.VerifyGrCmsExistence)
		//GET - LBE-8 - api/v1/user/gr-reg - verify GR user's profile pushed by CMS
		usersGroup.GET("/gr-reg/:reg_id", user.GetCachedGrCmsProfile)
	}

}
