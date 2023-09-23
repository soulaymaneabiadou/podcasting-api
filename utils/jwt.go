package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtPayload struct {
	jwt.RegisteredClaims
	ID uint `json:"id"`
}

type tokenOpts struct {
	secret string
	exp    time.Time
	sub    string
	p      JwtPayload
}

func createToken(opts tokenOpts) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtPayload{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("DOMAIN"),
			Audience:  jwt.ClaimStrings{os.Getenv("DOMAIN")},
			Subject:   opts.sub,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(opts.exp),
		},
		ID: opts.p.ID,
	})

	return token.SignedString([]byte(opts.secret))
}

func CreateAccessToken(sub string, p JwtPayload) (string, error) {
	expDur, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRE"))
	exp := time.Now().Add(time.Minute * time.Duration(expDur)) // minutes

	return createToken(tokenOpts{
		secret: os.Getenv("JWT_ACCESS_SECRET"),
		exp:    exp,
		sub:    sub,
		p:      p,
	})
}

func CreateRefreshToken(sub string, p JwtPayload) (string, error) {
	expDur, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRE"))
	exp := time.Now().Add(time.Hour * 24 * time.Duration(expDur)) // days

	token, err := createToken(tokenOpts{
		secret: os.Getenv("JWT_REFRESH_SECRET"),
		exp:    exp,
		sub:    sub,
		p:      p,
	})

	// we can save token the token in order to perform whitelisting
	return token, err
}

func CreateTokens(sub string, p JwtPayload) (string, string, error) {
	accessToken, err := CreateAccessToken(sub, p)
	if err != nil {
		return "", "", fmt.Errorf("error creating access jwt: %v", err)
	}

	refreshToken, err := CreateRefreshToken(sub, p)
	if err != nil {
		return "", "", fmt.Errorf("error creating refresh jwt: %v", err)
	}

	return accessToken, refreshToken, nil
}

func parseJwt(t string, secret string) (JwtPayload, error) {
	var jwtPayload JwtPayload

	token, err := jwt.ParseWithClaims(t, &jwtPayload, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return JwtPayload{}, err
	}

	return jwtPayload, nil
}

func ParseAccessJWT(t string) (JwtPayload, error) {
	return parseJwt(t, os.Getenv("JWT_ACCESS_SECRET"))
}

func ParseRefreshJWT(t string) (JwtPayload, error) {
	return parseJwt(t, os.Getenv("JWT_REFRESH_SECRET"))
}
