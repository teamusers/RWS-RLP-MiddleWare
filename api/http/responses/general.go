package responses

import (
	"fmt"
	"lbe/codes"
)

// APIResponse is the standard envelope for successful operations.
// The Data field contains the payload, which varies by endpoint.
type ApiResponse[T any] struct {
	Code int64 `json:"code"`
	// Message provides a human‑readable status or result description.
	// Example: "user created", "email found"
	Message string `json:"message"`
	// Data holds the response payload. Its type depends on the endpoint:
	// e.g. AuthResponse for /auth, LoginResponse for /user/login, etc.
	Data T `json:"data"`
}

func DefaultResponse(code int64, message string) ApiResponse[any] {
	return ApiResponse[any]{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

func InternalErrorResponse() ApiResponse[any] {
	return DefaultResponse(codes.INTERNAL_ERROR, "internal error")
}

func InvalidRequestBodyErrorResponse() ApiResponse[any] {
	return DefaultResponse(codes.INVALID_REQUEST_BODY, "invalid json request body")
}

func InvalidRequestBodySpecificErrorResponse(errString string) ApiResponse[any] {
	return DefaultResponse(codes.INVALID_REQUEST_BODY, fmt.Sprintf("invalid json request body:%s", errString))
}

func InvalidQueryParametersErrorResponse() ApiResponse[any] {
	return DefaultResponse(codes.INVALID_QUERY_PARAMETERS, "invalid query parameters")
}

func MissingAppIdErrorResponse() ApiResponse[any] {
	return DefaultResponse(codes.MISSING_APP_ID, "missing appId header")
}

func InvalidAppIdErrorResponse() ApiResponse[any] {
	return DefaultResponse(codes.INVALID_APP_ID, "invalid appId header")
}

func MissingAuthTokenErrorResponse() ApiResponse[any] {
	return DefaultResponse(codes.MISSING_AUTH_TOKEN, "missing authorization token")
}

func InvalidAuthTokenErrorResponse() ApiResponse[any] {
	return DefaultResponse(codes.INVALID_AUTH_TOKEN, "invalid authorization token")
}

func InvalidSignatureErrorResponse() ApiResponse[any] {
	return DefaultResponse(codes.INVALID_SIGNATURE, "invalid signature")
}

func ExistingUserFoundErrorResponse() ApiResponse[any] {
	return DefaultResponse(codes.EXISTING_USER_FOUND, "existing user found")
}

func ExistingUserNotFoundErrorResponse() ApiResponse[any] {
	return DefaultResponse(codes.EXISTING_USER_NOT_FOUND, "existing user not found")
}

func GrMemberIdLinkedErrorResponse() ApiResponse[any] {
	return DefaultResponse(codes.GR_MEMBER_LINKED, "gr profile already linked to another email")
}

func InvalidGrMemberClassErrorResponse() ApiResponse[any] {
	return DefaultResponse(codes.INVALID_GR_MEMBER_CLASS, "invalid gr member class provided")
}

func CachedProfileNotFoundErrorResponse() ApiResponse[any] {
	return DefaultResponse(codes.CACHED_PROFILE_NOT_FOUND, "cached profile not found")
}
