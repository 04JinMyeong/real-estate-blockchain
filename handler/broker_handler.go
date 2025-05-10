package handler

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "realestate/database"
    "realestate/models"
)

// RegisterBrokerRequest payload for broker registration
type RegisterBrokerRequest struct {
    Name          string `json:"name"`           // Broker full name
    LicenseNo     string `json:"license_no"`     // Broker license number
    OfficeAddress string `json:"office_address"` // Optional office address
}

// RegisterBroker issues a DID and VC for a licensed broker
func RegisterBroker(c *gin.Context) {
    var req RegisterBrokerRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다"})
        return
    }

    // Create VC struct
    vc := models.BrokerVC{
        Name:      req.Name,
        LicenseNo: req.LicenseNo,
        Issuer:    "did:realestate:platform",
        IssuedAt:  time.Now().UTC(),
    }

    // Sign VC
    sig, err := signVC(vc)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "VC 서명 실패", "detail": err.Error()})
        return
    }
    vc.Signature = sig

    // Store VC in database
    if err := database.StoreBrokerVC(&vc); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "VC 저장 실패", "detail": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "✅ 공인중개사 VC 발급 완료", "vc": vc})
}

// VerifyBrokerRequest payload for VC verification
type VerifyBrokerRequest struct {
    ID string `json:"id"` // Broker DID
}

// VerifyBroker checks the validity of a broker's VC
func VerifyBroker(c *gin.Context) {
    var req VerifyBrokerRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청입니다"})
        return
    }

    // Retrieve stored VC
    vc, err := database.GetBrokerVC(req.ID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "VC를 찾을 수 없습니다"})
        return
    }

    // Recalculate signature
    expectedSig, err := signVC(*vc)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "서명 검증 실패", "detail": err.Error()})
        return
    }
    if vc.Signature != expectedSig {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "VC 검증 실패"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "VC 검증 성공", "vc": vc})
}

// signVC computes a SHA-256 hash of the VC content as a signature
func signVC(vc models.BrokerVC) (string, error) {
    temp := vc
    temp.Signature = ""
    data, err := json.Marshal(temp)
    if err != nil {
        return "", err
    }
    sum := sha256.Sum256(data)
    return hex.EncodeToString(sum[:]), nil
}
