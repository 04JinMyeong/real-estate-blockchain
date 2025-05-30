package main

import (
	"fmt"
	"log"

	"realestate/did"
)

func main() {
    didStr, err := did.GenerateDID()
    if err != nil {
        log.Fatalf("DID 생성 실패: %v", err)
    }
    fmt.Println("생성된 DID:", didStr)
}
