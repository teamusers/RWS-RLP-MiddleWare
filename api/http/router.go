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
		//POST - LBE-2 - api/v1/user/login - user login To be removed
		// usersGroup.POST("/login", user.Login)

		//POST - LBE-3 - api/v1/user/register/verify - verify if user email is new and unregistered
		usersGroup.POST("/register/verify", user.VerifyUserExistence)
		//POST - LBE-4 - api/v1/user/register - register user based on provided fields
		usersGroup.POST("/register", user.CreateUser)
		// LBE-5 To be removed
		// usersGroup.PUT("/pin", user.UpdateBurnPin)
		//POST - LBE-6 - api/v1/user/gr - GR user's profile verification
		usersGroup.POST("/gr", user.VerifyGrExistence)
		//POST - LBE-7 - api/v1/user/gr-cms - GR user's profile pushed by CMS
		usersGroup.POST("/gr-cms", user.VerifyGrCmsExistence)
		//GET - LBE-8 - api/v1/user/gr-reg - verify GR user's profile pushed by CMS
		usersGroup.GET("/gr-reg/:reg_id", user.GetCachedGrCmsProfile)
		usersGroup.GET("/gr-reg", v1.InvalidQueryParametersHandler)

		//GET - LBE-9 - api/v1/user/:external_id - get user profile from rlp
		usersGroup.GET("/:external_id", user.GetUserProfile)
		usersGroup.GET("", v1.InvalidQueryParametersHandler)
		//PUT - LBE-10 - api/v1/user/update/:external_id - update user profile
		usersGroup.PUT("/update/:external_id", user.UpdateUserProfile)
		usersGroup.PUT("/update", v1.InvalidQueryParametersHandler)
		//PUT - LBE-11 - api/v1/user/archive - withdraw user profile (active_status=0, previous email=current email, email=null)
		usersGroup.PUT("/archive/:external_id", user.WithdrawUserProfile)
		usersGroup.PUT("/archive", v1.InvalidQueryParametersHandler)
	}

}
