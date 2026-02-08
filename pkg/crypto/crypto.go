package crypto

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// CalculateHash 计算数据的 SHA-256 摘要
func CalculateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// Sign 使用 Ed25519 私钥执行数字签名
func Sign(privateKeyHex string, message []byte) (string, error) {
	seed, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key format: %w", err)
	}

	privateKey := ed25519.NewKeyFromSeed(seed)
	signature := ed25519.Sign(privateKey, message)

	return hex.EncodeToString(signature), nil
}

// Verify 验证签名合法性
func Verify(publicKeyHex string, message []byte, signatureHex string) (bool, error) {
	pubBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return false, fmt.Errorf("invalid public key: %w", err)
	}

	sigBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, fmt.Errorf("invalid signature format: %w", err)
	}

	return ed25519.Verify(pubBytes, message, sigBytes), nil
}
