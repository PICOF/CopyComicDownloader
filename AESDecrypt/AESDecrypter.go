package AESDecrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

func AesDecrypt(decodeStr string, key []byte, iv []byte) ([]byte, error) {
	//decodeBytes, err := base64.StdEncoding.DecodeString(decodeStr)//先解密base64
	decodeBytes, err := hex.DecodeString(decodeStr)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(decodeBytes))

	blockMode.CryptBlocks(origData, decodeBytes)
	origData = ZeroUnPadding(origData) // origData = PKCS5UnPadding(origData)
	return origData, nil
}
func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
