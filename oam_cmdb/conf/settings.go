package conf

import (
	"crypto/rsa"
)

var GlobalCfg settings

const STATIC_EXT_BASE_URL = "/static_ext"
const INSTALL_FILE_NAME = ".install"

type settings struct {
	//WORK_PATH               string
	DATA_STORE_PATH         string
	CACHE_TYPE              string
	RSA_DEFAULT_PRIVATE_KEY *rsa.PrivateKey
	RSA_DEFAULT_PUBLIC_KEY  *rsa.PublicKey
	SYMMETRIC_KEY           string
	IS_INSTALLED            bool
}
