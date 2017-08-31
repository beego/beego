package session

import (
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"time"
)

type IdGeneration interface {
	UpgradeRawSessionID(sid string) (string, error) //raw session
	GetSessionID(rawSid string) (string, error)
	GenerateSessionID() (string, error)
}

type JwtTokenGenerator struct {
	config *ManagerConfig
	hmac   *jwt.SigningMethodHMAC
}

func NewJwtTokenGenerator(config *ManagerConfig) *JwtTokenGenerator {
	return &JwtTokenGenerator{
		config: config,
		hmac:   jwt.SigningMethodHS256,
	}
}

//通过session id 生成 jwt session
func (j *JwtTokenGenerator) UpgradeRawSessionID(sid string) (string, error) {
	life := int64(j.config.CookieLifeTime) * int64(time.Second)
	now := jwt.TimeFunc()
	expTime := now.Add(time.Duration(life)).Unix()
	claims := &jwt.StandardClaims{ExpiresAt: expTime, Id: sid, IssuedAt: now.Unix()}
	token := jwt.NewWithClaims(j.hmac, claims)
	return token.SignedString([]byte(j.config.JwtSecKey))
}

//检查jwt session ，并且返回session id
func (j *JwtTokenGenerator) GetSessionID(rawSid string) (string, error) {
	token, err := jwt.ParseWithClaims(rawSid, &jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.config.JwtSecKey), nil
		})
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims.Id, nil
	}
	return "", err
}

func (j *JwtTokenGenerator) GenerateSessionID() (string, error) {
	b := make([]byte, j.config.SessionIDLength)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return "", fmt.Errorf("Could not successfully read from the system CSPRNG")
	}
	sid := hex.EncodeToString(b)
	return j.UpgradeRawSessionID(sid)
}
