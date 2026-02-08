package compiler

import (
	"encoding/json"
	"fmt"

	"github.com/originbeat-inc/runly-cli/internal/i18n"
	"github.com/originbeat-inc/runly-cli/pkg/crypto"
	"github.com/originbeat-inc/runly-cli/pkg/protocol"
)

// BuildArtifact 对协议执行标准化分发处理
func BuildArtifact(proto *protocol.RunlyProtocol, privateKeyHex string) (string, error) {
	// 1. 签名预处理：清空签名位，设定算法
	proto.Security.Signature = ""
	proto.Security.HashAlgo = "SHA-256"

	// 2. 确定性标准化 (Canonicalization)
	// 使用有序 JSON 序列化确保不同环境下哈希计算结果的一致性
	standardizedData, err := json.Marshal(proto)
	if err != nil {
		return "", fmt.Errorf(i18n.T("errors.yaml_unmarshal_fail"), err)
	}

	// 3. 计算内容摘要
	contentHash := crypto.CalculateHash(standardizedData)

	// 4. 执行数字签名（背书）
	signature, err := crypto.Sign(privateKeyHex, []byte(contentHash))
	if err != nil {
		// 返回多语言错误：未找到签名私钥。请使用 'runly-cli keys generate' 创建身份
		return "", fmt.Errorf(i18n.T("errors.no_key"), err)
	}

	// 5. 注入指纹
	proto.Security.Signature = signature

	return contentHash, nil
}

// VerifyIntegrity 验证协议资产的完整性与身份真实性
func VerifyIntegrity(proto *protocol.RunlyProtocol) (bool, error) {
	storedSignature := proto.Security.Signature
	pubKey := proto.Manifest.Creator.PubKey

	// 检查是否存在签名
	if storedSignature == "" {
		// 返回多语言错误：资产未编译签名，请先执行 runly-cli build
		return false, fmt.Errorf(i18n.T("errors.no_sig"))
	}

	// 还原计算态：清空签名位
	proto.Security.Signature = ""
	standardizedData, _ := json.Marshal(proto)
	currentHash := crypto.CalculateHash(standardizedData)

	// 恢复原始对象
	proto.Security.Signature = storedSignature

	// 校验数字签名
	isValid, err := crypto.Verify(pubKey, []byte(currentHash), storedSignature)
	if err != nil || !isValid {
		// 返回多语言错误：安全签名验证未通过：资产已被篡改或公钥不匹配
		return false, fmt.Errorf(i18n.T("errors.sign_verify_fail"))
	}

	return true, nil
}
