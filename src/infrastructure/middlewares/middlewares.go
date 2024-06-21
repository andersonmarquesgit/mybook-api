package middlewares

import (
	"log"
	"mybook-api/src/infrastructure/autenticacao"
	"mybook-api/src/response"
	"net/http"
)

func Logger(nextFunction http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n %s %s %s", r.Method, r.RequestURI, r.Host)
		nextFunction(w, r)
	}
}

func Autenticar(nextFunction http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := autenticacao.ValidarToken(r); err != nil {
			response.Erro(w, http.StatusUnauthorized, err)
			return
		}
		nextFunction(w, r)
	}
}
