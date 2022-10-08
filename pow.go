package main

import (
	"log"
	"net/http"
	"time"

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

func calculateHash(b Block) string {
	return ""
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		time := time.Now()
		genesisBlock := Block{}
		genesisBlock = Block{0, time.String(), "", calculateHash(genesisBlock), "", "", difficulty}
		blockchain = append(blockchain, genesisBlock)
	}()
}
