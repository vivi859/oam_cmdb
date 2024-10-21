package util

import (
	conf "OAM/conf"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// RSA加密字符串,返回base64字符明文
func RSAEncryptBase64Str(src string, _publicKey string) (string, error) {
	cryptBytes, err := RSAEncrypt([]byte(src), _publicKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cryptBytes), nil
}

//RSA解密base64字符密文
func RSADecryptBase64Str(cryptBase64Str string, _privateKey string) (string, error) {
	cryptBytes, err := base64.StdEncoding.DecodeString(cryptBase64Str)
	if err != nil {
		return "", err
	}

	plainBytes, err := RSADecrypt(cryptBytes, _privateKey)
	if err != nil {
		return "", err
	}
	return string(plainBytes), nil
}

// RSA公钥加密.
// @param _publicKey 表示公钥文件路径或者bas64公钥串,如果是密钥文件必须是绝对路径
func RSAEncrypt(src []byte, _publicKey string) ([]byte, error) {
	var publicKey *rsa.PublicKey
	var err error
	if _publicKey == "" {
		publicKey = conf.GlobalCfg.RSA_DEFAULT_PUBLIC_KEY
	} else if isFile(_publicKey) {
		publicKey, err = LoadPublicKeyFromFile(_publicKey)
	} else {
		publicKey, err = LoadBase64PublicKey(_publicKey)
	}
	if err != nil {
		return nil, err
	}
	// 公钥加密
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, src)
}

func isFile(f string) bool {
	if filepath.IsAbs(f) {
		fi, err := os.Stat(f)
		if err != nil {
			return false
		}
		return !fi.IsDir()
	}
	return false
}

// RSA私钥解密.
// @param _privateKey: 表示公钥文件路径或者bas64公钥串,如果是密钥文件必须是绝对路径
func RSADecrypt(src []byte, _privateKey string) ([]byte, error) {
	var privateKey *rsa.PrivateKey
	var err error
	//私钥为空时使用系统默认私钥
	if _privateKey == "" {
		privateKey = conf.GlobalCfg.RSA_DEFAULT_PRIVATE_KEY
	} else if isFile(_privateKey) {
		privateKey, err = LoadPrivateKeyFromFile(_privateKey)
	} else {
		privateKey, err = LoadBase64PrivateKey(_privateKey)
	}

	if err != nil {
		return nil, err
	}

	// 私钥解密
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, src)
}

// 生成 RSA 公私钥
func RSAGenerateKeyPair(keyLen int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privatekey, err := rsa.GenerateKey(rand.Reader, keyLen)
	if err != nil {
		return nil, nil, err
	}

	publickey := &privatekey.PublicKey
	return privatekey, publickey, nil
}

// 生成RSA公私钥并保存到文件
func RSAGenerateKeyFile(keyLen int, privateKeyPath, publicKeyPath string) error {
	//1. 生成密钥
	priKey, pubKey, err := RSAGenerateKeyPair(keyLen)
	if err != nil {
		return err
	}
	// 2.私钥序列化为ASN.1 PKCS#1 DER编码
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(priKey)

	// 3. Block代表PEM编码的结构, 对其进行设置
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// 4. 创建文件
	privateFile, err := os.Create(privateKeyPath)
	if err == nil {
		defer privateFile.Close()
	} else {
		return err
	}

	// 5. 使用pem编码, 并将数据写入文件中
	err = pem.Encode(privateFile, &block)
	if err != nil {
		return err
	}

	// 生成公钥文件
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return err
	}

	block = pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	publicFile, err := os.Create(publicKeyPath)
	if err == nil {
		defer publicFile.Close()
	} else {
		return err
	}

	// 编码公钥 写入文件
	err = pem.Encode(publicFile, &block)
	if err != nil {
		return err
	}
	return nil
}

//生成rsa公私钥,结果以base64编码字符返回
//@return 公钥,私钥,error
func RSAGenerateKeyStr(keyLen int) (string, string, error) {
	priKey, pubKey, err := RSAGenerateKeyPair(keyLen)
	if err != nil {
		return "", "", err
	}
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(priKey)
	priKeybase64 := base64.StdEncoding.EncodeToString(privateKeyBytes)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", "", err
	}
	pubKeybase64 := base64.StdEncoding.EncodeToString(publicKeyBytes)
	return pubKeybase64, priKeybase64, nil
}

func PublicKeyToBase64(pubKey *rsa.PublicKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(publicKeyBytes), nil
}

func PrivateKeyToBase64(pubKey *rsa.PrivateKey) (string, error) {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(pubKey)
	return base64.StdEncoding.EncodeToString(privateKeyBytes), nil
}

// 加载私钥文件
func LoadPrivateKeyFromFile(keyfile string) (*rsa.PrivateKey, error) {
	keybuffer, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(keybuffer))
	if block == nil {
		return nil, errors.New("private key error")
	}

	privatekey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("parse private key error")
	}

	return privatekey, nil
}

//加载公钥文件
func LoadPublicKeyFromFile(keyfile string) (*rsa.PublicKey, error) {
	keybuffer, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keybuffer)
	if block == nil {
		return nil, errors.New("public key error")
	}

	pubkeyinterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publickey := pubkeyinterface.(*rsa.PublicKey)
	return publickey, nil
}

// 私钥base64字符生成私钥
func LoadBase64PrivateKey(base64Privatekey string) (*rsa.PrivateKey, error) {
	keybytes, err := base64.StdEncoding.DecodeString(base64Privatekey)
	if err != nil {
		return nil, err
	}
	privatekey, err := x509.ParsePKCS1PrivateKey(keybytes)
	if err != nil {
		return nil, err
	}

	return privatekey, nil
}

//公钥bas64字符生成公钥
func LoadBase64PublicKey(base64key string) (*rsa.PublicKey, error) {
	keybytes, err := base64.StdEncoding.DecodeString(base64key)
	if err != nil {
		return nil, err
	}

	pubkeyinterface, err := x509.ParsePKIXPublicKey(keybytes)
	if err != nil {
		return nil, err
	}

	publickey := pubkeyinterface.(*rsa.PublicKey)
	return publickey, nil
}
