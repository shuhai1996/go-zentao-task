package util

import (
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm3"
	"github.com/tjfoc/gmsm/sm4"
)

func SM3(plainText string) string {
	h := sm3.New()
	if _, err := h.Write([]byte(plainText)); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func SM2GenerateKey() (publicKey string, privateKey string) {
	priv, _ := sm2.GenerateKey()
	pub := &priv.PublicKey

	var b []byte
	var err error
	b, err = sm2.MarshalSm2PublicKey(pub)
	if err != nil {
		return "", ""
	}
	publicKey = base64.StdEncoding.EncodeToString(b)
	b, err = sm2.MarshalSm2PrivateKey(priv, nil)
	if err != nil {
		return "", ""
	}
	privateKey = base64.StdEncoding.EncodeToString(b)
	return
}

func SM2Encrypt(publicKey, plainText string) (string, error) {
	var b []byte
	var err error
	var pub *sm2.PublicKey

	b, err = base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", err
	}
	pub, err = sm2.ParseSm2PublicKey(b)
	if err != nil {
		return "", err
	}
	b, err = pub.Encrypt([]byte(plainText))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func SM2Decrypt(privateKey, cipherText string) (string, error) {
	var b []byte
	var err error
	var priv *sm2.PrivateKey

	b, err = base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	priv, err = sm2.ParsePKCS8UnecryptedPrivateKey(b)
	if err != nil {
		return "", err
	}
	b, err = base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	b, err = priv.Decrypt(b)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func SM2Sign(privateKey, msg string) (string, error) {
	var b []byte
	var err error
	var priv *sm2.PrivateKey

	b, err = base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	priv, err = sm2.ParsePKCS8UnecryptedPrivateKey(b)
	if err != nil {
		return "", err
	}
	r, s, e := sm2.Sign(priv, []byte(msg))
	if e != nil {
		return "", e
	}
	b, err = sm2.SignDigitToSignData(r, s)
	if err != nil {
		return "", nil
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func SM2Verify(publicKey, sig, msg string) (bool, error) {
	var b []byte
	var err error
	var pub *sm2.PublicKey

	b, err = base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return false, err
	}
	pub, err = sm2.ParseSm2PublicKey(b)
	if err != nil {
		return false, err
	}
	b, err = base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return false, err
	}
	r, s, e := sm2.SignDataToSignDigit(b)
	if e != nil {
		return false, err
	}
	return sm2.Verify(pub, []byte(msg), r, s), nil
}

func SM4Encrypt(key, plainText string) (string, error) {
	kb := []byte(key)
	pb := []byte(plainText)
	block, err := sm4.NewCipher(kb)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origData := PKCS7Padding(pb, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, kb[:blockSize])
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted), nil
}

func SM4Decrypt(key, cipherText string) (string, error) {
	kb := []byte(key)
	cb, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	block, e := sm4.NewCipher(kb)
	if e != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, kb[:blockSize])
	origData := make([]byte, len(cipherText))
	blockMode.CryptBlocks(origData, cb)
	return string(PKCS7UnPadding(origData)), nil
}
