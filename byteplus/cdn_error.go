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
