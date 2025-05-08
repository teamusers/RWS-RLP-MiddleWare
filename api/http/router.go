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
	v1Group.POST("/email", v1.Email)
	v1Group.POST("/verify", v1.VerifyUserHandler)
	v1Group.POST("/register", v1.RegisterUserHandler)

	usersGroup := v1Group.Group("/user", interceptor.HttpInterceptor())
	{
		// The endpoints below will all require a valid access token.
		//POST - LBE-2 - api/v1/user/login - user login
		usersGroup.POST("/login", user.Login)
		//POST - LBE-3 - api/v1/user/register/verify - verify if user email is new and unregistered
		usersGroup.POST("/register/verify", user.VerifyUserExistence)
		//POST - LBE-4 - api/v1/user/register - register user based on provided fields
		usersGroup.POST("/register", user.CreateUser)

		//GET - LBE-9 - api/v1/user/:external_id - get user profile from rlp
		usersGroup.GET("/:external_id", user.GetUserProfile)
		usersGroup.PUT("/pin", user.UpdateBurnPin)
		usersGroup.PUT("/update/:external_id", user.UpdateUserProfile)

		//archive not ready yet for RLP - SessionM API
		//PUT - LBE-11 - api/v1/user/archive - withdraw user profile (active_status=0, previous email=current email, email=null)

		//POST - LBE-6 - api/v1/user/gr - GR user's profile verification
		usersGroup.POST("/gr", user.VerifyGrExistence)
		//POST - LBE-7 - api/v1/user/gr-cms - GR user's profile pushed by CMS
		usersGroup.POST("/gr-cms", user.VerifyGrCmsExistence)
		//GET - LBE-8 - api/v1/user/gr-reg - verify GR user's profile pushed by CMS
		usersGroup.GET("/gr-reg/:reg_id", user.GetCachedGrCmsProfile)
	}

}
