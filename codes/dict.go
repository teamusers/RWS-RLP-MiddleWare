package codes

const (
	// enums

	// signUpType enums
	SignUpTypeNew   = "NEW"
	SignUpTypeGR    = "GR"
	SignUpTypeGRCMS = "GR_CMS"
	SignUpTypeTM    = "TM"

	// codes

	CODE_SUCCESS              = 0
	CODE_ERR_METHOD_UNSUPPORT = 1
	CODE_ERR_REQFORMAT        = 2
	CODE_ERR_APPID_INVALID    = 3
	CODE_ERR_SIGMETHOD_UNSUPP = 4
	CODE_ERR_AUTHTOKEN_FAIL   = 5

	SUCCESSFUL   int64 = 1000
	UNSUCCESSFUL int64 = 1001
	FOUND        int64 = 1002
	NOT_FOUND    int64 = 1003

	INTERNAL_ERROR           int64 = 4000
	INVALID_REQUEST_BODY     int64 = 4001
	INVALID_AUTH_TOKEN       int64 = 4002
	MISSING_AUTH_TOKEN       int64 = 4003
	INVALID_SIGNATURE        int64 = 4004
	MISSING_SIGNATURE        int64 = 4005
	INVALID_APP_ID           int64 = 4006
	MISSING_APP_ID           int64 = 4007
	INVALID_QUERY_PARAMETERS int64 = 4008
	EXISTING_USER_NOT_FOUND  int64 = 4009
	EXISTING_USER_FOUND      int64 = 4010
	CACHED_PROFILE_NOT_FOUND int64 = 4011
	GR_MEMBER_LINKED         int64 = 4012
	GR_MEMBER_NOT_FOUND      int64 = 4013
	INVALID_GR_MEMBER_CLASS  int64 = 4014
)

func IsValidSignUpType(t string) bool {
	switch t {
	case SignUpTypeNew, SignUpTypeGRCMS, SignUpTypeGR, SignUpTypeTM:
		return true
	default:
		return false
	}
}
