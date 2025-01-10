package byteplus

const (
	ERR_CODE_NOT_FOUND_DOMAIN                                = "NotFound.Domain"
	ERR_CODE_IAM_UNAUTHORIZED                                = "AccessDenied.IAMUnauthorized"
	ERR_CODE_QUOTA_EXCEEDED_TODAY                            = "QuotaExceeded.UrlsToday"
	ERR_CODE_INVALID_PARAMETER                               = "InvalidParameter."
	ERR_CODE_INVALID_PARAMETER_URLS                          = "InvalidParameter.Urls"
	ERR_CODE_INVALID_PARAMETER_CERTIFICATE                   = "InvalidParameter.Certificate"
	ERR_CODE_INVALID_PARAMETER_CERTIFICATE_KEY_NOT_MATCH     = "InvalidParameter.Certificate.KeyNotMatch"
	ERR_CODE_INVALID_PARAMETER_HTTPS_CERT_INFO_CHAIN_MISSING = "InvalidParameter.Https.CertInfo.ChainMissing"
	ERR_CODE_INVALID_PARAMETER_BILLING_CODE                  = "InvalidParameter.BillingCode"
	ERR_CODE_INVALID_PARAMETER_BILLING_REGION                = "InvalidParameter.BillingRegion"
	ERR_CODE_SERVICE_STOPPED                                 = "OperationDenied.ServiceStopped"

	ERR_CLOSE_DNS_SLB_FAILED  = "CloseDnsSlbFailed"
	ERR_DISABLE_DNS_SLB       = "DisableDNSSLB"
	ERR_ENABLE_DNS_SLB_FAILED = "EnableDnsSlbFailed"
	ERR_DNS_SYSTEM_BUSYNESS   = "DnsSystemBusyness"
	ERR_SERVICE_UNAVAILABLE   = "ServiceUnavailable"
	ERR_THROTTLING_USER       = "Throttling.User"
	ERR_THROTTLING_API        = "Throttling.API"
	ERR_THROTTLING            = "Throttling"
	ERR_UNKNOWN_ERROR         = "UnknownError"
	ERR_INTERNAL_ERROR        = "InternalError"
	ERR_BACKEND_TIMEOUT       = "D504TO"
)

func isPermanentCdnError(errCode string) bool {
	switch errCode {
	case
		ERR_CODE_NOT_FOUND_DOMAIN,
		ERR_CODE_IAM_UNAUTHORIZED,
		ERR_CODE_SERVICE_STOPPED,
		ERR_CODE_QUOTA_EXCEEDED_TODAY:
		return true
	default:
		return false
	}
}

func isAbleToRetry(errCode string) bool {
	switch errCode {
	case ERR_CLOSE_DNS_SLB_FAILED,
		ERR_DISABLE_DNS_SLB,
		ERR_ENABLE_DNS_SLB_FAILED,
		ERR_DNS_SYSTEM_BUSYNESS,
		ERR_SERVICE_UNAVAILABLE,
		ERR_THROTTLING_USER,
		ERR_THROTTLING_API,
		ERR_THROTTLING,
		ERR_UNKNOWN_ERROR,
		ERR_INTERNAL_ERROR:
		return true
	default:
		return false
	}
	// return false
}
