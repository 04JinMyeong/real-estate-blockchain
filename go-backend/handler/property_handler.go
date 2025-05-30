// 📁 realestate/handler/property_handler.go
package handler

import (
	"encoding/json"
	"net/http"
	"realestate/vc"
	"realestate/blockchain"
	"realestate/utils"
	"strings"
	"time"
	"log"

	"github.com/gin-gonic/gin"
)

// 요청 구조체: 매물 등록 및 조회에 사용
type PropertyRequest struct {
	User     string `json:"user"`
	Address  string `json:"address,omitempty"`
	Owner    string `json:"owner,omitempty"`
	Price    string `json:"price,omitempty"`
	PhotoURL string `json:"photoUrl,omitempty"`
	VC       string `json:"vc" binding:"required"`       // 제출된 VC (JSON 문자열 형태), 필수 필드
}

// ✅ 매물 등록
func AddProperty(c *gin.Context) {
	var req PropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청 형식입니다."})
		return
	}

	log.Printf("매물 등록 요청 수신: User=%s, Address=%s, VC 소지 여부=%t\n", req.User, req.Address, req.VC != "")

	// propertyID := 부분까지 VC 검증 로직 추가 -----

	// --- 1. VC 기본 유효성 검증 (서명 등) ---
	isValidVC, verificationError := vc.VerifyVC(req.VC) // vc.VerifyVC 함수 호출
	if verificationError != nil {
		log.Printf("VC 검증 오류 (User: %s): %v\n", req.User, verificationError)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "VC 검증 중 오류 발생", "detail": verificationError.Error()})
		return
	}
	if !isValidVC {
		log.Printf("유효하지 않은 VC 제출 (User: %s)\n", req.User)
		c.JSON(http.StatusForbidden, gin.H{"error": "제출된 VC가 유효하지 않습니다. 위변조되었거나 잘못된 VC일 수 있습니다."})
		return
	}
	log.Printf("✅ VC 기본 유효성 검증 통과 (User: %s)\n", req.User)

	// --- 2. VC 클레임 기반 자격 확인 ---
	var vcData map[string]interface{}
	if err := json.Unmarshal([]byte(req.VC), &vcData); err != nil { // VC 문자열을 JSON 객체로 파싱
		log.Printf("VC JSON 파싱 오류 (User: %s): %v\n", req.User, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "VC 내용 분석 중 오류 발생", "detail": "VC JSON 파싱 실패"})
		return
	}

	credentialSubject, ok := vcData["credentialSubject"].(map[string]interface{})
	if !ok {
		log.Printf("VC credentialSubject 형식 오류 (User: %s)\n", req.User)
		c.JSON(http.StatusBadRequest, gin.H{"error": "제출된 VC의 credentialSubject 형식이 올바르지 않습니다."})
		return
	}

	statusClaim, ok := credentialSubject["status"].(string) // status 클레임 값 추출
	if !ok {
		log.Printf("VC status 클레임 누락 또는 타입 오류 (User: %s)\n", req.User)
		c.JSON(http.StatusBadRequest, gin.H{"error": "제출된 VC의 status 클레임이 없거나 형식이 올바르지 않습니다."})
		return
	}
	log.Printf("VC 'status' 클레임 값 (User: %s): %s\n", req.User, statusClaim)

	// 정책 적용: statusClaim 값을 기준으로 매물 등록 허용/거부
	if statusClaim == "fraudRecord_Exists" { // "전과 기록 보유" VC의 status 값
		log.Printf("매물 등록 거부 (User: %s, Status: %s): 전과 기록 확인\n", req.User, statusClaim)
		c.JSON(http.StatusForbidden, gin.H{"error": "자격 증명 확인 결과, 매물 등록이 불가능한 상태입니다 (사유: 전과 기록)."})
		return
	}
	if statusClaim != "valid" { // "정상" VC의 status 값이 아니면 거부
		log.Printf("매물 등록 거부 (User: %s, Status: %s): 유효하지 않은 자격 상태\n", req.User, statusClaim)
		c.JSON(http.StatusForbidden, gin.H{"error": "자격 증명 상태(" + statusClaim + ")가 유효하지 않아 매물 등록이 불가능합니다."})
		return
	}
	log.Printf("✅ VC 클레임 기반 자격 확인 통과 (User: %s, Status: %s)\n", req.User, statusClaim)

	// --- 3. 모든 검증 통과 시 매물 정보 블록체인 등록 ---

	propertyID := utils.GeneratePropertyID(req.Address)

	err := blockchain.SubmitAddListing(req.User, propertyID, req.Address, req.Owner, req.Price, req.PhotoURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "매물 등록 실패", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "✅ 매물 등록 완료!", "id": propertyID})
}

// ✅ 매물 조회 (ID로)
func GetProperty(c *gin.Context) {
	id := c.Param("id")
	user := c.Query("user")

	if user == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user 쿼리 파라미터가 필요합니다"})
		return
	}

	result, err := blockchain.QueryProperty(user, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "매물 조회 실패", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"property": result})
}

// ✅ 전체 매물 조회
func GetAllProperties(c *gin.Context) {
	result, err := blockchain.QueryAllProperties("admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "전체 매물 조회 실패",
			"detail": err.Error(),
		})
		return
	}

	c.Data(http.StatusOK, "application/json", []byte(result))
	/*user := c.Query("user")
	if user == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user 쿼리 파라미터가 필요합니다"})
		return
	}

	result, err := blockchain.QueryAllProperties(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "매물 전체 조회 실패", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)*/
}

// 매물 수정
func UpdateProperty(c *gin.Context) {
	var req struct {
		User  string `json:"user"`
		ID    string `json:"id"`
		Owner string `json:"owner"`
		Price string `json:"price"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청 형식입니다"})
		return
	}

	err := blockchain.SubmitUpdateListing(req.User, req.ID, req.Owner, req.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "매물 수정 실패", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "✅ 매물 수정 완료!"})
}

// 매물 예약
func ReserveProperty(c *gin.Context) {
	var req struct {
		User string `json:"user"`
		ID   string `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "요청 형식이 올바르지 않습니다"})
		return
	}

	expiresAt := time.Now().Add(12 * time.Hour).Unix()

	err := blockchain.ReserveListing(req.User, req.ID, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "예약 실패",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "✅ 예약 완료!",
		"expiresAt": expiresAt,
	})
}

// ✅ 내가 올린 매물만 조회 (CreatedBy 기준)
func GetMyProperties(c *gin.Context) {
	user := c.Query("user")
	if user == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user 쿼리 파라미터가 필요합니다"})
		return
	}

	resultStr, err := blockchain.QueryAllProperties(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "조회 실패", "detail": err.Error()})
		return
	}

	trimmed := strings.TrimSpace(resultStr)
	if trimmed == "" || trimmed == "null" {
		c.JSON(http.StatusOK, gin.H{"properties": []map[string]interface{}{}})
		return
	}

	var listings []map[string]interface{}
	if err := json.Unmarshal([]byte(resultStr), &listings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JSON 파싱 실패", "detail": err.Error()})
		return
	}

	var myListings []map[string]interface{}
	for _, l := range listings {
		if createdBy, ok := l["createdBy"].(string); ok && createdBy == user {
			myListings = append(myListings, l)
		}
	}
	c.JSON(http.StatusOK, gin.H{"properties": myListings})
}
