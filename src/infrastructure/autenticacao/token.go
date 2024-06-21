package autenticacao

import (
	"mybook-api/src/infrastructure/config"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func CriarToken(usuarioID string) (string, error) {
	permissoes := jwt.MapClaims{}
	permissoes["authorized"] = true
	permissoes["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissoes["usuarioId"] = usuarioID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissoes)
	return token.SignedString([]byte(config.SecretKey))
}
