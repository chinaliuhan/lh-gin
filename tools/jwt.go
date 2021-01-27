package tools

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTClaims struct {
	jwt.StandardClaims
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Mobile   string `json:"mobile"`
}

var (
	Secret     = "To be or not to be, that's a question." //salt
	ExpireTime = 3600                                     //token expire
)

type jwtUtil struct {
}

func NewJwtUtil() *jwtUtil {

	return &jwtUtil{}
}

/**
创建token
*/
func (r *jwtUtil) GenerateToken(claims *JWTClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(Secret))
	if err != nil {
		return ""
	}
	return signedToken
}

/**
解析token
*/
func (r *jwtUtil) ParseToken(tokenString string) *JWTClaims {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		return nil
	}

	//断言
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil
	}
	if err := token.Claims.Valid(); err != nil {
		return nil
	}

	return claims
}

func (r *jwtUtil) Refresh(token string) string {
	claims := r.ParseToken(token)
	if claims != nil {
		return ""
	}
	claims.ExpiresAt = time.Now().Unix() + (claims.ExpiresAt - claims.IssuedAt)
	newToken := r.GenerateToken(claims)
	if newToken == "" {
		return ""
	}

	return newToken
}
