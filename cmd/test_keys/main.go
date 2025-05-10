// File: go-backend/crypto/test_keys.go

package main

import (
	"fmt"
	"realestate/go-backend/crypto"
)



func main() {
    data := []byte("테스트 데이터")
    sig, err := crypto.Sign(data)
    if err != nil {
        panic("Sign 실패: " + err.Error())
    }
    ok, err := crypto.Verify(data, sig)
    if err != nil {
        panic("Verify 실패: " + err.Error())
    }
    fmt.Println("서명 검증 결과:", ok) // → true 가 나와야 합니다.
}
