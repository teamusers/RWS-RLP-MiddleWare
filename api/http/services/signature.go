package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"lbe/api/http/requests"
)

func GenerateSignature(appID, secretKey string) (*requests.AuthRequest, error) {
	// Initialize random seed.
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random nonce (16 characters).
	nonce := randomNonce(16)
	// Get the current Unix timestamp as a string.
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	return computeSignature(appID, nonce, timestamp, secretKey)
}

// GenerateSignatureWithParams computes the signature using the provided parameters.
// You can supply your own AppID, nonce, timestamp, and secretKey.
func GenerateSignatureWithParams(appID, nonce, timestamp, secretKey string) (*requests.AuthRequest, error) {
	return computeSignature(appID, nonce, timestamp, secretKey)
}

// computeSignature is a helper function which computes the HMAC-SHA256 signature from the given parameters.
func computeSignature(appID, nonce, timestamp, secretKey string) (*requests.AuthRequest, error) {
	// Build the base string by concatenating AppID, timestamp, and nonce.
	baseString := appID + timestamp + nonce

	// Compute HMAC-SHA256 with the provided secret key.
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err := mac.Write([]byte(baseString))
	if err != nil {
		return nil, fmt.Errorf("error computing MAC: %w", err)
	}
	expectedMAC := mac.Sum(nil)
	signature := hex.EncodeToString(expectedMAC)

	return &requests.AuthRequest{
		Nonce:     nonce,
		Timestamp: timestamp,
		Signature: signature,
	}, nil
}

// randomNonce returns a random alphanumeric string of length n.
func randomNonce(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
