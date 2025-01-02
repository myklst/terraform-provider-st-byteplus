package byteplus

const (
	ERR_CODE_MISSING_PARAMETER            = "MissingParameter"
	ERR_CODE_MISSING_AUTHENTICATION_TOKEN = "MissingAuthenticationToken"
	ERR_CODE_MISSING_REQUEST_INFO         = "MissingRequestInfo"
	ERR_CODE_MISSING_SIGNATURE            = "MissingSignature"
	ERR_CODE_INVALID_TIMESTAMP            = "InvalidTimestamp"
	ERR_CODE_SERVICE_NOT_FOUND            = "ServiceNotFound"
	ERR_CODE_INVALID_ACTION_OR_VERSION    = "InvalidActionOrVersion"
	ERR_CODE_INVALID_ACCESS_KEY           = "InvalidAccessKey"
	ERR_CODE_SIGNATURE_DOES_NOT_MATCH     = "SignatureDoesNotMatch"
	ERR_CODE_METHOD_NOT_ALLOWED           = "MethodNotAllowed"
	ERR_CODE_INVALID_AUTHORIZATION        = "InvalidAuthorization"
	ERR_CODE_INVALID_CREDENTIAL           = "InvalidCredential"
	ERR_CODE_UNDEFINED_ERROR              = "UndefinedError"
	ERR_CODE_INTERNAL_ERROR               = "InternalError"
	ERR_CODE_INTERNAL_SERVICE_ERROR       = "InternalServiceError"
	ERR_CODE_FAIL_TO_CONNECT              = "FailToConnect"
	ERR_CODE_INTERNAL_SERVICE_TIMEOUT     = "InternalServiceTimeout"
	ERR_CODE_SERVICE_UNAVAILABLE_TEMP     = "ServiceUnavailableTemp"
)

func isPermanentCommonError(errCode string) bool {
	switch errCode {
	case
		ERR_CODE_MISSING_PARAMETER,
		ERR_CODE_MISSING_AUTHENTICATION_TOKEN,
		ERR_CODE_MISSING_REQUEST_INFO,
		ERR_CODE_MISSING_SIGNATURE,
		ERR_CODE_INVALID_TIMESTAMP,
		ERR_CODE_SERVICE_NOT_FOUND,
		ERR_CODE_INVALID_ACTION_OR_VERSION,
		ERR_CODE_INVALID_ACCESS_KEY,
		ERR_CODE_SIGNATURE_DOES_NOT_MATCH,
		ERR_CODE_METHOD_NOT_ALLOWED,
		ERR_CODE_INVALID_AUTHORIZATION,
		ERR_CODE_INVALID_CREDENTIAL,
		ERR_CODE_UNDEFINED_ERROR,
		ERR_CODE_INTERNAL_ERROR,
		ERR_CODE_INTERNAL_SERVICE_ERROR:
		return true
	default:
		return false
	}
}
