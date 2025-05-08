package v1

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func contentHash(body []byte) string {
	h := sha256.Sum256(body)
	return base64.StdEncoding.EncodeToString(h[:])
}

func computeSignature(stringToSign, base64Key string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, keyBytes)
	mac.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

func Email(c *gin.Context) {
	// ————————————————————————————————————
	// 1) Parse the connection string
	conn := "endpoint=https://lbe-acs-service.asiapacific.communication.azure.com/;accesskey=EycfZ5VwayONYQoEkffje6YImDC8pNvq7UhdHY0I0mSKLOSM2WwkJQQJ99BEACULyCpE2V5NAAAAAZCSTyPv"
	var (
		endpoint, key string
	)
	for _, kv := range strings.Split(conn, ";") {
		if strings.HasPrefix(kv, "endpoint=") {
			endpoint = strings.TrimSuffix(strings.TrimPrefix(kv, "endpoint="), "/")
		}
		if strings.HasPrefix(kv, "accesskey=") {
			key = strings.TrimPrefix(kv, "accesskey=")
		}
	}

	// ————————————————————————————————————
	// 2) Build the URL parts
	u, err := url.Parse(endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid endpoint URL"})
		return
	}
	host := u.Host
	pathAndQuery := "/emails:send?api-version=2023-03-31"

	// ————————————————————————————————————
	// 3) Build the JSON body
	bodyObj := map[string]interface{}{
		"senderAddress": "DoNotReply@0efe76d7-4c16-4d97-9fd1-adc32956d527.azurecomm.net",
		"content": map[string]string{
			"subject":   "Test HMAC",
			"plainText": "Hello via HMAC!",
			"html":      "<p>Hello via <strong>HMAC</strong>!</p>",
		},
		"recipients": map[string][]map[string]string{
			"to": {{"address": "perrylzx@gmail.com"}},
		},
	}
	bodyBytes, _ := json.Marshal(bodyObj)

	// ————————————————————————————————————
	// 4) Compute date and body hash
	gmt := time.FixedZone("GMT", 0)
	date := time.Now().In(gmt).Format(time.RFC1123)
	hash := contentHash(bodyBytes)

	// ————————————————————————————————————
	// 5) Build the string-to-sign
	stringToSign := strings.Join([]string{
		"POST",
		pathAndQuery,
		date + ";" + host + ";" + hash,
	}, "\n")

	// ————————————————————————————————————
	// 6) Sign it
	sig, err := computeSignature(stringToSign, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	authHeader := "HMAC-SHA256 SignedHeaders=x-ms-date;host;x-ms-content-sha256&Signature=" + sig

	// ————————————————————————————————————
	// 7) Create and send the outbound request
	req, err := http.NewRequestWithContext(context.Background(),
		"POST",
		endpoint+pathAndQuery,
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// override Host header so signing works
	req.Host = host
	req.Header.Set("x-ms-date", date)
	req.Header.Set("x-ms-content-sha256", hash)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.JSON(resp.StatusCode, gin.H{
		"status": resp.Status,
		"body":   string(body),
	})
}
