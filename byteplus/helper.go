package byteplus

import (
	byteplusCdnClient "github.com/byteplus-sdk/byteplus-sdk-golang/service/cdn"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func initNewClient(providerConfig *byteplusCdnClient.CDN, planConfig *clientConfig) (initClient bool, diag diag.Diagnostics) {
	initClient = false
	return
}
