package helpers

import (
	"GoAPIBackEnd/internal/database"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	CtxUserId = "userId"
	CtxEmail  = "email"
	CtxRole   = "role"
)

/*func GetUserId(c *gin.Context) string {
	v, _ := c.Get(CtxUserId)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func GetUserEmail(c *gin.Context) string {
	v, _ := c.Get(CtxEmail)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}*/

func GetUserRole(c *gin.Context) string {
	v, _ := c.Get(CtxRole)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

var jwtKey []byte

type Claims struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`

	jwt.StandardClaims
}

func SetJWTKey(key string) {
	jwtKey = []byte(key)

}

func GetJWTKey() []byte {
	return []byte(jwtKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func HashPassword(password string) *string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	hasHedPwd := string(bytes)
	return &hasHedPwd
}

func VerifyPassword(userPwd, pwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(userPwd), []byte(pwd))
	return err == nil, err
}

func GenerateToken(email, userId, userType string) (string, string) {
	tokenExpiry := time.Now().Add(24 * time.Hour).Unix()
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour).Unix()

	claims := &Claims{
		Email:  email,
		UserId: userId,
		Role:   userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiry,
		},
	}
	refreshClaims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpiry,
		},
	}

	// Generating the tokens
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedAccessRefreshToken, err := accessToken.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	return signedAccessRefreshToken, signedRefreshToken
}

func UpdatedAllToken(signedToken, signedRefreshToken, userId string) error {
	var dbConn = database.Conn

	user, err := database.GetUser(dbConn, userId, "", "")
	if err != nil {
		return err
	}
	user[0].Token = &signedToken
	user[0].RefreshToken = &signedRefreshToken

	err = database.PutUser(dbConn, *user[0])
	if err != nil {
		return err
	}

	return nil
}
