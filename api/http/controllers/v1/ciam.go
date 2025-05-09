package v1

import (
	"log"
	"net/http"

	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/config"

	"lbe/api/http/services"

	"github.com/gin-gonic/gin"
)

// VerifyUserHandler handles POST /ciam/users/verify
func VerifyUserHandler(c *gin.Context) {
	var req requests.VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// coll, err := services.VerifyCIAMExistenceByEmail(c.Request.Context(), req.Email)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// if len(coll.Value) == 0 {
	// 	c.JSON(http.StatusNotFound, gin.H{"message": "no user found"})
	// 	return
	// }

	c.JSON(http.StatusOK, nil)
}

// RegisterUserHandler handles POST /ciam/users/register
func RegisterUserHandler(c *gin.Context) {
	var req requests.GraphCreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("verify if user exists before creation")
	if respData, err := services.GetCIAMUserByEmail(c.Request.Context(), req.Mail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if len(respData.Value) != 0 {
		log.Println("user found")
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	log.Println("starting register")
	var userId = ""
	if respData, err := services.PostCIAMRegisterUser(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Println("user id: " + respData.Id)
		userId = respData.Id
	}

	log.Println("verify if user exists after creation")
	if respData, err := services.GetCIAMUserByEmail(c.Request.Context(), req.Mail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if len(respData.Value) == 0 {
		log.Println("user not found")
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	log.Println("register schema")
	schemaExtensionsPayload := map[string]any{
		config.GetConfig().Api.Eeid.UserIdLinkExtensionKey: requests.UserIdLinkSchemaExtensionFields{
			RlpId: "1101",
			RlpNo: "1100",
			GrId:  "0000",
		},
	}

	log.Println("schema" + userId)
	log.Println(schemaExtensionsPayload)

	if err := services.PatchCIAMAddUserSchemaExtensions(c, userId, schemaExtensionsPayload); err != nil {
		log.Printf("CIAM Patch User Schema Extensions failed: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	log.Println("get user by grId - success")
	if respData, err := services.GetCIAMUserByGrId(c, "0000"); err != nil {
		log.Printf("error encountered verifying user existence: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	} else if len(respData.Value) == 0 {
		log.Println("user not found")
		c.JSON(http.StatusConflict, nil)
		return
	}

	log.Println("get user by grId - fail")
	if respData, err := services.GetCIAMUserByGrId(c, "00100"); err != nil {
		log.Printf("error encountered verifying user existence: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	} else if len(respData.Value) != 0 {
		log.Println("user found")
		c.JSON(http.StatusConflict, nil)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}
