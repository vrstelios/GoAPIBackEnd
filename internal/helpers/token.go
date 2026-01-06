package helpers

import (
	config2 "GoAPIBackEnd/internal/config"
	"GoAPIBackEnd/internal/database"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Claims struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`

	jwt.StandardClaims
}

var jwtKey []byte

func SetJWTKey(key string) {
	jwtKey = []byte(key)

}

func GetJWTKey() []byte {
	return []byte(jwtKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	secretKey := GetJWTKey()

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
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
	var dbConn = config2.Conn

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
