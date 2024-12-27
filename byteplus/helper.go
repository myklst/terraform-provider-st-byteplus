package byteplus

import (
	"github.com/alibabacloud-go/tea/tea"
	byteplusBaseClient "github.com/byteplus-sdk/byteplus-sdk-golang/base"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func initNewClient(providerConfig *byteplusBaseClient.Credentials, planConfig *clientConfig) (initClient bool, clientConfig *byteplusBaseClient.Credentials, diag diag.Diagnostics) {
	initClient = false
	clientConfig = &byteplusBaseClient.Credentials{}
	region := planConfig.Region.ValueString()
	accessKey := planConfig.AccessKey.ValueString()
	secretKey := planConfig.SecretKey.ValueString()

	if region != "" || accessKey != "" || secretKey != "" {
		initClient = true
	}

	if initClient {
		if region == "" {
			region = tea.StringValue(&providerConfig.Region)
		}
		if accessKey == "" {
			clientAccessKey := providerConfig.AccessKeyID
			if clientAccessKey == "" {
				diag.AddError(
					"Failed to retrieve client Access Key.",
					"This is an error in provider, please contact the provider developers.",
				)
			} else {
				accessKey = clientAccessKey
			}
		}
		if secretKey == "" {
			clientSecretKey := providerConfig.SecretAccessKey
			if clientSecretKey == "" {
				diag.AddError(
					"Failed to retrieve client Secret Key.",
					"This is an error in provider, please contact the provider developers.",
				)
			} else {
				secretKey = clientSecretKey
			}
		}
		if diag.HasError() {
			return
		}

		clientConfig = &byteplusBaseClient.Credentials{
			AccessKeyID:     accessKey,
			SecretAccessKey: secretKey,
			Region:          region,
		}
	}

	return
}
