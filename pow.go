package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
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

var mutex = &sync.Mutex{}

func run() error {

}

func makeMuxRoutes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", handleGetBlockchain).Methods("GET")
	router.HandleFunc("/", handlePostBlock).Methods("POST")
	return router
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
		mutex.Lock()
		blockchain = append(blockchain, genesisBlock)
		mutex.Unlock()
	}()
	log.Fatal(run())
}
