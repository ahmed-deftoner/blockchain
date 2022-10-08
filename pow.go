package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
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

type Message struct {
	data string
}

func run() error {
	mux := makeMuxRoutes()
	httpPort := os.Getenv("PORT")
	log.Printf("Listening on port : %s", httpPort)
	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func makeMuxRoutes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", handleGetBlockchain).Methods("GET")
	router.HandleFunc("/", handlePostBlock).Methods("POST")
	return router
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(blockchain, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handlePostBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m Message
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&m); err != nil {
		respndWithJson(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	mutex.Lock()
	newblock := generateBlock(blockchain[len(blockchain)-1], m.data)
	mutex.Unlock()
	if isBlockValid(blockchain[len(blockchain)-1], newblock) {
		blockchain = append(blockchain, newblock)
		spew.Dump(blockchain)
	}
	respndWithJson(w, r, http.StatusCreated, newblock)
}

func respndWithJson(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func calculateHash(b Block) string {
	return ""
}

func isBlockValid(oldBlock, newBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.Prevhash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

func generateBlock(oldblock Block, data string) Block {
	newblock := Block{}
	t := time.Now()
	newblock.Index = oldblock.Index + 1
	newblock.Timestamp = t.String()
	newblock.Data = data
	newblock.Prevhash = oldblock.Hash
	newblock.Difficulty = difficulty

	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newblock.Nonce = hex
		if !isValidHash(calculateHash(newblock), newblock.Difficulty) {
			fmt.Println(calculateHash(newblock), "keep mining")
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(calculateHash(newblock), "mining complete")
			newblock.Hash = calculateHash(newblock)
			break
		}
	}
	return newblock
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		t := time.Now()
		genesisBlock := Block{}
		genesisBlock = Block{0, t.String(), "", calculateHash(genesisBlock), "", "", difficulty}
		spew.Dump(genesisBlock)

		mutex.Lock()
		blockchain = append(blockchain, genesisBlock)
		mutex.Unlock()
	}()
	log.Fatal(run())
}
