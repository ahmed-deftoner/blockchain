package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

type Block struct {
}

var blockchain []Block

func run() http.Handler {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

}
