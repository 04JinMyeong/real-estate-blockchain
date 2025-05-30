// File: go-backend/crypto/keys.go
package crypto

import (
    "crypto/ed25519"
    "encoding/base64"
    "errors"
    "io/ioutil"
    "os"
)

// LoadPrivateKey reads a Base64-encoded Ed25519 private key from the path in env var.
func LoadPrivateKey(envVar string) (ed25519.PrivateKey, error) {
    path := os.Getenv(envVar)
    if path == "" {
        return nil, errors.New(envVar + " 환경 변수가 설정되어 있지 않습니다")
    }
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    raw, err := base64.StdEncoding.DecodeString(string(data))
    if err != nil {
        return nil, err
    }
    if len(raw) != ed25519.PrivateKeySize {
        return nil, errors.New("잘못된 개인키 크기")
    }
    return ed25519.PrivateKey(raw), nil
}

// LoadPublicKey reads a Base64-encoded Ed25519 public key from the path in env var.
func LoadPublicKey(envVar string) (ed25519.PublicKey, error) {
    path := os.Getenv(envVar)
    if path == "" {
        return nil, errors.New(envVar + " 환경 변수가 설정되어 있지 않습니다")
    }
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    raw, err := base64.StdEncoding.DecodeString(string(data))
    if err != nil {
        return nil, err
    }
    if len(raw) != ed25519.PublicKeySize {
        return nil, errors.New("잘못된 공개키 크기")
    }
    return ed25519.PublicKey(raw), nil
}

// Sign 데이터에 개인키로 디지털 서명하고 서명 바이트를 반환합니다.
func Sign(data []byte) ([]byte, error) {
    priv, err := LoadPrivateKey("PRIVATE_KEY_PATH")
    if err != nil {
        return nil, err
    }
    sig := ed25519.Sign(priv, data)
    return sig, nil
}

// Verify 데이터와 서명이 공개키로부터 유효한지 검증합니다.
func Verify(data, sig []byte) (bool, error) {
    pub, err := LoadPublicKey("PUBLIC_KEY_PATH")
    if err != nil {
        return false, err
    }
    return ed25519.Verify(pub, data, sig), nil
}
