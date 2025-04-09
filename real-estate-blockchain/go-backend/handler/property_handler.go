package handler

import (
	"net/http"
	"realestate/blockchain"

	"github.com/gin-gonic/gin"
)

// 요청 구조체: 매물 등록 및 조회에 사용
type PropertyRequest struct {
	User    string `json:"user"` // ← 사용자 이름
	ID      string `json:"id"`
	Address string `json:"address,omitempty"`
	Owner   string `json:"owner,omitempty"`
	Price   string `json:"price,omitempty"`
}

// ✅ 매물 등록
func AddProperty(c *gin.Context) {
	var req PropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청 형식입니다."})
		return
	}

	err := blockchain.SubmitAddListing(req.User, req.ID, req.Address, req.Owner, req.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "매물 등록 실패", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "✅ 매물 등록 완료!"})
}

// ✅ 매물 조회
func GetProperty(c *gin.Context) {
	id := c.Param("id")
	user := c.Query("user") // 쿼리스트링으로 사용자 받기

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
