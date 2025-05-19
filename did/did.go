package did

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json" // DID Document JSON 변환용
	"fmt"

	"github.com/mr-tron/base58" // Base58 인코딩을 위해 추가
)

// --- 기존 함수 ---

// GenerateDIDFromPublicKey 주어진 공개키 바이트를 해시하여 DID 문자열을 반환합니다.
// 이 함수는 계속 사용합니다.
func GenerateDIDFromPublicKey(pubKey []byte) string {
	hash := sha256.Sum256(pubKey)
	// hex 인코딩된 해시값을 DID suffix 로 사용
	return fmt.Sprintf("did:realestate:%s", hex.EncodeToString(hash[:]))
}

// --- DID Document 관련 구조체 및 함수 추가 ---

type DIDDocument struct {
	Context            interface{}          `json:"@context"`
	ID                 string               `json:"id"`
	VerificationMethod []VerificationMethod `json:"verificationMethod,omitempty"`
	Authentication     []string             `json:"authentication,omitempty"`
	// AssertionMethod    []string             `json:"assertionMethod,omitempty"`
	// KeyAgreement       []string             `json:"keyAgreement,omitempty"`
	// Service            []Service            `json:"service,omitempty"`
}

type VerificationMethod struct {
	ID                 string `json:"id"`
	Type               string `json:"type"`
	Controller         string `json:"controller"`
	PublicKeyMultibase string `json:"publicKeyMultibase,omitempty"` // 공개키를 Multibase 형식(z 접두사 Base58BTC)으로 표현
}

type Service struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

func CreateAgentDIDDocument(agentDID string, pubKeyBytes []byte, keyType string) (DIDDocument, error) {
	// 공개키를 Base58BTC로 인코딩합니다.
	// Ed25519 공개키의 경우, 멀티코덱 prefix 0xed + version 0x01 = 0xed01
	// 그 후 Base58BTC 인코딩. 여기서는 단순 Base58 인코딩 후 'z'를 붙이는 간소화된 예시를 사용.
	// 더 정확한 Multibase 구현을 위해서는 multiformats/go-multibase 와 multiformats/go-multicodec 라이브러리 사용 고려.
	// 지금은 DID Core의 일반적인 Base58 표현(z로 시작)을 모방합니다.
	publicKeyBase58 := base58.Encode(pubKeyBytes)
	publicKeyMultibase := "z" + publicKeyBase58 // 'z'는 Base58BTC를 나타내는 Multibase 접두사

	vmID := fmt.Sprintf("%s#keys-1", agentDID)

	// Verification Method의 Type에 따라 적절한 컨텍스트를 추가할 수 있습니다.
	contexts := []string{"https://www.w3.org/ns/did/v1"}
	if keyType == "Ed25519VerificationKey2020" {
		contexts = append(contexts, "https://w3id.org/security/suites/ed25519-2020/v1")
	} else if keyType == "Ed25519VerificationKey2018" { // 2018 스펙을 사용한다면
		// contexts = append(contexts, "https://w3id.org/security/suites/ed25519-2018/v1") // 이 컨텍스트 URL이 정확한지 확인 필요
	}

	doc := DIDDocument{
		Context: contexts,
		ID:      agentDID,
		VerificationMethod: []VerificationMethod{
			{
				ID:                 vmID,
				Type:               keyType,
				Controller:         agentDID,
				PublicKeyMultibase: publicKeyMultibase,
			},
		},
		Authentication: []string{
			vmID,
		},
	}
	return doc, nil
}

func (doc *DIDDocument) ToJson() (string, error) {
	jsonBytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling DIDDocument to JSON: %w", err)
	}
	return string(jsonBytes), nil
}

func FromJson(jsonString string) (DIDDocument, error) {
	var doc DIDDocument
	err := json.Unmarshal([]byte(jsonString), &doc)
	if err != nil {
		return DIDDocument{}, fmt.Errorf("error unmarshalling JSON to DIDDocument: %w", err)
	}
	return doc, nil
}
