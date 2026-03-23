package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

const (
	apiURL        = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	dbFile        = "cotacoes.db"
	serverAddress = ":8080"
)

// Estrutura do retorno da API externa
type ExchangeAPIResponse struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

// Estrutura do retorno da nossa API
type ExchangeResponse struct {
	Bid string `json:"bid"`
}

func main() {
	// Inicializa banco SQLite (modernc.org/sqlite)
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatalf("Erro ao abrir banco SQLite: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		handleCotacao(w, r, db)
	})

	log.Printf("Servidor rodando em http://localhost%s/cotacao", serverAddress)
	if err := http.ListenAndServe(serverAddress, nil); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}

func handleCotacao(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Timeout para consulta na API externa (200ms)
	ctxAPI, cancelAPI := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancelAPI()

	req, err := http.NewRequestWithContext(ctxAPI, "GET", apiURL, nil)
	if err != nil {
		http.Error(w, "Erro ao criar requisição", http.StatusInternalServerError)
		log.Printf("Erro request: %v", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Erro ao consultar API externa", http.StatusGatewayTimeout)
		log.Printf("Erro chamada API externa: %v", err)
		return
	}
	defer resp.Body.Close()

	var apiResp ExchangeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		http.Error(w, "Erro ao decodificar resposta da API externa", http.StatusInternalServerError)
		log.Printf("Erro decode: %v", err)
		return
	}

	// Timeout para inserção no banco (10ms)
	ctxDB, cancelDB := context.WithTimeout(r.Context(), 10*time.Millisecond)
	defer cancelDB()

	_, err = db.ExecContext(ctxDB, "INSERT INTO cotacoes (bid) VALUES (?)", apiResp.USDBRL.Bid)
	if err != nil {
		log.Printf("Erro ao inserir no banco (timeout ou falha): %v", err)
	}

	// Retorna resposta JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ExchangeResponse{Bid: apiResp.USDBRL.Bid})
}
