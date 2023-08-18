package main

import (
	"context"
	"log"
	"net/http"
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

}
