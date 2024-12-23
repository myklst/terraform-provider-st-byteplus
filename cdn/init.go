package cdn

import (
	"github.com/byteplus-sdk/byteplus-sdk-golang/service/cdn"
)

var (
	DefaultInstance = cdn.DefaultInstance
	ak              = "AKAPMWMzODdjNGVhZjhmNDYyN2FhYWIyM2RjNDdjMDBiODE"
	sk              = "TXpaa05HUmlPV0k1T0dNeU5EUTRaR0l5TjJZeFlUVm1NMlF6T0dZd1pqaw=="
)

func init() {
	DefaultInstance.Client.SetAccessKey(ak)
	DefaultInstance.Client.SetSecretKey(sk)
}
