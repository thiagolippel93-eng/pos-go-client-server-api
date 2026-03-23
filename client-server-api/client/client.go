package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	serverURL  = "http://localhost:8080/cotacao"
	outputFile = "cotacao.txt"
)

// Estrutura esperada da resposta do servidor
type ExchangeResponse struct {
	Bid string `json:"bid"`
}

func main() {
	// Timeout de 300ms para requisição ao servidor
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL, nil)
	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Erro ao chamar servidor: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Resposta inválida do servidor: %s", string(body))
	}

	var exchange ExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&exchange); err != nil {
		log.Fatalf("Erro ao decodificar resposta: %v", err)
	}

	// Salva em arquivo
	content := fmt.Sprintf("Dólar: %s", exchange.Bid)
	if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
		log.Fatalf("Erro ao salvar arquivo: %v", err)
	}

	fmt.Println("Cotação recebida e salva em", outputFile)
}
