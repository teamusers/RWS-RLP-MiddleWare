package interceptor

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"rlp-middleware/log"

	"github.com/gin-gonic/gin"
)

// GetSecretForAppID retrieves the secret key for the given AppID.
// In a real-world scenario, this might query a database or configuration store.
func GetSecretForAppID(appID string) (string, error) {
	// For demonstration purposes, we use a hard-coded map.
	// Replace this with your actual lookup logic.
	secretStore := map[string]string{
		"1": "mySecretKeyForApp1",
		// Add other AppIDs and secrets as needed.
	}

	secret, exists := secretStore[appID]
	if !exists {
		return "", ErrInvalidAppID
	}
	return secret, nil
}

// ErrInvalidAppID is returned when the provided AppID is not recognized.
var ErrInvalidAppID = &gin.Error{
	Err:  http.ErrNoCookie,
	Type: gin.ErrorTypePublic,
	Meta: "Invalid AppID",
}

// APITokenInterceptor validates the API token using HMAC-SHA256.
func APITokenInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve required headers.
		appID := c.Request.Header.Get("AppID")
		timestampStr := c.Request.Header.Get("Timestamp")
		nonce := c.Request.Header.Get("Nonce")
		clientSignature := c.Request.Header.Get("Signature")

		if appID == "" || timestampStr == "" || nonce == "" || clientSignature == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing required authentication headers",
			})
			c.Abort()
			return
		}

		// Optionally, validate that the timestamp is recent (e.g., within 5 minutes).
		timestampInt, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid timestamp format",
			})
			c.Abort()
			return
		}

		requestTime := time.Unix(timestampInt, 0)
		if time.Since(requestTime) > 5*time.Minute {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Request expired",
			})
			c.Abort()
			return
		}

		// Retrieve the secret for the given AppID.
		secret, err := GetSecretForAppID(appID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid AppID",
			})
			c.Abort()
			return
		}

		// Build the signature base string.
		// You can include additional fields if necessary.
		signatureBase := appID + timestampStr + nonce

		// Compute the HMAC-SHA256 signature.
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(signatureBase))
		expectedSignature := hex.EncodeToString(mac.Sum(nil))

		// Log details for debugging (ensure to remove in production or secure logging).
		log.Info("Expected Signature: ", expectedSignature)
		log.Info("Client Signature: ", clientSignature)

		// Compare the computed signature with the client's signature.
		if expectedSignature != clientSignature {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid signature",
			})
			c.Abort()
			return
		}

		// Optionally, set validated values in the context for use in handlers.
		c.Set("APPID", appID)
		c.Set("TIMESTAMP", timestampStr)
		c.Set("NONCE", nonce)

		c.Next()
	}
}
