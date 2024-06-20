package rotas

import (
	"mybook-api/src/controllers"
	"net/http"
)

var rotaLogin = Rota{
	URI:                "/login",
	Metodo:             http.MethodPost,
	Funcao:             controllers.Login,
	RequerAutenticacao: false,
}
