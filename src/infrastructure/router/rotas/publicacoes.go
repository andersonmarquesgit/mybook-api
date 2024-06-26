package rotas

import (
	controllers "mybook-api/src/controllers"
	"net/http"
)

var rotasPublicacoes = []Rota{
	{
		URI:                "/publications",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarPublicacoes,
		RequerAutenticacao: false,
	},
	{
		URI:                "/publications",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarPublicacoes,
		RequerAutenticacao: false,
	},
	{
		URI:                "/publications/{id}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarPublicacao,
		RequerAutenticacao: false,
	},
	{
		URI:                "/publicacoes/{id}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarPublicacao,
		RequerAutenticacao: false,
	},
	{
		URI:                "/publicacoes/{id}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeletarPublicacao,
		RequerAutenticacao: false,
	},
}
