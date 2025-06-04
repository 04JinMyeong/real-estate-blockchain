// handler/photo.go
package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// UploadPhoto handles incoming multipart/form-data with a "photo" field,
// saves the file under "./uploads", and returns a URL that clients can use to fetch it.

func UploadPhoto(c *gin.Context) {
	// 1. 클라이언트가 보낸 "photo" 파일 추출
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "파일을 받지 못했습니다."})
		return
	}

	// 2. uploads 디렉터리가 없으면 생성
	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "업로드 디렉터리 생성 실패"})
		return
	}

	// 3. 파일 이름 안전하게 추출 (원래 이름을 그대로 사용하되, 경로 구분자는 제거)
	filename := filepath.Base(file.Filename)
	savePath := fmt.Sprintf("./uploads/%s", filename)

	// 4. 실제로 디스크에 "./uploads/filename" 위치로 저장
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "파일 저장 실패"})
		return
	}

	// 5. 요청이 들어온 호스트와 프로토콜에 맞춰 접근 가능한 URL 생성
	//    - c.Request.Host 예: "localhost:8080" 또는 "abcd1234.ngrok.io"
	//    - c.Request.TLS가 nil이 아니면 https, nil이면 http
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	photoURL := fmt.Sprintf("%s://%s/uploads/%s", scheme, c.Request.Host, filename)

	// 6. 성공 시 JSON으로 photoUrl 반환
	c.JSON(http.StatusOK, gin.H{"photoUrl": photoURL})
}
