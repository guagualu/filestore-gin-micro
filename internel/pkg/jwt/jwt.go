package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

type Claims struct { //载荷
	UserUuid string
	jwt.StandardClaims
}

var jwtpwd = "thegua"

func Getjwtpwd() []byte {
	return []byte(jwtpwd)
}

func GenerateToken(userUuid string) (string, error) {
	Claims := Claims{
		UserUuid: userUuid,
		StandardClaims: jwt.StandardClaims{ //可以直接使用包含的结构
			//ExpiresAt: 60 * 60 * 5 * 10,
			Issuer: "gua",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	tokenstring, err := token.SignedString(Getjwtpwd())
	if err != nil {
		fmt.Println("token genrate err:", err)
		return "", err
	}
	return tokenstring, nil
}

func ParseToken(tokenstring string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenstring, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return Getjwtpwd(), nil
	})
	if err != nil {
		fmt.Println("token parse err:", err)
		return "", err
	}
	if claims, ok := token.Claims.(Claims); ok && token.Valid { //ok 是指interface转化成功 token.valid 是token自身带的属性 可以用于判断是否过期
		return claims.UserUuid, nil
	}
	return "", err
}
