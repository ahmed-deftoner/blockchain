package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Block struct {
	Pos       int
	Data      AlbumCheckout
	Timestamp string
	Hash      string
	PrevHash  string
}

type Blockchain struct {
	blocks []*Block
}

type AlbumCheckout struct {
	AlbumID      string `json:"album_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"`
}

type Album struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Artist      string `json:"artist"`
	ReleaseDate string `json:"release_date"`
	Genre       string `json:"genre"`
}

var blockchain *Blockchain

func CreateBlock(prevBlock *Block, checkoutItem AlbumCheckout) *Block {
	block := &Block{}
	block.Pos = prevBlock.Pos + 1
	block.Timestamp = time.Now().String()
	block.Data = checkoutItem
	block.PrevHash = prevBlock.Hash
	block.generateHash()
	return block
}

func (bc *Blockchain) AddBlock(data AlbumCheckout) {

	prevBlock := bc.blocks[len(bc.blocks)-1]

	block := CreateBlock(prevBlock, data)

	if validBlock(block, prevBlock) {
		bc.blocks = append(bc.blocks, block)
	}
}

func writeBlock(w http.ResponseWriter, r *http.Request) {
	var checkoutItem AlbumCheckout
	if err := json.NewDecoder(r.Body).Decode(&checkoutItem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not write Block: %v", err)
		w.Write([]byte("could not write block"))
		return
	}

	blockchain.AddBlock(checkoutItem)
	resp, err := json.MarshalIndent(checkoutItem, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not marshal payload: %v", err)
		w.Write([]byte("could not write block"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func newAlbum(w http.ResponseWriter, r *http.Request) {
	var album Album
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not create: %v", err)
		w.Write([]byte("could not create new Book"))
		return
	}

	h := md5.New()
	io.WriteString(h, album.Genre+album.Name+album.Artist)
	album.ID = fmt.Sprintf("%x", h.Sum(nil))

	resp, err := json.MarshalIndent(album, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not marshal payload: %v", err)
		w.Write([]byte("could not save book data"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newAlbum).Methods("POST")

	log.Println("listening on port 3000 ...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
