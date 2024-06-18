package main

import (
	"fmt"
	"log"
	"mybook-api/src/infrastructure/banco"
	"mybook-api/src/infrastructure/config"
	"mybook-api/src/infrastructure/router"
	"net/http"
)

func main() {
	config.Carregar()
	banco.ConectarMongoDB()
	defer banco.DesconectarMongoDB()

	log.Printf("Escutando na porta %d", config.Port)

	r := router.Gerar()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))

}
