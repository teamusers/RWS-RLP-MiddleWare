package http

import (
	v1 "lbe/api/http/controllers/v1"
	user "lbe/api/http/controllers/v1/user"
	"lbe/api/interceptor"

	"github.com/gin-gonic/gin"
)

func Routers(e *gin.RouterGroup) {

	v1Group := e.Group("/v1")
	v1Group.GET("/auth", v1.AuthHandler)
	usersGroup := v1Group.Group("/user", interceptor.HttpInterceptor())
	{
		// The endpoints below will all require a valid access token.
		usersGroup.GET("/login/:email", user.Login)
		usersGroup.GET("/register/:email/:sign_up_type", user.GetUser)
		usersGroup.POST("/register", user.CreateUser)

		//GET - LBE-5 - /api/v1/user/pin - burn PIN update
		//GET - LBE-6 - api/v1/user/gr - GR user's profile verification
		//GET - LBE-7 - api/v1/user/gr-cms - GR user's profile pushed by CMS
		//GET - LBE-8 - api/v1/user/gr-reg - verify GR user's profile pushed by CMS
		//GET - LBE-9 - api/v1/member - view member profile
		//PUT - LBE-10 - api/v1/member/update - update member profile (name, phone, marketing consent, burn PIN)
		//PUT - LBE-11 - api/v1/member/archive - withdraw member profile (active_status=0, previous email=current email, email=null)

	}

}
