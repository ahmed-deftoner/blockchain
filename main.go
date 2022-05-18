package main

import (
	"log"
	"net/http"

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

var Blockchain *Blockchain

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newBook).Methods("POST")

	log.Println("listening on port 3000 ...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
