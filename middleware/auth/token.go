package auth

import (
    "time"
    "github.com/dgrijalva/jwt-go"
    "github.com/pkg/errors"
)


type Claims struct {
    Uid uint64
    jwt.StandardClaims
}

func GenAuthToken(uid uint64, expireDuration time.Duration, jwtSecret string) (string, error) {
    nowTime := time.Now()
    expireTime := nowTime.Add(expireDuration)

    claims := Claims{
        uid,
        jwt.StandardClaims{
            ExpiresAt: expireTime.Unix(),
            IssuedAt: nowTime.Unix(),
        },
    }

    tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return tokenClaims.SignedString([]byte(jwtSecret))
}

func ParseAuthToken(token string, secret string) (uint64, error) {
    claims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    if err != nil {
        return 0, err
    }

    if claims != nil {
        if c, ok := claims.Claims.(*Claims); ok && claims.Valid {
            return c.Uid, nil
        }
    }
    return 0, errors.New("invalid token")
}
