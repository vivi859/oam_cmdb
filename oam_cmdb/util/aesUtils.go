package util

import (
	conf "OAM/conf"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"github.com/beego/beego/v2/core/logs"
)

type EncryptMode string

const (
	CBC EncryptMode = "cbc"
	CFB EncryptMode = "cfb"
)

//aes加密字符串,使用默认参数:默认密钥,CBC模式,BASE64编码
func AesEncryptStr(rawData string) (string, error) {
	ciphertext, err := AesEncryptToBase64Str(rawData, conf.GlobalCfg.SYMMETRIC_KEY, CBC)
	if err != nil {
		return "", errors.New("aes加密失败")
	}
	return ciphertext, nil
}

//aes解密,使用默认参数:默认密钥,CBC模式,BASE64编码
func AesDecryptStr(base64Ciphertext string) (string, error) {
	plaintext, err := AesDecryptBase64Str(base64Ciphertext, conf.GlobalCfg.SYMMETRIC_KEY, CBC)
	if err != nil {
		return "", errors.New("aes加密失败")
	}
	return plaintext, nil
}

//aes加密字符,返回base64编码的字符串
func AesEncryptToBase64Str(rawData string, key string, mode EncryptMode) (string, error) {
	bs, err := AesEncrypt([]byte(rawData), []byte(key), mode)
	if err != nil {
		return "", nil
	}
	return base64.StdEncoding.EncodeToString(bs), err
}

//aes解密字符串,字符串是base64编码后的密文
func AesDecryptBase64Str(base64Ciphertext string, key string, mode EncryptMode) (string, error) {
	bs, err := base64.StdEncoding.DecodeString(base64Ciphertext)
	if err != nil {
		return "", err
	}
	cs, err := AesDecrypt(bs, []byte(key), mode)
	if err != nil {
		return "", err
	}
	return string(cs), nil
}

//aes加密
func AesEncrypt(rawData []byte, key []byte, mode EncryptMode) ([]byte, error) {
	var keySize = len(key)
	checkSize(keySize)
	if mode == CBC {
		return AesEncryptCBC(rawData, key)
	} else if mode == CFB {
		return AesEncryptCFB(rawData, key)
	}
	return nil, errors.New("不支持的加密模式")
}

//aes解密
func AesDecrypt(rawData []byte, key []byte, mode EncryptMode) ([]byte, error) {
	var keySize = len(key)
	checkSize(keySize)
	if mode == CBC {
		return AesDecryptCBC(rawData, key)
	} else if mode == CFB {
		return AesDecryptCFB(rawData, key)
	}
	return nil, errors.New("不支持的加密模式")
}

// aes/cbc加密
func AesEncryptCBC(rawData []byte, key []byte) ([]byte, error) {
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize() // 获取秘钥块的长度
	defer func() {
		err1 := recover()
		if err1 != nil {
			logs.Error(err1)
			err = errors.New("aes加密异常")
		}
	}()
	rawData = pkcs5Padding(rawData, blockSize)                  // pcks5填充
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式,向量iv从key中取
	encryptedData := make([]byte, len(rawData))                 // 创建数组
	blockMode.CryptBlocks(encryptedData, rawData)               // 加密
	return encryptedData, nil
}

// aes/cbc解密
func AesDecryptCBC(encryptedData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	defer func() {
		err1 := recover()
		if err1 != nil {
			logs.Error(err1)
			err = errors.New("aes解密异常")
		}
	}()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	rawData := make([]byte, len(encryptedData))
	blockMode.CryptBlocks(rawData, encryptedData)
	rawData = pkcs5UnPadding(rawData)
	return rawData, nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// aes/cfb加密
func AesEncryptCFB(rawData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	encryptedData := make([]byte, aes.BlockSize+len(rawData))
	iv := encryptedData[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encryptedData[aes.BlockSize:], rawData)
	return encryptedData, nil
}

//aes/cfb解密
func AesDecryptCFB(encryptedData []byte, key []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key)
	if len(encryptedData) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := encryptedData[:aes.BlockSize]
	encryptedData = encryptedData[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encryptedData, encryptedData)
	return encryptedData, nil
}

func checkSize(keySize int) {
	if keySize == 0 || (keySize != 16 && keySize != 24 && keySize != 32) {
		panic("密钥不正确")
	}
}
