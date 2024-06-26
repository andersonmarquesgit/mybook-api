package autenticacao

import (
	"errors"
	"fmt"
	"log"
	"mybook-api/src/infrastructure/config"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type Token struct {
	AccessToken string `json:"access_token,omitempty"`
}

func CriarToken(usuarioID string) (Token, error) {
	permissoes := jwt.MapClaims{}
	permissoes["authorized"] = true
	permissoes["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissoes["usuarioId"] = usuarioID

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, permissoes)
	accessTokenString, err := accessToken.SignedString([]byte(config.SecretKey))
	token := Token{AccessToken: accessTokenString}
	return token, err
}

func ValidarToken(r *http.Request) error {
	tokenString := extrairToken(r)
	token, err := jwt.Parse(tokenString, retornarChaveDeVerificacao)

	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("Token inválido")
}

func extrairToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	//Bearer 123
	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	log.Println("Token inválido!")
	return ""
}

func ExtrairUsuarioID(r *http.Request) (string, error) {
	tokenString := extrairToken(r)
	token, err := jwt.Parse(tokenString, retornarChaveDeVerificacao)
	if err != nil {
		return "", err
	}

	if permissoes, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		usuarioID := permissoes["usuarioId"]
		if usuarioID == nil {
			return "", nil
		}

		return fmt.Sprintf("%v", usuarioID), nil

	}

	return "", errors.New("Token Inválido")
}

func retornarChaveDeVerificacao(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Método de assinatura inesperado! %v", token.Header["alg"])
	}

	return config.SecretKey, nil
}
