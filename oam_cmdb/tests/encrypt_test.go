package test

import (
	"OAM/models"
	fn "OAM/util"
	"fmt"
	"testing"
)

func TestAes(t *testing.T) {
	//origData := []byte("Hello World中国字") // 待加密的数据
	key := "123456789qwertyu" // 加密的密钥
	//log.Println("原文：", string(origData))
	plainStr := "中国字pwd123456789qwertyu"
	cphierStr, _ := fn.AesEncryptToBase64Str(plainStr, key, fn.CBC)
	fmt.Println("cbc密文:" + cphierStr)

	ps, err := fn.AesDecryptBase64Str(cphierStr, key, fn.CBC)
	if err == nil {
		fmt.Println("cbc解密后明文:" + ps)
	} else {
		fmt.Println(err)
	}
	// cfb
	cphierStr1, _ := fn.AesEncryptToBase64Str(plainStr, key, fn.CFB)
	fmt.Println("cfb密文:" + cphierStr1)

	ps1, _ := fn.AesDecryptBase64Str(cphierStr1, key, fn.CFB)
	fmt.Println("cfb解密后明文:" + ps1)
}

func TestSha(t *testing.T) {
	s := fn.SHA256Hex("中国人民解放军123456789qwertyuiopasd@#@")
	fmt.Println(s)

	u := models.UserInfo{UserName: "root", Passwd: "2022@00"}
	u.EncryptPasswd()
	fmt.Println("加密后密码:" + u.Passwd)
}

func TestRsa(t *testing.T) {
	err := fn.RSAGenerateKeyFile(1024, "d:/zppri_.pem", "d:/zppub_.pem")
	fmt.Print(err)

	/* var cryptb64 = "c3Ed62FuEhVrS6xnXd0BY0gON1F3rg8L0VJOa7RuajcpnUNwXnMQCtST2x4KJjzvM3fMFLWDaT6SJk75qAib8EiQSzjXIPTl909+i1fjZfWsoolecmRRKISxr+AO5fHgU7Hr0Oxx/tBCifxdve+UBilfW6D6anTKZvUm7BtJofM="

	plainbytes, err := fn.RSADecryptBase64Str(cryptb64, "d:/pri.pem")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(string(plainbytes)) */
}
