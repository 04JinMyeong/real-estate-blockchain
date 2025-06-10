package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// key: 중개사무소 등록번호, value: 대표자 이름
var validAgents = map[string]string{
	//"11110-2015-00123": "김중개",
	//"41287-2023-00045": "박성실",
	//"26290-2020-00088": "이신뢰",
	"12345-2019-00001": "최정직",
	"67890-2021-00099": "홍범인",
	//"54321-2022-00022": "전과자",
	"98765-2020-00033": "전과쟈",
	"13579-2018-00044": "윤희망", // did:realestate:bd601d135552360ff6924828551db554f6cc171ad49f9aea69510a6e98d7045a   OPHblkozTPPS8Ddsf5C2ChVAcVdAo5WUG3gHMtKiyjwU86TlJuNCCi0yzv7VPS0tFmStFOKzmyFx8vHWmyNKww==
}

func main() {
	router := gin.Default()
	// CORS 문제를 피하기 위해 간단한 CORS 미들웨어 추가
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	router.GET("/verify", func(c *gin.Context) {
		agentName := c.Query("agentName")
		regNum := c.Query("registrationNumber")

		if name, ok := validAgents[regNum]; ok && name == agentName {
			// 유효한 공인중개사일 경우
			c.JSON(http.StatusOK, gin.H{"totalCount": 1})
		} else {
			// 유효하지 않을 경우
			c.JSON(http.StatusOK, gin.H{"totalCount": 0})
		}
	})
	// 메인 서버와 다른 포트(8081)에서 실행
	router.Run(":8081")
}
