package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

const difficulty = 1

type Block struct {
	Index      int
	Timestamp  string
	Data       string
	Hash       string
	Prevhash   string
	Nonce      string
	Difficulty int
}

var blockchain []Block

func run() http.Handler {

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}
