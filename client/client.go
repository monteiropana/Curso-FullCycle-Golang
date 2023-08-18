package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	//Criando e Preparando a requisicao, usando o contexto TimeOut
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Microsecond)
	defer cancel()
	requuest, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080", nil)
	if err != nil {
		log.Print("Erro ao fazer a requisicao do client", err)
	}

	//Executando a requisicao
	resposta, err := http.DefaultClient.Do(requuest)
	if err != nil {
		log.Print("Erro ao executar a requisicao do client", err)
	}
	defer resposta.Body.Close()

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Print("Erro ao criar o arquivo.txt")
	}
	defer file.Close()

	respBody, err := io.ReadAll(resposta.Body)
	if err != nil {
		log.Print("Erro ao ler os dados")
	}

	_, err = file.Write(respBody)
	if err != nil {
		log.Fatal(err)
	}
}
