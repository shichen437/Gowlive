package utils

import (
	"context"
	"encoding/hex"

	"github.com/gogf/gf/v2/crypto/gaes"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/crypto/gsha1"
	"github.com/gogf/gf/v2/frame/g"
)

// Decrypt 使用 SM4-CBC 模式解密数据
func Decrypt(ctx context.Context, ciphertextHex string) (string, error) {
	ciphertext, _ := hex.DecodeString(ciphertextHex)
	p, err := gaes.DecryptCBC(ciphertext, getKeyBytes())
	if err != nil {
		g.Log().Error(ctx, "Decrypt error", err)
		return "", nil
	}
	return string(p), nil
}

// Decrypt 使用 SM4-CBC 模式解密数据
func Encrypt(ctx context.Context, ciphertext string) (string, error) {
	p, err := gaes.EncryptCBC(convertToBytes(ciphertext), getKeyBytes())
	if err != nil {
		g.Log().Error(ctx, "Encrypt error", err)
		return "", err
	}
	return hex.EncodeToString(p), nil
}

func getKeyBytes() []byte {
	return convertToBytes(SM4_KEY)
}

func convertToBytes(key string) []byte {
	bytes := make([]byte, len(key))
	copy(bytes, key)
	return bytes
}

// sha1加密
func Sha1En(data string) string {
	t := gsha1.Encrypt(data)
	return t
}

// 对字符串进行MD5哈希
func Md5En(data string) string {
	t, _ := gmd5.Encrypt(data)
	return t
}
