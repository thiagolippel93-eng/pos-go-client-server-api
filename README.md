▶️ Como Executar
1. Clonar o projeto
git clone https://github.com/thiagolippel93-eng/pos-go-client-server-api
cd client-server-api

2. Inicializar módulo Go
go mod init client-server-api
go mod tidy

3. Executar o servidor
go run server.go

O servidor estará disponível em:

http://localhost:8080/cotacao
4. Executar o cliente

Em outro terminal:

go run client.go
