package main

import (
	"fmt"
	"log"
	"mybook-api/src/router"
	"net/http"
)

func main() {
	fmt.Println("Executando API")

	r := router.Gerar()

	log.Fatal(http.ListenAndServe(":5000", r))
}
