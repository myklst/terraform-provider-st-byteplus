package cdn

import (
	"github.com/byteplus-sdk/byteplus-sdk-golang/service/cdn"
)

var (
	DefaultInstance = cdn.DefaultInstance
	ak              = "AKAPMWMzODdjNGVhZjhmNDYyN2FhYWIyM2RjNDdjMDBiODE"
	sk              = ""
)

func init() {
	DefaultInstance.Client.SetAccessKey(ak)
	DefaultInstance.Client.SetSecretKey(sk)
}
