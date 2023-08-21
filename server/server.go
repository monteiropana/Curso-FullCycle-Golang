package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	Usdbrl Response `json: USDBRL"`
}

type Response struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type Database struct {
	DB *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{
		DB: db,
	}
}

var datab *Database

func main() {
	db, err := sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	datab = NewDatabase(db)
	http.HandleFunc("/cotacao", cotacao)
	http.ListenAndServe(":8080", nil)

}

func cotacao(w http.ResponseWriter, r *http.Request) {
	log.Println("Request iniciada")
	defer log.Println("Request finalizada")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Println("Erro ao realizar a requisicao no server", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro ao executar a requisicao no server", err)
	}
	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Print(string(responseData))

	var ObjectResponse Cotacao
	json.Unmarshal(responseData, &ObjectResponse)
	fmt.Println(ObjectResponse.Usdbrl.Bid)
	datab.SaveCotacao(ObjectResponse.Usdbrl.Bid)

	w.Write([]byte("{'bid':'" + ObjectResponse.Usdbrl.Bid + "'}"))
}

func (db *Database) SaveCotacao(bid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	conexao, err := db.DB.Conn(ctx)
	if err != nil {
		log.Println("Error creating connection dataBase")
	}

	var id = uuid.New().String()
	_, err = conexao.ExecContext(ctx, "INSERT INTO cotacao (id,bid) VALUES($1,$2)", id, bid)
	if err != nil {
		log.Println("Error inserting")
	}
	return nil
}
