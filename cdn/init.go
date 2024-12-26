package cdn

import (
	"github.com/byteplus-sdk/byteplus-sdk-golang/service/cdn"
)

var (
	DefaultInstance = cdn.DefaultInstance
	ak              = "AKAPZjJjZWQwNzdkNTg1NGQwNzgyYTdhNzM4MjRiM2RmMzc"
	sk              = ""
)

func init() {
	DefaultInstance.Client.SetAccessKey(ak)
	DefaultInstance.Client.SetSecretKey(sk)
}
