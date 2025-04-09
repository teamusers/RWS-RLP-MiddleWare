package services

import (
	"context"
	"fmt"
	"math/rand"
	"rlp-member-service/system"
	"strconv"
	"time"
)

// OTPResponse contains the OTP code and its expiration timestamp.
type OTPResponse struct {
	OTP       string `json:"otp_code"`
	ExpiresAt int64  `json:"expires_at"` // Unix timestamp
}

// OTPService is responsible for generating and validating OTP tokens.
type OTPService interface {
	GenerateOTP(ctx context.Context, identifier string) (OTPResponse, error)
	ValidateOTP(ctx context.Context, identifier string, otp string) (bool, error)
}

type otpService struct{}

// NewOTPService creates an instance of OTPService.
func NewOTPService() OTPService {
	return &otpService{}
}

// GenerateOTP generates a 6-digit OTP, stores it in Redis with a 30-minute expiration,
// and returns the OTP along with its expiration time.
func (s *otpService) GenerateOTP(ctx context.Context, identifier string) (OTPResponse, error) {
	// Seed the random number generator (consider seeding once in your application's startup in production)
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a 6-digit OTP.
	otp := rand.Intn(900000) + 100000
	otpStr := strconv.Itoa(otp)

	// Define the Redis key (e.g., "otp:user@example.com").
	key := "otp:" + identifier

	// Set an expiration duration, e.g., 30 minutes.
	expiration := 30 * time.Minute

	// Store the OTP in Redis.
	err := system.GetRedis().Set(ctx, key, otpStr, expiration).Err()
	if err != nil {
		return OTPResponse{}, fmt.Errorf("failed to store OTP in Redis: %v", err)
	}

	// Calculate the expiration timestamp.
	expiresAt := time.Now().Add(expiration).Unix()

	// Return the OTP and the expiration time.
	return OTPResponse{
		OTP:       otpStr,
		ExpiresAt: expiresAt,
	}, nil
}

// ValidateOTP retrieves the stored OTP from Redis for the given identifier, compares it with the provided OTP,
// and optionally deletes it after successful validation.
func (s *otpService) ValidateOTP(ctx context.Context, identifier string, providedOTP string) (bool, error) {
	// Define the Redis key.
	key := "otp:" + identifier

	// Retrieve the stored OTP from Redis.
	storedOTP, err := system.GetRedis().Get(ctx, key).Result()
	if err != nil {
		if err == system.Nil {
			// OTP not found or expired.
			return false, fmt.Errorf("OTP not found or expired for identifier: %s", identifier)
		}
		return false, fmt.Errorf("failed to get OTP from Redis: %v", err)
	}

	// Compare the stored OTP with the provided OTP.
	if storedOTP != providedOTP {
		return false, nil
	}

	// Optionally delete the OTP from Redis after successful validation.
	_, delErr := system.GetRedis().Del(ctx, key).Result()
	if delErr != nil {
		// Log a warning if deletion fails.
		fmt.Printf("Warning: failed to delete OTP from Redis: %v\n", delErr)
	}

	return true, nil
}
