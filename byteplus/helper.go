package byteplus

import (
	byteplusBaseClient "github.com/byteplus-sdk/byteplus-sdk-golang/base"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func initNewClient(providerConfig *byteplusBaseClient.Client, planConfig *clientConfig) (initClient bool, clientConfig *byteplusBaseClient.Credentials, diag diag.Diagnostics) {
	initClient = false
	clientConfig = &byteplusBaseClient.Credentials{}
	region := planConfig.Region.ValueString()
	accessKey := planConfig.AccessKey.ValueString()
	secretKey := planConfig.SecretKey.ValueString()

	if region != "" || accessKey != "" || secretKey != "" {
		initClient = true
	}

	if initClient {

	}

	return
}
