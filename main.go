package main

import (
	"fmt"
	"log"
	"mybook-api/src/config"
	"mybook-api/src/router"
	"net/http"
)

func main() {
	config.Carregar()
	log.Printf("Escutando na porta %d", config.Port)

	r := router.Gerar()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
