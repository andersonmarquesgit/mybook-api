package router

import (
	"mybook-api/src/infrastructure/router/rotas"

	"github.com/gorilla/mux"
)

func Gerar() *mux.Router {
	r := mux.NewRouter()
	return rotas.Configurar(r)
}
