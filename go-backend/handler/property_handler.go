// 📁 realestate/handler/property_handler.go
package handler

import (
	"encoding/json"
	"net/http"
	"realestate/blockchain"
	"realestate/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 요청 구조체: 매물 등록 및 조회에 사용
type PropertyRequest struct {
	User     string `json:"user"`
	Address  string `json:"address,omitempty"`
	Owner    string `json:"owner,omitempty"`
	Price    string `json:"price,omitempty"`
	PhotoURL string `json:"photoUrl,omitempty"`
}

// ✅ 매물 등록
func AddProperty(c *gin.Context) {
	var req PropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청 형식입니다."})
		return
	}

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
