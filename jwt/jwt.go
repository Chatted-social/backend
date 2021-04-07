package jwt

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/mitchellh/mapstructure"
)

func NewWithClaims(claims Claims) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

func From(v interface{}) Claims {
	mc := v.(*jwt.Token).Claims.(jwt.MapClaims)

	var claims Claims
	err := mapstructure.Decode(mc, &claims)

	if err != nil {
		return Claims{}
	}

	return claims

}

type Claims struct {
	jwt.StandardClaims
	UserID int `json:"usr,omitempty"`
}
