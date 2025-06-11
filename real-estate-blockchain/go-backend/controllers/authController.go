package controllers

import (
	"net/http"
	"realestate/database"
	"realestate/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	// Fabric SDK
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	// for JWT if needed
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

func SignUp(c *gin.Context) {
	// --- 1. HTTP Body → User 구조체 바인딩
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// --- 2. 비밀번호 해시
	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashed)

	// --- 3. 로컬 DB에 저장
	db := database.GetDB()
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// --- 4. Fabric CA 등록(Register) & 인증서 발급(Enroll)
	//    * TLS cert 경로는 환경 변수로 세팅하거나 connection YAML에서 읽힘
	sdk, err := fabsdk.New(config.FromFile("./connection-org1.yaml"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fabric SDK init failed"})
		return
	}
	defer sdk.Close()

	mspClient, err := msp.New(sdk.Context(), msp.WithOrg("Org1"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MSP client init failed"})
		return
	}

	// Register 사용자
	regReq := &msp.RegistrationRequest{
		Name:        user.ID,            // DB의 id 컬럼
		Secret:      user.Password,      // 해시된 비밀번호가 아니라 평문 secret이어야 함
		Affiliation: "org1.department1", // 조직 설정에 맞춰 조정
		Type:        "client",
	}
	if _, err := mspClient.Register(regReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fabric CA register failed", "detail": err.Error()})
		return
	}

	// Enroll 사용자
	if err := mspClient.Enroll(user.ID, msp.WithSecret(user.Password)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fabric CA enroll failed", "detail": err.Error()})
		return
	}

	// Wallet에 identity 저장
	wallet, err := gateway.NewFileSystemWallet("./wallet")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wallet open failed", "detail": err.Error()})
		return
	}
	signingID, err := mspClient.GetSigningIdentity(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "GetSigningIdentity failed", "detail": err.Error()})
		return
	}
	cert := signingID.EnrollmentCertificate()
	key, _ := signingID.PrivateKey().Bytes()
	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))
	if err := wallet.Put(user.ID, identity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wallet put failed", "detail": err.Error()})
		return
	}

	// --- 5. (선택) 바로 로그인 토큰 발급
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.ID,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	// --- 6. 성공 응답
	c.JSON(http.StatusOK, gin.H{
		"message": "Signup and wallet registration successful",
		"token":   tokenString,
	})
}

// func Login(c *gin.Context) {
// 	var user, foundUser models.User
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	db := database.GetDB()
// 	if err := db.Where("username = ?", user.Username).First(&foundUser).Error; err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
// 		return
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
// 		return
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"username": foundUser.Username,
// 		"exp":      time.Now().Add(24 * time.Hour).Unix(),
// 	})

// 	tokenString, err := token.SignedString(jwtKey)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"token": tokenString})
// }
