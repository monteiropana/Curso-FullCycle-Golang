package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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

func main() {
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

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Print(string(responseData))

	var ObjectResponse Cotacao
	json.Unmarshal(responseData, &ObjectResponse)
	fmt.Println(ObjectResponse.Usdbrl.Bid)

	w.Write([]byte("{'bid':'" + ObjectResponse.Usdbrl.Bid + "'}"))
}
