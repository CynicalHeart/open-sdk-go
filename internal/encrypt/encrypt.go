package encrypt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// rsaEncode 使用 RSA 公钥加密数据并返回 Base64 编码的密文；
//
// data: 要加密的数据, json字符串
func RsaEncode(data string) (string, error) {
	// 从文件加载公钥
	_, currPath, _, _ := runtime.Caller(0)
	rsaPath := filepath.Join(filepath.Dir(currPath), "RSA-PublicKey.pem")
	pemCont, err := os.ReadFile(rsaPath)
	if err != nil {
		return "", fmt.Errorf("读取公钥文件失败: %v。", err)
	}

	// 解析公钥
	block, _ := pem.Decode(pemCont)
	if block == nil || block.Type != "PUBLIC KEY" {
		return "", fmt.Errorf("从PEM文件中没有找到公钥。")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("解析公钥失败: %v。", err)
	}
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("RSA公钥不符合规范。")
	}

	// 使用RSA公钥加密数据[兼容长文本]
	keySize := rsaPubKey.Size()
	maxChunkSize := keySize - 11    // PKCS1v15 填充需要的空间
	var encryptedBytes bytes.Buffer // 存储加密后的数据
	plaintextBytes := []byte(data)  // 要加密的字节数据

	for len(plaintextBytes) > 0 {
		chunkSize := maxChunkSize
		if len(plaintextBytes) < chunkSize {
			chunkSize = len(plaintextBytes)
		}

		chunk, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, plaintextBytes[:chunkSize])
		if err != nil {
			return "", fmt.Errorf("RSA加密数据块失败: %v", err)
		}
		encryptedBytes.Write(chunk)                 // 向缓冲区写入加密后的数据块
		plaintextBytes = plaintextBytes[chunkSize:] // 更新待加密的数据
	}

	return base64.StdEncoding.EncodeToString(encryptedBytes.Bytes()), nil
}
