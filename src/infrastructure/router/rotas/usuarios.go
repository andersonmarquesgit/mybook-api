package rotas

import (
	controllers "mybook-api/src/controllers"
	"net/http"
)

var rotasUsuarios = []Rota{
	{
		URI:                "/usuarios",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarUsuario,
		RequerAutenticacao: false,
	},
	{
		URI:                "/usuarios",
		Metodo:             http.MethodGet,
		Funcao:             controllers.ListarUsuarios,
		RequerAutenticacao: true,
	},
	{
		URI:                "/usuarios/{id}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarUsuario,
		RequerAutenticacao: true,
	},
	{
		URI:                "/usuarios/{id}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarUsuario,
		RequerAutenticacao: true,
	},
	{
		URI:                "/usuarios/{id}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeletarUsuario,
		RequerAutenticacao: true,
	},
	{
		URI:                "/usuarios/{id}/follow",
		Metodo:             http.MethodPost,
		Funcao:             controllers.FollowUser,
		RequerAutenticacao: true,
	},
	{
		URI:                "/usuarios/{id}/unfollow",
		Metodo:             http.MethodPost,
		Funcao:             controllers.UnfollowUser,
		RequerAutenticacao: true,
	},
	{
		URI:                "/usuarios/{id}/followers",
		Metodo:             http.MethodGet,
		Funcao:             controllers.Followers,
		RequerAutenticacao: true,
	},
	{
		URI:                "/usuarios/{id}/following",
		Metodo:             http.MethodGet,
		Funcao:             controllers.Following,
		RequerAutenticacao: true,
	},
}
