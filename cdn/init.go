package cdn

import (
	"github.com/byteplus-sdk/byteplus-sdk-golang/service/cdn"
)

var (
	DefaultInstance = cdn.DefaultInstance
	ak              = ""
	sk              = ""
	domainName      = "www.example.com"
)

func init() {
	DefaultInstance.Client.SetAccessKey(ak)
	DefaultInstance.Client.SetSecretKey(sk)
}
