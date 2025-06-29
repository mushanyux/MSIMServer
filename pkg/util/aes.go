package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// AesEncryptSimple 加密
func AesEncryptSimple(origData []byte, key string, iv string) ([]byte, error) {
	return AesDecryptPkcs5(origData, []byte(key), []byte(iv))
}

// AesEncryptPkcs5 加密Pkcs5
func AesEncryptPkcs5(origData []byte, key []byte, iv []byte) ([]byte, error) {
	return AesEncrypt(origData, key, iv, PKCS5Padding)
}

// AesEncryptPkcs7 加密Pkcs7
func AesEncryptPkcs7(origData []byte, key []byte, iv []byte) ([]byte, error) {
	return AesEncrypt(origData, key, iv, PKCS7Padding)
}

// AesEncryptPkcs7Base64 加密Pkcs7Base64
func AesEncryptPkcs7Base64(origData []byte, key []byte, iv []byte) ([]byte, error) {
	data, err := AesEncrypt(origData, key, iv, PKCS7Padding)
	if err != nil {
		return nil, err
	}
	dataStr := base64.StdEncoding.EncodeToString(data)
	return []byte(dataStr), nil
}

// AesEncrypt 加密
func AesEncrypt(origData []byte, key []byte, iv []byte, paddingFunc func([]byte, int) []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = paddingFunc(origData, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AesDecryptSimple 解密
func AesDecryptSimple(crypted []byte, key string, iv string) ([]byte, error) {
	return AesDecryptPkcs5(crypted, []byte(key), []byte(iv))
}

// AesDecryptPkcs5 解密
func AesDecryptPkcs5(crypted []byte, key []byte, iv []byte) ([]byte, error) {
	return AesDecrypt(crypted, key, iv, PKCS5UnPadding)
}

// AesDecryptPkcs7 解密
func AesDecryptPkcs7(crypted []byte, key []byte, iv []byte) ([]byte, error) {
	return AesDecrypt(crypted, key, iv, PKCS7UnPadding)
}

// AesDecryptPkcs7Base64 解密
func AesDecryptPkcs7Base64(crypted []byte, key []byte, iv []byte) ([]byte, error) {
	cryptedData, err := base64.StdEncoding.DecodeString(string(crypted))
	if err != nil {
		return nil, err
	}
	return AesDecrypt(cryptedData, key, iv, PKCS7UnPadding)
}

// AesDecrypt 解密
func AesDecrypt(crypted, key []byte, iv []byte, unPaddingFunc func([]byte) []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = unPaddingFunc(origData)
	return origData, nil
}

// PKCS5Padding
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5UnPadding
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length < unpadding {
		return []byte("unpadding error")
	}
	return origData[:(length - unpadding)]
}

// PKCS7Padding
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length < unpadding {
		return []byte("unpadding error")
	}
	return origData[:(length - unpadding)]
}
