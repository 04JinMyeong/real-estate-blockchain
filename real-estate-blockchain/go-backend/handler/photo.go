package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadPhoto(c *gin.Context) {
	// 파일 받기
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "파일을 받지 못했습니다."})
		return
	}

	// 저장 경로 및 파일 이름 설정 (원래 이름 사용)
	filename := filepath.Base(file.Filename)
	savePath := fmt.Sprintf("./uploads/%s", filename)

	// 파일 저장
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "파일 저장 실패"})
		return
	}

	// 접근 가능한 URL 리턴
	photoURL := fmt.Sprintf("http://localhost:8080/uploads/%s", filename)
	c.JSON(http.StatusOK, gin.H{"photoUrl": photoURL})
}
